package repositories

import (
	"context"
	"fmt"
	"strings"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	catalogMap := make(map[string]models.CatalogSourceConfig)

	if raw, ok := defaultCM.Data[k8s.CatalogSourceKey]; ok {
		defaultCatalogSources, err := parseCatalogYaml(raw, true)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default catalogs: %w", err)
		}

		for _, catalog := range defaultCatalogSources {
			catalogMap[catalog.Id] = catalog
		}
	}

	if raw, ok := userCM.Data[k8s.CatalogSourceKey]; ok {
		userManagedSources, err := parseCatalogYaml(raw, false)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user managed catalogs: %w", err)
		}

		for _, userCatalogSource := range userManagedSources {
			if existingSource, exist := catalogMap[userCatalogSource.Id]; exist {
				mergedCatalogSources := mergeCatalogSourceConfigs(existingSource, userCatalogSource)
				catalogMap[userCatalogSource.Id] = mergedCatalogSources
			} else {
				catalogMap[userCatalogSource.Id] = userCatalogSource
			}
		}
	}

	catalogSources := &models.CatalogSourceConfigList{
		Catalogs: make([]models.CatalogSourceConfig, 0),
	}

	for _, c := range catalogMap {
		catalogSources.Catalogs = append(catalogSources.Catalogs, c)
	}

	return catalogSources, nil
}

func (r *ModelCatalogSettingsRepository) GetCatalogSourceConfig(ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	catalogSourceId string,
) (*models.CatalogSourceConfig, error) {
	defaultCM, userCM, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	defaultSource := findCatalogSourceById(defaultCM.Data[k8s.CatalogSourceKey], catalogSourceId, true)
	userSource := findCatalogSourceById(userCM.Data[k8s.CatalogSourceKey], catalogSourceId, false)

	var result *models.CatalogSourceConfig

	if userSource != nil {
		if defaultSource != nil {
			merged := mergeCatalogSourceConfigs(*defaultSource, *userSource)
			result = &merged
		} else {
			result = userSource
		}
	} else if defaultSource != nil {
		result = defaultSource
	} else {
		return nil, fmt.Errorf("catalog source not found: %s", catalogSourceId)
	}

	secretName, yamlFilePath := findCatalogSourceProperties(userCM.Data[k8s.CatalogSourceKey], catalogSourceId)
	if secretName == "" && yamlFilePath == "" {
		secretName, yamlFilePath = findCatalogSourceProperties(defaultCM.Data[k8s.CatalogSourceKey], catalogSourceId)
	}

	switch result.Type {
	case "yaml":
		if yamlFilePath != "" {

			if yamlContent, ok := userCM.Data[yamlFilePath]; ok {
				result.Yaml = &yamlContent
			} else if yamlContent, ok := defaultCM.Data[yamlFilePath]; ok {
				result.Yaml = &yamlContent
			}
		}

	case "huggingface":
		if secretName != "" {
			apiKey, err := client.GetSecretValue(ctx, namespace, secretName, "apiKey")
			if err == nil && apiKey != "" {
				result.ApiKey = &apiKey
			}
		}
	}

	return result, nil
}

