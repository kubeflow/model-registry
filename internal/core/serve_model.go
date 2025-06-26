package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"gorm.io/gorm"
)

func (b *ModelRegistryService) UpsertServeModel(serveModel *openapi.ServeModel, inferenceServiceId *string) (*openapi.ServeModel, error) {
	if serveModel == nil {
		return nil, fmt.Errorf("invalid serve model pointer, cannot be nil: %w", api.ErrBadRequest)
	}

	if serveModel.Id != nil {
		existing, err := b.GetServeModelById(*serveModel.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := b.mapper.OverrideNotEditableForServeModel(converter.NewOpenapiUpdateWrapper(existing, serveModel))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		serveModel = &withNotEditable
	}

	srvModel, err := b.mapper.MapFromServeModel(serveModel)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	name := ""

	if srvModel.GetAttributes().Name != nil {
		name = *srvModel.GetAttributes().Name
	}

	prefixedName := converter.PrefixWhenOwned(inferenceServiceId, name)
	srvModel.GetAttributes().Name = &prefixedName

	if inferenceServiceId == nil && srvModel.GetID() == nil {
		return nil, fmt.Errorf("missing inferenceServiceId, cannot create ServeModel without parent resource InferenceService: %w", api.ErrBadRequest)
	}

	var inferenceServiceID *int32

	if inferenceServiceId != nil {
		convertedId, err := strconv.ParseInt(*inferenceServiceId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		id := int32(convertedId)

		inferenceServiceID = &id

		_, err = b.inferenceServiceRepository.GetByID(*inferenceServiceID)
		if err != nil {
			return nil, fmt.Errorf("no InferenceService found for id %d: %w", *inferenceServiceID, api.ErrNotFound)
		}
	}

	savedSrvModel, err := b.serveModelRepository.Save(srvModel, inferenceServiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("serve model with name %s already exists: %w", *serveModel.Name, api.ErrConflict)
		}

		return nil, err
	}

	toReturn, err := b.mapper.MapToServeModel(savedSrvModel)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetServeModelById(id string) (*openapi.ServeModel, error) {
	glog.Infof("Getting ServeModel by id %s", id)

	convertedId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	serveModel, err := b.serveModelRepository.GetByID(int32(convertedId))
	if err != nil {
		return nil, fmt.Errorf("no serve model found for id %s: %w", id, api.ErrNotFound)
	}

	toReturn, err := b.mapper.MapToServeModel(serveModel)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

func (b *ModelRegistryService) GetServeModels(listOptions api.ListOptions, inferenceServiceId *string) (*openapi.ServeModelList, error) {
	var inferenceServiceID *int32

	if inferenceServiceId != nil {
		convertedId, err := strconv.ParseInt(*inferenceServiceId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		id := int32(convertedId)

		inferenceServiceID = &id
	}

	serveModels, err := b.serveModelRepository.List(models.ServeModelListOptions{
		Pagination: models.Pagination{
			PageSize:      listOptions.PageSize,
			OrderBy:       listOptions.OrderBy,
			SortOrder:     listOptions.SortOrder,
			NextPageToken: listOptions.NextPageToken,
		},
		InferenceServiceID: inferenceServiceID,
	})
	if err != nil {
		return nil, err
	}

	serveModelList := &openapi.ServeModelList{
		Items: []openapi.ServeModel{},
	}

	for _, serveModel := range serveModels.Items {
		serveModel, err := b.mapper.MapToServeModel(serveModel)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		serveModelList.Items = append(serveModelList.Items, *serveModel)
	}

	serveModelList.NextPageToken = serveModels.NextPageToken
	serveModelList.PageSize = serveModels.PageSize
	serveModelList.Size = int32(serveModels.Size)

	return serveModelList, nil
}
