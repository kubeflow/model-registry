package service

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/db/scopes"
)

// CatalogOrderByColumns includes NAME in addition to standard columns.
// This is specific to catalog and not available in model registry.
var CatalogOrderByColumns = map[string]string{
	"ID":               "id",
	"CREATE_TIME":      "create_time_since_epoch",
	"LAST_UPDATE_TIME": "last_update_time_since_epoch",
	"NAME":             "name",
	"id":               "id", // default fallback
}

// CreateNamePaginationToken creates a pagination token for NAME ordering.
// If name is nil, it falls back to using the entity ID.
func CreateNamePaginationToken(entityID int32, name *string) string {
	if name != nil {
		return scopes.CreateNextPageToken(entityID, name)
	}
	// Fallback to ID if name is nil
	return scopes.CreateNextPageToken(entityID, fmt.Sprintf("%d", entityID))
}
