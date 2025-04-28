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

// ARTIFACTS

// UpsertModelVersionArtifact creates a new artifact if the provided artifact's ID is nil, or updates an existing artifact if the
// ID is provided.
// Upon creation, new artifacts will be associated with their corresponding model version.
func (serv *ModelRegistryService) UpsertModelVersionArtifact(artifact *openapi.Artifact, modelVersionId string) (*openapi.Artifact, error) {
	if artifact == nil {
		return nil, fmt.Errorf("invalid artifact pointer, can't upsert nil: %w", api.ErrBadRequest)
	}
	art, err := serv.upsertArtifact(artifact, &modelVersionId)
	if err != nil {
		return nil, err
	}
	// upsertArtifact already validates modelVersion

	var id *string
	if art.ModelArtifact != nil {
		id = art.ModelArtifact.Id
	} else if art.DocArtifact != nil {
		id = art.DocArtifact.Id
	} else {
		return nil, fmt.Errorf("unexpected artifact type: %v", art)
	}

	mv, _ := serv.getModelVersionByArtifactId(*id)

	if mv == nil {
		// add explicit Attribution between Artifact and ModelVersion
		modelVersionId, err := converter.StringToInt64(&modelVersionId)
		if err != nil {
			// unreachable
			return nil, fmt.Errorf("%v", err)
		}
		artifactId, err := converter.StringToInt64(id)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		attributions := []*proto.Attribution{}
		attributions = append(attributions, &proto.Attribution{
			ContextId:  modelVersionId,
			ArtifactId: artifactId,
		})
		_, err = serv.mlmdClient.PutAttributionsAndAssociations(context.Background(), &proto.PutAttributionsAndAssociationsRequest{
			Attributions: attributions,
			Associations: make([]*proto.Association, 0),
		})
		if err != nil {
			return nil, err
		}
	} else if *mv.Id != modelVersionId {
		return nil, fmt.Errorf("artifact %s is already associated with a different model version %s: %w", *id, *mv.Id, api.ErrBadRequest)
	}
	return art, nil
}

func (serv *ModelRegistryService) upsertArtifact(artifact *openapi.Artifact, modelVersionId *string) (*openapi.Artifact, error) {
	if artifact == nil {
		return nil, fmt.Errorf("invalid artifact pointer, can't upsert nil: %w", api.ErrBadRequest)
	}
	if ma := artifact.ModelArtifact; ma != nil {
		if ma.Id == nil {
			glog.Info("Creating model artifact")
		} else {
			glog.Info("Updating model artifact")
			existing, err := serv.GetModelArtifactById(*ma.Id)
			if err != nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a model artifact: %w", *ma.Id, api.ErrBadRequest)
			}

			withNotEditable, err := serv.openapiConv.OverrideNotEditableForModelArtifact(converter.NewOpenapiUpdateWrapper(existing, ma))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			artifact.ModelArtifact = &withNotEditable
		}
	} else if da := artifact.DocArtifact; da != nil {
		if da.Id == nil {
			glog.Info("Creating doc artifact")
		} else {
			glog.Info("Updating doc artifact")
			existing, err := serv.GetArtifactById(*da.Id)
			if err != nil {
				return nil, err
			}
			if existing.DocArtifact == nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a doc artifact: %w", *da.Id, api.ErrBadRequest)
			}

			withNotEditable, err := serv.openapiConv.OverrideNotEditableForDocArtifact(converter.NewOpenapiUpdateWrapper(existing.DocArtifact, da))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			artifact.DocArtifact = &withNotEditable
		}
	} else {
		return nil, fmt.Errorf("invalid artifact type, must be either ModelArtifact or DocArtifact: %w", api.ErrBadRequest)
	}
	if modelVersionId != nil {
		if _, err := serv.GetModelVersionById(*modelVersionId); err != nil {
			return nil, fmt.Errorf("no model version found for id %s: %w", *modelVersionId, api.ErrNotFound)
		}
	}
	pa, err := serv.mapper.MapFromArtifact(artifact, modelVersionId)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	artifactsResp, err := serv.mlmdClient.PutArtifacts(context.Background(), &proto.PutArtifactsRequest{
		Artifacts: []*proto.Artifact{pa},
	})
	if err != nil {
		return nil, err
	}

	idAsString := converter.Int64ToString(&artifactsResp.ArtifactIds[0])
	return serv.GetArtifactById(*idAsString)
}

