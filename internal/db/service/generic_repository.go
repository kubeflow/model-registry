package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/kubeflow/model-registry/internal/db/dbutil"
	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"github.com/kubeflow/model-registry/pkg/api"
	"gorm.io/gorm"
)

// Generic constraints for different entity types
type SchemaEntity interface {
	schema.Artifact | schema.Context | schema.Execution
}

type PropertyEntity interface {
	schema.ArtifactProperty | schema.ContextProperty | schema.ExecutionProperty
}

// Mapper function types
type EntityToSchemaMapper[TEntity any, TSchema SchemaEntity] func(TEntity) TSchema
type SchemaToEntityMapper[TSchema SchemaEntity, TProp PropertyEntity, TEntity any] func(TSchema, []TProp) TEntity
type EntityToPropertiesMapper[TEntity any, TProp PropertyEntity] func(TEntity, int32) []TProp

// List options interface
type BaseListOptions interface {
	GetPageSize() int32
	GetNextPageToken() string
	SetNextPageToken(*string)
	GetOrderBy() string
	GetSortOrder() string
	GetFilterQuery() string
}

// Filter applier interface for entities that support advanced filtering
type FilterApplier interface {
	GetRestEntityType() filter.RestEntityType
}

// ApplyFilterQuery applies advanced filter query processing to a GORM query
// This function encapsulates the common pattern used by both GenericRepository and custom repositories
func ApplyFilterQuery(query *gorm.DB, listOptions any, mappingFuncs filter.EntityMappingFunctions) (*gorm.DB, error) {
	if filterQueryGetter, ok := listOptions.(interface{ GetFilterQuery() string }); ok {
		if filterQuery := filterQueryGetter.GetFilterQuery(); filterQuery != "" {
			if filterApplier, ok := listOptions.(FilterApplier); ok {
				filterExpr, err := filter.Parse(filterQuery)
				if err != nil {
					// Enhance error message with helpful hints for common mistakes
					enhancedErr := dbutil.EnhanceFilterQueryError(err, filterQuery)
					return nil, fmt.Errorf("%v: %w", enhancedErr, api.ErrBadRequest)
				}

				if filterExpr != nil {
					queryBuilder := filter.NewQueryBuilderForRestEntity(filterApplier.GetRestEntityType(), mappingFuncs)
					query = queryBuilder.BuildQuery(query, filterExpr)
				}
			}
		}
	}
	return query, nil
}

// applyFilterQuery is a legacy alias for backward compatibility
func applyFilterQuery(query *gorm.DB, listOptions any, mappingFuncs filter.EntityMappingFunctions) (*gorm.DB, error) {
	return ApplyFilterQuery(query, listOptions, mappingFuncs)
}

// Generic repository configuration
type GenericRepositoryConfig[TEntity any, TSchema SchemaEntity, TProp PropertyEntity, TListOpts BaseListOptions] struct {
	DB                      *gorm.DB
	TypeID                  int32
	EntityToSchema          EntityToSchemaMapper[TEntity, TSchema]
	SchemaToEntity          SchemaToEntityMapper[TSchema, TProp, TEntity]
	EntityToProperties      EntityToPropertiesMapper[TEntity, TProp]
	NotFoundError           error
	EntityName              string
	PropertyFieldName       string // "artifact_id", "context_id", or "execution_id"
	ApplyListFilters        func(*gorm.DB, TListOpts) *gorm.DB
	CreatePaginationToken   func(TSchema, TListOpts) string    // Optional - defaults to standard implementation
	ApplyCustomOrdering     func(*gorm.DB, TListOpts) *gorm.DB // Optional - custom ordering logic that bypasses standard pagination
	IsNewEntity             func(TEntity) bool
	HasCustomProperties     func(TEntity) bool
	EntityMappingFuncs      filter.EntityMappingFunctions // Optional - custom entity mappings for filtering
	PreserveHistoricalTimes bool                          // Optional - when true, preserves timestamps from source data (e.g. YAML catalog loading). Default false (Model Registry behavior - always auto-generate timestamps)
}

// Generic repository implementation
type GenericRepository[TEntity any, TSchema SchemaEntity, TProp PropertyEntity, TListOpts BaseListOptions] struct {
	config GenericRepositoryConfig[TEntity, TSchema, TProp, TListOpts]
}

