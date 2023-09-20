package db

const TableNameContextProperty = "ContextProperty"

// ContextProperty mapped from table <ContextProperty>
type ContextProperty struct {
	ContextID        int64    `gorm:"column:context_id;primaryKey" json:"-"`
	Name             string   `gorm:"column:name;primaryKey;index:idx_context_property_int,priority:1;index:idx_context_property_string,priority:1;index:idx_context_property_double,priority:1" json:"-"`
	IsCustomProperty bool     `gorm:"column:is_custom_property;primaryKey;index:idx_context_property_int,priority:2;index:idx_context_property_string,priority:2;index:idx_context_property_double,priority:2" json:"-"`
	IntValue         *int64   `gorm:"column:int_value;index:idx_context_property_int,priority:3" json:"-"`
	DoubleValue      *float64 `gorm:"column:double_value;index:idx_context_property_double,priority:3" json:"-"`
	StringValue      *string  `gorm:"column:string_value;index:idx_context_property_string,priority:3" json:"-"`
	ByteValue        *[]byte  `gorm:"column:byte_value" json:"-"`
	ProtoValue       *[]byte  `gorm:"column:proto_value" json:"-"`
	BoolValue        *bool    `gorm:"column:bool_value" json:"-"`
	TypeURL          *string  `gorm:"column:type_url" json:"-"`
}

// TableName ContextProperty's table name
func (*ContextProperty) TableName() string {
	return TableNameContextProperty
}

func (p *ContextProperty) GetID() int64 {
	return p.ContextID
}

func (p *ContextProperty) SetID(i int64) {
	p.ContextID = i
}

func (p *ContextProperty) GetName() string {
	return p.Name
}

func (p *ContextProperty) SetName(s string) {
	p.Name = s
}

func (p *ContextProperty) GetIsCustomProperty() bool {
	return p.IsCustomProperty
}

func (p *ContextProperty) SetIsCustomProperty(b bool) {
	p.IsCustomProperty = b
}

func (p *ContextProperty) GetIntValue() *int64 {
	return p.IntValue
}

func (p *ContextProperty) SetIntValue(i *int64) {
	p.IntValue = i
}

func (p *ContextProperty) GetDoubleValue() *float64 {
	return p.DoubleValue
}

func (p *ContextProperty) SetDoubleValue(f *float64) {
	p.DoubleValue = f
}

func (p *ContextProperty) GetStringValue() *string {
	return p.StringValue
}

func (p *ContextProperty) SetStringValue(s *string) {
	p.StringValue = s
}

func (p *ContextProperty) GetByteValue() *[]byte {
	return p.ByteValue
}

func (p *ContextProperty) SetByteValue(b *[]byte) {
	p.ByteValue = b
}

func (p *ContextProperty) GetProtoValue() *[]byte {
	return p.ProtoValue
}

func (p *ContextProperty) SetProtoValue(b *[]byte) {
	p.ProtoValue = b
}

func (p *ContextProperty) GetBoolValue() *bool {
	return p.BoolValue
}

func (p *ContextProperty) SetBoolValue(b *bool) {
	p.BoolValue = b
}

func (p *ContextProperty) GetTypeURL() *string {
	return p.TypeURL
}

func (p *ContextProperty) SetTypeURL(s *string) {
	p.TypeURL = s
}
