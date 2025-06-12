package models

type TypeAttributes struct {
	Name        *string
	Version     *string
	TypeKind    *int32
	Description *string
	InputType   *string
	OutputType  *string
	ExternalID  *string
}

type Type interface {
	GetID() *int32
	GetAttributes() *TypeAttributes
}

type TypeImpl struct {
	ID         *int32
	Attributes *TypeAttributes
}

func (t *TypeImpl) GetID() *int32 {
	return t.ID
}

func (t *TypeImpl) GetAttributes() *TypeAttributes {
	return t.Attributes
}

type TypeRepository interface {
	GetAll() ([]Type, error)
}
