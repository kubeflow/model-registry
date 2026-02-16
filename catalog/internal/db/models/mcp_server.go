package models

import (
	catalogfilter "github.com/kubeflow/model-registry/catalog/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/models"
)

// MCPServerListOptions holds the options for listing MCP servers.
type MCPServerListOptions struct {
	models.Pagination
	Name        *string
	SourceIDs   *[]string
	Query       *string
	FilterQuery *string
	NamedQuery  *string
}

// GetRestEntityType implements the FilterApplier interface.
func (c *MCPServerListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityType(catalogfilter.RestEntityMCPServer)
}

// GetFilterQuery returns the filter query string for advanced filtering.
func (c *MCPServerListOptions) GetFilterQuery() string {
	if c.FilterQuery == nil {
		return ""
	}
	return *c.FilterQuery
}

// MCPServerAttributes holds the attributes for an MCP server record.
type MCPServerAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

// MCPServer represents an MCP server stored in the database.
type MCPServer interface {
	models.Entity[MCPServerAttributes]
}

// MCPServerImpl is the concrete implementation of MCPServer.
type MCPServerImpl = models.BaseEntity[MCPServerAttributes]

// MCPServerRepository defines the interface for MCP server persistence.
type MCPServerRepository interface {
	GetByID(id int32) (MCPServer, error)
	GetByNameAndVersion(name string, version string) (MCPServer, error)
	List(listOptions MCPServerListOptions) (*models.ListWrapper[MCPServer], error)
	Save(server MCPServer) (MCPServer, error)
	DeleteBySource(sourceID string) error
	DeleteByID(id int32) error
	GetDistinctSourceIDs() ([]string, error)
}