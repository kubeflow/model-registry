package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/opendatahub-io/model-registry/internal/apiutils"
	"github.com/opendatahub-io/model-registry/internal/constants"
	"github.com/opendatahub-io/model-registry/internal/converter"
	"github.com/opendatahub-io/model-registry/internal/converter/generated"
	"github.com/opendatahub-io/model-registry/internal/mapper"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/pkg/api"
	"github.com/opendatahub-io/model-registry/pkg/openapi"
	"google.golang.org/grpc"
)

var (
	registeredModelTypeName    = apiutils.Of(constants.RegisteredModelTypeName)
	modelVersionTypeName       = apiutils.Of(constants.ModelVersionTypeName)
	modelArtifactTypeName      = apiutils.Of(constants.ModelArtifactTypeName)
	servingEnvironmentTypeName = apiutils.Of(constants.ServingEnvironmentTypeName)
	inferenceServiceTypeName   = apiutils.Of(constants.InferenceServiceTypeName)
	serveModelTypeName         = apiutils.Of(constants.ServeModelTypeName)
)

// ModelRegistryService is the core library of the model registry
type ModelRegistryService struct {
	mlmdClient  proto.MetadataStoreServiceClient
	typesMap    map[string]int64
	mapper      *mapper.Mapper
	openapiConv *generated.OpenAPIConverterImpl
}

// NewModelRegistryService creates a new instance of the ModelRegistryService, initializing it with the provided gRPC client connection.
// It _assumes_ the necessary MLMD's Context, Artifact, Execution types etc. are already setup in the underlying MLMD service.
//
// Parameters:
//   - cc: A gRPC client connection to the underlying MLMD service
func NewModelRegistryService(cc grpc.ClientConnInterface) (api.ModelRegistryApi, error) {
	typesMap, err := BuildTypesMap(cc)
	if err != nil { // early return in case type Ids cannot be retrieved
		return nil, err
	}

	client := proto.NewMetadataStoreServiceClient(cc)

	return &ModelRegistryService{
		mlmdClient:  client,
		typesMap:    typesMap,
		openapiConv: &generated.OpenAPIConverterImpl{},
		mapper:      mapper.NewMapper(typesMap),
	}, nil
}

func BuildTypesMap(cc grpc.ClientConnInterface) (map[string]int64, error) {
	client := proto.NewMetadataStoreServiceClient(cc)

	registeredModelContextTypeReq := proto.GetContextTypeRequest{
		TypeName: registeredModelTypeName,
	}
	registeredModelResp, err := client.GetContextType(context.Background(), &registeredModelContextTypeReq)
	if err != nil {
		return nil, fmt.Errorf("error getting context type %s: %v", *registeredModelTypeName, err)
	}
	modelVersionContextTypeReq := proto.GetContextTypeRequest{
		TypeName: modelVersionTypeName,
	}
	modelVersionResp, err := client.GetContextType(context.Background(), &modelVersionContextTypeReq)
	if err != nil {
		return nil, fmt.Errorf("error getting context type %s: %v", *modelVersionTypeName, err)
	}
	modelArtifactArtifactTypeReq := proto.GetArtifactTypeRequest{
		TypeName: modelArtifactTypeName,
	}
	modelArtifactResp, err := client.GetArtifactType(context.Background(), &modelArtifactArtifactTypeReq)
	if err != nil {
		return nil, fmt.Errorf("error getting artifact type %s: %v", *modelArtifactTypeName, err)
	}
	servingEnvironmentContextTypeReq := proto.GetContextTypeRequest{
		TypeName: servingEnvironmentTypeName,
	}
	servingEnvironmentResp, err := client.GetContextType(context.Background(), &servingEnvironmentContextTypeReq)
	if err != nil {
		return nil, fmt.Errorf("error getting context type %s: %v", *servingEnvironmentTypeName, err)
	}
	inferenceServiceContextTypeReq := proto.GetContextTypeRequest{
		TypeName: inferenceServiceTypeName,
	}
	inferenceServiceResp, err := client.GetContextType(context.Background(), &inferenceServiceContextTypeReq)
	if err != nil {
		return nil, fmt.Errorf("error getting context type %s: %v", *inferenceServiceTypeName, err)
	}
	serveModelExecutionReq := proto.GetExecutionTypeRequest{
		TypeName: serveModelTypeName,
	}
	serveModelResp, err := client.GetExecutionType(context.Background(), &serveModelExecutionReq)
	if err != nil {
		return nil, fmt.Errorf("error getting execution type %s: %v", *serveModelTypeName, err)
	}

	typesMap := map[string]int64{
		constants.RegisteredModelTypeName:    registeredModelResp.ContextType.GetId(),
		constants.ModelVersionTypeName:       modelVersionResp.ContextType.GetId(),
		constants.ModelArtifactTypeName:      modelArtifactResp.ArtifactType.GetId(),
		constants.ServingEnvironmentTypeName: servingEnvironmentResp.ContextType.GetId(),
		constants.InferenceServiceTypeName:   inferenceServiceResp.ContextType.GetId(),
		constants.ServeModelTypeName:         serveModelResp.ExecutionType.GetId(),
	}
	return typesMap, nil
}

