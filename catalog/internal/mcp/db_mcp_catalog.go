package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	mrmodels "github.com/kubeflow/model-registry/internal/db/models"
)

// NamedQueryResolver is a function that returns named query definitions.
type NamedQueryResolver func() map[string]map[string]model.FieldFilter

// DbMcpCatalogProvider implements MCP server catalog using database storage.
type DbMcpCatalogProvider struct {
	repository         dbmodels.McpServerRepository
	namedQueryResolver NamedQueryResolver
}

// NewDbMcpCatalogProvider creates a new database-backed MCP catalog provider.
func NewDbMcpCatalogProvider(repository dbmodels.McpServerRepository) *DbMcpCatalogProvider {
	return &DbMcpCatalogProvider{
		repository: repository,
	}
}

// SetNamedQueryResolver sets the function to resolve named queries.
func (p *DbMcpCatalogProvider) SetNamedQueryResolver(resolver NamedQueryResolver) {
	p.namedQueryResolver = resolver
}

// ListMcpServers returns all MCP servers with optional filtering by name, filterQuery, and namedQuery.
func (p *DbMcpCatalogProvider) ListMcpServers(ctx context.Context, name string, filterQuery string, namedQuery string) ([]model.McpServer, error) {
	listOptions := dbmodels.McpServerListOptions{}
	if name != "" {
		listOptions.Query = &name
	}

	// Resolve named query to filter conditions
	resolvedFilterQuery := filterQuery
	if namedQuery != "" && p.namedQueryResolver != nil {
		namedQueries := p.namedQueryResolver()
		if queryFilters, exists := namedQueries[namedQuery]; exists {
			// Convert named query fields to SQL-like filter conditions
			namedFilterQuery := convertNamedQueryToFilterQuery(queryFilters)
			if namedFilterQuery != "" {
				if resolvedFilterQuery != "" {
					// Combine with existing filter query using AND
					resolvedFilterQuery = fmt.Sprintf("(%s) AND (%s)", resolvedFilterQuery, namedFilterQuery)
				} else {
					resolvedFilterQuery = namedFilterQuery
				}
			}
		} else {
			return nil, fmt.Errorf("named query %q not found", namedQuery)
		}
	}

	if resolvedFilterQuery != "" {
		// Transform license display names to SPDX identifiers before filtering
		transformedQuery := transformLicenseInFilterQuery(resolvedFilterQuery)
		listOptions.FilterQuery = &transformedQuery
	}

	result, err := p.repository.List(listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list MCP servers: %w", err)
	}

	servers := make([]model.McpServer, 0, len(result.Items))
	for _, item := range result.Items {
		server := convertDbToApiMcpServer(item)
		servers = append(servers, server)
	}

	return servers, nil
}

