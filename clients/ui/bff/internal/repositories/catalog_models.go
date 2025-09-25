package repositories

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

const catalogModelsPath = "/models"

type CatalogModelsInterface interface {
	GetAllCatalogModelsAcrossSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.CatalogModelList, error)
}

type CatalogModels struct {
	CatalogModelsInterface
}

func (a CatalogModels) GetAllCatalogModelsAcrossSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.CatalogModelList, error) {
	responseData, err := client.GET(UrlWithPageParams(catalogModelsPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var models models.CatalogModelList

	if err := json.Unmarshal(responseData, &models); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &models, nil
}
