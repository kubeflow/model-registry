package utils

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/platform/db/schema"
	"gorm.io/gorm"
)

func getTableName(db *gorm.DB, model any) string {
	stmt := &gorm.Statement{DB: db}
	err := stmt.Parse(model)
	if err != nil {
		switch model.(type) {
		case *schema.ParentContext:
			return db.NamingStrategy.TableName("ParentContext")
		case *schema.Attribution:
			return db.NamingStrategy.TableName("Attribution")
		case *schema.Context:
			return db.NamingStrategy.TableName("Context")
		case *schema.Artifact:
			return db.NamingStrategy.TableName("Artifact")
		case *schema.ContextProperty:
			return db.NamingStrategy.TableName("ContextProperty")
		case *schema.Association:
			return db.NamingStrategy.TableName("Association")
		case *schema.Execution:
			return db.NamingStrategy.TableName("Execution")
		default:
			return "unknown_table"
		}
	}
	return stmt.Quote(stmt.Schema.Table)
}

func BuildParentContextJoin(db *gorm.DB) string {
	parentTable := getTableName(db, &schema.ParentContext{})
	contextTable := getTableName(db, &schema.Context{})
	return fmt.Sprintf("JOIN %s ON %s.context_id = %s.id",
		parentTable, parentTable, contextTable)
}

func BuildAttributionJoin(db *gorm.DB) string {
	attributionTable := getTableName(db, &schema.Attribution{})
	artifactTable := getTableName(db, &schema.Artifact{})
	return fmt.Sprintf("JOIN %s ON %s.artifact_id = %s.id",
		attributionTable, attributionTable, artifactTable)
}

func BuildAssociationJoin(db *gorm.DB) string {
	associationTable := getTableName(db, &schema.Association{})
	executionTable := getTableName(db, &schema.Execution{})
	return fmt.Sprintf("JOIN %s ON %s.execution_id = %s.id",
		associationTable, associationTable, executionTable)
}

func BuildContextPropertyJoin(db *gorm.DB, propertyName string) string {
	propertyTable := getTableName(db, &schema.ContextProperty{})
	contextTable := getTableName(db, &schema.Context{})
	return fmt.Sprintf("JOIN %s ON %s.context_id = %s.id AND %s.name = '%s'",
		propertyTable, propertyTable, contextTable, propertyTable, propertyName)
}

func GetColumnRef(db *gorm.DB, model any, column string) string {
	tableName := getTableName(db, model)
	return fmt.Sprintf("%s.%s", tableName, db.NamingStrategy.ColumnName("", column))
}

func GetTableName(db *gorm.DB, model any) string {
	return getTableName(db, model)
}