func NewGenericRepository[TEntity any, TSchema SchemaEntity, TProp PropertyEntity, TListOpts BaseListOptions](
	config GenericRepositoryConfig[TEntity, TSchema, TProp, TListOpts],
) *GenericRepository[TEntity, TSchema, TProp, TListOpts] {
	return &GenericRepository[TEntity, TSchema, TProp, TListOpts]{
		config: config,
	}
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) GetByID(id int32) (TEntity, error) {
	var entity TSchema
	var properties []TProp
	var zeroEntity TEntity

	// Query main entity
	if err := r.config.DB.Where("id = ? AND type_id = ?", id, r.config.TypeID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return zeroEntity, fmt.Errorf("%w: %v", r.config.NotFoundError, err)
		}
		return zeroEntity, fmt.Errorf("error getting %s by id: %w", r.config.EntityName, err)
	}

	// Query properties
	entityID := r.getEntityID(entity)
	if err := r.config.DB.Where(r.config.PropertyFieldName+" = ?", entityID).Find(&properties).Error; err != nil {
		return zeroEntity, fmt.Errorf("error getting properties by %s id: %w", r.config.EntityName, err)
	}

	// Map to domain model
	return r.config.SchemaToEntity(entity, properties), nil
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) GetByName(name string) (TEntity, error) {
	var entity TSchema
	var properties []TProp
	var zeroEntity TEntity

	// Query main entity
	if err := r.config.DB.Where("name = ? AND type_id = ?", name, r.config.TypeID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return zeroEntity, fmt.Errorf("%w: %v", r.config.NotFoundError, err)
		}
		return zeroEntity, fmt.Errorf("error getting %s by name: %w", r.config.EntityName, err)
	}

	// Query properties
	entityID := r.getEntityID(entity)
	if err := r.config.DB.Where(r.config.PropertyFieldName+" = ?", entityID).Find(&properties).Error; err != nil {
		return zeroEntity, fmt.Errorf("error getting properties by %s id: %w", r.config.EntityName, err)
	}

	// Map to domain model
	return r.config.SchemaToEntity(entity, properties), nil
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) List(listOptions TListOpts) (*models.ListWrapper[TEntity], error) {
	pageSize := listOptions.GetPageSize()

	list := models.ListWrapper[TEntity]{
		PageSize: pageSize,
	}

	var entities []TEntity
	var schemaEntities []TSchema

	// Build base query
	query := r.buildBaseQuery()

	// Apply type-specific filters
	if r.config.ApplyListFilters != nil {
		query = r.config.ApplyListFilters(query, listOptions)
	}

	// Apply advanced filter query if supported
	query, err := applyFilterQuery(query, listOptions, r.config.EntityMappingFuncs)
	if err != nil {
		return nil, err
	}

	// Apply ordering and pagination
	if r.config.ApplyCustomOrdering != nil {
		// Use custom ordering logic if provided
		query = r.config.ApplyCustomOrdering(query, listOptions)
	} else {
		// Apply standard pagination
		query = r.ApplyStandardPagination(query, listOptions, entities)
	}

	// Execute query
	if err := query.Find(&schemaEntities).Error; err != nil {
		// Sanitize database errors to avoid exposing internal details to users
		err = dbutil.SanitizeDatabaseError(err)
		return nil, fmt.Errorf("error listing %ss: %w", r.config.EntityName, err)
	}

	// Handle pagination
	hasMore := false
	if pageSize > 0 {
		hasMore = len(schemaEntities) > int(pageSize)
		if hasMore {
			schemaEntities = schemaEntities[:len(schemaEntities)-1]
		}
	}

	// Load properties and map to domain models
	for _, schemaEntity := range schemaEntities {
		var properties []TProp
		entityID := r.getEntityID(schemaEntity)
		if err := r.config.DB.Where(r.config.PropertyFieldName+" = ?", entityID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by %s id: %w", r.config.EntityName, err)
		}

		entity := r.config.SchemaToEntity(schemaEntity, properties)
		entities = append(entities, entity)
	}

	// Set pagination token
	if hasMore && len(schemaEntities) > 0 {
		lastEntity := schemaEntities[len(schemaEntities)-1]
		var nextToken string
		if r.config.CreatePaginationToken != nil {
			nextToken = r.config.CreatePaginationToken(lastEntity, listOptions)
		} else {
			nextToken = r.CreateDefaultPaginationToken(lastEntity, listOptions)
		}
		listOptions.SetNextPageToken(&nextToken)
	} else {
		listOptions.SetNextPageToken(nil)
	}

	list.Items = entities
	nextPageToken := listOptions.GetNextPageToken()
	list.NextPageToken = nextPageToken
	list.Size = int32(len(entities))

	return &list, nil
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) Save(entity TEntity, parentResourceID *int32) (TEntity, error) {
	now := time.Now().UnixMilli()
	var zeroEntity TEntity

	schemaEntity := r.config.EntityToSchema(entity)
	var finalProperties []TProp

	// Determine if this is a new entity or an update
	isNewEntity := r.config.IsNewEntity != nil && r.config.IsNewEntity(entity)

	// Set timestamps based on configuration and entity state
	existingCreateTime := r.getCreateTime(schemaEntity)
	existingUpdateTime := r.getLastUpdateTime(schemaEntity)

	if r.config.PreserveHistoricalTimes {
		// Catalog mode: Preserve historical timestamps from source data (e.g. YAML)
		// - For new entities: only set if not already present (preserves YAML timestamps)
		// - For updates: only set if not already present (preserves YAML timestamps)
		if isNewEntity {
			if existingUpdateTime == 0 {
				r.setLastUpdateTime(&schemaEntity, now)
			}
			if existingCreateTime == 0 {
				r.setCreateTime(&schemaEntity, now)
			}
		} else {
			// Update: preserve timestamps from YAML if present
			if existingUpdateTime == 0 {
				r.setLastUpdateTime(&schemaEntity, now)
			}
			if existingCreateTime == 0 {
				r.setCreateTime(&schemaEntity, now)
			}
		}
	} else {
		// Model Registry mode (default): Always auto-generate timestamps
		// - For new entities: always set both timestamps to current time
		// - For updates: always update LastUpdateTime, preserve CreateTime if present
		if isNewEntity {
			r.setLastUpdateTime(&schemaEntity, now)
			r.setCreateTime(&schemaEntity, now)
		} else {
			// Always update LastUpdateTime for existing entities being updated
			r.setLastUpdateTime(&schemaEntity, now)
			// Only set CreateTime if it's not already set (preserve historical CreateTime)
			if existingCreateTime == 0 {
				r.setCreateTime(&schemaEntity, now)
			}
		}
	}

	hasCustomProperties := r.config.HasCustomProperties != nil && r.config.HasCustomProperties(entity)

	err := r.config.DB.Transaction(func(tx *gorm.DB) error {
		// Save main entity with smart field handling
		if isNewEntity {
			// For new entities, save all fields
			if err := tx.Save(&schemaEntity).Error; err != nil {
				return fmt.Errorf("error saving %s: %w", r.config.EntityName, err)
			}
		} else {
			// For updates, use Updates() to only update changed fields
			// Updates() automatically handles zero values correctly and respects omitted fields
			omitFields := r.getNonUpdatableFields(schemaEntity)
			if err := tx.Model(&schemaEntity).Omit(omitFields...).Updates(&schemaEntity).Error; err != nil {
				return fmt.Errorf("error saving %s: %w", r.config.EntityName, err)
			}
		}

		// Handle parent relationship if applicable
		if parentResourceID != nil {
			if err := r.handleParentRelationship(tx, schemaEntity, parentResourceID); err != nil {
				return err
			}
		}

		// Handle properties
		entityID := r.getEntityID(schemaEntity)
		properties := r.config.EntityToProperties(entity, entityID)

		if err := r.handleProperties(tx, entityID, properties, hasCustomProperties); err != nil {
			return err
		}

		// Get final properties for return object
		if err := tx.Where(r.config.PropertyFieldName+" = ?", entityID).Find(&finalProperties).Error; err != nil {
			return fmt.Errorf("error getting final properties by %s id: %w", r.config.EntityName, err)
		}

		return nil
	})
	if err != nil {
		return zeroEntity, err
	}

	// Return the updated entity
	return r.config.SchemaToEntity(schemaEntity, finalProperties), nil
}

