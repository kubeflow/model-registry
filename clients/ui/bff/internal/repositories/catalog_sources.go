package repositories

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

const sourcesPath = "/sources"

type CatalogSourcesInterface interface {
	GetAllCatalogSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.CatalogSourceList, error)
	GetCatalogSourceModel(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*models.CatalogModel, error)
	GetCatalogModelArtifacts(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*models.CatalogModelArtifactList, error)
}

type CatalogSources struct {
	CatalogSourcesInterface
}

func (a CatalogSources) GetAllCatalogSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.CatalogSourceList, error) {
	responseData, err := client.GET(UrlWithPageParams(sourcesPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var sources models.CatalogSourceList

	if err := json.Unmarshal(responseData, &sources); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &sources, nil
}

func (a CatalogSources) GetCatalogSourceModel(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*models.CatalogModel, error) {
	path, err := url.JoinPath(sourcesPath, sourceId, "models", modelName)

	if err != nil {
		return nil, err
	}
	responseData, err := client.GET(path)
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var catalogModel models.CatalogModel

	if err := json.Unmarshal(responseData, &catalogModel); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &catalogModel, nil
}

func (a CatalogSources) GetCatalogModelArtifacts(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*models.CatalogModelArtifactList, error) {
	path, err := url.JoinPath(sourcesPath, sourceId, "models", modelName, "artifacts")
	if err != nil {
		return nil, err
	}
	responseData, err := client.GET(path)
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var catalogModelArtifacts models.CatalogModelArtifactList

	if err := json.Unmarshal(responseData, &catalogModelArtifacts); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}
	return &catalogModelArtifacts, nil
}
