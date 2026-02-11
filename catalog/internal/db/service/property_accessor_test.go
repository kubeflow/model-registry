package service_test

import (
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/apiutils"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropertyAccessor_GetString(t *testing.T) {
	t.Run("ExistingProperty", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "key1", StringValue: apiutils.Of("value1")},
			{Name: "key2", StringValue: apiutils.Of("value2")},
		}

		pa := service.NewPropertyAccessor(props)
		assert.Equal(t, "value1", pa.GetString("key1"))
		assert.Equal(t, "value2", pa.GetString("key2"))
	})

	t.Run("NonExistentProperty", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "key1", StringValue: apiutils.Of("value1")},
		}

		pa := service.NewPropertyAccessor(props)
		assert.Equal(t, "", pa.GetString("nonexistent"))
	})

	t.Run("NilValue", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "key1", StringValue: nil},
		}

		pa := service.NewPropertyAccessor(props)
		assert.Equal(t, "", pa.GetString("key1"))
	})

	t.Run("NilPropertiesSlice", func(t *testing.T) {
		pa := service.NewPropertyAccessor(nil)
		assert.Equal(t, "", pa.GetString("any"))
	})

	t.Run("EmptyPropertiesSlice", func(t *testing.T) {
		props := &[]dbmodels.Properties{}
		pa := service.NewPropertyAccessor(props)
		assert.Equal(t, "", pa.GetString("any"))
	})
}

func TestPropertyAccessor_GetStringPtr(t *testing.T) {
	t.Run("ExistingProperty", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "key1", StringValue: apiutils.Of("value1")},
		}

		pa := service.NewPropertyAccessor(props)
		result := pa.GetStringPtr("key1")
		require.NotNil(t, result)
		assert.Equal(t, "value1", *result)
	})

	t.Run("NonExistentProperty_ReturnsNil", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "key1", StringValue: apiutils.Of("value1")},
		}

		pa := service.NewPropertyAccessor(props)
		assert.Nil(t, pa.GetStringPtr("nonexistent"))
	})

	t.Run("EmptyString_ReturnsNil", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "key1", StringValue: apiutils.Of("")},
		}

		pa := service.NewPropertyAccessor(props)
		assert.Nil(t, pa.GetStringPtr("key1"))
	})
}

func TestPropertyAccessor_GetBoolPtr(t *testing.T) {
	t.Run("TrueValue", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "enabled", BoolValue: apiutils.Of(true)},
		}

		pa := service.NewPropertyAccessor(props)
		result := pa.GetBoolPtr("enabled")
		require.NotNil(t, result)
		assert.True(t, *result)
	})

	t.Run("FalseValue", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "disabled", BoolValue: apiutils.Of(false)},
		}

		pa := service.NewPropertyAccessor(props)
		result := pa.GetBoolPtr("disabled")
		require.NotNil(t, result)
		assert.False(t, *result)
	})

	t.Run("NonExistentProperty", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "enabled", BoolValue: apiutils.Of(true)},
		}

		pa := service.NewPropertyAccessor(props)
		assert.Nil(t, pa.GetBoolPtr("nonexistent"))
	})
}

func TestPropertyAccessor_GetInt(t *testing.T) {
	t.Run("ExistingProperty", func(t *testing.T) {
		intVal := int32(42)
		props := &[]dbmodels.Properties{
			{Name: "count", IntValue: &intVal},
		}

		pa := service.NewPropertyAccessor(props)
		assert.Equal(t, int64(42), pa.GetInt("count"))
	})

	t.Run("NonExistentProperty_ReturnsZero", func(t *testing.T) {
		props := &[]dbmodels.Properties{}
		pa := service.NewPropertyAccessor(props)
		assert.Equal(t, int64(0), pa.GetInt("nonexistent"))
	})
}

