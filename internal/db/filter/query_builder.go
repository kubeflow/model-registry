package filter

import (
	"fmt"
	"strings"

	"github.com/kubeflow/model-registry/internal/db/constants"
	"gorm.io/gorm"
)

// EntityType represents the type of entity for proper query building
type EntityType string

const (
	EntityTypeContext   EntityType = "context"
	EntityTypeArtifact  EntityType = "artifact"
	EntityTypeExecution EntityType = "execution"
)

// QueryBuilder builds GORM queries from filter expressions
// It handles special cases like prefixed names for child entities (e.g., ModelVersion, ExperimentRun)
// where names are stored as "parentId:actualName" in the database
type QueryBuilder struct {
	entityType     EntityType
	restEntityType RestEntityType
	tablePrefix    string
	joinCounter    int
	db             *gorm.DB // Added to access naming strategy
}

// NewQueryBuilderForRestEntity creates a new query builder for the specified REST entity type
func NewQueryBuilderForRestEntity(restEntityType RestEntityType) *QueryBuilder {
	// Get the underlying MLMD entity type
	entityType := GetMLMDEntityType(restEntityType)

	var tablePrefix string
	switch entityType {
	case EntityTypeContext:
		tablePrefix = "Context"
	case EntityTypeArtifact:
		tablePrefix = "Artifact"
	case EntityTypeExecution:
		tablePrefix = "Execution"
	default:
		tablePrefix = "Context" // default
	}

	return &QueryBuilder{
		entityType:     entityType,
		restEntityType: restEntityType,
		tablePrefix:    tablePrefix,
		joinCounter:    0,
	}
}

// BuildQuery builds a GORM query from a filter expression
func (qb *QueryBuilder) BuildQuery(db *gorm.DB, expr *FilterExpression) *gorm.DB {
	if expr == nil {
		return db
	}

	// Store db reference for table name quoting
	qb.db = db
	qb.applyDatabaseQuoting()

	return qb.buildExpression(db, expr)
}

// applyDatabaseQuoting updates tablePrefix with proper quoting based on database dialect
func (qb *QueryBuilder) applyDatabaseQuoting() {
	if qb.db == nil {
		return
	}
	switch qb.db.Name() {
	case "mysql":
		qb.tablePrefix = "`" + qb.tablePrefix + "`"
	case "postgres":
		qb.tablePrefix = `"` + qb.tablePrefix + `"`
	default:
		// Keep unquoted for other databases
	}
}

// quoteTableName quotes a table name based on database dialect
func (qb *QueryBuilder) quoteTableName(tableName string) string {
	if qb.db == nil {
		return tableName
	}
	switch qb.db.Name() {
	case "mysql":
		return "`" + tableName + "`"
	case "postgres":
		return `"` + tableName + `"`
	default:
		return tableName
	}
}

// buildExpression recursively builds query conditions from filter expressions
func (qb *QueryBuilder) buildExpression(db *gorm.DB, expr *FilterExpression) *gorm.DB {
	if expr.IsLeaf {
		return qb.buildLeafExpression(db, expr)
	}

	// Handle logical operators (AND, OR)
	switch expr.Operator {
	case "AND":
		leftQuery := qb.buildExpression(db, expr.Left)
		return qb.buildExpression(leftQuery, expr.Right)

	case "OR":
		// For OR conditions, we need to group them properly
		leftCondition := qb.buildConditionString(expr.Left)
		rightCondition := qb.buildConditionString(expr.Right)

		condition := fmt.Sprintf("(%s OR %s)", leftCondition.condition, rightCondition.condition)
		args := append(leftCondition.args, rightCondition.args...)

		return db.Where(condition, args...)

	default:
		return db
	}
}

// conditionResult holds a condition string and its arguments
type conditionResult struct {
	condition string
	args      []any
}

// buildConditionString builds a condition string from an expression (for OR grouping)
func (qb *QueryBuilder) buildConditionString(expr *FilterExpression) conditionResult {
	if expr.IsLeaf {
		return qb.buildLeafConditionString(expr)
	}

	switch expr.Operator {
	case "AND":
		left := qb.buildConditionString(expr.Left)
		right := qb.buildConditionString(expr.Right)

		condition := fmt.Sprintf("(%s AND %s)", left.condition, right.condition)
		args := append(left.args, right.args...)

		return conditionResult{condition: condition, args: args}

	case "OR":
		left := qb.buildConditionString(expr.Left)
		right := qb.buildConditionString(expr.Right)

		condition := fmt.Sprintf("(%s OR %s)", left.condition, right.condition)
		args := append(left.args, right.args...)

		return conditionResult{condition: condition, args: args}
	}

	return conditionResult{condition: "1=1", args: []any{}}
}

