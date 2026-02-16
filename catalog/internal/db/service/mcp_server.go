package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/db/filter"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/db/dbutil"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"gorm.io/gorm"
)

var (
	ErrMCPServerNotFound       = errors.New("MCP server by id not found")
	ErrInvalidBaseName         = errors.New("base_name validation failed")
	ErrBaseNameContainsAtSign  = errors.New("base_name cannot contain '@' character")
	ErrBaseNameEmpty           = errors.New("base_name cannot be empty")
	ErrBaseNameTooLong         = errors.New("base_name exceeds maximum length of 255 characters")
	ErrVersionTooLong          = errors.New("version exceeds maximum length of 100 characters")
	ErrVersionContainsAtSign   = errors.New("version cannot contain '@' character")
)

// MCPServerRepositoryImpl implements MCPServerRepository using GORM.
type MCPServerRepositoryImpl struct {
	*service.GenericRepository[models.MCPServer, schema.Context, schema.ContextProperty, *models.MCPServerListOptions]
}

// NewMCPServerRepository creates a new MCPServerRepository.
func NewMCPServerRepository(db *gorm.DB, typeID int32) models.MCPServerRepository {
	r := &MCPServerRepositoryImpl{}

	r.GenericRepository = service.NewGenericRepository(service.GenericRepositoryConfig[models.MCPServer, schema.Context, schema.ContextProperty, *models.MCPServerListOptions]{
		DB:                      db,
		TypeID:                  typeID,
		EntityToSchema:          mapMCPServerToContext,
		SchemaToEntity:          mapDataLayerToMCPServer,
		EntityToProperties:      mapMCPServerToContextProperties,
		NotFoundError:           ErrMCPServerNotFound,
		EntityName:              "MCP server",
		PropertyFieldName:       "context_id",
		ApplyListFilters:        applyMCPServerListFilters,
		CreatePaginationToken:   r.createPaginationToken,
		ApplyCustomOrdering:     r.applyCustomOrdering,
		IsNewEntity:             func(entity models.MCPServer) bool { return entity.GetID() == nil },
		HasCustomProperties:     func(entity models.MCPServer) bool { return entity.GetCustomProperties() != nil },
		EntityMappingFuncs:      filter.NewCatalogEntityMappings(),
		PreserveHistoricalTimes: true, // Preserve timestamps from YAML source data
	})

	return r
}

// Save creates or updates an MCP server.
// Uses (base_name, version) as the unique identifier.
// Stores composite name (base_name@version) in Context.name field.
func (r *MCPServerRepositoryImpl) Save(server models.MCPServer) (models.MCPServer, error) {
	config := r.GetConfig()
	if server.GetTypeID() == nil {
		if config.TypeID > 0 {
			server.SetTypeID(config.TypeID)
		}
	}

	attr := server.GetAttributes()
	if attr != nil && attr.Name != nil {
		baseName := strings.TrimSpace(*attr.Name)
		version := extractVersionProperty(server)

		// Validate base_name
		if baseName == "" {
			return nil, ErrBaseNameEmpty
		}
		if len(baseName) > 255 {
			return nil, fmt.Errorf("%w: length %d", ErrBaseNameTooLong, len(baseName))
		}
		if strings.Contains(baseName, "@") {
			return nil, fmt.Errorf("%w: %s", ErrBaseNameContainsAtSign, baseName)
		}

		// Validate version
		if len(version) > 100 {
			return nil, fmt.Errorf("%w: length %d", ErrVersionTooLong, len(version))
		}
		if strings.Contains(version, "@") {
			return nil, fmt.Errorf("%w: %s", ErrVersionContainsAtSign, version)
		}

		// Build composite name for Context.name field
		compositeName := buildCompositeName(baseName, version)

		// Set the composite name in attributes
		attr.Name = &compositeName

		// Add or update base_name property
		props := server.GetProperties()
		if props != nil {
			hasBaseName := false
			var baseNameProp *dbmodels.Properties
			for i := range *props {
				if (*props)[i].Name == "base_name" {
					hasBaseName = true
					baseNameProp = &(*props)[i]
					break
				}
			}
			if hasBaseName {
				// Update if different
				if baseNameProp.StringValue == nil || *baseNameProp.StringValue != baseName {
					baseNameProp.StringValue = &baseName
				}
			} else {
				// Add new property
				*props = append(*props, dbmodels.Properties{
					Name:        "base_name",
					StringValue: &baseName,
				})
			}
		} else {
			// Initialize properties if nil
			newProps := []dbmodels.Properties{
				{
					Name:        "base_name",
					StringValue: &baseName,
				},
			}
			// Use type assertion to set properties on the impl
			if impl, ok := server.(*models.MCPServerImpl); ok {
				impl.Properties = &newProps
			}
		}

		// Check for existing server with same (base_name, version)
		if server.GetID() == nil {
			existing, err := r.GetByNameAndVersion(baseName, version)
			if err != nil {
				if !errors.Is(err, ErrMCPServerNotFound) {
					return nil, fmt.Errorf("error finding existing MCP server named %s version %s: %w", baseName, version, err)
				}
				// If not found, continue with create
			} else {
				// Found existing - update it
				if existing.GetID() != nil {
					server.SetID(*existing.GetID())
				}
			}
		}
	}

	// Attempt to save the server
	saved, err := r.GenericRepository.Save(server, nil)
	if err != nil {
		// Handle race condition: if unique constraint violation occurs,
		// retry by fetching the existing record and updating it
		if dbutil.IsDuplicateKeyError(err) && attr != nil && attr.Name != nil {
			// Extract base name and version from the composite name
			compositeName := *attr.Name
			baseName, version := parseCompositeName(compositeName)

			// Try to get the existing server
			existing, getErr := r.GetByNameAndVersion(baseName, version)
			if getErr == nil && existing.GetID() != nil {
				// Found it - set the ID and retry the save as an update
				server.SetID(*existing.GetID())
				return r.GenericRepository.Save(server, nil)
			}
		}
		return nil, err
	}

	return saved, nil
}

