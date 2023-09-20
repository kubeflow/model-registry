package db

const TableNameTypeProperty = "TypeProperty"

// TypeProperty mapped from table <TypeProperty>
type TypeProperty struct {
	TypeID   int64  `gorm:"column:type_id;primaryKey" json:"-"`
	Name     string `gorm:"column:name;primaryKey" json:"-"`
	DataType int32 `gorm:"column:data_type;not null" json:"-"`
}

// TableName TypeProperty's table name
func (*TypeProperty) TableName() string {
	return TableNameTypeProperty
}