// Helper methods

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) buildBaseQuery() *gorm.DB {
	var schemaEntity TSchema
	var tableName string
	var model interface{}

	// Determine table name and model based on schema entity type
	switch any(schemaEntity).(type) {
	case schema.Artifact:
		tableName = "Artifact"
		model = &schema.Artifact{}
	case schema.Context:
		tableName = "Context"
		model = &schema.Context{}
	case schema.Execution:
		tableName = "Execution"
		model = &schema.Execution{}
	default:
		panic(fmt.Sprintf("unsupported schema entity type: %T", schemaEntity))
	}

	// Quote table name based on database dialect and build WHERE clause
	tableNameQuoted := dbutil.QuoteTableName(r.config.DB, tableName)
	whereClause := fmt.Sprintf("%s.type_id = ?", tableNameQuoted)

	return r.config.DB.Model(model).Where(whereClause, r.config.TypeID)
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) getEntityID(entity TSchema) int32 {
	switch e := any(entity).(type) {
	case schema.Artifact:
		return e.ID
	case schema.Context:
		return e.ID
	case schema.Execution:
		return e.ID
	default:
		panic(fmt.Sprintf("unsupported entity type: %T", entity))
	}
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) setLastUpdateTime(entity *TSchema, timestamp int64) {
	switch e := any(entity).(type) {
	case *schema.Artifact:
		e.LastUpdateTimeSinceEpoch = timestamp
	case *schema.Context:
		e.LastUpdateTimeSinceEpoch = timestamp
	case *schema.Execution:
		e.LastUpdateTimeSinceEpoch = timestamp
	}
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) setCreateTime(entity *TSchema, timestamp int64) {
	switch e := any(entity).(type) {
	case *schema.Artifact:
		e.CreateTimeSinceEpoch = timestamp
	case *schema.Context:
		e.CreateTimeSinceEpoch = timestamp
	case *schema.Execution:
		e.CreateTimeSinceEpoch = timestamp
	}
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) handleParentRelationship(tx *gorm.DB, entity TSchema, parentResourceID *int32) error {
	// Handle Attribution for artifacts, ParentContext for contexts, or Association for executions
	entityID := r.getEntityID(entity)

	switch any(entity).(type) {
	case schema.Artifact:
		// Check if attribution already exists to avoid duplicate key errors
		var existingAttribution schema.Attribution
		result := tx.Where("context_id = ? AND artifact_id = ?", *parentResourceID, entityID).First(&existingAttribution)

		if result.Error == gorm.ErrRecordNotFound {
			// Attribution doesn't exist, create it
			attribution := schema.Attribution{
				ArtifactID: entityID,
				ContextID:  *parentResourceID,
			}
			if err := tx.Create(&attribution).Error; err != nil {
				return fmt.Errorf("error creating attribution: %w", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("error checking existing attribution: %w", result.Error)
		}
		// If attribution already exists, do nothing

	case schema.Context:
		// Check if parent context already exists to avoid duplicate key errors
		var existingParentContext schema.ParentContext
		result := tx.Where("parent_context_id = ? AND context_id = ?", *parentResourceID, entityID).First(&existingParentContext)

		if result.Error == gorm.ErrRecordNotFound {
			// Parent context doesn't exist, create it
			parentContext := schema.ParentContext{
				ContextID:       entityID,
				ParentContextID: *parentResourceID,
			}
			if err := tx.Create(&parentContext).Error; err != nil {
				return fmt.Errorf("error creating parent context: %w", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("error checking existing parent context: %w", result.Error)
		}
		// If parent context already exists, do nothing

	case schema.Execution:
		// Check if association already exists to avoid duplicate key errors
		var existingAssociation schema.Association
		result := tx.Where("context_id = ? AND execution_id = ?", *parentResourceID, entityID).First(&existingAssociation)

		if result.Error == gorm.ErrRecordNotFound {
			// Association doesn't exist, create it
			association := schema.Association{
				ExecutionID: entityID,
				ContextID:   *parentResourceID,
			}
			if err := tx.Create(&association).Error; err != nil {
				return fmt.Errorf("error creating association: %w", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("error checking existing association: %w", result.Error)
		}
		// If association already exists, do nothing
	}

	return nil
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) handleProperties(tx *gorm.DB, entityID int32, properties []TProp, hasCustomProperties bool) error {
	// Get existing custom properties if we have custom properties
	if hasCustomProperties {
		var existingCustomProperties []TProp
		if err := tx.Where(r.config.PropertyFieldName+" = ? AND is_custom_property = ?", entityID, true).Find(&existingCustomProperties).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties: %w", err)
		}

		// Delete removed custom properties
		for _, existingProp := range existingCustomProperties {
			found := false
			for _, prop := range properties {
				if r.propertiesMatch(prop, existingProp) {
					found = true
					break
				}
			}

			if !found {
				if err := tx.Delete(&existingProp).Error; err != nil {
					return fmt.Errorf("error deleting property: %w", err)
				}
			}
		}
	}

	// Upsert properties
	for _, prop := range properties {
		var existingProp TProp
		result := tx.Where(r.config.PropertyFieldName+" = ? AND name = ? AND is_custom_property = ?",
			entityID, r.getPropertyName(prop), r.getPropertyIsCustom(prop)).First(&existingProp)

		switch result.Error {
		case nil:
			// Update existing property
			r.copyPropertyValues(&prop, &existingProp)
			if err := tx.Model(&existingProp).Updates(prop).Error; err != nil {
				return fmt.Errorf("error updating property %s: %w", r.getPropertyName(prop), err)
			}
		case gorm.ErrRecordNotFound:
			// Create new property
			if err := tx.Create(&prop).Error; err != nil {
				return fmt.Errorf("error creating property %s: %w", r.getPropertyName(prop), err)
			}
		default:
			return fmt.Errorf("error checking existing property %s: %w", r.getPropertyName(prop), result.Error)
		}
	}

	return nil
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) propertiesMatch(prop1, prop2 TProp) bool {
	return r.getPropertyName(prop1) == r.getPropertyName(prop2) &&
		r.getPropertyIsCustom(prop1) == r.getPropertyIsCustom(prop2) &&
		r.getPropertyEntityID(prop1) == r.getPropertyEntityID(prop2)
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) getPropertyName(prop TProp) string {
	switch p := any(prop).(type) {
	case schema.ArtifactProperty:
		return p.Name
	case schema.ContextProperty:
		return p.Name
	case schema.ExecutionProperty:
		return p.Name
	default:
		panic(fmt.Sprintf("unsupported property type: %T", prop))
	}
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) getPropertyIsCustom(prop TProp) bool {
	switch p := any(prop).(type) {
	case schema.ArtifactProperty:
		return p.IsCustomProperty
	case schema.ContextProperty:
		return p.IsCustomProperty
	case schema.ExecutionProperty:
		return p.IsCustomProperty
	default:
		panic(fmt.Sprintf("unsupported property type: %T", prop))
	}
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) getPropertyEntityID(prop TProp) int32 {
	switch p := any(prop).(type) {
	case schema.ArtifactProperty:
		return p.ArtifactID
	case schema.ContextProperty:
		return p.ContextID
	case schema.ExecutionProperty:
		return p.ExecutionID
	default:
		panic(fmt.Sprintf("unsupported property type: %T", prop))
	}
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) copyPropertyValues(src, dst *TProp) {
	switch srcProp := any(src).(type) {
	case *schema.ArtifactProperty:
		dstProp := any(dst).(*schema.ArtifactProperty)
		dstProp.IntValue = srcProp.IntValue
		dstProp.DoubleValue = srcProp.DoubleValue
		dstProp.StringValue = srcProp.StringValue
		dstProp.BoolValue = srcProp.BoolValue
		dstProp.ByteValue = srcProp.ByteValue
		dstProp.ProtoValue = srcProp.ProtoValue
	case *schema.ContextProperty:
		dstProp := any(dst).(*schema.ContextProperty)
		dstProp.IntValue = srcProp.IntValue
		dstProp.DoubleValue = srcProp.DoubleValue
		dstProp.StringValue = srcProp.StringValue
		dstProp.BoolValue = srcProp.BoolValue
		dstProp.ByteValue = srcProp.ByteValue
		dstProp.ProtoValue = srcProp.ProtoValue
	case *schema.ExecutionProperty:
		dstProp := any(dst).(*schema.ExecutionProperty)
		dstProp.IntValue = srcProp.IntValue
		dstProp.DoubleValue = srcProp.DoubleValue
		dstProp.StringValue = srcProp.StringValue
		dstProp.BoolValue = srcProp.BoolValue
		dstProp.ByteValue = srcProp.ByteValue
		dstProp.ProtoValue = srcProp.ProtoValue
	}
}

// CreateDefaultPaginationToken provides a standard implementation that works for all entities
// with ID, CreateTimeSinceEpoch, and LastUpdateTimeSinceEpoch fields
func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) CreateDefaultPaginationToken(entity TSchema, listOptions TListOpts) string {
	entityID := r.getEntityID(entity)
	orderBy := listOptions.GetOrderBy()
	value := ""

	if orderBy != "" {
		switch orderBy {
		case "ID":
			value = fmt.Sprintf("%d", entityID)
		case "CREATE_TIME":
			value = fmt.Sprintf("%d", r.getCreateTime(entity))
		case "LAST_UPDATE_TIME":
			value = fmt.Sprintf("%d", r.getLastUpdateTime(entity))
		default:
			value = fmt.Sprintf("%d", entityID)
		}
	}

	return scopes.CreateNextPageToken(entityID, value)
}

// getCreateTime extracts CreateTimeSinceEpoch from any schema entity
func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) getCreateTime(entity TSchema) int64 {
	switch e := any(entity).(type) {
	case schema.Artifact:
		return e.CreateTimeSinceEpoch
	case schema.Context:
		return e.CreateTimeSinceEpoch
	case schema.Execution:
		return e.CreateTimeSinceEpoch
	default:
		return 0
	}
}

// getLastUpdateTime extracts LastUpdateTimeSinceEpoch from any schema entity
func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) getLastUpdateTime(entity TSchema) int64 {
	switch e := any(entity).(type) {
	case schema.Artifact:
		return e.LastUpdateTimeSinceEpoch
	case schema.Context:
		return e.LastUpdateTimeSinceEpoch
	case schema.Execution:
		return e.LastUpdateTimeSinceEpoch
	default:
		return 0
	}
}

// getNonUpdatableFields returns the list of fields that should be omitted during updates
func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) getNonUpdatableFields(entity TSchema) []string {
	var omitFields []string

	switch any(entity).(type) {
	case schema.Artifact:
		// Non-updatable fields for artifacts: id, name, type_id, create_time_since_epoch
		omitFields = []string{"id", "name", "type_id", "create_time_since_epoch"}
	case schema.Context:
		// Non-updatable fields for contexts: id, name, type_id, create_time_since_epoch
		omitFields = []string{"id", "name", "type_id", "create_time_since_epoch"}
	case schema.Execution:
		// Non-updatable fields for executions: id, name, type_id, create_time_since_epoch
		omitFields = []string{"id", "name", "type_id", "create_time_since_epoch"}
	default:
		// Default case: omit common non-updatable fields
		omitFields = []string{"id", "name", "type_id", "create_time_since_epoch"}
	}

	return omitFields
}

func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) GetConfig() GenericRepositoryConfig[TEntity, TSchema, TProp, TListOpts] {
	return r.config
}

