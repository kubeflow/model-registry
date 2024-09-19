package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// INFERENCE SERVICE

// UpsertInferenceService creates a new inference service if the provided inference service's ID is nil,
// or updates an existing inference service if the ID is provided.
func (serv *ModelRegistryService) UpsertInferenceService(inferenceService *openapi.InferenceService) (*openapi.InferenceService, error) {
	if inferenceService == nil {
		return nil, fmt.Errorf("invalid inference service pointer, can't upsert nil: %w", api.ErrBadRequest)
	}
	var err error
	var existing *openapi.InferenceService
	var servingEnvironment *openapi.ServingEnvironment

	if inferenceService.Id == nil {
		// create
		glog.Info("Creating new InferenceService")
		servingEnvironment, err = serv.GetServingEnvironmentById(inferenceService.ServingEnvironmentId)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		glog.Infof("Updating InferenceService %s", *inferenceService.Id)

		existing, err = serv.GetInferenceServiceById(*inferenceService.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := serv.openapiConv.OverrideNotEditableForInferenceService(converter.NewOpenapiUpdateWrapper(existing, inferenceService))
		if err != nil {
			return nil, err
		}
		inferenceService = &withNotEditable

		servingEnvironment, err = serv.getServingEnvironmentByInferenceServiceId(*inferenceService.Id)
		if err != nil {
			return nil, err
		}
	}

	// validate RegisteredModelId is also valid
	if _, err := serv.GetRegisteredModelById(inferenceService.RegisteredModelId); err != nil {
		return nil, err
	}

	// if already existing assure the name is the same
	if existing != nil && inferenceService.Name == nil {
		// user did not provide it
		// need to set it to avoid mlmd error "context name should not be empty"
		inferenceService.Name = existing.Name
	}

	protoCtx, err := serv.mapper.MapFromInferenceService(inferenceService, *servingEnvironment.Id)
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

	inferenceServiceId := &protoCtxResp.ContextIds[0]
	if inferenceService.Id == nil {
		servingEnvironmentId, err := converter.StringToInt64(servingEnvironment.Id)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		_, err = serv.mlmdClient.PutParentContexts(context.Background(), &proto.PutParentContextsRequest{
			ParentContexts: []*proto.ParentContext{{
				ChildId:  inferenceServiceId,
				ParentId: servingEnvironmentId,
			}},
			TransactionOptions: &proto.TransactionOptions{},
		})
		if err != nil {
			return nil, err
		}
	}

	idAsString := converter.Int64ToString(inferenceServiceId)
	toReturn, err := serv.GetInferenceServiceById(*idAsString)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

// getServingEnvironmentByInferenceServiceId retrieves the serving environment associated with the specified inference service ID.
func (serv *ModelRegistryService) getServingEnvironmentByInferenceServiceId(id string) (*openapi.ServingEnvironment, error) {
	glog.Infof("Getting ServingEnvironment for InferenceService %s", id)

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
		return nil, fmt.Errorf("multiple ServingEnvironments found for InferenceService %s: %w", id, api.ErrNotFound)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no ServingEnvironments found for InferenceService %s: %w", id, api.ErrNotFound)
	}

	toReturn, err := serv.mapper.MapToServingEnvironment(getParentResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

// GetInferenceServiceById retrieves an inference service by its unique identifier (ID).
func (serv *ModelRegistryService) GetInferenceServiceById(id string) (*openapi.InferenceService, error) {
	glog.Infof("Getting InferenceService by id %s", id)

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
		return nil, fmt.Errorf("multiple InferenceServices found for id %s: %w", id, api.ErrNotFound)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no InferenceService found for id %s: %w", id, api.ErrNotFound)
	}

	toReturn, err := serv.mapper.MapToInferenceService(getByIdResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

// GetInferenceServiceByParams retrieves an inference service based on specified parameters, such as (name and serving environment ID), or external ID.
// If multiple or no serving environments are found, an error is returned accordingly.
func (serv *ModelRegistryService) GetInferenceServiceByParams(name *string, servingEnvironmentId *string, externalId *string) (*openapi.InferenceService, error) {
	filterQuery := ""
	if name != nil && servingEnvironmentId != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", converter.PrefixWhenOwned(servingEnvironmentId, *name))
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (name and servingEnvironmentId), or externalId: %w", api.ErrBadRequest)
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.InferenceServiceTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple inference services found for name=%v, servingEnvironmentId=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(servingEnvironmentId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(getByParamsResp.Contexts) == 0 {
		return nil, fmt.Errorf("no inference services found for name=%v, servingEnvironmentId=%v, externalId=%v: %w", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(servingEnvironmentId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	toReturn, err := serv.mapper.MapToInferenceService(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	return toReturn, nil
}

// GetInferenceServices retrieves a list of inference services based on the provided list options and optional serving environment ID and runtime.
func (serv *ModelRegistryService) GetInferenceServices(listOptions api.ListOptions, servingEnvironmentId *string, runtime *string) (*openapi.InferenceServiceList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	queries := []string{}
	if servingEnvironmentId != nil {
		queryParentCtxId := fmt.Sprintf("parent_contexts_a.id = %s", *servingEnvironmentId)
		queries = append(queries, queryParentCtxId)
	}

	if runtime != nil {
		queryRuntimeProp := fmt.Sprintf("properties.runtime.string_value = \"%s\"", *runtime)
		queries = append(queries, queryRuntimeProp)
	}

	query := strings.Join(queries, " and ")
	listOperationOptions.FilterQuery = &query

	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.InferenceServiceTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.InferenceService{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToInferenceService(c)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.InferenceServiceList{
		NextPageToken: apiutils.ZeroIfNil(contextsResp.NextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}