func TestPropertyAccessor_GetStringArray(t *testing.T) {
	t.Run("ValidJSONArray", func(t *testing.T) {
		jsonStr := `["tag1","tag2","tag3"]`
		props := &[]dbmodels.Properties{
			{Name: "tags", StringValue: &jsonStr},
		}

		pa := service.NewPropertyAccessor(props)
		result := pa.GetStringArray("tags")
		require.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.Equal(t, "tag1", result[0])
		assert.Equal(t, "tag2", result[1])
		assert.Equal(t, "tag3", result[2])
	})

	t.Run("EmptyArray", func(t *testing.T) {
		jsonStr := `[]`
		props := &[]dbmodels.Properties{
			{Name: "tags", StringValue: &jsonStr},
		}

		pa := service.NewPropertyAccessor(props)
		result := pa.GetStringArray("tags")
		require.NotNil(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("InvalidJSON_ReturnsNil", func(t *testing.T) {
		jsonStr := `{invalid json}`
		props := &[]dbmodels.Properties{
			{Name: "tags", StringValue: &jsonStr},
		}

		pa := service.NewPropertyAccessor(props)
		result := pa.GetStringArray("tags")
		assert.Nil(t, result)
	})

	t.Run("NonExistentProperty_ReturnsNil", func(t *testing.T) {
		props := &[]dbmodels.Properties{}
		pa := service.NewPropertyAccessor(props)
		assert.Nil(t, pa.GetStringArray("tags"))
	})

	t.Run("EmptyString_ReturnsNil", func(t *testing.T) {
		emptyStr := ""
		props := &[]dbmodels.Properties{
			{Name: "tags", StringValue: &emptyStr},
		}

		pa := service.NewPropertyAccessor(props)
		assert.Nil(t, pa.GetStringArray("tags"))
	})
}

func TestPropertyAccessor_HasAny(t *testing.T) {
	t.Run("OnePropertyExists", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "prop1", StringValue: apiutils.Of("value1")},
			{Name: "prop2", StringValue: apiutils.Of("value2")},
		}

		pa := service.NewPropertyAccessor(props)
		assert.True(t, pa.HasAny("prop1", "nonexistent"))
	})

	t.Run("MultiplePropertiesExist", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "prop1", StringValue: apiutils.Of("value1")},
			{Name: "prop2", StringValue: apiutils.Of("value2")},
		}

		pa := service.NewPropertyAccessor(props)
		assert.True(t, pa.HasAny("prop1", "prop2"))
	})

	t.Run("NoPropertiesExist", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "prop1", StringValue: apiutils.Of("value1")},
		}

		pa := service.NewPropertyAccessor(props)
		assert.False(t, pa.HasAny("nonexistent1", "nonexistent2"))
	})

	t.Run("EmptyNamesList", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "prop1", StringValue: apiutils.Of("value1")},
		}

		pa := service.NewPropertyAccessor(props)
		assert.False(t, pa.HasAny())
	})
}

func TestPropertyAccessor_MultipleAccesses(t *testing.T) {
	t.Run("CanAccessSamePropertyMultipleTimes", func(t *testing.T) {
		props := &[]dbmodels.Properties{
			{Name: "key1", StringValue: apiutils.Of("value1")},
		}

		pa := service.NewPropertyAccessor(props)

		// Access the same property multiple times
		assert.Equal(t, "value1", pa.GetString("key1"))
		assert.Equal(t, "value1", pa.GetString("key1"))
		assert.Equal(t, "value1", pa.GetString("key1"))
	})

	t.Run("MixedPropertyTypes", func(t *testing.T) {
		jsonArray := `["a","b","c"]`
		intVal := int32(100)
		props := &[]dbmodels.Properties{
			{Name: "stringProp", StringValue: apiutils.Of("text")},
			{Name: "boolProp", BoolValue: apiutils.Of(true)},
			{Name: "intProp", IntValue: &intVal},
			{Name: "arrayProp", StringValue: &jsonArray},
		}

		pa := service.NewPropertyAccessor(props)

		assert.Equal(t, "text", pa.GetString("stringProp"))
		assert.True(t, *pa.GetBoolPtr("boolProp"))
		assert.Equal(t, int64(100), pa.GetInt("intProp"))
		assert.Len(t, pa.GetStringArray("arrayProp"), 3)
	})
}
