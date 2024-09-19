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

// SERVE MODEL

// UpsertServeModel creates a new serve model if the provided serve model's ID is nil,
// or updates an existing serve model if the ID is provided.
func (serv *ModelRegistryService) UpsertServeModel(serveModel *openapi.ServeModel, inferenceServiceId *string) (*openapi.ServeModel, error) {
	if serveModel == nil {
		return nil, fmt.Errorf("invalid serve model pointer, can't upsert nil: %w", api.ErrBadRequest)
	}
	var err error
	var existing *openapi.ServeModel

	if serveModel.Id == nil {
		// create
		glog.Info("Creating new ServeModel")
		if inferenceServiceId == nil {
			return nil, fmt.Errorf("missing inferenceServiceId, cannot create ServeModel without parent resource InferenceService: %w", api.ErrBadRequest)
		}
		_, err = serv.GetInferenceServiceById(*inferenceServiceId)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		glog.Infof("Updating ServeModel %s", *serveModel.Id)

		existing, err = serv.GetServeModelById(*serveModel.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := serv.openapiConv.OverrideNotEditableForServeModel(converter.NewOpenapiUpdateWrapper(existing, serveModel))
		if err != nil {
			return nil, err
		}
		serveModel = &withNotEditable

		_, err = serv.getInferenceServiceByServeModel(*serveModel.Id)
		if err != nil {
			return nil, err
		}
	}
	_, err = serv.GetModelVersionById(serveModel.ModelVersionId)
	if err != nil {
		return nil, err
	}

	// if already existing assure the name is the same
	if existing != nil && serveModel.Name == nil {
		// user did not provide it
		// need to set it to avoid mlmd error "artifact name should not be empty"
		serveModel.Name = existing.Name
	}

	execution, err := serv.mapper.MapFromServeModel(serveModel, *inferenceServiceId)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	executionsResp, err := serv.mlmdClient.PutExecutions(context.Background(), &proto.PutExecutionsRequest{
		Executions: []*proto.Execution{execution},
	})
	if err != nil {
		return nil, err
	}

	// add explicit Association between ServeModel and InferenceService
	if inferenceServiceId != nil && serveModel.Id == nil {
		inferenceServiceId, err := converter.StringToInt64(inferenceServiceId)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		associations := []*proto.Association{}
		for _, a := range executionsResp.ExecutionIds {
			associations = append(associations, &proto.Association{
				ContextId:   inferenceServiceId,
				ExecutionId: &a,
			})
		}
		_, err = serv.mlmdClient.PutAttributionsAndAssociations(context.Background(), &proto.PutAttributionsAndAssociationsRequest{
			Attributions: make([]*proto.Attribution, 0),
			Associations: associations,
		})
		if err != nil {
			return nil, err
		}
	}

	idAsString := converter.Int64ToString(&executionsResp.ExecutionIds[0])
	mapped, err := serv.GetServeModelById(*idAsString)
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

// getInferenceServiceByServeModel retrieves the inference service associated with the specified serve model ID.
func (serv *ModelRegistryService) getInferenceServiceByServeModel(id string) (*openapi.InferenceService, error) {
	glog.Infof("Getting InferenceService for ServeModel %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	getParentResp, err := serv.mlmdClient.GetContextsByExecution(context.Background(), &proto.GetContextsByExecutionRequest{
		ExecutionId: idAsInt,
	})
	if err != nil {
		return nil, err
	}

	if len(getParentResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple InferenceService found for ServeModel %s: %w", id, api.ErrNotFound)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no InferenceService found for ServeModel %s: %w", id, api.ErrNotFound)
	}

	toReturn, err := serv.mapper.MapToInferenceService(getParentResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return toReturn, nil
}

// GetServeModelById retrieves a serve model by its unique identifier (ID).
func (serv *ModelRegistryService) GetServeModelById(id string) (*openapi.ServeModel, error) {
	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	executionsResp, err := serv.mlmdClient.GetExecutionsByID(context.Background(), &proto.GetExecutionsByIDRequest{
		ExecutionIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}

	if len(executionsResp.Executions) > 1 {
		return nil, fmt.Errorf("multiple ServeModels found for id %s: %w", id, api.ErrNotFound)
	}

	if len(executionsResp.Executions) == 0 {
		return nil, fmt.Errorf("no ServeModel found for id %s: %w", id, api.ErrNotFound)
	}

	result, err := serv.mapper.MapToServeModel(executionsResp.Executions[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return result, nil
}

// GetServeModels retrieves a list of serve models based on the provided list options and optional inference service ID.
func (serv *ModelRegistryService) GetServeModels(listOptions api.ListOptions, inferenceServiceId *string) (*openapi.ServeModelList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	var executions []*proto.Execution
	var nextPageToken *string
	if inferenceServiceId != nil {
		ctxId, err := converter.StringToInt64(inferenceServiceId)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		executionsResp, err := serv.mlmdClient.GetExecutionsByContext(context.Background(), &proto.GetExecutionsByContextRequest{
			ContextId: ctxId,
			Options:   listOperationOptions,
		})
		if err != nil {
			return nil, err
		}
		executions = executionsResp.Executions
		nextPageToken = executionsResp.NextPageToken
	} else {
		executionsResp, err := serv.mlmdClient.GetExecutionsByType(context.Background(), &proto.GetExecutionsByTypeRequest{
			TypeName: &serv.nameConfig.ServeModelTypeName,
			Options:  listOperationOptions,
		})
		if err != nil {
			return nil, err
		}
		executions = executionsResp.Executions
		nextPageToken = executionsResp.NextPageToken
	}

	results := []openapi.ServeModel{}
	for _, a := range executions {
		mapped, err := serv.mapper.MapToServeModel(a)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ServeModelList{
		NextPageToken: apiutils.ZeroIfNil(nextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}
