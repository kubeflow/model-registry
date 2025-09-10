package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"net/url"
)

const catalogModelsPath = "/models"

type CatalogModelsInterface interface {
	GetAllCatalogModelsAcrossSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*openapi.CatalogModelList, error)
}

type CatalogModels struct {
	CatalogModelsInterface
}

func (a CatalogModels) GetAllCatalogModelsAcrossSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*openapi.CatalogModelList, error) {
	responseData, err := client.GET(UrlWithPageParams(catalogModelsPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var models openapi.CatalogModelList

	if err := json.Unmarshal(responseData, &models); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &models, nil
}
