package modelregistry

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/kubeflow/model-registry/cmd/csi/internal/constants"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

func NewAPIClient(cfg *openapi.Configuration, storageUri string, serviceName string) *openapi.APIClient {
	client := openapi.NewAPIClient(cfg)

	u, err := url.Parse(storageUri)
	if err == nil {
		ns := u.Query().Get("namespace")
		if ns != "" {
			// Extract port from cfg.Host
			hostParts := strings.Split(cfg.Host, ":")
			port := "8080"
			if len(hostParts) == 2 {
				port = hostParts[1]
			}
			
			newCfg := openapi.NewConfiguration()
			newCfg.Host = fmt.Sprintf("%s.%s.svc.cluster.local:%s", serviceName, ns, port)
			newCfg.Scheme = cfg.Scheme
			
			log.Printf("Using explicit namespace=%s from storageUri query parameter. Target host override: %s", ns, newCfg.Host)
			return openapi.NewAPIClient(newCfg)
		}
	}

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
