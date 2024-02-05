package mlmdtypes

import (
	"context"
	"fmt"

	"github.com/opendatahub-io/model-registry/internal/apiutils"
	"github.com/opendatahub-io/model-registry/internal/constants"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"google.golang.org/grpc"
)

var (
	registeredModelTypeName    = apiutils.Of(constants.RegisteredModelTypeName)
	modelVersionTypeName       = apiutils.Of(constants.ModelVersionTypeName)
	modelArtifactTypeName      = apiutils.Of(constants.ModelArtifactTypeName)
	docArtifactTypeName        = apiutils.Of(constants.DocArtifactTypeName)
	servingEnvironmentTypeName = apiutils.Of(constants.ServingEnvironmentTypeName)
	inferenceServiceTypeName   = apiutils.Of(constants.InferenceServiceTypeName)
	serveModelTypeName         = apiutils.Of(constants.ServeModelTypeName)
	canAddFields               = apiutils.Of(true)
)

// Utility method that created the necessary Model Registry's logical-model types
// as the necessary MLMD's Context, Artifact, Execution types etc. in the underlying MLMD service
func CreateMLMDTypes(cc grpc.ClientConnInterface) (map[string]int64, error) {
	client := proto.NewMetadataStoreServiceClient(cc)

	registeredModelReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: registeredModelTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
				"state":       proto.PropertyType_STRING,
			},
		},
	}

	modelVersionReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: modelVersionTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
				"model_name":  proto.PropertyType_STRING,
				"version":     proto.PropertyType_STRING,
				"author":      proto.PropertyType_STRING,
				"state":       proto.PropertyType_STRING,
			},
		},
	}

	docArtifactReq := proto.PutArtifactTypeRequest{
		CanAddFields: canAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: docArtifactTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
			},
		},
	}

	modelArtifactReq := proto.PutArtifactTypeRequest{
		CanAddFields: canAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: modelArtifactTypeName,
			Properties: map[string]proto.PropertyType{
				"description":          proto.PropertyType_STRING,
				"model_format_name":    proto.PropertyType_STRING,
				"model_format_version": proto.PropertyType_STRING,
				"storage_key":          proto.PropertyType_STRING,
				"storage_path":         proto.PropertyType_STRING,
				"service_account_name": proto.PropertyType_STRING,
			},
		},
	}

	servingEnvironmentReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: servingEnvironmentTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
			},
		},
	}

	inferenceServiceReq := proto.PutContextTypeRequest{
		CanAddFields: canAddFields,
		ContextType: &proto.ContextType{
			Name: inferenceServiceTypeName,
			Properties: map[string]proto.PropertyType{
				"description":         proto.PropertyType_STRING,
				"model_version_id":    proto.PropertyType_INT,
				"registered_model_id": proto.PropertyType_INT,
				// same information tracked using ParentContext association
				"serving_environment_id": proto.PropertyType_INT,
				"runtime":                proto.PropertyType_STRING,
				"desired_state":          proto.PropertyType_STRING,
			},
		},
	}

	serveModelReq := proto.PutExecutionTypeRequest{
		CanAddFields: canAddFields,
		ExecutionType: &proto.ExecutionType{
			Name: serveModelTypeName,
			Properties: map[string]proto.PropertyType{
				"description":      proto.PropertyType_STRING,
				"model_version_id": proto.PropertyType_INT,
			},
		},
	}

	registeredModelResp, err := client.PutContextType(context.Background(), &registeredModelReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up context type %s: %v", *registeredModelTypeName, err)
	}

	modelVersionResp, err := client.PutContextType(context.Background(), &modelVersionReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up context type %s: %v", *modelVersionTypeName, err)
	}

	docArtifactResp, err := client.PutArtifactType(context.Background(), &docArtifactReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up artifact type %s: %v", *docArtifactTypeName, err)
	}

	modelArtifactResp, err := client.PutArtifactType(context.Background(), &modelArtifactReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up artifact type %s: %v", *modelArtifactTypeName, err)
	}

	servingEnvironmentResp, err := client.PutContextType(context.Background(), &servingEnvironmentReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up context type %s: %v", *servingEnvironmentTypeName, err)
	}

	inferenceServiceResp, err := client.PutContextType(context.Background(), &inferenceServiceReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up context type %s: %v", *inferenceServiceTypeName, err)
	}

	serveModelResp, err := client.PutExecutionType(context.Background(), &serveModelReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up execution type %s: %v", *serveModelTypeName, err)
	}

	typesMap := map[string]int64{
		constants.RegisteredModelTypeName:    registeredModelResp.GetTypeId(),
		constants.ModelVersionTypeName:       modelVersionResp.GetTypeId(),
		constants.DocArtifactTypeName:        docArtifactResp.GetTypeId(),
		constants.ModelArtifactTypeName:      modelArtifactResp.GetTypeId(),
		constants.ServingEnvironmentTypeName: servingEnvironmentResp.GetTypeId(),
		constants.InferenceServiceTypeName:   inferenceServiceResp.GetTypeId(),
		constants.ServeModelTypeName:         serveModelResp.GetTypeId(),
	}
	return typesMap, nil
}
