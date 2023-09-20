package db

const TableNameAttribution = "Attribution"

// Attribution mapped from table <Attribution>
type Attribution struct {
	ID         int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"-"`
	ContextID  int64 `gorm:"column:context_id;not null;uniqueIndex:UniqueAttribution,priority:1" json:"-"`
	ArtifactID int64 `gorm:"column:artifact_id;not null;uniqueIndex:UniqueAttribution,priority:2" json:"-"`
	Context    Context
	Artifact   Artifact
}

// TableName Attribution's table name
func (*Attribution) TableName() string {
	return TableNameAttribution
}
