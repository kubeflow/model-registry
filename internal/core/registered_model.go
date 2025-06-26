package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"gorm.io/gorm"
)

func (b *ModelRegistryService) UpsertRegisteredModel(registeredModel *openapi.RegisteredModel) (*openapi.RegisteredModel, error) {
	if registeredModel == nil {
		return nil, fmt.Errorf("invalid registered model pointer, cannot be nil: %w", api.ErrBadRequest)
	}

	if registeredModel.Id != nil {
		existing, err := b.GetRegisteredModelById(*registeredModel.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := b.mapper.OverrideNotEditableForRegisteredModel(converter.NewOpenapiUpdateWrapper(existing, registeredModel))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		registeredModel = &withNotEditable
	}

	model, err := b.mapper.MapFromRegisteredModel(registeredModel)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	savedModel, err := b.registeredModelRepository.Save(model)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("registered model with name %s already exists: %w", registeredModel.Name, api.ErrConflict)
		}

		return nil, err
	}

	return b.mapper.MapToRegisteredModel(savedModel)
}

func (b *ModelRegistryService) GetRegisteredModelById(id string) (*openapi.RegisteredModel, error) {
	convertedId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	model, err := b.registeredModelRepository.GetByID(int32(convertedId))
	if err != nil {
		return nil, fmt.Errorf("no registered model found for id %s: %w", id, api.ErrNotFound)
	}

	return b.mapper.MapToRegisteredModel(model)
}

func (b *ModelRegistryService) GetRegisteredModelByInferenceService(inferenceServiceId string) (*openapi.RegisteredModel, error) {
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
		return nil, fmt.Errorf("no registered model found for inference service: %w", api.ErrNotFound)
	}

	var registeredModelId *int32

	for _, prop := range *infSvcProps {
		if prop.Name == "registered_model_id" {
			registeredModelId = prop.IntValue
			break
		}
	}

	if registeredModelId == nil {
		return nil, fmt.Errorf("no registered model id found for inference service: %w", api.ErrNotFound)
	}

	model, err := b.registeredModelRepository.GetByID(*registeredModelId)
	if err != nil {
		return nil, err
	}

	return b.mapper.MapToRegisteredModel(model)
}

func (b *ModelRegistryService) GetRegisteredModelByParams(name *string, externalId *string) (*openapi.RegisteredModel, error) {
	if name == nil && externalId == nil {
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId: %w", api.ErrBadRequest)
	}

	modelsList, err := b.registeredModelRepository.List(models.RegisteredModelListOptions{
		Name:       name,
		ExternalID: externalId,
	})
	if err != nil {
		return nil, err
	}

	if len(modelsList.Items) == 0 {
		return nil, fmt.Errorf("no registered models found for name=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(modelsList.Items) > 1 {
		return nil, fmt.Errorf("multiple registered models found for name=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	registeredModel, err := b.mapper.MapToRegisteredModel(modelsList.Items[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return registeredModel, nil
}

func (b *ModelRegistryService) GetRegisteredModels(listOptions api.ListOptions) (*openapi.RegisteredModelList, error) {
	modelsList, err := b.registeredModelRepository.List(models.RegisteredModelListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
		},
	})
	if err != nil {
		return nil, err
	}

	registeredModelList := &openapi.RegisteredModelList{
		Items: []openapi.RegisteredModel{},
	}

	for _, model := range modelsList.Items {
		registeredModel, err := b.mapper.MapToRegisteredModel(model)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		registeredModelList.Items = append(registeredModelList.Items, *registeredModel)
	}

	registeredModelList.NextPageToken = modelsList.NextPageToken
	registeredModelList.PageSize = modelsList.PageSize
	registeredModelList.Size = int32(modelsList.Size)

	return registeredModelList, nil
}