// REGISTERED MODELS

// UpsertRegisteredModel creates a new registered model if the given registered model's ID is nil,
// or updates an existing registered model if the ID is provided.
func (serv *ModelRegistryService) UpsertRegisteredModel(registeredModel *openapi.RegisteredModel) (*openapi.RegisteredModel, error) {
	var err error
	var existing *openapi.RegisteredModel

	if registeredModel.Id == nil {
		glog.Info("Creating new registered model")
	} else {
		glog.Infof("Updating registered model %s", *registeredModel.Id)
		existing, err = serv.GetRegisteredModelById(*registeredModel.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := serv.openapiConv.OverrideNotEditableForRegisteredModel(converter.NewOpenapiUpdateWrapper(existing, registeredModel))
		if err != nil {
			return nil, err
		}
		registeredModel = &withNotEditable
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

	idAsString := converter.Int64ToString(&modelCtxResp.ContextIds[0])
	model, err := serv.GetRegisteredModelById(*idAsString)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// GetRegisteredModelById retrieves a registered model by its unique identifier (ID).
func (serv *ModelRegistryService) GetRegisteredModelById(id string) (*openapi.RegisteredModel, error) {
	glog.Infof("Getting registered model %s", id)

	idAsInt, err := converter.StringToInt64(&id)
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

// GetRegisteredModelByInferenceService retrieves a registered model associated with the specified inference service ID.
func (serv *ModelRegistryService) GetRegisteredModelByInferenceService(inferenceServiceId string) (*openapi.RegisteredModel, error) {
	is, err := serv.GetInferenceServiceById(inferenceServiceId)
	if err != nil {
		return nil, err
	}
	return serv.GetRegisteredModelById(is.RegisteredModelId)
}

// getRegisteredModelByVersionId retrieves a registered model associated with the specified model version ID.
func (serv *ModelRegistryService) getRegisteredModelByVersionId(id string) (*openapi.RegisteredModel, error) {
	glog.Infof("Getting registered model for model version %s", id)

	idAsInt, err := converter.StringToInt64(&id)
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

// GetRegisteredModelByParams retrieves a registered model based on specified parameters, such as name or external ID.
// If multiple or no registered models are found, an error is returned accordingly.
func (serv *ModelRegistryService) GetRegisteredModelByParams(name *string, externalId *string) (*openapi.RegisteredModel, error) {
	glog.Infof("Getting registered model by params name=%v, externalId=%v", name, externalId)

	filterQuery := ""
	if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId")
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: registeredModelTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple registered models found for name=%v, externalId=%v", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId))
	}

	if len(getByParamsResp.Contexts) == 0 {
		return nil, fmt.Errorf("no registered models found for name=%v, externalId=%v", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId))
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return regModel, nil
}

// GetRegisteredModels retrieves a list of registered models based on the provided list options.
func (serv *ModelRegistryService) GetRegisteredModels(listOptions api.ListOptions) (*openapi.RegisteredModelList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}
	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: registeredModelTypeName,
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
		NextPageToken: apiutils.ZeroIfNil(contextsResp.NextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// MODEL VERSIONS

// UpsertModelVersion creates a new model version if the provided model version's ID is nil,
// or updates an existing model version if the ID is provided.
func (serv *ModelRegistryService) UpsertModelVersion(modelVersion *openapi.ModelVersion, registeredModelId *string) (*openapi.ModelVersion, error) {
	var err error
	var existing *openapi.ModelVersion
	var registeredModel *openapi.RegisteredModel

	if modelVersion.Id == nil {
		// create
		glog.Info("Creating new model version")
		if registeredModelId == nil {
			return nil, fmt.Errorf("missing registered model id, cannot create model version without registered model")
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
			return nil, err
		}
		modelVersion = &withNotEditable

		registeredModel, err = serv.getRegisteredModelByVersionId(*modelVersion.Id)
		if err != nil {
			return nil, err
		}
	}

	modelCtx, err := serv.mapper.MapFromModelVersion(modelVersion, *registeredModel.Id, registeredModel.Name)
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
		registeredModelId, err := converter.StringToInt64(registeredModel.Id)
		if err != nil {
			return nil, err
		}

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
		return nil, fmt.Errorf("no model versions found for id %s", is.RegisteredModelId)
	}
	return &versions.Items[0], nil
}

// getModelVersionByArtifactId retrieves the model version associated with the specified model artifact ID.
func (serv *ModelRegistryService) getModelVersionByArtifactId(id string) (*openapi.ModelVersion, error) {
	glog.Infof("Getting model version for model artifact %s", id)

	idAsInt, err := converter.StringToInt64(&id)
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

// GetModelVersionByParams retrieves a model version based on specified parameters, such as (version name and registered model ID), or external ID.
// If multiple or no model versions are found, an error is returned.
func (serv *ModelRegistryService) GetModelVersionByParams(versionName *string, registeredModelId *string, externalId *string) (*openapi.ModelVersion, error) {
	filterQuery := ""
	if versionName != nil && registeredModelId != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", converter.PrefixWhenOwned(registeredModelId, *versionName))
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (versionName and registeredModelId), or externalId")
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: modelVersionTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple model versions found for versionName=%v, registeredModelId=%v, externalId=%v", apiutils.ZeroIfNil(versionName), apiutils.ZeroIfNil(registeredModelId), apiutils.ZeroIfNil(externalId))
	}

	if len(getByParamsResp.Contexts) == 0 {
		return nil, fmt.Errorf("no model versions found for versionName=%v, registeredModelId=%v, externalId=%v", apiutils.ZeroIfNil(versionName), apiutils.ZeroIfNil(registeredModelId), apiutils.ZeroIfNil(externalId))
	}

	modelVer, err := serv.mapper.MapToModelVersion(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return modelVer, nil
}

// GetModelVersions retrieves a list of model versions based on the provided list options and optional registered model ID.
func (serv *ModelRegistryService) GetModelVersions(listOptions api.ListOptions, registeredModelId *string) (*openapi.ModelVersionList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	if registeredModelId != nil {
		queryParentCtxId := fmt.Sprintf("parent_contexts_a.id = %s", *registeredModelId)
		listOperationOptions.FilterQuery = &queryParentCtxId
	}

	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: modelVersionTypeName,
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
		NextPageToken: apiutils.ZeroIfNil(contextsResp.NextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// MODEL ARTIFACTS

// UpsertModelArtifact creates a new model artifact if the provided model artifact's ID is nil,
// or updates an existing model artifact if the ID is provided.
// If a model version ID is provided and the model artifact is newly created, establishes an
// explicit attribution between the model version and the created model artifact.
func (serv *ModelRegistryService) UpsertModelArtifact(modelArtifact *openapi.ModelArtifact, modelVersionId *string) (*openapi.ModelArtifact, error) {
	var err error
	var existing *openapi.ModelArtifact

	if modelArtifact.Id == nil {
		// create
		glog.Info("Creating new model artifact")
		if modelVersionId == nil {
			return nil, fmt.Errorf("missing model version id, cannot create model artifact without model version")
		}
		_, err = serv.GetModelVersionById(*modelVersionId)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		glog.Infof("Updating model artifact %s", *modelArtifact.Id)
		existing, err = serv.GetModelArtifactById(*modelArtifact.Id)
		if err != nil {
			return nil, err
		}

		withNotEditable, err := serv.openapiConv.OverrideNotEditableForModelArtifact(converter.NewOpenapiUpdateWrapper(existing, modelArtifact))
		if err != nil {
			return nil, err
		}
		modelArtifact = &withNotEditable

		_, err = serv.getModelVersionByArtifactId(*modelArtifact.Id)
		if err != nil {
			return nil, err
		}
	}

	artifact, err := serv.mapper.MapFromModelArtifact(modelArtifact, modelVersionId)
	if err != nil {
		return nil, err
	}

	artifactsResp, err := serv.mlmdClient.PutArtifacts(context.Background(), &proto.PutArtifactsRequest{
		Artifacts: []*proto.Artifact{artifact},
	})
	if err != nil {
		return nil, err
	}

	// add explicit Attribution between Artifact and ModelVersion
	if modelVersionId != nil && modelArtifact.Id == nil {
		modelVersionId, err := converter.StringToInt64(modelVersionId)
		if err != nil {
			return nil, err
		}
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

	idAsString := converter.Int64ToString(&artifactsResp.ArtifactIds[0])
	mapped, err := serv.GetModelArtifactById(*idAsString)
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

// GetModelArtifactById retrieves a model artifact by its unique identifier (ID).
func (serv *ModelRegistryService) GetModelArtifactById(id string) (*openapi.ModelArtifact, error) {
	idAsInt, err := converter.StringToInt64(&id)
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
		return nil, fmt.Errorf("no artifacts found for model version %s", *mv.Id)
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
		return nil, fmt.Errorf("invalid parameters call, supply either (artifactName and modelVersionId), or externalId")
	}

	artifactsResponse, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
		TypeName: modelArtifactTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(artifactsResponse.Artifacts) > 1 {
		return nil, fmt.Errorf("multiple model artifacts found for artifactName=%v, modelVersionId=%v, externalId=%v", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId))
	}

	if len(artifactsResponse.Artifacts) == 0 {
		return nil, fmt.Errorf("no model artifacts found for artifactName=%v, modelVersionId=%v, externalId=%v", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId))
	}

	artifact0 = artifactsResponse.Artifacts[0]

	result, err := serv.mapper.MapToModelArtifact(artifact0)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetModelArtifacts retrieves a list of model artifacts based on the provided list options and optional model version ID.
func (serv *ModelRegistryService) GetModelArtifacts(listOptions api.ListOptions, modelVersionId *string) (*openapi.ModelArtifactList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	var artifacts []*proto.Artifact
	var nextPageToken *string
	if modelVersionId != nil {
		ctxId, err := converter.StringToInt64(modelVersionId)
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
			TypeName: modelArtifactTypeName,
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
		NextPageToken: apiutils.ZeroIfNil(nextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

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
			return nil, err
		}
		servingEnvironment = &withNotEditable
	}

	protoCtx, err := serv.mapper.MapFromServingEnvironment(servingEnvironment)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return openapiModel, nil
}

// GetServingEnvironmentById retrieves a serving environment by its unique identifier (ID).
func (serv *ModelRegistryService) GetServingEnvironmentById(id string) (*openapi.ServingEnvironment, error) {
	glog.Infof("Getting serving environment %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*idAsInt},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple serving environments found for id %s", id)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no serving environment found for id %s", id)
	}

	openapiModel, err := serv.mapper.MapToServingEnvironment(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId")
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: servingEnvironmentTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple serving environments found for name=%v, externalId=%v", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId))
	}

	if len(getByParamsResp.Contexts) == 0 {
		return nil, fmt.Errorf("no serving environments found for name=%v, externalId=%v", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(externalId))
	}

	openapiModel, err := serv.mapper.MapToServingEnvironment(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return openapiModel, nil
}

// GetServingEnvironments retrieves a list of serving environments based on the provided list options.
func (serv *ModelRegistryService) GetServingEnvironments(listOptions api.ListOptions) (*openapi.ServingEnvironmentList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}
	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: servingEnvironmentTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.ServingEnvironment{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToServingEnvironment(c)
		if err != nil {
			return nil, err
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

// INFERENCE SERVICE

// UpsertInferenceService creates a new inference service if the provided inference service's ID is nil,
// or updates an existing inference service if the ID is provided.
func (serv *ModelRegistryService) UpsertInferenceService(inferenceService *openapi.InferenceService) (*openapi.InferenceService, error) {
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
		return nil, err
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
			return nil, err
		}

		_, err = serv.mlmdClient.PutParentContexts(context.Background(), &proto.PutParentContextsRequest{
			ParentContexts: []*proto.ParentContext{{
				ChildId:  inferenceServiceId,
				ParentId: servingEnvironmentId}},
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
		return nil, err
	}

	getParentResp, err := serv.mlmdClient.GetParentContextsByContext(context.Background(), &proto.GetParentContextsByContextRequest{
		ContextId: idAsInt,
	})
	if err != nil {
		return nil, err
	}

	if len(getParentResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple ServingEnvironments found for InferenceService %s", id)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no ServingEnvironments found for InferenceService %s", id)
	}

	toReturn, err := serv.mapper.MapToServingEnvironment(getParentResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

// GetInferenceServiceById retrieves an inference service by its unique identifier (ID).
func (serv *ModelRegistryService) GetInferenceServiceById(id string) (*openapi.InferenceService, error) {
	glog.Infof("Getting InferenceService by id %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*idAsInt},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple InferenceServices found for id %s", id)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no InferenceService found for id %s", id)
	}

	toReturn, err := serv.mapper.MapToInferenceService(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("invalid parameters call, supply either (name and servingEnvironmentId), or externalId")
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: inferenceServiceTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple inference services found for name=%v, servingEnvironmentId=%v, externalId=%v", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(servingEnvironmentId), apiutils.ZeroIfNil(externalId))
	}

	if len(getByParamsResp.Contexts) == 0 {
		return nil, fmt.Errorf("no inference services found for name=%v, servingEnvironmentId=%v, externalId=%v", apiutils.ZeroIfNil(name), apiutils.ZeroIfNil(servingEnvironmentId), apiutils.ZeroIfNil(externalId))
	}

	toReturn, err := serv.mapper.MapToInferenceService(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return toReturn, nil
}

// GetInferenceServices retrieves a list of inference services based on the provided list options and optional serving environment ID and runtime.
func (serv *ModelRegistryService) GetInferenceServices(listOptions api.ListOptions, servingEnvironmentId *string, runtime *string) (*openapi.InferenceServiceList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
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
		TypeName: inferenceServiceTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.InferenceService{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToInferenceService(c)
		if err != nil {
			return nil, err
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

// SERVE MODEL

// UpsertServeModel creates a new serve model if the provided serve model's ID is nil,
// or updates an existing serve model if the ID is provided.
func (serv *ModelRegistryService) UpsertServeModel(serveModel *openapi.ServeModel, inferenceServiceId *string) (*openapi.ServeModel, error) {
	var err error
	var existing *openapi.ServeModel

	if serveModel.Id == nil {
		// create
		glog.Info("Creating new ServeModel")
		if inferenceServiceId == nil {
			return nil, fmt.Errorf("missing inferenceServiceId, cannot create ServeModel without parent resource InferenceService")
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
		return nil, err
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
			return nil, err
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
		return nil, err
	}

	getParentResp, err := serv.mlmdClient.GetContextsByExecution(context.Background(), &proto.GetContextsByExecutionRequest{
		ExecutionId: idAsInt,
	})
	if err != nil {
		return nil, err
	}

	if len(getParentResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple InferenceService found for ServeModel %s", id)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no InferenceService found for ServeModel %s", id)
	}

	toReturn, err := serv.mapper.MapToInferenceService(getParentResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

// GetServeModelById retrieves a serve model by its unique identifier (ID).
func (serv *ModelRegistryService) GetServeModelById(id string) (*openapi.ServeModel, error) {
	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	executionsResp, err := serv.mlmdClient.GetExecutionsByID(context.Background(), &proto.GetExecutionsByIDRequest{
		ExecutionIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}

	if len(executionsResp.Executions) > 1 {
		return nil, fmt.Errorf("multiple ServeModels found for id %s", id)
	}

	if len(executionsResp.Executions) == 0 {
		return nil, fmt.Errorf("no ServeModel found for id %s", id)
	}

	result, err := serv.mapper.MapToServeModel(executionsResp.Executions[0])
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetServeModels retrieves a list of serve models based on the provided list options and optional inference service ID.
func (serv *ModelRegistryService) GetServeModels(listOptions api.ListOptions, inferenceServiceId *string) (*openapi.ServeModelList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	var executions []*proto.Execution
	var nextPageToken *string
	if inferenceServiceId != nil {
		ctxId, err := converter.StringToInt64(inferenceServiceId)
		if err != nil {
			return nil, err
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
			TypeName: serveModelTypeName,
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
			return nil, err
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
