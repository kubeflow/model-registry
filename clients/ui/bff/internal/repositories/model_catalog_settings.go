package repositories

import (
	"context"
	"fmt"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"gopkg.in/yaml.v3"
)

type ModelCatalogSettingsRepository struct {
}

func NewModelCatalogSettingsRepository() *ModelCatalogSettingsRepository {
	return &ModelCatalogSettingsRepository{}
}

func (r *ModelCatalogSettingsRepository) GetAllCatalogSourceConfigs(ctx context.Context, client k8s.KubernetesClientInterface, namespace string) (*models.CatalogSourceConfigList, error) {
	// TODO ppadti we need to merge this catalog source with the other one
	defaultCM, _, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	raw, ok := defaultCM.Data[k8s.CatalogSourceKey]
	if !ok {
		return nil, fmt.Errorf("default catalog sources configmap missing key")
	}

	// Internal struct to match YAML structure
	var parsed struct {
		Catalogs []struct {
			Name       string            `yaml:"name"`
			Id         string            `yaml:"id"`
			Type       string            `yaml:"type"`
			Enabled    *bool             `yaml:"enabled"`
			Properties map[string]string `yaml:"properties"`
			Labels     []string          `yaml:"labels"`
		} `yaml:"catalogs"`
	}

	if err := yaml.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse catalogs yaml: %w", err)
	}

	catalogSources := &models.CatalogSourceConfigList{
		Catalogs: make([]models.CatalogSourceConfig, 0, len(parsed.Catalogs)),
	}

	for _, c := range parsed.Catalogs {
		isDefault := true

		entry := models.CatalogSourceConfig{
			Id:        c.Id,
			Name:      c.Name,
			Type:      c.Type,
			Enabled:   c.Enabled,
			Labels:    c.Labels,
			IsDefault: &isDefault,
		}

		catalogSources.Catalogs = append(catalogSources.Catalogs, entry)
	}

	return catalogSources, nil
}

func (r *ModelCatalogSettingsRepository) GetCatalogSourceConfig(ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	catalogSourceId string,
) (*models.CatalogSourceConfig, error) {
	// TODO ppadti write the real implementation here calling k8s client
	return nil, fmt.Errorf("not implemented yet")
}

func (r *ModelCatalogSettingsRepository) CreateCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	payload models.CatalogSourceConfigPayload,
) (*models.CatalogSourceConfig, error) {

	// TODO ppadti write the real implementation here calling k8s client
	return nil, fmt.Errorf("not implemented yet")
}

func (r *ModelCatalogSettingsRepository) UpdateCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	payload models.CatalogSourceConfigPayload,
) (*models.CatalogSourceConfig, error) {
	// TODO ppadti write the real implementation here calling k8s client
	return nil, fmt.Errorf("not implemented yet")
}

func (r *ModelCatalogSettingsRepository) DeleteCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	catalogSourceId string,
) (*models.CatalogSourceConfig, error) {
	// TODO ppadti write the real implementation here calling k8s client
	return nil, fmt.Errorf("not implemented yet")
}
