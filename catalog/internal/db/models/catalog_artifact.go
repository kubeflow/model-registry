package models

import (
	"sync"

	dbfilter "github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
)

type CatalogArtifactListOptions struct {
	models.Pagination
	Name                *string
	ExternalID          *string
	ParentResourceID    *int32
	ArtifactType        *string
	ArtifactTypesFilter []string
}

// GetRestEntityType implements the FilterApplier interface
// This enables advanced filtering support for catalog artifacts
func (c *CatalogArtifactListOptions) GetRestEntityType() dbfilter.RestEntityType {
	return dbfilter.RestEntityType(RestEntityCatalogArtifact)
}

// ArtifactMapperFunc defines the signature for artifact mapping functions
type ArtifactMapperFunc func(artifact schema.Artifact, properties []schema.ArtifactProperty) interface{}

// Global registry for artifact mappers
var (
	artifactMappersMu sync.RWMutex
	artifactMappers   = make(map[string]ArtifactMapperFunc)
)

// RegisterArtifactMapper registers a mapping function for a specific artifact type
func RegisterArtifactMapper(typeName string, mapper ArtifactMapperFunc) {
	artifactMappersMu.Lock()
	defer artifactMappersMu.Unlock()
	artifactMappers[typeName] = mapper
}

// GetArtifactMapper retrieves a mapping function for a specific artifact type
func GetArtifactMapper(typeName string) (ArtifactMapperFunc, bool) {
	artifactMappersMu.RLock()
	defer artifactMappersMu.RUnlock()
	mapper, exists := artifactMappers[typeName]
	return mapper, exists
}

// CatalogArtifactEntity defines the common interface that all catalog artifacts must implement.
// This allows the shared infrastructure to work with catalog-specific types without import cycles.
// Note: GetAttributes() is intentionally excluded because concrete types return different
// typed pointers, which Go's interface matching does not allow.
type CatalogArtifactEntity interface {
	GetID() *int32
	SetID(int32)
	GetProperties() *[]models.Properties
	GetCustomProperties() *[]models.Properties
}

// CatalogModelArtifact represents the interface for model artifacts
type CatalogModelArtifact CatalogArtifactEntity

// CatalogMetricsArtifact represents the interface for metrics artifacts
type CatalogMetricsArtifact CatalogArtifactEntity

// CatalogArtifact is a discriminated union that can hold different catalog artifact types
type CatalogArtifact struct {
	CatalogModelArtifact   CatalogModelArtifact
	CatalogMetricsArtifact CatalogMetricsArtifact
}

type CatalogArtifactRepository interface {
	GetByID(id int32) (CatalogArtifact, error)
	List(listOptions CatalogArtifactListOptions) (*models.ListWrapper[CatalogArtifact], error)
	DeleteByParentID(artifactType string, parentResourceID int32) error
}
