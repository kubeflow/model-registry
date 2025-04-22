package core

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// REGISTERED MODELS

// UpsertRegisteredModel creates a new registered model if the given registered model's ID is nil,
// or updates an existing registered model if the ID is provided.
func (serv *ModelRegistryService) UpsertRegisteredModel(registeredModel *openapi.RegisteredModel) (*openapi.RegisteredModel, error) {
	var err error
	var existing *openapi.RegisteredModel

	if registeredModel.Id == nil {
		glog.Info("Creating new registered model")
	} else {
		glog.Infof("Updating registered model %s", *registeredModel.Id)
		existing, err = serv.GetRegisteredModelById(*registeredModel.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := serv.openapiConv.OverrideNotEditableForRegisteredModel(converter.NewOpenapiUpdateWrapper(existing, registeredModel))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		registeredModel = &withNotEditable
	}

	modelCtx, err := serv.mapper.MapFromRegisteredModel(registeredModel)
	if err != nil {
		return nil, err
	}

	modelCtxResp, err := serv.mlmdClient.PutContexts(context.Background(), &proto.PutContextsRequest{
		Contexts: []*proto.Context{
			modelCtx,
		},
	})
	if err != nil {
		return nil, err
	}

	idAsString := converter.Int64ToString(&modelCtxResp.ContextIds[0])
	model, err := serv.GetRegisteredModelById(*idAsString)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// GetRegisteredModelById retrieves a registered model by its unique identifier (ID).
func (serv *ModelRegistryService) GetRegisteredModelById(id string) (*openapi.RegisteredModel, error) {
	glog.Infof("Getting registered model %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple registered models found for id %s: %w", id, api.ErrNotFound)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no registered model found for id %s: %w", id, api.ErrNotFound)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getByIdResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return regModel, nil
}

// GetRegisteredModelByInferenceService retrieves a registered model associated with the specified inference service ID.
func (serv *ModelRegistryService) GetRegisteredModelByInferenceService(inferenceServiceId string) (*openapi.RegisteredModel, error) {
	is, err := serv.GetInferenceServiceById(inferenceServiceId)
	if err != nil {
		return nil, err
	}
	return serv.GetRegisteredModelById(is.RegisteredModelId)
}

// getRegisteredModelByVersionId retrieves a registered model associated with the specified model version ID.
func (serv *ModelRegistryService) getRegisteredModelByVersionId(id string) (*openapi.RegisteredModel, error) {
	glog.Infof("Getting registered model for model version %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	getParentResp, err := serv.mlmdClient.GetParentContextsByContext(context.Background(), &proto.GetParentContextsByContextRequest{
		ContextId: idAsInt,
	})
	if err != nil {
		return nil, err
	}

	if len(getParentResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple registered models found for model version %s: %w", id, api.ErrNotFound)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no registered model found for model version %s: %w", id, api.ErrNotFound)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getParentResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return regModel, nil
}

// GetRegisteredModelByParams retrieves a registered model based on specified parameters, such as name or external ID.
// If multiple or no registered models are found, an error is returned accordingly.
func (serv *ModelRegistryService) GetRegisteredModelByParams(name *string, externalId *string) (*openapi.RegisteredModel, error) {
	glog.Infof("Getting registered model by params name=%v, externalId=%v", name, externalId)

	filterQuery := ""
	if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId: %w", api.ErrBadRequest)
	}
	glog.Info("FilterQuery ", filterQuery)

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.RegisteredModelTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple registered models found for name=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(getByParamsResp.Contexts) == 0 {
		return nil, fmt.Errorf("no registered models found for name=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	return regModel, nil
}

// GetRegisteredModels retrieves a list of registered models based on the provided list options.
func (serv *ModelRegistryService) GetRegisteredModels(listOptions api.ListOptions) (*openapi.RegisteredModelList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}
	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.RegisteredModelTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.RegisteredModel{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToRegisteredModel(c)
		if err != nil {
			return nil, err
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.RegisteredModelList{
		NextPageToken: apiutils.ZeroIfNil(contextsResp.NextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}
