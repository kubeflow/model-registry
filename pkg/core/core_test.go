package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/converter"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/testutils"
	"github.com/opendatahub-io/model-registry/pkg/api"
	"github.com/opendatahub-io/model-registry/pkg/openapi"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// common utility test variables
var (
	// generic
	ascOrderDirection  string
	descOrderDirection string
	customString       string
	// registered models
	modelName        string
	modelDescription string
	modelExternalId  string
	owner            string
	// model version
	modelVersionName        string
	modelVersionDescription string
	versionExternalId       string
	author                  string
	// model artifact
	artifactName        string
	artifactDescription string
	artifactExtId       string
	artifactState       string
	artifactUri         string
	// entity under test
	entityName        string
	entityExternalId  string
	entityExternalId2 string
	entityDescription string
	// ServeModel
	executionState string
)

type CoreTestSuite struct {
	suite.Suite
	grpcConn   *grpc.ClientConn
	mlmdClient proto.MetadataStoreServiceClient
}

func TestRunCoreTestSuite(t *testing.T) {
	// before all
	grpcConn, mlmdClient, teardown := testutils.SetupMLMetadataTestContainer(t)
	defer teardown(t)

	coreTestSuite := CoreTestSuite{
		grpcConn:   grpcConn,
		mlmdClient: mlmdClient,
	}
	suite.Run(t, &coreTestSuite)
}

// before each test case
func (suite *CoreTestSuite) SetupTest() {
	// initialize test variable before each test
	ascOrderDirection = "ASC"
	descOrderDirection = "DESC"
	customString = "this is a customString value"
	modelName = "MyAwesomeModel"
	modelDescription = "reg model description"
	modelExternalId = "org.myawesomemodel"
	owner = "owner"
	modelVersionName = "v1"
	modelVersionDescription = "model version description"
	versionExternalId = "org.myawesomemodel@v1"
	author = "author1"
	artifactName = "Pickle model"
	artifactDescription = "artifact description"
	artifactExtId = "org.myawesomemodel@v1:pickle"
	artifactState = "LIVE"
	artifactUri = "path/to/model/v1"
	entityName = "MyAwesomeEntity"
	entityExternalId = "entityExternalID"
	entityExternalId2 = "entityExternalID2"
	entityDescription = "lorem ipsum entity description"
	executionState = "RUNNING"
}

// after each test
//   - remove the metadata sqlite file used by mlmd, this way mlmd will recreate it
func (suite *CoreTestSuite) AfterTest(suiteName, testName string) {
	if err := testutils.ClearMetadataSqliteDB(); err != nil {
		suite.Error(err)
	}
}

func (suite *CoreTestSuite) setupModelRegistryService() *ModelRegistryService {
	// setup model registry service
	service, err := NewModelRegistryService(suite.grpcConn)
	suite.Nilf(err, "error creating core service: %v", err)
	mrService, ok := service.(*ModelRegistryService)
	suite.True(ok)
	return mrService
}

// utility function that register a new simple model and return its ID
func (suite *CoreTestSuite) registerModel(service api.ModelRegistryApi, overrideModelName *string, overrideExternalId *string) string {
	registeredModel := &openapi.RegisteredModel{
		Name:        &modelName,
		ExternalID:  &modelExternalId,
		Description: &modelDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	if overrideModelName != nil {
		registeredModel.Name = overrideModelName
	}

	if overrideExternalId != nil {
		registeredModel.ExternalID = overrideExternalId
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	return *createdModel.Id
}

// utility function that register a new simple ServingEnvironment and return its ID
func (suite *CoreTestSuite) registerServingEnvironment(service api.ModelRegistryApi, overrideName *string, overrideExternalId *string) string {
	eutName := "Simple ServingEnvironment"
	eutExtID := "Simple ServingEnvironment ExtID"
	eut := &openapi.ServingEnvironment{
		Name:        &eutName,
		ExternalID:  &eutExtID,
		Description: &entityDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	if overrideName != nil {
		eut.Name = overrideName
	}

	if overrideExternalId != nil {
		eut.ExternalID = overrideExternalId
	}

	// test
	createdEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	return *createdEntity.Id
}

// utility function that register a new simple model and return its ID
func (suite *CoreTestSuite) registerModelVersion(
	service api.ModelRegistryApi,
	overrideModelName *string,
	overrideExternalId *string,
	overrideVersionName *string,
	overrideVersionExtId *string,
) string {
	registeredModelId := suite.registerModel(service, overrideModelName, overrideExternalId)

	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}

	if overrideVersionName != nil {
		modelVersion.Name = overrideVersionName
	}

	if overrideVersionExtId != nil {
		modelVersion.ExternalID = overrideVersionExtId
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating model version: %v", err)

	return *createdVersion.Id
}

// utility function that register a new simple ServingEnvironment and return its ID
func (suite *CoreTestSuite) registerInferenceService(service api.ModelRegistryApi, registerdModelId string, overrideParentResourceName *string, overrideParentResourceExternalId *string, overrideName *string, overrideExternalId *string) string {
	servingEnvironmentId := suite.registerServingEnvironment(service, overrideParentResourceName, overrideParentResourceExternalId)

	eutName := "simpleInferenceService"
	eutExtID := "simpleInferenceService ExtID"
	eut := &openapi.InferenceService{
		Name:                 &eutName,
		ExternalID:           &eutExtID,
		RegisteredModelId:    registerdModelId,
		ServingEnvironmentId: servingEnvironmentId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	if overrideName != nil {
		eut.Name = overrideName
	}
	if overrideExternalId != nil {
		eut.ExternalID = overrideExternalId
	}

	// test
	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating InferenceService: %v", err)

	return *createdEntity.Id
}

func (suite *CoreTestSuite) TestModelRegistryStartupWithExistingEmptyTypes() {
	ctx := context.Background()

	// create all types without props
	registeredModelReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: registeredModelTypeName,
		},
	}
	modelVersionReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: modelVersionTypeName,
		},
	}
	modelArtifactReq := proto.PutArtifactTypeRequest{
		CanAddFields: canAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: modelArtifactTypeName,
		},
	}
	servingEnvironmentReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: servingEnvironmentTypeName,
		},
	}
	inferenceServiceReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: inferenceServiceTypeName,
		},
	}
	serveModelReq := proto.PutExecutionTypeRequest{
		CanAddFields: canAddFields,
		ExecutionType: &proto.ExecutionType{
			Name: serveModelTypeName,
		},
	}

	_, err := suite.mlmdClient.PutContextType(context.Background(), &registeredModelReq)
	suite.Nil(err)
	_, err = suite.mlmdClient.PutContextType(context.Background(), &modelVersionReq)
	suite.Nil(err)
	_, err = suite.mlmdClient.PutArtifactType(context.Background(), &modelArtifactReq)
	suite.Nil(err)
	_, err = suite.mlmdClient.PutContextType(context.Background(), &servingEnvironmentReq)
	suite.Nil(err)
	_, err = suite.mlmdClient.PutContextType(context.Background(), &inferenceServiceReq)
	suite.Nil(err)
	_, err = suite.mlmdClient.PutExecutionType(context.Background(), &serveModelReq)
	suite.Nil(err)

	// check empty props
	regModelResp, _ := suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: registeredModelTypeName,
	})
	modelVersionResp, _ := suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: modelVersionTypeName,
	})
	modelArtifactResp, _ := suite.mlmdClient.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: modelArtifactTypeName,
	})
	servingEnvResp, _ := suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: servingEnvironmentTypeName,
	})
	inferenceServiceResp, _ := suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: inferenceServiceTypeName,
	})
	serveModelResp, _ := suite.mlmdClient.GetExecutionType(ctx, &proto.GetExecutionTypeRequest{
		TypeName: serveModelTypeName,
	})

	suite.Equal(0, len(regModelResp.ContextType.Properties))
	suite.Equal(0, len(modelVersionResp.ContextType.Properties))
	suite.Equal(0, len(modelArtifactResp.ArtifactType.Properties))
	suite.Equal(0, len(servingEnvResp.ContextType.Properties))
	suite.Equal(0, len(inferenceServiceResp.ContextType.Properties))
	suite.Equal(0, len(serveModelResp.ExecutionType.Properties))

	// create model registry service
	_, err = NewModelRegistryService(suite.grpcConn)
	suite.Nil(err)

	// assure the types have been correctly setup at startup
	// check NOT empty props
	regModelResp, _ = suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: registeredModelTypeName,
	})
	suite.NotNilf(regModelResp.ContextType, "registered model type %s should exists", *registeredModelTypeName)
	suite.Equal(*registeredModelTypeName, *regModelResp.ContextType.Name)
	suite.Equal(2, len(regModelResp.ContextType.Properties))

	modelVersionResp, _ = suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: modelVersionTypeName,
	})
	suite.NotNilf(modelVersionResp.ContextType, "model version type %s should exists", *modelVersionTypeName)
	suite.Equal(*modelVersionTypeName, *modelVersionResp.ContextType.Name)
	suite.Equal(5, len(modelVersionResp.ContextType.Properties))

	modelArtifactResp, _ = suite.mlmdClient.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: modelArtifactTypeName,
	})
	suite.NotNilf(modelArtifactResp.ArtifactType, "model artifact type %s should exists", *modelArtifactTypeName)
	suite.Equal(*modelArtifactTypeName, *modelArtifactResp.ArtifactType.Name)
	suite.Equal(6, len(modelArtifactResp.ArtifactType.Properties))

	servingEnvResp, _ = suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: servingEnvironmentTypeName,
	})
	suite.NotNilf(servingEnvResp.ContextType, "serving environment type %s should exists", *servingEnvironmentTypeName)
	suite.Equal(*servingEnvironmentTypeName, *servingEnvResp.ContextType.Name)
	suite.Equal(1, len(servingEnvResp.ContextType.Properties))

	inferenceServiceResp, _ = suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: inferenceServiceTypeName,
	})
	suite.NotNilf(inferenceServiceResp.ContextType, "inference service type %s should exists", *inferenceServiceTypeName)
	suite.Equal(*inferenceServiceTypeName, *inferenceServiceResp.ContextType.Name)
	suite.Equal(6, len(inferenceServiceResp.ContextType.Properties))

	serveModelResp, _ = suite.mlmdClient.GetExecutionType(ctx, &proto.GetExecutionTypeRequest{
		TypeName: serveModelTypeName,
	})
	suite.NotNilf(serveModelResp.ExecutionType, "serve model type %s should exists", *serveModelTypeName)
	suite.Equal(*serveModelTypeName, *serveModelResp.ExecutionType.Name)
	suite.Equal(2, len(serveModelResp.ExecutionType.Properties))
}