func (r *ModelCatalogSettingsRepository) CreateCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	payload models.CatalogSourceConfigPayload,
) (*models.CatalogSourceConfig, error) {
	if err := validateCatalogSourceConfigPayload(payload); err != nil {
		return nil, err
	}

	defaultCM, userCM, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	if findCatalogSourceById(defaultCM.Data[k8s.CatalogSourceKey], payload.Id, true) != nil {
		return nil, fmt.Errorf("catalog source '%s' already exists in default sources", payload.Id)
	}

	if findCatalogSourceById(userCM.Data[k8s.CatalogSourceKey], payload.Id, false) != nil {
		return nil, fmt.Errorf("catalog source '%s' already exists in user managed sources", payload.Id)
	}

	var secretName string
	var yamlFileName string
	yamlContent := make(map[string]string)

	switch payload.Type {
	case "yaml":
		yamlFileName = fmt.Sprintf("%s.yaml", payload.Id)
		yamlContent[yamlFileName] = *payload.Yaml
	case "huggingface":
		secretName, err = createSecretForHuggingFace(ctx, client, namespace, payload.Id, *payload.ApiKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create secret for huggingface source: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported catalog type: %s", payload.Type)
	}

	newEntry := convertSourceConfigToYamlEntry(payload, yamlFileName, secretName)

	existingConfigMapEntry := userCM.Data[k8s.CatalogSourceKey]

	updatedConfigMapEntry, err := appendCatalogSourceToYaml(existingConfigMapEntry, newEntry)
	if err != nil {
		if secretName != "" {
			deleteSecretForHuggingFace(ctx, client, namespace, secretName)
		}
		return nil, fmt.Errorf("failed to append catalog to yaml: %w", err)
	}

	if userCM.Data == nil {
		userCM.Data = make(map[string]string)
	}
	userCM.Data[k8s.CatalogSourceKey] = updatedConfigMapEntry

	for key, value := range yamlContent {
		userCM.Data[key] = value
	}

	err = client.UpdateCatalogSourceConfig(ctx, namespace, &userCM)
	if err != nil {
		if secretName != "" {
			deleteSecretForHuggingFace(ctx, client, namespace, secretName)
		}
		return nil, fmt.Errorf("failed to update user configmap: %w", err)
	}

	result := models.CatalogSourceConfig(payload)
	isDefault := false
	result.IsDefault = &isDefault

	return &result, nil
}

func (r *ModelCatalogSettingsRepository) UpdateCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	sourceId string,
	payload models.CatalogSourceConfigPayload,
) (*models.CatalogSourceConfig, error) {
	defaultCM, userCM, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	existingUserSource := findCatalogSourceById(userCM.Data[k8s.CatalogSourceKey], sourceId, false)
	existingDefaultSource := findCatalogSourceById(defaultCM.Data[k8s.CatalogSourceKey], sourceId, true)

	var existingCatalog *models.CatalogSourceConfig
	var isOverridingDefault bool

	if existingUserSource != nil {
		existingCatalog = existingUserSource
		isOverridingDefault = existingDefaultSource != nil
	} else if existingDefaultSource != nil {
		existingCatalog = existingDefaultSource
		isOverridingDefault = true
	} else {
		return nil, fmt.Errorf("catalog source '%s' not found", sourceId)
	}

	if payload.Type != "" && payload.Type != existingCatalog.Type {
		return nil, fmt.Errorf(
			"cannot change catalog source type from '%s' to '%s'",
			existingCatalog.Type,
			payload.Type,
		)
	}

	catalogType := existingCatalog.Type

	if isOverridingDefault {
		if err := validateUpdatePayloadForDefaultOverride(payload); err != nil {
			return nil, err
		}
	}

	var secretName, yamlFilePath string
	if !isOverridingDefault || existingUserSource != nil {
		secretName, yamlFilePath = findCatalogSourceProperties(userCM.Data[k8s.CatalogSourceKey], sourceId)
		if secretName == "" || yamlFilePath == "" {
			defaultSecretName, defaultYamlPath := findCatalogSourceProperties(defaultCM.Data[k8s.CatalogSourceKey], sourceId)
			if secretName == "" {
				secretName = defaultSecretName
			}
			if yamlFilePath == "" {
				yamlFilePath = defaultYamlPath
			}
		}
	}

	if userCM.Data == nil {
		userCM.Data = make(map[string]string)
	}

	if !isOverridingDefault || existingUserSource != nil {
		switch catalogType {
		case "yaml":
			if payload.Yaml != nil && *payload.Yaml != "" {
				if yamlFilePath == "" {
					yamlFilePath = fmt.Sprintf("%s.yaml", sourceId)
				}
				userCM.Data[yamlFilePath] = *payload.Yaml
			}

		case "huggingface":
			if payload.ApiKey != nil && *payload.ApiKey != "" {
				if secretName != "" {
					err := client.PatchSecret(ctx, namespace, secretName, map[string]string{
						"apiKey": *payload.ApiKey,
					})
					if err != nil {
						if isNotFoundError(err) {
							secretName, err = createSecretForHuggingFace(ctx, client, namespace, sourceId, *payload.ApiKey)
							if err != nil {
								return nil, fmt.Errorf("failed to create replacement secret: %w", err)
							}
						} else {
							return nil, fmt.Errorf("failed to patch secret '%s': %w", secretName, err)
						}
					}
				} else {
					secretName, err = createSecretForHuggingFace(ctx, client, namespace, sourceId, *payload.ApiKey)
					if err != nil {
						return nil, fmt.Errorf("failed to create secret: %w", err)
					}
				}
			}
		}
	}

	if existingUserSource != nil {
		updatedYAML, err := updateCatalogSourceInYAML(
			userCM.Data[k8s.CatalogSourceKey],
			sourceId,
			payload,
			secretName,
			yamlFilePath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update catalog in yaml: %w", err)
		}
		userCM.Data[k8s.CatalogSourceKey] = updatedYAML
	} else {
		overrideEntry := buildOverrideEntryForDefaultSource(sourceId, payload)
		updatedYAML, err := appendCatalogSourceToYaml(userCM.Data[k8s.CatalogSourceKey], overrideEntry)
		if err != nil {
			return nil, fmt.Errorf("failed to append override entry: %w", err)
		}
		userCM.Data[k8s.CatalogSourceKey] = updatedYAML
	}

	err = client.UpdateCatalogSourceConfig(ctx, namespace, &userCM)
	if err != nil {
		return nil, fmt.Errorf("failed to update user configmap: %w", err)
	}

	return r.GetCatalogSourceConfig(ctx, client, namespace, sourceId)

}

