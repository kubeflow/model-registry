package modelcatalog

import (
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// GetCatalogModelTypeIDForDBTest retrieves the CatalogModel type ID for testing
func GetCatalogModelTypeIDForDBTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModel type")
	}
	return typeRecord.ID
}

// GetCatalogModelArtifactTypeIDForDBTest retrieves the CatalogModelArtifact type ID for testing
func GetCatalogModelArtifactTypeIDForDBTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModelArtifact type")
	}
	return typeRecord.ID
}

// GetCatalogMetricsArtifactTypeIDForDBTest retrieves the CatalogMetricsArtifact type ID for testing
func GetCatalogMetricsArtifactTypeIDForDBTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogMetricsArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogMetricsArtifact type")
	}
	return typeRecord.ID
}

// GetCatalogSourceTypeIDForDBTest retrieves the CatalogSource type ID for testing
func GetCatalogSourceTypeIDForDBTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogSourceTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogSource type")
	}
	return typeRecord.ID
}
