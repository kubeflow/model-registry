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
	defaultCM, userCM, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	catalogSources := &models.CatalogSourceConfigList{
		Catalogs: make([]models.CatalogSourceConfig, 0),
	}

	if raw, ok := defaultCM.Data[k8s.CatalogSourceKey]; ok {
		defaulfCatalogSources, err := parseCatalogYaml(raw, true)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default catalogs: %w", err)
		}
		catalogSources.Catalogs = append(catalogSources.Catalogs, defaulfCatalogSources...)
	}

	if raw, ok := userCM.Data[k8s.CatalogSourceKey]; ok {
		userManagedSources, err := parseCatalogYaml(raw, false)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default catalogs: %w", err)
		}
		catalogSources.Catalogs = append(catalogSources.Catalogs, userManagedSources...)
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

func parseCatalogYaml(raw string, isDefault bool) ([]models.CatalogSourceConfig, error) {
	// Internal struct to match YAML structure
	var parsed struct {
		Catalogs []struct {
			Name       string                 `yaml:"name"`
			Id         string                 `yaml:"id"`
			Type       string                 `yaml:"type"`
			Enabled    *bool                  `yaml:"enabled"`
			Properties map[string]interface{} `yaml:"properties"`
			Labels     []string               `yaml:"labels"`
		} `yaml:"catalogs"`
	}

	if err := yaml.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse catalogs yaml: %w", err)
	}

	catalogs := make([]models.CatalogSourceConfig, 0, len(parsed.Catalogs))
	for _, c := range parsed.Catalogs {
		entry := models.CatalogSourceConfig{
			Id:        c.Id,
			Name:      c.Name,
			Type:      c.Type,
			Enabled:   c.Enabled,
			Labels:    c.Labels,
			IsDefault: &isDefault,
		}

		if c.Properties != nil {
			if includedModels, ok := c.Properties["includedModels"]; ok {
				entry.IncludedModels = extractStringSlice(includedModels)
			}

			if excludedModels, ok := c.Properties["excludedModels"]; ok {
				entry.ExcludedModels = extractStringSlice(excludedModels)
			}

			if apiKey, ok := c.Properties["apiKey"].(string); ok {
				entry.ApiKey = &apiKey
			}

			if allowedOrganization, ok := c.Properties["allowedOrganization"].(string); ok {
				entry.AllowedOrganization = &allowedOrganization
			}
		}
		catalogs = append(catalogs, entry)
	}

	return catalogs, nil
}

func extractStringSlice(value interface{}) []string {
	if arr, ok := value.([]interface{}); ok {
		result := make([]string, 0, len(arr))
		for _, item := range arr {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	if strSlice, ok := value.([]string); ok {
		return strSlice
	}
	return []string{}

}
