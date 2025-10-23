package filter

import (
	"fmt"
	"strings"

	"github.com/kubeflow/model-registry/internal/db/constants"
	"github.com/kubeflow/model-registry/internal/db/dbutil"
	"gorm.io/gorm"
)

// EntityType represents the type of entity for proper query building
type EntityType string

const (
	EntityTypeContext   EntityType = "context"
	EntityTypeArtifact  EntityType = "artifact"
	EntityTypeExecution EntityType = "execution"
)

// EntityMappingFunctions defines the interface for entity type mapping functions
// This allows different packages (like catalog) to provide their own entity mappings
type EntityMappingFunctions interface {
	// GetMLMDEntityType maps a REST entity type to its underlying MLMD entity type
	GetMLMDEntityType(restEntityType RestEntityType) EntityType

	// GetPropertyDefinitionForRestEntity returns property definition for a REST entity type
	GetPropertyDefinitionForRestEntity(restEntityType RestEntityType, propertyName string) PropertyDefinition

	// IsChildEntity returns true if the REST entity type uses prefixed names (parentId:name)
	IsChildEntity(entityType RestEntityType) bool
}

// QueryBuilder builds GORM queries from filter expressions
// It handles special cases like prefixed names for child entities (e.g., ModelVersion, ExperimentRun)
// where names are stored as "parentId:actualName" in the database
type QueryBuilder struct {
	entityType     EntityType
	restEntityType RestEntityType
	tablePrefix    string
	joinCounter    int
	db             *gorm.DB               // Added to access naming strategy
	mappingFuncs   EntityMappingFunctions // Entity mapping functions
}

// NewQueryBuilderForRestEntity creates a new query builder for the specified REST entity type
// If mappingFuncs is nil, it falls back to the global functions
func NewQueryBuilderForRestEntity(restEntityType RestEntityType, mappingFuncs EntityMappingFunctions) *QueryBuilder {
	// Use default mappings if none provided
	if mappingFuncs == nil {
		mappingFuncs = &defaultEntityMappings{}
	}

	// Get the underlying MLMD entity type
	entityType := mappingFuncs.GetMLMDEntityType(restEntityType)

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
		mappingFuncs:   mappingFuncs,
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
	// Extract unquoted table name if it was already quoted
	unquotedName := strings.Trim(qb.tablePrefix, "`\"")
	qb.tablePrefix = dbutil.QuoteTableName(qb.db, unquotedName)
}

