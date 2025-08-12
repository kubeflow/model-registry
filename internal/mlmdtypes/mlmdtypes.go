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
	DataSetTypeName            string
	MetricTypeName             string
	MetricHistoryTypeName      string
	ParameterTypeName          string
	ServingEnvironmentTypeName string
	InferenceServiceTypeName   string
	ServeModelTypeName         string
	ExperimentTypeName         string
	ExperimentRunTypeName      string
	CanAddFields               bool
}

func NewMLMDTypeNamesConfigFromDefaults() MLMDTypeNamesConfig {
	return MLMDTypeNamesConfig{
		RegisteredModelTypeName:    defaults.RegisteredModelTypeName,
		ModelVersionTypeName:       defaults.ModelVersionTypeName,
		ModelArtifactTypeName:      defaults.ModelArtifactTypeName,
		DocArtifactTypeName:        defaults.DocArtifactTypeName,
		DataSetTypeName:            defaults.DataSetTypeName,
		MetricTypeName:             defaults.MetricTypeName,
		MetricHistoryTypeName:      defaults.MetricHistoryTypeName,
		ParameterTypeName:          defaults.ParameterTypeName,
		ServingEnvironmentTypeName: defaults.ServingEnvironmentTypeName,
		InferenceServiceTypeName:   defaults.InferenceServiceTypeName,
		ServeModelTypeName:         defaults.ServeModelTypeName,
		ExperimentTypeName:         defaults.ExperimentTypeName,
		ExperimentRunTypeName:      defaults.ExperimentRunTypeName,
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
				"description":  proto.PropertyType_STRING,
				"language":     proto.PropertyType_STRUCT,
				"library_name": proto.PropertyType_STRING,
				"license_link": proto.PropertyType_STRING,
				"license":      proto.PropertyType_STRING,
				"logo":         proto.PropertyType_STRING,
				"maturity":     proto.PropertyType_STRING,
				"owner":        proto.PropertyType_STRING,
				"provider":     proto.PropertyType_STRING,
				"readme":       proto.PropertyType_STRING,
				"state":        proto.PropertyType_STRING,
				"tasks":        proto.PropertyType_STRUCT,
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
				"model_source_kind":    proto.PropertyType_STRING,
				"model_source_class":   proto.PropertyType_STRING,
				"model_source_group":   proto.PropertyType_STRING,
				"model_source_id":      proto.PropertyType_STRING,
				"model_source_name":    proto.PropertyType_STRING,
			},
		},
	}

	dataSetReq := proto.PutArtifactTypeRequest{
		CanAddFields: &nameConfig.CanAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: &nameConfig.DataSetTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
				"digest":      proto.PropertyType_STRING,
				"source_type": proto.PropertyType_STRING,
				"source":      proto.PropertyType_STRING,
				"schema":      proto.PropertyType_STRING,
				"profile":     proto.PropertyType_STRING,
			},
		},
	}

	metricArtifactReq := proto.PutArtifactTypeRequest{
		CanAddFields: &nameConfig.CanAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: &nameConfig.MetricTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
				"value":       proto.PropertyType_DOUBLE,
				"timestamp":   proto.PropertyType_STRING,
				"step":        proto.PropertyType_INT,
			},
		},
	}

	metricHistoryArtifactReq := proto.PutArtifactTypeRequest{
		CanAddFields: &nameConfig.CanAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: &nameConfig.MetricHistoryTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
				"value":       proto.PropertyType_DOUBLE,
				"timestamp":   proto.PropertyType_STRING,
				"step":        proto.PropertyType_INT,
			},
		},
	}

	parameterReq := proto.PutArtifactTypeRequest{
		CanAddFields: &nameConfig.CanAddFields,
		ArtifactType: &proto.ArtifactType{
			Name: &nameConfig.ParameterTypeName,
			Properties: map[string]proto.PropertyType{
				"description":    proto.PropertyType_STRING,
				"value":          proto.PropertyType_STRING,
				"parameter_type": proto.PropertyType_STRING,
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

	experimentReq := proto.PutContextTypeRequest{
		CanAddFields: &nameConfig.CanAddFields,
		ContextType: &proto.ContextType{
			Name: &nameConfig.ExperimentTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
				"owner":       proto.PropertyType_STRING,
				"state":       proto.PropertyType_STRING,
			},
		},
	}

	experimentRunReq := proto.PutContextTypeRequest{
		CanAddFields: &nameConfig.CanAddFields,
		ContextType: &proto.ContextType{
			Name: &nameConfig.ExperimentRunTypeName,
			Properties: map[string]proto.PropertyType{
				"description":            proto.PropertyType_STRING,
				"owner":                  proto.PropertyType_STRING,
				"state":                  proto.PropertyType_STRING,
				"status":                 proto.PropertyType_STRING,
				"start_time_since_epoch": proto.PropertyType_STRING,
				"end_time_since_epoch":   proto.PropertyType_STRING,
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

	dataSetResp, err := client.PutArtifactType(context.Background(), &dataSetReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up artifact type %s: %w", nameConfig.DataSetTypeName, err)
	}

	metricArtifactResp, err := client.PutArtifactType(context.Background(), &metricArtifactReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up artifact type %s: %w", nameConfig.MetricTypeName, err)
	}

	metricHistoryArtifactResp, err := client.PutArtifactType(context.Background(), &metricHistoryArtifactReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up artifact type %s: %w", nameConfig.MetricHistoryTypeName, err)
	}

	parameterResp, err := client.PutArtifactType(context.Background(), &parameterReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up artifact type %s: %w", nameConfig.ParameterTypeName, err)
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

	experimentResp, err := client.PutContextType(context.Background(), &experimentReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up context type %s: %w", nameConfig.ExperimentTypeName, err)
	}

	experimentRunResp, err := client.PutContextType(context.Background(), &experimentRunReq)
	if err != nil {
		return nil, fmt.Errorf("error setting up context type %s: %w", nameConfig.ExperimentRunTypeName, err)
	}

	typesMap := map[string]int64{
		defaults.RegisteredModelTypeName:    registeredModelResp.GetTypeId(),
		defaults.ModelVersionTypeName:       modelVersionResp.GetTypeId(),
		defaults.DocArtifactTypeName:        docArtifactResp.GetTypeId(),
		defaults.ModelArtifactTypeName:      modelArtifactResp.GetTypeId(),
		defaults.DataSetTypeName:            dataSetResp.GetTypeId(),
		defaults.MetricTypeName:             metricArtifactResp.GetTypeId(),
		defaults.MetricHistoryTypeName:      metricHistoryArtifactResp.GetTypeId(),
		defaults.ParameterTypeName:          parameterResp.GetTypeId(),
		defaults.ServingEnvironmentTypeName: servingEnvironmentResp.GetTypeId(),
		defaults.InferenceServiceTypeName:   inferenceServiceResp.GetTypeId(),
		defaults.ServeModelTypeName:         serveModelResp.GetTypeId(),
		defaults.ExperimentTypeName:         experimentResp.GetTypeId(),
		defaults.ExperimentRunTypeName:      experimentRunResp.GetTypeId(),
	}
	return typesMap, nil
}
