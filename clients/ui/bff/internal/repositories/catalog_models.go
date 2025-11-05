package repositories

import (
	"encoding/json"
	"fmt"
	"net/url"

	catalogOpenapi "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
)

const catalogModelsPath = "/models"

type CatalogModelsInterface interface {
	GetAllCatalogModelsAcrossSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*catalogOpenapi.CatalogModelList, error)
}

type CatalogModels struct {
	CatalogModelsInterface
}

func (a CatalogModels) GetAllCatalogModelsAcrossSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*catalogOpenapi.CatalogModelList, error) {
	responseData, err := client.GET(UrlWithPageParams(catalogModelsPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var models catalogOpenapi.CatalogModelList

	if err := json.Unmarshal(responseData, &models); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &models, nil
}
