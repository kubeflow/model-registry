package service

import (
	"errors"
	"fmt"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"gorm.io/gorm"
)

var (
	ErrMCPServerToolNotFound  = errors.New("MCP server tool by id not found")
	ErrMCPServerToolNameEmpty = errors.New("MCP server tool name cannot be empty")
)

// MCPServerToolRepositoryImpl implements MCPServerToolRepository using GORM with Execution schema.
// Note: Uses models.Pagination for GenericRepository type parameter, though List method is overridden.
type MCPServerToolRepositoryImpl struct {
	*service.GenericRepository[models.MCPServerTool, schema.Execution, schema.ExecutionProperty, *dbmodels.Pagination]
}

// NewMCPServerToolRepository creates a new MCPServerToolRepository.
func NewMCPServerToolRepository(db *gorm.DB, typeID int32) models.MCPServerToolRepository {
	config := service.GenericRepositoryConfig[models.MCPServerTool, schema.Execution, schema.ExecutionProperty, *dbmodels.Pagination]{
		DB:                      db,
		TypeID:                  typeID,
		EntityToSchema:          mapMCPServerToolToExecution,
		SchemaToEntity:          mapDataLayerToMCPServerTool,
		EntityToProperties:      mapMCPServerToolToExecutionProperties,
		NotFoundError:           ErrMCPServerToolNotFound,
		EntityName:              "MCP server tool",
		PropertyFieldName:       "execution_id",
		ApplyListFilters:        nil, // No filters needed for simple list
		IsNewEntity:             func(entity models.MCPServerTool) bool { return entity.GetID() == nil },
		HasCustomProperties:     func(entity models.MCPServerTool) bool { return entity.GetCustomProperties() != nil },
		PreserveHistoricalTimes: true, // Preserve timestamps from YAML source data
	}

	return &MCPServerToolRepositoryImpl{
		GenericRepository: service.NewGenericRepository(config),
	}
}

// GetByID retrieves an MCP server tool by its ID.
func (r *MCPServerToolRepositoryImpl) GetByID(id int32) (models.MCPServerTool, error) {
	return r.GenericRepository.GetByID(id)
}

// List retrieves all tools for a given parent MCP server.
func (r *MCPServerToolRepositoryImpl) List(parentID int32) ([]models.MCPServerTool, error) {
	config := r.GetConfig()

	// Query executions linked to parent via Association table
	var executions []schema.Execution
	associationTable := utils.GetTableName(config.DB, &schema.Association{})
	executionTable := utils.GetTableName(config.DB, &schema.Execution{})

	err := config.DB.Table(executionTable).
		Joins(fmt.Sprintf("INNER JOIN %s ON %s.execution_id = %s.id",
			associationTable, associationTable, executionTable)).
		Where(fmt.Sprintf("%s.context_id = ? AND %s.type_id = ?",
			associationTable, executionTable), parentID, config.TypeID).
		Find(&executions).Error

	if err != nil {
		return nil, fmt.Errorf("error listing %s by parent: %w", config.EntityName, err)
	}

	// Load properties for each execution
	var tools []models.MCPServerTool
	for _, exec := range executions {
		var properties []schema.ExecutionProperty
		if err := config.DB.Where("execution_id = ?", exec.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties for %s: %w", config.EntityName, err)
		}

		tool := config.SchemaToEntity(exec, properties)
		tools = append(tools, tool)
	}

	return tools, nil
}

// Save creates or updates an MCP server tool.
func (r *MCPServerToolRepositoryImpl) Save(tool models.MCPServerTool, parentID *int32) (models.MCPServerTool, error) {
	config := r.GetConfig()
	if tool.GetTypeID() == nil {
		if config.TypeID > 0 {
			tool.SetTypeID(config.TypeID)
		}
	}

	// Validate tool name
	attrs := tool.GetAttributes()
	if attrs == nil || attrs.Name == nil || *attrs.Name == "" {
		return nil, ErrMCPServerToolNameEmpty
	}

	return r.GenericRepository.Save(tool, parentID)
}

