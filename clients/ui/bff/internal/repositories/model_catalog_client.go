package repositories

import (
	"log/slog"
)

type ModelCatalogClientInterface interface {
	CatalogSourcesInterface
	CatalogModelsInterface
	CatalogSourcePreviewInterface
}

type ModelCatalogClient struct {
	logger *slog.Logger
	CatalogSources
	CatalogModels
	CatalogSourcePreview
}

func NewModelCatalogClient(logger *slog.Logger) (ModelCatalogClientInterface, error) {
	return &ModelCatalogClient{logger: logger}, nil
}