// GetMcpServer returns a specific MCP server by ID.
func (p *DbMcpCatalogProvider) GetMcpServer(ctx context.Context, serverId string) (*model.McpServer, error) {
	// First try to parse as integer ID
	if id, err := strconv.ParseInt(serverId, 10, 32); err == nil {
		dbServer, err := p.repository.GetByID(int32(id))
		if err != nil {
			if err.Error() == "MCP server by id not found" {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to get MCP server by ID: %w", err)
		}
		server := convertDbToApiMcpServer(dbServer)
		return &server, nil
	}

	// Otherwise try by name
	dbServer, err := p.repository.GetByName(serverId)
	if err != nil {
		if err.Error() == "MCP server by id not found" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get MCP server by name: %w", err)
	}

	server := convertDbToApiMcpServer(dbServer)
	return &server, nil
}

// convertDbToApiMcpServer converts a database MCP server to API model.
func convertDbToApiMcpServer(dbServer dbmodels.McpServer) model.McpServer {
	server := model.McpServer{}

	// Get basic attributes
	attrs := dbServer.GetAttributes()
	if attrs != nil {
		if attrs.Name != nil {
			server.Name = *attrs.Name
		}
		if attrs.CreateTimeSinceEpoch != nil {
			t := time.UnixMilli(*attrs.CreateTimeSinceEpoch)
			server.LastUpdated = &t
		}
	}

	// Use database ID as the server ID
	if id := dbServer.GetID(); id != nil {
		server.Id = strconv.Itoa(int(*id))
	}

	// Initialize customProperties map for Phase 2.11 alignment
	customProperties := make(map[string]model.MetadataValue)

	// Convert properties
	props := dbServer.GetProperties()
	if props != nil {
		for _, prop := range *props {
			convertPropertyToApiField(&server, prop)
			// Also populate customProperties for Phase 2.11 alignment
			addPropertyToCustomProperties(customProperties, prop)
		}
	}

	// Populate customProperties if any entries were added
	if len(customProperties) > 0 {
		server.CustomProperties = customProperties
	}

	// Derive transports from endpoints for remote servers if not explicitly set
	if len(server.Transports) == 0 {
		server.Transports = deriveTransportsFromEndpoints(&server)
	}

	return server
}

// deriveTransportsFromEndpoints derives transport types from the endpoints object.
// For remote servers, transports are derived from available endpoints.
// For local servers, defaults to ["stdio"].
func deriveTransportsFromEndpoints(server *model.McpServer) []model.McpTransportType {
	// For remote servers, derive from endpoints
	if server.DeploymentMode != nil && *server.DeploymentMode == model.MCPDEPLOYMENTMODE_REMOTE {
		transports := []model.McpTransportType{}
		if server.Endpoints != nil {
			if server.Endpoints.Http != nil && *server.Endpoints.Http != "" {
				transports = append(transports, model.MCPTRANSPORTTYPE_HTTP)
			}
			if server.Endpoints.Sse != nil && *server.Endpoints.Sse != "" {
				transports = append(transports, model.MCPTRANSPORTTYPE_SSE)
			}
		}
		// If no endpoints found, default to http for remote servers
		if len(transports) == 0 {
			return []model.McpTransportType{model.MCPTRANSPORTTYPE_HTTP}
		}
		return transports
	}

	// For local servers, default to stdio
	return []model.McpTransportType{model.MCPTRANSPORTTYPE_STDIO}
}

// convertPropertyToApiField maps a database property to an API server field.
func convertPropertyToApiField(server *model.McpServer, prop mrmodels.Properties) {
	switch prop.Name {
	case "source_id":
		if prop.StringValue != nil {
			server.SourceId = prop.StringValue
		}
	case "description":
		if prop.StringValue != nil {
			server.Description = prop.StringValue
		}
	case "logo":
		if prop.StringValue != nil {
			server.Logo = prop.StringValue
		}
	case "provider":
		if prop.StringValue != nil {
			server.Provider = prop.StringValue
		}
	case "version":
		if prop.StringValue != nil {
			server.Version = prop.StringValue
		}
	case "license":
		if prop.StringValue != nil {
			displayName := formatLicenseDisplayName(*prop.StringValue)
			server.License = &displayName
		}
	case "license_link":
		if prop.StringValue != nil {
			server.LicenseLink = prop.StringValue
		}
	case "transports":
		if prop.StringValue != nil {
			var transports []string
			if err := json.Unmarshal([]byte(*prop.StringValue), &transports); err == nil {
				server.Transports = convertStringsToTransports(transports)
			}
		}
	case "artifacts":
		if prop.StringValue != nil {
			var artifacts []yamlMcpArtifact
			if err := json.Unmarshal([]byte(*prop.StringValue), &artifacts); err == nil {
				server.Artifacts = make([]model.McpArtifact, 0, len(artifacts))
				for _, a := range artifacts {
					artifact := model.McpArtifact{
						Uri: a.Uri,
					}
					if a.CreateTimeSinceEpoch != "" {
						artifact.CreateTimeSinceEpoch = &a.CreateTimeSinceEpoch
					}
					if a.LastUpdateTimeSinceEpoch != "" {
						artifact.LastUpdateTimeSinceEpoch = &a.LastUpdateTimeSinceEpoch
					}
					server.Artifacts = append(server.Artifacts, artifact)
				}
			}
		}
	case "documentationUrl":
		if prop.StringValue != nil {
			server.DocumentationUrl = prop.StringValue
		}
	case "repositoryUrl":
		if prop.StringValue != nil {
			server.RepositoryUrl = prop.StringValue
		}
	case "sourceCode":
		if prop.StringValue != nil {
			server.SourceCode = prop.StringValue
		}
	case "readme":
		if prop.StringValue != nil {
			server.Readme = prop.StringValue
		}
	case "publishedDate":
		if prop.StringValue != nil {
			server.PublishedDate = prop.StringValue
		}
	case "tags":
		if prop.StringValue != nil {
			var tags []string
			if err := json.Unmarshal([]byte(*prop.StringValue), &tags); err == nil {
				server.Tags = tags
			}
		}
	case "tools":
		if prop.StringValue != nil {
			var tools []yamlMcpTool
			if err := json.Unmarshal([]byte(*prop.StringValue), &tools); err == nil {
				server.Tools = make([]model.McpTool, 0, len(tools))
				for _, t := range tools {
					server.Tools = append(server.Tools, convertYamlToolToApiTool(t))
				}
			}
		}
	case "verifiedSource":
		if prop.BoolValue != nil {
			if server.SecurityIndicators == nil {
				server.SecurityIndicators = &model.McpSecurityIndicator{}
			}
			server.SecurityIndicators.VerifiedSource = prop.BoolValue
		}
	case "secureEndpoint":
		if prop.BoolValue != nil {
			if server.SecurityIndicators == nil {
				server.SecurityIndicators = &model.McpSecurityIndicator{}
			}
			server.SecurityIndicators.SecureEndpoint = prop.BoolValue
		}
	case "sast":
		if prop.BoolValue != nil {
			if server.SecurityIndicators == nil {
				server.SecurityIndicators = &model.McpSecurityIndicator{}
			}
			server.SecurityIndicators.Sast = prop.BoolValue
		}
	case "readOnlyTools":
		if prop.BoolValue != nil {
			if server.SecurityIndicators == nil {
				server.SecurityIndicators = &model.McpSecurityIndicator{}
			}
			server.SecurityIndicators.ReadOnlyTools = prop.BoolValue
		}
	case "deploymentMode":
		if prop.StringValue != nil {
			deploymentMode := convertStringToDeploymentMode(*prop.StringValue)
			server.DeploymentMode = &deploymentMode
		}
	case "endpoints":
		if prop.StringValue != nil {
			var endpoints yamlMcpEndpoints
			if err := json.Unmarshal([]byte(*prop.StringValue), &endpoints); err == nil {
				server.Endpoints = convertYamlEndpointsToApiEndpoints(&endpoints)
			}
		}
	}
}

// addPropertyToCustomProperties adds a database property to the customProperties map.
// This function implements Phase 2.11 alignment:
// - Tags are stored as MetadataStringValue entries with empty string_value (label pattern)
// - Security indicators are stored as MetadataBoolValue entries
// - Other properties are stored with their appropriate MetadataValue types
func addPropertyToCustomProperties(customProperties map[string]model.MetadataValue, prop mrmodels.Properties) {
	metadataBool := "MetadataBoolValue"
	metadataString := "MetadataStringValue"

	switch prop.Name {
	case "tags":
		// Convert tags JSON array to individual entries with empty string_value (label pattern)
		if prop.StringValue != nil {
			var tags []string
			if err := json.Unmarshal([]byte(*prop.StringValue), &tags); err == nil {
				emptyString := ""
				for _, tag := range tags {
					customProperties[tag] = model.MetadataStringValueAsMetadataValue(&model.MetadataStringValue{
						StringValue:  emptyString,
						MetadataType: metadataString,
					})
				}
			}
		}
	case "verifiedSource", "secureEndpoint", "sast", "readOnlyTools":
		// Security indicators stored as MetadataBoolValue
		if prop.BoolValue != nil {
			customProperties[prop.Name] = model.MetadataBoolValueAsMetadataValue(&model.MetadataBoolValue{
				BoolValue:    *prop.BoolValue,
				MetadataType: metadataBool,
			})
		}
	case "provider", "license", "version", "deploymentMode":
		// String properties that might be useful for filtering
		if prop.StringValue != nil {
			customProperties[prop.Name] = model.MetadataStringValueAsMetadataValue(&model.MetadataStringValue{
				StringValue:  *prop.StringValue,
				MetadataType: metadataString,
			})
		}
	// Skip properties that are already first-class fields or internal
	case "description", "logo", "documentationUrl", "repositoryUrl", "sourceCode",
		"readme", "publishedDate", "source_id", "tools", "transports", "endpoints":
		// These are already first-class fields or handled elsewhere
		return
	}
}

// convertYamlToolToApiTool converts a YAML tool to API model.
func convertYamlToolToApiTool(t yamlMcpTool) model.McpTool {
	tool := model.McpTool{
		Name:       t.Name,
		AccessType: convertStringToAccessType(t.AccessType),
	}

	if t.Description != "" {
		tool.Description = &t.Description
	}

	if len(t.Parameters) > 0 {
		tool.Parameters = make([]model.McpToolParameter, 0, len(t.Parameters))
		for _, p := range t.Parameters {
			param := model.McpToolParameter{
				Name:     p.Name,
				Type:     p.Type,
				Required: p.Required,
			}
			if p.Description != "" {
				param.Description = &p.Description
			}
			tool.Parameters = append(tool.Parameters, param)
		}
	}

	// Set revoked status (defaults to false)
	tool.Revoked = &t.Revoked
	if t.RevokedReason != "" {
		tool.RevokedReason = &t.RevokedReason
	}

	return tool
}

// convertStringToTransport converts a string to McpTransportType.
func convertStringToTransport(transport string) model.McpTransportType {
	switch transport {
	case "http":
		return model.MCPTRANSPORTTYPE_HTTP
	case "sse":
		return model.MCPTRANSPORTTYPE_SSE
	case "stdio":
		return model.MCPTRANSPORTTYPE_STDIO
	default:
		return model.MCPTRANSPORTTYPE_STDIO
	}
}

// convertStringsToTransports converts a slice of strings to McpTransportType slice.
func convertStringsToTransports(transports []string) []model.McpTransportType {
	result := make([]model.McpTransportType, 0, len(transports))
	for _, t := range transports {
		result = append(result, convertStringToTransport(t))
	}
	return result
}

// convertStringToAccessType converts a string to McpToolAccessType.
func convertStringToAccessType(accessType string) model.McpToolAccessType {
	switch accessType {
	case "read_only":
		return model.MCPTOOLACCESSTYPE_READ_ONLY
	case "read_write":
		return model.MCPTOOLACCESSTYPE_READ_WRITE
	case "execute":
		return model.MCPTOOLACCESSTYPE_EXECUTE
	default:
		return model.MCPTOOLACCESSTYPE_READ_ONLY
	}
}

// convertStringToDeploymentMode converts a string to McpDeploymentMode.
func convertStringToDeploymentMode(mode string) model.McpDeploymentMode {
	switch mode {
	case "remote":
		return model.MCPDEPLOYMENTMODE_REMOTE
	case "local":
		return model.MCPDEPLOYMENTMODE_LOCAL
	default:
		return model.MCPDEPLOYMENTMODE_LOCAL
	}
}

// convertYamlEndpointsToApiEndpoints converts YAML endpoints to API model.
func convertYamlEndpointsToApiEndpoints(endpoints *yamlMcpEndpoints) *model.McpEndpoints {
	if endpoints == nil {
		return nil
	}
	apiEndpoints := &model.McpEndpoints{}
	if endpoints.Http != "" {
		apiEndpoints.Http = &endpoints.Http
	}
	if endpoints.Sse != "" {
		apiEndpoints.Sse = &endpoints.Sse
	}
	return apiEndpoints
}

// GetFilterOptions returns all available filter options for MCP servers.
// This includes field names and available values for each filterable field.
func (p *DbMcpCatalogProvider) GetFilterOptions(ctx context.Context) (*model.FilterOptionsList, error) {
	// Get all servers (without filtering) to extract unique values
	allServers, err := p.ListMcpServers(ctx, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to list MCP servers for filter options: %w", err)
	}

	// Use maps to track unique values for each filterable field
	providers := make(map[string]struct{})
	licenses := make(map[string]struct{})
	tags := make(map[string]struct{})
	transports := make(map[string]struct{})
	deploymentModes := make(map[string]struct{})

	for _, server := range allServers {
		// Provider
		if server.Provider != nil && *server.Provider != "" {
			providers[*server.Provider] = struct{}{}
		}

		// License (display name)
		if server.License != nil && *server.License != "" {
			licenses[*server.License] = struct{}{}
		}

		// Tags
		for _, tag := range server.Tags {
			if tag != "" {
				tags[tag] = struct{}{}
			}
		}

		// Transports
		for _, transport := range server.Transports {
			transports[string(transport)] = struct{}{}
		}

		// Deployment mode
		if server.DeploymentMode != nil {
			deploymentModes[string(*server.DeploymentMode)] = struct{}{}
		}
	}

	// Build filter options map
	options := make(map[string]model.FilterOption)

	if len(providers) > 0 {
		options["provider"] = model.FilterOption{
			Type:   "string",
			Values: mapKeysToInterface(providers),
		}
	}

	if len(licenses) > 0 {
		options["license"] = model.FilterOption{
			Type:   "string",
			Values: mapKeysToInterface(licenses),
		}
	}

	if len(tags) > 0 {
		options["tags"] = model.FilterOption{
			Type:   "string",
			Values: mapKeysToInterface(tags),
		}
	}

	if len(transports) > 0 {
		options["transports"] = model.FilterOption{
			Type:   "string",
			Values: mapKeysToInterface(transports),
		}
	}

	if len(deploymentModes) > 0 {
		options["deploymentMode"] = model.FilterOption{
			Type:   "string",
			Values: mapKeysToInterface(deploymentModes),
		}
	}

	result := &model.FilterOptionsList{
		Filters: &options,
	}

	// Include named queries if resolver is configured
	if p.namedQueryResolver != nil {
		namedQueries := p.namedQueryResolver()
		if len(namedQueries) > 0 {
			result.NamedQueries = &namedQueries
		}
	}

	return result, nil
}

// mapKeysToInterface converts a map's keys to a sorted slice of interface{}.
func mapKeysToInterface(m map[string]struct{}) []interface{} {
	// Extract keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort keys for deterministic output
	sort.Strings(keys)

	// Convert to interface slice
	result := make([]interface{}, len(keys))
	for i, k := range keys {
		result[i] = k
	}
	return result
}

// spdxToDisplayName maps SPDX license identifiers to user-friendly display names.
var spdxToDisplayName = map[string]string{
	"apache-2.0":   "Apache 2.0",
	"mit":          "MIT",
	"gpl-2.0":      "GPL 2.0",
	"gpl-3.0":      "GPL 3.0",
	"lgpl-2.1":     "LGPL 2.1",
	"lgpl-3.0":     "LGPL 3.0",
	"bsd-2-clause": "BSD 2-Clause",
	"bsd-3-clause": "BSD 3-Clause",
	"mpl-2.0":      "MPL 2.0",
	"elastic-2.0":  "Elastic 2.0",
	"cc-by-4.0":    "CC BY 4.0",
	"cc-by-sa-4.0": "CC BY-SA 4.0",
	"cc0-1.0":      "CC0 1.0",
	"unlicense":    "Unlicense",
	"isc":          "ISC",
	"wtfpl":        "WTFPL",
	"gemma":        "Gemma",
	"llama-3.1":    "Llama 3.1",
	"llama-3.3":    "Llama 3.3",
	"llama3.1":     "Llama 3.1",
	"llama3.3":     "Llama 3.3",
	"llama4":       "Llama 4",
	"modified-mit": "Modified MIT",
}

// formatLicenseDisplayName converts an SPDX license identifier to a user-friendly display name.
func formatLicenseDisplayName(spdxLicense string) string {
	if displayName, exists := spdxToDisplayName[spdxLicense]; exists {
		return displayName
	}
	// Return the original value if not found in the mapping
	return spdxLicense
}

// displayNameToSpdx is the reverse mapping from display names to SPDX identifiers.
// This is built once at init time from spdxToDisplayName.
var displayNameToSpdx map[string]string

func init() {
	displayNameToSpdx = make(map[string]string, len(spdxToDisplayName))
	for spdx, display := range spdxToDisplayName {
		displayNameToSpdx[display] = spdx
	}
}

// transformLicenseInFilterQuery transforms license display names in a filter query to SPDX identifiers.
// This handles filter queries like "license='Apache 2.0'" and converts them to "license='apache-2.0'".
func transformLicenseInFilterQuery(filterQuery string) string {
	if filterQuery == "" {
		return filterQuery
	}

	result := filterQuery

	// Replace each display name with its SPDX identifier
	for display, spdx := range displayNameToSpdx {
		// Handle single quotes: license='Apache 2.0'
		oldSingle := "'" + display + "'"
		newSingle := "'" + spdx + "'"
		result = strings.ReplaceAll(result, oldSingle, newSingle)
	}

	return result
}

// convertNamedQueryToFilterQuery converts named query field filters to a SQL-like filter query string.
// This enables named queries to work with the existing filterQuery parser in the database layer.
func convertNamedQueryToFilterQuery(fieldFilters map[string]model.FieldFilter) string {
	if len(fieldFilters) == 0 {
		return ""
	}

	var conditions []string

	for fieldName, filter := range fieldFilters {
		condition := convertFieldFilterToCondition(fieldName, filter)
		if condition != "" {
			conditions = append(conditions, condition)
		}
	}

	if len(conditions) == 0 {
		return ""
	}

	// Join all conditions with AND
	return strings.Join(conditions, " AND ")
}

// convertFieldFilterToCondition converts a single field filter to a SQL-like condition.
func convertFieldFilterToCondition(fieldName string, filter model.FieldFilter) string {
	operator := strings.ToUpper(filter.Operator)

	switch operator {
	case "=", "EQUALS":
		return formatEqualityCondition(fieldName, "=", filter.Value)
	case "!=", "<>", "NOT_EQUALS":
		return formatEqualityCondition(fieldName, "!=", filter.Value)
	case ">", ">=", "<", "<=":
		return formatComparisonCondition(fieldName, operator, filter.Value)
	case "IN", "ANYOF":
		return formatInCondition(fieldName, filter.Value)
	case "LIKE", "ILIKE":
		return formatLikeCondition(fieldName, operator, filter.Value)
	default:
		// Default to equality
		return formatEqualityCondition(fieldName, "=", filter.Value)
	}
}

// formatEqualityCondition formats an equality/inequality condition.
func formatEqualityCondition(fieldName, operator string, value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%s %s '%s'", fieldName, operator, escapeString(v))
	case bool:
		return fmt.Sprintf("%s %s %t", fieldName, operator, v)
	case float64:
		// JSON numbers come as float64
		if float64(int64(v)) == v {
			return fmt.Sprintf("%s %s %d", fieldName, operator, int64(v))
		}
		return fmt.Sprintf("%s %s %f", fieldName, operator, v)
	case int, int32, int64:
		return fmt.Sprintf("%s %s %v", fieldName, operator, v)
	default:
		return fmt.Sprintf("%s %s '%v'", fieldName, operator, v)
	}
}