// buildPropertyReference creates a property reference from a filter expression
func (qb *QueryBuilder) buildPropertyReference(expr *FilterExpression) *PropertyReference {
	var propDef PropertyDefinition
	propertyName := expr.Property

	// Check if the property has an explicit type suffix (e.g., "budget.double_value")
	var explicitType string
	if parts := strings.Split(propertyName, "."); len(parts) == 2 {
		propertyName = parts[0]
		explicitType = parts[1]
	}

	// Use REST entity type-aware property mapping if available
	if qb.restEntityType != "" {
		propDef = GetPropertyDefinitionForRestEntity(qb.restEntityType, propertyName)
	} else {
		// Fallback to MLMD entity type only
		propDef = GetPropertyDefinition(qb.entityType, propertyName)
	}

	// For property table properties, use the Column field as the database property name
	propName := propertyName
	if propDef.Location == PropertyTable && propDef.Column != "" {
		propName = propDef.Column
	}

	propRef := &PropertyReference{
		Name:      propName,
		IsCustom:  propDef.Location == Custom,
		ValueType: propDef.ValueType,
	}

	// If explicit type was specified, use it
	if explicitType != "" {
		propRef.ValueType = explicitType
	} else if propRef.IsCustom {
		// For custom properties without explicit type, infer from value
		propRef.ValueType = qb.inferValueTypeFromInterface(expr.Value)
	}

	return propRef
}

// buildLeafExpression builds a GORM query for a leaf expression (property comparison)
func (qb *QueryBuilder) buildLeafExpression(db *gorm.DB, expr *FilterExpression) *gorm.DB {
	propRef := qb.buildPropertyReference(expr)
	return qb.buildPropertyCondition(db, propRef, expr.Operator, expr.Value)
}

// buildLeafConditionString builds a condition string for a leaf expression
func (qb *QueryBuilder) buildLeafConditionString(expr *FilterExpression) conditionResult {
	propRef := qb.buildPropertyReference(expr)
	return qb.buildPropertyConditionString(propRef, expr.Operator, expr.Value)
}

// inferValueTypeFromInterface infers the value type from an any value
func (qb *QueryBuilder) inferValueTypeFromInterface(value any) string {
	switch value.(type) {
	case string:
		return StringValueType
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return IntValueType
	case float32, float64:
		return DoubleValueType
	case bool:
		return BoolValueType
	default:
		return StringValueType // fallback
	}
}

// buildPropertyCondition builds a GORM query condition for a property
func (qb *QueryBuilder) buildPropertyCondition(db *gorm.DB, propRef *PropertyReference, operator string, value any) *gorm.DB {
	propDef := GetPropertyDefinition(qb.entityType, propRef.Name)

	switch propDef.Location {
	case EntityTable:
		return qb.buildEntityTablePropertyCondition(db, propRef, operator, value)
	case PropertyTable, Custom:
		return qb.buildPropertyTableCondition(db, propRef, operator, value)
	default:
		return qb.buildEntityTablePropertyCondition(db, propRef, operator, value)
	}
}

// buildPropertyConditionString builds a condition string for a property
func (qb *QueryBuilder) buildPropertyConditionString(propRef *PropertyReference, operator string, value any) conditionResult {
	propDef := GetPropertyDefinition(qb.entityType, propRef.Name)

	switch propDef.Location {
	case EntityTable:
		return qb.buildEntityTablePropertyConditionString(propRef, operator, value)
	case PropertyTable, Custom:
		return qb.buildPropertyTableConditionString(propRef, operator, value)
	default:
		return qb.buildEntityTablePropertyConditionString(propRef, operator, value)
	}
}

// ConvertStateValue converts string state values to integers based on entity type
func (qb *QueryBuilder) ConvertStateValue(propertyName string, value any) any {
	// Only convert for state properties
	if propertyName == "state" {
		if strValue, ok := value.(string); ok {
			switch qb.entityType {
			case EntityTypeArtifact:
				if intValue, exists := constants.ArtifactStateMapping[strValue]; exists {
					return int32(intValue)
				}
				// Invalid artifact state - return value that matches no records
				return int32(-1) // No artifact has state=-1, so this returns empty results
			case EntityTypeExecution:
				if intValue, exists := constants.ExecutionStateMapping[strValue]; exists {
					return int32(intValue)
				}
				// Invalid execution state - return value that matches no records
				return int32(-1) // No execution has state=-1, so this returns empty results
			case EntityTypeContext:
				// Context entities (RegisteredModel, ModelVersion, etc.) use string states
				// These are stored as string properties, so no conversion needed
				return value
			}
		}
		// If conversion fails or value is not a string, return original value
	}
	return value
}

