package core

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
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

		withNotEditable, err := b.mapper.UpdateExistingServeModel(converter.NewOpenapiUpdateWrapper(existing, serveModel))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		serveModel = &withNotEditable
	}

	srvModel, err := b.mapper.MapFromServeModel(serveModel, inferenceServiceId)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	if inferenceServiceId == nil && srvModel.GetID() == nil {
		return nil, fmt.Errorf("missing inferenceServiceId, cannot create ServeModel without parent resource InferenceService: %w", api.ErrBadRequest)
	}

	var inferenceServiceID *int32

	if inferenceServiceId != nil {
		var err error
		inferenceServiceID, err = apiutils.ValidateIDAsInt32Ptr(inferenceServiceId, "inference service")
		if err != nil {
			return nil, err
		}

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

	convertedId, err := apiutils.ValidateIDAsInt32(id, "serve model")
	if err != nil {
		return nil, err
	}
	serveModel, err := b.serveModelRepository.GetByID(convertedId)
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
		var err error
		inferenceServiceID, err = apiutils.ValidateIDAsInt32Ptr(inferenceServiceId, "inference service")
		if err != nil {
			return nil, err
		}
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
