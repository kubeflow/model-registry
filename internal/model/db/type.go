package db

const TableNameType = "Type"

// Type mapped from table <Type>
type Type struct {
	ID          int64   `gorm:"column:id;not null;primaryKey;autoIncrement:true" json:"-"`
	Name        string  `gorm:"column:name;type:varchar(255);not null;uniqueIndex:idx_type_name,priority:1" json:"-"`
	Version     *string `gorm:"column:version;type:varchar(255)" json:"-"`
	TypeKind    int8    `gorm:"column:type_kind;not null" json:"-"`
	Description *string `gorm:"column:description;type:text" json:"-"`
	InputType   *string `gorm:"column:input_type;type:text" json:"-"`
	OutputType  *string `gorm:"column:output_type;type:text" json:"-"`
	ExternalID  *string `gorm:"column:external_id;type:varchar(255);uniqueIndex:idx_type_external_id,priority:1" json:"-"`

	// relationships
	Properties []TypeProperty
}

// TableName Type's table name
func (*Type) TableName() string {
	return TableNameType
}

type TypeKind int

// artifact type values from ml-metadata table values
const (
	EXECUTION_TYPE TypeKind = iota
	ARTIFACT_TYPE
	CONTEXT_TYPE
)
