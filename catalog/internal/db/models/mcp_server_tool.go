package models

import (
	"github.com/kubeflow/model-registry/internal/db/models"
)

// McpServerToolAttributes holds the attributes for an MCP server tool record.
type McpServerToolAttributes struct {
	Name                     *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

// McpServerTool represents an MCP server tool stored in the database.
type McpServerTool interface {
	models.Entity[McpServerToolAttributes]
}

// McpServerToolImpl is the concrete implementation of McpServerTool.
type McpServerToolImpl = models.BaseEntity[McpServerToolAttributes]

// McpServerToolRepository defines the interface for MCP server tool persistence.
type McpServerToolRepository interface {
	GetByID(id int32) (McpServerTool, error)
	List(parentID int32) ([]McpServerTool, error)
	Save(tool McpServerTool, parentID *int32) (McpServerTool, error)
	DeleteByParentID(parentID int32) error
	DeleteByID(id int32) error
}