// extractVersionProperty extracts the version property value from an MCP server.
// Returns empty string if no version property exists.
func extractVersionProperty(server models.MCPServer) string {
	if server.GetProperties() == nil {
		return ""
	}

	for _, prop := range *server.GetProperties() {
		if prop.Name == "version" && prop.StringValue != nil {
			return *prop.StringValue
		}
	}

	return ""
}

// buildCompositeName constructs the composite name (name@version) for storage in Context.name.
// If version is empty, returns just the base name.
func buildCompositeName(baseName string, version string) string {
	if version == "" {
		return baseName
	}
	return fmt.Sprintf("%s@%s", baseName, version)
}

// parseCompositeName parses a composite name into base name and version.
// Returns (baseName, version). If no @ separator, returns (name, "").
//
// Note: The validation logic (Save method) now forbids "@" in both base_name and version,
// so this function will only ever encounter a single "@" separator in valid data.
// The multi-@ handling is preserved for backward compatibility with any historical data.
func parseCompositeName(compositeName string) (string, string) {
	parts := strings.Split(compositeName, "@")
	if len(parts) == 1 {
		return compositeName, ""
	}
	// Handle case like "name@v1.0@extra" - take first @ as separator
	// (This case should not occur in new data due to validation)
	baseName := parts[0]
	version := strings.Join(parts[1:], "@")
	return baseName, version
}

// List returns a paginated list of MCP servers.
func (r *MCPServerRepositoryImpl) List(listOptions models.MCPServerListOptions) (*dbmodels.ListWrapper[models.MCPServer], error) {
	return r.GenericRepository.List(&listOptions)
}

// GetByNameAndVersion retrieves an MCP server by its base name and version.
// Uses composite name (base_name@version) stored in Context.name field.
func (r *MCPServerRepositoryImpl) GetByNameAndVersion(baseName string, version string) (models.MCPServer, error) {
	var zeroEntity models.MCPServer
	entity, err := r.lookupServerByNameAndVersion(baseName, version)
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

// lookupServerByNameAndVersion finds an MCP server by base name and version using composite name.
func (r *MCPServerRepositoryImpl) lookupServerByNameAndVersion(baseName string, version string) (*schema.Context, error) {
	config := r.GetConfig()

	// Build composite name
	compositeName := buildCompositeName(baseName, version)

	var entity schema.Context
	if err := config.DB.Where("name = ? AND type_id = ?", compositeName, config.TypeID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: name=%s, version=%s", config.NotFoundError, baseName, version)
		}
		err = dbutil.SanitizeDatabaseError(err)
		return nil, fmt.Errorf("error getting %s by name and version: %w", config.EntityName, err)
	}

	return &entity, nil
}