// UpsertArtifact creates a new artifact if the provided artifact's ID is nil, or updates an existing artifact if the
// ID is provided.
func (serv *ModelRegistryService) UpsertArtifact(artifact *openapi.Artifact) (*openapi.Artifact, error) {
	return serv.upsertArtifact(artifact, nil)
}

func (serv *ModelRegistryService) GetArtifactById(id string) (*openapi.Artifact, error) {
	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	artifactsResp, err := serv.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}
	if len(artifactsResp.Artifacts) > 1 {
		return nil, fmt.Errorf("multiple artifacts found for id %s: %w", id, api.ErrNotFound)
	}
	if len(artifactsResp.Artifacts) == 0 {
		return nil, fmt.Errorf("no artifact found for id %s: %w", id, api.ErrNotFound)
	}
	return serv.mapper.MapToArtifact(artifactsResp.Artifacts[0])
}

// GetArtifactByParams retrieves an artifact based on specified parameters, such as (artifact name and model version ID), or external ID.
// If multiple or no model artifacts are found, an error is returned.
func (serv *ModelRegistryService) GetArtifactByParams(artifactName *string, modelVersionId *string, externalId *string) (*openapi.Artifact, error) {
	var artifact0 *proto.Artifact

	filterQuery := ""
	if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else if artifactName != nil && modelVersionId != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", converter.PrefixWhenOwned(modelVersionId, *artifactName))
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (artifactName and modelVersionId), or externalId: %w", api.ErrBadRequest)
	}
	glog.Info("FilterQuery ", filterQuery)

	artifactsResponse, err := serv.mlmdClient.GetArtifacts(context.Background(), &proto.GetArtifactsRequest{
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(artifactsResponse.Artifacts) > 1 {
		return nil, fmt.Errorf("multiple model artifacts found for artifactName=%v, modelVersionId=%v, externalId=%v: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(artifactsResponse.Artifacts) == 0 {
		return nil, fmt.Errorf("no model artifacts found for artifactName=%v, modelVersionId=%v, externalId=%v: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	artifact0 = artifactsResponse.Artifacts[0]

	result, err := serv.mapper.MapToArtifact(artifact0)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return result, nil
}

// GetArtifacts retrieves a list of artifacts based on the provided list options and optional model version ID.
func (serv *ModelRegistryService) GetArtifacts(listOptions api.ListOptions, modelVersionId *string) (*openapi.ArtifactList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	var artifacts []*proto.Artifact
	var nextPageToken *string
	if modelVersionId != nil {
		ctxId, err := converter.StringToInt64(modelVersionId)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artifactsResp, err := serv.mlmdClient.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
			ContextId: ctxId,
			Options:   listOperationOptions,
		})
		if err != nil {
			return nil, err
		}
		artifacts = artifactsResp.Artifacts
		nextPageToken = artifactsResp.NextPageToken
	} else {
		artifactsResp, err := serv.mlmdClient.GetArtifacts(context.Background(), &proto.GetArtifactsRequest{
			Options: listOperationOptions,
		})
		if err != nil {
			return nil, err
		}
		artifacts = artifactsResp.Artifacts
		nextPageToken = artifactsResp.NextPageToken
	}

	results := []openapi.Artifact{}
	for _, a := range artifacts {
		mapped, err := serv.mapper.MapToArtifact(a)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ArtifactList{
		NextPageToken: apiutils.ZeroIfNil(nextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// MODEL ARTIFACTS

// UpsertModelArtifact creates a new model artifact if the provided model artifact's ID is nil,
// or updates an existing model artifact if the ID is provided.
func (serv *ModelRegistryService) UpsertModelArtifact(modelArtifact *openapi.ModelArtifact) (*openapi.ModelArtifact, error) {
	if modelArtifact == nil {
		return nil, fmt.Errorf("invalid artifact pointer, can't upsert nil: %w", api.ErrBadRequest)
	}
	art, err := serv.UpsertArtifact(&openapi.Artifact{
		ModelArtifact: modelArtifact,
	})
	if err != nil {
		return nil, err
	}
	return art.ModelArtifact, err
}

// GetModelArtifactById retrieves a model artifact by its unique identifier (ID).
func (serv *ModelRegistryService) GetModelArtifactById(id string) (*openapi.ModelArtifact, error) {
	art, err := serv.GetArtifactById(id)
	if err != nil {
		return nil, err
	}
	ma := art.ModelArtifact
	if ma == nil {
		return nil, fmt.Errorf("artifact with id %s is not a model artifact: %w", id, api.ErrNotFound)
	}
	return ma, err
}

// GetModelArtifactByInferenceService retrieves the model artifact associated with the specified inference service ID.
func (serv *ModelRegistryService) GetModelArtifactByInferenceService(inferenceServiceId string) (*openapi.ModelArtifact, error) {
	mv, err := serv.GetModelVersionByInferenceService(inferenceServiceId)
	if err != nil {
		return nil, err
	}

	artifactList, err := serv.GetModelArtifacts(api.ListOptions{}, mv.Id)
	if err != nil {
		return nil, err
	}

	if artifactList.Size == 0 {
		return nil, fmt.Errorf("no artifacts found for model version %s: %w", *mv.Id, api.ErrNotFound)
	}

	return &artifactList.Items[0], nil
}

// GetModelArtifactByParams retrieves a model artifact based on specified parameters, such as (artifact name and model version ID), or external ID.
// If multiple or no model artifacts are found, an error is returned.
func (serv *ModelRegistryService) GetModelArtifactByParams(artifactName *string, modelVersionId *string, externalId *string) (*openapi.ModelArtifact, error) {
	var artifact0 *proto.Artifact

	filterQuery := ""
	if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else if artifactName != nil && modelVersionId != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", converter.PrefixWhenOwned(modelVersionId, *artifactName))
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (artifactName and modelVersionId), or externalId: %w", api.ErrBadRequest)
	}
	glog.Info("FilterQuery ", filterQuery)

	artifactsResponse, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
		TypeName: &serv.nameConfig.ModelArtifactTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(artifactsResponse.Artifacts) > 1 {
		return nil, fmt.Errorf("multiple model artifacts found for artifactName=%v, modelVersionId=%v, externalId=%v: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(artifactsResponse.Artifacts) == 0 {
		return nil, fmt.Errorf("no model artifacts found for artifactName=%v, modelVersionId=%v, externalId=%v: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	artifact0 = artifactsResponse.Artifacts[0]

	result, err := serv.mapper.MapToModelArtifact(artifact0)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return result, nil
}

// GetModelArtifacts retrieves a list of model artifacts based on the provided list options and optional model version ID.
func (serv *ModelRegistryService) GetModelArtifacts(listOptions api.ListOptions, modelVersionId *string) (*openapi.ModelArtifactList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	var artifacts []*proto.Artifact
	var nextPageToken *string
	if modelVersionId != nil {
		ctxId, err := converter.StringToInt64(modelVersionId)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		typeQuery := fmt.Sprintf("type = '%v'", serv.nameConfig.ModelArtifactTypeName)
		listOperationOptions.FilterQuery = &typeQuery
		artifactsResp, err := serv.mlmdClient.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
			ContextId: ctxId,
			Options:   listOperationOptions,
		})
		if err != nil {
			return nil, err
		}
		artifacts = artifactsResp.Artifacts
		nextPageToken = artifactsResp.NextPageToken
	} else {
		artifactsResp, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
			TypeName: &serv.nameConfig.ModelArtifactTypeName,
			Options:  listOperationOptions,
		})
		if err != nil {
			return nil, err
		}
		artifacts = artifactsResp.Artifacts
		nextPageToken = artifactsResp.NextPageToken
	}

	results := []openapi.ModelArtifact{}
	for _, a := range artifacts {
		mapped, err := serv.mapper.MapToModelArtifact(a)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ModelArtifactList{
		NextPageToken: apiutils.ZeroIfNil(nextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}
