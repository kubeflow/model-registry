package modelregistry

import (
	"context"
	"log"
	"strings"

	"github.com/kubeflow/model-registry/internal/csi/constants"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

func NewAPIClient(cfg *openapi.Configuration, storageUri string) *openapi.APIClient {
	client := openapi.NewAPIClient(cfg)

	// Parse the URI to retrieve the needed information to query model registry (modelArtifact)
	mrUri := strings.TrimPrefix(storageUri, string(constants.MR))

	tokens := strings.SplitN(mrUri, "/", 3)

	if len(tokens) < 2 {
		return client
	}

	newCfg := openapi.NewConfiguration()
	newCfg.Host = tokens[0]
	newCfg.Scheme = cfg.Scheme

	newClient := openapi.NewAPIClient(newCfg)

	if len(tokens) == 2 {
		// Check if the model registry service is available
		_, _, err := newClient.ModelRegistryServiceAPI.GetRegisteredModels(context.Background()).Execute()
		if err != nil {
			log.Printf("Falling back to base url %s for model registry service", cfg.Host)

			return client
		}
	}

	return newClient
}
