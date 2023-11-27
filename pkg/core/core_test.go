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
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// common utility test variables
var (
	// generic
	ascOrderDirection  string
	descOrderDirection string
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

func setup(t *testing.T) (*assert.Assertions, *grpc.ClientConn, proto.MetadataStoreServiceClient, func(t *testing.T)) {
	if testing.Short() {
		// skip these tests using -short param
		t.Skip("skipping testing in short mode")
	}

	// initialize test variable before each test
	ascOrderDirection = "ASC"
	descOrderDirection = "DESC"
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

	conn, client, teardown := testutils.SetupMLMDTestContainer(t)
	return assert.New(t), conn, client, teardown
}

// initialize model registry service and assert no error is thrown
func initModelRegistryService(assertion *assert.Assertions, conn *grpc.ClientConn) api.ModelRegistryApi {
	service, err := NewModelRegistryService(conn)
	assertion.Nilf(err, "error creating core service: %v", err)
	return service
}

// utility function that register a new simple model and return its ID
func registerModel(assertion *assert.Assertions, service api.ModelRegistryApi, overrideModelName *string, overrideExternalId *string) string {
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
	assertion.Nilf(err, "error creating registered model: %v", err)

	return *createdModel.Id
}

// utility function that register a new simple ServingEnvironment and return its ID
func registerServingEnvironment(assertion *assert.Assertions, service api.ModelRegistryApi, overrideName *string, overrideExternalId *string) string {
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
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	return *createdEntity.Id
}

// utility function that register a new simple model and return its ID
func registerModelVersion(
	assertion *assert.Assertions,
	service api.ModelRegistryApi,
	overrideModelName *string,
	overrideExternalId *string,
	overrideVersionName *string,
	overrideVersionExtId *string,
) string {
	registeredModelId := registerModel(assertion, service, overrideModelName, overrideExternalId)

	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	if overrideVersionName != nil {
		modelVersion.Name = overrideVersionName
	}

	if overrideVersionExtId != nil {
		modelVersion.ExternalID = overrideVersionExtId
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating model version: %v", err)

	return *createdVersion.Id
}

// utility function that register a new simple ServingEnvironment and return its ID
func registerInferenceService(assertion *assert.Assertions, service api.ModelRegistryApi, registerdModelId string, overrideParentResourceName *string, overrideParentResourceExternalId *string, overrideName *string, overrideExternalId *string) string {
	servingEnvironmentId := registerServingEnvironment(assertion, service, overrideParentResourceName, overrideParentResourceExternalId)

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
	assertion.Nilf(err, "error creating InferenceService: %v", err)

	return *createdEntity.Id
}

func TestModelRegistryStartupWithExistingEmptyTypes(t *testing.T) {
	ctx := context.Background()
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

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

	_, err := client.PutContextType(context.Background(), &registeredModelReq)
	assertion.Nil(err)
	_, err = client.PutContextType(context.Background(), &modelVersionReq)
	assertion.Nil(err)
	_, err = client.PutArtifactType(context.Background(), &modelArtifactReq)
	assertion.Nil(err)
	_, err = client.PutContextType(context.Background(), &servingEnvironmentReq)
	assertion.Nil(err)
	_, err = client.PutContextType(context.Background(), &inferenceServiceReq)
	assertion.Nil(err)
	_, err = client.PutExecutionType(context.Background(), &serveModelReq)
	assertion.Nil(err)

	// check empty props
	regModelResp, _ := client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: registeredModelTypeName,
	})
	modelVersionResp, _ := client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: modelVersionTypeName,
	})
	modelArtifactResp, _ := client.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: modelArtifactTypeName,
	})
	servingEnvResp, _ := client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: servingEnvironmentTypeName,
	})
	inferenceServiceResp, _ := client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: inferenceServiceTypeName,
	})
	serveModelResp, _ := client.GetExecutionType(ctx, &proto.GetExecutionTypeRequest{
		TypeName: serveModelTypeName,
	})

	assertion.Equal(0, len(regModelResp.ContextType.Properties))
	assertion.Equal(0, len(modelVersionResp.ContextType.Properties))
	assertion.Equal(0, len(modelArtifactResp.ArtifactType.Properties))
	assertion.Equal(0, len(servingEnvResp.ContextType.Properties))
	assertion.Equal(0, len(inferenceServiceResp.ContextType.Properties))
	assertion.Equal(0, len(serveModelResp.ExecutionType.Properties))

	// create model registry service
	_, err = NewModelRegistryService(conn)
	assertion.Nil(err)

	// assure the types have been correctly setup at startup
	// check NOT empty props
	regModelResp, _ = client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: registeredModelTypeName,
	})
	assertion.NotNilf(regModelResp.ContextType, "registered model type %s should exists", *registeredModelTypeName)
	assertion.Equal(*registeredModelTypeName, *regModelResp.ContextType.Name)
	assertion.Equal(2, len(regModelResp.ContextType.Properties))

	modelVersionResp, _ = client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: modelVersionTypeName,
	})
	assertion.NotNilf(modelVersionResp.ContextType, "model version type %s should exists", *modelVersionTypeName)
	assertion.Equal(*modelVersionTypeName, *modelVersionResp.ContextType.Name)
	assertion.Equal(5, len(modelVersionResp.ContextType.Properties))

	modelArtifactResp, _ = client.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: modelArtifactTypeName,
	})
	assertion.NotNilf(modelArtifactResp.ArtifactType, "model artifact type %s should exists", *modelArtifactTypeName)
	assertion.Equal(*modelArtifactTypeName, *modelArtifactResp.ArtifactType.Name)
	assertion.Equal(6, len(modelArtifactResp.ArtifactType.Properties))

	servingEnvResp, _ = client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: servingEnvironmentTypeName,
	})
	assertion.NotNilf(servingEnvResp.ContextType, "serving environment type %s should exists", *servingEnvironmentTypeName)
	assertion.Equal(*servingEnvironmentTypeName, *servingEnvResp.ContextType.Name)
	assertion.Equal(1, len(servingEnvResp.ContextType.Properties))

	inferenceServiceResp, _ = client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: inferenceServiceTypeName,
	})
	assertion.NotNilf(inferenceServiceResp.ContextType, "inference service type %s should exists", *inferenceServiceTypeName)
	assertion.Equal(*inferenceServiceTypeName, *inferenceServiceResp.ContextType.Name)
	assertion.Equal(6, len(inferenceServiceResp.ContextType.Properties))

	serveModelResp, _ = client.GetExecutionType(ctx, &proto.GetExecutionTypeRequest{
		TypeName: serveModelTypeName,
	})
	assertion.NotNilf(serveModelResp.ExecutionType, "serve model type %s should exists", *serveModelTypeName)
	assertion.Equal(*serveModelTypeName, *serveModelResp.ExecutionType.Name)
	assertion.Equal(2, len(serveModelResp.ExecutionType.Properties))
}

func TestModelRegistryTypes(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create model registry service
	_ = initModelRegistryService(assertion, conn)

	// assure the types have been correctly setup at startup
	ctx := context.Background()
	regModelResp, _ := client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: registeredModelTypeName,
	})
	assertion.NotNilf(regModelResp.ContextType, "registered model type %s should exists", *registeredModelTypeName)
	assertion.Equal(*registeredModelTypeName, *regModelResp.ContextType.Name)

	modelVersionResp, _ := client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: modelVersionTypeName,
	})
	assertion.NotNilf(modelVersionResp.ContextType, "model version type %s should exists", *modelVersionTypeName)
	assertion.Equal(*modelVersionTypeName, *modelVersionResp.ContextType.Name)

	modelArtifactResp, _ := client.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: modelArtifactTypeName,
	})
	assertion.NotNilf(modelArtifactResp.ArtifactType, "model version type %s should exists", *modelArtifactTypeName)
	assertion.Equal(*modelArtifactTypeName, *modelArtifactResp.ArtifactType.Name)

	servingEnvResp, _ := client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: servingEnvironmentTypeName,
	})
	assertion.NotNilf(servingEnvResp.ContextType, "serving environment type %s should exists", *servingEnvironmentTypeName)
	assertion.Equal(*servingEnvironmentTypeName, *servingEnvResp.ContextType.Name)

	inferenceServiceResp, _ := client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: inferenceServiceTypeName,
	})
	assertion.NotNilf(inferenceServiceResp.ContextType, "inference service type %s should exists", *inferenceServiceTypeName)
	assertion.Equal(*inferenceServiceTypeName, *inferenceServiceResp.ContextType.Name)

	serveModelResp, _ := client.GetExecutionType(ctx, &proto.GetExecutionTypeRequest{
		TypeName: serveModelTypeName,
	})
	assertion.NotNilf(serveModelResp.ExecutionType, "serve model type %s should exists", *serveModelTypeName)
	assertion.Equal(*serveModelTypeName, *serveModelResp.ExecutionType.Name)
}

