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
	pageSize := pagination.GetPageSize()
	orderBy := pagination.GetOrderBy()
	sortOrder := pagination.GetSortOrder()
	nextPageToken := pagination.GetNextPageToken()

	return func(db *gorm.DB) *gorm.DB {
		if pageSize != nil && *pageSize > 0 {
			db = db.Limit(int(*pageSize) + 1)
		}

		if orderBy != nil && sortOrder != nil {
			orderByStr := ""

			if *orderBy == "CREATE_TIME" {
				orderByStr = "create_time_since_epoch"
			} else if *orderBy == "LAST_UPDATE_TIME" {
				orderByStr = "last_update_time_since_epoch"
			} else {
				orderByStr = "id"
			}

			orderClause := fmt.Sprintf("%s %s", orderByStr, *sortOrder)
			db = db.Order(orderClause)
		}

		if nextPageToken != nil && *nextPageToken != "" {
			decodedCursor, err := decodeCursor(*nextPageToken)
			if err == nil {
				whereClause := buildWhereClause(decodedCursor, orderBy, sortOrder)
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

func buildWhereClause(cursor *cursor, orderBy *string, sortOrder *string) string {
	if orderBy == nil || sortOrder == nil {
		return ""
	}

	if *orderBy == "" {
		if *sortOrder == "ASC" {
			return fmt.Sprintf("id > %d", cursor.ID)
		}
		return fmt.Sprintf("id < %d", cursor.ID)
	}

	if *sortOrder == "ASC" {
		return fmt.Sprintf("(%s > '%s' OR (%s = '%s' AND id > %d))",
			*orderBy, cursor.Value, *orderBy, cursor.Value, cursor.ID)
	}
	return fmt.Sprintf("(%s < '%s' OR (%s = '%s' AND id < %d))",
		*orderBy, cursor.Value, *orderBy, cursor.Value, cursor.ID)
}

func CreateNextPageToken(id int32, value string) string {
	cursor := fmt.Sprintf("%d:%s", id, value)
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}