func (suite *CoreTestSuite) TestModelRegistryTypes() {
	// create model registry service
	_ = suite.setupModelRegistryService()

	// assure the types have been correctly setup at startup
	ctx := context.Background()
	regModelResp, _ := suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: registeredModelTypeName,
	})
	suite.NotNilf(regModelResp.ContextType, "registered model type %s should exists", *registeredModelTypeName)
	suite.Equal(*registeredModelTypeName, *regModelResp.ContextType.Name)

	modelVersionResp, _ := suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: modelVersionTypeName,
	})
	suite.NotNilf(modelVersionResp.ContextType, "model version type %s should exists", *modelVersionTypeName)
	suite.Equal(*modelVersionTypeName, *modelVersionResp.ContextType.Name)

	modelArtifactResp, _ := suite.mlmdClient.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: modelArtifactTypeName,
	})
	suite.NotNilf(modelArtifactResp.ArtifactType, "model version type %s should exists", *modelArtifactTypeName)
	suite.Equal(*modelArtifactTypeName, *modelArtifactResp.ArtifactType.Name)

	servingEnvResp, _ := suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: servingEnvironmentTypeName,
	})
	suite.NotNilf(servingEnvResp.ContextType, "serving environment type %s should exists", *servingEnvironmentTypeName)
	suite.Equal(*servingEnvironmentTypeName, *servingEnvResp.ContextType.Name)

	inferenceServiceResp, _ := suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: inferenceServiceTypeName,
	})
	suite.NotNilf(inferenceServiceResp.ContextType, "inference service type %s should exists", *inferenceServiceTypeName)
	suite.Equal(*inferenceServiceTypeName, *inferenceServiceResp.ContextType.Name)

	serveModelResp, _ := suite.mlmdClient.GetExecutionType(ctx, &proto.GetExecutionTypeRequest{
		TypeName: serveModelTypeName,
	})
	suite.NotNilf(serveModelResp.ExecutionType, "serve model type %s should exists", *serveModelTypeName)
	suite.Equal(*serveModelTypeName, *serveModelResp.ExecutionType.Name)
}

func (suite *CoreTestSuite) TestModelRegistryFailureForOmittedFieldInRegisteredModel() {
	registeredModelReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: registeredModelTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}

	_, err := suite.mlmdClient.PutContextType(context.Background(), &registeredModelReq)
	suite.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(suite.grpcConn)
	suite.NotNil(err)
	suite.Regexp("error setting up context type odh.RegisteredModel: rpc error: code = AlreadyExists.*", err.Error())
}

func (suite *CoreTestSuite) TestModelRegistryFailureForOmittedFieldInModelVersion() {
	modelVersionReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: modelVersionTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}

	_, err := suite.mlmdClient.PutContextType(context.Background(), &modelVersionReq)
	suite.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(suite.grpcConn)
	suite.NotNil(err)
	suite.Regexp("error setting up context type odh.ModelVersion: rpc error: code = AlreadyExists.*", err.Error())
}

func (suite *CoreTestSuite) TestModelRegistryFailureForOmittedFieldInModelArtifact() {
	modelArtifactReq := proto.PutArtifactTypeRequest{
		CanAddFields: canAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: modelArtifactTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}

	_, err := suite.mlmdClient.PutArtifactType(context.Background(), &modelArtifactReq)
	suite.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(suite.grpcConn)
	suite.NotNil(err)
	suite.Regexp("error setting up artifact type odh.ModelArtifact: rpc error: code = AlreadyExists.*", err.Error())
}

func (suite *CoreTestSuite) TestModelRegistryFailureForOmittedFieldInServingEnvironment() {
	servingEnvironmentReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: servingEnvironmentTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}
	_, err := suite.mlmdClient.PutContextType(context.Background(), &servingEnvironmentReq)
	suite.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(suite.grpcConn)
	suite.NotNil(err)
	suite.Regexp("error setting up context type odh.ServingEnvironment: rpc error: code = AlreadyExists.*", err.Error())
}

func (suite *CoreTestSuite) TestModelRegistryFailureForOmittedFieldInInferenceService() {
	inferenceServiceReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: inferenceServiceTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}

	_, err := suite.mlmdClient.PutContextType(context.Background(), &inferenceServiceReq)
	suite.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(suite.grpcConn)
	suite.NotNil(err)
	suite.Regexp("error setting up context type odh.InferenceService: rpc error: code = AlreadyExists.*", err.Error())
}

