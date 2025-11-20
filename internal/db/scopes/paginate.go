package scopes

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kubeflow/model-registry/internal/db/dbutil"
	"github.com/kubeflow/model-registry/internal/db/models"
	"gorm.io/gorm"
)

// Allowed column names for orderBy to prevent SQL injection
var allowedOrderByColumns = map[string]string{
	"ID":               "id",
	"CREATE_TIME":      "create_time_since_epoch",
	"LAST_UPDATE_TIME": "last_update_time_since_epoch",
	"id":               "id", // default fallback
}

// Allowed sort orders to prevent SQL injection
var allowedSortOrders = map[string]string{
	"ASC":  models.SortOrderAsc,
	"DESC": models.SortOrderDesc,
}

// isValidTablePrefix validates table prefix to prevent SQL injection
func isValidTablePrefix(tablePrefix string) bool {
	if tablePrefix == "" {
		return true
	}
	// Only allow alphanumeric characters and underscores, common table naming convention
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, tablePrefix)
	return matched
}

func Paginate(value any, pagination *models.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return PaginateWithTablePrefix(value, pagination, db, "")
}

func PaginateWithTablePrefix(value any, pagination *models.Pagination, db *gorm.DB, tablePrefix string) func(db *gorm.DB) *gorm.DB {
	return PaginateWithOptions(value, pagination, db, tablePrefix, nil)
}

// PaginateWithOptions provides full control over pagination with custom allowed columns
func PaginateWithOptions(value any, pagination *models.Pagination, db *gorm.DB, tablePrefix string, customAllowedColumns map[string]string) func(db *gorm.DB) *gorm.DB {
	pageSize := pagination.GetPageSize()
	orderBy := pagination.GetOrderBy()
	sortOrder := pagination.GetSortOrder()
	nextPageToken := pagination.GetNextPageToken()

	// Use custom allowed columns if provided, otherwise use default
	columnsMap := allowedOrderByColumns
	if customAllowedColumns != nil {
		columnsMap = customAllowedColumns
	}

	return func(db *gorm.DB) *gorm.DB {
		if pageSize > 0 {
			db = db.Limit(int(pageSize) + 1)
		}

		if orderBy != "" && sortOrder != "" {
			// Validate and sanitize orderBy
			sanitizedOrderBy, ok := columnsMap[orderBy]
			if !ok {
				sanitizedOrderBy = models.DefaultOrderBy
			}

			sanitizedSortOrder := models.DefaultSortOrder

			// Validate and sanitize so
			if so, ok := allowedSortOrders[sortOrder]; ok {
				sanitizedSortOrder = so
			}

			// Apply table prefix to ORDER BY column for consistency with WHERE clauses
			orderByColumn := sanitizedOrderBy
			if tablePrefix != "" {
				quotedPrefix := dbutil.QuoteTableName(db, tablePrefix)
				orderByColumn = quotedPrefix + "." + sanitizedOrderBy
			}

			db = db.Order(fmt.Sprintf("%s %s", orderByColumn, sanitizedSortOrder))
		}

		if nextPageToken != "" {
			decodedCursor, err := DecodeCursor(nextPageToken)
			if err == nil {
				db = buildWhereClause(db, decodedCursor, orderBy, sortOrder, tablePrefix)
			}
		}

		return db
	}
}

// buildWhereClause now returns a *gorm.DB with properly parameterized queries instead of raw SQL strings
func buildWhereClause(db *gorm.DB, cursor *Cursor, orderBy string, sortOrder string, tablePrefix string) *gorm.DB {
	// Validate table prefix to prevent SQL injection
	if !isValidTablePrefix(tablePrefix) {
		// If invalid table prefix, ignore it and use no prefix
		tablePrefix = ""
	}

	// Apply database-specific quoting to table prefix
	if tablePrefix != "" {
		tablePrefix = dbutil.QuoteTableName(db, tablePrefix)
	}

	// Build column names with proper validation
	idColumn := "id"
	if tablePrefix != "" {
		idColumn = tablePrefix + ".id"
	}

	// Validate and get the actual column name for orderBy
	orderByColumn, ok := allowedOrderByColumns[orderBy]
	if !ok {
		orderByColumn = models.DefaultOrderBy
	}

	if tablePrefix != "" && orderByColumn != "" {
		orderByColumn = tablePrefix + "." + orderByColumn
	}

	// Validate sort order
	sanitizedSortOrder := models.DefaultSortOrder
	if so, ok := allowedSortOrders[sortOrder]; ok {
		sanitizedSortOrder = so
	}

	if orderBy == "" {
		if sanitizedSortOrder == "ASC" {
			return db.Where(idColumn+" > ?", cursor.ID)
		}
		return db.Where(idColumn+" < ?", cursor.ID)
	}

	if sanitizedSortOrder == "ASC" {
		return db.Where("("+orderByColumn+" > ? OR ("+orderByColumn+" = ? AND "+idColumn+" > ?))",
			cursor.Value, cursor.Value, cursor.ID)
	}
	return db.Where("("+orderByColumn+" < ? OR ("+orderByColumn+" = ? AND "+idColumn+" < ?))",
		cursor.Value, cursor.Value, cursor.ID)
}

type Cursor struct {
	ID    int32
	Value string
}

// DecodeCursor parses a next page token.
func DecodeCursor(token string) (*Cursor, error) {
	// Sanity check the token size
	if len(token) > 1024 {
		return nil, fmt.Errorf("invalid cursor format")
	}

	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cursor format")
	}

	id, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return nil, err
	}

	return &Cursor{
		ID:    int32(id),
		Value: parts[1],
	}, nil
}

func CreateNextPageToken(id int32, value any) string {
	var valueString string

	switch tval := value.(type) {
	case string:
		valueString = tval
	case *string:
		if tval != nil {
			valueString = *tval
		}
	case *float64:
		if tval != nil {
			valueString = fmt.Sprintf("%.15f", *tval)
		}
	case *int64:
		if tval != nil {
			valueString = fmt.Sprintf("%d", *tval)
		}
	default:
		valueString = fmt.Sprintf("%v", value)
	}
	cursor := fmt.Sprintf("%d:%s", id, valueString)
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}
