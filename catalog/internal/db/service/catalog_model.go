package service

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/kubeflow/model-registry/catalog/internal/db/filter"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"gorm.io/gorm"
)

var ErrCatalogModelNotFound = errors.New("catalog model by id not found")

type CatalogModelRepositoryImpl struct {
	*service.GenericRepository[models.CatalogModel, schema.Context, schema.ContextProperty, *models.CatalogModelListOptions]
}

func NewCatalogModelRepository(db *gorm.DB, typeID int64) models.CatalogModelRepository {
	r := &CatalogModelRepositoryImpl{}

	r.GenericRepository = service.NewGenericRepository(service.GenericRepositoryConfig[models.CatalogModel, schema.Context, schema.ContextProperty, *models.CatalogModelListOptions]{
		DB:                  db,
		TypeID:              typeID,
		EntityToSchema:      mapCatalogModelToContext,
		SchemaToEntity:      mapDataLayerToCatalogModel,
		EntityToProperties:  mapCatalogModelToContextProperties,
		NotFoundError:       ErrCatalogModelNotFound,
		EntityName:          "catalog model",
		PropertyFieldName:   "context_id",
		ApplyListFilters:    applyCatalogModelListFilters,
		IsNewEntity:         func(entity models.CatalogModel) bool { return entity.GetID() == nil },
		HasCustomProperties: func(entity models.CatalogModel) bool { return entity.GetCustomProperties() != nil },
		EntityMappingFuncs:  filter.NewCatalogEntityMappings(),
	})

	return r
}

func (r *CatalogModelRepositoryImpl) Save(model models.CatalogModel) (models.CatalogModel, error) {
	config := r.GetConfig()
	if model.GetTypeID() == nil {
		if config.TypeID > 0 && config.TypeID < math.MaxInt32 {
			model.SetTypeID(int32(config.TypeID))
		}
	}

	attr := model.GetAttributes()
	if model.GetID() == nil && attr != nil && attr.Name != nil {
		existing, err := r.lookupModelByName(*attr.Name)
		if err != nil {
			if !errors.Is(err, ErrCatalogModelNotFound) {
				return nil, fmt.Errorf("error finding existing model named %s: %w", *attr.Name, err)
			}
		} else {
			model.SetID(existing.ID)
		}
	}

	return r.GenericRepository.Save(model, nil)
}

func (r *CatalogModelRepositoryImpl) List(listOptions models.CatalogModelListOptions) (*dbmodels.ListWrapper[models.CatalogModel], error) {
	return r.GenericRepository.List(&listOptions)
}

func (r *CatalogModelRepositoryImpl) GetByName(name string) (models.CatalogModel, error) {
	var zeroEntity models.CatalogModel
	entity, err := r.lookupModelByName(name)
	if err != nil {
		return zeroEntity, err
	}

	config := r.GetConfig()

	// Query properties
	var properties []schema.ContextProperty
	if err := config.DB.Where(config.PropertyFieldName+" = ?", entity.ID).Find(&properties).Error; err != nil {
		return zeroEntity, fmt.Errorf("error getting properties by %s id: %w", config.EntityName, err)
	}

	// Map to domain model
	return config.SchemaToEntity(*entity, properties), nil
}

func (r *CatalogModelRepositoryImpl) lookupModelByName(name string) (*schema.Context, error) {
	var entity schema.Context

	config := r.GetConfig()

	if err := config.DB.Where("name = ? AND type_id = ?", name, config.TypeID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", config.NotFoundError, err)
		}
		return nil, fmt.Errorf("error getting %s by name: %w", config.EntityName, err)
	}

	return &entity, nil
}