// ApplyStandardPagination applies the standard pagination logic using scopes.PaginateWithTablePrefix
func (r *GenericRepository[TEntity, TSchema, TProp, TListOpts]) ApplyStandardPagination(query *gorm.DB, listOptions TListOpts, entities any) *gorm.DB {
	pageSize := listOptions.GetPageSize()
	orderBy := listOptions.GetOrderBy()
	sortOrder := listOptions.GetSortOrder()
	nextPageToken := listOptions.GetNextPageToken()

	pagination := &models.Pagination{
		PageSize:      &pageSize,
		OrderBy:       &orderBy,
		SortOrder:     &sortOrder,
		NextPageToken: &nextPageToken,
	}

	// Use table prefix for pagination to handle JOINs properly
	var tablePrefix string
	var schemaEntity TSchema
	switch any(schemaEntity).(type) {
	case schema.Artifact:
		tablePrefix = "Artifact"
	case schema.Context:
		tablePrefix = "Context"
	case schema.Execution:
		tablePrefix = "Execution"
	default:
		tablePrefix = ""
	}

	return query.Scopes(scopes.PaginateWithTablePrefix(entities, pagination, r.config.DB, tablePrefix))
}

// Shared mapping functions for common property conversions

// MapPropertiesToArtifactProperty converts models.Properties to schema.ArtifactProperty
func MapPropertiesToArtifactProperty(prop models.Properties, artifactID int32, isCustomProperty bool) schema.ArtifactProperty {
	return schema.ArtifactProperty{
		ArtifactID:       artifactID,
		Name:             prop.Name,
		IsCustomProperty: isCustomProperty,
		IntValue:         prop.IntValue,
		DoubleValue:      prop.DoubleValue,
		StringValue:      prop.StringValue,
		BoolValue:        prop.BoolValue,
		ByteValue:        prop.ByteValue,
		ProtoValue:       prop.ProtoValue,
	}
}

