package core

import (
	"context"
	"fmt"

	"github.com/kubeflow/model-registry/internal/converter/generated"
	"github.com/kubeflow/model-registry/internal/mapper"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/internal/mlmdtypes"
	"github.com/kubeflow/model-registry/pkg/api"
	"google.golang.org/grpc"
)

// ModelRegistryService is the core library of the model registry
type ModelRegistryService struct {
	mlmdClient  proto.MetadataStoreServiceClient
	typesMap    map[string]int64
	mapper      *mapper.Mapper
	openapiConv *generated.OpenAPIConverterImpl
	nameConfig  mlmdtypes.MLMDTypeNamesConfig
}

// NewModelRegistryService creates a new instance of the ModelRegistryService, initializing it with the provided gRPC client connection.
// It _assumes_ the necessary MLMD's Context, Artifact, Execution types etc. are already setup in the underlying MLMD service.
//
// Parameters:
//   - cc: A gRPC client connection to the underlying MLMD service
func NewModelRegistryService(cc grpc.ClientConnInterface, nameConfig mlmdtypes.MLMDTypeNamesConfig) (api.ModelRegistryApi, error) {
	typesMap, err := BuildTypesMap(cc, nameConfig)
	if err != nil { // early return in case type Ids cannot be retrieved
		return nil, err
	}

	client := proto.NewMetadataStoreServiceClient(cc)

	return &ModelRegistryService{
		mlmdClient:  client,
		nameConfig:  nameConfig,
		typesMap:    typesMap,
		openapiConv: &generated.OpenAPIConverterImpl{},
		mapper:      mapper.NewMapper(typesMap),
	}, nil
}

func BuildTypesMap(cc grpc.ClientConnInterface, nameConfig mlmdtypes.MLMDTypeNamesConfig) (map[string]int64, error) {
	client := proto.NewMetadataStoreServiceClient(cc)

	registeredModelContextTypeReq := proto.GetContextTypeRequest{
		TypeName: &nameConfig.RegisteredModelTypeName,
	}
	registeredModelResp, err := client.GetContextType(context.Background(), &registeredModelContextTypeReq)
	if err != nil {
		return nil, fmt.Errorf("error getting context type %s: %w", nameConfig.RegisteredModelTypeName, err)
	}
	modelVersionContextTypeReq := proto.GetContextTypeRequest{
		TypeName: &nameConfig.ModelVersionTypeName,
	}
	modelVersionResp, err := client.GetContextType(context.Background(), &modelVersionContextTypeReq)
	if err != nil {
		return nil, fmt.Errorf("error getting context type %s: %w", nameConfig.ModelVersionTypeName, err)
	}
	docArtifactResp, err := client.GetArtifactType(context.Background(), &proto.GetArtifactTypeRequest{
		TypeName: &nameConfig.DocArtifactTypeName,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting artifact type %s: %w", nameConfig.DocArtifactTypeName, err)
	}
	modelArtifactArtifactTypeReq := proto.GetArtifactTypeRequest{
		TypeName: &nameConfig.ModelArtifactTypeName,
	}
	modelArtifactResp, err := client.GetArtifactType(context.Background(), &modelArtifactArtifactTypeReq)
	if err != nil {
		return nil, fmt.Errorf("error getting artifact type %s: %w", nameConfig.ModelArtifactTypeName, err)
	}
	servingEnvironmentContextTypeReq := proto.GetContextTypeRequest{
		TypeName: &nameConfig.ServingEnvironmentTypeName,
	}
	servingEnvironmentResp, err := client.GetContextType(context.Background(), &servingEnvironmentContextTypeReq)
	if err != nil {
		return nil, fmt.Errorf("error getting context type %s: %w", nameConfig.ServingEnvironmentTypeName, err)
	}
	inferenceServiceContextTypeReq := proto.GetContextTypeRequest{
		TypeName: &nameConfig.InferenceServiceTypeName,
	}
	inferenceServiceResp, err := client.GetContextType(context.Background(), &inferenceServiceContextTypeReq)
	if err != nil {
		return nil, fmt.Errorf("error getting context type %s: %w", nameConfig.InferenceServiceTypeName, err)
	}
	serveModelExecutionReq := proto.GetExecutionTypeRequest{
		TypeName: &nameConfig.ServeModelTypeName,
	}
	serveModelResp, err := client.GetExecutionType(context.Background(), &serveModelExecutionReq)
	if err != nil {
		return nil, fmt.Errorf("error getting execution type %s: %w", nameConfig.ServeModelTypeName, err)
	}

	typesMap := map[string]int64{
		nameConfig.RegisteredModelTypeName:    registeredModelResp.ContextType.GetId(),
		nameConfig.ModelVersionTypeName:       modelVersionResp.ContextType.GetId(),
		nameConfig.DocArtifactTypeName:        docArtifactResp.ArtifactType.GetId(),
		nameConfig.ModelArtifactTypeName:      modelArtifactResp.ArtifactType.GetId(),
		nameConfig.ServingEnvironmentTypeName: servingEnvironmentResp.ContextType.GetId(),
		nameConfig.InferenceServiceTypeName:   inferenceServiceResp.ContextType.GetId(),
		nameConfig.ServeModelTypeName:         serveModelResp.ExecutionType.GetId(),
	}
	return typesMap, nil
}