func applyCatalogModelListFilters(query *gorm.DB, listOptions *models.CatalogModelListOptions) *gorm.DB {
	contextTable := utils.GetTableName(query.Statement.DB, &schema.Context{})

	if listOptions.Name != nil {
		query = query.Where(fmt.Sprintf("%s.name LIKE ?", contextTable), listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where(fmt.Sprintf("%s.external_id = ?", contextTable), listOptions.ExternalID)
	}

	if listOptions.Query != nil && *listOptions.Query != "" {
		queryPattern := fmt.Sprintf("%%%s%%", strings.ToLower(*listOptions.Query))
		propertyTable := utils.GetTableName(query.Statement.DB, &schema.ContextProperty{})

		// Search in name (context table)
		nameCondition := fmt.Sprintf("LOWER(%s.name) LIKE ?", contextTable)

		// Search in description, provider, libraryName properties
		propertyCondition := fmt.Sprintf("EXISTS (SELECT 1 FROM %s cp WHERE cp.context_id = %s.id AND cp.name IN (?, ?, ?) AND LOWER(cp.string_value) LIKE ?)",
			propertyTable, contextTable)

		// Search in tasks (assuming tasks are stored as comma-separated or multiple properties)
		tasksCondition := fmt.Sprintf("EXISTS (SELECT 1 FROM %s cp WHERE cp.context_id = %s.id AND cp.name = ? AND LOWER(cp.string_value) LIKE ?)",
			propertyTable, contextTable)

		query = query.Where(fmt.Sprintf("(%s OR %s OR %s)", nameCondition, propertyCondition, tasksCondition),
			queryPattern,                                           // for name
			"description", "provider", "libraryName", queryPattern, // for properties
			"tasks", queryPattern, // for tasks
		)
	}

	// Filter out empty strings from SourceIDs, for some reason it's passed if no sources are specified
	var nonEmptySourceIDs []string
	if listOptions.SourceIDs != nil {
		for _, sourceID := range *listOptions.SourceIDs {
			if sourceID != "" {
				nonEmptySourceIDs = append(nonEmptySourceIDs, sourceID)
			}
		}
	}

	if len(nonEmptySourceIDs) > 0 {
		propertyTable := utils.GetTableName(query.Statement.DB, &schema.ContextProperty{})

		joinClause := fmt.Sprintf("JOIN %s cp ON cp.context_id = %s.id", propertyTable, contextTable)
		query = query.Joins(joinClause).
			Where("cp.name = ? AND cp.string_value IN ?", "source_id", nonEmptySourceIDs)
	}

	return query
}

func mapCatalogModelToContext(model models.CatalogModel) schema.Context {
	attrs := model.GetAttributes()
	context := schema.Context{}

	if typeID := model.GetTypeID(); typeID != nil {
		context.TypeID = *typeID
	}

	if model.GetID() != nil {
		context.ID = *model.GetID()
	}

	if attrs != nil {
		if attrs.Name != nil {
			context.Name = *attrs.Name
		}
		context.ExternalID = attrs.ExternalID
		if attrs.CreateTimeSinceEpoch != nil {
			context.CreateTimeSinceEpoch = *attrs.CreateTimeSinceEpoch
		}
		if attrs.LastUpdateTimeSinceEpoch != nil {
			context.LastUpdateTimeSinceEpoch = *attrs.LastUpdateTimeSinceEpoch
		}
	}

	return context
}

func mapCatalogModelToContextProperties(model models.CatalogModel, contextID int32) []schema.ContextProperty {
	var properties []schema.ContextProperty

	if model.GetProperties() != nil {
		for _, prop := range *model.GetProperties() {
			properties = append(properties, service.MapPropertiesToContextProperty(prop, contextID, false))
		}
	}

	if model.GetCustomProperties() != nil {
		for _, prop := range *model.GetCustomProperties() {
			properties = append(properties, service.MapPropertiesToContextProperty(prop, contextID, true))
		}
	}

	return properties
}

func mapDataLayerToCatalogModel(modelCtx schema.Context, propertiesCtx []schema.ContextProperty) models.CatalogModel {
	catalogModel := &models.CatalogModelImpl{
		ID:     &modelCtx.ID,
		TypeID: &modelCtx.TypeID,
		Attributes: &models.CatalogModelAttributes{
			Name:                     &modelCtx.Name,
			ExternalID:               modelCtx.ExternalID,
			CreateTimeSinceEpoch:     &modelCtx.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &modelCtx.LastUpdateTimeSinceEpoch,
		},
	}

	properties := []dbmodels.Properties{}
	customProperties := []dbmodels.Properties{}

	for _, prop := range propertiesCtx {
		mappedProperty := service.MapContextPropertyToProperties(prop)

		if prop.IsCustomProperty {
			customProperties = append(customProperties, mappedProperty)
		} else {
			properties = append(properties, mappedProperty)
		}
	}

	catalogModel.Properties = &properties
	catalogModel.CustomProperties = &customProperties

	return catalogModel
}

// GetFilterableProperties returns property names and their unique values
// Only includes properties where ALL values are shorter than maxLength
func (r *CatalogModelRepositoryImpl) GetFilterableProperties(maxLength int) (map[string][]string, error) {
	config := r.GetConfig()

	// Get table names using GORM utilities for database compatibility
	contextTable := utils.GetTableName(config.DB, &schema.Context{})
	propertyTable := utils.GetTableName(config.DB, &schema.ContextProperty{})

	// Simplified query: get distinct property name/value pairs
	query := fmt.Sprintf(`
		SELECT DISTINCT cp.name, cp.string_value
		FROM %s cp
		WHERE cp.context_id IN (
			SELECT id FROM %s WHERE type_id = ?
		)
		AND cp.name IN (
			SELECT name FROM (
				SELECT name, MAX(CHAR_LENGTH(string_value)) as max_len
				FROM %s
				WHERE context_id IN (
					SELECT id FROM %s WHERE type_id = ?
				)
				AND string_value IS NOT NULL
				AND string_value != ''
				GROUP BY name
			) AS field_lengths
			WHERE max_len <= ?
		)
		AND cp.string_value IS NOT NULL
		AND cp.string_value != ''
		ORDER BY cp.name, cp.string_value
	`, propertyTable, contextTable, propertyTable, contextTable)

	type propertyRow struct {
		Name        string
		StringValue string
	}

	var rows []propertyRow
	if err := config.DB.Raw(query, config.TypeID, config.TypeID, maxLength).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("error querying filterable properties: %w", err)
	}

	// Aggregate values by property name in Go
	result := make(map[string][]string)
	for _, row := range rows {
		result[row.Name] = append(result[row.Name], row.StringValue)
	}

	return result, nil
}
