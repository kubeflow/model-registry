package service

import (
	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/platform/db/repository"
)

// Generic constraints
type SchemaEntity = repository.SchemaEntity
type PropertyEntity = repository.PropertyEntity

// Mapper function types
type EntityToSchemaMapper[TEntity any, TSchema SchemaEntity] = repository.EntityToSchemaMapper[TEntity, TSchema]
type SchemaToEntityMapper[TSchema SchemaEntity, TProp PropertyEntity, TEntity any] = repository.SchemaToEntityMapper[TSchema, TProp, TEntity]
type EntityToPropertiesMapper[TEntity any, TProp PropertyEntity] = repository.EntityToPropertiesMapper[TEntity, TProp]

// Interfaces
type BaseListOptions = repository.BaseListOptions
type FilterApplier = repository.FilterApplier

// Generic repository types
type GenericRepositoryConfig[TEntity any, TSchema SchemaEntity, TProp PropertyEntity, TListOpts BaseListOptions] = repository.GenericRepositoryConfig[TEntity, TSchema, TProp, TListOpts]
type GenericRepository[TEntity any, TSchema SchemaEntity, TProp PropertyEntity, TListOpts BaseListOptions] = repository.GenericRepository[TEntity, TSchema, TProp, TListOpts]

// Functions
var ApplyFilterQuery = repository.ApplyFilterQuery

// NewGenericRepository creates a GenericRepository, injecting the default MR entity
// mapping functions when EntityMappingFuncs is not set in the config.
func NewGenericRepository[TEntity any, TSchema SchemaEntity, TProp PropertyEntity, TListOpts BaseListOptions](
	config GenericRepositoryConfig[TEntity, TSchema, TProp, TListOpts],
) *GenericRepository[TEntity, TSchema, TProp, TListOpts] {
	if config.EntityMappingFuncs == nil {
		config.EntityMappingFuncs = filter.DefaultEntityMappingFuncs()
	}
	return repository.NewGenericRepository(config)
}

// applyFilterQuery is a legacy alias for backward compatibility within this package
var applyFilterQuery = ApplyFilterQuery

// Property mapping functions
var (
	MapPropertiesToArtifactProperty   = repository.MapPropertiesToArtifactProperty
	MapPropertiesToContextProperty    = repository.MapPropertiesToContextProperty
	MapPropertiesToExecutionProperty  = repository.MapPropertiesToExecutionProperty
	MapArtifactPropertyToProperties   = repository.MapArtifactPropertyToProperties
	MapContextPropertyToProperties    = repository.MapContextPropertyToProperties
	MapExecutionPropertyToProperties  = repository.MapExecutionPropertyToProperties
)
