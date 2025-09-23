package service_test

import (
	"testing"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	repo := service.NewTypeRepository(sharedDB)

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
			defaults.ExperimentTypeName,
			defaults.ExperimentRunTypeName,
			defaults.DataSetTypeName,
			defaults.MetricTypeName,
			defaults.ParameterTypeName,
			defaults.MetricHistoryTypeName,
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

	t.Run("TestSave", func(t *testing.T) {
		// Test saving a new type
		newTypeName := "test-custom-type"
		newTypeKind := int32(1)
		newVersion := "1.0.0"
		newDescription := "Test custom type description"
		newInputType := "application/json"
		newOutputType := "application/json"
		newExternalID := "external-123"

		newType := &models.TypeImpl{
			Attributes: &models.TypeAttributes{
				Name:        &newTypeName,
				TypeKind:    &newTypeKind,
				Version:     &newVersion,
				Description: &newDescription,
				InputType:   &newInputType,
				OutputType:  &newOutputType,
				ExternalID:  &newExternalID,
			},
		}

		savedType, err := repo.Save(newType)
		require.NoError(t, err)
		require.NotNil(t, savedType)

		// Verify the saved type has an ID
		assert.NotNil(t, savedType.GetID())
		assert.Greater(t, *savedType.GetID(), int32(0))

		// Verify all attributes are preserved
		attrs := savedType.GetAttributes()
		assert.Equal(t, newTypeName, *attrs.Name)
		assert.Equal(t, newTypeKind, *attrs.TypeKind)
		assert.Equal(t, newVersion, *attrs.Version)
		assert.Equal(t, newDescription, *attrs.Description)
		assert.Equal(t, newInputType, *attrs.InputType)
		assert.Equal(t, newOutputType, *attrs.OutputType)
		assert.Equal(t, newExternalID, *attrs.ExternalID)
	})

	t.Run("TestSaveExisting", func(t *testing.T) {
		// Test saving a type that already exists
		existingTypeName := "test-existing-type"
		existingTypeKind := int32(2)

		// First, save the type
		firstType := &models.TypeImpl{
			Attributes: &models.TypeAttributes{
				Name:     &existingTypeName,
				TypeKind: &existingTypeKind,
			},
		}

		savedType1, err := repo.Save(firstType)
		require.NoError(t, err)
		require.NotNil(t, savedType1)

		// Now try to save the same type again
		secondType := &models.TypeImpl{
			Attributes: &models.TypeAttributes{
				Name:     &existingTypeName,
				TypeKind: &existingTypeKind,
			},
		}

		savedType2, err := repo.Save(secondType)
		require.NoError(t, err)
		require.NotNil(t, savedType2)

		// Should return the existing type with same ID
		assert.Equal(t, *savedType1.GetID(), *savedType2.GetID())
		assert.Equal(t, *savedType1.GetAttributes().Name, *savedType2.GetAttributes().Name)
		assert.Equal(t, *savedType1.GetAttributes().TypeKind, *savedType2.GetAttributes().TypeKind)
	})

	t.Run("TestSaveInvalidTypeKindChange", func(t *testing.T) {
		// Test trying to save a type with different kind than existing
		typeName := "test-kind-change-type"
		originalKind := int32(1)
		newKind := int32(2)

		// First, save the type with original kind
		originalType := &models.TypeImpl{
			Attributes: &models.TypeAttributes{
				Name:     &typeName,
				TypeKind: &originalKind,
			},
		}

		_, err := repo.Save(originalType)
		require.NoError(t, err)

		// Now try to save with different kind - should fail
		changedType := &models.TypeImpl{
			Attributes: &models.TypeAttributes{
				Name:     &typeName,
				TypeKind: &newKind,
			},
		}

		_, err = repo.Save(changedType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot change to kind")
	})

	t.Run("TestSaveValidationErrors", func(t *testing.T) {
		// Test saving type without attributes
		emptyType := &models.TypeImpl{}
		_, err := repo.Save(emptyType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing attributes")

		// Test saving type without name
		noNameType := &models.TypeImpl{
			Attributes: &models.TypeAttributes{
				TypeKind: apiutils.Of(int32(1)),
			},
		}
		_, err = repo.Save(noNameType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing name")

		// Test saving type without kind
		noKindType := &models.TypeImpl{
			Attributes: &models.TypeAttributes{
				Name: apiutils.Of("test-no-kind"),
			},
		}
		_, err = repo.Save(noKindType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing kind")
	})

	t.Run("TestSaveMinimalType", func(t *testing.T) {
		// Test saving type with only required fields
		typeName := "test-minimal-type"
		typeKind := int32(3)

		minimalType := &models.TypeImpl{
			Attributes: &models.TypeAttributes{
				Name:     &typeName,
				TypeKind: &typeKind,
			},
		}

		savedType, err := repo.Save(minimalType)
		require.NoError(t, err)
		require.NotNil(t, savedType)

		// Verify required fields
		assert.NotNil(t, savedType.GetID())
		assert.Equal(t, typeName, *savedType.GetAttributes().Name)
		assert.Equal(t, typeKind, *savedType.GetAttributes().TypeKind)

		// Optional fields should be nil
		attrs := savedType.GetAttributes()
		assert.Nil(t, attrs.Version)
		assert.Nil(t, attrs.Description)
		assert.Nil(t, attrs.InputType)
		assert.Nil(t, attrs.OutputType)
		assert.Nil(t, attrs.ExternalID)
	})
}
