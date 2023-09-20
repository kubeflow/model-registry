package db

const TableNameArtifactProperty = "ArtifactProperty"

// ArtifactProperty mapped from table <ArtifactProperty>
type ArtifactProperty struct {
	ArtifactID       int64    `gorm:"column:artifact_id;primaryKey" json:"-"`
	Name             string   `gorm:"column:name;primaryKey;index:idx_artifact_property_double,priority:1;index:idx_artifact_property_string,priority:1;index:idx_artifact_property_int,priority:1" json:"-"`
	IsCustomProperty bool     `gorm:"column:is_custom_property;primaryKey;index:idx_artifact_property_double,priority:2;index:idx_artifact_property_string,priority:2;index:idx_artifact_property_int,priority:2" json:"-"`
	IntValue         *int64   `gorm:"column:int_value;index:idx_artifact_property_int,priority:3" json:"-"`
	DoubleValue      *float64 `gorm:"column:double_value;index:idx_artifact_property_double,priority:3" json:"-"`
	StringValue      *string  `gorm:"column:string_value;index:idx_artifact_property_string,priority:3" json:"-"`
	ByteValue        *[]byte  `gorm:"column:byte_value" json:"-"`
	ProtoValue       *[]byte  `gorm:"column:proto_value" json:"-"`
	BoolValue        *bool    `gorm:"column:bool_value" json:"-"`
	TypeURL          *string  `gorm:"column:type_url" json:"-"`
}

// TableName ArtifactProperty's table name
func (*ArtifactProperty) TableName() string {
	return TableNameArtifactProperty
}

func (p *ArtifactProperty) GetID() int64 {
	return p.ArtifactID
}

func (p *ArtifactProperty) SetID(i int64) {
	p.ArtifactID = i
}

func (p *ArtifactProperty) GetName() string {
	return p.Name
}

func (p *ArtifactProperty) SetName(s string) {
	p.Name = s
}

func (p *ArtifactProperty) GetIsCustomProperty() bool {
	return p.IsCustomProperty
}

func (p *ArtifactProperty) SetIsCustomProperty(b bool) {
	p.IsCustomProperty = b
}

func (p *ArtifactProperty) GetIntValue() *int64 {
	return p.IntValue
}

func (p *ArtifactProperty) SetIntValue(i *int64) {
	p.IntValue = i
}

func (p *ArtifactProperty) GetDoubleValue() *float64 {
	return p.DoubleValue
}

func (p *ArtifactProperty) SetDoubleValue(f *float64) {
	p.DoubleValue = f
}

func (p *ArtifactProperty) GetStringValue() *string {
	return p.StringValue
}

func (p *ArtifactProperty) SetStringValue(s *string) {
	p.StringValue = s
}

func (p *ArtifactProperty) GetByteValue() *[]byte {
	return p.ByteValue
}

func (p *ArtifactProperty) SetByteValue(b *[]byte) {
	p.ByteValue = b
}

func (p *ArtifactProperty) GetProtoValue() *[]byte {
	return p.ProtoValue
}

func (p *ArtifactProperty) SetProtoValue(b *[]byte) {
	p.ProtoValue = b
}

func (p *ArtifactProperty) GetBoolValue() *bool {
	return p.BoolValue
}

func (p *ArtifactProperty) SetBoolValue(b *bool) {
	p.BoolValue = b
}

func (p *ArtifactProperty) GetTypeURL() *string {
	return p.TypeURL
}

func (p *ArtifactProperty) SetTypeURL(s *string) {
	p.TypeURL = s
}
