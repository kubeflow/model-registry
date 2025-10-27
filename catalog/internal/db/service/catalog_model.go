package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/db/filter"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// accuracyProperty is the property of a metrics artifact to use when sorting by accuracy.
const accuracyProperty = "overall_average"

var ErrCatalogModelNotFound = errors.New("catalog model by id not found")

type CatalogModelRepositoryImpl struct {
	*service.GenericRepository[models.CatalogModel, schema.Context, schema.ContextProperty, *models.CatalogModelListOptions]
	metricsArtifactTypeID int32
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
		ApplyCustomOrdering:   r.applyAccuracyOrdering,
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

// getMetricsArtifactTypeID looks up the type ID for CatalogMetricsArtifact dynamically
func (r *CatalogModelRepositoryImpl) getMetricsArtifactTypeID() (int32, error) {
	if r.metricsArtifactTypeID != 0 {
		return r.metricsArtifactTypeID, nil
	}

	// Look up the type ID dynamically from the database
	var typeRecord struct {
		ID int32 `gorm:"column:id"`
	}

	err := r.GetConfig().DB.
		Table("\"Type\"").
		Select("id").
		Where("name = ?", CatalogMetricsArtifactTypeName).
		First(&typeRecord).Error

	if err != nil {
		return 0, fmt.Errorf("failed to lookup CatalogMetricsArtifact type ID: %w", err)
	}

	// Cache the type ID for future use
	r.metricsArtifactTypeID = typeRecord.ID
	return typeRecord.ID, nil
}

// applyAccuracyOrdering applies custom ordering logic for ACCURACY orderBy field
func (r *CatalogModelRepositoryImpl) applyAccuracyOrdering(query *gorm.DB, listOptions *models.CatalogModelListOptions) *gorm.DB {
	orderBy := listOptions.GetOrderBy()

	// Only apply custom ordering for ACCURACY orderBy
	if orderBy != "ACCURACY" {
		// Fall back to standard pagination for non-ACCURACY ordering
		return r.ApplyStandardPagination(query, listOptions, []models.CatalogModel{})
	}

	// Get the metrics artifact type ID
	metricsTypeID, err := r.getMetricsArtifactTypeID()
	if err != nil {
		// Fall back to standard pagination if we can't get the type ID
		return r.ApplyStandardPagination(query, listOptions, []models.CatalogModel{})
	}

	db := r.GetConfig().DB
	contextTable := utils.GetTableName(db, &schema.Context{})
	attributionTable := utils.GetTableName(db, &schema.Attribution{})
	artifactTable := utils.GetTableName(db, &schema.Artifact{})
	propertyTable := utils.GetTableName(db, &schema.ArtifactProperty{})

	sortOrder := listOptions.GetSortOrder()
	pageSize := listOptions.GetPageSize()

	// Build the accuracy subquery
	// This gets the accuracy score for each model from its AccuracyMetric artifacts
	accuracySubquery := db.
		Select(fmt.Sprintf("%s.id, max(%s.double_value) AS accuracy", contextTable, propertyTable)).
		Table(contextTable).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.id=%s.context_id", attributionTable, contextTable, attributionTable)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.artifact_id=%s.id AND %s.type_id=?", artifactTable, attributionTable, artifactTable, artifactTable), metricsTypeID).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.id=%s.artifact_id AND %s.name=?", propertyTable, artifactTable, propertyTable, propertyTable), accuracyProperty).
		Where(contextTable+".type_id=?", r.GetConfig().TypeID).
		Group(contextTable + ".id")

	// Join the main query with the accuracy subquery
	query = query.
		Joins("LEFT JOIN (?) accuracy ON \"Context\".id=accuracy.id", accuracySubquery)

	// Apply sorting order
	if sortOrder == "ASC" {
		query = query.Order("accuracy ASC NULLS LAST")
	} else {
		// Default to DESC for ACCURACY sorting
		query = query.Order("accuracy DESC NULLS LAST")
	}

	// Handle cursor-based pagination with nextPageToken
	nextPageToken := listOptions.GetNextPageToken()
	if nextPageToken != "" {
		// Parse the cursor from the token
		if cursor, err := r.parseNextPageToken(nextPageToken); err == nil {
			// Apply WHERE clause for cursor-based pagination with ACCURACY
			query = r.applyCursorPagination(query, cursor, sortOrder)
		}
		// If token parsing fails, fall back to no cursor (first page)
	}

	// Apply pagination limit
	if pageSize > 0 {
		query = query.Limit(int(pageSize) + 1) // +1 to detect if there are more pages
	}

	return query
}

// cursor represents a pagination cursor with ID and accuracy value
type accuracyCursor struct {
	ID       int32
	Accuracy *float64
}

