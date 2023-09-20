package db

type MetadataProperty interface {
	GetID() int64
	SetID(int64)
	GetName() string
	SetName(string)
	GetIsCustomProperty() bool
	SetIsCustomProperty(bool)
	GetIntValue() *int64
	SetIntValue(*int64)
	GetDoubleValue() *float64
	SetDoubleValue(*float64)
	GetStringValue() *string
	SetStringValue(*string)
	GetByteValue() *[]byte
	SetByteValue(*[]byte)
	GetProtoValue() *[]byte
	SetProtoValue(*[]byte)
	GetBoolValue() *bool
	SetBoolValue(*bool)
	GetTypeURL() *string
	SetTypeURL(*string)
}
