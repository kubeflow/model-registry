package core

import (
	"context"
	"fmt"

	"github.com/kubeflow/model-registry/internal/defaults"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// ARTIFACTS

// upsertArtifactWithAssociation creates a new artifact if the provided artifact's ID is nil, or updates an existing artifact if the
// ID is provided.
// Upon creation, new artifacts will be associated with their corresponding parent context.
func (serv *ModelRegistryService) upsertArtifactWithAssociation(artifact *openapi.Artifact, parentContextID string, parentTypeName string) (*openapi.Artifact, error) {
	if artifact == nil {
		return nil, fmt.Errorf("invalid artifact pointer, can't upsert nil: %w", api.ErrBadRequest)
	}
	art, err := serv.upsertArtifactWithParentContext(artifact, &parentContextID, parentTypeName)
	if err != nil {
		return nil, err
	}
	// upsertArtifactWithParentContext already validates parentContextID

	var id *string
	if art.ModelArtifact != nil {
		id = art.ModelArtifact.Id
	} else if art.DocArtifact != nil {
		id = art.DocArtifact.Id
	} else if art.DataSet != nil {
		id = art.DataSet.Id
	} else if art.Metric != nil {
		id = art.Metric.Id
	} else if art.Parameter != nil {
		id = art.Parameter.Id
	} else {
		return nil, fmt.Errorf("unexpected artifact type: %v", art)
	}

	// get all contexts that the artifact is attributed to of the given type
	contexts, err := serv.getContextsByArtifactId(*id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	contextId, err := converter.StringToInt64(&parentContextID)
	if err != nil {
		// unreachable
		return nil, fmt.Errorf("%v", err)
	}
	// find parent context in contexts
	var parentContext *proto.Context
	for _, c := range contexts {
		if c.Id == contextId {
			parentContext = c
			break
		}
	}

	if parentContext == nil {
		// add explicit Attribution between Artifact and Context
		artifactId, err := converter.StringToInt64(id)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		attributions := []*proto.Attribution{}
		attributions = append(attributions, &proto.Attribution{
			ContextId:  contextId,
			ArtifactId: artifactId,
		})
		_, err = serv.mlmdClient.PutAttributionsAndAssociations(context.Background(), &proto.PutAttributionsAndAssociationsRequest{
			Attributions: attributions,
			Associations: make([]*proto.Association, 0),
		})
		if err != nil {
			return nil, err
		}
	}
	return art, nil
}

func (serv *ModelRegistryService) upsertArtifactWithParentContext(artifact *openapi.Artifact, parentContextId *string, parentTypeName string) (*openapi.Artifact, error) {
	if artifact == nil {
		return nil, fmt.Errorf("invalid artifact pointer, can't upsert nil: %w", api.ErrBadRequest)
	}

	var existingArtifact *openapi.Artifact
	if ma := artifact.ModelArtifact; ma != nil {
		if ma.Id == nil {
			glog.Info("Creating model artifact")
			if ma.Name == nil {
				ma.Name = converter.GenerateNewName() // uuid name
			}
		} else {
			glog.Info("Updating model artifact")
			existing, err := serv.GetModelArtifactById(*ma.Id)
			if err != nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a model artifact: %w", *ma.Id, api.ErrBadRequest)
			}
			existingArtifact = &openapi.Artifact{ModelArtifact: existing}

			withNotEditable, err := serv.openapiConv.OverrideNotEditableForModelArtifact(converter.NewOpenapiUpdateWrapper(existing, ma))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			artifact.ModelArtifact = &withNotEditable
		}
	} else if da := artifact.DocArtifact; da != nil {
		if da.Id == nil {
			glog.Info("Creating doc artifact")
			if da.Name == nil {
				da.Name = converter.GenerateNewName() // uuid name
			}
		} else {
			glog.Info("Updating doc artifact")
			existing, err := serv.GetArtifactById(*da.Id)
			if err != nil {
				return nil, err
			}
			if existing.DocArtifact == nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a doc artifact: %w", *da.Id, api.ErrBadRequest)
			}
			existingArtifact = existing

			withNotEditable, err := serv.openapiConv.OverrideNotEditableForDocArtifact(converter.NewOpenapiUpdateWrapper(existing.DocArtifact, da))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			artifact.DocArtifact = &withNotEditable
		}
	} else if ds := artifact.DataSet; ds != nil {
		if ds.Id == nil {
			glog.Info("Creating dataset artifact")
			if ds.Name == nil {
				ds.Name = converter.GenerateNewName() // uuid name
			}
		} else {
			glog.Info("Updating dataset artifact")
			existing, err := serv.GetArtifactById(*ds.Id)
			if err != nil {
				return nil, err
			}
			if existing.DataSet == nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a dataset artifact: %w", *ds.Id, api.ErrBadRequest)
			}
			existingArtifact = existing

			withNotEditable, err := serv.openapiConv.OverrideNotEditableForDataSet(converter.NewOpenapiUpdateWrapper(existing.DataSet, ds))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			artifact.DataSet = &withNotEditable
		}
	} else if m := artifact.Metric; m != nil {
		if m.Id == nil {
			// parentContextId and name is required for creation
			if parentContextId == nil {
				return nil, fmt.Errorf("missing parent context id for metric artifact: %w", api.ErrBadRequest)
			}
			if m.Name == nil {
				return nil, fmt.Errorf("missing name for metric artifact: %w", api.ErrBadRequest)
			}
			// check if metric name already exists
			existing, err := serv.GetArtifactByParams(m.Name, parentContextId, nil)
			if api.IgnoreNotFound(err) != nil {
				return nil, fmt.Errorf("error getting metric by name: %w", err)
			}
			if existing != nil {
				glog.Info("Updating metric artifact")
				existingArtifact = existing

				withNotEditable, err := serv.openapiConv.OverrideNotEditableForMetric(converter.NewOpenapiUpdateWrapper(existing.Metric, m))
				if err != nil {
					return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
				}
				artifact.Metric = &withNotEditable
			} else {
				glog.Info("Creating metric artifact")
			}
		} else {
			glog.Info("Updating metric artifact")
			existing, err := serv.GetArtifactById(*m.Id)
			if err != nil {
				return nil, err
			}
			if existing.Metric == nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a metric artifact: %w", *m.Id, api.ErrBadRequest)
			}
			existingArtifact = existing

			withNotEditable, err := serv.openapiConv.OverrideNotEditableForMetric(converter.NewOpenapiUpdateWrapper(existing.Metric, m))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			artifact.Metric = &withNotEditable
		}
	} else if p := artifact.Parameter; p != nil {
		if p.Id == nil {
			// parentContextId is required for creation
			if parentContextId == nil {
				return nil, fmt.Errorf("missing parent context id for parameter artifact: %w", api.ErrBadRequest)
			}
			if p.Name == nil {
				return nil, fmt.Errorf("missing name for parameter artifact: %w", api.ErrBadRequest)
			}
			glog.Info("Creating parameter artifact")
		} else {
			glog.Info("Updating parameter artifact")
			existing, err := serv.GetArtifactById(*p.Id)
			if err != nil {
				return nil, err
			}
			if existing.Parameter == nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a parameter artifact: %w", *p.Id, api.ErrBadRequest)
			}
			existingArtifact = existing

			withNotEditable, err := serv.openapiConv.OverrideNotEditableForParameter(converter.NewOpenapiUpdateWrapper(existing.Parameter, p))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			artifact.Parameter = &withNotEditable
		}
	} else {
		return nil, fmt.Errorf("invalid artifact type, must be either ModelArtifact or DocArtifact: %w", api.ErrBadRequest)
	}

	if parentContextId != nil {
		if _, err := serv.GetContextByID(*parentContextId); err != nil {
			return nil, fmt.Errorf("no %s found for id %s: %w", parentTypeName, *parentContextId, api.ErrNotFound)
		}
	}

	if existingArtifact != nil {
		newArtifact, err := converter.UpdateExistingArtifact(serv.reconciler, converter.NewOpenapiUpdateWrapper(existingArtifact, artifact))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artifact = &newArtifact
	}

	pa, err := serv.mapper.MapFromArtifact(artifact, parentContextId)
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