func (r *ModelCatalogSettingsRepository) DeleteCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	catalogSourceId string,
) (*models.CatalogSourceConfig, error) {
	defaultCM, userCM, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	if findCatalogSourceById(defaultCM.Data[k8s.CatalogSourceKey], catalogSourceId, true) != nil {
		return nil, fmt.Errorf(
			"cannot delete catalog source '%s': it is a default source. "+
				"Default sources cannot be deleted, only customized via PATCH",
			catalogSourceId,
		)
	}

	catalogSourceToDelete := findCatalogSourceById(userCM.Data[k8s.CatalogSourceKey], catalogSourceId, false)
	if catalogSourceToDelete == nil {
		return nil, fmt.Errorf("catalog source '%s' not found in user sources", catalogSourceId)
	}

	secretName, yamlFilePath := findCatalogSourceProperties(userCM.Data[k8s.CatalogSourceKey], catalogSourceId)

	if catalogSourceToDelete.Type == "huggingface" && secretName != "" {
		err := client.DeleteSecret(ctx, namespace, secretName)
		if err != nil {
			if !isNotFoundError(err) {
				return nil, fmt.Errorf(
					"failed to delete secret '%s' for catalog '%s': %w. Source was not deleted",
					secretName,
					catalogSourceId,
					err,
				)
			}
		}
	}

	if catalogSourceToDelete.Type == "yaml" && yamlFilePath != "" {
		delete(userCM.Data, yamlFilePath)
	}

	updatedYAML, err := removeCatalogSourceFromYAML(userCM.Data[k8s.CatalogSourceKey], catalogSourceId)
	if err != nil {
		return nil, fmt.Errorf("failed to remove catalog from sources.yaml: %w", err)
	}
	userCM.Data[k8s.CatalogSourceKey] = updatedYAML

	err = client.UpdateCatalogSourceConfig(ctx, namespace, &userCM)
	if err != nil {
		return nil, fmt.Errorf("failed to update configmap after deletion: %w", err)
	}

	return catalogSourceToDelete, nil
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

func mergeCatalogSourceConfigs(defaultCatalog models.CatalogSourceConfig, userCatalog models.CatalogSourceConfig) models.CatalogSourceConfig {
	mergedSource := defaultCatalog

	if userCatalog.Name != "" {
		mergedSource.Name = userCatalog.Name
	}

	if userCatalog.Type != "" {
		mergedSource.Type = userCatalog.Type
	}

	if userCatalog.Enabled != nil {
		mergedSource.Enabled = userCatalog.Enabled
	}

	if userCatalog.AllowedOrganization != nil {
		mergedSource.AllowedOrganization = userCatalog.AllowedOrganization
	}

	if userCatalog.ApiKey != nil {
		mergedSource.ApiKey = userCatalog.ApiKey
	}

	if len(userCatalog.IncludedModels) > 0 {
		mergedSource.IncludedModels = userCatalog.IncludedModels
	}

	if len(userCatalog.ExcludedModels) > 0 {
		mergedSource.ExcludedModels = userCatalog.ExcludedModels
	}

	if userCatalog.Yaml != nil {
		mergedSource.Yaml = userCatalog.Yaml
	}

	return mergedSource
}

func validateCatalogSourceConfigPayload(payload models.CatalogSourceConfigPayload) error {
	if payload.Id == "" {
		return fmt.Errorf("id is required")
	}

	if payload.Name == "" {
		return fmt.Errorf("name is required")
	}

	if payload.Type == "" {
		return fmt.Errorf("type is required")
	}

	switch payload.Type {
	case "yaml":
		if payload.Yaml == nil || *payload.Yaml == "" {
			return fmt.Errorf("yaml field is required for yaml-type sources")
		}
	case "huggingface":
		if payload.ApiKey == nil || *payload.ApiKey == "" {
			return fmt.Errorf("apiKey is required for huggingface-type sources")
		}
	default:
		return fmt.Errorf("unsupported catalog type: %s (supported: yaml, huggingface)", payload.Type)
	}
	return nil
}

func findCatalogSourceById(sourceYAML string, catalogId string, isDefault bool) *models.CatalogSourceConfig {
	if sourceYAML == "" {
		return nil
	}

	var parsed struct {
		Catalogs []struct {
			Id         string                 `yaml:"id"`
			Name       string                 `yaml:"name"`
			Type       string                 `yaml:"type"`
			Enabled    *bool                  `yaml:"enabled"`
			Labels     []string               `yaml:"labels"`
			Properties map[string]interface{} `yaml:"properties"`
		} `yaml:"catalogs"`
	}

	if err := yaml.Unmarshal([]byte(sourceYAML), &parsed); err != nil {
		return nil
	}

	for _, catalog := range parsed.Catalogs {
		if catalog.Id == catalogId {
			isDefaultVal := isDefault
			result := &models.CatalogSourceConfig{
				Id:        catalog.Id,
				Name:      catalog.Name,
				Type:      catalog.Type,
				Enabled:   catalog.Enabled,
				Labels:    catalog.Labels,
				IsDefault: &isDefaultVal,
			}

			if catalog.Properties != nil {
				if included, ok := catalog.Properties["includedModels"]; ok {
					result.IncludedModels = extractStringSlice(included)
				}
				if excluded, ok := catalog.Properties["excludedModels"]; ok {
					result.ExcludedModels = extractStringSlice(excluded)
				}
				if org, ok := catalog.Properties["allowedOrganization"].(string); ok {
					result.AllowedOrganization = &org
				}
			}

			return result
		}
	}

	return nil
}

func convertSourceConfigToYamlEntry(payload models.CatalogSourceConfigPayload,
	yamlFileName string,
	secretName string) map[string]interface{} {
	entry := map[string]interface{}{
		"id":      payload.Id,
		"name":    payload.Name,
		"type":    payload.Type,
		"enabled": payload.Enabled,
	}

	if len(payload.Labels) > 0 {
		entry["labels"] = payload.Labels
	}

	properties := make(map[string]interface{})

	switch payload.Type {
	case "yaml":
		properties["yamlCatalogPath"] = yamlFileName
	case "huggingface":
		properties["apiKey"] = secretName
		if payload.AllowedOrganization != nil {
			properties["allowedOrganization"] = *payload.AllowedOrganization
		}

	}

	if len(payload.IncludedModels) > 0 {
		properties["includedModels"] = payload.IncludedModels
	}
	if len(payload.ExcludedModels) > 0 {
		properties["excludedModels"] = payload.ExcludedModels
	}

	if len(properties) > 0 {
		entry["properties"] = properties
	}
	return entry
}

func appendCatalogSourceToYaml(existingConfigMapEntry string, newEntry map[string]interface{}) (string, error) {
	var parsed struct {
		Catalogs []map[string]interface{} `yaml:"catalogs"`
	}

	if existingConfigMapEntry != "" {
		if err := yaml.Unmarshal([]byte(existingConfigMapEntry), &parsed); err != nil {
			return "", fmt.Errorf("failed to parse existing sources.yaml: %w", err)
		}
	} else {
		parsed.Catalogs = []map[string]interface{}{}
	}
	parsed.Catalogs = append(parsed.Catalogs, newEntry)

	updatedBytes, err := yaml.Marshal(parsed)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated sources.yaml: %w", err)
	}

	return string(updatedBytes), nil
}

