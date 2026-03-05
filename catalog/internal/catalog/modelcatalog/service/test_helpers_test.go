package service

import (
	"errors"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	testCatalogModelTypeName           = "kf.CatalogModel"
	testCatalogModelArtifactTypeName   = "kf.CatalogModelArtifact"
	testCatalogMetricsArtifactTypeName = "kf.CatalogMetricsArtifact"
)

// testDatastoreSpec returns a minimal datastore spec for modelcatalog tests.
// This avoids importing catalog/internal/db/service which would cause an import cycle.
func testDatastoreSpec() *datastore.Spec {
	return datastore.NewSpec().
		AddContext(testCatalogModelTypeName, datastore.NewSpecType(NewCatalogModelRepository).
			AddString("source_id").
			AddString("description").
			AddString("owner").
			AddString("state").
			AddStruct("language").
			AddString("library_name").
			AddString("license_link").
			AddString("license").
			AddString("logo").
			AddString("maturity").
			AddString("provider").
			AddString("readme").
			AddStruct("tasks"),
		).
		AddArtifact(testCatalogModelArtifactTypeName, datastore.NewSpecType(NewCatalogModelArtifactRepository).
			AddString("uri"),
		).
		AddArtifact(testCatalogMetricsArtifactTypeName, datastore.NewSpecType(NewCatalogMetricsArtifactRepository).
			AddString("metricsType"),
		)
}

// getCatalogModelTypeID gets the CatalogModel type ID from the database
func getCatalogModelTypeID(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", testCatalogModelTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModel type")
	}

	return typeRecord.ID
}

// getCatalogModelArtifactTypeID gets the CatalogModelArtifact type ID from the database
func getCatalogModelArtifactTypeID(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", testCatalogModelArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModelArtifact type")
	}

	return typeRecord.ID
}

// getCatalogMetricsArtifactTypeID gets the CatalogMetricsArtifact type ID from the database
func getCatalogMetricsArtifactTypeID(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", testCatalogMetricsArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogMetricsArtifact type")
	}

	return typeRecord.ID
}

// getOrCreateExecutionTypeID creates or gets an execution type ID for testing
func getOrCreateExecutionTypeID(t *testing.T, db *gorm.DB) int32 {
	typeName := "test.Execution"
	var typeRecord schema.Type
	err := db.Where("name = ?", typeName).First(&typeRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			typeRecord = schema.Type{
				Name: typeName,
			}
			err = db.Create(&typeRecord).Error
			require.NoError(t, err, "Failed to create test execution type")
		} else {
			require.NoError(t, err)
		}
	}
	return typeRecord.ID
}

// createTestArtifact creates a test artifact in the database
func createTestArtifact(t *testing.T, db *gorm.DB, typeID int32, name string) *schema.Artifact {
	artifact := &schema.Artifact{
		TypeID:                   typeID,
		Name:                     &name,
		CreateTimeSinceEpoch:     1000,
		LastUpdateTimeSinceEpoch: 1000,
	}
	err := db.Create(artifact).Error
	require.NoError(t, err, "Failed to create test artifact")
	return artifact
}

// createAttribution creates an attribution between context and artifact
func createAttribution(t *testing.T, db *gorm.DB, contextID, artifactID int32) {
	attribution := &schema.Attribution{
		ContextID:  contextID,
		ArtifactID: artifactID,
	}
	err := db.Create(attribution).Error
	require.NoError(t, err, "Failed to create attribution")
}

// createTestExecution creates a test execution in the database
func createTestExecution(t *testing.T, db *gorm.DB, typeID int32, name string) *schema.Execution {
	execution := &schema.Execution{
		TypeID:                   typeID,
		Name:                     &name,
		CreateTimeSinceEpoch:     1000,
		LastUpdateTimeSinceEpoch: 1000,
	}
	err := db.Create(execution).Error
	require.NoError(t, err, "Failed to create test execution")
	return execution
}

// createTestEvent creates a test event in the database
func createTestEvent(t *testing.T, db *gorm.DB, artifactID, executionID int32) *schema.Event {
	event := &schema.Event{
		ArtifactID:  artifactID,
		ExecutionID: executionID,
		Type:        1, // INPUT event type
	}
	err := db.Create(event).Error
	require.NoError(t, err, "Failed to create test event")
	return event
}

// findIndex finds the index of a string in a slice
func findIndex(slice []string, target string) int {
	for i, item := range slice {
		if item == target {
			return i
		}
	}
	return -1
}
