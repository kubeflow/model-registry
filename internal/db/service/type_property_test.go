package service_test

import (
	"fmt"
	"testing"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypePropertyRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	repo := service.NewTypePropertyRepository(sharedDB)
	typeRepo := service.NewTypeRepository(sharedDB)

	// Create a test type to use for properties
	testTypeName := "test-type-for-properties"
	testTypeKind := int32(1)
	testType := &models.TypeImpl{
		Attributes: &models.TypeAttributes{
			Name:     &testTypeName,
			TypeKind: &testTypeKind,
		},
	}

	savedType, err := typeRepo.Save(testType)
	require.NoError(t, err)
	require.NotNil(t, savedType)
	typeID := *savedType.GetID()

	t.Run("TestSave", func(t *testing.T) {
		// Test saving a new type property
		propertyName := "test-property"
		dataType := int32(1)

		property := &models.TypePropertyImpl{
			TypeID:   typeID,
			Name:     propertyName,
			DataType: &dataType,
		}

		savedProperty, err := repo.Save(property)
		require.NoError(t, err)
		require.NotNil(t, savedProperty)

		// Verify the saved property
		assert.Equal(t, typeID, savedProperty.GetTypeID())
		assert.Equal(t, propertyName, savedProperty.GetName())
		assert.Equal(t, dataType, *savedProperty.GetDataType())
	})

	t.Run("TestSaveExisting", func(t *testing.T) {
		// Test saving a property that already exists with same data type
		propertyName := "test-existing-property"
		dataType := int32(2)

		// First, save the property
		firstProperty := &models.TypePropertyImpl{
			TypeID:   typeID,
			Name:     propertyName,
			DataType: &dataType,
		}

		savedProperty1, err := repo.Save(firstProperty)
		require.NoError(t, err)
		require.NotNil(t, savedProperty1)

		// Now try to save the same property again
		secondProperty := &models.TypePropertyImpl{
			TypeID:   typeID,
			Name:     propertyName,
			DataType: &dataType,
		}

		savedProperty2, err := repo.Save(secondProperty)
		require.NoError(t, err)
		require.NotNil(t, savedProperty2)

		// Should return the existing property
		assert.Equal(t, savedProperty1.GetTypeID(), savedProperty2.GetTypeID())
		assert.Equal(t, savedProperty1.GetName(), savedProperty2.GetName())
		assert.Equal(t, *savedProperty1.GetDataType(), *savedProperty2.GetDataType())
	})

	t.Run("TestSaveInvalidDataTypeChange", func(t *testing.T) {
		// Test trying to save a property with different data type than existing
		propertyName := "test-datatype-change-property"
		originalDataType := int32(1)
		newDataType := int32(2)

		// First, save the property with original data type
		originalProperty := &models.TypePropertyImpl{
			TypeID:   typeID,
			Name:     propertyName,
			DataType: &originalDataType,
		}

		_, err := repo.Save(originalProperty)
		require.NoError(t, err)

		// Now try to save with different data type - should fail
		changedProperty := &models.TypePropertyImpl{
			TypeID:   typeID,
			Name:     propertyName,
			DataType: &newDataType,
		}

		_, err = repo.Save(changedProperty)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot change to")
	})

	t.Run("TestSaveMultiplePropertiesForSameType", func(t *testing.T) {
		// Test saving multiple properties for the same type
		properties := []struct {
			name     string
			dataType int32
		}{
			{"property1", 1},
			{"property2", 2},
			{"property3", 3},
		}

		for _, prop := range properties {
			property := &models.TypePropertyImpl{
				TypeID:   typeID,
				Name:     prop.name,
				DataType: &prop.dataType,
			}

			savedProperty, err := repo.Save(property)
			require.NoError(t, err, "Failed to save property %s", prop.name)
			require.NotNil(t, savedProperty)

			assert.Equal(t, typeID, savedProperty.GetTypeID())
			assert.Equal(t, prop.name, savedProperty.GetName())
			assert.Equal(t, prop.dataType, *savedProperty.GetDataType())
		}
	})

	t.Run("TestSaveWithDifferentTypes", func(t *testing.T) {
		// Create another test type
		anotherTypeName := "another-test-type"
		anotherTypeKind := int32(2)
		anotherType := &models.TypeImpl{
			Attributes: &models.TypeAttributes{
				Name:     &anotherTypeName,
				TypeKind: &anotherTypeKind,
			},
		}

		savedAnotherType, err := typeRepo.Save(anotherType)
		require.NoError(t, err)
		require.NotNil(t, savedAnotherType)
		anotherTypeID := *savedAnotherType.GetID()

		// Test saving properties with same name for different types
		propertyName := "shared-property-name"
		dataType := int32(1)

		// Save property for first type
		property1 := &models.TypePropertyImpl{
			TypeID:   typeID,
			Name:     propertyName,
			DataType: &dataType,
		}

		savedProperty1, err := repo.Save(property1)
		require.NoError(t, err)
		require.NotNil(t, savedProperty1)

		// Save property with same name for second type - should work
		property2 := &models.TypePropertyImpl{
			TypeID:   anotherTypeID,
			Name:     propertyName,
			DataType: &dataType,
		}

		savedProperty2, err := repo.Save(property2)
		require.NoError(t, err)
		require.NotNil(t, savedProperty2)

		// Verify both properties exist with different type IDs
		assert.Equal(t, typeID, savedProperty1.GetTypeID())
		assert.Equal(t, anotherTypeID, savedProperty2.GetTypeID())
		assert.Equal(t, propertyName, savedProperty1.GetName())
		assert.Equal(t, propertyName, savedProperty2.GetName())
		assert.Equal(t, dataType, *savedProperty1.GetDataType())
		assert.Equal(t, dataType, *savedProperty2.GetDataType())
	})

	t.Run("TestSaveNilDataType", func(t *testing.T) {
		// Test saving a property with nil data type for a new property
		propertyName := "test-nil-datatype-property"

		property := &models.TypePropertyImpl{
			TypeID:   typeID,
			Name:     propertyName,
			DataType: nil,
		}

		savedProperty, err := repo.Save(property)
		require.NoError(t, err)
		require.NotNil(t, savedProperty)

		// Verify the saved property
		assert.Equal(t, typeID, savedProperty.GetTypeID())
		assert.Equal(t, propertyName, savedProperty.GetName())
		assert.Nil(t, savedProperty.GetDataType())
	})

	t.Run("TestSaveValidDataTypes", func(t *testing.T) {
		// Test saving properties with various valid data types
		validDataTypes := []int32{0, 1, 2, 3, 4, 5, 10, 100}

		for i, dataType := range validDataTypes {
			propertyName := fmt.Sprintf("test-datatype-%d-property", i)
			property := &models.TypePropertyImpl{
				TypeID:   typeID,
				Name:     propertyName,
				DataType: &dataType,
			}

			savedProperty, err := repo.Save(property)
			require.NoError(t, err, "Failed to save property with data type %d", dataType)
			require.NotNil(t, savedProperty)

			assert.Equal(t, typeID, savedProperty.GetTypeID())
			assert.Equal(t, propertyName, savedProperty.GetName())
			assert.Equal(t, dataType, *savedProperty.GetDataType())
		}
	})
}
