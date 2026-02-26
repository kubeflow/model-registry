package models

import (
	catalogfilter "github.com/kubeflow/model-registry/catalog/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/models"
)

// MCPServerToolListOptions holds the options for listing MCP server tools.
type MCPServerToolListOptions struct {
	models.Pagination
	ParentID int32
}

// GetRestEntityType implements the FilterApplier interface.
func (o *MCPServerToolListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityType(catalogfilter.RestEntityMCPServerTool)
}

// MCPServerToolAttributes holds the attributes for an MCP server tool record.
type MCPServerToolAttributes struct {
	Name                     *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

// MCPServerTool represents an MCP server tool stored in the database.
type MCPServerTool interface {
	models.Entity[MCPServerToolAttributes]
}

// MCPServerToolImpl is the concrete implementation of MCPServerTool.
type MCPServerToolImpl = models.BaseEntity[MCPServerToolAttributes]

// MCPServerToolRepository defines the interface for MCP server tool persistence.
type MCPServerToolRepository interface {
	GetByID(id int32) (MCPServerTool, error)
	List(listOptions MCPServerToolListOptions) ([]MCPServerTool, error)
	Save(tool MCPServerTool, parentID *int32) (MCPServerTool, error)
	DeleteByParentID(parentID int32) error
	DeleteByID(id int32) error
}
