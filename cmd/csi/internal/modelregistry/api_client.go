package modelregistry

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/kubeflow/hub/cmd/csi/internal/constants"
	"github.com/kubeflow/hub/pkg/openapi"
)

func NewAPIClient(cfg *openapi.Configuration, storageUri string, serviceName string, clusterDomain string) *openapi.APIClient {
	client := openapi.NewAPIClient(cfg)

	u, err := url.Parse(storageUri)
	if err == nil {
		ns := u.Query().Get("namespace")
		if ns != "" {
			if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(ns) {
				log.Printf("Invalid namespace parameter: %s", ns)
				return client
			}

			// Extract port from cfg.Host
			_, port, err := net.SplitHostPort(cfg.Host)
			if err != nil {
				port = "8080"
			}

			if !regexp.MustCompile(`^[0-9]+$`).MatchString(port) {
				log.Printf("Invalid port derived from config host: %s", port)
				port = "8080"
			}

			newCfg := openapi.NewConfiguration()
			newCfg.Host = fmt.Sprintf("%s.%s.%s:%s", serviceName, ns, clusterDomain, port)
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