func TestModelRegistryFailureForOmittedFieldInRegisteredModel(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	registeredModelReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: registeredModelTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}

	_, err := client.PutContextType(context.Background(), &registeredModelReq)
	assertion.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(conn)
	assertion.NotNil(err)
	assertion.Regexp("error setting up context type odh.RegisteredModel: rpc error: code = AlreadyExists.*", err.Error())
}

func TestModelRegistryFailureForOmittedFieldInModelVersion(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	modelVersionReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: modelVersionTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}

	_, err := client.PutContextType(context.Background(), &modelVersionReq)
	assertion.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(conn)
	assertion.NotNil(err)
	assertion.Regexp("error setting up context type odh.ModelVersion: rpc error: code = AlreadyExists.*", err.Error())
}

func TestModelRegistryFailureForOmittedFieldInModelArtifact(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	modelArtifactReq := proto.PutArtifactTypeRequest{
		CanAddFields: canAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: modelArtifactTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}

	_, err := client.PutArtifactType(context.Background(), &modelArtifactReq)
	assertion.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(conn)
	assertion.NotNil(err)
	assertion.Regexp("error setting up artifact type odh.ModelArtifact: rpc error: code = AlreadyExists.*", err.Error())
}

func TestModelRegistryFailureForOmittedFieldInServingEnvironment(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	servingEnvironmentReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: servingEnvironmentTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}
	_, err := client.PutContextType(context.Background(), &servingEnvironmentReq)
	assertion.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(conn)
	assertion.NotNil(err)
	assertion.Regexp("error setting up context type odh.ServingEnvironment: rpc error: code = AlreadyExists.*", err.Error())
}

func TestModelRegistryFailureForOmittedFieldInInferenceService(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	inferenceServiceReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: inferenceServiceTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}

	_, err := client.PutContextType(context.Background(), &inferenceServiceReq)
	assertion.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(conn)
	assertion.NotNil(err)
	assertion.Regexp("error setting up context type odh.InferenceService: rpc error: code = AlreadyExists.*", err.Error())
}

func TestModelRegistryFailureForOmittedFieldInServeModel(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	serveModelReq := proto.PutExecutionTypeRequest{
		CanAddFields: canAddFields,
		ExecutionType: &proto.ExecutionType{
			Name: serveModelTypeName,
			Properties: map[string]proto.PropertyType{
				"deprecated": proto.PropertyType_STRING,
			},
		},
	}

	_, err := client.PutExecutionType(context.Background(), &serveModelReq)
	assertion.Nil(err)

	// create model registry service
	_, err = NewModelRegistryService(conn)
	assertion.NotNil(err)
	assertion.Regexp("error setting up execution type odh.ServeModel: rpc error: code = AlreadyExists.*", err.Error())
}

// REGISTERED MODELS