func createSecretForHuggingFace(ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	catalogId string,
	apiKey string) (string, error) {
	modifiedSecretName := strings.ReplaceAll(catalogId, "_", "-")
	secretName := fmt.Sprintf("catalog-%s-apikey", modifiedSecretName)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/component": "model-catalog",
				"catalog-source-id":           catalogId,
			},
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"apiKey": apiKey,
		},
	}
	err := client.CreateSecret(ctx, namespace, secret)
	if err != nil {
		return "", fmt.Errorf("failed to create secret '%s': %w", secretName, err)
	}

	return secretName, nil

}

func deleteSecretForHuggingFace(ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	secretName string) {
	_ = client.DeleteSecret(ctx, namespace, secretName)
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "NotFound")
}

func removeCatalogSourceFromYAML(existingYAML string, sourceId string) (string, error) {
	var parsed struct {
		Catalogs []map[string]interface{} `yaml:"catalogs"`
	}

	if err := yaml.Unmarshal([]byte(existingYAML), &parsed); err != nil {
		return "", fmt.Errorf("failed to parse sources.yaml: %w", err)
	}

	filteredCatalogs := make([]map[string]interface{}, 0)
	for _, catalogSource := range parsed.Catalogs {
		if id, ok := catalogSource["id"].(string); ok && id != sourceId {
			filteredCatalogs = append(filteredCatalogs, catalogSource)
		}
	}

	parsed.Catalogs = filteredCatalogs
	updatedBytes, err := yaml.Marshal(parsed)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated sources.yaml: %w", err)
	}

	return string(updatedBytes), nil
}