// formatComparisonCondition formats a comparison condition.
func formatComparisonCondition(fieldName, operator string, value interface{}) string {
	switch v := value.(type) {
	case float64:
		if float64(int64(v)) == v {
			return fmt.Sprintf("%s %s %d", fieldName, operator, int64(v))
		}
		return fmt.Sprintf("%s %s %f", fieldName, operator, v)
	case int, int32, int64:
		return fmt.Sprintf("%s %s %v", fieldName, operator, v)
	default:
		return fmt.Sprintf("%s %s %v", fieldName, operator, v)
	}
}

// formatInCondition formats an IN condition.
func formatInCondition(fieldName string, value interface{}) string {
	switch v := value.(type) {
	case []interface{}:
		values := make([]string, 0, len(v))
		for _, item := range v {
			switch itemVal := item.(type) {
			case string:
				values = append(values, fmt.Sprintf("'%s'", escapeString(itemVal)))
			default:
				values = append(values, fmt.Sprintf("%v", itemVal))
			}
		}
		return fmt.Sprintf("%s IN (%s)", fieldName, strings.Join(values, ", "))
	case []string:
		values := make([]string, 0, len(v))
		for _, item := range v {
			values = append(values, fmt.Sprintf("'%s'", escapeString(item)))
		}
		return fmt.Sprintf("%s IN (%s)", fieldName, strings.Join(values, ", "))
	default:
		return formatEqualityCondition(fieldName, "=", value)
	}
}

// formatLikeCondition formats a LIKE/ILIKE condition.
func formatLikeCondition(fieldName, operator string, value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%s %s '%s'", fieldName, operator, escapeString(v))
	default:
		return fmt.Sprintf("%s %s '%v'", fieldName, operator, v)
	}
}

// escapeString escapes single quotes in a string for use in SQL-like conditions.
func escapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
