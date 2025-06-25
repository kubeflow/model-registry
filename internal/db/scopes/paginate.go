package scopes

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/kubeflow/model-registry/internal/db/models"
	"gorm.io/gorm"
)

type cursor struct {
	ID    int32
	Value string
}

func Paginate(value any, pagination *models.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return PaginateWithTablePrefix(value, pagination, db, "")
}

func PaginateWithTablePrefix(value any, pagination *models.Pagination, db *gorm.DB, tablePrefix string) func(db *gorm.DB) *gorm.DB {
	pageSize := pagination.GetPageSize()
	orderBy := pagination.GetOrderBy()
	sortOrder := pagination.GetSortOrder()
	nextPageToken := pagination.GetNextPageToken()

	return func(db *gorm.DB) *gorm.DB {
		if pageSize > 0 {
			db = db.Limit(int(pageSize) + 1)
		}

		if orderBy != "" && sortOrder != "" {
			orderByStr := ""
			sortOrderStr := ""

			switch orderBy {
			case "CREATE_TIME":
				orderByStr = "create_time_since_epoch"
			case "LAST_UPDATE_TIME":
				orderByStr = "last_update_time_since_epoch"
			default:
				orderByStr = models.DefaultOrderBy
			}

			switch sortOrder {
			case "ASC":
				sortOrderStr = "ASC"
			case "DESC":
				sortOrderStr = "DESC"
			default:
				sortOrderStr = models.DefaultSortOrder
			}

			orderClause := fmt.Sprintf("%s %s", orderByStr, sortOrderStr)
			db = db.Order(orderClause)
		}

		if nextPageToken != "" {
			decodedCursor, err := decodeCursor(nextPageToken)
			if err == nil {
				whereClause := buildWhereClause(decodedCursor, orderBy, sortOrder, tablePrefix)
				if whereClause != "" {
					db = db.Where(whereClause)
				}
			}
		}

		return db
	}
}

func decodeCursor(token string) (*cursor, error) {
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

	return &cursor{
		ID:    int32(id),
		Value: parts[1],
	}, nil
}

func buildWhereClause(cursor *cursor, orderBy string, sortOrder string, tablePrefix string) string {
	// Add table prefix to column names if provided
	idColumn := "id"
	orderByColumn := orderBy
	if tablePrefix != "" {
		idColumn = tablePrefix + ".id"
		if orderBy != "" {
			orderByColumn = tablePrefix + "." + orderBy
		}
	}

	if orderBy == "" {
		if sortOrder == "ASC" {
			return fmt.Sprintf("%s > %d", idColumn, cursor.ID)
		}
		return fmt.Sprintf("%s < %d", idColumn, cursor.ID)
	}

	if sortOrder == "ASC" {
		return fmt.Sprintf("(%s > '%s' OR (%s = '%s' AND %s > %d))",
			orderByColumn, cursor.Value, orderByColumn, cursor.Value, idColumn, cursor.ID)
	}
	return fmt.Sprintf("(%s < '%s' OR (%s = '%s' AND %s < %d))",
		orderByColumn, cursor.Value, orderByColumn, cursor.Value, idColumn, cursor.ID)
}

func CreateNextPageToken(id int32, value string) string {
	cursor := fmt.Sprintf("%d:%s", id, value)
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}