func findCatalogSourceProperties(sourceYAML string, sourceId string) (secretName string, yamlPath string) {
	var parsed struct {
		Catalogs []struct {
			Id         string                 `yaml:"id"`
			Properties map[string]interface{} `yaml:"properties"`
		} `yaml:"catalogs"`
	}

	if err := yaml.Unmarshal([]byte(sourceYAML), &parsed); err != nil {
		return "", ""
	}

	for _, catalogSource := range parsed.Catalogs {
		if catalogSource.Id == sourceId && catalogSource.Properties != nil {
			secretName, _ = catalogSource.Properties["apiKey"].(string)
			yamlPath, _ = catalogSource.Properties["yamlCatalogPath"].(string)
			return
		}
	}
	return "", ""
}

func validateUpdatePayloadForDefaultOverride(payload models.CatalogSourceConfigPayload) error {
	if payload.Name != "" {
		return fmt.Errorf("cannot change 'name' of a default catalog source")
	}

	if len(payload.Labels) > 0 {
		return fmt.Errorf("cannot change 'labels' of a default catalog source")
	}

	if payload.Yaml != nil && *payload.Yaml != "" {
		return fmt.Errorf("cannot change 'yaml' content of a default catalog source")
	}

	return nil
}

func updateCatalogSourceInYAML(
	existingYAML string,
	catalogId string,
	payload models.CatalogSourceConfigPayload,
	secretName string,
	yamlFilePath string,
) (string, error) {
	var parsed struct {
		Catalogs []map[string]interface{} `yaml:"catalogs"`
	}

	if existingYAML == "" {
		return "", fmt.Errorf("no existing yaml to update")
	}

	if err := yaml.Unmarshal([]byte(existingYAML), &parsed); err != nil {
		return "", fmt.Errorf("failed to parse sources.yaml: %w", err)
	}

	found := false
	for i, catalogSource := range parsed.Catalogs {
		if id, ok := catalogSource["id"].(string); ok && id == catalogId {
			found = true

			if payload.Name != "" {
				catalogSource["name"] = payload.Name
			}
			if len(payload.Labels) > 0 {
				catalogSource["labels"] = payload.Labels
			}
			if payload.Enabled != nil {
				catalogSource["enabled"] = *payload.Enabled
			}

			properties, _ := catalogSource["properties"].(map[string]interface{})
			if properties == nil {
				properties = make(map[string]interface{})
			}

			if len(payload.IncludedModels) > 0 {
				properties["includedModels"] = payload.IncludedModels
			}
			if len(payload.ExcludedModels) > 0 {
				properties["excludedModels"] = payload.ExcludedModels
			}
			if payload.AllowedOrganization != nil {
				properties["allowedOrganization"] = *payload.AllowedOrganization
			}
			if secretName != "" {
				properties["apiKey"] = secretName
			}
			if yamlFilePath != "" && payload.Yaml != nil {
				properties["yamlCatalogPath"] = yamlFilePath
			}

			if len(properties) > 0 {
				catalogSource["properties"] = properties
			}

			parsed.Catalogs[i] = catalogSource
			break
		}
	}

	if !found {
		return "", fmt.Errorf("catalog '%s' not found in yaml", catalogId)
	}

	updatedBytes, err := yaml.Marshal(parsed)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated sources.yaml: %w", err)
	}

	return string(updatedBytes), nil
}

func buildOverrideEntryForDefaultSource(catalogId string, payload models.CatalogSourceConfigPayload) map[string]interface{} {
	entry := map[string]interface{}{
		"id": catalogId,
	}

	if payload.Enabled != nil {
		entry["enabled"] = *payload.Enabled
	}

	properties := make(map[string]interface{})

	if len(payload.IncludedModels) > 0 {
		properties["includedModels"] = payload.IncludedModels
	}
	if len(payload.ExcludedModels) > 0 {
		properties["excludedModels"] = payload.ExcludedModels
	}

	if len(properties) > 0 {
		entry["properties"] = properties
	}

	return entry
}
