package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/url"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type CatalogSourcePreviewInterface interface {
	CreateCatalogSourcePreview(client httpclient.HTTPClientInterface, sourcePreviewPayload models.CatalogSourcePreviewRequest, pageValues url.Values) (*models.CatalogSourcePreviewResult, error)
}

type CatalogSourcePreview struct {
	CatalogSourcePreviewInterface
}

func (a CatalogSourcePreview) CreateCatalogSourcePreview(client httpclient.HTTPClientInterface, sourcePreviewPayload models.CatalogSourcePreviewRequest, pageValues url.Values) (*models.CatalogSourcePreviewResult, error) {
	path, err := url.JoinPath(sourcesPath, "preview")
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	configPart, err := writer.CreateFormFile("config", "config.json")
	if err != nil {
		return nil, fmt.Errorf("error creating config form file: %w", err)
	}

	configData := map[string]interface{}{
		"type":           sourcePreviewPayload.Type,
		"includedModels": sourcePreviewPayload.IncludedModels,
		"excludedModels": sourcePreviewPayload.ExcludedModels,
	}

	properties := make(map[string]interface{})
	var yamlContent string
	var hasYamlContent bool

	if sourcePreviewPayload.Type == CatalogTypeHuggingFace {
		if org, ok := sourcePreviewPayload.Properties["allowedOrganization"]; ok {
			properties["allowedOrganization"] = org
		}
		if apiKey, ok := sourcePreviewPayload.Properties["apiKey"]; ok {
			properties["apiKey"] = apiKey
		}
	}

	if sourcePreviewPayload.Type == CatalogTypeYaml {
		// Check if yaml content is provided
		if content, ok := sourcePreviewPayload.Properties["yaml"].(string); ok && content != "" {
			yamlContent = content
			hasYamlContent = true
		} else {
			// If no yaml content, check for yamlCatalogPath
			if yamlPath, ok := sourcePreviewPayload.Properties["yamlCatalogPath"].(string); ok && yamlPath != "" {
				properties["yamlCatalogPath"] = yamlPath
			} else {
				return nil, fmt.Errorf("for yaml catalog type, either 'yaml' content or 'yamlCatalogPath' must be provided")
			}
		}
	}

	if len(properties) > 0 {
		configData["properties"] = properties
	}

	configJSON, err := json.Marshal(configData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling config: %w", err)
	}

	if _, err := configPart.Write(configJSON); err != nil {
		return nil, fmt.Errorf("error writing config data: %w", err)
	}

	if hasYamlContent {
		catalogDataPart, err := writer.CreateFormFile("catalogData", "catalog.yaml")
		if err != nil {
			return nil, fmt.Errorf("error creating catalogData form file: %w", err)
		}
		if _, err := catalogDataPart.Write([]byte(yamlContent)); err != nil {
			return nil, fmt.Errorf("error writing catalog data: %w", err)
		}
	}

	writer.Close()

	responseData, err := client.POSTWithContentType(UrlWithPageParams(path, pageValues), &body, writer.FormDataContentType())
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var sourcePreview models.CatalogSourcePreviewResult

	if err := json.Unmarshal(responseData, &sourcePreview); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &sourcePreview, nil
}
