package repositories

import (
	"encoding/json"
	"fmt"
	"net/url"

	catalogOpenapi "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
)

const sourcesPath = "/sources"
const filterOptionPath = "/models/filter_options"

type CatalogSourcesInterface interface {
	GetAllCatalogSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*catalogOpenapi.CatalogSourceList, error)
	GetCatalogSourceModel(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*catalogOpenapi.CatalogModel, error)
	GetCatalogModelArtifacts(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*catalogOpenapi.CatalogArtifactList, error)
	GetCatalogFilterOptions(client httpclient.HTTPClientInterface) (*catalogOpenapi.FilterOptionsList, error)
}

type CatalogSources struct {
	CatalogSourcesInterface
}

func (a CatalogSources) GetAllCatalogSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*catalogOpenapi.CatalogSourceList, error) {
	responseData, err := client.GET(UrlWithPageParams(sourcesPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var sources catalogOpenapi.CatalogSourceList

	if err := json.Unmarshal(responseData, &sources); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &sources, nil
}

func (a CatalogSources) GetCatalogSourceModel(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*catalogOpenapi.CatalogModel, error) {
	path, err := url.JoinPath(sourcesPath, sourceId, "models", modelName)

	if err != nil {
		return nil, err
	}
	responseData, err := client.GET(path)
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var catalogModel catalogOpenapi.CatalogModel

	if err := json.Unmarshal(responseData, &catalogModel); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &catalogModel, nil
}

func (a CatalogSources) GetCatalogModelArtifacts(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*catalogOpenapi.CatalogArtifactList, error) {
	path, err := url.JoinPath(sourcesPath, sourceId, "models", modelName, "artifacts")
	if err != nil {
		return nil, err
	}
	responseData, err := client.GET(path)
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var catalogModelArtifacts catalogOpenapi.CatalogArtifactList

	if err := json.Unmarshal(responseData, &catalogModelArtifacts); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}
	return &catalogModelArtifacts, nil
}

func (a CatalogSources) GetCatalogFilterOptions(client httpclient.HTTPClientInterface) (*catalogOpenapi.FilterOptionsList, error) {
	responseData, err := client.GET(filterOptionPath)

	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var sources catalogOpenapi.FilterOptionsList

	if err := json.Unmarshal(responseData, &sources); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &sources, nil
}
