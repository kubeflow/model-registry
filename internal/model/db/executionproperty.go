package db

const TableNameExecutionProperty = "ExecutionProperty"

// ExecutionProperty mapped from table <ExecutionProperty>
type ExecutionProperty struct {
	ExecutionID      int64    `gorm:"column:execution_id;primaryKey" json:"-"`
	Name             string   `gorm:"column:name;primaryKey;index:idx_execution_property_string,priority:1;index:idx_execution_property_int,priority:1;index:idx_execution_property_double,priority:1" json:"-"`
	IsCustomProperty bool     `gorm:"column:is_custom_property;primaryKey;index:idx_execution_property_string,priority:2;index:idx_execution_property_int,priority:2;index:idx_execution_property_double,priority:2" json:"-"`
	IntValue         *int64   `gorm:"column:int_value;index:idx_execution_property_int,priority:3" json:"-"`
	DoubleValue      *float64 `gorm:"column:double_value;index:idx_execution_property_double,priority:3" json:"-"`
	StringValue      *string  `gorm:"column:string_value;index:idx_execution_property_string,priority:3" json:"-"`
	ByteValue        *[]byte  `gorm:"column:byte_value" json:"-"`
	ProtoValue       *[]byte  `gorm:"column:proto_value" json:"-"`
	BoolValue        *bool    `gorm:"column:bool_value" json:"-"`
	TypeURL          *string  `gorm:"column:type_url" json:"-"`
}

// TableName ExecutionProperty's table name
func (*ExecutionProperty) TableName() string {
	return TableNameExecutionProperty
}

func (p *ExecutionProperty) GetID() int64 {
	return p.ExecutionID
}

func (p *ExecutionProperty) SetID(i int64) {
	p.ExecutionID = i
}

func (p *ExecutionProperty) GetName() string {
	return p.Name
}

func (p *ExecutionProperty) SetName(s string) {
	p.Name = s
}

func (p *ExecutionProperty) GetIsCustomProperty() bool {
	return p.IsCustomProperty
}

func (p *ExecutionProperty) SetIsCustomProperty(b bool) {
	p.IsCustomProperty = b
}

func (p *ExecutionProperty) GetIntValue() *int64 {
	return p.IntValue
}

func (p *ExecutionProperty) SetIntValue(i *int64) {
	p.IntValue = i
}

func (p *ExecutionProperty) GetDoubleValue() *float64 {
	return p.DoubleValue
}

func (p *ExecutionProperty) SetDoubleValue(f *float64) {
	p.DoubleValue = f
}

func (p *ExecutionProperty) GetStringValue() *string {
	return p.StringValue
}

func (p *ExecutionProperty) SetStringValue(s *string) {
	p.StringValue = s
}

func (p *ExecutionProperty) GetByteValue() *[]byte {
	return p.ByteValue
}

func (p *ExecutionProperty) SetByteValue(b *[]byte) {
	p.ByteValue = b
}

func (p *ExecutionProperty) GetProtoValue() *[]byte {
	return p.ProtoValue
}

func (p *ExecutionProperty) SetProtoValue(b *[]byte) {
	p.ProtoValue = b
}

func (p *ExecutionProperty) GetBoolValue() *bool {
	return p.BoolValue
}

func (p *ExecutionProperty) SetBoolValue(b *bool) {
	p.BoolValue = b
}

func (p *ExecutionProperty) GetTypeURL() *string {
	return p.TypeURL
}

func (p *ExecutionProperty) SetTypeURL(s *string) {
	p.TypeURL = s
}