// buildEntityTablePropertyCondition builds a condition for properties stored in the entity table
func (qb *QueryBuilder) buildEntityTablePropertyCondition(db *gorm.DB, propRef *PropertyReference, operator string, value any) *gorm.DB {
	propDef := GetPropertyDefinition(qb.entityType, propRef.Name)
	column := fmt.Sprintf("%s.%s", qb.tablePrefix, propDef.Column)

	// Convert state string values to integers based on entity type
	value = qb.ConvertStateValue(propRef.Name, value)

	// Handle prefixed names for child entities
	if qb.restEntityType != "" && propRef.Name == "name" && isChildEntity(qb.restEntityType) {
		if strValue, ok := value.(string); ok {
			// For exact match, convert to LIKE pattern with prefix
			if operator == "=" {
				operator = "LIKE"
				value = "%:" + strValue
			} else if operator == "LIKE" && !strings.Contains(strValue, ":") {
				// For LIKE patterns without ':', add prefix handling
				if !strings.HasPrefix(strValue, "%") {
					// Pattern like 'pattern%' -> needs prefix wildcard -> '%:pattern%'
					value = "%:" + strValue
				}
				// Pattern like '%something' or '%-beta' -> keep as is
				// because names are stored as 'parentId:actualName' and '%' will match 'parentId:'
			}
			// If pattern already contains ':', assume it's already properly formatted
		}
	}

	// Use cross-database case-insensitive LIKE for ILIKE operator
	if operator == "ILIKE" {
		return qb.buildCaseInsensitiveLikeCondition(db, column, value)
	}

	condition := qb.buildOperatorCondition(column, operator, value)
	return db.Where(condition.condition, condition.args...)
}

// buildEntityTablePropertyConditionString builds a condition string for properties stored in the entity table
func (qb *QueryBuilder) buildEntityTablePropertyConditionString(propRef *PropertyReference, operator string, value any) conditionResult {
	propDef := GetPropertyDefinition(qb.entityType, propRef.Name)
	column := fmt.Sprintf("%s.%s", qb.tablePrefix, propDef.Column)

	// Convert state string values to integers based on entity type
	value = qb.ConvertStateValue(propRef.Name, value)

	// Handle prefixed names for child entities
	if qb.restEntityType != "" && propRef.Name == "name" && isChildEntity(qb.restEntityType) {
		if strValue, ok := value.(string); ok {
			// For exact match, convert to LIKE pattern with prefix
			if operator == "=" {
				operator = "LIKE"
				value = "%:" + strValue
			} else if operator == "LIKE" && !strings.Contains(strValue, ":") {
				// For LIKE patterns without ':', add prefix handling
				if !strings.HasPrefix(strValue, "%") {
					// Pattern like 'pattern%' -> needs prefix wildcard -> '%:pattern%'
					value = "%:" + strValue
				}
				// Pattern like '%something' or '%-beta' -> keep as is
				// because names are stored as 'parentId:actualName' and '%' will match 'parentId:'
			}
			// If pattern already contains ':', assume it's already properly formatted
		}
	}

	return qb.buildOperatorCondition(column, operator, value)
}

// buildPropertyTableCondition builds a condition for properties stored in the property table (requires join)
func (qb *QueryBuilder) buildPropertyTableCondition(db *gorm.DB, propRef *PropertyReference, operator string, value any) *gorm.DB {
	qb.joinCounter++
	alias := fmt.Sprintf("prop_%d", qb.joinCounter)

	// Determine the property table based on entity type
	var propertyTable string
	var joinCondition string

	switch qb.entityType {
	case EntityTypeContext:
		propertyTable = qb.quoteTableName("ContextProperty")
		joinCondition = fmt.Sprintf("%s.context_id = %s.id", alias, qb.tablePrefix)
	case EntityTypeArtifact:
		propertyTable = qb.quoteTableName("ArtifactProperty")
		joinCondition = fmt.Sprintf("%s.artifact_id = %s.id", alias, qb.tablePrefix)
	case EntityTypeExecution:
		propertyTable = qb.quoteTableName("ExecutionProperty")
		joinCondition = fmt.Sprintf("%s.execution_id = %s.id", alias, qb.tablePrefix)
	}

	// Join the property table
	joinClause := fmt.Sprintf("JOIN %s %s ON %s", propertyTable, alias, joinCondition)
	db = db.Joins(joinClause)

	// Add conditions for property name
	db = db.Where(fmt.Sprintf("%s.name = ?", alias), propRef.Name)

	// Use the specific value type column based on inferred type
	valueColumn := fmt.Sprintf("%s.%s", alias, propRef.ValueType)

	// Use cross-database case-insensitive LIKE for ILIKE operator
	if operator == "ILIKE" {
		return qb.buildCaseInsensitiveLikeCondition(db, valueColumn, value)
	}

	condition := qb.buildOperatorCondition(valueColumn, operator, value)
	return db.Where(condition.condition, condition.args...)
}