// UpsertArtifact creates a new standalone artifact if the provided artifact's ID is nil, or updates an existing artifact if the
// ID is provided.
func (serv *ModelRegistryService) UpsertArtifact(artifact *openapi.Artifact) (*openapi.Artifact, error) {
	return serv.upsertArtifactWithParentContext(artifact, nil, "")
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

// GetArtifactByParams retrieves an artifact based on specified parameters, such as (artifact name and parent resource ID), or external ID.
// If multiple or no model artifacts are found, an error is returned.
func (serv *ModelRegistryService) GetArtifactByParams(artifactName *string, parentResourceId *string, externalId *string) (*openapi.Artifact, error) {
	var artifact0 *proto.Artifact

	filterQuery := ""
	if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else if artifactName != nil && parentResourceId != nil {
		// search for the artifact name or the artifact name with the parent resource id as a prefix
		// and the parent resource id is the same as parentResourceId
		filterQuery = fmt.Sprintf("(name = \"%s\" or name like \"%%:%s\") and contexts_a.id = %s",
			*artifactName, *artifactName, *parentResourceId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (artifactName and parentResourceId), or externalId: %w", api.ErrBadRequest)
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
		return nil, fmt.Errorf("multiple model artifacts found for artifactName=%v, modelVersionId=%v, externalId=%v: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(parentResourceId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(artifactsResponse.Artifacts) == 0 {
		return nil, fmt.Errorf("no model artifacts found for artifactName=%v, modelVersionId=%v, externalId=%v: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(parentResourceId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	artifact0 = artifactsResponse.Artifacts[0]

	result, err := serv.mapper.MapToArtifact(artifact0)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return result, nil
}

// GetArtifacts retrieves a list of artifacts based on the provided list options and optional parent context ID.
func (serv *ModelRegistryService) GetArtifacts(artifactType openapi.ArtifactTypeQueryParam, listOptions api.ListOptions, parentContextId *string) (*openapi.ArtifactList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	// handle artifactType filter
	if artifactType != "" {
		mlmdType, err := toMlmdArtifactType(artifactType)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		listOperationOptions.FilterQuery = apiutils.Of(fmt.Sprintf("type = '%v'", mlmdType))
	}

	var artifacts []*proto.Artifact
	var nextPageToken *string
	if parentContextId != nil {
		ctxId, err := converter.StringToInt64(parentContextId)
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

var artifactTypeNameMap = map[openapi.ArtifactTypeQueryParam]string{
	openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT:   defaults.ModelArtifactTypeName,
	openapi.ARTIFACTTYPEQUERYPARAM_DOC_ARTIFACT:     defaults.DocArtifactTypeName,
	openapi.ARTIFACTTYPEQUERYPARAM_DATASET_ARTIFACT: defaults.DataSetTypeName,
	openapi.ARTIFACTTYPEQUERYPARAM_METRIC:           defaults.MetricTypeName,
	openapi.ARTIFACTTYPEQUERYPARAM_PARAMETER:        defaults.ParameterTypeName,
}

func toMlmdArtifactType(artifactType openapi.ArtifactTypeQueryParam) (string, error) {
	t, ok := artifactTypeNameMap[artifactType]
	if !ok {
		return "", fmt.Errorf("unknown artifact type: %v", artifactType)
	}
	return t, nil
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

// GetModelArtifactByParams retrieves a model artifact based on specified parameters, such as (artifact name and parent resource ID), or external ID.
// If multiple or no model artifacts are found, an error is returned.
func (serv *ModelRegistryService) GetModelArtifactByParams(artifactName *string, parentResourceId *string, externalId *string) (*openapi.ModelArtifact, error) {
	var artifact0 *proto.Artifact

	filterQuery := ""
	if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else if artifactName != nil && parentResourceId != nil {
		// search for the artifact name or the artifact name with the parent resource id as a prefix
		// and the parent resource id is the same as parentResourceId
		filterQuery = fmt.Sprintf("(name = \"%s\" or name like \"%%:%s\") and contexts_a.id = %s",
			*artifactName, *artifactName, *parentResourceId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (artifactName and parentResourceId), or externalId: %w", api.ErrBadRequest)
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
		return nil, fmt.Errorf("multiple model artifacts found for artifactName=%v, parentResourceId=%v, externalId=%v: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(parentResourceId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(artifactsResponse.Artifacts) == 0 {
		return nil, fmt.Errorf("no model artifacts found for artifactName=%v, parentResourceId=%v, externalId=%v: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(parentResourceId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	artifact0 = artifactsResponse.Artifacts[0]

	result, err := serv.mapper.MapToModelArtifact(artifact0)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return result, nil
}

// GetModelArtifacts retrieves a list of model artifacts based on the provided list options and optional parent resource ID.
func (serv *ModelRegistryService) GetModelArtifacts(listOptions api.ListOptions, parentResourceId *string) (*openapi.ModelArtifactList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	var artifacts []*proto.Artifact
	var nextPageToken *string
	if parentResourceId != nil {
		ctxId, err := converter.StringToInt64(parentResourceId)
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

// GetContextByID retrieves a context by its unique identifier (ID).
func (serv *ModelRegistryService) GetContextByID(id string) (*proto.Context, error) {
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
		return nil, fmt.Errorf("multiple contexts found for id %s: %w", id, api.ErrNotFound)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no context found for id %s: %w", id, api.ErrNotFound)
	}

	return getByIdResp.Contexts[0], nil
}

func (serv *ModelRegistryService) getContextsByArtifactId(id string) ([]*proto.Context, error) {
	if id == "" {
		return nil, fmt.Errorf("invalid artifact id: %w", api.ErrBadRequest)
	}

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	response, err := serv.mlmdClient.GetContextsByArtifact(context.Background(), &proto.GetContextsByArtifactRequest{
		ArtifactId: idAsInt,
	})
	if api.IgnoreNotFound(err) != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	if response == nil {
		return nil, nil
	}
	return response.Contexts, nil
}