func (suite *CoreTestSuite) TestModelRegistryFailureForOmittedFieldInServeModel() {
	serveModelReq := proto.PutExecutionTypeRequest{
		CanAddFields: canAddFields,
		ExecutionType: &proto.ExecutionType{
			Name: serveModelTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}

	_, err := suite.mlmdClient.PutExecutionType(context.Background(), &serveModelReq)
	suite.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(suite.grpcConn)
	suite.NotNil(err)
	suite.Regexp("error setting up execution type odh.ServeModel: rpc error: code = AlreadyExists.*", err.Error())
}

// REGISTERED MODELS

func (suite *CoreTestSuite) TestCreateRegisteredModel() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	state := openapi.REGISTEREDMODELSTATE_ARCHIVED
	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:        &modelName,
		ExternalID:  &modelExternalId,
		Description: &modelDescription,
		State:       &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	suite.Nilf(err, "error creating registered model: %v", err)
	suite.NotNilf(createdModel.Id, "created registered model should not have nil Id")

	createdModelId, _ := converter.StringToInt64(createdModel.Id)
	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdModelId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	suite.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	suite.Equal(modelName, *ctx.Name, "saved model name should match the provided one")
	suite.Equal(modelExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(modelDescription, ctx.Properties["description"].GetStringValue(), "saved description should match the provided one")
	suite.Equal(string(state), ctx.Properties["state"].GetStringValue(), "saved state should match the provided one")
	suite.Equal(owner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")
}

func (suite *CoreTestSuite) TestUpdateRegisteredModel() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	suite.Nilf(err, "error creating registered model: %v", err)
	suite.NotNilf(createdModel.Id, "created registered model should not have nil Id")
	createdModelId, _ := converter.StringToInt64(createdModel.Id)

	// checks created model matches original one except for Id
	suite.Equal(*registeredModel.Name, *createdModel.Name, "returned model name should match the original one")
	suite.Equal(*registeredModel.ExternalID, *createdModel.ExternalID, "returned model external id should match the original one")
	suite.Equal(*registeredModel.CustomProperties, *createdModel.CustomProperties, "returned model custom props should match the original one")

	// update existing model
	newModelExternalId := "newExternalId"
	newOwner := "newOwner"

	createdModel.ExternalID = &newModelExternalId
	(*createdModel.CustomProperties)["owner"] = openapi.MetadataValue{
		MetadataStringValue: &openapi.MetadataStringValue{
			StringValue: &newOwner,
		},
	}

	// update the model
	createdModel, err = service.UpsertRegisteredModel(createdModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	// still one registered model
	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdModelId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	suite.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	suite.Equal(modelName, *ctx.Name, "saved model name should match the provided one")
	suite.Equal(newModelExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(newOwner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	// update the model keeping nil name
	newModelExternalId = "newNewExternalId"
	createdModel.ExternalID = &newModelExternalId
	createdModel.Name = nil
	createdModel, err = service.UpsertRegisteredModel(createdModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	// still one registered model
	getAllResp, err = suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err = suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdModelId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx = ctxById.Contexts[0]
	ctxId = converter.Int64ToString(ctx.Id)
	suite.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	suite.Equal(modelName, *ctx.Name, "saved model name should match the provided one")
	suite.Equal(newModelExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(newOwner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")
}

func (suite *CoreTestSuite) TestGetRegisteredModelById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	state := openapi.REGISTEREDMODELSTATE_LIVE
	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
		State:      &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	suite.Nilf(err, "error creating registered model: %v", err)

	getModelById, err := service.GetRegisteredModelById(*createdModel.Id)
	suite.Nilf(err, "error getting registered model by id %s: %v", *createdModel.Id, err)

	// checks created model matches original one except for Id
	suite.Equal(*registeredModel.Name, *getModelById.Name, "saved model name should match the original one")
	suite.Equal(*registeredModel.ExternalID, *getModelById.ExternalID, "saved model external id should match the original one")
	suite.Equal(*registeredModel.State, *getModelById.State, "saved model state should match the original one")
	suite.Equal(*registeredModel.CustomProperties, *getModelById.CustomProperties, "saved model custom props should match the original one")
}

func (suite *CoreTestSuite) TestGetRegisteredModelByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	_, err := service.GetRegisteredModelByParams(of("not-present"), nil)
	suite.NotNil(err)
	suite.Equal("no registered models found for name=not-present, externalId=", err.Error())
}

func (suite *CoreTestSuite) TestGetRegisteredModelByParamsName() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	byName, err := service.GetRegisteredModelByParams(&modelName, nil)
	suite.Nilf(err, "error getting registered model by name: %v", err)

	suite.Equalf(*createdModel.Id, *byName.Id, "the returned model id should match the retrieved by name")
}

func (suite *CoreTestSuite) TestGetRegisteredModelByParamsExternalId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	byName, err := service.GetRegisteredModelByParams(nil, &modelExternalId)
	suite.Nilf(err, "error getting registered model by external id: %v", err)

	suite.Equalf(*createdModel.Id, *byName.Id, "the returned model id should match the retrieved by name")
}

func (suite *CoreTestSuite) TestGetRegisteredModelByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	_, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	_, err = service.GetRegisteredModelByParams(nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either name or externalId", err.Error())
}

func (suite *CoreTestSuite) TestGetRegisteredModelsOrderedById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	orderBy := "ID"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	_, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	_, err = service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	_, err = service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	orderedById, err := service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		suite.Less(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}

	orderedById, err = service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		suite.Greater(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}
}

func (suite *CoreTestSuite) TestGetRegisteredModelsOrderedByLastUpdate() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	orderBy := "LAST_UPDATE_TIME"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	firstModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	secondModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	thirdModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	// update second model
	secondModel.ExternalID = nil
	_, err = service.UpsertRegisteredModel(secondModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	orderedById, err := service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(3, int(orderedById.Size))
	suite.Equal(*firstModel.Id, *orderedById.Items[0].Id)
	suite.Equal(*thirdModel.Id, *orderedById.Items[1].Id)
	suite.Equal(*secondModel.Id, *orderedById.Items[2].Id)

	orderedById, err = service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(3, int(orderedById.Size))
	suite.Equal(*secondModel.Id, *orderedById.Items[0].Id)
	suite.Equal(*thirdModel.Id, *orderedById.Items[1].Id)
	suite.Equal(*firstModel.Id, *orderedById.Items[2].Id)
}

func (suite *CoreTestSuite) TestGetRegisteredModelsWithPageSize() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	pageSize := int32(1)
	pageSize2 := int32(2)
	modelName := "PricingModel1"
	modelExternalId := "myExternalId1"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	firstModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	secondModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	thirdModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	truncatedList, err := service.GetRegisteredModels(api.ListOptions{
		PageSize: &pageSize,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(1, int(truncatedList.Size))
	suite.NotEqual("", truncatedList.NextPageToken, "next page token should not be empty")
	suite.Equal(*firstModel.Id, *truncatedList.Items[0].Id)

	truncatedList, err = service.GetRegisteredModels(api.ListOptions{
		PageSize:      &pageSize2,
		NextPageToken: &truncatedList.NextPageToken,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(2, int(truncatedList.Size))
	suite.Equal("", truncatedList.NextPageToken, "next page token should be empty as list item returned")
	suite.Equal(*secondModel.Id, *truncatedList.Items[0].Id)
	suite.Equal(*thirdModel.Id, *truncatedList.Items[1].Id)
}

// MODEL VERSIONS

func (suite *CoreTestSuite) TestCreateModelVersion() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	state := openapi.MODELVERSIONSTATE_LIVE
	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		State:       &state,
		Author:      &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")

	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*createdVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	suite.Equal(versionExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(author, byId.Contexts[0].Properties["author"].GetStringValue(), "saved author property should match the provided one")
	suite.Equal(modelVersionDescription, byId.Contexts[0].Properties["description"].GetStringValue(), "saved description should match the provided one")
	suite.Equal(string(state), byId.Contexts[0].Properties["state"].GetStringValue(), "saved state should match the provided one")
	suite.Equalf(*modelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *modelVersionTypeName)

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(2, len(getAllResp.Contexts), "there should be two contexts saved in mlmd")
}

func (suite *CoreTestSuite) TestCreateModelVersionFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := "9999"

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		Author:     &author,
	}

	_, err := service.UpsertModelVersion(modelVersion, nil)
	suite.NotNil(err)
	suite.Equal("missing registered model id, cannot create model version without registered model", err.Error())

	_, err = service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.NotNil(err)
	suite.Equal("no registered model found for id 9999", err.Error())
}

func (suite *CoreTestSuite) TestUpdateModelVersion() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	newExternalId := "org.my_awesome_model@v1"
	newScore := 0.95

	createdVersion.ExternalID = &newExternalId
	(*createdVersion.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: &openapi.MetadataDoubleValue{
			DoubleValue: &newScore,
		},
	}

	updatedVersion, err := service.UpsertModelVersion(createdVersion, &registeredModelId)
	suite.Nilf(err, "error updating new model version for %s: %v", registeredModelId, err)

	updateVersionId, _ := converter.StringToInt64(updatedVersion.Id)
	suite.Equal(*createdVersionId, *updateVersionId, "created and updated model version should have same id")

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*updateVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	suite.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(author, byId.Contexts[0].Properties["author"].GetStringValue(), "saved author property should match the provided one")
	suite.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	suite.Equalf(*modelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *modelVersionTypeName)

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(2, len(getAllResp.Contexts), "there should be two contexts saved in mlmd")

	// update with nil name
	newExternalId = "org.my_awesome_model_@v1"
	updatedVersion.ExternalID = &newExternalId
	updatedVersion.Name = nil
	updatedVersion, err = service.UpsertModelVersion(updatedVersion, &registeredModelId)
	suite.Nilf(err, "error updating new model version for %s: %v", registeredModelId, err)

	updateVersionId, _ = converter.StringToInt64(updatedVersion.Id)
	suite.Equal(*createdVersionId, *updateVersionId, "created and updated model version should have same id")

	byId, err = suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*updateVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	suite.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(author, byId.Contexts[0].Properties["author"].GetStringValue(), "saved author property should match the provided one")
	suite.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	suite.Equalf(*modelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *modelVersionTypeName)
}

func (suite *CoreTestSuite) TestUpdateModelVersionFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %s", registeredModelId)
	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")

	newExternalId := "org.my_awesome_model@v1"
	newScore := 0.95

	createdVersion.ExternalID = &newExternalId
	(*createdVersion.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: &openapi.MetadataDoubleValue{
			DoubleValue: &newScore,
		},
	}

	wrongId := "9999"
	createdVersion.Id = &wrongId
	_, err = service.UpsertModelVersion(createdVersion, &registeredModelId)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no model version found for id %s", wrongId), err.Error())
}

func (suite *CoreTestSuite) TestGetModelVersionById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	state := openapi.MODELVERSIONSTATE_ARCHIVED
	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		State:      &state,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	getById, err := service.GetModelVersionById(*createdVersion.Id)
	suite.Nilf(err, "error getting model version with id %d", *createdVersionId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getById.Id, "returned model version id should match the mlmd context one")
	suite.Equal(*modelVersion.Name, *getById.Name, "saved model name should match the provided one")
	suite.Equal(*modelVersion.ExternalID, *getById.ExternalID, "saved external id should match the provided one")
	suite.Equal(*modelVersion.State, *getById.State, "saved model state should match the original one")
	suite.Equal(*getById.Author, author, "saved author property should match the provided one")
}