// DeleteBySource deletes all MCP servers from a given source.
func (r *MCPServerRepositoryImpl) DeleteBySource(sourceID string) error {
	config := r.GetConfig()

	// Build subquery to find matching context IDs
	subQuery := config.DB.Table(utils.GetTableName(config.DB, &schema.Context{})).
		Select(utils.GetTableName(config.DB, &schema.Context{}) + ".id").
		Joins("INNER JOIN "+utils.GetTableName(config.DB, &schema.ContextProperty{})+" ON "+
			utils.GetTableName(config.DB, &schema.Context{})+".id = "+
			utils.GetTableName(config.DB, &schema.ContextProperty{})+".context_id").
		Where(utils.GetTableName(config.DB, &schema.ContextProperty{})+".name = ? AND "+
			utils.GetTableName(config.DB, &schema.ContextProperty{})+".string_value = ? AND "+
			utils.GetTableName(config.DB, &schema.Context{})+".type_id = ?",
			"source_id", sourceID, config.TypeID)

	// Delete contexts with matching IDs
	return config.DB.Where("id IN (?)", subQuery).Delete(&schema.Context{}).Error
}

// DeleteByID deletes an MCP server by its ID.
func (r *MCPServerRepositoryImpl) DeleteByID(id int32) error {
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
func (r *MCPServerRepositoryImpl) GetDistinctSourceIDs() ([]string, error) {
	config := r.GetConfig()

	var sourceIDs []string

	contextPropertyTable := utils.GetTableName(config.DB, &schema.ContextProperty{})
	contextTable := utils.GetTableName(config.DB, &schema.Context{})

	err := config.DB.Table(contextPropertyTable+" cp").
		Select("DISTINCT cp.string_value").
		Joins("INNER JOIN "+contextTable+" c ON cp.context_id = c.id").
		Where("cp.name = ? AND c.type_id = ?", "source_id", config.TypeID).
		Pluck("string_value", &sourceIDs).Error

	if err != nil {
		err = dbutil.SanitizeDatabaseError(err)
		return nil, fmt.Errorf("error querying distinct source IDs: %w", err)
	}

	return sourceIDs, nil
}

// applyMCPServerListFilters applies list filters to the query.
func applyMCPServerListFilters(query *gorm.DB, listOptions *models.MCPServerListOptions) *gorm.DB {
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

// mapMCPServerToContext maps an MCPServer entity to a Context schema.
func mapMCPServerToContext(server models.MCPServer) schema.Context {
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

// mapMCPServerToContextProperties maps an MCPServer entity to ContextProperty schema.
func mapMCPServerToContextProperties(server models.MCPServer, contextID int32) []schema.ContextProperty {
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

// mapDataLayerToMCPServer maps database schema to an MCPServer entity.
func mapDataLayerToMCPServer(serverCtx schema.Context, propertiesCtx []schema.ContextProperty) models.MCPServer {
	// Parse composite name to get base name
	baseName, _ := parseCompositeName(serverCtx.Name)

	mcpServer := &models.MCPServerImpl{
		ID:     &serverCtx.ID,
		TypeID: &serverCtx.TypeID,
		Attributes: &models.MCPServerAttributes{
			Name:                     &baseName, // Use base name, not composite
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
func (r *MCPServerRepositoryImpl) applyCustomOrdering(query *gorm.DB, listOptions *models.MCPServerListOptions) *gorm.DB {
	db := r.GetConfig().DB
	contextTable := utils.GetTableName(db, &schema.Context{})
	orderBy := listOptions.GetOrderBy()

	// Handle NAME ordering
	if orderBy == "NAME" {
		return ApplyNameOrdering(query, contextTable, listOptions.GetSortOrder(), listOptions.GetNextPageToken(), listOptions.GetPageSize())
	}

	// Fall back to standard pagination
	return r.ApplyStandardPagination(query, listOptions, []models.MCPServer{})
}

// ApplyStandardPagination overrides the base implementation.
func (r *MCPServerRepositoryImpl) ApplyStandardPagination(query *gorm.DB, listOptions *models.MCPServerListOptions, entities any) *gorm.DB {
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

	return query.Scopes(scopes.PaginateWithOptions(entities, pagination, r.GetConfig().DB, "Context", McpOrderByColumns))
}

// createPaginationToken creates a pagination token for the last item.
func (r *MCPServerRepositoryImpl) createPaginationToken(lastItem schema.Context, listOptions *models.MCPServerListOptions) string {
	if listOptions.GetOrderBy() == "NAME" {
		return CreateNamePaginationToken(lastItem.ID, &lastItem.Name)
	}

	return r.CreateDefaultPaginationToken(lastItem, listOptions)
}

// McpOrderByColumns are the allowed orderBy columns for MCP servers.
var McpOrderByColumns = map[string]string{
	"ID":               "id",
	"CREATE_TIME":      "create_time_since_epoch",
	"LAST_UPDATE_TIME": "last_update_time_since_epoch",
	"NAME":             "name",
	"id":               "id", // default fallback
}

// ==============================================================================
// OpenAPI <-> Database Model Converters
// ==============================================================================

// ConvertOpenapiMCPServerToDb converts an OpenAPI MCPServer to a database MCPServer model.
// This sets all properties from the OpenAPI struct into the database model's Properties field.
func ConvertOpenapiMCPServerToDb(openapiServer *openapi.MCPServer) models.MCPServer {
	dbServer := &models.MCPServerImpl{
		Attributes: &models.MCPServerAttributes{
			Name:                     &openapiServer.Name,
			ExternalID:               openapiServer.ExternalId,
			CreateTimeSinceEpoch:     parseEpochString(openapiServer.CreateTimeSinceEpoch),
			LastUpdateTimeSinceEpoch: parseEpochString(openapiServer.LastUpdateTimeSinceEpoch),
		},
	}

	// Build properties array from OpenAPI fields
	properties := []dbmodels.Properties{}

	// Simple string fields
	addStringProperty(&properties, "source_id", openapiServer.SourceId)
	addStringProperty(&properties, "provider", openapiServer.Provider)
	addStringProperty(&properties, "logo", openapiServer.Logo)
	addStringProperty(&properties, "version", openapiServer.Version)
	addStringProperty(&properties, "license", openapiServer.License)
	addStringProperty(&properties, "license_link", openapiServer.LicenseLink)
	addStringProperty(&properties, "readme", openapiServer.Readme)
	addStringProperty(&properties, "deploymentMode", openapiServer.DeploymentMode)
	addStringProperty(&properties, "documentationUrl", openapiServer.DocumentationUrl)
	addStringProperty(&properties, "repositoryUrl", openapiServer.RepositoryUrl)
	addStringProperty(&properties, "sourceCode", openapiServer.SourceCode)
	addStringProperty(&properties, "description", openapiServer.Description)

	// Array fields (stored as JSON)
	addArrayProperty(&properties, "tags", openapiServer.Tags)
	addArrayProperty(&properties, "transports", openapiServer.Transports)

	// Time fields (convert time.Time to ISO 8601 string for content timestamps)
	if publishedDate, ok := openapiServer.GetPublishedDateOk(); ok && publishedDate != nil {
		timeStr := formatTimeToISO8601(publishedDate)
		addStringProperty(&properties, "publishedDate", &timeStr)
	}
	if lastUpdated, ok := openapiServer.GetLastUpdatedOk(); ok && lastUpdated != nil {
		timeStr := formatTimeToISO8601(lastUpdated)
		addStringProperty(&properties, "lastUpdated", &timeStr)
	}

	// Security indicators (store as individual boolean properties)
	if securityIndicators, ok := openapiServer.GetSecurityIndicatorsOk(); ok && securityIndicators != nil {
		addBoolProperty(&properties, "verifiedSource", securityIndicators.VerifiedSource)
		addBoolProperty(&properties, "secureEndpoint", securityIndicators.SecureEndpoint)
		addBoolProperty(&properties, "sast", securityIndicators.Sast)
		addBoolProperty(&properties, "readOnlyTools", securityIndicators.ReadOnlyTools)
	}

	// Complex objects (store as JSON)
	if endpoints, ok := openapiServer.GetEndpointsOk(); ok && endpoints != nil {
		addJSONProperty(&properties, "endpoints", endpoints)
	}
	if len(openapiServer.Artifacts) > 0 {
		addJSONProperty(&properties, "artifacts", openapiServer.Artifacts)
	}
	if runtimeMetadata, ok := openapiServer.GetRuntimeMetadataOk(); ok && runtimeMetadata != nil {
		addJSONProperty(&properties, "runtimeMetadata", runtimeMetadata)
	}

	// Set properties on the server
	dbServer.Properties = &properties

	// Handle custom properties
	if len(openapiServer.CustomProperties) > 0 {
		customProps := []dbmodels.Properties{}
		for name, value := range openapiServer.CustomProperties {
			customProps = append(customProps, convertMetadataValueToProperty(name, value))
		}
		dbServer.CustomProperties = &customProps
	}

	return dbServer
}

// ConvertDbMCPServerToOpenapi converts a database MCPServer model to an OpenAPI MCPServer.
// This extracts all properties from the database model and populates the OpenAPI struct.
//
// NOTE: The returned MCPServer will have ToolCount=0. Use ConvertDbMCPServerWithToolsToOpenapi
// if you have loaded the associated tools and need an accurate tool count.
func ConvertDbMCPServerToOpenapi(dbServer models.MCPServer) *openapi.MCPServer {
	return convertDbMCPServerToOpenapiInternal(dbServer, nil)
}

// ConvertDbMCPServerWithToolsToOpenapi converts a database MCPServer to OpenAPI representation
// with accurate tool count and tools array populated from the provided tools.
//
// Parameters:
//   - dbServer: The database server entity to convert
//   - tools: Associated tools (can be nil or empty for toolCount=0)
//
// This is the preferred method when tools have been loaded via repository queries.
func ConvertDbMCPServerWithToolsToOpenapi(dbServer models.MCPServer, tools []models.MCPServerTool) *openapi.MCPServer {
	openapiTools := make([]openapi.MCPTool, 0, len(tools))
	for _, tool := range tools {
		if converted := ConvertDbMCPToolToOpenapi(tool); converted != nil {
			openapiTools = append(openapiTools, *converted)
		}
	}
	return convertDbMCPServerToOpenapiInternal(dbServer, openapiTools)
}

// convertDbMCPServerToOpenapiInternal is the shared implementation for server conversion.
// If tools is nil, toolCount is set to 0. Otherwise, it's set to len(tools).
func convertDbMCPServerToOpenapiInternal(dbServer models.MCPServer, tools []openapi.MCPTool) *openapi.MCPServer {
	attrs := dbServer.GetAttributes()
	props := dbServer.GetProperties()

	// Create property accessor for O(1) lookups
	pa := NewPropertyAccessor(props)

	// Extract base name and version
	baseName := ""
	if attrs != nil && attrs.Name != nil {
		baseName = *attrs.Name
	}
	version := pa.GetString("version")

	// Compute tool count from provided tools
	toolCount := int32(0)
	if tools != nil {
		toolCount = int32(len(tools))
	}

	openapiServer := &openapi.MCPServer{
		Name:      baseName,
		ToolCount: toolCount,
	}

	// Set tools if provided
	if len(tools) > 0 {
		openapiServer.Tools = tools
	}

	// Set core attributes
	if attrs != nil {
		openapiServer.ExternalId = attrs.ExternalID
		openapiServer.CreateTimeSinceEpoch = formatEpochToString(attrs.CreateTimeSinceEpoch)
		openapiServer.LastUpdateTimeSinceEpoch = formatEpochToString(attrs.LastUpdateTimeSinceEpoch)
	}

	// Set ID
	if dbServer.GetID() != nil {
		idStr := fmt.Sprintf("%d", *dbServer.GetID())
		openapiServer.Id = &idStr
	}

	// Extract simple string properties
	openapiServer.SourceId = pa.GetStringPtr("source_id")
	openapiServer.Provider = pa.GetStringPtr("provider")
	openapiServer.Logo = pa.GetStringPtr("logo")
	openapiServer.Version = &version
	openapiServer.License = pa.GetStringPtr("license")
	openapiServer.LicenseLink = pa.GetStringPtr("license_link")
	openapiServer.Readme = pa.GetStringPtr("readme")
	openapiServer.DeploymentMode = pa.GetStringPtr("deploymentMode")
	openapiServer.DocumentationUrl = pa.GetStringPtr("documentationUrl")
	openapiServer.RepositoryUrl = pa.GetStringPtr("repositoryUrl")
	openapiServer.SourceCode = pa.GetStringPtr("sourceCode")
	openapiServer.Description = pa.GetStringPtr("description")

	// Extract array properties
	openapiServer.Tags = pa.GetStringArray("tags")
	openapiServer.Transports = pa.GetStringArray("transports")

	// Extract time properties (parse ISO 8601 strings for content timestamps)
	openapiServer.PublishedDate = parseISO8601ToTime(pa.GetString("publishedDate"))
	openapiServer.LastUpdated = parseISO8601ToTime(pa.GetString("lastUpdated"))

	// Extract security indicators
	if pa.HasAny("verifiedSource", "secureEndpoint", "sast", "readOnlyTools") {
		openapiServer.SecurityIndicators = &openapi.MCPSecurityIndicator{
			VerifiedSource: pa.GetBoolPtr("verifiedSource"),
			SecureEndpoint: pa.GetBoolPtr("secureEndpoint"),
			Sast:           pa.GetBoolPtr("sast"),
			ReadOnlyTools:  pa.GetBoolPtr("readOnlyTools"),
		}
	}

	// Extract complex objects from JSON
	if endpointsJSON := pa.GetString("endpoints"); endpointsJSON != "" {
		var endpoints openapi.MCPEndpoints
		if err := json.Unmarshal([]byte(endpointsJSON), &endpoints); err != nil {
			glog.Warningf("Failed to unmarshal endpoints JSON: %v", err)
		} else {
			openapiServer.Endpoints = &endpoints
		}
	}
	if artifactsJSON := pa.GetString("artifacts"); artifactsJSON != "" {
		var artifacts []openapi.MCPArtifact
		if err := json.Unmarshal([]byte(artifactsJSON), &artifacts); err != nil {
			glog.Warningf("Failed to unmarshal artifacts JSON: %v", err)
		} else {
			openapiServer.Artifacts = artifacts
		}
	}
	if runtimeJSON := pa.GetString("runtimeMetadata"); runtimeJSON != "" {
		var runtimeMetadata openapi.MCPRuntimeMetadata
		if err := json.Unmarshal([]byte(runtimeJSON), &runtimeMetadata); err == nil {
			openapiServer.RuntimeMetadata = &runtimeMetadata
		}
	}

	// Convert custom properties
	if dbServer.GetCustomProperties() != nil {
		customProps := make(map[string]openapi.MetadataValue)
		for _, prop := range *dbServer.GetCustomProperties() {
			customProps[prop.Name] = convertPropertyToMetadataValue(prop)
		}
		openapiServer.CustomProperties = customProps
	}

	return openapiServer
}

// ==============================================================================
// Property Helper Functions
// ==============================================================================
//
// NOTE: Two timestamp patterns are used in this codebase:
//
// 1. SYSTEM TIMESTAMPS (CreateTimeSinceEpoch, LastUpdateTimeSinceEpoch):
//    - Stored as int64 milliseconds in entity attributes (Context/Artifact tables)
//    - Converted to numeric strings for OpenAPI (e.g., "1704067200000")
//    - Managed automatically by the system
//
// 2. CONTENT TIMESTAMPS (PublishedDate, LastUpdated):
//    - Stored as ISO 8601 strings in properties (ContextProperty/ArtifactProperty tables)
//    - Parsed to/from time.Time for Go API layer
//    - Serialized as RFC3339 in JSON API
//    - Sourced from external data (catalogs, user input)
//
// ==============================================================================

// addStringProperty adds a string property if the value is not nil.
func addStringProperty(props *[]dbmodels.Properties, name string, value *string) {
	if value != nil && *value != "" {
		*props = append(*props, dbmodels.Properties{
			Name:        name,
			StringValue: value,
		})
	}
}

// addBoolProperty adds a boolean property if the value is not nil.
func addBoolProperty(props *[]dbmodels.Properties, name string, value *bool) {
	if value != nil {
		*props = append(*props, dbmodels.Properties{
			Name:      name,
			BoolValue: value,
		})
	}
}

// addArrayProperty adds a string array property as JSON.
func addArrayProperty(props *[]dbmodels.Properties, name string, value []string) {
	if len(value) > 0 {
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			glog.Warningf("Failed to marshal array property '%s': %v", name, err)
			return
		}
		jsonStr := string(jsonBytes)
		*props = append(*props, dbmodels.Properties{
			Name:        name,
			StringValue: &jsonStr,
		})
	}
}

// addJSONProperty adds a complex object property as JSON.
func addJSONProperty(props *[]dbmodels.Properties, name string, value any) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		glog.Warningf("Failed to marshal JSON property '%s': %v", name, err)
		return
	}
	jsonStr := string(jsonBytes)
	*props = append(*props, dbmodels.Properties{
		Name:        name,
		StringValue: &jsonStr,
	})
}


// parseEpochString converts a string epoch to *int64.
// Used for SYSTEM TIMESTAMPS (CreateTimeSinceEpoch, LastUpdateTimeSinceEpoch).
func parseEpochString(epochStr *string) *int64 {
	if epochStr == nil {
		return nil
	}
	var epoch int64
	if _, err := fmt.Sscanf(*epochStr, "%d", &epoch); err == nil {
		return &epoch
	}
	glog.Warningf("Failed to parse epoch string '%s': invalid format", *epochStr)
	return nil
}

// formatEpochToString converts an int64 epoch to a numeric string.
// Used for SYSTEM TIMESTAMPS (CreateTimeSinceEpoch, LastUpdateTimeSinceEpoch).
func formatEpochToString(epoch *int64) *string {
	if epoch == nil {
		return nil
	}
	str := fmt.Sprintf("%d", *epoch)
	return &str
}

// parseISO8601ToTime parses an ISO 8601 string to time.Time.
// Used for CONTENT TIMESTAMPS (PublishedDate, LastUpdated).
// Returns nil if the string is empty or parsing fails.
func parseISO8601ToTime(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		glog.Warningf("Failed to parse ISO 8601 timestamp '%s': %v", timeStr, err)
		return nil
	}
	return &parsedTime
}

// formatTimeToISO8601 formats a time.Time to an ISO 8601 string.
// Used for CONTENT TIMESTAMPS (PublishedDate, LastUpdated).
func formatTimeToISO8601(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

// convertMetadataValueToProperty converts an OpenAPI MetadataValue to a database Property.
func convertMetadataValueToProperty(name string, value openapi.MetadataValue) dbmodels.Properties {
	prop := dbmodels.Properties{Name: name}

	if value.MetadataStringValue != nil {
		prop.StringValue = &value.MetadataStringValue.StringValue
	} else if value.MetadataIntValue != nil {
		// MetadataIntValue stores as string, parse to int32
		var intVal int32
		if _, err := fmt.Sscanf(value.MetadataIntValue.IntValue, "%d", &intVal); err != nil {
			glog.Warningf("Failed to parse int metadata value '%s': %v", value.MetadataIntValue.IntValue, err)
			// Don't set prop.IntValue - leave it nil on parse failure
		} else {
			prop.IntValue = &intVal
		}
	} else if value.MetadataDoubleValue != nil {
		prop.DoubleValue = &value.MetadataDoubleValue.DoubleValue
	} else if value.MetadataBoolValue != nil {
		prop.BoolValue = &value.MetadataBoolValue.BoolValue
	} else if value.MetadataStructValue != nil {
		prop.StringValue = &value.MetadataStructValue.StructValue
	} else if value.MetadataProtoValue != nil {
		prop.StringValue = &value.MetadataProtoValue.ProtoValue
	}

	return prop
}

// convertPropertyToMetadataValue converts a database Property to an OpenAPI MetadataValue.
func convertPropertyToMetadataValue(prop dbmodels.Properties) openapi.MetadataValue {
	metadataValue := openapi.MetadataValue{}

	if prop.StringValue != nil {
		metadataValue.MetadataStringValue = openapi.NewMetadataStringValueWithDefaults()
		metadataValue.MetadataStringValue.StringValue = *prop.StringValue
	} else if prop.IntValue != nil {
		metadataValue.MetadataIntValue = openapi.NewMetadataIntValueWithDefaults()
		metadataValue.MetadataIntValue.IntValue = fmt.Sprintf("%d", *prop.IntValue)
	} else if prop.DoubleValue != nil {
		metadataValue.MetadataDoubleValue = openapi.NewMetadataDoubleValueWithDefaults()
		metadataValue.MetadataDoubleValue.DoubleValue = *prop.DoubleValue
	} else if prop.BoolValue != nil {
		metadataValue.MetadataBoolValue = openapi.NewMetadataBoolValueWithDefaults()
		metadataValue.MetadataBoolValue.BoolValue = *prop.BoolValue
	}

	return metadataValue
}

// ConvertOpenapiMCPToolToDb converts an OpenAPI MCPTool to database MCPServerTool.
func ConvertOpenapiMCPToolToDb(openapiTool *openapi.MCPTool) models.MCPServerTool {
	dbTool := &models.MCPServerToolImpl{
		Attributes: &models.MCPServerToolAttributes{
			Name:                     &openapiTool.Name,
			CreateTimeSinceEpoch:     parseEpochString(openapiTool.CreateTimeSinceEpoch),
			LastUpdateTimeSinceEpoch: parseEpochString(openapiTool.LastUpdateTimeSinceEpoch),
		},
	}

	// Build properties
	properties := []dbmodels.Properties{}

	// Required field: accessType
	accessType := openapiTool.AccessType
	addStringProperty(&properties, "accessType", &accessType)

	// Optional field: description
	addStringProperty(&properties, "description", openapiTool.Description)

	// Optional field: parameters (JSON array)
	if len(openapiTool.Parameters) > 0 {
		addJSONProperty(&properties, "parameters", openapiTool.Parameters)
	}

	// Handle ExternalID if present (store as property since MCPServerToolAttributes doesn't have this field)
	if openapiTool.ExternalId != nil {
		addStringProperty(&properties, "externalId", openapiTool.ExternalId)
	}

	// Handle custom properties
	if customProps := openapiTool.GetCustomProperties(); len(customProps) > 0 {
		customProperties := []dbmodels.Properties{}
		for key, value := range customProps {
			customProperties = append(customProperties, convertMetadataValueToProperty(key, value))
		}
		dbTool.CustomProperties = &customProperties
	}

	dbTool.Properties = &properties
	return dbTool
}

// ConvertDbMCPToolToOpenapi converts a database MCPServerTool to OpenAPI MCPTool.
func ConvertDbMCPToolToOpenapi(dbTool models.MCPServerTool) *openapi.MCPTool {
	attr := dbTool.GetAttributes()
	if attr == nil || attr.Name == nil {
		return nil
	}

	// Create property accessor for O(1) lookups
	props := dbTool.GetProperties()
	pa := NewPropertyAccessor(props)

	// Get required accessType (required field in OpenAPI spec)
	accessType := pa.GetString("accessType")
	if accessType == "" {
		glog.Warningf("Required field 'accessType' missing for tool '%s', defaulting to 'read-only'", *attr.Name)
		accessType = "read-only" // default fallback to prevent API contract violation
	}

	// Create OpenAPI tool with required fields
	openapiTool := openapi.NewMCPTool(*attr.Name, accessType)

	// Set ID if available
	if id := dbTool.GetID(); id != nil {
		idStr := fmt.Sprintf("%d", *id)
		openapiTool.Id = &idStr
	}

	// Set timestamps (system timestamps converted to numeric strings)
	openapiTool.CreateTimeSinceEpoch = formatEpochToString(attr.CreateTimeSinceEpoch)
	openapiTool.LastUpdateTimeSinceEpoch = formatEpochToString(attr.LastUpdateTimeSinceEpoch)

	// Extract optional description
	openapiTool.Description = pa.GetStringPtr("description")

	// Extract externalId if present
	openapiTool.ExternalId = pa.GetStringPtr("externalId")

	// Extract parameters (JSON array)
	if paramsJSON := pa.GetString("parameters"); paramsJSON != "" {
		var parameters []openapi.MCPToolParameter
		if err := json.Unmarshal([]byte(paramsJSON), &parameters); err != nil {
			glog.Warningf("Failed to unmarshal tool parameters JSON: %v", err)
		} else {
			openapiTool.Parameters = parameters
		}
	}

	// Extract custom properties
	if customProps := dbTool.GetCustomProperties(); customProps != nil {
		customProperties := make(map[string]openapi.MetadataValue)
		for _, prop := range *customProps {
			customProperties[prop.Name] = convertPropertyToMetadataValue(prop)
		}
		openapiTool.CustomProperties = customProperties
	}

	return openapiTool
}
