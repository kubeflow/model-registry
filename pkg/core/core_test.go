package core

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/internal/mlmdtypes"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
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
	modelOwner       string
	modelExternalId  string
	myCustomProp     string
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

// test defaults
var (
	registeredModelTypeName    = apiutils.Of(defaults.RegisteredModelTypeName)
	modelVersionTypeName       = apiutils.Of(defaults.ModelVersionTypeName)
	modelArtifactTypeName      = apiutils.Of(defaults.ModelArtifactTypeName)
	docArtifactTypeName        = apiutils.Of(defaults.DocArtifactTypeName)
	servingEnvironmentTypeName = apiutils.Of(defaults.ServingEnvironmentTypeName)
	inferenceServiceTypeName   = apiutils.Of(defaults.InferenceServiceTypeName)
	serveModelTypeName         = apiutils.Of(defaults.ServeModelTypeName)
	canAddFields               = apiutils.Of(true)
)

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
	modelOwner = "reg model owner"
	modelExternalId = "org.myawesomemodel"
	myCustomProp = "myCustomPropValue"
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

	// sanity check before each test: connect to MLMD directly, and dry-run any of the gRPC (read) operations;
	// on newer Podman might delay in recognising volume mount files for sqlite3 db,
	// hence in case of error "Cannot connect sqlite3 database: unable to open database file" make some retries.
	maxRetries := 3
	var err error
	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err = suite.mlmdClient.GetContextTypes(context.Background(), &proto.GetContextTypesRequest{})
		if err == nil {
			break
		} else if !strings.Contains(err.Error(), "Cannot connect sqlite3 database: unable to open database file") {
			break // err is different than expected
		}
		time.Sleep(1 * time.Second)
	}
	suite.Nilf(err, "error connecting to MLMD and dry-run any of the gRPC operations: %v", err)
}

// after each test
//   - remove the metadata sqlite file used by mlmd, this way mlmd will recreate it
func (suite *CoreTestSuite) AfterTest(suiteName, testName string) {
	if err := testutils.ClearMetadataSqliteDB(); err != nil {
		suite.Error(err)
	}
}

func (suite *CoreTestSuite) setupModelRegistryService() *ModelRegistryService {
	mlmdtypeNames := mlmdtypes.NewMLMDTypeNamesConfigFromDefaults()
	_, err := mlmdtypes.CreateMLMDTypes(suite.grpcConn, mlmdtypeNames)
	suite.Nilf(err, "error creating MLMD types: %v", err)
	// setup model registry service
	service, err := NewModelRegistryService(suite.grpcConn, mlmdtypeNames)
	suite.Nilf(err, "error creating core service: %v", err)
	mrService, ok := service.(*ModelRegistryService)
	suite.True(ok)
	return mrService
}