// buildPropertyTableConditionString builds a condition string for properties stored in the property table
func (qb *QueryBuilder) buildPropertyTableConditionString(propRef *PropertyReference, operator string, value any) conditionResult {
	// This is more complex for OR conditions - we need to handle joins differently
	// For now, we'll create a subquery-based approach

	var propertyTable string
	var joinColumn string

	switch qb.entityType {
	case EntityTypeContext:
		propertyTable = qb.quoteTableName("ContextProperty")
		joinColumn = "context_id"
	case EntityTypeArtifact:
		propertyTable = qb.quoteTableName("ArtifactProperty")
		joinColumn = "artifact_id"
	case EntityTypeExecution:
		propertyTable = qb.quoteTableName("ExecutionProperty")
		joinColumn = "execution_id"
	}

	// Use the specific value type column based on inferred type
	valueColumn := propRef.ValueType
	condition := qb.buildOperatorCondition(valueColumn, operator, value)

	subquery := fmt.Sprintf("EXISTS (SELECT 1 FROM %s WHERE %s.%s = %s.id AND %s.name = ? AND %s)",
		propertyTable, propertyTable, joinColumn, qb.tablePrefix, propertyTable, condition.condition)

	args := []any{propRef.Name}
	args = append(args, condition.args...)

	return conditionResult{condition: subquery, args: args}
}

// buildOperatorCondition builds a condition string for an operator
func (qb *QueryBuilder) buildOperatorCondition(column string, operator string, value any) conditionResult {
	switch operator {
	case "=":
		return conditionResult{condition: fmt.Sprintf("%s = ?", column), args: []any{value}}
	case "!=":
		return conditionResult{condition: fmt.Sprintf("%s != ?", column), args: []any{value}}
	case ">":
		return conditionResult{condition: fmt.Sprintf("%s > ?", column), args: []any{value}}
	case ">=":
		return conditionResult{condition: fmt.Sprintf("%s >= ?", column), args: []any{value}}
	case "<":
		return conditionResult{condition: fmt.Sprintf("%s < ?", column), args: []any{value}}
	case "<=":
		return conditionResult{condition: fmt.Sprintf("%s <= ?", column), args: []any{value}}
	case "LIKE":
		return conditionResult{condition: fmt.Sprintf("%s LIKE ?", column), args: []any{value}}
	case "ILIKE":
		// Cross-database case-insensitive LIKE using UPPER()
		// This works across MySQL, PostgreSQL, SQLite, and most other databases
		if strValue, ok := value.(string); ok {
			return conditionResult{
				condition: fmt.Sprintf("UPPER(%s) LIKE UPPER(?)", column),
				args:      []any{strValue},
			}
		}
		// Fallback to regular LIKE if value is not a string
		return conditionResult{condition: fmt.Sprintf("%s LIKE ?", column), args: []any{value}}
	case "IN":
		// Handle IN operator with array values
		if valueSlice, ok := value.([]interface{}); ok {
			if len(valueSlice) == 0 {
				// Empty list should return false condition
				return conditionResult{condition: "1 = 0", args: []any{}}
			}
			// Create placeholders for each value
			condition := fmt.Sprintf("%s IN (?%s)", column, strings.Repeat(",?", len(valueSlice)-1))
			return conditionResult{condition: condition, args: valueSlice}
		}
		// Fallback to single value (shouldn't normally happen with proper parsing)
		return conditionResult{condition: fmt.Sprintf("%s IN (?)", column), args: []any{value}}
	default:
		// Default to equality
		return conditionResult{condition: fmt.Sprintf("%s = ?", column), args: []any{value}}
	}
}

// buildCaseInsensitiveLikeCondition builds a cross-database case-insensitive LIKE condition
// This method provides different implementations based on the database type for optimal performance
//
//nolint:staticcheck // Embedded field access is intentional for database dialect checking
func (qb *QueryBuilder) buildCaseInsensitiveLikeCondition(db *gorm.DB, column string, value any) *gorm.DB {
	if strValue, ok := value.(string); ok {
		// Check database type for optimal implementation
		switch db.Dialector.Name() {
		case "postgres":
			// PostgreSQL supports ILIKE natively (most efficient)
			return db.Where(fmt.Sprintf("%s ILIKE ?", column), strValue)
		case "mysql":
			// MySQL: use COLLATE for case-insensitive comparison
			return db.Where(fmt.Sprintf("%s LIKE ? COLLATE utf8mb4_unicode_ci", column), strValue)
		case "sqlite":
			// SQLite: LIKE is case-insensitive by default for ASCII characters
			return db.Where(fmt.Sprintf("%s LIKE ?", column), strValue)
		default:
			// Fallback: use UPPER() function (works on most databases)
			return db.Where(fmt.Sprintf("UPPER(%s) LIKE UPPER(?)", column), strValue)
		}
	}

	// Fallback to regular LIKE if value is not a string
	return db.Where(fmt.Sprintf("%s LIKE ?", column), value)
}
