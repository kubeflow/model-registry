package mlmdtypes

import (
	"context"
	"fmt"

	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"google.golang.org/grpc"
)

type MLMDTypeNamesConfig struct {
	RegisteredModelTypeName    string
	ModelVersionTypeName       string
	ModelArtifactTypeName      string
	DocArtifactTypeName        string
	ServingEnvironmentTypeName string
	InferenceServiceTypeName   string
	ServeModelTypeName         string
	CanAddFields               bool
}

func NewMLMDTypeNamesConfigFromDefaults() MLMDTypeNamesConfig {
	return MLMDTypeNamesConfig{
		RegisteredModelTypeName:    defaults.RegisteredModelTypeName,
		ModelVersionTypeName:       defaults.ModelVersionTypeName,
		ModelArtifactTypeName:      defaults.ModelArtifactTypeName,
		DocArtifactTypeName:        defaults.DocArtifactTypeName,
		ServingEnvironmentTypeName: defaults.ServingEnvironmentTypeName,
		InferenceServiceTypeName:   defaults.InferenceServiceTypeName,
		ServeModelTypeName:         defaults.ServeModelTypeName,
		CanAddFields:               true,
	}
}

// Utility method that created the necessary Model Registry's logical-model types
// as the necessary MLMD's Context, Artifact, Execution types etc. in the underlying MLMD service
func CreateMLMDTypes(cc grpc.ClientConnInterface, nameConfig MLMDTypeNamesConfig) (map[string]int64, error) {
	client := proto.NewMetadataStoreServiceClient(cc)

	registeredModelReq := proto.PutContextTypeRequest{
		CanAddFields: &nameConfig.CanAddFields,
		ContextType: &proto.ContextType{
			Name: &nameConfig.RegisteredModelTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
				"owner":       proto.PropertyType_STRING,
				"state":       proto.PropertyType_STRING,
			},
		},
	}

	modelVersionReq := proto.PutContextTypeRequest{
		CanAddFields: &nameConfig.CanAddFields,
		ContextType: &proto.ContextType{
			Name: &nameConfig.ModelVersionTypeName,
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
		CanAddFields: &nameConfig.CanAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: &nameConfig.DocArtifactTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
			},
		},
	}

	modelArtifactReq := proto.PutArtifactTypeRequest{
		CanAddFields: &nameConfig.CanAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: &nameConfig.ModelArtifactTypeName,
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
		CanAddFields: &nameConfig.CanAddFields,
		ContextType: &proto.ContextType{
			Name: &nameConfig.ServingEnvironmentTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
			},
		},
	}

	inferenceServiceReq := proto.PutContextTypeRequest{
		CanAddFields: &nameConfig.CanAddFields,
		ContextType: &proto.ContextType{
			Name: &nameConfig.InferenceServiceTypeName,
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
		CanAddFields: &nameConfig.CanAddFields,
		ExecutionType: &proto.ExecutionType{
			Name: &nameConfig.ServeModelTypeName,
			Properties: map[string]proto.PropertyType{
				"description":      proto.PropertyType_STRING,
				"model_version_id": proto.PropertyType_INT,
			},
		},
	}

	registeredModelResp, err := client.PutContextType(context.Background(), &registeredModelReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up context type %s: %w", nameConfig.RegisteredModelTypeName, err)
	}

	modelVersionResp, err := client.PutContextType(context.Background(), &modelVersionReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up context type %s: %w", nameConfig.ModelVersionTypeName, err)
	}

	docArtifactResp, err := client.PutArtifactType(context.Background(), &docArtifactReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up artifact type %s: %w", nameConfig.DocArtifactTypeName, err)
	}

	modelArtifactResp, err := client.PutArtifactType(context.Background(), &modelArtifactReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up artifact type %s: %w", nameConfig.ModelArtifactTypeName, err)
	}

	servingEnvironmentResp, err := client.PutContextType(context.Background(), &servingEnvironmentReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up context type %s: %w", nameConfig.ServingEnvironmentTypeName, err)
	}

	inferenceServiceResp, err := client.PutContextType(context.Background(), &inferenceServiceReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up context type %s: %w", nameConfig.InferenceServiceTypeName, err)
	}

	serveModelResp, err := client.PutExecutionType(context.Background(), &serveModelReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up execution type %s: %w", nameConfig.ServeModelTypeName, err)
	}

	typesMap := map[string]int64{
		defaults.RegisteredModelTypeName:    registeredModelResp.GetTypeId(),
		defaults.ModelVersionTypeName:       modelVersionResp.GetTypeId(),
		defaults.DocArtifactTypeName:        docArtifactResp.GetTypeId(),
		defaults.ModelArtifactTypeName:      modelArtifactResp.GetTypeId(),
		defaults.ServingEnvironmentTypeName: servingEnvironmentResp.GetTypeId(),
		defaults.InferenceServiceTypeName:   inferenceServiceResp.GetTypeId(),
		defaults.ServeModelTypeName:         serveModelResp.GetTypeId(),
	}
	return typesMap, nil
}
