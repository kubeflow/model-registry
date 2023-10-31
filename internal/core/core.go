package core

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/glog"
	"github.com/opendatahub-io/model-registry/internal/core/mapper"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"google.golang.org/grpc"
)

var (
	RegisteredModelTypeName = "odh.RegisteredModel"
	ModelVersionTypeName    = "odh.ModelVersion"
	ModelArtifactTypeName   = "odh.ModelArtifact"
)

// modelRegistryService is the core library of the model registry
type modelRegistryService struct {
	mlmdClient proto.MetadataStoreServiceClient
	mapper     *mapper.Mapper
}

// NewModelRegistryService create a fresh instance of ModelRegistryService, taking care of setting up needed MLMD Types
func NewModelRegistryService(cc grpc.ClientConnInterface) (ModelRegistryApi, error) {

	client := proto.NewMetadataStoreServiceClient(cc)

	// Setup the needed Type instances if not existing already

	registeredModelReq := proto.PutContextTypeRequest{
		ContextType: &proto.ContextType{
			Name: &RegisteredModelTypeName,
		},
	}

	modelVersionReq := proto.PutContextTypeRequest{
		ContextType: &proto.ContextType{
			Name: &ModelVersionTypeName,
			Properties: map[string]proto.PropertyType{
				"model_name": proto.PropertyType_STRING,
				"version":    proto.PropertyType_STRING,
				"author":     proto.PropertyType_STRING,
			},
		},
	}

	modelArtifactReq := proto.PutArtifactTypeRequest{
		ArtifactType: &proto.ArtifactType{
			Name: &ModelArtifactTypeName,
			Properties: map[string]proto.PropertyType{
				"model_format": proto.PropertyType_STRING,
			},
		},
	}

	registeredModelResp, err := client.PutContextType(context.Background(), &registeredModelReq)
	if err != nil {
		log.Fatalf("Error setting up context type %s: %v", RegisteredModelTypeName, err)
	}

	modelVersionResp, err := client.PutContextType(context.Background(), &modelVersionReq)
	if err != nil {
		log.Fatalf("Error setting up context type %s: %v", ModelVersionTypeName, err)
	}
	modelArtifactResp, err := client.PutArtifactType(context.Background(), &modelArtifactReq)
	if err != nil {
		log.Fatalf("Error setting up artifact type %s: %v", ModelArtifactTypeName, err)
	}

	return &modelRegistryService{
		mlmdClient: client,
		mapper:     mapper.NewMapper(registeredModelResp.GetTypeId(), modelVersionResp.GetTypeId(), modelArtifactResp.GetTypeId()),
	}, nil
}

// REGISTERED MODELS

