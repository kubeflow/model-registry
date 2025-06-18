package service_test

import (
	"testing"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := service.NewTypeRepository(db)

	t.Run("TestGetAll", func(t *testing.T) {
		// Test retrieving all types
		types, err := repo.GetAll()
		require.NoError(t, err)
		require.NotNil(t, types)

		// Verify we get the expected number of types (should be at least the default types)
		assert.GreaterOrEqual(t, len(types), 7, "Should have at least 7 default types")

		// Create a map for easier lookup
		typeMap := make(map[string]models.Type)
		for _, typeItem := range types {
			typeMap[*typeItem.GetAttributes().Name] = typeItem
		}

		// Verify all expected default types are present
		expectedTypes := []string{
			defaults.RegisteredModelTypeName,
			defaults.ModelVersionTypeName,
			defaults.ServingEnvironmentTypeName,
			defaults.InferenceServiceTypeName,
			defaults.ServeModelTypeName,
			defaults.ModelArtifactTypeName,
			defaults.DocArtifactTypeName,
		}

		for _, expectedType := range expectedTypes {
			foundType, exists := typeMap[expectedType]
			assert.True(t, exists, "Should find type: %s", expectedType)
			if exists {
				// Verify the type has required fields
				assert.NotNil(t, foundType.GetID(), "Type %s should have an ID", expectedType)
				assert.NotNil(t, foundType.GetAttributes(), "Type %s should have attributes", expectedType)
				assert.NotNil(t, foundType.GetAttributes().Name, "Type %s should have a name", expectedType)
				assert.Equal(t, expectedType, *foundType.GetAttributes().Name, "Type name should match")

				// Verify ID is positive
				assert.Greater(t, *foundType.GetID(), int32(0), "Type %s should have positive ID", expectedType)
			}
		}
	})

	t.Run("TestGetAllStructure", func(t *testing.T) {
		// Test the structure and content of returned types
		types, err := repo.GetAll()
		require.NoError(t, err)
		require.NotEmpty(t, types)

		for _, typeItem := range types {
			// Verify each type has the required structure
			assert.NotNil(t, typeItem.GetID(), "Each type should have an ID")
			assert.NotNil(t, typeItem.GetAttributes(), "Each type should have attributes")

			attrs := typeItem.GetAttributes()
			assert.NotNil(t, attrs.Name, "Each type should have a name")
			assert.NotEmpty(t, *attrs.Name, "Type name should not be empty")
		}
	})

	t.Run("TestGetAllConsistency", func(t *testing.T) {
		// Test that multiple calls return consistent results
		types1, err1 := repo.GetAll()
		require.NoError(t, err1)

		types2, err2 := repo.GetAll()
		require.NoError(t, err2)

		// Should return the same number of types
		assert.Equal(t, len(types1), len(types2), "Multiple calls should return same number of types")

		// Create maps for comparison
		typeMap1 := make(map[int32]string)
		typeMap2 := make(map[int32]string)

		for _, typeItem := range types1 {
			typeMap1[*typeItem.GetID()] = *typeItem.GetAttributes().Name
		}

		for _, typeItem := range types2 {
			typeMap2[*typeItem.GetID()] = *typeItem.GetAttributes().Name
		}

		// Should have the same types with same IDs
		assert.Equal(t, typeMap1, typeMap2, "Multiple calls should return consistent results")
	})

	t.Run("TestGetAllSpecificTypes", func(t *testing.T) {
		// Test specific type details for known types
		types, err := repo.GetAll()
		require.NoError(t, err)

		// Find RegisteredModel type and verify its properties
		var registeredModelType models.Type
		for _, typeItem := range types {
			if *typeItem.GetAttributes().Name == defaults.RegisteredModelTypeName {
				registeredModelType = typeItem
				break
			}
		}

		require.NotNil(t, registeredModelType, "Should find RegisteredModel type")
		assert.Equal(t, defaults.RegisteredModelTypeName, *registeredModelType.GetAttributes().Name)
		assert.NotNil(t, registeredModelType.GetID())

		// Find ModelArtifact type and verify its properties
		var modelArtifactType models.Type
		for _, typeItem := range types {
			if *typeItem.GetAttributes().Name == defaults.ModelArtifactTypeName {
				modelArtifactType = typeItem
				break
			}
		}

		require.NotNil(t, modelArtifactType, "Should find ModelArtifact type")
		assert.Equal(t, defaults.ModelArtifactTypeName, *modelArtifactType.GetAttributes().Name)
		assert.NotNil(t, modelArtifactType.GetID())

		// Verify different types have different IDs
		assert.NotEqual(t, *registeredModelType.GetID(), *modelArtifactType.GetID(), "Different types should have different IDs")
	})

	t.Run("TestGetAllOptionalFields", func(t *testing.T) {
		// Test that optional fields are handled correctly
		types, err := repo.GetAll()
		require.NoError(t, err)

		for _, typeItem := range types {
			attrs := typeItem.GetAttributes()

			// Name is required and should always be present
			assert.NotNil(t, attrs.Name, "Name should always be present")

			// Optional fields may be nil, but if present should be valid
			if attrs.Version != nil {
				// Version can be empty string, but if present should be a string
				assert.IsType(t, "", *attrs.Version, "Version should be string if present")
			}

			if attrs.Description != nil {
				// Description can be empty, but should be a string
				assert.IsType(t, "", *attrs.Description, "Description should be string if present")
			}

			if attrs.ExternalID != nil {
				// ExternalID should be a string if present
				assert.IsType(t, "", *attrs.ExternalID, "ExternalID should be string if present")
			}

			if attrs.TypeKind != nil {
				// TypeKind should be a valid integer if present
				assert.IsType(t, int32(0), *attrs.TypeKind, "TypeKind should be int32 if present")
			}
		}
	})

	t.Run("TestGetAllEmptyDatabase", func(t *testing.T) {
		// This test verifies behavior with the migrated database
		// Since migrations create default types, we should always have types
		types, err := repo.GetAll()
		require.NoError(t, err)

		// Even with a "fresh" database, migrations should have created default types
		assert.NotEmpty(t, types, "Migrated database should have default types")
		assert.GreaterOrEqual(t, len(types), 1, "Should have at least one type after migration")
	})
}
