package basecatalog

import (
	"testing"

	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommonContextProperties(t *testing.T) {
	props := CommonContextProperties()
	assert.Len(t, props, 5)
	for _, key := range []string{"id", "name", "externalId", "createTimeSinceEpoch", "lastUpdateTimeSinceEpoch"} {
		_, ok := props[key]
		assert.True(t, ok, "missing key %q", key)
	}
}

func TestCommonArtifactProperties(t *testing.T) {
	props := CommonArtifactProperties()
	assert.Len(t, props, 7)
	for _, key := range []string{"id", "name", "externalId", "createTimeSinceEpoch", "lastUpdateTimeSinceEpoch", "uri", "state"} {
		_, ok := props[key]
		assert.True(t, ok, "missing key %q", key)
	}
}

func TestCommonCatalogContextProperties(t *testing.T) {
	props := CommonCatalogContextProperties()
	assert.Len(t, props, 7)
	for _, key := range []string{"source_id", "description", "provider", "license", "license_link", "logo", "readme"} {
		_, ok := props[key]
		assert.True(t, ok, "missing key %q", key)
	}
}

func TestMergeProperties(t *testing.T) {
	base := map[string]filter.PropertyDefinition{
		"id":   {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
		"name": {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	}
	override := map[string]filter.PropertyDefinition{
		"name":  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "custom_name"},
		"extra": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "extra"},
	}

	merged := MergeProperties(base, override)

	assert.Len(t, merged, 3)
	// "id" from base
	assert.Equal(t, filter.EntityTable, merged["id"].Location)
	// "name" overridden by second map
	assert.Equal(t, filter.PropertyTable, merged["name"].Location)
	assert.Equal(t, "custom_name", merged["name"].Column)
	// "extra" from override
	assert.Equal(t, "extra", merged["extra"].Column)
}

func TestMergePropertiesEmpty(t *testing.T) {
	merged := MergeProperties()
	assert.Empty(t, merged)
}

func buildTestRegistry() *CatalogEntityRegistry {
	reg := NewCatalogEntityRegistry()
	reg.Register("TestContext", EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeContext,
		Properties: MergeProperties(
			CommonContextProperties(),
			map[string]filter.PropertyDefinition{
				"owner": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "owner"},
			},
		),
		RelatedEntityPrefix: "artifacts.",
		RelatedEntityType:   filter.RelatedEntityArtifact,
	})
	reg.Register("TestArtifact", EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeArtifact,
		Properties:     CommonArtifactProperties(),
	})
	reg.Register("TestChild", EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeArtifact,
		Properties:     CommonArtifactProperties(),
		IsChild:        true,
	})
	return reg
}

func TestGetMLMDEntityType(t *testing.T) {
	reg := buildTestRegistry()

	tests := []struct {
		name     string
		entity   filter.RestEntityType
		expected filter.EntityType
	}{
		{"registered context", "TestContext", filter.EntityTypeContext},
		{"registered artifact", "TestArtifact", filter.EntityTypeArtifact},
		{"unknown defaults to context", "UnknownEntity", filter.EntityTypeContext},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, reg.GetMLMDEntityType(tt.entity))
		})
	}
}

func TestGetPropertyDefinitionForRestEntity(t *testing.T) {
	reg := buildTestRegistry()

	t.Run("entity table property", func(t *testing.T) {
		def := reg.GetPropertyDefinitionForRestEntity("TestContext", "id")
		assert.Equal(t, filter.EntityTable, def.Location)
		assert.Equal(t, filter.IntValueType, def.ValueType)
		assert.Equal(t, "id", def.Column)
	})

	t.Run("property table property", func(t *testing.T) {
		def := reg.GetPropertyDefinitionForRestEntity("TestContext", "owner")
		assert.Equal(t, filter.PropertyTable, def.Location)
		assert.Equal(t, filter.StringValueType, def.ValueType)
		assert.Equal(t, "owner", def.Column)
	})

	t.Run("related entity prefix", func(t *testing.T) {
		def := reg.GetPropertyDefinitionForRestEntity("TestContext", "artifacts.modelFormatName")
		assert.Equal(t, filter.RelatedEntity, def.Location)
		assert.Equal(t, filter.RelatedEntityArtifact, def.RelatedEntityType)
		assert.Equal(t, "modelFormatName", def.RelatedProperty)
		assert.Equal(t, "Attribution", def.JoinTable)
		assert.Empty(t, def.ValueType) // runtime inference
	})

	t.Run("related entity with custom join table", func(t *testing.T) {
		reg2 := NewCatalogEntityRegistry()
		reg2.Register("WithJoin", EntityTypeDefinition{
			MLMDEntityType:         filter.EntityTypeContext,
			Properties:             CommonContextProperties(),
			RelatedEntityPrefix:    "children.",
			RelatedEntityType:      filter.RelatedEntityContext,
			RelatedEntityJoinTable: "ParentContext",
		})
		def := reg2.GetPropertyDefinitionForRestEntity("WithJoin", "children.status")
		assert.Equal(t, filter.RelatedEntity, def.Location)
		assert.Equal(t, "ParentContext", def.JoinTable)
		assert.Equal(t, "status", def.RelatedProperty)
	})

	t.Run("custom fallback", func(t *testing.T) {
		def := reg.GetPropertyDefinitionForRestEntity("TestContext", "unknownProp")
		assert.Equal(t, filter.Custom, def.Location)
		assert.Equal(t, filter.StringValueType, def.ValueType)
		assert.Equal(t, "unknownProp", def.Column)
	})

	t.Run("unknown entity type", func(t *testing.T) {
		def := reg.GetPropertyDefinitionForRestEntity("NoSuchEntity", "anything")
		assert.Equal(t, filter.Custom, def.Location)
	})

	t.Run("artifact entity table property", func(t *testing.T) {
		def := reg.GetPropertyDefinitionForRestEntity("TestArtifact", "uri")
		require.Equal(t, filter.EntityTable, def.Location)
		assert.Equal(t, filter.StringValueType, def.ValueType)
		assert.Equal(t, "uri", def.Column)
	})
}

func TestIsChildEntity(t *testing.T) {
	reg := buildTestRegistry()

	assert.False(t, reg.IsChildEntity("TestContext"))
	assert.False(t, reg.IsChildEntity("TestArtifact"))
	assert.True(t, reg.IsChildEntity("TestChild"))
	assert.False(t, reg.IsChildEntity("UnknownEntity"))
}

func TestRegistryImplementsInterface(t *testing.T) {
	// Compile-time check that CatalogEntityRegistry implements EntityMappingFunctions.
	var _ filter.EntityMappingFunctions = (*CatalogEntityRegistry)(nil)
}
