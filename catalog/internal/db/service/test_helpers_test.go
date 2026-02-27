package service_test

import (
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// getCatalogModelTypeID gets the CatalogModel type ID from the database
func getCatalogModelTypeID(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModel type")
	}

	return typeRecord.ID
}

// getCatalogModelArtifactTypeID gets the CatalogModelArtifact type ID from the database
func getCatalogModelArtifactTypeID(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModelArtifact type")
	}

	return typeRecord.ID
}

// getCatalogMetricsArtifactTypeID gets the CatalogMetricsArtifact type ID from the database
func getCatalogMetricsArtifactTypeID(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogMetricsArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogMetricsArtifact type")
	}

	return typeRecord.ID
}