func (suite *CoreTestSuite) TestGetModelVersionByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	_, err := service.GetModelVersionByParams(of("not-present"), &registeredModelId, nil)
	suite.NotNil(err)
	suite.Equal("no model versions found for versionName=not-present, parentResourceId=1, externalId=", err.Error())
}

func (suite *CoreTestSuite) TestGetModelVersionByParamsName() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	getByName, err := service.GetModelVersionByParams(&modelVersionName, &registeredModelId, nil)
	suite.Nilf(err, "error getting model version by name %d", *createdVersionId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getByName.Id, "returned model version id should match the mlmd context one")
	suite.Equal(fmt.Sprintf("%s:%s", registeredModelId, *getByName.Name), *ctx.Name, "saved model name should match the provided one")
	suite.Equal(*ctx.ExternalId, *getByName.ExternalID, "saved external id should match the provided one")
	suite.Equal(ctx.Properties["author"].GetStringValue(), *getByName.Author, "saved author property should match the provided one")
}

func (suite *CoreTestSuite) TestGetModelVersionByParamsExternalId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	getByExternalId, err := service.GetModelVersionByParams(nil, nil, modelVersion.ExternalID)
	suite.Nilf(err, "error getting model version by external id %d", *modelVersion.ExternalID)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getByExternalId.Id, "returned model version id should match the mlmd context one")
	suite.Equal(fmt.Sprintf("%s:%s", registeredModelId, *getByExternalId.Name), *ctx.Name, "saved model name should match the provided one")
	suite.Equal(*ctx.ExternalId, *getByExternalId.ExternalID, "saved external id should match the provided one")
	suite.Equal(ctx.Properties["author"].GetStringValue(), *getByExternalId.Author, "saved author property should match the provided one")
}

func (suite *CoreTestSuite) TestGetModelVersionByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")

	_, err = service.GetModelVersionByParams(nil, nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either (versionName and parentResourceId), or externalId", err.Error())
}

func (suite *CoreTestSuite) TestGetModelVersions() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion1 := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
	}

	secondModelVersionName := "v2"
	secondModelVersionExtId := "org.myawesomemodel@v2"
	modelVersion2 := &openapi.ModelVersion{
		Name:       &secondModelVersionName,
		ExternalID: &secondModelVersionExtId,
	}

	thirdModelVersionName := "v3"
	thirdModelVersionExtId := "org.myawesomemodel@v3"
	modelVersion3 := &openapi.ModelVersion{
		Name:       &thirdModelVersionName,
		ExternalID: &thirdModelVersionExtId,
	}

	createdVersion1, err := service.UpsertModelVersion(modelVersion1, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	createdVersion3, err := service.UpsertModelVersion(modelVersion3, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	anotherRegModelName := "AnotherModel"
	anotherRegModelExtId := "org.another"
	anotherRegisteredModelId := suite.registerModel(service, &anotherRegModelName, &anotherRegModelExtId)

	anotherModelVersionName := "v1.0"
	anotherModelVersionExtId := "org.another@v1.0"
	modelVersionAnother := &openapi.ModelVersion{
		Name:       &anotherModelVersionName,
		ExternalID: &anotherModelVersionExtId,
	}

	_, err = service.UpsertModelVersion(modelVersionAnother, &anotherRegisteredModelId)
	suite.Nilf(err, "error creating new model version for %d", anotherRegisteredModelId)

	createdVersionId1, _ := converter.StringToInt64(createdVersion1.Id)
	createdVersionId2, _ := converter.StringToInt64(createdVersion2.Id)
	createdVersionId3, _ := converter.StringToInt64(createdVersion3.Id)

	getAll, err := service.GetModelVersions(api.ListOptions{}, nil)
	suite.Nilf(err, "error getting all model versions")
	suite.Equal(int32(4), getAll.Size, "expected four model versions across all registered models")

	getAllByRegModel, err := service.GetModelVersions(api.ListOptions{}, &registeredModelId)
	suite.Nilf(err, "error getting all model versions")
	suite.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	suite.Equal(*converter.Int64ToString(createdVersionId1), *getAllByRegModel.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId2), *getAllByRegModel.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId3), *getAllByRegModel.Items[2].Id)

	// order by last update time, expecting last created as first
	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByRegModel, err = service.GetModelVersions(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &registeredModelId)
	suite.Nilf(err, "error getting all model versions")
	suite.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	suite.Equal(*converter.Int64ToString(createdVersionId1), *getAllByRegModel.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId2), *getAllByRegModel.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId3), *getAllByRegModel.Items[0].Id)

	// update the second version
	newVersionExternalId := "updated.org:v2"
	createdVersion2.ExternalID = &newVersionExternalId
	createdVersion2, err = service.UpsertModelVersion(createdVersion2, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.Equal(newVersionExternalId, *createdVersion2.ExternalID)

	getAllByRegModel, err = service.GetModelVersions(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &registeredModelId)
	suite.Nilf(err, "error getting all model versions")
	suite.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	suite.Equal(*converter.Int64ToString(createdVersionId1), *getAllByRegModel.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId2), *getAllByRegModel.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId3), *getAllByRegModel.Items[1].Id)
}

// MODEL ARTIFACTS

func (suite *CoreTestSuite) TestCreateModelArtifact() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:               &artifactName,
		State:              (*openapi.ArtifactState)(&artifactState),
		Uri:                &artifactUri,
		Description:        &artifactDescription,
		ModelFormatName:    of("onnx"),
		ModelFormatVersion: of("1"),
		StorageKey:         of("aws-connection-models"),
		StoragePath:        of("bucket"),
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *createdArtifact.Name)
	suite.Equal(*state, *createdArtifact.State)
	suite.Equal(artifactUri, *createdArtifact.Uri)
	suite.Equal(artifactDescription, *createdArtifact.Description)
	suite.Equal("onnx", *createdArtifact.ModelFormatName)
	suite.Equal("1", *createdArtifact.ModelFormatVersion)
	suite.Equal("aws-connection-models", *createdArtifact.StorageKey)
	suite.Equal("bucket", *createdArtifact.StoragePath)
	suite.Equal(customString, *(*createdArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)
	getById, err := suite.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	suite.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	suite.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	suite.Equal(fmt.Sprintf("%s:%s", modelVersionId, *createdArtifact.Name), *getById.Artifacts[0].Name)
	suite.Equal(string(*createdArtifact.State), getById.Artifacts[0].State.String())
	suite.Equal(*createdArtifact.Uri, *getById.Artifacts[0].Uri)
	suite.Equal(*createdArtifact.Description, getById.Artifacts[0].Properties["description"].GetStringValue())
	suite.Equal(*createdArtifact.ModelFormatName, getById.Artifacts[0].Properties["model_format_name"].GetStringValue())
	suite.Equal(*createdArtifact.ModelFormatVersion, getById.Artifacts[0].Properties["model_format_version"].GetStringValue())
	suite.Equal(*createdArtifact.StorageKey, getById.Artifacts[0].Properties["storage_key"].GetStringValue())
	suite.Equal(*createdArtifact.StoragePath, getById.Artifacts[0].Properties["storage_path"].GetStringValue())
	suite.Equal(*(*createdArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["custom_string_prop"].GetStringValue())

	modelVersionIdAsInt, _ := converter.StringToInt64(&modelVersionId)
	byCtx, _ := suite.mlmdClient.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
		ContextId: (*int64)(modelVersionIdAsInt),
	})
	suite.Equal(1, len(byCtx.Artifacts))
	suite.Equal(*createdArtifactId, *byCtx.Artifacts[0].Id)
}

func (suite *CoreTestSuite) TestCreateModelArtifactFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := "9998"

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	_, err := service.UpsertModelArtifact(modelArtifact, nil)
	suite.NotNil(err)
	suite.Equal("missing model version id, cannot create model artifact without model version", err.Error())

	_, err = service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.NotNil(err)
	suite.Equal("no model version found for id 9998", err.Error())
}

