package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/db/filter"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var ErrCatalogModelNotFound = errors.New("catalog model by id not found")

type CatalogModelRepositoryImpl struct {
	*service.GenericRepository[models.CatalogModel, schema.Context, schema.ContextProperty, *models.CatalogModelListOptions]
}

func NewCatalogModelRepository(db *gorm.DB, typeID int32) models.CatalogModelRepository {
	r := &CatalogModelRepositoryImpl{}

	r.GenericRepository = service.NewGenericRepository(service.GenericRepositoryConfig[models.CatalogModel, schema.Context, schema.ContextProperty, *models.CatalogModelListOptions]{
		DB:                    db,
		TypeID:                typeID,
		EntityToSchema:        mapCatalogModelToContext,
		SchemaToEntity:        mapDataLayerToCatalogModel,
		EntityToProperties:    mapCatalogModelToContextProperties,
		NotFoundError:         ErrCatalogModelNotFound,
		EntityName:            "catalog model",
		PropertyFieldName:     "context_id",
		ApplyListFilters:      applyCatalogModelListFilters,
		CreatePaginationToken: r.createPaginationToken,
		ApplyCustomOrdering:   r.applyCustomOrdering,
		IsNewEntity:           func(entity models.CatalogModel) bool { return entity.GetID() == nil },
		HasCustomProperties:   func(entity models.CatalogModel) bool { return entity.GetCustomProperties() != nil },
		EntityMappingFuncs:    filter.NewCatalogEntityMappings(),
	})

	return r
}

