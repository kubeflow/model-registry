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

// UpsertExperimentRun creates a new experiment run or updates an existing one
func (serv *ModelRegistryService) UpsertExperimentRun(experimentRun *openapi.ExperimentRun, experimentId *string) (*openapi.ExperimentRun, error) {
	if experimentRun == nil {
		return nil, fmt.Errorf("invalid experiment run pointer, can't upsert nil: %w", api.ErrBadRequest)
	}
	var err error
	var existing *openapi.ExperimentRun
	var experiment *openapi.Experiment

	if experimentRun.Id == nil {
		// create
		glog.Info("Creating new experiment run")
		if experimentId == nil {
			return nil, fmt.Errorf("missing experiment id, cannot create experiment run without experiment: %w", api.ErrBadRequest)
		}
		experiment, err = serv.GetExperimentById(*experimentId)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		glog.Infof("Updating experiment run %s", *experimentRun.Id)
		existing, err = serv.GetExperimentRunById(*experimentRun.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := serv.openapiConv.OverrideNotEditableForExperimentRun(converter.NewOpenapiUpdateWrapper(existing, experimentRun))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		experimentRun = &withNotEditable

		experiment, err = serv.getExperimentByExperimentRunId(*experimentRun.Id)
		if err != nil {
			return nil, err
		}
	}

	modelCtx, err := serv.mapper.MapFromExperimentRun(experimentRun, experiment.Id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	modelCtxResp, err := serv.mlmdClient.PutContexts(context.Background(), &proto.PutContextsRequest{
		Contexts: []*proto.Context{
			modelCtx,
		},
	})
	if err != nil {
		return nil, err
	}

	experimentRunId := &modelCtxResp.ContextIds[0]
	if experimentRun.Id == nil {
		experimentId, err := converter.StringToInt64(experiment.Id)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		_, err = serv.mlmdClient.PutParentContexts(context.Background(), &proto.PutParentContextsRequest{
			ParentContexts: []*proto.ParentContext{{
				ChildId:  experimentRunId,
				ParentId: experimentId,
			}},
			TransactionOptions: &proto.TransactionOptions{},
		})
		if err != nil {
			return nil, err
		}
	}

	idAsString := converter.Int64ToString(experimentRunId)
	experimentRun, err = serv.GetExperimentRunById(*idAsString)
	if err != nil {
		return nil, err
	}

	return experimentRun, nil
}

// GetExperimentRunById retrieves an experiment run by its ID
func (serv *ModelRegistryService) GetExperimentRunById(id string) (*openapi.ExperimentRun, error) {
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
		return nil, fmt.Errorf("multiple experiment runs found for id %s: %w", id, api.ErrNotFound)
	}

	if len(resp.Contexts) == 0 {
		return nil, fmt.Errorf("no experiment run found for id %s: %w", id, api.ErrNotFound)
	}

	experimentRun, err := serv.mapper.MapToExperimentRun(resp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return experimentRun, nil
}

// GetExperimentRunByParams finds an experiment run by name, experiment ID, and/or external ID
func (serv *ModelRegistryService) GetExperimentRunByParams(name *string, experimentId *string, externalId *string) (*openapi.ExperimentRun, error) {
	filterQuery := ""
	if name != nil && experimentId != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", converter.PrefixWhenOwned(experimentId, *name))
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (name and experimentId), or externalId: %w", api.ErrBadRequest)
	}

	resp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.ExperimentRunTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple experiment runs found for provided parameters: %w", api.ErrNotFound)
	}

	if len(resp.Contexts) == 0 {
		return nil, fmt.Errorf("no experiment run found for provided parameters: %w", api.ErrNotFound)
	}

	experimentRun, err := serv.mapper.MapToExperimentRun(resp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return experimentRun, nil
}

// GetExperimentRuns retrieves all experiment runs with optional filtering by experiment ID
func (serv *ModelRegistryService) GetExperimentRuns(listOptions api.ListOptions, experimentId *string) (*openapi.ExperimentRunList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	// Add experiment ID filter if provided
	if experimentId != nil {
		filterQuery := fmt.Sprintf("parent_contexts_a.id = %s", *experimentId)
		if listOperationOptions.FilterQuery != nil {
			existingFilter := *listOperationOptions.FilterQuery
			filterQuery = fmt.Sprintf("(%s) AND (%s)", existingFilter, filterQuery)
		}
		listOperationOptions.FilterQuery = &filterQuery
	}

	resp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.ExperimentRunTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.ExperimentRun{}
	for _, c := range resp.Contexts {
		mapped, err := serv.mapper.MapToExperimentRun(c)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ExperimentRunList{
		NextPageToken: apiutils.ZeroIfNil(resp.NextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}

	return &toReturn, nil
}

// UpsertExperimentRunArtifact creates or updates an artifact associated with an experiment run
func (serv *ModelRegistryService) UpsertExperimentRunArtifact(artifact *openapi.Artifact, experimentRunId string) (*openapi.Artifact, error) {
	return serv.upsertArtifactWithAssociation(artifact, experimentRunId, "experiment run")
}

// GetExperimentRunArtifacts retrieves all artifacts associated with an experiment run
func (serv *ModelRegistryService) GetExperimentRunArtifacts(listOptions api.ListOptions, experimentRunId *string) (*openapi.ArtifactList, error) {
	return serv.GetArtifacts(listOptions, experimentRunId)
}
