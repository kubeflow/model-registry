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

// UpsertExperiment creates a new experiment or updates an existing one
func (serv *ModelRegistryService) UpsertExperiment(experiment *openapi.Experiment) (*openapi.Experiment, error) {
	var err error
	var existing *openapi.Experiment

	if experiment.Id != nil {
		existing, err = serv.GetExperimentById(*experiment.Id)
		if err != nil {
			return nil, err
		}
		overridden, err := serv.openapiConv.OverrideNotEditableForExperiment(converter.NewOpenapiUpdateWrapper(existing, experiment))
		if err != nil {
			return nil, err
		}
		experiment = &overridden
	}

	experimentCtx, err := serv.mapper.MapFromExperiment(experiment)
	if err != nil {
		return nil, err
	}

	contexts, err := serv.mlmdClient.PutContexts(context.Background(), &proto.PutContextsRequest{
		Contexts: []*proto.Context{experimentCtx},
	})
	if err != nil {
		return nil, err
	}

	idAsString := converter.Int64ToString(&contexts.ContextIds[0])
	experiment, err = serv.GetExperimentById(*idAsString)
	if err != nil {
		return nil, err
	}

	return experiment, nil
}

// GetExperimentById retrieves an experiment by its ID
func (serv *ModelRegistryService) GetExperimentById(id string) (*openapi.Experiment, error) {
	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	getByIdReq := proto.GetContextsByIDRequest{
		ContextIds: []int64{*idAsInt},
	}

	resp, err := serv.mlmdClient.GetContextsByID(context.Background(), &getByIdReq)
	if err != nil {
		return nil, err
	}

	if len(resp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple experiments found for id %s: %w", id, api.ErrNotFound)
	}

	if len(resp.Contexts) == 0 {
		return nil, fmt.Errorf("no experiment found for id %s: %w", id, api.ErrNotFound)
	}

	experiment, err := serv.mapper.MapToExperiment(resp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	experiment.Id = &id

	return experiment, nil
}

// getExperimentByVersionId retrieves a registered model associated with the specified model version ID.
func (serv *ModelRegistryService) getExperimentByExperimentRunId(id string) (*openapi.Experiment, error) {
	glog.Infof("Getting experiment for experiment run %s", id)

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
		return nil, fmt.Errorf("multiple experiments found for experiment run %s: %w", id, api.ErrNotFound)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no experiment found for experiment run %s: %w", id, api.ErrNotFound)
	}

	experiment, err := serv.mapper.MapToExperiment(getParentResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return experiment, nil
}

// GetExperimentByParams finds an experiment by name and/or external ID
func (serv *ModelRegistryService) GetExperimentByParams(name *string, externalId *string) (*openapi.Experiment, error) {
	filterQuery := ""
	if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	}
	if externalId != nil {
		if filterQuery != "" {
			filterQuery += " AND "
		}
		filterQuery += fmt.Sprintf("external_id = \"%s\"", *externalId)
	}

	if filterQuery == "" {
		return nil, fmt.Errorf("at least one parameter (name or externalId) must be provided")
	}

	getByParamsReq := proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.ExperimentTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	}

	resp, err := serv.mlmdClient.GetContextsByType(context.Background(), &getByParamsReq)
	if err != nil {
		return nil, err
	}

	if len(resp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple experiments found for provided parameters")
	}

	if len(resp.Contexts) == 0 {
		return nil, api.ErrNotFound
	}

	result, err := serv.mapper.MapToExperiment(resp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetExperiments retrieves all experiments with pagination
func (serv *ModelRegistryService) GetExperiments(listOptions api.ListOptions) (*openapi.ExperimentList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	getExperimentsReq := proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.ExperimentTypeName,
		Options:  listOperationOptions,
	}

	resp, err := serv.mlmdClient.GetContextsByType(context.Background(), &getExperimentsReq)
	if err != nil {
		return nil, err
	}

	results := []openapi.Experiment{}
	for _, c := range resp.Contexts {
		mapped, err := serv.mapper.MapToExperiment(c)
		if err != nil {
			return nil, err
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ExperimentList{
		NextPageToken: apiutils.ZeroIfNil(resp.NextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}

	return &toReturn, nil
}
