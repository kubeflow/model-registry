package models

type Entity[T any] interface {
	GetID() *int32
	GetTypeID() *int32
	GetAttributes() *T
	GetProperties() *[]Properties
	GetCustomProperties() *[]Properties
}

type BaseEntity[T any] struct {
	ID               *int32
	TypeID           *int32
	Attributes       *T
	Properties       *[]Properties
	CustomProperties *[]Properties
}

func (r *BaseEntity[T]) GetID() *int32 {
	return r.ID
}

func (r *BaseEntity[T]) GetTypeID() *int32 {
	return r.TypeID
}

func (r *BaseEntity[T]) GetAttributes() *T {
	return r.Attributes
}

func (r *BaseEntity[T]) GetProperties() *[]Properties {
	return r.Properties
}

func (r *BaseEntity[T]) GetCustomProperties() *[]Properties {
	return r.CustomProperties
}

// GetCustomProperty returns the value of a custom property by name
func (r *BaseEntity[T]) GetCustomProperty(name string) *string {
	if r.CustomProperties == nil {
		return nil
	}

	for _, prop := range *r.CustomProperties {
		if prop.Name == name && prop.StringValue != nil {
			return prop.StringValue
		}
	}
	return nil
}
