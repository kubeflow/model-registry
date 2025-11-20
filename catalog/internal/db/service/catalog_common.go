package service

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/db/scopes"
	"gorm.io/gorm"
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

// ApplyNameOrdering applies NAME-based ordering with cursor pagination to a query.
// This handles the catalog-specific NAME ordering which requires string comparison
// in WHERE clauses (not integer casting like standard pagination).
//
// Parameters:
//   - query: The GORM query to modify
//   - tableName: The table name to use in SQL (e.g., "Context" or "Artifact")
//   - sortOrder: The sort order ("ASC" or "DESC")
//   - nextPageToken: Optional pagination token for cursor-based pagination
//   - pageSize: The page size (0 means no limit)
//
// Returns the modified query with NAME ordering and pagination applied.
func ApplyNameOrdering(query *gorm.DB, tableName string, sortOrder string, nextPageToken string, pageSize int32) *gorm.DB {
	// Normalize sort order
	order := "ASC"
	if sortOrder == "DESC" {
		order = "DESC"
	}

	// Apply name-based ordering with ID as tie-breaker
	query = query.Order(fmt.Sprintf("%s.name %s, %s.id ASC", tableName, order, tableName))

	// Handle cursor-based pagination for NAME
	if nextPageToken != "" {
		if cursor, err := scopes.DecodeCursor(nextPageToken); err == nil {
			// Cursor pagination based on name (string comparison)
			cmp := ">"
			if order == "DESC" {
				cmp = "<"
			}
			// Use proper string comparison with name and ID as tie-breaker
			query = query.Where(fmt.Sprintf("(%s.name %s ? OR (%s.name = ? AND %s.id > ?))", tableName, cmp, tableName, tableName),
				cursor.Value, cursor.Value, cursor.ID)
		}
	}

	// Apply pagination limit
	if pageSize > 0 {
		query = query.Limit(int(pageSize) + 1) // +1 to detect if there are more pages
	}

	return query
}
