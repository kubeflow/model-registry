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

// SERVING ENVIRONMENT

// UpsertServingEnvironment creates a new serving environment if the provided serving environment's ID is nil,
// or updates an existing serving environment if the ID is provided.
func (serv *ModelRegistryService) UpsertServingEnvironment(servingEnvironment *openapi.ServingEnvironment) (*openapi.ServingEnvironment, error) {
	var err error
	var existing *openapi.ServingEnvironment

	if servingEnvironment.Id == nil {
		glog.Info("Creating new serving environment")
	} else {
		glog.Infof("Updating serving environment %s", *servingEnvironment.Id)
		existing, err = serv.GetServingEnvironmentById(*servingEnvironment.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := serv.openapiConv.OverrideNotEditableForServingEnvironment(converter.NewOpenapiUpdateWrapper(existing, servingEnvironment))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		servingEnvironment = &withNotEditable
	}

	protoCtx, err := serv.mapper.MapFromServingEnvironment(servingEnvironment)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	protoCtxResp, err := serv.mlmdClient.PutContexts(context.Background(), &proto.PutContextsRequest{
		Contexts: []*proto.Context{
			protoCtx,
		},
	})
	if err != nil {
		return nil, err
	}

	idAsString := converter.Int64ToString(&protoCtxResp.ContextIds[0])
	openapiModel, err := serv.GetServingEnvironmentById(*idAsString)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return openapiModel, nil
}

// GetServingEnvironmentById retrieves a serving environment by its unique identifier (ID).
func (serv *ModelRegistryService) GetServingEnvironmentById(id string) (*openapi.ServingEnvironment, error) {
	glog.Infof("Getting serving environment %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*idAsInt},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple serving environments found for id %s: %w", id, api.ErrNotFound)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no serving environment found for id %s: %w", id, api.ErrNotFound)
	}

	openapiModel, err := serv.mapper.MapToServingEnvironment(getByIdResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return openapiModel, nil
}

// GetServingEnvironmentByParams retrieves a serving environment based on specified parameters, such as name or external ID.
// If multiple or no serving environments are found, an error is returned accordingly.
func (serv *ModelRegistryService) GetServingEnvironmentByParams(name *string, externalId *string) (*openapi.ServingEnvironment, error) {
	glog.Infof("Getting serving environment by params name=%v, externalId=%v", name, externalId)

	filterQuery := ""
	if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId: %w", api.ErrBadRequest)
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.ServingEnvironmentTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple serving environments found for name=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(getByParamsResp.Contexts) == 0 {
		return nil, fmt.Errorf("no serving environments found for name=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	openapiModel, err := serv.mapper.MapToServingEnvironment(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	return openapiModel, nil
}

// GetServingEnvironments retrieves a list of serving environments based on the provided list options.
func (serv *ModelRegistryService) GetServingEnvironments(listOptions api.ListOptions) (*openapi.ServingEnvironmentList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.ServingEnvironmentTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.ServingEnvironment{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToServingEnvironment(c)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ServingEnvironmentList{
		NextPageToken: apiutils.ZeroIfNil(contextsResp.NextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}
