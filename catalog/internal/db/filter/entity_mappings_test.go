package filter

import (
	"testing"

	catalogmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// expectedCatalogModelProperties captures every property the old hand-coded
// catalogModelProperties map defined, so we can verify the registry is identical.
var expectedCatalogModelProperties = map[string]filter.PropertyDefinition{
	"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
	"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	"externalId":               {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},
	"source_id":                {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "source_id"},
	"description":              {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "description"},
	"owner":                    {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "owner"},
	"state":                    {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "state"},
	"language":                 {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "language"},
	"library_name":             {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "library_name"},
	"license_link":             {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license_link"},
	"license":                  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license"},
	"logo":                     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "logo"},
	"maturity":                 {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "maturity"},
	"provider":                 {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "provider"},
	"readme":                   {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "readme"},
	"tasks":                    {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "tasks"},
}

var expectedCatalogArtifactProperties = map[string]filter.PropertyDefinition{
	"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
	"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	"externalId":               {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},
	"uri":                      {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "uri"},
	"state":                    {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "state"},
	"artifactType":             {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "artifactType"},
}

var expectedMCPServerProperties = map[string]filter.PropertyDefinition{
	"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
	"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	"externalId":               {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},
	"source_id":                {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "source_id"},
	"base_name":                {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "base_name"},
	"description":              {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "description"},
	"provider":                 {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "provider"},
	"license":                  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license"},
	"license_link":             {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license_link"},
	"logo":                     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "logo"},
	"readme":                   {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "readme"},
	"version":                  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "version"},
	"tags":                     {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "tags"},
	"transports":               {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "transports"},
	"deploymentMode":           {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "deploymentMode"},
	"documentationUrl":         {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "documentationUrl"},
	"repositoryUrl":            {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "repositoryUrl"},
	"sourceCode":               {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "sourceCode"},
	"publishedDate":            {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "publishedDate"},
	"lastUpdated":              {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "lastUpdated"},
	"verifiedSource":           {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "verifiedSource"},
	"secureEndpoint":           {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "secureEndpoint"},
	"sast":                     {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "sast"},
	"readOnlyTools":            {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "readOnlyTools"},
	"endpoints":                {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "endpoints"},
	"artifacts":                {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "artifacts"},
	"runtimeMetadata":          {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "runtimeMetadata"},
}

var expectedMCPServerToolProperties = map[string]filter.PropertyDefinition{
	"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
	"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},
	"description":              {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "description"},
	"accessType":               {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "accessType"},
}

func TestBackwardCompat_CatalogModel(t *testing.T) {
	mappings := NewCatalogEntityMappings()
	entityType := filter.RestEntityType(catalogmodels.RestEntityCatalogModel)

	// Verify MLMD entity type
	assert.Equal(t, filter.EntityTypeContext, mappings.GetMLMDEntityType(entityType))

	// Verify every well-known property returns the expected definition
	for prop, expected := range expectedCatalogModelProperties {
		t.Run(prop, func(t *testing.T) {
			got := mappings.GetPropertyDefinitionForRestEntity(entityType, prop)
			assert.Equal(t, expected, got)
		})
	}

	// Verify unknown property returns custom fallback
	got := mappings.GetPropertyDefinitionForRestEntity(entityType, "nonExistentProp")
	assert.Equal(t, filter.Custom, got.Location)
	assert.Equal(t, filter.StringValueType, got.ValueType)
	assert.Equal(t, "nonExistentProp", got.Column)

	// Verify IsChildEntity
	assert.False(t, mappings.IsChildEntity(entityType))
}

func TestBackwardCompat_CatalogArtifact(t *testing.T) {
	mappings := NewCatalogEntityMappings()
	entityType := filter.RestEntityType(catalogmodels.RestEntityCatalogArtifact)

	assert.Equal(t, filter.EntityTypeArtifact, mappings.GetMLMDEntityType(entityType))

	for prop, expected := range expectedCatalogArtifactProperties {
		t.Run(prop, func(t *testing.T) {
			got := mappings.GetPropertyDefinitionForRestEntity(entityType, prop)
			assert.Equal(t, expected, got)
		})
	}

	assert.False(t, mappings.IsChildEntity(entityType))
}

func TestBackwardCompat_MCPServer(t *testing.T) {
	mappings := NewCatalogEntityMappings()
	entityType := filter.RestEntityType(catalogmodels.RestEntityMCPServer)

	assert.Equal(t, filter.EntityTypeContext, mappings.GetMLMDEntityType(entityType))

	for prop, expected := range expectedMCPServerProperties {
		t.Run(prop, func(t *testing.T) {
			got := mappings.GetPropertyDefinitionForRestEntity(entityType, prop)
			assert.Equal(t, expected, got)
		})
	}

	assert.False(t, mappings.IsChildEntity(entityType))
}

func TestBackwardCompat_MCPServerTool(t *testing.T) {
	mappings := NewCatalogEntityMappings()
	entityType := filter.RestEntityType(catalogmodels.RestEntityMCPServerTool)

	assert.Equal(t, filter.EntityTypeArtifact, mappings.GetMLMDEntityType(entityType))

	for prop, expected := range expectedMCPServerToolProperties {
		t.Run(prop, func(t *testing.T) {
			got := mappings.GetPropertyDefinitionForRestEntity(entityType, prop)
			assert.Equal(t, expected, got)
		})
	}

	assert.False(t, mappings.IsChildEntity(entityType))
}

func TestBackwardCompat_ArtifactPrefix(t *testing.T) {
	mappings := NewCatalogEntityMappings()
	entityType := filter.RestEntityType(catalogmodels.RestEntityCatalogModel)

	tests := []struct {
		name            string
		propertyPath    string
		expectedColumn  string
		expectedRelProp string
	}{
		{
			name:            "simple artifact property",
			propertyPath:    "artifacts.modelFormatName",
			expectedColumn:  "modelFormatName",
			expectedRelProp: "modelFormatName",
		},
		{
			name:            "nested artifact custom property",
			propertyPath:    "artifacts.customProperties.myProp",
			expectedColumn:  "customProperties.myProp",
			expectedRelProp: "customProperties.myProp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mappings.GetPropertyDefinitionForRestEntity(entityType, tt.propertyPath)
			require.Equal(t, filter.RelatedEntity, got.Location)
			assert.Empty(t, got.ValueType, "ValueType should be empty for runtime inference")
			assert.Equal(t, tt.expectedColumn, got.Column)
			assert.Equal(t, filter.RelatedEntityArtifact, got.RelatedEntityType)
			assert.Equal(t, tt.expectedRelProp, got.RelatedProperty)
			assert.Equal(t, "Attribution", got.JoinTable)
		})
	}
}

func TestBackwardCompat_UnknownEntityDefaultsToContext(t *testing.T) {
	mappings := NewCatalogEntityMappings()
	assert.Equal(t, filter.EntityTypeContext, mappings.GetMLMDEntityType("SomeUnknownEntity"))
}