// MapPropertiesToContextProperty converts models.Properties to schema.ContextProperty
func MapPropertiesToContextProperty(prop models.Properties, contextID int32, isCustomProperty bool) schema.ContextProperty {
	return schema.ContextProperty{
		ContextID:        contextID,
		Name:             prop.Name,
		IsCustomProperty: isCustomProperty,
		IntValue:         prop.IntValue,
		DoubleValue:      prop.DoubleValue,
		StringValue:      prop.StringValue,
		BoolValue:        prop.BoolValue,
		ByteValue:        prop.ByteValue,
		ProtoValue:       prop.ProtoValue,
	}
}

// MapPropertiesToExecutionProperty converts models.Properties to schema.ExecutionProperty
func MapPropertiesToExecutionProperty(prop models.Properties, executionID int32, isCustomProperty bool) schema.ExecutionProperty {
	return schema.ExecutionProperty{
		ExecutionID:      executionID,
		Name:             prop.Name,
		IsCustomProperty: isCustomProperty,
		IntValue:         prop.IntValue,
		DoubleValue:      prop.DoubleValue,
		StringValue:      prop.StringValue,
		BoolValue:        prop.BoolValue,
		ByteValue:        prop.ByteValue,
		ProtoValue:       prop.ProtoValue,
	}
}

