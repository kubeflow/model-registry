package pagination

import (
	"fmt"

	"github.com/kubeflow/hub/internal/db/scopes"
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
		return scopes.CreateNextPageToken(entityID, *name)
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
//   - stripSourcePrefix: When true, strips the "sourceId:" prefix before sorting.
//     Only models use prefixed names; artifacts and MCP servers do not.
//
// Returns the modified query with NAME ordering and pagination applied.
func ApplyNameOrdering(query *gorm.DB, tableName string, sortOrder string, nextPageToken string, pageSize int32, stripSourcePrefix bool) *gorm.DB {
	// Normalize sort order
	order := "ASC"
	if sortOrder == "DESC" {
		order = "DESC"
	}

	var nameExpr string
	if stripSourcePrefix {
		nameExpr = fmt.Sprintf("COALESCE(NULLIF(SUBSTRING(%s.name FROM STRPOS(%s.name, ':') + 1), ''), %s.name)", tableName, tableName, tableName)
	} else {
		nameExpr = fmt.Sprintf("%s.name", tableName)
	}
	query = query.Order(fmt.Sprintf("%s %s, %s.id ASC", nameExpr, order, tableName))

	// Handle cursor-based pagination for NAME
	if nextPageToken != "" {
		cursor, err := scopes.DecodeCursor(nextPageToken)
		if err != nil {
			_ = query.AddError(fmt.Errorf("invalid nextPageToken: %w", err))
			return query
		}
		// Cursor pagination based on display name (string comparison).
		// cursor.Value holds the display name (after the colon prefix).
		cmp := ">"
		if order == "DESC" {
			cmp = "<"
		}
		// Use proper string comparison with display name and ID as tie-breaker
		query = query.Where(
			fmt.Sprintf("(%s %s ? OR (%s = ? AND %s.id > ?))", nameExpr, cmp, nameExpr, tableName),
			cursor.Value, cursor.Value, cursor.ID)
	}

	// Apply pagination limit
	if pageSize > 0 {
		query = query.Limit(int(pageSize) + 1) // +1 to detect if there are more pages
	}

	return query
}
