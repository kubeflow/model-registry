// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package schema

const TableNameParentContext = "ParentContext"

// ParentContext mapped from table <ParentContext>
type ParentContext struct {
	ContextID       int32 `gorm:"column:context_id;primaryKey" json:"context_id"`
	ParentContextID int32 `gorm:"column:parent_context_id;primaryKey" json:"parent_context_id"`
}

// TableName ParentContext's table name
func (*ParentContext) TableName() string {
	return TableNameParentContext
}