// MapArtifactPropertyToProperties converts schema.ArtifactProperty to models.Properties
func MapArtifactPropertyToProperties(artProperty schema.ArtifactProperty) models.Properties {
	return models.Properties{
		Name:             artProperty.Name,
		IsCustomProperty: artProperty.IsCustomProperty,
		IntValue:         artProperty.IntValue,
		DoubleValue:      artProperty.DoubleValue,
		StringValue:      artProperty.StringValue,
		BoolValue:        artProperty.BoolValue,
		ByteValue:        artProperty.ByteValue,
		ProtoValue:       artProperty.ProtoValue,
	}
}

// MapContextPropertyToProperties converts schema.ContextProperty to models.Properties
func MapContextPropertyToProperties(contextProperty schema.ContextProperty) models.Properties {
	return models.Properties{
		Name:             contextProperty.Name,
		IsCustomProperty: contextProperty.IsCustomProperty,
		IntValue:         contextProperty.IntValue,
		DoubleValue:      contextProperty.DoubleValue,
		StringValue:      contextProperty.StringValue,
		BoolValue:        contextProperty.BoolValue,
		ByteValue:        contextProperty.ByteValue,
		ProtoValue:       contextProperty.ProtoValue,
	}
}

// MapExecutionPropertyToProperties converts schema.ExecutionProperty to models.Properties
func MapExecutionPropertyToProperties(executionProperty schema.ExecutionProperty) models.Properties {
	return models.Properties{
		Name:             executionProperty.Name,
		IsCustomProperty: executionProperty.IsCustomProperty,
		IntValue:         executionProperty.IntValue,
		DoubleValue:      executionProperty.DoubleValue,
		StringValue:      executionProperty.StringValue,
		BoolValue:        executionProperty.BoolValue,
		ByteValue:        executionProperty.ByteValue,
		ProtoValue:       executionProperty.ProtoValue,
	}
}
