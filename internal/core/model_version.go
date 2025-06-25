package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"gorm.io/gorm"
)

func (b *ModelRegistryService) UpsertModelVersion(modelVersion *openapi.ModelVersion, registeredModelId *string) (*openapi.ModelVersion, error) {
	if modelVersion == nil {
		return nil, fmt.Errorf("invalid model version pointer, cannot be nil: %w", api.ErrBadRequest)
	}

	if modelVersion.Id != nil {
		existing, err := b.GetModelVersionById(*modelVersion.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := b.mapper.OverrideNotEditableForModelVersion(converter.NewOpenapiUpdateWrapper(existing, modelVersion))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		modelVersion = &withNotEditable
	}

	if registeredModelId != nil {
		modelVersion.RegisteredModelId = *registeredModelId
	}

	model, err := b.mapper.MapFromModelVersion(modelVersion)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	modelVersion.Name = converter.PrefixWhenOwned(&modelVersion.RegisteredModelId, modelVersion.Name)

	savedModel, err := b.modelVersionRepository.Save(model)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("model version with name %s already exists: %w", modelVersion.Name, api.ErrConflict)
		}

		return nil, err
	}

	toReturn, err := b.mapper.MapToModelVersion(savedModel)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetModelVersionById(id string) (*openapi.ModelVersion, error) {
	glog.Infof("Getting ModelVersion by id %s", id)

	convertedId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	model, err := b.modelVersionRepository.GetByID(int32(convertedId))
	if err != nil {
		return nil, fmt.Errorf("no model version found for id %s: %w", id, api.ErrNotFound)
	}

	toReturn, err := b.mapper.MapToModelVersion(model)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetModelVersionByInferenceService(inferenceServiceId string) (*openapi.ModelVersion, error) {
	convertedId, err := strconv.ParseInt(inferenceServiceId, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	infSvc, err := b.inferenceServiceRepository.GetByID(int32(convertedId))
	if err != nil {
		return nil, fmt.Errorf("no inference service found for id %s: %w", inferenceServiceId, api.ErrNotFound)
	}

	infSvcProps := infSvc.GetProperties()

	if infSvcProps == nil {
		return nil, fmt.Errorf("no registered model found for inference service")
	}

	modelVersionID := int32(0)

	for _, prop := range *infSvcProps {
		if prop.Name == "model_version_id" {
			modelVersionID = *prop.IntValue
			break
		}
	}

	if modelVersionID != 0 {
		return b.GetModelVersionById(strconv.Itoa(int(modelVersionID)))
	}

	registeredModelID := ""

	for _, prop := range *infSvcProps {
		if prop.Name == "registered_model_id" {
			registeredModelID = strconv.Itoa(int(*prop.IntValue))
			break
		}
	}
	// modelVersionId: ID of the ModelVersion to serve. If it's unspecified, then the latest ModelVersion by creation order will be served.
	orderByCreateTime := "CREATE_TIME"
	sortOrderDesc := "DESC"
	versions, err := b.GetModelVersions(api.ListOptions{OrderBy: &orderByCreateTime, SortOrder: &sortOrderDesc}, &registeredModelID)
	if err != nil {
		return nil, err
	}

	if len(versions.Items) == 0 {
		return nil, fmt.Errorf("no model versions found for id %s: %w", inferenceServiceId, api.ErrNotFound)
	}

	return &versions.Items[0], nil
}

func (b *ModelRegistryService) GetModelVersionByParams(versionName *string, registeredModelId *string, externalId *string) (*openapi.ModelVersion, error) {
	var combinedName *string

	if versionName != nil && registeredModelId != nil {
		n := converter.PrefixWhenOwned(registeredModelId, *versionName)
		combinedName = &n
	} else if externalId == nil {
		return nil, fmt.Errorf("invalid parameters call, supply either (versionName and registeredModelId), or externalId: %w", api.ErrBadRequest)
	}

	versionsList, err := b.modelVersionRepository.List(models.ModelVersionListOptions{
		Name:       combinedName,
		ExternalID: externalId,
	})
	if err != nil {
		return nil, err
	}

	if len(versionsList.Items) > 1 {
		return nil, fmt.Errorf("multiple model versions found for versionName=%v, registeredModelId=%v, externalId=%v: %w", apiutils.ZeroIfNil(versionName), apiutils.ZeroIfNil(registeredModelId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(versionsList.Items) == 0 {
		return nil, fmt.Errorf("no model versions found for versionName=%v, registeredModelId=%v, externalId=%v: %w", apiutils.ZeroIfNil(versionName), apiutils.ZeroIfNil(registeredModelId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	toReturn, err := b.mapper.MapToModelVersion(versionsList.Items[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetModelVersions(listOptions api.ListOptions, registeredModelId *string) (*openapi.ModelVersionList, error) {
	var parentResourceID *int32

	if registeredModelId != nil {
		convertedId, err := strconv.ParseInt(*registeredModelId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		id := int32(convertedId)
		parentResourceID = &id
	}

	versionsList, err := b.modelVersionRepository.List(models.ModelVersionListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
		},
		ParentResourceID: parentResourceID,
	})
	if err != nil {
		return nil, err
	}

	modelVersionList := &openapi.ModelVersionList{
		Items: []openapi.ModelVersion{},
	}

	for _, model := range versionsList.Items {
		modelVersion, err := b.mapper.MapToModelVersion(model)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		modelVersionList.Items = append(modelVersionList.Items, *modelVersion)
	}

	modelVersionList.NextPageToken = versionsList.NextPageToken
	modelVersionList.PageSize = versionsList.PageSize
	modelVersionList.Size = int32(versionsList.Size)

	return modelVersionList, nil
}