func (suite *CoreTestSuite) TestUpdateModelArtifact() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	newState := "MARKED_FOR_DELETION"
	createdArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertModelArtifact(createdArtifact, &modelVersionId)
	suite.Nilf(err, "error updating model artifact for %d: %v", modelVersionId, err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)
	updatedArtifactId, _ := converter.StringToInt64(updatedArtifact.Id)
	suite.Equal(createdArtifactId, updatedArtifactId)

	getById, err := suite.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	suite.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	suite.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	suite.Equal(fmt.Sprintf("%s:%s", modelVersionId, *createdArtifact.Name), *getById.Artifacts[0].Name)
	suite.Equal(string(newState), getById.Artifacts[0].State.String())
	suite.Equal(*createdArtifact.Uri, *getById.Artifacts[0].Uri)
	suite.Equal(*(*createdArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["custom_string_prop"].GetStringValue())
}

func (suite *CoreTestSuite) TestUpdateModelArtifactFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for model version %s", modelVersionId)
	suite.NotNilf(createdArtifact.Id, "created model artifact should not have nil Id")

	newState := "MARKED_FOR_DELETION"
	createdArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertModelArtifact(createdArtifact, &modelVersionId)
	suite.Nilf(err, "error updating model artifact for %d: %v", modelVersionId, err)

	wrongId := "9998"
	updatedArtifact.Id = &wrongId
	_, err = service.UpsertModelArtifact(updatedArtifact, &modelVersionId)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no model artifact found for id %s", wrongId), err.Error())
}

func (suite *CoreTestSuite) TestGetModelArtifactById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)

	getById, err := service.GetModelArtifactById(*createdArtifact.Id)
	suite.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getById.Name)
	suite.Equal(*state, *getById.State)
	suite.Equal(artifactUri, *getById.Uri)
	suite.Equal(customString, *(*getById.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdArtifact, *getById, "artifacts returned during creation and on get by id should be equal")
}

func (suite *CoreTestSuite) TestGetModelArtifactByParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalID: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)

	getByName, err := service.GetModelArtifactByParams(&artifactName, &modelVersionId, nil)
	suite.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	suite.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getByName.Name)
	suite.Equal(artifactExtId, *getByName.ExternalID)
	suite.Equal(*state, *getByName.State)
	suite.Equal(artifactUri, *getByName.Uri)
	suite.Equal(customString, *(*getByName.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdArtifact, *getByName, "artifacts returned during creation and on get by name should be equal")

	getByExtId, err := service.GetModelArtifactByParams(nil, nil, &artifactExtId)
	suite.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	suite.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getByExtId.Name)
	suite.Equal(artifactExtId, *getByExtId.ExternalID)
	suite.Equal(*state, *getByExtId.State)
	suite.Equal(artifactUri, *getByExtId.Uri)
	suite.Equal(customString, *(*getByExtId.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdArtifact, *getByExtId, "artifacts returned during creation and on get by ext id should be equal")
}

func (suite *CoreTestSuite) TestGetModelArtifactByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalID: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	_, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	_, err = service.GetModelArtifactByParams(nil, nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either (artifactName and parentResourceId), or externalId", err.Error())
}

func (suite *CoreTestSuite) TestGetModelArtifactByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	_, err := service.GetModelArtifactByParams(of("not-present"), &modelVersionId, nil)
	suite.NotNil(err)
	suite.Equal("no model artifacts found for artifactName=not-present, parentResourceId=2, externalId=", err.Error())
}

func (suite *CoreTestSuite) TestGetModelArtifacts() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact1 := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalID: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	secondArtifactName := "second-name"
	secondArtifactExtId := "second-ext-id"
	secondArtifactUri := "second-uri"
	modelArtifact2 := &openapi.ModelArtifact{
		Name:       &secondArtifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &secondArtifactUri,
		ExternalID: &secondArtifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	thirdArtifactName := "third-name"
	thirdArtifactExtId := "third-ext-id"
	thirdArtifactUri := "third-uri"
	modelArtifact3 := &openapi.ModelArtifact{
		Name:       &thirdArtifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &thirdArtifactUri,
		ExternalID: &thirdArtifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdArtifact1, err := service.UpsertModelArtifact(modelArtifact1, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)
	createdArtifact2, err := service.UpsertModelArtifact(modelArtifact2, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)
	createdArtifact3, err := service.UpsertModelArtifact(modelArtifact3, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId1, _ := converter.StringToInt64(createdArtifact1.Id)
	createdArtifactId2, _ := converter.StringToInt64(createdArtifact2.Id)
	createdArtifactId3, _ := converter.StringToInt64(createdArtifact3.Id)

	getAll, err := service.GetModelArtifacts(api.ListOptions{}, nil)
	suite.Nilf(err, "error getting all model artifacts")
	suite.Equalf(int32(3), getAll.Size, "expected three model artifacts")

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAll.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAll.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId3), *getAll.Items[2].Id)

	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByModelVersion, err := service.GetModelArtifacts(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &modelVersionId)
	suite.Nilf(err, "error getting all model artifacts for %d", modelVersionId)
	suite.Equalf(int32(3), getAllByModelVersion.Size, "expected three model artifacts for model version %d", modelVersionId)

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAllByModelVersion.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAllByModelVersion.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId3), *getAllByModelVersion.Items[0].Id)
}

// SERVING ENVIRONMENT

func (suite *CoreTestSuite) TestCreateServingEnvironment() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:        &entityName,
		ExternalID:  &entityExternalId,
		Description: &entityDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	// test
	createdEntity, err := service.UpsertServingEnvironment(eut)

	// checks
	suite.Nilf(err, "error creating uut: %v", err)
	suite.NotNilf(createdEntity.Id, "created uut should not have nil Id")

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)
	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdEntityId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	suite.Equal(*createdEntity.Id, *ctxId, "returned id should match the mlmd one")
	suite.Equal(entityName, *ctx.Name, "saved name should match the provided one")
	suite.Equal(entityExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(entityDescription, ctx.Properties["description"].GetStringValue(), "saved description should match the provided one")
	suite.Equal(owner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")
}

func (suite *CoreTestSuite) TestUpdateServingEnvironment() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	// test
	createdEntity, err := service.UpsertServingEnvironment(eut)

	// checks
	suite.Nilf(err, "error creating uut: %v", err)
	suite.NotNilf(createdEntity.Id, "created uut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	// checks created entity matches original one except for Id
	suite.Equal(*eut.Name, *createdEntity.Name, "returned entity should match the original one")
	suite.Equal(*eut.ExternalID, *createdEntity.ExternalID, "returned entity external id should match the original one")
	suite.Equal(*eut.CustomProperties, *createdEntity.CustomProperties, "returned entity custom props should match the original one")

	// update existing entity
	newExternalId := "newExternalId"
	newOwner := "newOwner"

	createdEntity.ExternalID = &newExternalId
	(*createdEntity.CustomProperties)["owner"] = openapi.MetadataValue{
		MetadataStringValue: &openapi.MetadataStringValue{
			StringValue: &newOwner,
		},
	}

	// update the entity
	createdEntity, err = service.UpsertServingEnvironment(createdEntity)
	suite.Nilf(err, "error creating uut: %v", err)

	// still one expected MLMD type
	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdEntityId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	suite.Equal(*createdEntity.Id, *ctxId, "returned entity id should match the mlmd one")
	suite.Equal(entityName, *ctx.Name, "saved entity name should match the provided one")
	suite.Equal(newExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(newOwner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	// update the entity under test, keeping nil name
	newExternalId = "newNewExternalId"
	createdEntity.ExternalID = &newExternalId
	createdEntity.Name = nil
	createdEntity, err = service.UpsertServingEnvironment(createdEntity)
	suite.Nilf(err, "error creating entity: %v", err)

	// still one registered entity
	getAllResp, err = suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err = suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdEntityId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx = ctxById.Contexts[0]
	ctxId = converter.Int64ToString(ctx.Id)
	suite.Equal(*createdEntity.Id, *ctxId, "returned entity id should match the mlmd one")
	suite.Equal(entityName, *ctx.Name, "saved entity name should match the provided one")
	suite.Equal(newExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(newOwner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")
}

func (suite *CoreTestSuite) TestGetServingEnvironmentById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new entity
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	// test
	createdEntity, err := service.UpsertServingEnvironment(eut)

	// checks
	suite.Nilf(err, "error creating eut: %v", err)

	getEntityById, err := service.GetServingEnvironmentById(*createdEntity.Id)
	suite.Nilf(err, "error getting eut by id %s: %v", *createdEntity.Id, err)

	// checks created entity matches original one except for Id
	suite.Equal(*eut.Name, *getEntityById.Name, "saved name should match the original one")
	suite.Equal(*eut.ExternalID, *getEntityById.ExternalID, "saved external id should match the original one")
	suite.Equal(*eut.CustomProperties, *getEntityById.CustomProperties, "saved custom props should match the original one")
}

func (suite *CoreTestSuite) TestGetServingEnvironmentByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	_, err := service.GetServingEnvironmentByParams(of("not-present"), nil)
	suite.NotNil(err)
	suite.Equal("no serving environments found for name=not-present, externalId=", err.Error())
}

func (suite *CoreTestSuite) TestGetServingEnvironmentByParamsName() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	createdEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	byName, err := service.GetServingEnvironmentByParams(&entityName, nil)
	suite.Nilf(err, "error getting ServingEnvironment by name: %v", err)

	suite.Equalf(*createdEntity.Id, *byName.Id, "the returned entity id should match the retrieved by name")
}

func (suite *CoreTestSuite) TestGetServingEnvironmentByParamsExternalId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	createdEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	byName, err := service.GetServingEnvironmentByParams(nil, &entityExternalId)
	suite.Nilf(err, "error getting ServingEnvironment by external id: %v", err)

	suite.Equalf(*createdEntity.Id, *byName.Id, "the returned entity id should match the retrieved by name")
}

func (suite *CoreTestSuite) TestGetServingEnvironmentByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	_, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	_, err = service.GetServingEnvironmentByParams(nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either name or externalId", err.Error())
}

func (suite *CoreTestSuite) TestGetServingEnvironmentsOrderedById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	orderBy := "ID"

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	_, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName := "Pricingentity2"
	newExternalId := "myExternalId2"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	_, err = service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName = "Pricingentity3"
	newExternalId = "myExternalId3"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	_, err = service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	orderedById, err := service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	suite.Nilf(err, "error getting ServingEnvironment: %v", err)

	suite.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		suite.Less(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}

	orderedById, err = service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	suite.Nilf(err, "error getting ServingEnvironments: %v", err)

	suite.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		suite.Greater(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}
}

