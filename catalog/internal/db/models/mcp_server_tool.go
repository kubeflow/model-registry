package models

import (
	"github.com/kubeflow/model-registry/internal/db/models"
)

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
	List(parentID int32) ([]MCPServerTool, error)
	Save(tool MCPServerTool, parentID *int32) (MCPServerTool, error)
	DeleteByParentID(parentID int32) error
	DeleteByID(id int32) error
}
