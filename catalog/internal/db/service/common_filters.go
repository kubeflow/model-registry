package service

import (
	"fmt"
	"strings"

	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/utils"
	"gorm.io/gorm"
)

// ApplySourceIDFilter adds a WHERE clause filtering Context entities by source_id property.
// Reusable by any Context-based catalog entity (models, MCP servers, future entity types).
func ApplySourceIDFilter(query *gorm.DB, sourceIDs []string, entityTable string) *gorm.DB {
	// Filter out empty strings from SourceIDs, for some reason it's passed if no sources are specified
	var nonEmptySourceIDs []string
	for _, sourceID := range sourceIDs {
		if sourceID != "" {
			nonEmptySourceIDs = append(nonEmptySourceIDs, sourceID)
		}
	}

	if len(nonEmptySourceIDs) > 0 {
		propertyTable := utils.GetTableName(query.Statement.DB, &schema.ContextProperty{})

		joinClause := fmt.Sprintf("JOIN %s cp ON cp.context_id = %s.id", propertyTable, entityTable)
		query = query.Joins(joinClause).
			Where("cp.name = ? AND cp.string_value IN ?", "source_id", nonEmptySourceIDs)
	}

	return query
}

// ApplyTextQueryFilter adds a WHERE clause searching entity name and configurable properties.
// searchableProperties lists property-table property names to search (e.g., "description", "provider").
func ApplyTextQueryFilter(query *gorm.DB, searchTerm string, entityTable string, searchableProperties []string) *gorm.DB {
	if searchTerm == "" {
		return query
	}

	queryPattern := fmt.Sprintf("%%%s%%", strings.ToLower(searchTerm))
	propertyTable := utils.GetTableName(query.Statement.DB, &schema.ContextProperty{})

	// Search in name (context table)
	nameCondition := fmt.Sprintf("LOWER(%s.name) LIKE ?", entityTable)

	conditions := []string{nameCondition}
	args := []any{queryPattern}

	// Search in configurable properties
	if len(searchableProperties) > 0 {
		// Build placeholders for property names
		placeholders := strings.Repeat("?,", len(searchableProperties))
		placeholders = placeholders[:len(placeholders)-1] // remove last comma

		propertyCondition := fmt.Sprintf("EXISTS (SELECT 1 FROM %s cp WHERE cp.context_id = %s.id AND cp.name IN (%s) AND LOWER(cp.string_value) LIKE ?)",
			propertyTable, entityTable, placeholders)

		conditions = append(conditions, propertyCondition)

		// Add property names as arguments
		for _, prop := range searchableProperties {
			args = append(args, prop)
		}
		// Add query pattern for property values
		args = append(args, queryPattern)
	}

	// Combine all conditions with OR
	whereClause := fmt.Sprintf("(%s)", strings.Join(conditions, " OR "))
	query = query.Where(whereClause, args...)

	return query
}
