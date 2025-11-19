package repositories

import (
	"fmt"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type CatalogSourcePreviewInterface interface {
	CreateCatalogSourcePreview(client httpclient.HTTPClientInterface, sourcePreviewPaylod models.CatalogSourcePreviewRequest) (*models.CatalogSourcePreviewResult, error)
}

type CatalogSourcePreview struct {
	CatalogSourcePreviewInterface
}

func (a CatalogSourcePreview) CreateCatalogSourcePreview(client httpclient.HTTPClientInterface, sourcePreviewPaylod models.CatalogSourcePreviewRequest) (*models.CatalogSourcePreviewResult, error) {
	return nil, fmt.Errorf("not implemented yet")
}