// utility function that register a new simple model and return its ID
func (suite *CoreTestSuite) registerModel(service api.ModelRegistryApi, overrideModelName *string, overrideExternalId *string) string {
	registeredModel := &openapi.RegisteredModel{
		Name:        modelName,
		ExternalId:  &modelExternalId,
		Description: &modelDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	if overrideModelName != nil {
		registeredModel.Name = *overrideModelName
	}

	if overrideExternalId != nil {
		registeredModel.ExternalId = overrideExternalId
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	return *createdModel.Id
}

// utility function that register a new simple ServingEnvironment and return its ID
func (suite *CoreTestSuite) registerServingEnvironment(service api.ModelRegistryApi, overrideName string, overrideExternalId *string) string {
	eutName := "Simple ServingEnvironment"
	eutExtID := "Simple ServingEnvironment ExtID"
	eut := &openapi.ServingEnvironment{
		Name:        eutName,
		ExternalId:  &eutExtID,
		Description: &entityDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	if overrideName != "" {
		eut.Name = overrideName
	}

	if overrideExternalId != nil {
		eut.ExternalId = overrideExternalId
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
		Name:        modelVersionName,
		ExternalId:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}

	if overrideVersionName != nil {
		modelVersion.Name = *overrideVersionName
	}

	if overrideVersionExtId != nil {
		modelVersion.ExternalId = overrideVersionExtId
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating model version: %v", err)

	return *createdVersion.Id
}

// utility function that register a new simple ServingEnvironment and return its ID
func (suite *CoreTestSuite) registerInferenceService(service api.ModelRegistryApi, registerdModelId string, overrideParentResourceName string, overrideParentResourceExternalId *string, overrideName *string, overrideExternalId *string) string {
	servingEnvironmentId := suite.registerServingEnvironment(service, overrideParentResourceName, overrideParentResourceExternalId)

	eutName := "simpleInferenceService"
	eutExtID := "simpleInferenceService ExtID"
	eut := &openapi.InferenceService{
		Name:                 &eutName,
		ExternalId:           &eutExtID,
		RegisteredModelId:    registerdModelId,
		ServingEnvironmentId: servingEnvironmentId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	if overrideName != nil {
		eut.Name = overrideName
	}
	if overrideExternalId != nil {
		eut.ExternalId = overrideExternalId
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
	docArtifactReq := proto.PutArtifactTypeRequest{
		CanAddFields: canAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: docArtifactTypeName,
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
	_, err = suite.mlmdClient.PutArtifactType(context.Background(), &docArtifactReq)
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
	docArtifactResp, _ := suite.mlmdClient.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: docArtifactTypeName,
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
	suite.Equal(0, len(docArtifactResp.ArtifactType.Properties))
	suite.Equal(0, len(modelArtifactResp.ArtifactType.Properties))
	suite.Equal(0, len(servingEnvResp.ContextType.Properties))
	suite.Equal(0, len(inferenceServiceResp.ContextType.Properties))
	suite.Equal(0, len(serveModelResp.ExecutionType.Properties))

	// create model registry service
	_ = suite.setupModelRegistryService()

	// assure the types have been correctly setup at startup
	// check NOT empty props
	regModelResp, _ = suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: registeredModelTypeName,
	})
	suite.NotNilf(regModelResp.ContextType, "registered model type %s should exists", *registeredModelTypeName)
	suite.Equal(*registeredModelTypeName, *regModelResp.ContextType.Name)
	suite.Equal(3, len(regModelResp.ContextType.Properties))

	modelVersionResp, _ = suite.mlmdClient.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: modelVersionTypeName,
	})
	suite.NotNilf(modelVersionResp.ContextType, "model version type %s should exists", *modelVersionTypeName)
	suite.Equal(*modelVersionTypeName, *modelVersionResp.ContextType.Name)
	suite.Equal(5, len(modelVersionResp.ContextType.Properties))

	docArtifactResp, _ = suite.mlmdClient.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: docArtifactTypeName,
	})
	suite.NotNilf(docArtifactResp.ArtifactType, "doc artifact type %s should exists", *docArtifactTypeName)
	suite.Equal(*docArtifactTypeName, *docArtifactResp.ArtifactType.Name)
	suite.Equal(1, len(docArtifactResp.ArtifactType.Properties))

	modelArtifactResp, _ = suite.mlmdClient.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: modelArtifactTypeName,
	})
	suite.NotNilf(modelArtifactResp.ArtifactType, "model artifact type %s should exists", *modelArtifactTypeName)
	suite.Equal(*modelArtifactTypeName, *modelArtifactResp.ArtifactType.Name)
	suite.Equal(11, len(modelArtifactResp.ArtifactType.Properties))

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

	docArtifactResp, _ := suite.mlmdClient.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: docArtifactTypeName,
	})
	suite.NotNilf(docArtifactResp.ArtifactType, "doc artifact type %s should exists", *docArtifactTypeName)
	suite.Equal(*docArtifactTypeName, *docArtifactResp.ArtifactType.Name)

	modelArtifactResp, _ := suite.mlmdClient.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: modelArtifactTypeName,
	})
	suite.NotNilf(modelArtifactResp.ArtifactType, "model artifact type %s should exists", *modelArtifactTypeName)
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

	// steps to create model registry service
	_, err = mlmdtypes.CreateMLMDTypes(suite.grpcConn, mlmdtypes.NewMLMDTypeNamesConfigFromDefaults())
	suite.NotNil(err)
	suite.Regexp("error setting up context type "+*registeredModelTypeName+": rpc error: code = AlreadyExists.*", err.Error())
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

	// steps to create model registry service
	_, err = mlmdtypes.CreateMLMDTypes(suite.grpcConn, mlmdtypes.NewMLMDTypeNamesConfigFromDefaults())
	suite.NotNil(err)
	suite.Regexp("error setting up context type "+*modelVersionTypeName+": rpc error: code = AlreadyExists.*", err.Error())
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

	// steps to create model registry service
	_, err = mlmdtypes.CreateMLMDTypes(suite.grpcConn, mlmdtypes.NewMLMDTypeNamesConfigFromDefaults())
	suite.NotNil(err)
	suite.Regexp("error setting up artifact type "+*modelArtifactTypeName+": rpc error: code = AlreadyExists.*", err.Error())
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

	// steps to create model registry service
	_, err = mlmdtypes.CreateMLMDTypes(suite.grpcConn, mlmdtypes.NewMLMDTypeNamesConfigFromDefaults())
	suite.NotNil(err)
	suite.Regexp("error setting up context type "+*servingEnvironmentTypeName+": rpc error: code = AlreadyExists.*", err.Error())
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

	// steps to create model registry service
	_, err = mlmdtypes.CreateMLMDTypes(suite.grpcConn, mlmdtypes.NewMLMDTypeNamesConfigFromDefaults())
	suite.NotNil(err)
	suite.Regexp("error setting up context type "+*inferenceServiceTypeName+": rpc error: code = AlreadyExists.*", err.Error())
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

	// steps to create model registry service
	_, err = mlmdtypes.CreateMLMDTypes(suite.grpcConn, mlmdtypes.NewMLMDTypeNamesConfigFromDefaults())
	suite.NotNil(err)
	suite.Regexp("error setting up execution type "+*serveModelTypeName+": rpc error: code = AlreadyExists.*", err.Error())
}