func (suite *CoreTestSuite) TestGetServingEnvironmentsOrderedByLastUpdate() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	orderBy := "LAST_UPDATE_TIME"

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	firstEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName := "Pricingentity2"
	newExternalId := "myExternalId2"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	secondEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName = "Pricingentity3"
	newExternalId = "myExternalId3"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	thirdEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	// update second entity
	secondEntity.ExternalID = nil
	_, err = service.UpsertServingEnvironment(secondEntity)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	orderedById, err := service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	suite.Nilf(err, "error getting ServingEnvironments: %v", err)

	suite.Equal(3, int(orderedById.Size))
	suite.Equal(*firstEntity.Id, *orderedById.Items[0].Id)
	suite.Equal(*thirdEntity.Id, *orderedById.Items[1].Id)
	suite.Equal(*secondEntity.Id, *orderedById.Items[2].Id)

	orderedById, err = service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	suite.Nilf(err, "error getting ServingEnvironments: %v", err)

	suite.Equal(3, int(orderedById.Size))
	suite.Equal(*secondEntity.Id, *orderedById.Items[0].Id)
	suite.Equal(*thirdEntity.Id, *orderedById.Items[1].Id)
	suite.Equal(*firstEntity.Id, *orderedById.Items[2].Id)
}

func (suite *CoreTestSuite) TestGetServingEnvironmentsWithPageSize() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	pageSize := int32(1)
	pageSize2 := int32(2)
	entityName := "Pricingentity1"
	entityExternalId := "myExternalId1"

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	firstEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating registered entity: %v", err)

	newName := "Pricingentity2"
	newExternalId := "myExternalId2"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	secondEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName = "Pricingentity3"
	newExternalId = "myExternalId3"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	thirdEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	truncatedList, err := service.GetServingEnvironments(api.ListOptions{
		PageSize: &pageSize,
	})
	suite.Nilf(err, "error getting ServingEnvironments: %v", err)

	suite.Equal(1, int(truncatedList.Size))
	suite.NotEqual("", truncatedList.NextPageToken, "next page token should not be empty")
	suite.Equal(*firstEntity.Id, *truncatedList.Items[0].Id)

	truncatedList, err = service.GetServingEnvironments(api.ListOptions{
		PageSize:      &pageSize2,
		NextPageToken: &truncatedList.NextPageToken,
	})
	suite.Nilf(err, "error getting ServingEnvironments: %v", err)

	suite.Equal(2, int(truncatedList.Size))
	suite.Equal("", truncatedList.NextPageToken, "next page token should be empty as list item returned")
	suite.Equal(*secondEntity.Id, *truncatedList.Items[0].Id)
	suite.Equal(*thirdEntity.Id, *truncatedList.Items[1].Id)
}

// INFERENCE SERVICE

func (suite *CoreTestSuite) TestCreateInferenceService() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	registeredModelId := suite.registerModel(service, nil, nil)
	runtime := "model-server"
	state := openapi.INFERENCESERVICESTATE_DEPLOYED

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		Runtime:              &runtime,
		State:                &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s: %v", parentResourceId, err)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*createdEntityId, *byId.Contexts[0].Id, "returned id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", parentResourceId, entityName), *byId.Contexts[0].Name, "saved name should match the provided one")
	suite.Equal(entityExternalId2, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(customString, byId.Contexts[0].CustomProperties["custom_string_prop"].GetStringValue(), "saved custom_string_prop custom property should match the provided one")
	suite.Equal(entityDescription, byId.Contexts[0].Properties["description"].GetStringValue(), "saved description should match the provided one")
	suite.Equal(runtime, byId.Contexts[0].Properties["runtime"].GetStringValue(), "saved runtime should match the provided one")
	suite.Equal(string(state), byId.Contexts[0].Properties["state"].GetStringValue(), "saved state should match the provided one")
	suite.Equalf(*inferenceServiceTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *inferenceServiceTypeName)

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(3, len(getAllResp.Contexts), "there should be 3 contexts (RegisteredModel, ServingEnvironment, InferenceService) saved in mlmd")
}

func (suite *CoreTestSuite) TestCreateInferenceServiceFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		ServingEnvironmentId: "9999",
		RegisteredModelId:    "9998",
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	_, err := service.UpsertInferenceService(eut)
	suite.NotNil(err)
	suite.Equal("no serving environment found for id 9999", err.Error())

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	eut.ServingEnvironmentId = parentResourceId

	_, err = service.UpsertInferenceService(eut)
	suite.NotNil(err)
	suite.Equal("no registered model found for id 9998", err.Error())
}

func (suite *CoreTestSuite) TestUpdateInferenceService() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	newExternalId := "org.my_awesome_entity@v1"
	newScore := 0.95

	createdEntity.ExternalID = &newExternalId
	(*createdEntity.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: &openapi.MetadataDoubleValue{
			DoubleValue: &newScore,
		},
	}

	updatedEntity, err := service.UpsertInferenceService(createdEntity)
	suite.Nilf(err, "error updating new entity for %s: %v", registeredModelId, err)

	updateEntityId, _ := converter.StringToInt64(updatedEntity.Id)
	suite.Equal(*createdEntityId, *updateEntityId, "created and updated should have same id")

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be 1 context saved in mlmd by id")

	suite.Equal(*updateEntityId, *byId.Contexts[0].Id, "returned id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", parentResourceId, *eut.Name), *byId.Contexts[0].Name, "saved name should match the provided one")
	suite.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(customString, byId.Contexts[0].CustomProperties["custom_string_prop"].GetStringValue(), "saved custom_string_prop custom property should match the provided one")
	suite.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	suite.Equalf(*inferenceServiceTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *inferenceServiceTypeName)

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(3, len(getAllResp.Contexts), "there should be 3 contexts saved in mlmd")

	// update with nil name
	newExternalId = "org.my_awesome_entity_@v1"
	updatedEntity.ExternalID = &newExternalId
	updatedEntity.Name = nil
	updatedEntity, err = service.UpsertInferenceService(updatedEntity)
	suite.Nilf(err, "error updating new model version for %s: %v", updateEntityId, err)

	updateEntityId, _ = converter.StringToInt64(updatedEntity.Id)
	suite.Equal(*createdEntityId, *updateEntityId, "created and updated should have same id")

	byId, err = suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be 1 context saved in mlmd by id")

	suite.Equal(*updateEntityId, *byId.Contexts[0].Id, "returned id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", parentResourceId, *eut.Name), *byId.Contexts[0].Name, "saved name should match the provided one")
	suite.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(customString, byId.Contexts[0].CustomProperties["custom_string_prop"].GetStringValue(), "saved custom_string_prop custom property should match the provided one")
	suite.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	suite.Equalf(*inferenceServiceTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *inferenceServiceTypeName)

	// update with empty registeredModelId
	newExternalId = "org.my_awesome_entity_@v1"
	prevRegModelId := updatedEntity.RegisteredModelId
	updatedEntity.RegisteredModelId = ""
	updatedEntity, err = service.UpsertInferenceService(updatedEntity)
	suite.Nil(err)
	suite.Equal(prevRegModelId, updatedEntity.RegisteredModelId)
}

