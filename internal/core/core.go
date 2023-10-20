package core

import (
	"context"
	"fmt"
	"log"
	"strconv"

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
	log.Printf("Creating or updating registered model for %s", *registeredModel.Name)

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

	modelId := &modelCtxResp.ContextIds[0]
	model, err := serv.GetRegisteredModelById((*BaseResourceId)(modelId))
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (serv *modelRegistryService) GetRegisteredModelById(id *BaseResourceId) (*openapi.RegisteredModel, error) {
	log.Printf("Getting registered model %d", *id)

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{int64(*id)},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) != 1 {
		return nil, fmt.Errorf("multiple registered models found for id %d", *id)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return regModel, nil
}

func (serv *modelRegistryService) GetRegisteredModelByParams(name *string, externalId *string) (*openapi.RegisteredModel, error) {
	log.Printf("Getting registered model by params name=%v, externalId=%v", name, externalId)

	filterQuery := ""
	if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
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

func (serv *modelRegistryService) UpsertModelVersion(modelVersion *openapi.ModelVersion, registeredModelId *BaseResourceId) (*openapi.ModelVersion, error) {
	registeredModel, err := serv.GetRegisteredModelById(registeredModelId)
	if err != nil {
		return nil, fmt.Errorf("not a valid registered model id: %d", *registeredModelId)
	}
	registeredModelIdCtxID, err := mapper.IdToInt64(*registeredModel.Id)
	if err != nil {
		return nil, err
	}
	registeredModelName := registeredModel.Name
	modelCtx, err := serv.mapper.MapFromModelVersion(modelVersion, *registeredModelIdCtxID, registeredModelName)
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
	_, err = serv.mlmdClient.PutParentContexts(context.Background(), &proto.PutParentContextsRequest{
		ParentContexts: []*proto.ParentContext{{
			ChildId:  modelId,
			ParentId: registeredModelIdCtxID}},
		TransactionOptions: &proto.TransactionOptions{},
	})
	if err != nil {
		return nil, err
	}

	model, err := serv.GetModelVersionById((*BaseResourceId)(modelId))
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (serv *modelRegistryService) GetModelVersionById(id *BaseResourceId) (*openapi.ModelVersion, error) {
	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{int64(*id)},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) != 1 {
		return nil, fmt.Errorf("multiple model versions found for id %d", *id)
	}

	modelVer, err := serv.mapper.MapToModelVersion(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return modelVer, nil
}

func (serv *modelRegistryService) GetModelVersionByParams(name *string, externalId *string) (*openapi.ModelVersion, error) {
	filterQuery := ""
	if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
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
		return nil, fmt.Errorf("multiple registered models found for name=%v, externalId=%v", *name, *externalId)
	}

	modelVer, err := serv.mapper.MapToModelVersion(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return modelVer, nil
}

func (serv *modelRegistryService) GetModelVersions(listOptions ListOptions, registeredModelId *BaseResourceId) (*openapi.ModelVersionList, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	if registeredModelId != nil {
		queryParentCtxId := fmt.Sprintf("parent_contexts_a.type = %d", *registeredModelId)
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

func (serv *modelRegistryService) UpsertModelArtifact(modelArtifact *openapi.ModelArtifact, modelVersionId *BaseResourceId) (*openapi.ModelArtifact, error) {
	artifact := serv.mapper.MapFromModelArtifact(*modelArtifact)

	artifactsResp, err := serv.mlmdClient.PutArtifacts(context.Background(), &proto.PutArtifactsRequest{
		Artifacts: []*proto.Artifact{artifact},
	})
	if err != nil {
		return nil, err
	}
	idString := strconv.FormatInt(artifactsResp.ArtifactIds[0], 10)
	modelArtifact.Id = &idString

	// add explicit association between artifacts and model version
	if modelVersionId != nil {
		modelVersionIdCtx := int64(*modelVersionId)
		attributions := []*proto.Attribution{}
		for _, a := range artifactsResp.ArtifactIds {
			attributions = append(attributions, &proto.Attribution{
				ContextId:  &modelVersionIdCtx,
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

	return modelArtifact, nil
}

func (serv *modelRegistryService) GetModelArtifactById(id *BaseResourceId) (*openapi.ModelArtifact, error) {
	artifactsResp, err := serv.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{int64(*id)},
	})
	if err != nil {
		return nil, err
	}

	result, err := serv.mapper.MapToModelArtifact(artifactsResp.Artifacts[0])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (serv *modelRegistryService) GetModelArtifactByParams(name *string, externalId *string) (*openapi.ModelArtifact, error) {
	var artifact0 *proto.Artifact

	filterQuery := ""
	if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId")
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

func (serv *modelRegistryService) GetModelArtifacts(listOptions ListOptions, modelVersionId *BaseResourceId) (*openapi.ModelArtifactList, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	var artifacts []*proto.Artifact
	var nextPageToken *string
	if modelVersionId != nil {
		ctxId := int64(*modelVersionId)
		artifactsResp, err := serv.mlmdClient.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
			ContextId: &ctxId,
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
