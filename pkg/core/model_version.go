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

// MODEL VERSIONS

// UpsertModelVersion creates a new model version if the provided model version's ID is nil,
// or updates an existing model version if the ID is provided.
func (serv *ModelRegistryService) UpsertModelVersion(modelVersion *openapi.ModelVersion, registeredModelId *string) (*openapi.ModelVersion, error) {
	if modelVersion == nil {
		return nil, fmt.Errorf("invalid model version pointer, can't upsert nil: %w", api.ErrBadRequest)
	}
	var err error
	var existing *openapi.ModelVersion
	var registeredModel *openapi.RegisteredModel

	if modelVersion.Id == nil {
		// create
		glog.Info("Creating new model version")
		if registeredModelId == nil {
			return nil, fmt.Errorf("missing registered model id, cannot create model version without registered model: %w", api.ErrBadRequest)
		}
		registeredModel, err = serv.GetRegisteredModelById(*registeredModelId)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		glog.Infof("Updating model version %s", *modelVersion.Id)
		existing, err = serv.GetModelVersionById(*modelVersion.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := serv.openapiConv.OverrideNotEditableForModelVersion(converter.NewOpenapiUpdateWrapper(existing, modelVersion))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		modelVersion = &withNotEditable

		registeredModel, err = serv.getRegisteredModelByVersionId(*modelVersion.Id)
		if err != nil {
			return nil, err
		}
	}

	modelCtx, err := serv.mapper.MapFromModelVersion(modelVersion, *registeredModel.Id, registeredModel.Name)
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

	modelId := &modelCtxResp.ContextIds[0]
	if modelVersion.Id == nil {
		registeredModelId, err := converter.StringToInt64(registeredModel.Id)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}

		_, err = serv.mlmdClient.PutParentContexts(context.Background(), &proto.PutParentContextsRequest{
			ParentContexts: []*proto.ParentContext{{
				ChildId:  modelId,
				ParentId: registeredModelId,
			}},
			TransactionOptions: &proto.TransactionOptions{},
		})
		if err != nil {
			return nil, err
		}
	}

	idAsString := converter.Int64ToString(modelId)
	model, err := serv.GetModelVersionById(*idAsString)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// GetModelVersionById retrieves a model version by its unique identifier (ID).
func (serv *ModelRegistryService) GetModelVersionById(id string) (*openapi.ModelVersion, error) {
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
		return nil, fmt.Errorf("multiple model versions found for id %s: %w", id, api.ErrNotFound)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no model version found for id %s: %w", id, api.ErrNotFound)
	}

	modelVer, err := serv.mapper.MapToModelVersion(getByIdResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return modelVer, nil
}

// GetModelVersionByInferenceService retrieves the model version associated with the specified inference service ID.
func (serv *ModelRegistryService) GetModelVersionByInferenceService(inferenceServiceId string) (*openapi.ModelVersion, error) {
	is, err := serv.GetInferenceServiceById(inferenceServiceId)
	if err != nil {
		return nil, err
	}
	if is.ModelVersionId != nil {
		return serv.GetModelVersionById(*is.ModelVersionId)
	}
	// modelVersionId: ID of the ModelVersion to serve. If it's unspecified, then the latest ModelVersion by creation order will be served.
	orderByCreateTime := "CREATE_TIME"
	sortOrderDesc := "DESC"
	versions, err := serv.GetModelVersions(api.ListOptions{OrderBy: &orderByCreateTime, SortOrder: &sortOrderDesc}, &is.RegisteredModelId)
	if err != nil {
		return nil, err
	}
	if len(versions.Items) == 0 {
		return nil, fmt.Errorf("no model versions found for id %s: %w", is.RegisteredModelId, api.ErrNotFound)
	}
	return &versions.Items[0], nil
}

// getModelVersionByArtifactId retrieves the model version associated with the specified model artifact ID.
func (serv *ModelRegistryService) getModelVersionByArtifactId(id string) (*openapi.ModelVersion, error) {
	glog.Infof("Getting model version for model artifact %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	getParentResp, err := serv.mlmdClient.GetContextsByArtifact(context.Background(), &proto.GetContextsByArtifactRequest{
		ArtifactId: idAsInt,
	})
	if err != nil {
		return nil, err
	}

	if len(getParentResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple model versions found for artifact %s: %w", id, api.ErrNotFound)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no model version found for artifact %s: %w", id, api.ErrNotFound)
	}

	modelVersion, err := serv.mapper.MapToModelVersion(getParentResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return modelVersion, nil
}

// GetModelVersionByParams retrieves a model version based on specified parameters, such as (version name and registered model ID), or external ID.
// If multiple or no model versions are found, an error is returned.
func (serv *ModelRegistryService) GetModelVersionByParams(versionName *string, registeredModelId *string, externalId *string) (*openapi.ModelVersion, error) {
	filterQuery := ""
	if versionName != nil && registeredModelId != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", converter.PrefixWhenOwned(registeredModelId, *versionName))
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (versionName and registeredModelId), or externalId: %w", api.ErrBadRequest)
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.ModelVersionTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple model versions found for versionName=%v, registeredModelId=%v, externalId=%v: %w", apiutils.ZeroIfNil(versionName), apiutils.ZeroIfNil(registeredModelId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(getByParamsResp.Contexts) == 0 {
		return nil, fmt.Errorf("no model versions found for versionName=%v, registeredModelId=%v, externalId=%v: %w", apiutils.ZeroIfNil(versionName), apiutils.ZeroIfNil(registeredModelId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	modelVer, err := serv.mapper.MapToModelVersion(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	return modelVer, nil
}

// GetModelVersions retrieves a list of model versions based on the provided list options and optional registered model ID.
func (serv *ModelRegistryService) GetModelVersions(listOptions api.ListOptions, registeredModelId *string) (*openapi.ModelVersionList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	if registeredModelId != nil {
		queryParentCtxId := fmt.Sprintf("parent_contexts_a.id = %s", *registeredModelId)
		listOperationOptions.FilterQuery = &queryParentCtxId
	}

	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &serv.nameConfig.ModelVersionTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.ModelVersion{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToModelVersion(c)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ModelVersionList{
		NextPageToken: apiutils.ZeroIfNil(contextsResp.NextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}
