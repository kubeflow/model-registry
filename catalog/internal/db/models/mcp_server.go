package models

import (
	catalogfilter "github.com/kubeflow/model-registry/catalog/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/models"
)

// McpServerListOptions holds the options for listing MCP servers.
type McpServerListOptions struct {
	models.Pagination
	Name        *string
	SourceIDs   *[]string
	Query       *string
	FilterQuery *string
	NamedQuery  *string
}

// GetRestEntityType implements the FilterApplier interface.
func (c *McpServerListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityType(catalogfilter.RestEntityMcpServer)
}

// GetFilterQuery returns the filter query string for advanced filtering.
func (c *McpServerListOptions) GetFilterQuery() string {
	if c.FilterQuery == nil {
		return ""
	}
	return *c.FilterQuery
}

// McpServerAttributes holds the attributes for an MCP server record.
type McpServerAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

// McpServer represents an MCP server stored in the database.
type McpServer interface {
	models.Entity[McpServerAttributes]
}

// McpServerImpl is the concrete implementation of McpServer.
type McpServerImpl = models.BaseEntity[McpServerAttributes]

// McpServerRepository defines the interface for MCP server persistence.
type McpServerRepository interface {
	GetByID(id int32) (McpServer, error)
	GetByNameAndVersion(name string, version string) (McpServer, error)
	List(listOptions McpServerListOptions) (*models.ListWrapper[McpServer], error)
	Save(server McpServer) (McpServer, error)
	DeleteBySource(sourceID string) error
	DeleteByID(id int32) error
	GetDistinctSourceIDs() ([]string, error)
}