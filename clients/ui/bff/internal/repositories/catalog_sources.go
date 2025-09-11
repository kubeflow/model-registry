package repositories

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
)

const sourcesPath = "/sources"

type CatalogSourcesInterface interface {
	GetAllCatalogSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*openapi.CatalogSourceList, error)
	GetCatalogSourceModel(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*openapi.CatalogModel, error)
	GetCatalogModelArtifacts(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*openapi.CatalogModelArtifactList, error)
}

type CatalogSources struct {
	CatalogSourcesInterface
}

func (a CatalogSources) GetAllCatalogSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*openapi.CatalogSourceList, error) {
	responseData, err := client.GET(UrlWithPageParams(sourcesPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var sources openapi.CatalogSourceList

	if err := json.Unmarshal(responseData, &sources); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &sources, nil
}

func (a CatalogSources) GetCatalogSourceModel(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*openapi.CatalogModel, error) {
	path, err := url.JoinPath(sourcesPath, sourceId, "models", modelName)

	if err != nil {
		return nil, err
	}
	responseData, err := client.GET(path)
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var catalogModel openapi.CatalogModel

	if err := json.Unmarshal(responseData, &catalogModel); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &catalogModel, nil
}

func (a CatalogSources) GetCatalogModelArtifacts(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*openapi.CatalogModelArtifactList, error) {
	path, err := url.JoinPath(sourcesPath, sourceId, "models", modelName, "artifacts")
	if err != nil {
		return nil, err
	}
	responseData, err := client.GET(path)
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var catalogModelArtifacts openapi.CatalogModelArtifactList

	if err := json.Unmarshal(responseData, &catalogModelArtifacts); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}
	return &catalogModelArtifacts, nil
}