func (serv *modelRegistryService) UpsertRegisteredModel(registeredModel *openapi.RegisteredModel) (*openapi.RegisteredModel, error) {
	var err error
	var existing *openapi.RegisteredModel

	if registeredModel.Id == nil {
		glog.Info("Creating new registered model")
	} else {
		glog.Info("Updating registered model %s", *registeredModel.Id)
		existing, err = serv.GetRegisteredModelById(*registeredModel.Id)
		if err != nil {
			return nil, err
		}
	}

	// if already existing assure the name is the same
	if existing != nil && registeredModel.Name == nil {
		// user did not provide it
		// need to set it to avoid mlmd error "context name should not be empty"
		registeredModel.Name = existing.Name
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

	idAsString := mapper.IdToString(modelCtxResp.ContextIds[0])
	model, err := serv.GetRegisteredModelById(*idAsString)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (serv *modelRegistryService) GetRegisteredModelById(id string) (*openapi.RegisteredModel, error) {
	glog.Info("Getting registered model %s", id)

	idAsInt, err := mapper.IdToInt64(id)
	if err != nil {
		return nil, err
	}

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple registered models found for id %s", id)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no registered model found for id %s", id)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return regModel, nil
}

func (serv *modelRegistryService) getRegisteredModelByVersionId(id string) (*openapi.RegisteredModel, error) {
	glog.Info("Getting registered model for model version %s", id)

	idAsInt, err := mapper.IdToInt64(id)
	if err != nil {
		return nil, err
	}

	getParentResp, err := serv.mlmdClient.GetParentContextsByContext(context.Background(), &proto.GetParentContextsByContextRequest{
		ContextId: idAsInt,
	})
	if err != nil {
		return nil, err
	}

	if len(getParentResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple registered models found for model version %s", id)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no registered model found for model version %s", id)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getParentResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return regModel, nil
}

func (serv *modelRegistryService) GetRegisteredModelByParams(name *string, externalId *string) (*openapi.RegisteredModel, error) {
	glog.Info("Getting registered model by params name=%v, externalId=%v", name, externalId)

	filterQuery := ""
	if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId")
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &RegisteredModelTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) != 1 {
		return nil, fmt.Errorf("multiple registered models found for name=%v, externalId=%v", *name, *externalId)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return regModel, nil
}

func (serv *modelRegistryService) GetRegisteredModels(listOptions ListOptions) (*openapi.RegisteredModelList, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}
	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &RegisteredModelTypeName,
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
		NextPageToken: zeroIfNil(contextsResp.NextPageToken),
		PageSize:      zeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// MODEL VERSIONS

func (serv *modelRegistryService) UpsertModelVersion(modelVersion *openapi.ModelVersion, parentResourceId *string) (*openapi.ModelVersion, error) {
	var err error
	var existing *openapi.ModelVersion
	var registeredModel *openapi.RegisteredModel

	if modelVersion.Id == nil {
		// create
		glog.Info("Creating new model version")
		if parentResourceId == nil {
			return nil, fmt.Errorf("missing registered model id, cannot create model version without registered model")
		}
		registeredModel, err = serv.GetRegisteredModelById(*parentResourceId)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		glog.Info("Updating model version %s", *modelVersion.Id)
		existing, err = serv.GetModelVersionById(*modelVersion.Id)
		if err != nil {
			return nil, err
		}
		registeredModel, err = serv.getRegisteredModelByVersionId(*modelVersion.Id)
		if err != nil {
			return nil, err
		}
	}

	registeredModelId, err := mapper.IdToInt64(*registeredModel.Id)
	if err != nil {
		return nil, err
	}

	// if already existing assure the name is the same
	if existing != nil && modelVersion.Name == nil {
		// user did not provide it
		// need to set it to avoid mlmd error "context name should not be empty"
		modelVersion.Name = existing.Name
	}

	modelCtx, err := serv.mapper.MapFromModelVersion(modelVersion, *registeredModelId, registeredModel.Name)
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

	modelId := &modelCtxResp.ContextIds[0]
	if modelVersion.Id == nil {
		_, err = serv.mlmdClient.PutParentContexts(context.Background(), &proto.PutParentContextsRequest{
			ParentContexts: []*proto.ParentContext{{
				ChildId:  modelId,
				ParentId: registeredModelId}},
			TransactionOptions: &proto.TransactionOptions{},
		})
		if err != nil {
			return nil, err
		}
	}

	idAsString := mapper.IdToString(*modelId)
	model, err := serv.GetModelVersionById(*idAsString)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (serv *modelRegistryService) GetModelVersionById(id string) (*openapi.ModelVersion, error) {
	idAsInt, err := mapper.IdToInt64(id)
	if err != nil {
		return nil, err
	}

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple model versions found for id %s", id)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no model version found for id %s", id)
	}

	modelVer, err := serv.mapper.MapToModelVersion(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return modelVer, nil
}

func (serv *modelRegistryService) getModelVersionByArtifactId(id string) (*openapi.ModelVersion, error) {
	glog.Info("Getting model version for model artifact %s", id)

	idAsInt, err := mapper.IdToInt64(id)
	if err != nil {
		return nil, err
	}

	getParentResp, err := serv.mlmdClient.GetContextsByArtifact(context.Background(), &proto.GetContextsByArtifactRequest{
		ArtifactId: idAsInt,
	})
	if err != nil {
		return nil, err
	}

	if len(getParentResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple model versions found for model artifact %s", id)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no model version found for model artifact %s", id)
	}

	modelVersion, err := serv.mapper.MapToModelVersion(getParentResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return modelVersion, nil
}

func (serv *modelRegistryService) GetModelVersionByParams(versionName *string, parentResourceId *string, externalId *string) (*openapi.ModelVersion, error) {
	filterQuery := ""
	if versionName != nil && parentResourceId != nil {
		idAsInt, err := mapper.IdToInt64(*parentResourceId)
		if err != nil {
			return nil, err
		}
		filterQuery = fmt.Sprintf("name = \"%s\"", mapper.PrefixWhenOwned(idAsInt, *versionName))
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (versionName and parentResourceId), or externalId")
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &ModelVersionTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) != 1 {
		return nil, fmt.Errorf("multiple model versions found for versionName=%v, parentResourceId=%v, externalId=%v", zeroIfNil(versionName), zeroIfNil(parentResourceId), zeroIfNil(externalId))
	}

	modelVer, err := serv.mapper.MapToModelVersion(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return modelVer, nil
}

func (serv *modelRegistryService) GetModelVersions(listOptions ListOptions, parentResourceId *string) (*openapi.ModelVersionList, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	if parentResourceId != nil {
		queryParentCtxId := fmt.Sprintf("parent_contexts_a.id = %s", *parentResourceId)
		listOperationOptions.FilterQuery = &queryParentCtxId
	}

	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &ModelVersionTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.ModelVersion{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToModelVersion(c)
		if err != nil {
			return nil, err
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ModelVersionList{
		NextPageToken: zeroIfNil(contextsResp.NextPageToken),
		PageSize:      zeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// MODEL ARTIFACTS

func (serv *modelRegistryService) UpsertModelArtifact(modelArtifact *openapi.ModelArtifact, parentResourceId *string) (*openapi.ModelArtifact, error) {
	var err error
	var existing *openapi.ModelArtifact

	if modelArtifact.Id == nil {
		// create
		glog.Info("Creating new model artifact")
		if parentResourceId == nil {
			return nil, fmt.Errorf("missing model version id, cannot create model artifact without model version")
		}
		_, err = serv.GetModelVersionById(*parentResourceId)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		glog.Info("Updating model artifact %s", *modelArtifact.Id)
		existing, err = serv.GetModelArtifactById(*modelArtifact.Id)
		if err != nil {
			return nil, err
		}
		_, err = serv.getModelVersionByArtifactId(*modelArtifact.Id)
		if err != nil {
			return nil, err
		}
	}

	modelVersionId, err := mapper.IdToInt64(*parentResourceId)
	if err != nil {
		return nil, err
	}

	// if already existing assure the name is the same
	if existing != nil {
		if modelArtifact.Name == nil {
			// user did not provide it
			// need to set it to avoid mlmd error "artifact name should not be empty"
			modelArtifact.Name = existing.Name
		}
	}

	artifact := serv.mapper.MapFromModelArtifact(*modelArtifact, modelVersionId)

	artifactsResp, err := serv.mlmdClient.PutArtifacts(context.Background(), &proto.PutArtifactsRequest{
		Artifacts: []*proto.Artifact{artifact},
	})
	if err != nil {
		return nil, err
	}

	// add explicit association between artifacts and model version
	if parentResourceId != nil && modelArtifact.Id == nil {
		attributions := []*proto.Attribution{}
		for _, a := range artifactsResp.ArtifactIds {
			attributions = append(attributions, &proto.Attribution{
				ContextId:  modelVersionId,
				ArtifactId: &a,
			})
		}
		_, err = serv.mlmdClient.PutAttributionsAndAssociations(context.Background(), &proto.PutAttributionsAndAssociationsRequest{
			Attributions: attributions,
			Associations: make([]*proto.Association, 0),
		})
		if err != nil {
			return nil, err
		}
	}

	idAsString := mapper.IdToString(artifactsResp.ArtifactIds[0])
	mapped, err := serv.GetModelArtifactById(*idAsString)
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

func (serv *modelRegistryService) GetModelArtifactById(id string) (*openapi.ModelArtifact, error) {
	idAsInt, err := mapper.IdToInt64(id)
	if err != nil {
		return nil, err
	}

	artifactsResp, err := serv.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}

	if len(artifactsResp.Artifacts) > 1 {
		return nil, fmt.Errorf("multiple model artifacts found for id %s", id)
	}

	if len(artifactsResp.Artifacts) == 0 {
		return nil, fmt.Errorf("no model artifact found for id %s", id)
	}

	result, err := serv.mapper.MapToModelArtifact(artifactsResp.Artifacts[0])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (serv *modelRegistryService) GetModelArtifactByParams(artifactName *string, parentResourceId *string, externalId *string) (*openapi.ModelArtifact, error) {
	var artifact0 *proto.Artifact

	filterQuery := ""
	if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else if artifactName != nil && parentResourceId != nil {
		idAsInt, err := mapper.IdToInt64(*parentResourceId)
		if err != nil {
			return nil, err
		}
		filterQuery = fmt.Sprintf("name = \"%s\"", mapper.PrefixWhenOwned(idAsInt, *artifactName))
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (artifactName and parentResourceId), or externalId")
	}

	artifactsResponse, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
		TypeName: &ModelArtifactTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(artifactsResponse.Artifacts) > 1 {
		return nil, fmt.Errorf("more than an artifact detected matching criteria: %v", artifactsResponse.Artifacts)
	}
	artifact0 = artifactsResponse.Artifacts[0]

	result, err := serv.mapper.MapToModelArtifact(artifact0)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (serv *modelRegistryService) GetModelArtifacts(listOptions ListOptions, parentResourceId *string) (*openapi.ModelArtifactList, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	var artifacts []*proto.Artifact
	var nextPageToken *string
	if parentResourceId != nil {
		ctxId, err := mapper.IdToInt64(*parentResourceId)
		if err != nil {
			return nil, err
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
		artifactsResp, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
			TypeName: &ModelArtifactTypeName,
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
			return nil, err
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ModelArtifactList{
		NextPageToken: zeroIfNil(nextPageToken),
		PageSize:      zeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}