// parseNextPageToken decodes a nextPageToken and extracts the cursor information
func (r *CatalogModelRepositoryImpl) parseNextPageToken(token string) (*accuracyCursor, error) {
	// Sanity check the length before decoding
	if len(token) > 64 {
		return nil, fmt.Errorf("invalid nextPageToken")
	}

	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cursor format, expected 'ID:Value'")
	}

	id, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid ID in cursor: %w", err)
	}

	cursor := accuracyCursor{ID: int32(id)}

	// Parse accuracy value from cursor
	accuracy, err := strconv.ParseFloat(parts[1], 64)
	if err == nil {
		cursor.Accuracy = &accuracy
	}

	return &cursor, nil
}

// applyCursorPagination applies WHERE clause for cursor-based pagination with ACCURACY sorting
func (r *CatalogModelRepositoryImpl) applyCursorPagination(query *gorm.DB, cursor *accuracyCursor, sortOrder string) *gorm.DB {
	contextTable := utils.GetTableName(query, &schema.Context{})

	// Handle NULL accuracy values in cursor
	if cursor.Accuracy == nil {
		// For models without accuracy, just use ID-based pagination
		if sortOrder == "ASC" {
			return query.Where(contextTable+".id > ?", cursor.ID)
		}
		return query.Where(contextTable+".id > ?", cursor.ID) // In DESC, NULL comes last, so ID ordering is fine
	}

	accuracyValue := *cursor.Accuracy

	// Apply cursor pagination logic for ACCURACY sorting
	if sortOrder == "ASC" {
		// For ASC: get records where (accuracy > cursor_accuracy) OR (accuracy = cursor_accuracy AND id > cursor_id)
		// Also include NULL values at the end
		return query.Where("(accuracy > ? OR (accuracy = ? AND "+contextTable+".id > ?) OR accuracy IS NULL)",
			accuracyValue, accuracyValue, cursor.ID)
	} else {
		// For DESC: get records where (accuracy < cursor_accuracy) OR (accuracy = cursor_accuracy AND id > cursor_id)
		return query.Where("(accuracy < ? OR (accuracy = ? AND "+contextTable+".id > ?))",
			accuracyValue, accuracyValue, cursor.ID)
	}
}

func (r *CatalogModelRepositoryImpl) createPaginationToken(lastItem schema.Context, listOptions *models.CatalogModelListOptions) string {
	if listOptions.GetOrderBy() == "ACCURACY" {
		// The accuracy metric is not available from the context table,
		// so we'll need another query to get it.

		db := r.GetConfig().DB
		contextTable := utils.GetTableName(db, &schema.Context{})
		attributionTable := utils.GetTableName(db, &schema.Attribution{})
		artifactTable := utils.GetTableName(db, &schema.Artifact{})
		propertyTable := utils.GetTableName(db, &schema.ArtifactProperty{})
		metricsTypeID, err := r.getMetricsArtifactTypeID()
		if err != nil {
			glog.Warningf("Failed to get metrics artifact type ID: %v", err)
			return r.CreateDefaultPaginationToken(lastItem, listOptions)
		}

		query := db.
			Select("MAX(double_value) AS accuracy").
			Table(contextTable).
			Joins(fmt.Sprintf("LEFT JOIN %s ON %s.id=%s.context_id", attributionTable, contextTable, attributionTable)).
			Joins(fmt.Sprintf("LEFT JOIN %s ON %s.artifact_id=%s.id", artifactTable, attributionTable, artifactTable)).
			Joins(fmt.Sprintf("LEFT JOIN %s ON %s.id=%s.artifact_id", propertyTable, artifactTable, propertyTable)).
			Where(artifactTable+".type_id=?", metricsTypeID).
			Where(propertyTable+".name=?", accuracyProperty).
			Where(contextTable+".id=?", lastItem.ID)

		var result struct {
			Accuracy *float64 `gorm:"accuracy"`
		}
		err = query.Scan(&result).Error
		if err != nil {
			glog.Warningf("Failed to get accuracy score: %v", err)
			return r.CreateDefaultPaginationToken(lastItem, listOptions)
		}

		return createAccuracyPaginationToken(lastItem.ID, result.Accuracy)
	}

	return r.CreateDefaultPaginationToken(lastItem, listOptions)
}

// createAccuracyPaginationToken creates a pagination token for ACCURACY sorting
func createAccuracyPaginationToken(entityID int32, accuracyValue *float64) string {
	var valueStr string
	if accuracyValue != nil {
		valueStr = fmt.Sprintf("%.15f", *accuracyValue)
	} else {
		valueStr = "" // Represents NULL
	}

	cursor := fmt.Sprintf("%d:%s", entityID, valueStr)
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}
