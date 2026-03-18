package entity

import "math"

type Properties struct {
	Name             string
	IsCustomProperty bool
	IntValue         *int32
	DoubleValue      *float64
	StringValue      *string
	BoolValue        *bool
	ByteValue        *[]byte
	ProtoValue       *[]byte
}

func (p *Properties) SetInt64Value(n int64) {
	if n >= math.MinInt32 && n <= math.MaxInt32 {
		n32 := int32(n)
		p.IntValue = &n32
	} else {
		p.IntValue = nil
	}
}

func NewStringProperty(name string, value string, isCustom bool) Properties {
	return Properties{
		Name:             name,
		IsCustomProperty: isCustom,
		StringValue:      &value,
	}
}

func NewIntProperty(name string, value int32, isCustom bool) Properties {
	return Properties{
		Name:             name,
		IsCustomProperty: isCustom,
		IntValue:         &value,
	}
}

func NewDoubleProperty(name string, value float64, isCustom bool) Properties {
	return Properties{
		Name:             name,
		IsCustomProperty: isCustom,
		DoubleValue:      &value,
	}
}

func NewBoolProperty(name string, value bool, isCustom bool) Properties {
	return Properties{
		Name:             name,
		IsCustomProperty: isCustom,
		BoolValue:        &value,
	}
}

func NewByteProperty(name string, value []byte, isCustom bool) Properties {
	return Properties{
		Name:             name,
		IsCustomProperty: isCustom,
		ByteValue:        &value,
	}
}