func (suite *CoreTestSuite) TestUpdateInferenceServiceFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	newExternalId := "org.my_awesome_entity@v1"
	newScore := 0.95

	createdEntity.ExternalID = &newExternalId
	(*createdEntity.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: &openapi.MetadataDoubleValue{
			DoubleValue: &newScore,
		},
	}

	wrongId := "9999"
	createdEntity.Id = &wrongId
	_, err = service.UpsertInferenceService(createdEntity)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no InferenceService found for id %s", wrongId), err.Error())
}

func (suite *CoreTestSuite) TestGetInferenceServiceById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	state := openapi.INFERENCESERVICESTATE_UNDEPLOYED
	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		State:                &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getById, err := service.GetInferenceServiceById(*createdEntity.Id)
	suite.Nilf(err, "error getting model version with id %d", *createdEntityId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*getById.Id, *converter.Int64ToString(ctx.Id), "returned id should match the mlmd context one")
	suite.Equal(*eut.Name, *getById.Name, "saved name should match the provided one")
	suite.Equal(*eut.ExternalID, *getById.ExternalID, "saved external id should match the provided one")
	suite.Equal(*eut.State, *getById.State, "saved state should match the provided one")
	suite.Equal(*(*getById.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, customString, "saved custom_string_prop custom property should match the provided one")
}

func (suite *CoreTestSuite) TestGetRegisteredModelByInferenceServiceId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}
	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)
	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getRM, err := service.GetRegisteredModelByInferenceService(*createdEntity.Id)
	suite.Nilf(err, "error getting using id %d", *createdEntityId)

	suite.Equal(registeredModelId, *getRM.Id, "returned id should match the original registeredModelId")
}

func (suite *CoreTestSuite) TestGetModelVersionByInferenceServiceId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion1Name := "v1"
	modelVersion1 := &openapi.ModelVersion{Name: &modelVersion1Name, Description: &modelVersionDescription}
	createdVersion1, err := service.UpsertModelVersion(modelVersion1, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion1Id := *createdVersion1.Id

	modelVersion2Name := "v2"
	modelVersion2 := &openapi.ModelVersion{Name: &modelVersion2Name, Description: &modelVersionDescription}
	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion2Id := *createdVersion2.Id
	// end of data preparation

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		ModelVersionId:       nil, // first we test by unspecified
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}
	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getVModel, err := service.GetModelVersionByInferenceService(*createdEntity.Id)
	suite.Nilf(err, "error getting using id %d", *createdEntityId)
	suite.Equal(createdVersion2Id, *getVModel.Id, "returned id shall be the latest ModelVersion by creation order")

	// here we used the returned entity (so ID is populated), and we update to specify the "ID of the ModelVersion to serve"
	createdEntity.ModelVersionId = &createdVersion1Id
	_, err = service.UpsertInferenceService(createdEntity)
	suite.Nilf(err, "error updating eut for %v", parentResourceId)

	getVModel, err = service.GetModelVersionByInferenceService(*createdEntity.Id)
	suite.Nilf(err, "error getting using id %d", *createdEntityId)
	suite.Equal(createdVersion1Id, *getVModel.Id, "returned id shall be the specified one")
}

func (suite *CoreTestSuite) TestGetInferenceServiceByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)

	_, err := service.GetInferenceServiceByParams(of("not-present"), &parentResourceId, nil)
	suite.NotNil(err)
	suite.Equal("no inference services found for name=not-present, parentResourceId=1, externalId=", err.Error())
}

func (suite *CoreTestSuite) TestGetInferenceServiceByParamsName() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getByName, err := service.GetInferenceServiceByParams(&entityName, &parentResourceId, nil)
	suite.Nilf(err, "error getting model version by name %d", *createdEntityId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getByName.Id, "returned id should match the mlmd context one")
	suite.Equal(fmt.Sprintf("%s:%s", parentResourceId, *getByName.Name), *ctx.Name, "saved name should match the provided one")
	suite.Equal(*ctx.ExternalId, *getByName.ExternalID, "saved external id should match the provided one")
	suite.Equal(ctx.CustomProperties["custom_string_prop"].GetStringValue(), *(*getByName.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, "saved custom_string_prop custom property should match the provided one")
}

func (suite *CoreTestSuite) TestGetInfernenceServiceByParamsExternalId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getByExternalId, err := service.GetInferenceServiceByParams(nil, nil, eut.ExternalID)
	suite.Nilf(err, "error getting by external id %d", *eut.ExternalID)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getByExternalId.Id, "returned id should match the mlmd context one")
	suite.Equal(fmt.Sprintf("%s:%s", parentResourceId, *getByExternalId.Name), *ctx.Name, "saved name should match the provided one")
	suite.Equal(*ctx.ExternalId, *getByExternalId.ExternalID, "saved external id should match the provided one")
	suite.Equal(ctx.CustomProperties["custom_string_prop"].GetStringValue(), *(*getByExternalId.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, "saved custom_string_prop custom property should match the provided one")
}

func (suite *CoreTestSuite) TestGetInferenceServiceByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	_, err = service.GetInferenceServiceByParams(nil, nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either (name and parentResourceId), or externalId", err.Error())
}

func (suite *CoreTestSuite) TestGetInferenceServices() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, nil, nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut1 := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
	}

	secondName := "v2"
	secondExtId := "org.myawesomeentity@v2"
	eut2 := &openapi.InferenceService{
		Name:                 &secondName,
		ExternalID:           &secondExtId,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
	}

	thirdName := "v3"
	thirdExtId := "org.myawesomeentity@v3"
	eut3 := &openapi.InferenceService{
		Name:                 &thirdName,
		ExternalID:           &thirdExtId,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
	}

	createdEntity1, err := service.UpsertInferenceService(eut1)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	createdEntity2, err := service.UpsertInferenceService(eut2)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	createdEntity3, err := service.UpsertInferenceService(eut3)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	anotherParentResourceName := "AnotherModel"
	anotherParentResourceExtId := "org.another"
	anotherParentResourceId := suite.registerServingEnvironment(service, &anotherParentResourceName, &anotherParentResourceExtId)

	anotherName := "v1.0"
	anotherExtId := "org.another@v1.0"
	eutAnother := &openapi.InferenceService{
		Name:                 &anotherName,
		ExternalID:           &anotherExtId,
		ServingEnvironmentId: anotherParentResourceId,
		RegisteredModelId:    registeredModelId,
	}

	_, err = service.UpsertInferenceService(eutAnother)
	suite.Nilf(err, "error creating new model version for %d", anotherParentResourceId)

	createdId1, _ := converter.StringToInt64(createdEntity1.Id)
	createdId2, _ := converter.StringToInt64(createdEntity2.Id)
	createdId3, _ := converter.StringToInt64(createdEntity3.Id)

	getAll, err := service.GetInferenceServices(api.ListOptions{}, nil)
	suite.Nilf(err, "error getting all")
	suite.Equal(int32(4), getAll.Size, "expected 4 across all parent resources")

	getAllByParentResource, err := service.GetInferenceServices(api.ListOptions{}, &parentResourceId)
	suite.Nilf(err, "error getting all")
	suite.Equalf(int32(3), getAllByParentResource.Size, "expected 3 for parent resource %d", parentResourceId)

	suite.Equal(*converter.Int64ToString(createdId1), *getAllByParentResource.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdId2), *getAllByParentResource.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdId3), *getAllByParentResource.Items[2].Id)

	// order by last update time, expecting last created as first
	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByParentResource, err = service.GetInferenceServices(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &parentResourceId)
	suite.Nilf(err, "error getting all")
	suite.Equalf(int32(3), getAllByParentResource.Size, "expected 3 for parent resource %d", parentResourceId)

	suite.Equal(*converter.Int64ToString(createdId1), *getAllByParentResource.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdId2), *getAllByParentResource.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdId3), *getAllByParentResource.Items[0].Id)

	// update the second entity
	newExternalId := "updated.org:v2"
	createdEntity2.ExternalID = &newExternalId
	createdEntity2, err = service.UpsertInferenceService(createdEntity2)
	suite.Nilf(err, "error creating new eut2 for %d", parentResourceId)

	suite.Equal(newExternalId, *createdEntity2.ExternalID)

	getAllByParentResource, err = service.GetInferenceServices(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &parentResourceId)
	suite.Nilf(err, "error getting all")
	suite.Equalf(int32(3), getAllByParentResource.Size, "expected 3 for parent resource %d", parentResourceId)

	suite.Equal(*converter.Int64ToString(createdId1), *getAllByParentResource.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdId2), *getAllByParentResource.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdId3), *getAllByParentResource.Items[1].Id)
}