func TestCreateRegisteredModel(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

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
	assertion.Nilf(err, "error creating registered model: %v", err)
	assertion.NotNilf(createdModel.Id, "created registered model should not have nil Id")

	createdModelId, _ := converter.StringToInt64(createdModel.Id)
	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdModelId},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	assertion.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	assertion.Equal(modelName, *ctx.Name, "saved model name should match the provided one")
	assertion.Equal(modelExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	assertion.Equal(modelDescription, ctx.Properties["description"].GetStringValue(), "saved description should match the provided one")
	assertion.Equal(string(state), ctx.Properties["state"].GetStringValue(), "saved state should match the provided one")
	assertion.Equal(owner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")
}

func TestUpdateRegisteredModel(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

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
	assertion.Nilf(err, "error creating registered model: %v", err)
	assertion.NotNilf(createdModel.Id, "created registered model should not have nil Id")
	createdModelId, _ := converter.StringToInt64(createdModel.Id)

	// checks created model matches original one except for Id
	assertion.Equal(*registeredModel.Name, *createdModel.Name, "returned model name should match the original one")
	assertion.Equal(*registeredModel.ExternalID, *createdModel.ExternalID, "returned model external id should match the original one")
	assertion.Equal(*registeredModel.CustomProperties, *createdModel.CustomProperties, "returned model custom props should match the original one")

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
	assertion.Nilf(err, "error creating registered model: %v", err)

	// still one registered model
	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdModelId},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	assertion.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	assertion.Equal(modelName, *ctx.Name, "saved model name should match the provided one")
	assertion.Equal(newModelExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	assertion.Equal(newOwner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	// update the model keeping nil name
	newModelExternalId = "newNewExternalId"
	createdModel.ExternalID = &newModelExternalId
	createdModel.Name = nil
	createdModel, err = service.UpsertRegisteredModel(createdModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	// still one registered model
	getAllResp, err = client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err = client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdModelId},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx = ctxById.Contexts[0]
	ctxId = converter.Int64ToString(ctx.Id)
	assertion.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	assertion.Equal(modelName, *ctx.Name, "saved model name should match the provided one")
	assertion.Equal(newModelExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	assertion.Equal(newOwner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")
}

func TestGetRegisteredModelById(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

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
	assertion.Nilf(err, "error creating registered model: %v", err)

	getModelById, err := service.GetRegisteredModelById(*createdModel.Id)
	assertion.Nilf(err, "error getting registered model by id %s: %v", *createdModel.Id, err)

	// checks created model matches original one except for Id
	assertion.Equal(*registeredModel.Name, *getModelById.Name, "saved model name should match the original one")
	assertion.Equal(*registeredModel.ExternalID, *getModelById.ExternalID, "saved model external id should match the original one")
	assertion.Equal(*registeredModel.State, *getModelById.State, "saved model state should match the original one")
	assertion.Equal(*registeredModel.CustomProperties, *getModelById.CustomProperties, "saved model custom props should match the original one")
}

func TestGetRegisteredModelByParamsWithNoResults(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	_, err := service.GetRegisteredModelByParams(of("not-present"), nil)
	assertion.NotNil(err)
	assertion.Equal("no registered models found for name=not-present, externalId=", err.Error())
}

func TestGetRegisteredModelByParamsName(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	byName, err := service.GetRegisteredModelByParams(&modelName, nil)
	assertion.Nilf(err, "error getting registered model by name: %v", err)

	assertion.Equalf(*createdModel.Id, *byName.Id, "the returned model id should match the retrieved by name")
}

func TestGetRegisteredModelByParamsExternalId(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	byName, err := service.GetRegisteredModelByParams(nil, &modelExternalId)
	assertion.Nilf(err, "error getting registered model by external id: %v", err)

	assertion.Equalf(*createdModel.Id, *byName.Id, "the returned model id should match the retrieved by name")
}

func TestGetRegisteredModelByEmptyParams(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	_, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	_, err = service.GetRegisteredModelByParams(nil, nil)
	assertion.NotNil(err)
	assertion.Equal("invalid parameters call, supply either name or externalId", err.Error())
}

func TestGetRegisteredModelsOrderedById(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	orderBy := "ID"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	_, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	_, err = service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	_, err = service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	orderedById, err := service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		assertion.Less(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}

	orderedById, err = service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		assertion.Greater(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}
}

func TestGetRegisteredModelsOrderedByLastUpdate(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	orderBy := "LAST_UPDATE_TIME"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	firstModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	secondModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	thirdModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	// update second model
	secondModel.ExternalID = nil
	_, err = service.UpsertRegisteredModel(secondModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	orderedById, err := service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	assertion.Equal(*firstModel.Id, *orderedById.Items[0].Id)
	assertion.Equal(*thirdModel.Id, *orderedById.Items[1].Id)
	assertion.Equal(*secondModel.Id, *orderedById.Items[2].Id)

	orderedById, err = service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	assertion.Equal(*secondModel.Id, *orderedById.Items[0].Id)
	assertion.Equal(*thirdModel.Id, *orderedById.Items[1].Id)
	assertion.Equal(*firstModel.Id, *orderedById.Items[2].Id)
}

func TestGetRegisteredModelsWithPageSize(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

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
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	secondModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	thirdModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	truncatedList, err := service.GetRegisteredModels(api.ListOptions{
		PageSize: &pageSize,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(1, int(truncatedList.Size))
	assertion.NotEqual("", truncatedList.NextPageToken, "next page token should not be empty")
	assertion.Equal(*firstModel.Id, *truncatedList.Items[0].Id)

	truncatedList, err = service.GetRegisteredModels(api.ListOptions{
		PageSize:      &pageSize2,
		NextPageToken: &truncatedList.NextPageToken,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(2, int(truncatedList.Size))
	assertion.Equal("", truncatedList.NextPageToken, "next page token should be empty as list item returned")
	assertion.Equal(*secondModel.Id, *truncatedList.Items[0].Id)
	assertion.Equal(*thirdModel.Id, *truncatedList.Items[1].Id)
}

// MODEL VERSIONS

func TestCreateModelVersion(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	state := openapi.MODELVERSIONSTATE_LIVE
	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		State:       &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")

	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	byId, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	assertion.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	assertion.Equal(*createdVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	assertion.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	assertion.Equal(versionExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	assertion.Equal(author, byId.Contexts[0].CustomProperties["author"].GetStringValue(), "saved author custom property should match the provided one")
	assertion.Equal(modelVersionDescription, byId.Contexts[0].Properties["description"].GetStringValue(), "saved description should match the provided one")
	assertion.Equal(string(state), byId.Contexts[0].Properties["state"].GetStringValue(), "saved state should match the provided one")
	assertion.Equalf(*modelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *modelVersionTypeName)

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(2, len(getAllResp.Contexts), "there should be two contexts saved in mlmd")
}

func TestCreateModelVersionFailure(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := "9999"

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	_, err := service.UpsertModelVersion(modelVersion, nil)
	assertion.NotNil(err)
	assertion.Equal("missing registered model id, cannot create model version without registered model", err.Error())

	_, err = service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.NotNil(err)
	assertion.Equal("no registered model found for id 9999", err.Error())
}

func TestUpdateModelVersion(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")
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
	assertion.Nilf(err, "error updating new model version for %s: %v", registeredModelId, err)

	updateVersionId, _ := converter.StringToInt64(updatedVersion.Id)
	assertion.Equal(*createdVersionId, *updateVersionId, "created and updated model version should have same id")

	byId, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	assertion.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	assertion.Equal(*updateVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	assertion.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	assertion.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	assertion.Equal(author, byId.Contexts[0].CustomProperties["author"].GetStringValue(), "saved author custom property should match the provided one")
	assertion.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	assertion.Equalf(*modelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *modelVersionTypeName)

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	fmt.Printf("%+v", getAllResp.Contexts)
	assertion.Equal(2, len(getAllResp.Contexts), "there should be two contexts saved in mlmd")

	// update with nil name
	newExternalId = "org.my_awesome_model_@v1"
	updatedVersion.ExternalID = &newExternalId
	updatedVersion.Name = nil
	updatedVersion, err = service.UpsertModelVersion(updatedVersion, &registeredModelId)
	assertion.Nilf(err, "error updating new model version for %s: %v", registeredModelId, err)

	updateVersionId, _ = converter.StringToInt64(updatedVersion.Id)
	assertion.Equal(*createdVersionId, *updateVersionId, "created and updated model version should have same id")

	byId, err = client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	assertion.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	assertion.Equal(*updateVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	assertion.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	assertion.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	assertion.Equal(author, byId.Contexts[0].CustomProperties["author"].GetStringValue(), "saved author custom property should match the provided one")
	assertion.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	assertion.Equalf(*modelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *modelVersionTypeName)
}

func TestUpdateModelVersionFailure(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %s", registeredModelId)
	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")

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
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("no model version found for id %s", wrongId), err.Error())
}

func TestGetModelVersionById(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	state := openapi.MODELVERSIONSTATE_ARCHIVED
	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		State:      &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	getById, err := service.GetModelVersionById(*createdVersion.Id)
	assertion.Nilf(err, "error getting model version with id %d", *createdVersionId)

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	assertion.Equal(*converter.Int64ToString(ctx.Id), *getById.Id, "returned model version id should match the mlmd context one")
	assertion.Equal(*modelVersion.Name, *getById.Name, "saved model name should match the provided one")
	assertion.Equal(*modelVersion.ExternalID, *getById.ExternalID, "saved external id should match the provided one")
	assertion.Equal(*modelVersion.State, *getById.State, "saved model state should match the original one")
	assertion.Equal(author, *(*getById.CustomProperties)["author"].MetadataStringValue.StringValue, "saved author custom property should match the provided one")
}

func TestGetModelVersionByParamsWithNoResults(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	_, err := service.GetModelVersionByParams(of("not-present"), &registeredModelId, nil)
	assertion.NotNil(err)
	assertion.Equal("no model versions found for versionName=not-present, parentResourceId=1, externalId=", err.Error())
}

func TestGetModelVersionByParamsName(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	getByName, err := service.GetModelVersionByParams(&modelVersionName, &registeredModelId, nil)
	assertion.Nilf(err, "error getting model version by name %d", *createdVersionId)

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	assertion.Equal(*converter.Int64ToString(ctx.Id), *getByName.Id, "returned model version id should match the mlmd context one")
	assertion.Equal(fmt.Sprintf("%s:%s", registeredModelId, *getByName.Name), *ctx.Name, "saved model name should match the provided one")
	assertion.Equal(*ctx.ExternalId, *getByName.ExternalID, "saved external id should match the provided one")
	assertion.Equal(ctx.CustomProperties["author"].GetStringValue(), *(*getByName.CustomProperties)["author"].MetadataStringValue.StringValue, "saved author custom property should match the provided one")
}

func TestGetModelVersionByParamsExternalId(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	getByExternalId, err := service.GetModelVersionByParams(nil, nil, modelVersion.ExternalID)
	assertion.Nilf(err, "error getting model version by external id %d", *modelVersion.ExternalID)

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	assertion.Equal(*converter.Int64ToString(ctx.Id), *getByExternalId.Id, "returned model version id should match the mlmd context one")
	assertion.Equal(fmt.Sprintf("%s:%s", registeredModelId, *getByExternalId.Name), *ctx.Name, "saved model name should match the provided one")
	assertion.Equal(*ctx.ExternalId, *getByExternalId.ExternalID, "saved external id should match the provided one")
	assertion.Equal(ctx.CustomProperties["author"].GetStringValue(), *(*getByExternalId.CustomProperties)["author"].MetadataStringValue.StringValue, "saved author custom property should match the provided one")
}

func TestGetModelVersionByEmptyParams(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)
	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")

	_, err = service.GetModelVersionByParams(nil, nil, nil)
	assertion.NotNil(err)
	assertion.Equal("invalid parameters call, supply either (versionName and parentResourceId), or externalId", err.Error())
}

func TestGetModelVersions(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

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
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	createdVersion3, err := service.UpsertModelVersion(modelVersion3, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	anotherRegModelName := "AnotherModel"
	anotherRegModelExtId := "org.another"
	anotherRegisteredModelId := registerModel(assertion, service, &anotherRegModelName, &anotherRegModelExtId)

	anotherModelVersionName := "v1.0"
	anotherModelVersionExtId := "org.another@v1.0"
	modelVersionAnother := &openapi.ModelVersion{
		Name:       &anotherModelVersionName,
		ExternalID: &anotherModelVersionExtId,
	}

	_, err = service.UpsertModelVersion(modelVersionAnother, &anotherRegisteredModelId)
	assertion.Nilf(err, "error creating new model version for %d", anotherRegisteredModelId)

	createdVersionId1, _ := converter.StringToInt64(createdVersion1.Id)
	createdVersionId2, _ := converter.StringToInt64(createdVersion2.Id)
	createdVersionId3, _ := converter.StringToInt64(createdVersion3.Id)

	getAll, err := service.GetModelVersions(api.ListOptions{}, nil)
	assertion.Nilf(err, "error getting all model versions")
	assertion.Equal(int32(4), getAll.Size, "expected four model versions across all registered models")

	getAllByRegModel, err := service.GetModelVersions(api.ListOptions{}, &registeredModelId)
	assertion.Nilf(err, "error getting all model versions")
	assertion.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	assertion.Equal(*converter.Int64ToString(createdVersionId1), *getAllByRegModel.Items[0].Id)
	assertion.Equal(*converter.Int64ToString(createdVersionId2), *getAllByRegModel.Items[1].Id)
	assertion.Equal(*converter.Int64ToString(createdVersionId3), *getAllByRegModel.Items[2].Id)

	// order by last update time, expecting last created as first
	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByRegModel, err = service.GetModelVersions(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &registeredModelId)
	assertion.Nilf(err, "error getting all model versions")
	assertion.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	assertion.Equal(*converter.Int64ToString(createdVersionId1), *getAllByRegModel.Items[2].Id)
	assertion.Equal(*converter.Int64ToString(createdVersionId2), *getAllByRegModel.Items[1].Id)
	assertion.Equal(*converter.Int64ToString(createdVersionId3), *getAllByRegModel.Items[0].Id)

	// update the second version
	newVersionExternalId := "updated.org:v2"
	createdVersion2.ExternalID = &newVersionExternalId
	createdVersion2, err = service.UpsertModelVersion(createdVersion2, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.Equal(newVersionExternalId, *createdVersion2.ExternalID)

	getAllByRegModel, err = service.GetModelVersions(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &registeredModelId)
	assertion.Nilf(err, "error getting all model versions")
	assertion.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	assertion.Equal(*converter.Int64ToString(createdVersionId1), *getAllByRegModel.Items[2].Id)
	assertion.Equal(*converter.Int64ToString(createdVersionId2), *getAllByRegModel.Items[0].Id)
	assertion.Equal(*converter.Int64ToString(createdVersionId3), *getAllByRegModel.Items[1].Id)
}

// MODEL ARTIFACTS

func TestCreateModelArtifact(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

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
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	assertion.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	assertion.Equal(artifactName, *createdArtifact.Name)
	assertion.Equal(*state, *createdArtifact.State)
	assertion.Equal(artifactUri, *createdArtifact.Uri)
	assertion.Equal(artifactDescription, *createdArtifact.Description)
	assertion.Equal("onnx", *createdArtifact.ModelFormatName)
	assertion.Equal("1", *createdArtifact.ModelFormatVersion)
	assertion.Equal("aws-connection-models", *createdArtifact.StorageKey)
	assertion.Equal("bucket", *createdArtifact.StoragePath)
	assertion.Equal(author, *(*createdArtifact.CustomProperties)["author"].MetadataStringValue.StringValue)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)
	getById, err := client.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	assertion.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	assertion.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	assertion.Equal(fmt.Sprintf("%s:%s", modelVersionId, *createdArtifact.Name), *getById.Artifacts[0].Name)
	assertion.Equal(string(*createdArtifact.State), getById.Artifacts[0].State.String())
	assertion.Equal(*createdArtifact.Uri, *getById.Artifacts[0].Uri)
	assertion.Equal(*createdArtifact.Description, getById.Artifacts[0].Properties["description"].GetStringValue())
	assertion.Equal(*createdArtifact.ModelFormatName, getById.Artifacts[0].Properties["model_format_name"].GetStringValue())
	assertion.Equal(*createdArtifact.ModelFormatVersion, getById.Artifacts[0].Properties["model_format_version"].GetStringValue())
	assertion.Equal(*createdArtifact.StorageKey, getById.Artifacts[0].Properties["storage_key"].GetStringValue())
	assertion.Equal(*createdArtifact.StoragePath, getById.Artifacts[0].Properties["storage_path"].GetStringValue())
	assertion.Equal(*(*createdArtifact.CustomProperties)["author"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["author"].GetStringValue())

	modelVersionIdAsInt, _ := converter.StringToInt64(&modelVersionId)
	byCtx, _ := client.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
		ContextId: (*int64)(modelVersionIdAsInt),
	})
	assertion.Equal(1, len(byCtx.Artifacts))
	assertion.Equal(*createdArtifactId, *byCtx.Artifacts[0].Id)
}

func TestCreateModelArtifactFailure(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := "9998"

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	_, err := service.UpsertModelArtifact(modelArtifact, nil)
	assertion.NotNil(err)
	assertion.Equal("missing model version id, cannot create model artifact without model version", err.Error())

	_, err = service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.NotNil(err)
	assertion.Equal("no model version found for id 9998", err.Error())
}

func TestUpdateModelArtifact(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	newState := "MARKED_FOR_DELETION"
	createdArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertModelArtifact(createdArtifact, &modelVersionId)
	assertion.Nilf(err, "error updating model artifact for %d: %v", modelVersionId, err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)
	updatedArtifactId, _ := converter.StringToInt64(updatedArtifact.Id)
	assertion.Equal(createdArtifactId, updatedArtifactId)

	getById, err := client.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	assertion.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	assertion.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	assertion.Equal(fmt.Sprintf("%s:%s", modelVersionId, *createdArtifact.Name), *getById.Artifacts[0].Name)
	assertion.Equal(string(newState), getById.Artifacts[0].State.String())
	assertion.Equal(*createdArtifact.Uri, *getById.Artifacts[0].Uri)
	assertion.Equal(*(*createdArtifact.CustomProperties)["author"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["author"].GetStringValue())
}

func TestUpdateModelArtifactFailure(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for model version %s", modelVersionId)
	assertion.NotNilf(createdArtifact.Id, "created model artifact should not have nil Id")

	newState := "MARKED_FOR_DELETION"
	createdArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertModelArtifact(createdArtifact, &modelVersionId)
	assertion.Nilf(err, "error updating model artifact for %d: %v", modelVersionId, err)

	wrongId := "9998"
	updatedArtifact.Id = &wrongId
	_, err = service.UpsertModelArtifact(updatedArtifact, &modelVersionId)
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("no model artifact found for id %s", wrongId), err.Error())
}

func TestGetModelArtifactById(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)

	getById, err := service.GetModelArtifactById(*createdArtifact.Id)
	assertion.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	assertion.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	assertion.Equal(artifactName, *getById.Name)
	assertion.Equal(*state, *getById.State)
	assertion.Equal(artifactUri, *getById.Uri)
	assertion.Equal(author, *(*getById.CustomProperties)["author"].MetadataStringValue.StringValue)

	assertion.Equal(*createdArtifact, *getById, "artifacts returned during creation and on get by id should be equal")
}

func TestGetModelArtifactByParams(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalID: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)

	getByName, err := service.GetModelArtifactByParams(&artifactName, &modelVersionId, nil)
	assertion.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	assertion.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	assertion.Equal(artifactName, *getByName.Name)
	assertion.Equal(artifactExtId, *getByName.ExternalID)
	assertion.Equal(*state, *getByName.State)
	assertion.Equal(artifactUri, *getByName.Uri)
	assertion.Equal(author, *(*getByName.CustomProperties)["author"].MetadataStringValue.StringValue)

	assertion.Equal(*createdArtifact, *getByName, "artifacts returned during creation and on get by name should be equal")

	getByExtId, err := service.GetModelArtifactByParams(nil, nil, &artifactExtId)
	assertion.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	assertion.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	assertion.Equal(artifactName, *getByExtId.Name)
	assertion.Equal(artifactExtId, *getByExtId.ExternalID)
	assertion.Equal(*state, *getByExtId.State)
	assertion.Equal(artifactUri, *getByExtId.Uri)
	assertion.Equal(author, *(*getByExtId.CustomProperties)["author"].MetadataStringValue.StringValue)

	assertion.Equal(*createdArtifact, *getByExtId, "artifacts returned during creation and on get by ext id should be equal")
}

func TestGetModelArtifactByEmptyParams(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalID: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	_, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	_, err = service.GetModelArtifactByParams(nil, nil, nil)
	assertion.NotNil(err)
	assertion.Equal("invalid parameters call, supply either (artifactName and parentResourceId), or externalId", err.Error())
}

func TestGetModelArtifactByParamsWithNoResults(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	_, err := service.GetModelArtifactByParams(of("not-present"), &modelVersionId, nil)
	assertion.NotNil(err)
	assertion.Equal("no model artifacts found for artifactName=not-present, parentResourceId=2, externalId=", err.Error())
}

func TestGetModelArtifacts(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact1 := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalID: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
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
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
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
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact1, err := service.UpsertModelArtifact(modelArtifact1, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)
	createdArtifact2, err := service.UpsertModelArtifact(modelArtifact2, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)
	createdArtifact3, err := service.UpsertModelArtifact(modelArtifact3, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId1, _ := converter.StringToInt64(createdArtifact1.Id)
	createdArtifactId2, _ := converter.StringToInt64(createdArtifact2.Id)
	createdArtifactId3, _ := converter.StringToInt64(createdArtifact3.Id)

	getAll, err := service.GetModelArtifacts(api.ListOptions{}, nil)
	assertion.Nilf(err, "error getting all model artifacts")
	assertion.Equalf(int32(3), getAll.Size, "expected three model artifacts")

	assertion.Equal(*converter.Int64ToString(createdArtifactId1), *getAll.Items[0].Id)
	assertion.Equal(*converter.Int64ToString(createdArtifactId2), *getAll.Items[1].Id)
	assertion.Equal(*converter.Int64ToString(createdArtifactId3), *getAll.Items[2].Id)

	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByModelVersion, err := service.GetModelArtifacts(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &modelVersionId)
	assertion.Nilf(err, "error getting all model artifacts for %d", modelVersionId)
	assertion.Equalf(int32(3), getAllByModelVersion.Size, "expected three model artifacts for model version %d", modelVersionId)

	assertion.Equal(*converter.Int64ToString(createdArtifactId1), *getAllByModelVersion.Items[2].Id)
	assertion.Equal(*converter.Int64ToString(createdArtifactId2), *getAllByModelVersion.Items[1].Id)
	assertion.Equal(*converter.Int64ToString(createdArtifactId3), *getAllByModelVersion.Items[0].Id)
}

// SERVING ENVIRONMENT

func TestCreateServingEnvironment(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

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
	assertion.Nilf(err, "error creating uut: %v", err)
	assertion.NotNilf(createdEntity.Id, "created uut should not have nil Id")

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)
	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdEntityId},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	assertion.Equal(*createdEntity.Id, *ctxId, "returned id should match the mlmd one")
	assertion.Equal(entityName, *ctx.Name, "saved name should match the provided one")
	assertion.Equal(entityExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	assertion.Equal(entityDescription, ctx.Properties["description"].GetStringValue(), "saved description should match the provided one")
	assertion.Equal(owner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")
}

func TestUpdateServingEnvironment(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

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
	assertion.Nilf(err, "error creating uut: %v", err)
	assertion.NotNilf(createdEntity.Id, "created uut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	// checks created entity matches original one except for Id
	assertion.Equal(*eut.Name, *createdEntity.Name, "returned entity should match the original one")
	assertion.Equal(*eut.ExternalID, *createdEntity.ExternalID, "returned entity external id should match the original one")
	assertion.Equal(*eut.CustomProperties, *createdEntity.CustomProperties, "returned entity custom props should match the original one")

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
	assertion.Nilf(err, "error creating uut: %v", err)

	// still one expected MLMD type
	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdEntityId},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	assertion.Equal(*createdEntity.Id, *ctxId, "returned entity id should match the mlmd one")
	assertion.Equal(entityName, *ctx.Name, "saved entity name should match the provided one")
	assertion.Equal(newExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	assertion.Equal(newOwner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	// update the entity under test, keeping nil name
	newExternalId = "newNewExternalId"
	createdEntity.ExternalID = &newExternalId
	createdEntity.Name = nil
	createdEntity, err = service.UpsertServingEnvironment(createdEntity)
	assertion.Nilf(err, "error creating entity: %v", err)

	// still one registered entity
	getAllResp, err = client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err = client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdEntityId},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx = ctxById.Contexts[0]
	ctxId = converter.Int64ToString(ctx.Id)
	assertion.Equal(*createdEntity.Id, *ctxId, "returned entity id should match the mlmd one")
	assertion.Equal(entityName, *ctx.Name, "saved entity name should match the provided one")
	assertion.Equal(newExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	assertion.Equal(newOwner, ctx.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")
}

func TestGetServingEnvironmentById(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

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
	assertion.Nilf(err, "error creating eut: %v", err)

	getEntityById, err := service.GetServingEnvironmentById(*createdEntity.Id)
	assertion.Nilf(err, "error getting eut by id %s: %v", *createdEntity.Id, err)

	// checks created entity matches original one except for Id
	assertion.Equal(*eut.Name, *getEntityById.Name, "saved name should match the original one")
	assertion.Equal(*eut.ExternalID, *getEntityById.ExternalID, "saved external id should match the original one")
	assertion.Equal(*eut.CustomProperties, *getEntityById.CustomProperties, "saved custom props should match the original one")
}

func TestGetServingEnvironmentByParamsWithNoResults(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	_, err := service.GetServingEnvironmentByParams(of("not-present"), nil)
	assertion.NotNil(err)
	assertion.Equal("no serving environments found for name=not-present, externalId=", err.Error())
}

func TestGetServingEnvironmentByParamsName(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	createdEntity, err := service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	byName, err := service.GetServingEnvironmentByParams(&entityName, nil)
	assertion.Nilf(err, "error getting ServingEnvironment by name: %v", err)

	assertion.Equalf(*createdEntity.Id, *byName.Id, "the returned entity id should match the retrieved by name")
}

func TestGetServingEnvironmentByParamsExternalId(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	createdEntity, err := service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	byName, err := service.GetServingEnvironmentByParams(nil, &entityExternalId)
	assertion.Nilf(err, "error getting ServingEnvironment by external id: %v", err)

	assertion.Equalf(*createdEntity.Id, *byName.Id, "the returned entity id should match the retrieved by name")
}

func TestGetServingEnvironmentByEmptyParams(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	_, err := service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	_, err = service.GetServingEnvironmentByParams(nil, nil)
	assertion.NotNil(err)
	assertion.Equal("invalid parameters call, supply either name or externalId", err.Error())
}

func TestGetServingEnvironmentsOrderedById(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	orderBy := "ID"

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	_, err := service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName := "Pricingentity2"
	newExternalId := "myExternalId2"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	_, err = service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName = "Pricingentity3"
	newExternalId = "myExternalId3"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	_, err = service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	orderedById, err := service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	assertion.Nilf(err, "error getting ServingEnvironment: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		assertion.Less(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}

	orderedById, err = service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	assertion.Nilf(err, "error getting ServingEnvironments: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		assertion.Greater(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}
}

func TestGetServingEnvironmentsOrderedByLastUpdate(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	orderBy := "LAST_UPDATE_TIME"

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalID: &entityExternalId,
	}

	firstEntity, err := service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName := "Pricingentity2"
	newExternalId := "myExternalId2"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	secondEntity, err := service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName = "Pricingentity3"
	newExternalId = "myExternalId3"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	thirdEntity, err := service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	// update second entity
	secondEntity.ExternalID = nil
	_, err = service.UpsertServingEnvironment(secondEntity)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	orderedById, err := service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	assertion.Nilf(err, "error getting ServingEnvironments: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	assertion.Equal(*firstEntity.Id, *orderedById.Items[0].Id)
	assertion.Equal(*thirdEntity.Id, *orderedById.Items[1].Id)
	assertion.Equal(*secondEntity.Id, *orderedById.Items[2].Id)

	orderedById, err = service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	assertion.Nilf(err, "error getting ServingEnvironments: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	assertion.Equal(*secondEntity.Id, *orderedById.Items[0].Id)
	assertion.Equal(*thirdEntity.Id, *orderedById.Items[1].Id)
	assertion.Equal(*firstEntity.Id, *orderedById.Items[2].Id)
}

func TestGetServingEnvironmentsWithPageSize(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

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
	assertion.Nilf(err, "error creating registered entity: %v", err)

	newName := "Pricingentity2"
	newExternalId := "myExternalId2"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	secondEntity, err := service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName = "Pricingentity3"
	newExternalId = "myExternalId3"
	eut.Name = &newName
	eut.ExternalID = &newExternalId
	thirdEntity, err := service.UpsertServingEnvironment(eut)
	assertion.Nilf(err, "error creating ServingEnvironment: %v", err)

	truncatedList, err := service.GetServingEnvironments(api.ListOptions{
		PageSize: &pageSize,
	})
	assertion.Nilf(err, "error getting ServingEnvironments: %v", err)

	assertion.Equal(1, int(truncatedList.Size))
	assertion.NotEqual("", truncatedList.NextPageToken, "next page token should not be empty")
	assertion.Equal(*firstEntity.Id, *truncatedList.Items[0].Id)

	truncatedList, err = service.GetServingEnvironments(api.ListOptions{
		PageSize:      &pageSize2,
		NextPageToken: &truncatedList.NextPageToken,
	})
	assertion.Nilf(err, "error getting ServingEnvironments: %v", err)

	assertion.Equal(2, int(truncatedList.Size))
	assertion.Equal("", truncatedList.NextPageToken, "next page token should be empty as list item returned")
	assertion.Equal(*secondEntity.Id, *truncatedList.Items[0].Id)
	assertion.Equal(*thirdEntity.Id, *truncatedList.Items[1].Id)
}

// INFERENCE SERVICE

func TestCreateInferenceService(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	registeredModelId := registerModel(assertion, service, nil, nil)
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
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	assertion.Nilf(err, "error creating new eut for %s: %v", parentResourceId, err)

	assertion.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	byId, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	assertion.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	assertion.Equal(*createdEntityId, *byId.Contexts[0].Id, "returned id should match the mlmd one")
	assertion.Equal(fmt.Sprintf("%s:%s", parentResourceId, entityName), *byId.Contexts[0].Name, "saved name should match the provided one")
	assertion.Equal(entityExternalId2, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	assertion.Equal(author, byId.Contexts[0].CustomProperties["author"].GetStringValue(), "saved author custom property should match the provided one")
	assertion.Equal(entityDescription, byId.Contexts[0].Properties["description"].GetStringValue(), "saved description should match the provided one")
	assertion.Equal(runtime, byId.Contexts[0].Properties["runtime"].GetStringValue(), "saved runtime should match the provided one")
	assertion.Equal(string(state), byId.Contexts[0].Properties["state"].GetStringValue(), "saved state should match the provided one")
	assertion.Equalf(*inferenceServiceTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *inferenceServiceTypeName)

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(3, len(getAllResp.Contexts), "there should be 3 contexts (RegisteredModel, ServingEnvironment, InferenceService) saved in mlmd")
}

func TestCreateInferenceServiceFailure(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		ServingEnvironmentId: "9999",
		RegisteredModelId:    "9998",
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	_, err := service.UpsertInferenceService(eut)
	assertion.NotNil(err)
	assertion.Equal("no serving environment found for id 9999", err.Error())

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	eut.ServingEnvironmentId = parentResourceId

	_, err = service.UpsertInferenceService(eut)
	assertion.NotNil(err)
	assertion.Equal("no registered model found for id 9998", err.Error())
}

func TestUpdateInferenceService(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	registeredModelId := registerModel(assertion, service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)

	assertion.NotNilf(createdEntity.Id, "created eut should not have nil Id")

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
	assertion.Nilf(err, "error updating new entity for %s: %v", registeredModelId, err)

	updateEntityId, _ := converter.StringToInt64(updatedEntity.Id)
	assertion.Equal(*createdEntityId, *updateEntityId, "created and updated should have same id")

	byId, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateEntityId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	assertion.Equal(1, len(byId.Contexts), "there should be 1 context saved in mlmd by id")

	assertion.Equal(*updateEntityId, *byId.Contexts[0].Id, "returned id should match the mlmd one")
	assertion.Equal(fmt.Sprintf("%s:%s", parentResourceId, *eut.Name), *byId.Contexts[0].Name, "saved name should match the provided one")
	assertion.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	assertion.Equal(author, byId.Contexts[0].CustomProperties["author"].GetStringValue(), "saved author custom property should match the provided one")
	assertion.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	assertion.Equalf(*inferenceServiceTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *inferenceServiceTypeName)

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(3, len(getAllResp.Contexts), "there should be 3 contexts saved in mlmd")

	// update with nil name
	newExternalId = "org.my_awesome_entity_@v1"
	updatedEntity.ExternalID = &newExternalId
	updatedEntity.Name = nil
	updatedEntity, err = service.UpsertInferenceService(updatedEntity)
	assertion.Nilf(err, "error updating new model version for %s: %v", updateEntityId, err)

	updateEntityId, _ = converter.StringToInt64(updatedEntity.Id)
	assertion.Equal(*createdEntityId, *updateEntityId, "created and updated should have same id")

	byId, err = client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateEntityId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	assertion.Equal(1, len(byId.Contexts), "there should be 1 context saved in mlmd by id")

	assertion.Equal(*updateEntityId, *byId.Contexts[0].Id, "returned id should match the mlmd one")
	assertion.Equal(fmt.Sprintf("%s:%s", parentResourceId, *eut.Name), *byId.Contexts[0].Name, "saved name should match the provided one")
	assertion.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	assertion.Equal(author, byId.Contexts[0].CustomProperties["author"].GetStringValue(), "saved author custom property should match the provided one")
	assertion.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	assertion.Equalf(*inferenceServiceTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *inferenceServiceTypeName)

	// update with empty registeredModelId
	newExternalId = "org.my_awesome_entity_@v1"
	prevRegModelId := updatedEntity.RegisteredModelId
	updatedEntity.RegisteredModelId = ""
	updatedEntity, err = service.UpsertInferenceService(updatedEntity)
	assertion.Nil(err)
	assertion.Equal(prevRegModelId, updatedEntity.RegisteredModelId)
}

func TestUpdateInferenceServiceFailure(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	registeredModelId := registerModel(assertion, service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)

	assertion.NotNilf(createdEntity.Id, "created eut should not have nil Id")

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
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("no InferenceService found for id %s", wrongId), err.Error())
}

func TestGetInferenceServiceById(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	registeredModelId := registerModel(assertion, service, nil, nil)

	state := openapi.INFERENCESERVICESTATE_UNDEPLOYED
	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		State:                &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)

	assertion.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getById, err := service.GetInferenceServiceById(*createdEntity.Id)
	assertion.Nilf(err, "error getting model version with id %d", *createdEntityId)

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	assertion.Nilf(err, "error retrieving context, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	assertion.Equal(*getById.Id, *converter.Int64ToString(ctx.Id), "returned id should match the mlmd context one")
	assertion.Equal(*eut.Name, *getById.Name, "saved name should match the provided one")
	assertion.Equal(*eut.ExternalID, *getById.ExternalID, "saved external id should match the provided one")
	assertion.Equal(*eut.State, *getById.State, "saved state should match the provided one")
	assertion.Equal(*(*getById.CustomProperties)["author"].MetadataStringValue.StringValue, author, "saved author custom property should match the provided one")
}

func TestGetRegisteredModelByInferenceServiceId(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	registeredModelId := registerModel(assertion, service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}
	createdEntity, err := service.UpsertInferenceService(eut)
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)
	assertion.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getRM, err := service.GetRegisteredModelByInferenceService(*createdEntity.Id)
	assertion.Nilf(err, "error getting using id %d", *createdEntityId)

	assertion.Equal(registeredModelId, *getRM.Id, "returned id should match the original registeredModelId")
}

func TestGetModelVersionByInferenceServiceId(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion1Name := "v1"
	modelVersion1 := &openapi.ModelVersion{Name: &modelVersion1Name, Description: &modelVersionDescription}
	createdVersion1, err := service.UpsertModelVersion(modelVersion1, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion1Id := *createdVersion1.Id

	modelVersion2Name := "v2"
	modelVersion2 := &openapi.ModelVersion{Name: &modelVersion2Name, Description: &modelVersionDescription}
	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)
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
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}
	createdEntity, err := service.UpsertInferenceService(eut)
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getVModel, err := service.GetModelVersionByInferenceService(*createdEntity.Id)
	assertion.Nilf(err, "error getting using id %d", *createdEntityId)
	assertion.Equal(createdVersion2Id, *getVModel.Id, "returned id shall be the latest ModelVersion by creation order")

	// here we used the returned entity (so ID is populated), and we update to specify the "ID of the ModelVersion to serve"
	createdEntity.ModelVersionId = &createdVersion1Id
	_, err = service.UpsertInferenceService(createdEntity)
	assertion.Nilf(err, "error updating eut for %v", parentResourceId)

	getVModel, err = service.GetModelVersionByInferenceService(*createdEntity.Id)
	assertion.Nilf(err, "error getting using id %d", *createdEntityId)
	assertion.Equal(createdVersion1Id, *getVModel.Id, "returned id shall be the specified one")
}

func TestGetInferenceServiceByParamsWithNoResults(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)

	_, err := service.GetInferenceServiceByParams(of("not-present"), &parentResourceId, nil)
	assertion.NotNil(err)
	assertion.Equal("no inference services found for name=not-present, parentResourceId=1, externalId=", err.Error())
}

func TestGetInferenceServiceByParamsName(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	registeredModelId := registerModel(assertion, service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)

	assertion.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getByName, err := service.GetInferenceServiceByParams(&entityName, &parentResourceId, nil)
	assertion.Nilf(err, "error getting model version by name %d", *createdEntityId)

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	assertion.Nilf(err, "error retrieving context, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	assertion.Equal(*converter.Int64ToString(ctx.Id), *getByName.Id, "returned id should match the mlmd context one")
	assertion.Equal(fmt.Sprintf("%s:%s", parentResourceId, *getByName.Name), *ctx.Name, "saved name should match the provided one")
	assertion.Equal(*ctx.ExternalId, *getByName.ExternalID, "saved external id should match the provided one")
	assertion.Equal(ctx.CustomProperties["author"].GetStringValue(), *(*getByName.CustomProperties)["author"].MetadataStringValue.StringValue, "saved author custom property should match the provided one")
}

func TestGetInfernenceServiceByParamsExternalId(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	registeredModelId := registerModel(assertion, service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)

	assertion.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getByExternalId, err := service.GetInferenceServiceByParams(nil, nil, eut.ExternalID)
	assertion.Nilf(err, "error getting by external id %d", *eut.ExternalID)

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	assertion.Nilf(err, "error retrieving context, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	assertion.Equal(*converter.Int64ToString(ctx.Id), *getByExternalId.Id, "returned id should match the mlmd context one")
	assertion.Equal(fmt.Sprintf("%s:%s", parentResourceId, *getByExternalId.Name), *ctx.Name, "saved name should match the provided one")
	assertion.Equal(*ctx.ExternalId, *getByExternalId.ExternalID, "saved external id should match the provided one")
	assertion.Equal(ctx.CustomProperties["author"].GetStringValue(), *(*getByExternalId.CustomProperties)["author"].MetadataStringValue.StringValue, "saved author custom property should match the provided one")
}

func TestGetInferenceServiceByEmptyParams(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	registeredModelId := registerModel(assertion, service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalID:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)

	assertion.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	_, err = service.GetInferenceServiceByParams(nil, nil, nil)
	assertion.NotNil(err)
	assertion.Equal("invalid parameters call, supply either (name and parentResourceId), or externalId", err.Error())
}

func TestGetInferenceServices(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	parentResourceId := registerServingEnvironment(assertion, service, nil, nil)
	registeredModelId := registerModel(assertion, service, nil, nil)

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
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)

	createdEntity2, err := service.UpsertInferenceService(eut2)
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)

	createdEntity3, err := service.UpsertInferenceService(eut3)
	assertion.Nilf(err, "error creating new eut for %v", parentResourceId)

	anotherParentResourceName := "AnotherModel"
	anotherParentResourceExtId := "org.another"
	anotherParentResourceId := registerServingEnvironment(assertion, service, &anotherParentResourceName, &anotherParentResourceExtId)

	anotherName := "v1.0"
	anotherExtId := "org.another@v1.0"
	eutAnother := &openapi.InferenceService{
		Name:                 &anotherName,
		ExternalID:           &anotherExtId,
		ServingEnvironmentId: anotherParentResourceId,
		RegisteredModelId:    registeredModelId,
	}

	_, err = service.UpsertInferenceService(eutAnother)
	assertion.Nilf(err, "error creating new model version for %d", anotherParentResourceId)

	createdId1, _ := converter.StringToInt64(createdEntity1.Id)
	createdId2, _ := converter.StringToInt64(createdEntity2.Id)
	createdId3, _ := converter.StringToInt64(createdEntity3.Id)

	getAll, err := service.GetInferenceServices(api.ListOptions{}, nil)
	assertion.Nilf(err, "error getting all")
	assertion.Equal(int32(4), getAll.Size, "expected 4 across all parent resources")

	getAllByParentResource, err := service.GetInferenceServices(api.ListOptions{}, &parentResourceId)
	assertion.Nilf(err, "error getting all")
	assertion.Equalf(int32(3), getAllByParentResource.Size, "expected 3 for parent resource %d", parentResourceId)

	assertion.Equal(*converter.Int64ToString(createdId1), *getAllByParentResource.Items[0].Id)
	assertion.Equal(*converter.Int64ToString(createdId2), *getAllByParentResource.Items[1].Id)
	assertion.Equal(*converter.Int64ToString(createdId3), *getAllByParentResource.Items[2].Id)

	// order by last update time, expecting last created as first
	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByParentResource, err = service.GetInferenceServices(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &parentResourceId)
	assertion.Nilf(err, "error getting all")
	assertion.Equalf(int32(3), getAllByParentResource.Size, "expected 3 for parent resource %d", parentResourceId)

	assertion.Equal(*converter.Int64ToString(createdId1), *getAllByParentResource.Items[2].Id)
	assertion.Equal(*converter.Int64ToString(createdId2), *getAllByParentResource.Items[1].Id)
	assertion.Equal(*converter.Int64ToString(createdId3), *getAllByParentResource.Items[0].Id)

	// update the second entity
	newExternalId := "updated.org:v2"
	createdEntity2.ExternalID = &newExternalId
	createdEntity2, err = service.UpsertInferenceService(createdEntity2)
	assertion.Nilf(err, "error creating new eut2 for %d", parentResourceId)

	assertion.Equal(newExternalId, *createdEntity2.ExternalID)

	getAllByParentResource, err = service.GetInferenceServices(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &parentResourceId)
	assertion.Nilf(err, "error getting all")
	assertion.Equalf(int32(3), getAllByParentResource.Size, "expected 3 for parent resource %d", parentResourceId)

	assertion.Equal(*converter.Int64ToString(createdId1), *getAllByParentResource.Items[2].Id)
	assertion.Equal(*converter.Int64ToString(createdId2), *getAllByParentResource.Items[0].Id)
	assertion.Equal(*converter.Int64ToString(createdId3), *getAllByParentResource.Items[1].Id)
}

// SERVE MODEL

func TestCreateServeModel(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)
	inferenceServiceId := registerInferenceService(assertion, service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)
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
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	assertion.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	assertion.NotNil(createdEntity.Id, "created id should not be nil")

	state, _ := openapi.NewExecutionStateFromValue(executionState)
	assertion.Equal(entityName, *createdEntity.Name)
	assertion.Equal(*state, *createdEntity.LastKnownState)
	assertion.Equal(createdVersionId, createdEntity.ModelVersionId)
	assertion.Equal(entityDescription, *createdEntity.Description)
	assertion.Equal(author, *(*createdEntity.CustomProperties)["author"].MetadataStringValue.StringValue)

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)
	getById, err := client.GetExecutionsByID(context.Background(), &proto.GetExecutionsByIDRequest{
		ExecutionIds: []int64{*createdEntityId},
	})
	assertion.Nilf(err, "error getting Execution by id %d", createdEntityId)

	assertion.Equal(*createdEntityId, *getById.Executions[0].Id)
	assertion.Equal(fmt.Sprintf("%s:%s", inferenceServiceId, *createdEntity.Name), *getById.Executions[0].Name)
	assertion.Equal(string(*createdEntity.LastKnownState), getById.Executions[0].LastKnownState.String())
	assertion.Equal(*createdVersionIdAsInt, getById.Executions[0].Properties["model_version_id"].GetIntValue())
	assertion.Equal(*createdEntity.Description, getById.Executions[0].Properties["description"].GetStringValue())
	assertion.Equal(*(*createdEntity.CustomProperties)["author"].MetadataStringValue.StringValue, getById.Executions[0].CustomProperties["author"].GetStringValue())

	inferenceServiceIdAsInt, _ := converter.StringToInt64(&inferenceServiceId)
	byCtx, _ := client.GetExecutionsByContext(context.Background(), &proto.GetExecutionsByContextRequest{
		ContextId: (*int64)(inferenceServiceIdAsInt),
	})
	assertion.Equal(1, len(byCtx.Executions))
	assertion.Equal(*createdEntityId, *byCtx.Executions[0].Id)
}

func TestCreateServeModelFailure(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)
	inferenceServiceId := registerInferenceService(assertion, service, registeredModelId, nil, nil, nil, nil)
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalID:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: "9998",
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	_, err := service.UpsertServeModel(eut, nil)
	assertion.NotNil(err)
	assertion.Equal("missing parentResourceId, cannot create ServeModel without parent resource InferenceService", err.Error())

	_, err = service.UpsertServeModel(eut, &inferenceServiceId)
	assertion.NotNil(err)
	assertion.Equal("no model version found for id 9998", err.Error())
}

func TestUpdateServeModel(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)
	inferenceServiceId := registerInferenceService(assertion, service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)
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
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	assertion.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)

	newState := "UNKNOWN"
	createdEntity.LastKnownState = (*openapi.ExecutionState)(&newState)
	updatedEntity, err := service.UpsertServeModel(createdEntity, &inferenceServiceId)
	assertion.Nilf(err, "error updating entity for %d: %v", inferenceServiceId, err)

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)
	updatedEntityId, _ := converter.StringToInt64(updatedEntity.Id)
	assertion.Equal(createdEntityId, updatedEntityId)

	getById, err := client.GetExecutionsByID(context.Background(), &proto.GetExecutionsByIDRequest{
		ExecutionIds: []int64{*createdEntityId},
	})
	assertion.Nilf(err, "error getting by id %d", createdEntityId)

	assertion.Equal(*createdEntityId, *getById.Executions[0].Id)
	assertion.Equal(fmt.Sprintf("%s:%s", inferenceServiceId, *createdEntity.Name), *getById.Executions[0].Name)
	assertion.Equal(string(newState), getById.Executions[0].LastKnownState.String())
	assertion.Equal(*createdVersionIdAsInt, getById.Executions[0].Properties["model_version_id"].GetIntValue())
	assertion.Equal(*(*createdEntity.CustomProperties)["author"].MetadataStringValue.StringValue, getById.Executions[0].CustomProperties["author"].GetStringValue())

	prevModelVersionId := updatedEntity.ModelVersionId
	updatedEntity.ModelVersionId = ""
	updatedEntity, err = service.UpsertServeModel(updatedEntity, &inferenceServiceId)
	assertion.Nilf(err, "error updating entity for %d: %v", inferenceServiceId, err)
	assertion.Equal(prevModelVersionId, updatedEntity.ModelVersionId)

}

func TestUpdateServeModelFailure(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)
	inferenceServiceId := registerInferenceService(assertion, service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersionId := *createdVersion.Id
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalID:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	assertion.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	assertion.NotNil(createdEntity.Id, "created id should not be nil")

	newState := "UNKNOWN"
	createdEntity.LastKnownState = (*openapi.ExecutionState)(&newState)
	updatedEntity, err := service.UpsertServeModel(createdEntity, &inferenceServiceId)
	assertion.Nilf(err, "error updating entity for %d: %v", inferenceServiceId, err)

	wrongId := "9998"
	updatedEntity.Id = &wrongId
	_, err = service.UpsertServeModel(updatedEntity, &inferenceServiceId)
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("no ServeModel found for id %s", wrongId), err.Error())
}

func TestGetServeModelById(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)
	inferenceServiceId := registerInferenceService(assertion, service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        &modelVersionName,
		ExternalID:  &versionExternalId,
		Description: &modelVersionDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersionId := *createdVersion.Id
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalID:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	assertion.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)

	getById, err := service.GetServeModelById(*createdEntity.Id)
	assertion.Nilf(err, "error getting entity by id %d", *createdEntity.Id)

	state, _ := openapi.NewExecutionStateFromValue(executionState)
	assertion.NotNil(createdEntity.Id, "created artifact id should not be nil")
	assertion.Equal(entityName, *getById.Name)
	assertion.Equal(*state, *getById.LastKnownState)
	assertion.Equal(createdVersionId, getById.ModelVersionId)
	assertion.Equal(author, *(*getById.CustomProperties)["author"].MetadataStringValue.StringValue)

	assertion.Equal(*createdEntity, *getById, "artifacts returned during creation and on get by id should be equal")
}

func TestGetServeModels(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)
	inferenceServiceId := registerInferenceService(assertion, service, registeredModelId, nil, nil, nil, nil)

	modelVersion1Name := "v1"
	modelVersion1 := &openapi.ModelVersion{Name: &modelVersion1Name, Description: &modelVersionDescription}
	createdVersion1, err := service.UpsertModelVersion(modelVersion1, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion1Id := *createdVersion1.Id

	modelVersion2Name := "v2"
	modelVersion2 := &openapi.ModelVersion{Name: &modelVersion2Name, Description: &modelVersionDescription}
	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion2Id := *createdVersion2.Id

	modelVersion3Name := "v3"
	modelVersion3 := &openapi.ModelVersion{Name: &modelVersion3Name, Description: &modelVersionDescription}
	createdVersion3, err := service.UpsertModelVersion(modelVersion3, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion3Id := *createdVersion3.Id
	// end of data preparation

	eut1Name := "sm1"
	eut1 := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		Description:    &entityDescription,
		Name:           &eut1Name,
		ModelVersionId: createdVersion1Id,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
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
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
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
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdEntity1, err := service.UpsertServeModel(eut1, &inferenceServiceId)
	assertion.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	createdEntity2, err := service.UpsertServeModel(eut2, &inferenceServiceId)
	assertion.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	createdEntity3, err := service.UpsertServeModel(eut3, &inferenceServiceId)
	assertion.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)

	createdEntityId1, _ := converter.StringToInt64(createdEntity1.Id)
	createdEntityId2, _ := converter.StringToInt64(createdEntity2.Id)
	createdEntityId3, _ := converter.StringToInt64(createdEntity3.Id)

	getAll, err := service.GetServeModels(api.ListOptions{}, nil)
	assertion.Nilf(err, "error getting all ServeModel")
	assertion.Equalf(int32(3), getAll.Size, "expected three ServeModel")

	assertion.Equal(*converter.Int64ToString(createdEntityId1), *getAll.Items[0].Id)
	assertion.Equal(*converter.Int64ToString(createdEntityId2), *getAll.Items[1].Id)
	assertion.Equal(*converter.Int64ToString(createdEntityId3), *getAll.Items[2].Id)

	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByInferenceService, err := service.GetServeModels(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &inferenceServiceId)
	assertion.Nilf(err, "error getting all ServeModels for %d", inferenceServiceId)
	assertion.Equalf(int32(3), getAllByInferenceService.Size, "expected three ServeModels for InferenceServiceId %d", inferenceServiceId)

	assertion.Equal(*converter.Int64ToString(createdEntityId1), *getAllByInferenceService.Items[2].Id)
	assertion.Equal(*converter.Int64ToString(createdEntityId2), *getAllByInferenceService.Items[1].Id)
	assertion.Equal(*converter.Int64ToString(createdEntityId3), *getAllByInferenceService.Items[0].Id)
}