// DeleteByParentID deletes all tools belonging to a parent MCP server.
func (r *MCPServerToolRepositoryImpl) DeleteByParentID(parentID int32) error {
	config := r.GetConfig()

	// Find all execution IDs linked to this parent context via Association
	associationTable := utils.GetTableName(config.DB, &schema.Association{})
	executionTable := utils.GetTableName(config.DB, &schema.Execution{})

	subQuery := config.DB.Table(associationTable).
		Select(associationTable+".execution_id").
		Joins(fmt.Sprintf("INNER JOIN %s ON %s.execution_id = %s.id",
			executionTable, associationTable, executionTable)).
		Where(fmt.Sprintf("%s.context_id = ? AND %s.type_id = ?",
			associationTable, executionTable), parentID, config.TypeID)

	// Delete executions matching the subquery
	result := config.DB.Where("id IN (?)", subQuery).Delete(&schema.Execution{})
	if result.Error != nil {
		return fmt.Errorf("error deleting %s by parent: %w", config.EntityName, result.Error)
	}

	return nil
}

// DeleteByID deletes an MCP server tool by its ID.
func (r *MCPServerToolRepositoryImpl) DeleteByID(id int32) error {
	config := r.GetConfig()

	result := config.DB.Where("id = ? AND type_id = ?", id, config.TypeID).Delete(&schema.Execution{})

	if result.Error != nil {
		return fmt.Errorf("error deleting %s: %w", config.EntityName, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: id %d", config.NotFoundError, id)
	}

	return nil
}

// mapMCPServerToolToExecution maps an MCPServerTool entity to Execution schema.
// Note: Tool name validation is enforced in the Save method before this function is called.
func mapMCPServerToolToExecution(tool models.MCPServerTool) schema.Execution {
	attrs := tool.GetAttributes()
	execution := schema.Execution{}

	if typeID := tool.GetTypeID(); typeID != nil {
		execution.TypeID = *typeID
	}

	if tool.GetID() != nil {
		execution.ID = *tool.GetID()
	}

	if attrs != nil && attrs.Name != nil {
		execution.Name = attrs.Name
		if attrs.CreateTimeSinceEpoch != nil {
			execution.CreateTimeSinceEpoch = *attrs.CreateTimeSinceEpoch
		}
		if attrs.LastUpdateTimeSinceEpoch != nil {
			execution.LastUpdateTimeSinceEpoch = *attrs.LastUpdateTimeSinceEpoch
		}
	}

	return execution
}

// mapMCPServerToolToExecutionProperties maps an MCPServerTool entity to ExecutionProperty schema.
func mapMCPServerToolToExecutionProperties(tool models.MCPServerTool, executionID int32) []schema.ExecutionProperty {
	var properties []schema.ExecutionProperty

	if tool.GetProperties() != nil {
		for _, prop := range *tool.GetProperties() {
			properties = append(properties, service.MapPropertiesToExecutionProperty(prop, executionID, false))
		}
	}

	if tool.GetCustomProperties() != nil {
		for _, prop := range *tool.GetCustomProperties() {
			properties = append(properties, service.MapPropertiesToExecutionProperty(prop, executionID, true))
		}
	}

	return properties
}

// mapDataLayerToMCPServerTool maps database schema to an MCPServerTool entity.
func mapDataLayerToMCPServerTool(execution schema.Execution, properties []schema.ExecutionProperty) models.MCPServerTool {
	tool := &models.MCPServerToolImpl{
		ID:     &execution.ID,
		TypeID: &execution.TypeID,
		Attributes: &models.MCPServerToolAttributes{
			Name:                     execution.Name,
			CreateTimeSinceEpoch:     &execution.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &execution.LastUpdateTimeSinceEpoch,
		},
	}

	modelProperties := []dbmodels.Properties{}
	customProperties := []dbmodels.Properties{}

	for _, prop := range properties {
		mappedProperty := service.MapExecutionPropertyToProperties(prop)

		if prop.IsCustomProperty {
			customProperties = append(customProperties, mappedProperty)
		} else {
			modelProperties = append(modelProperties, mappedProperty)
		}
	}

	tool.Properties = &modelProperties
	tool.CustomProperties = &customProperties

	return tool
}
