package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/db/filter"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/dbutil"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"gorm.io/gorm"
)

var ErrMcpServerNotFound = errors.New("MCP server by id not found")

// McpServerRepositoryImpl implements McpServerRepository using GORM.
type McpServerRepositoryImpl struct {
	*service.GenericRepository[models.McpServer, schema.Context, schema.ContextProperty, *models.McpServerListOptions]
}

// NewMcpServerRepository creates a new McpServerRepository.
func NewMcpServerRepository(db *gorm.DB, typeID int32) models.McpServerRepository {
	r := &McpServerRepositoryImpl{}

	r.GenericRepository = service.NewGenericRepository(service.GenericRepositoryConfig[models.McpServer, schema.Context, schema.ContextProperty, *models.McpServerListOptions]{
		DB:                      db,
		TypeID:                  typeID,
		EntityToSchema:          mapMcpServerToContext,
		SchemaToEntity:          mapDataLayerToMcpServer,
		EntityToProperties:      mapMcpServerToContextProperties,
		NotFoundError:           ErrMcpServerNotFound,
		EntityName:              "MCP server",
		PropertyFieldName:       "context_id",
		ApplyListFilters:        applyMcpServerListFilters,
		CreatePaginationToken:   r.createPaginationToken,
		ApplyCustomOrdering:     r.applyCustomOrdering,
		IsNewEntity:             func(entity models.McpServer) bool { return entity.GetID() == nil },
		HasCustomProperties:     func(entity models.McpServer) bool { return entity.GetCustomProperties() != nil },
		EntityMappingFuncs:      filter.NewCatalogEntityMappings(),
		PreserveHistoricalTimes: true, // Preserve timestamps from YAML source data
	})

	return r
}

// Save creates or updates an MCP server.
func (r *McpServerRepositoryImpl) Save(server models.McpServer) (models.McpServer, error) {
	config := r.GetConfig()
	if server.GetTypeID() == nil {
		if config.TypeID > 0 {
			server.SetTypeID(config.TypeID)
		}
	}

	attr := server.GetAttributes()
	if server.GetID() == nil && attr != nil && attr.Name != nil {
		existing, err := r.lookupServerByName(*attr.Name)
		if err != nil {
			if !errors.Is(err, ErrMcpServerNotFound) {
				return nil, fmt.Errorf("error finding existing MCP server named %s: %w", *attr.Name, err)
			}
		} else {
			server.SetID(existing.ID)
		}
	}

	return r.GenericRepository.Save(server, nil)
}

// List returns a paginated list of MCP servers.
func (r *McpServerRepositoryImpl) List(listOptions models.McpServerListOptions) (*dbmodels.ListWrapper[models.McpServer], error) {
	return r.GenericRepository.List(&listOptions)
}

// GetByName retrieves an MCP server by its name.
func (r *McpServerRepositoryImpl) GetByName(name string) (models.McpServer, error) {
	var zeroEntity models.McpServer
	entity, err := r.lookupServerByName(name)
	if err != nil {
		return zeroEntity, err
	}

	config := r.GetConfig()

	// Query properties
	var properties []schema.ContextProperty
	if err := config.DB.Where(config.PropertyFieldName+" = ?", entity.ID).Find(&properties).Error; err != nil {
		err = dbutil.SanitizeDatabaseError(err)
		return zeroEntity, fmt.Errorf("error getting properties by %s id: %w", config.EntityName, err)
	}

	// Map to domain model
	return config.SchemaToEntity(*entity, properties), nil
}

// lookupServerByName finds an MCP server by name.
func (r *McpServerRepositoryImpl) lookupServerByName(name string) (*schema.Context, error) {
	var entity schema.Context

	config := r.GetConfig()

	if err := config.DB.Where("name = ? AND type_id = ?", name, config.TypeID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", config.NotFoundError, err)
		}
		err = dbutil.SanitizeDatabaseError(err)
		return nil, fmt.Errorf("error getting %s by name: %w", config.EntityName, err)
	}

	return &entity, nil
}

// DeleteBySource deletes all MCP servers from a given source.
func (r *McpServerRepositoryImpl) DeleteBySource(sourceID string) error {
	config := r.GetConfig()

	// Delete all Context records where there's a ContextProperty with name='source_id' and string_value=sourceID
	query := `DELETE FROM "Context" WHERE id IN (
		SELECT "Context".id
		FROM "Context"
		INNER JOIN "ContextProperty" ON "Context".id="ContextProperty".context_id
		AND "ContextProperty".name='source_id'
		WHERE "ContextProperty".string_value=?
		AND "Context".type_id=?
	)`

	return config.DB.Exec(query, sourceID, config.TypeID).Error
}