// quoteTableName quotes a table name based on database dialect
func (qb *QueryBuilder) quoteTableName(tableName string) string {
	return dbutil.QuoteTableName(qb.db, tableName)
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
	// Valid type suffixes: string_value, double_value, int_value, bool_value
	var explicitType string
	if parts := strings.Split(propertyName, "."); len(parts) >= 2 {
		lastPart := parts[len(parts)-1]
		// Only treat as type suffix if it's a valid value type
		if lastPart == "string_value" || lastPart == "double_value" || lastPart == "int_value" || lastPart == "bool_value" {
			// Reconstruct property name without the type suffix
			propertyName = strings.Join(parts[:len(parts)-1], ".")
			explicitType = lastPart
		}
		// Otherwise, keep the full path as property name
	}

	// Use REST entity type-aware property mapping if available
	if qb.restEntityType != "" {
		propDef = qb.mappingFuncs.GetPropertyDefinitionForRestEntity(qb.restEntityType, propertyName)
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
		Name:        propName,
		IsCustom:    propDef.Location == Custom,
		ValueType:   propDef.ValueType,
		PropertyDef: propDef, // Store full property definition for advanced handling
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
	// Use the property definition from the PropertyReference (already looked up with REST entity mappings)
	propDef := propRef.PropertyDef

	switch propDef.Location {
	case EntityTable:
		return qb.buildEntityTablePropertyCondition(db, propRef, operator, value)
	case PropertyTable, Custom:
		return qb.buildPropertyTableCondition(db, propRef, operator, value)
	case RelatedEntity:
		return qb.buildRelatedEntityPropertyCondition(db, propDef, propRef.ValueType, operator, value)
	default:
		return qb.buildEntityTablePropertyCondition(db, propRef, operator, value)
	}
}

// buildPropertyConditionString builds a condition string for a property
func (qb *QueryBuilder) buildPropertyConditionString(propRef *PropertyReference, operator string, value any) conditionResult {
	// Use the property definition from the PropertyReference (already looked up with REST entity mappings)
	propDef := propRef.PropertyDef

	switch propDef.Location {
	case EntityTable:
		return qb.buildEntityTablePropertyConditionString(propRef, operator, value)
	case PropertyTable, Custom:
		return qb.buildPropertyTableConditionString(propRef, operator, value)
	case RelatedEntity:
		return qb.buildRelatedEntityPropertyConditionString(propDef, propRef.ValueType, operator, value)
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
	if qb.restEntityType != "" && propRef.Name == "name" && qb.mappingFuncs.IsChildEntity(qb.restEntityType) {
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
	if qb.restEntityType != "" && propRef.Name == "name" && qb.mappingFuncs.IsChildEntity(qb.restEntityType) {
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
	var valueColumn string
	if propRef.ValueType == ArrayValueType {
		valueColumn = fmt.Sprintf("%s.%s", alias, StringValueType)
	} else {
		valueColumn = fmt.Sprintf("%s.%s", alias, propRef.ValueType)
	}

	// Use cross-database case-insensitive LIKE for ILIKE operator
	if operator == "ILIKE" {
		return qb.buildCaseInsensitiveLikeCondition(db, valueColumn, value)
	}

	var condition conditionResult
	if propRef.ValueType == ArrayValueType && db.Name() == "postgres" {
		condition = qb.buildJSONOperatorCondition(valueColumn, operator, value)
	} else {
		condition = qb.buildOperatorCondition(valueColumn, operator, value)
	}
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

// buildRelatedEntityPropertyCondition builds a condition for properties in related entities (requires multiple JOINs)
func (qb *QueryBuilder) buildRelatedEntityPropertyCondition(db *gorm.DB, propDef PropertyDefinition, explicitType string, operator string, value any) *gorm.DB {
	// Increment join counter for unique alias
	qb.joinCounter++

	// For artifact filtering from a Context (model), we need to:
	// 1. JOIN Attribution (relation table)
	// 2. JOIN Artifact (related entity)
	// 3. JOIN ArtifactProperty (property table of related entity)

	// Currently only supporting artifact filtering from Context entities
	if qb.entityType != EntityTypeContext || propDef.RelatedEntityType != RelatedEntityArtifact {
		// Fallback to regular property handling if not supported
		return db
	}

	// Create unique aliases for this join chain
	attributionAlias := fmt.Sprintf("attr_%d", qb.joinCounter)
	artifactAlias := fmt.Sprintf("art_%d", qb.joinCounter)
	propAlias := fmt.Sprintf("artprop_%d", qb.joinCounter)

	// Quote table names for database compatibility
	attributionTable := qb.quoteTableName("Attribution")
	artifactTable := qb.quoteTableName("Artifact")
	artifactPropertyTable := qb.quoteTableName("ArtifactProperty")

	// Build the JOIN chain:
	// JOIN Attribution attr_N ON attr_N.context_id = Context.id
	join1 := fmt.Sprintf("JOIN %s %s ON %s.context_id = %s.id",
		attributionTable, attributionAlias, attributionAlias, qb.tablePrefix)
	db = db.Joins(join1)

	// JOIN Artifact art_N ON art_N.id = attr_N.artifact_id
	join2 := fmt.Sprintf("JOIN %s %s ON %s.id = %s.artifact_id",
		artifactTable, artifactAlias, artifactAlias, attributionAlias)
	db = db.Joins(join2)

	// JOIN ArtifactProperty artprop_N ON artprop_N.artifact_id = art_N.id
	join3 := fmt.Sprintf("JOIN %s %s ON %s.artifact_id = %s.id",
		artifactPropertyTable, propAlias, propAlias, artifactAlias)
	db = db.Joins(join3)

	// Add condition for property name
	db = db.Where(fmt.Sprintf("%s.name = ?", propAlias), propDef.RelatedProperty)

	// Determine value type: use explicit type if provided, otherwise infer from value
	valueType := explicitType
	inferredAsInt := false
	if valueType == "" {
		// No explicit type, infer from value
		valueType = qb.inferValueTypeFromInterface(value)
		// Track if we inferred an integer type (for special handling below)
		if valueType == IntValueType {
			inferredAsInt = true
		}
	}

	// Special handling for integer literals without explicit type:
	// Query BOTH int_value and double_value to handle data stored in either column.
	// This prevents silent query failures when data type doesn't match user's expectation.
	if inferredAsInt {
		// Build conditions for both int_value and double_value
		intColumn := fmt.Sprintf("%s.int_value", propAlias)
		doubleColumn := fmt.Sprintf("%s.double_value", propAlias)

		// Use cross-database case-insensitive LIKE for ILIKE operator
		if operator == "ILIKE" {
			// ILIKE doesn't make sense for numeric types, but handle it anyway
			return qb.buildCaseInsensitiveLikeCondition(db, intColumn, value)
		}

		intCondition := qb.buildOperatorCondition(intColumn, operator, value)
		doubleCondition := qb.buildOperatorCondition(doubleColumn, operator, value)

		// Combine with OR to find values in either column
		combinedCondition := fmt.Sprintf("(%s OR %s)", intCondition.condition, doubleCondition.condition)
		combinedArgs := append(intCondition.args, doubleCondition.args...)

		return db.Where(combinedCondition, combinedArgs...)
	}

	// For explicit types or non-integer types, use the specified column
	valueColumn := fmt.Sprintf("%s.%s", propAlias, valueType)

	// Use cross-database case-insensitive LIKE for ILIKE operator
	if operator == "ILIKE" {
		return qb.buildCaseInsensitiveLikeCondition(db, valueColumn, value)
	}

	condition := qb.buildOperatorCondition(valueColumn, operator, value)
	return db.Where(condition.condition, condition.args...)
}

// buildRelatedEntityPropertyConditionString builds a condition string for properties in related entities
func (qb *QueryBuilder) buildRelatedEntityPropertyConditionString(propDef PropertyDefinition, explicitType string, operator string, value any) conditionResult {
	// Increment join counter for unique alias
	qb.joinCounter++

	// Currently only supporting artifact filtering from Context entities
	if qb.entityType != EntityTypeContext || propDef.RelatedEntityType != RelatedEntityArtifact {
		// Fallback - return empty condition
		return conditionResult{condition: "1=1", args: []any{}}
	}

	// Create unique aliases
	attributionAlias := fmt.Sprintf("attr_%d", qb.joinCounter)
	artifactAlias := fmt.Sprintf("art_%d", qb.joinCounter)
	propAlias := fmt.Sprintf("artprop_%d", qb.joinCounter)

	// Quote table names
	attributionTable := qb.quoteTableName("Attribution")
	artifactTable := qb.quoteTableName("Artifact")
	artifactPropertyTable := qb.quoteTableName("ArtifactProperty")

	// Build EXISTS subquery with the JOIN chain
	// EXISTS (
	//   SELECT 1 FROM Attribution attr_N
	//   JOIN Artifact art_N ON art_N.id = attr_N.artifact_id
	//   JOIN ArtifactProperty artprop_N ON artprop_N.artifact_id = art_N.id
	//   WHERE attr_N.context_id = Context.id
	//     AND artprop_N.name = ?
	//     AND artprop_N.value_column OPERATOR ?
	// )

	// Determine value type: use explicit type if provided, otherwise infer from value
	valueType := explicitType
	inferredAsInt := false
	if valueType == "" {
		// No explicit type, infer from value
		valueType = qb.inferValueTypeFromInterface(value)
		// Track if we inferred an integer type (for special handling below)
		if valueType == IntValueType {
			inferredAsInt = true
		}
	}

	// Build the value condition
	var valueCondition conditionResult

	// Special handling for integer literals without explicit type:
	// Query BOTH int_value and double_value to handle data stored in either column.
	// This prevents silent query failures when data type doesn't match user's expectation.
	if inferredAsInt {
		// Build conditions for both int_value and double_value
		intColumn := fmt.Sprintf("%s.int_value", propAlias)
		doubleColumn := fmt.Sprintf("%s.double_value", propAlias)

		intCondition := qb.buildOperatorCondition(intColumn, operator, value)
		doubleCondition := qb.buildOperatorCondition(doubleColumn, operator, value)

		// Combine with OR to find values in either column
		combinedCondition := fmt.Sprintf("(%s OR %s)", intCondition.condition, doubleCondition.condition)
		combinedArgs := append(intCondition.args, doubleCondition.args...)

		valueCondition = conditionResult{
			condition: combinedCondition,
			args:      combinedArgs,
		}
	} else {
		// For explicit types or non-integer types, use the specified column
		valueColumn := fmt.Sprintf("%s.%s", propAlias, valueType)
		valueCondition = qb.buildOperatorCondition(valueColumn, operator, value)
	}

	// Build the complete EXISTS subquery
	subquery := fmt.Sprintf(
		"EXISTS (SELECT 1 FROM %s %s "+
			"JOIN %s %s ON %s.id = %s.artifact_id "+
			"JOIN %s %s ON %s.artifact_id = %s.id "+
			"WHERE %s.context_id = %s.id AND %s.name = ? AND %s)",
		attributionTable, attributionAlias,
		artifactTable, artifactAlias, artifactAlias, attributionAlias,
		artifactPropertyTable, propAlias, propAlias, artifactAlias,
		attributionAlias, qb.tablePrefix, propAlias, valueCondition.condition)

	args := []any{propDef.RelatedProperty}
	args = append(args, valueCondition.args...)

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
		switch db.Name() {
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

// buildJSONOperatorCondition builds a condition string for an operator on a JSON array
func (qb *QueryBuilder) buildJSONOperatorCondition(column string, operator string, value any) conditionResult {
	switch operator {
	case "IN":
		// Handle IN operator with array values
		if valueSlice, ok := value.([]any); ok {
			if len(valueSlice) == 0 {
				// Empty list should return false condition
				return conditionResult{condition: "1 = 0", args: []any{}}
			}
			// Create placeholders for each value
			return conditionResult{
				condition: fmt.Sprintf("%s IS JSON ARRAY AND %s::jsonb ? array[?%s]", column, column, strings.Repeat(",?", len(valueSlice)-1)),
				args:      append([]any{gorm.Expr("?|")}, valueSlice...),
			}
		}
		// Fallback to single value (shouldn't normally happen with proper parsing)
		fallthrough
	case "=":
		return conditionResult{condition: fmt.Sprintf("%s IS JSON ARRAY AND %s::jsonb ? array[?]", column, column), args: []any{gorm.Expr("?|"), value}}
	case "!=":
		return conditionResult{condition: fmt.Sprintf("%s IS NOT JSON ARRAY OR NOT %s::jsonb ? array[?]", column, column), args: []any{gorm.Expr("?|"), value}}
	default:
		// Pass through anything else
		return qb.buildOperatorCondition(column, operator, value)
	}
}

// defaultEntityMappings implements EntityMappingFunctions for the model registry
type defaultEntityMappings struct{}

func (d *defaultEntityMappings) GetMLMDEntityType(restEntityType RestEntityType) EntityType {
	return GetMLMDEntityType(restEntityType)
}

func (d *defaultEntityMappings) GetPropertyDefinitionForRestEntity(restEntityType RestEntityType, propertyName string) PropertyDefinition {
	return GetPropertyDefinitionForRestEntity(restEntityType, propertyName)
}

func (d *defaultEntityMappings) IsChildEntity(entityType RestEntityType) bool {
	return isChildEntity(entityType)
}
