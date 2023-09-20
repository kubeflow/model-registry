package db

const TableNameAssociation = "Association"

// Association mapped from table <Association>
type Association struct {
	ID          int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"-"`
	ContextID   int64 `gorm:"column:context_id;not null;uniqueIndex:UniqueAssociation,priority:1" json:"-"`
	ExecutionID int64 `gorm:"column:execution_id;not null;uniqueIndex:UniqueAssociation,priority:2" json:"-"`
	Context     Context
	Execution   Execution
}

// TableName Association's table name
func (*Association) TableName() string {
	return TableNameAssociation
}