// DeleteByID deletes an MCP server by its ID.
func (r *McpServerRepositoryImpl) DeleteByID(id int32) error {
	config := r.GetConfig()

	result := config.DB.Where("id = ? AND type_id = ?", id, config.TypeID).Delete(&schema.Context{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: id %d", config.NotFoundError, id)
	}

	return nil
}

// GetDistinctSourceIDs retrieves all unique source_id values from MCP servers.
func (r *McpServerRepositoryImpl) GetDistinctSourceIDs() ([]string, error) {
	config := r.GetConfig()

	var sourceIDs []string

	query := `SELECT DISTINCT cp.string_value FROM "ContextProperty" cp
		INNER JOIN "Context" c ON cp.context_id = c.id
		WHERE cp.name='source_id' AND c.type_id=?`

	rows, err := config.DB.Raw(query, config.TypeID).Rows()
	if err != nil {
		err = dbutil.SanitizeDatabaseError(err)
		return nil, fmt.Errorf("error querying distinct source IDs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sourceID string
		if err := rows.Scan(&sourceID); err != nil {
			err = dbutil.SanitizeDatabaseError(err)
			return nil, fmt.Errorf("error scanning source ID: %w", err)
		}
		sourceIDs = append(sourceIDs, sourceID)
	}

	if err := rows.Err(); err != nil {
		err = dbutil.SanitizeDatabaseError(err)
		return nil, fmt.Errorf("error iterating source ID rows: %w", err)
	}

	return sourceIDs, nil
}

// applyMcpServerListFilters applies list filters to the query.
func applyMcpServerListFilters(query *gorm.DB, listOptions *models.McpServerListOptions) *gorm.DB {
	contextTable := utils.GetTableName(query.Statement.DB, &schema.Context{})

	if listOptions.Name != nil {
		query = query.Where(fmt.Sprintf("%s.name LIKE ?", contextTable), listOptions.Name)
	}

	if listOptions.Query != nil && *listOptions.Query != "" {
		queryPattern := fmt.Sprintf("%%%s%%", strings.ToLower(*listOptions.Query))
		propertyTable := utils.GetTableName(query.Statement.DB, &schema.ContextProperty{})

		// Search in name (context table)
		nameCondition := fmt.Sprintf("LOWER(%s.name) LIKE ?", contextTable)

		// Search in description, provider properties
		propertyCondition := fmt.Sprintf("EXISTS (SELECT 1 FROM %s cp WHERE cp.context_id = %s.id AND cp.name IN (?, ?) AND LOWER(cp.string_value) LIKE ?)",
			propertyTable, contextTable)

		query = query.Where(fmt.Sprintf("(%s OR %s)", nameCondition, propertyCondition),
			queryPattern,
			"description", "provider", queryPattern,
		)
	}

	// Filter by source IDs
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

// mapMcpServerToContext maps an McpServer entity to a Context schema.
func mapMcpServerToContext(server models.McpServer) schema.Context {
	attrs := server.GetAttributes()
	context := schema.Context{}

	if typeID := server.GetTypeID(); typeID != nil {
		context.TypeID = *typeID
	}

	if server.GetID() != nil {
		context.ID = *server.GetID()
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

// mapMcpServerToContextProperties maps an McpServer entity to ContextProperty schema.
func mapMcpServerToContextProperties(server models.McpServer, contextID int32) []schema.ContextProperty {
	var properties []schema.ContextProperty

	if server.GetProperties() != nil {
		for _, prop := range *server.GetProperties() {
			properties = append(properties, service.MapPropertiesToContextProperty(prop, contextID, false))
		}
	}

	if server.GetCustomProperties() != nil {
		for _, prop := range *server.GetCustomProperties() {
			properties = append(properties, service.MapPropertiesToContextProperty(prop, contextID, true))
		}
	}

	return properties
}

// mapDataLayerToMcpServer maps database schema to an McpServer entity.
func mapDataLayerToMcpServer(serverCtx schema.Context, propertiesCtx []schema.ContextProperty) models.McpServer {
	mcpServer := &models.McpServerImpl{
		ID:     &serverCtx.ID,
		TypeID: &serverCtx.TypeID,
		Attributes: &models.McpServerAttributes{
			Name:                     &serverCtx.Name,
			ExternalID:               serverCtx.ExternalID,
			CreateTimeSinceEpoch:     &serverCtx.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &serverCtx.LastUpdateTimeSinceEpoch,
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

	mcpServer.Properties = &properties
	mcpServer.CustomProperties = &customProperties

	return mcpServer
}

// applyCustomOrdering applies custom ordering logic.
func (r *McpServerRepositoryImpl) applyCustomOrdering(query *gorm.DB, listOptions *models.McpServerListOptions) *gorm.DB {
	db := r.GetConfig().DB
	contextTable := utils.GetTableName(db, &schema.Context{})
	orderBy := listOptions.GetOrderBy()

	// Handle NAME ordering
	if orderBy == "NAME" {
		return ApplyNameOrdering(query, contextTable, listOptions.GetSortOrder(), listOptions.GetNextPageToken(), listOptions.GetPageSize())
	}

	// Fall back to standard pagination
	return r.ApplyStandardPagination(query, listOptions, []models.McpServer{})
}

// ApplyStandardPagination overrides the base implementation.
func (r *McpServerRepositoryImpl) ApplyStandardPagination(query *gorm.DB, listOptions *models.McpServerListOptions, entities any) *gorm.DB {
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

	return query.Scopes(scopes.PaginateWithOptions(entities, pagination, r.GetConfig().DB, "Context", CatalogOrderByColumns))
}

// createPaginationToken creates a pagination token for the last item.
func (r *McpServerRepositoryImpl) createPaginationToken(lastItem schema.Context, listOptions *models.McpServerListOptions) string {
	if listOptions.GetOrderBy() == "NAME" {
		return CreateNamePaginationToken(lastItem.ID, &lastItem.Name)
	}

	return r.CreateDefaultPaginationToken(lastItem, listOptions)
}

// McpOrderByColumns are the allowed orderBy columns for MCP servers.
var McpOrderByColumns = map[string]bool{
	"NAME":                        true,
	"CREATE_TIME":                 true,
	"LAST_UPDATE_TIME":            true,
	"CREATE_TIME_SINCE_EPOCH":     true,
	"LAST_UPDATE_TIME_SINCE_EPOCH": true,
}

func init() {
	glog.Infof("MCP server repository initialized")
}