// SERVE MODEL

func (suite *CoreTestSuite) TestCreateServeModel() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersionId := *createdVersion.Id
	createdVersionIdAsInt, _ := converter.StringToInt64(&createdVersionId)
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalID:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	suite.NotNil(createdEntity.Id, "created id should not be nil")

	state, _ := openapi.NewExecutionStateFromValue(executionState)
	suite.Equal(entityName, *createdEntity.Name)
	suite.Equal(*state, *createdEntity.LastKnownState)
	suite.Equal(createdVersionId, createdEntity.ModelVersionId)
	suite.Equal(entityDescription, *createdEntity.Description)
	suite.Equal(customString, *(*createdEntity.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)
	getById, err := suite.mlmdClient.GetExecutionsByID(context.Background(), &proto.GetExecutionsByIDRequest{
		ExecutionIds: []int64{*createdEntityId},
	})
	suite.Nilf(err, "error getting Execution by id %d", createdEntityId)

	suite.Equal(*createdEntityId, *getById.Executions[0].Id)
	suite.Equal(fmt.Sprintf("%s:%s", inferenceServiceId, *createdEntity.Name), *getById.Executions[0].Name)
	suite.Equal(string(*createdEntity.LastKnownState), getById.Executions[0].LastKnownState.String())
	suite.Equal(*createdVersionIdAsInt, getById.Executions[0].Properties["model_version_id"].GetIntValue())
	suite.Equal(*createdEntity.Description, getById.Executions[0].Properties["description"].GetStringValue())
	suite.Equal(*(*createdEntity.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Executions[0].CustomProperties["custom_string_prop"].GetStringValue())

	inferenceServiceIdAsInt, _ := converter.StringToInt64(&inferenceServiceId)
	byCtx, _ := suite.mlmdClient.GetExecutionsByContext(context.Background(), &proto.GetExecutionsByContextRequest{
		ContextId: (*int64)(inferenceServiceIdAsInt),
	})
	suite.Equal(1, len(byCtx.Executions))
	suite.Equal(*createdEntityId, *byCtx.Executions[0].Id)
}

func (suite *CoreTestSuite) TestCreateServeModelFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalID:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: "9998",
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	_, err := service.UpsertServeModel(eut, nil)
	suite.NotNil(err)
	suite.Equal("missing parentResourceId, cannot create ServeModel without parent resource InferenceService", err.Error())

	_, err = service.UpsertServeModel(eut, &inferenceServiceId)
	suite.NotNil(err)
	suite.Equal("no model version found for id 9998", err.Error())
}

func (suite *CoreTestSuite) TestUpdateServeModel() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersionId := *createdVersion.Id
	createdVersionIdAsInt, _ := converter.StringToInt64(&createdVersionId)
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalID:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)

	newState := "UNKNOWN"
	createdEntity.LastKnownState = (*openapi.ExecutionState)(&newState)
	updatedEntity, err := service.UpsertServeModel(createdEntity, &inferenceServiceId)
	suite.Nilf(err, "error updating entity for %d: %v", inferenceServiceId, err)

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)
	updatedEntityId, _ := converter.StringToInt64(updatedEntity.Id)
	suite.Equal(createdEntityId, updatedEntityId)

	getById, err := suite.mlmdClient.GetExecutionsByID(context.Background(), &proto.GetExecutionsByIDRequest{
		ExecutionIds: []int64{*createdEntityId},
	})
	suite.Nilf(err, "error getting by id %d", createdEntityId)

	suite.Equal(*createdEntityId, *getById.Executions[0].Id)
	suite.Equal(fmt.Sprintf("%s:%s", inferenceServiceId, *createdEntity.Name), *getById.Executions[0].Name)
	suite.Equal(string(newState), getById.Executions[0].LastKnownState.String())
	suite.Equal(*createdVersionIdAsInt, getById.Executions[0].Properties["model_version_id"].GetIntValue())
	suite.Equal(*(*createdEntity.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Executions[0].CustomProperties["custom_string_prop"].GetStringValue())

	prevModelVersionId := updatedEntity.ModelVersionId
	updatedEntity.ModelVersionId = ""
	updatedEntity, err = service.UpsertServeModel(updatedEntity, &inferenceServiceId)
	suite.Nilf(err, "error updating entity for %d: %v", inferenceServiceId, err)
	suite.Equal(prevModelVersionId, updatedEntity.ModelVersionId)
}

func (suite *CoreTestSuite) TestUpdateServeModelFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersionId := *createdVersion.Id
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalID:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	suite.NotNil(createdEntity.Id, "created id should not be nil")

	newState := "UNKNOWN"
	createdEntity.LastKnownState = (*openapi.ExecutionState)(&newState)
	updatedEntity, err := service.UpsertServeModel(createdEntity, &inferenceServiceId)
	suite.Nilf(err, "error updating entity for %d: %v", inferenceServiceId, err)

	wrongId := "9998"
	updatedEntity.Id = &wrongId
	_, err = service.UpsertServeModel(updatedEntity, &inferenceServiceId)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no ServeModel found for id %s", wrongId), err.Error())
}

func (suite *CoreTestSuite) TestGetServeModelById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersionId := *createdVersion.Id
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalID:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)

	getById, err := service.GetServeModelById(*createdEntity.Id)
	suite.Nilf(err, "error getting entity by id %d", *createdEntity.Id)

	state, _ := openapi.NewExecutionStateFromValue(executionState)
	suite.NotNil(createdEntity.Id, "created artifact id should not be nil")
	suite.Equal(entityName, *getById.Name)
	suite.Equal(*state, *getById.LastKnownState)
	suite.Equal(createdVersionId, getById.ModelVersionId)
	suite.Equal(customString, *(*getById.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdEntity, *getById, "artifacts returned during creation and on get by id should be equal")
}

func (suite *CoreTestSuite) TestGetServeModels() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion1Name := "v1"
	modelVersion1 := &openapi.ModelVersion{Name: &modelVersion1Name, Description: &modelVersionDescription}
	createdVersion1, err := service.UpsertModelVersion(modelVersion1, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion1Id := *createdVersion1.Id

	modelVersion2Name := "v2"
	modelVersion2 := &openapi.ModelVersion{Name: &modelVersion2Name, Description: &modelVersionDescription}
	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion2Id := *createdVersion2.Id

	modelVersion3Name := "v3"
	modelVersion3 := &openapi.ModelVersion{Name: &modelVersion3Name, Description: &modelVersionDescription}
	createdVersion3, err := service.UpsertModelVersion(modelVersion3, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion3Id := *createdVersion3.Id
	// end of data preparation

	eut1Name := "sm1"
	eut1 := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		Description:    &entityDescription,
		Name:           &eut1Name,
		ModelVersionId: createdVersion1Id,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	eut2Name := "sm2"
	eut2 := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		Description:    &entityDescription,
		Name:           &eut2Name,
		ModelVersionId: createdVersion2Id,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	eut3Name := "sm3"
	eut3 := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		Description:    &entityDescription,
		Name:           &eut3Name,
		ModelVersionId: createdVersion3Id,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &customString,
				},
			},
		},
	}

	createdEntity1, err := service.UpsertServeModel(eut1, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	createdEntity2, err := service.UpsertServeModel(eut2, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	createdEntity3, err := service.UpsertServeModel(eut3, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)

	createdEntityId1, _ := converter.StringToInt64(createdEntity1.Id)
	createdEntityId2, _ := converter.StringToInt64(createdEntity2.Id)
	createdEntityId3, _ := converter.StringToInt64(createdEntity3.Id)

	getAll, err := service.GetServeModels(api.ListOptions{}, nil)
	suite.Nilf(err, "error getting all ServeModel")
	suite.Equalf(int32(3), getAll.Size, "expected three ServeModel")

	suite.Equal(*converter.Int64ToString(createdEntityId1), *getAll.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdEntityId2), *getAll.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdEntityId3), *getAll.Items[2].Id)

	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByInferenceService, err := service.GetServeModels(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &inferenceServiceId)
	suite.Nilf(err, "error getting all ServeModels for %d", inferenceServiceId)
	suite.Equalf(int32(3), getAllByInferenceService.Size, "expected three ServeModels for InferenceServiceId %d", inferenceServiceId)

	suite.Equal(*converter.Int64ToString(createdEntityId1), *getAllByInferenceService.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdEntityId2), *getAllByInferenceService.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdEntityId3), *getAllByInferenceService.Items[0].Id)
}
