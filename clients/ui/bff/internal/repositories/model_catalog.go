package repositories

import (
	"context"
	"fmt"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelCatalogRepository struct {
}

const ModelCatalogServiceName = "model-catalog"

func NewCatalogRepository() *ModelCatalogRepository {
	return &ModelCatalogRepository{}
}

func (m *ModelCatalogRepository) GetModelCatalogWithMode(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string, isFederatedMode bool) (models.ModelCatalogModel, error) {

	s, err := client.GetServiceDetailsByName(sessionCtx, namespace, ModelCatalogServiceName)
	if err != nil {
		return models.ModelCatalogModel{}, fmt.Errorf("error fetching model catalog: %w", err)
	}

	modelCatalog := models.ModelCatalogModel{
		Name:          s.Name,
		Description:   s.Description,
		DisplayName:   s.DisplayName,
		ServerAddress: m.ResolveServerAddress(s.ClusterIP, s.HTTPPort, s.IsHTTPS, s.ExternalAddressRest, isFederatedMode),
		IsHTTPS:       s.IsHTTPS,
	}
	return modelCatalog, nil
}

func (m *ModelCatalogRepository) ResolveServerAddress(clusterIP string, httpPort int32, isHTTPS bool, externalAddressRest string, isFederatedMode bool) string {
	// Default behavior - use cluster IP and port
	protocol := "http"
	if isHTTPS {
		protocol = "https"
	}
	// In federated mode, if external address is available, use it
	if isFederatedMode && externalAddressRest != "" {
		// External address is assumed to be HTTPS
		url := fmt.Sprintf("%s://%s/api/model_catalog/v1alpha1", protocol, externalAddressRest)
		return url
	}

	url := fmt.Sprintf("%s://%s:%d/api/model_catalog/v1alpha1", protocol, clusterIP, httpPort)
	return url
}