func (r *CatalogModelRepositoryImpl) Save(model models.CatalogModel) (models.CatalogModel, error) {
	config := r.GetConfig()
	if model.GetTypeID() == nil {
		if config.TypeID > 0 {
			model.SetTypeID(config.TypeID)
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

// ApplyStandardPagination overrides the base implementation to use catalog-specific allowed columns
func (r *CatalogModelRepositoryImpl) ApplyStandardPagination(query *gorm.DB, listOptions *models.CatalogModelListOptions, entities any) *gorm.DB {
	pageSize := listOptions.GetPageSize()
	orderBy := listOptions.GetOrderBy()
	sortOrder := listOptions.GetSortOrder()
	nextPageToken := listOptions.GetNextPageToken()

	pagination := &dbmodels.Pagination{
		PageSize:      &pageSize,
		OrderBy:       &orderBy,
		SortOrder:     &sortOrder,
		NextPageToken: &nextPageToken,
	}

	// Use catalog-specific allowed columns (includes NAME)
	return query.Scopes(scopes.PaginateWithOptions(entities, pagination, r.GetConfig().DB, "Context", CatalogOrderByColumns))
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

	if config.DB.Name() != "postgres" {
		return nil, fmt.Errorf("GetFilterableProperties is only supported on PostgreSQL")
	}

	db, err := config.DB.DB()
	if err != nil {
		return nil, err
	}

	// Get table names using GORM utilities for database compatibility
	contextTable := utils.GetTableName(config.DB, &schema.Context{})
	propertyTable := utils.GetTableName(config.DB, &schema.ContextProperty{})

	query := fmt.Sprintf(`
		SELECT name, array_agg(string_value) FROM (
			SELECT
				name,
				string_value
			FROM %s WHERE
				context_id IN (
					SELECT id FROM %s WHERE type_id=$1
				)
				AND string_value IS NOT NULL
				AND string_value != ''
				AND string_value IS NOT JSON ARRAY

			UNION

			SELECT
				name,
				json_array_elements_text(string_value::json) AS string_value
			FROM %s WHERE
				context_id IN (
					SELECT id FROM %s WHERE type_id=$1
				)
				AND string_value IS JSON ARRAY
		)
		GROUP BY name HAVING MAX(CHAR_LENGTH(string_value)) <= $2
	`, propertyTable, contextTable, propertyTable, contextTable)

	rows, err := db.Query(query, config.TypeID, maxLength)
	if err != nil {
		return nil, fmt.Errorf("error querying filterable properties: %w", err)
	}
	defer rows.Close()

	result := map[string][]string{}
	for rows.Next() {
		var name string
		var values pq.StringArray

		err = rows.Scan(&name, &values)
		if err != nil {
			return nil, fmt.Errorf("error scanning filterable property row: %w", err)
		}

		result[name] = []string(values)
	}

	return result, nil
}

// applyCustomOrdering applies custom ordering logic for non-standard orderBy field
func (r *CatalogModelRepositoryImpl) applyCustomOrdering(query *gorm.DB, listOptions *models.CatalogModelListOptions) *gorm.DB {

	db := r.GetConfig().DB
	contextTable := utils.GetTableName(db, &schema.Context{})
	orderBy := listOptions.GetOrderBy()

	// Handle NAME ordering specially (catalog-specific)
	if orderBy == "NAME" {
		return ApplyNameOrdering(query, contextTable, listOptions.GetSortOrder(), listOptions.GetNextPageToken(), listOptions.GetPageSize())
	}

	subquery, sortColumn := r.sortValueQuery(listOptions, contextTable+".id")
	if subquery == nil {
		// Fall back to standard pagination with catalog-specific allowed columns
		return r.ApplyStandardPagination(query, listOptions, []models.CatalogModel{})
	}
	subquery = subquery.Group(contextTable + ".id")

	// Join the main query with the subquery
	query = query.
		Joins(fmt.Sprintf("LEFT JOIN (?) sort_value ON %s.id=sort_value.id", contextTable), subquery)

	// Apply sorting order
	sortOrder := listOptions.GetSortOrder()
	if sortOrder != "ASC" {
		sortOrder = "DESC"
	}
	query = query.Order(fmt.Sprintf("sort_value.%s %s NULLS LAST, %s.id", sortColumn, sortOrder, contextTable))

	// Handle cursor-based pagination with nextPageToken
	nextPageToken := listOptions.GetNextPageToken()
	if nextPageToken != "" {
		// Parse the cursor from the token
		if cursor, err := scopes.DecodeCursor(nextPageToken); err == nil {
			// Apply WHERE clause for cursor-based pagination with ACCURACY
			query = r.applyCursorPagination(query, cursor, sortColumn, sortOrder)
		}
		// If token parsing fails, fall back to no cursor (first page)
	}

	// Apply pagination limit
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		query = query.Limit(int(pageSize) + 1) // +1 to detect if there are more pages
	}

	return query
}

// applyCursorPagination applies WHERE clause for cursor-based pagination with ACCURACY sorting
func (r *CatalogModelRepositoryImpl) applyCursorPagination(query *gorm.DB, cursor *scopes.Cursor, sortColumn, sortOrder string) *gorm.DB {
	contextTable := utils.GetTableName(query, &schema.Context{})

	// Handle NULL values in cursor
	if cursor.Value == "" {
		// Items without the sort value will be sorted to the bottom, just use ID-based pagination.
		return query.Where(fmt.Sprintf("sort_value.%s IS NULL AND %s.id > ?", sortColumn, contextTable), cursor.ID)
	}

	cmp := "<"
	if sortOrder == "ASC" {
		cmp = ">"
	}

	// Note that we sort ID ASCENDING as a tie-breaker, so ">" is correct below.
	return query.Where(fmt.Sprintf("(sort_value.%s %s ? OR (sort_value.%s = ? AND %s.id > ?) OR sort_value.%s IS NULL)", sortColumn, cmp, sortColumn, contextTable, sortColumn),
		cursor.Value, cursor.Value, cursor.ID)
}

func (r *CatalogModelRepositoryImpl) createPaginationToken(lastItem schema.Context, listOptions *models.CatalogModelListOptions) string {
	// Handle NAME ordering (catalog-specific)
	if listOptions.GetOrderBy() == "NAME" {
		return CreateNamePaginationToken(lastItem.ID, &lastItem.Name)
	}

	sortValueQuery, column := r.sortValueQuery(listOptions)
	if sortValueQuery != nil {
		contextTable := utils.GetTableName(r.GetConfig().DB, &schema.Context{})
		sortValueQuery = sortValueQuery.Where(contextTable+".id=?", lastItem.ID)

		var result struct {
			IntValue    *int64   `gorm:"int_value"`
			DoubleValue *float64 `gorm:"double_value"`
			StringValue *string  `gorm:"string_value"`
		}
		err := sortValueQuery.Scan(&result).Error
		if err != nil {
			glog.Warningf("Failed to get sort value: %v", err)
		} else {
			switch column {
			case "int_value":
				return scopes.CreateNextPageToken(lastItem.ID, result.IntValue)
			case "double_value":
				return scopes.CreateNextPageToken(lastItem.ID, result.DoubleValue)
			case "string_value":
				fallthrough
			default:
				return scopes.CreateNextPageToken(lastItem.ID, result.StringValue)
			}
		}
	}

	return r.CreateDefaultPaginationToken(lastItem, listOptions)
}

// sortValueQuery returns a query that will produce the value to sort on for
// the List response. The returned string is the column name.
//
// If the sort does not require a subquery, sortValueQuery returns nil.
func (r *CatalogModelRepositoryImpl) sortValueQuery(listOptions *models.CatalogModelListOptions, extraColumns ...any) (*gorm.DB, string) {
	db := r.GetConfig().DB
	contextTable := utils.GetTableName(db, &schema.Context{})

	query := db.Table(contextTable).
		Where(contextTable+".type_id=?", r.GetConfig().TypeID)

	orderBy := strings.Split(listOptions.GetOrderBy(), ".")

	var valueColumn string

	switch {
	case len(orderBy) == 3 && orderBy[0] == "artifacts":
		// artifacts.<property>.<value_column> e.g. artifacts.ttft_p90.double_value

		attributionTable := utils.GetTableName(db, &schema.Attribution{})
		propertyTable := utils.GetTableName(db, &schema.ArtifactProperty{})

		aggFn := "max"
		if listOptions.GetSortOrder() == "ASC" {
			aggFn = "min"
		}
		valueColumn = orderBy[2]

		query = query.
			Select(fmt.Sprintf("%s(%s.%s) AS %s", aggFn, propertyTable, valueColumn, valueColumn), extraColumns...).
			Joins(fmt.Sprintf("LEFT JOIN %s ON %s.id=%s.context_id", attributionTable, contextTable, attributionTable)).
			Joins(fmt.Sprintf("LEFT JOIN %s ON %s.artifact_id=%s.artifact_id AND %s.name=?", propertyTable, attributionTable, propertyTable, propertyTable), orderBy[1])
	case len(orderBy) == 2:
		// <property>.<value_column> e.g. provider.string_value
		propertyTable := utils.GetTableName(db, &schema.ContextProperty{})
		valueColumn = orderBy[1]
		query = query.
			Select(fmt.Sprintf("max(%s.%s) AS %s", propertyTable, valueColumn, valueColumn), extraColumns...).
			Joins(fmt.Sprintf("LEFT JOIN %s ON %s.id=%s.context_id AND %s.name=?", propertyTable, contextTable, propertyTable, propertyTable), orderBy[0])
	default:
		// Standard sort will work
		return nil, ""
	}

	// The query is built, but verify that the value column is valid before
	// returning it.
	switch valueColumn {
	case "int_value", "double_value", "string_value":
		// OK
	default:
		return nil, ""
	}

	return query, valueColumn
}
