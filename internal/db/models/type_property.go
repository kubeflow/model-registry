package models

type TypeProperty interface {
	GetTypeID() int32
	GetName() string
	GetDataType() *int32
}

var _ TypeProperty = (*TypePropertyImpl)(nil)

type TypePropertyImpl struct {
	TypeID   int32
	Name     string
	DataType *int32
}

func (tp *TypePropertyImpl) GetTypeID() int32 {
	return tp.TypeID
}

func (tp *TypePropertyImpl) GetName() string {
	return tp.Name
}

func (tp *TypePropertyImpl) GetDataType() *int32 {
	return tp.DataType
}

type TypePropertyRepository interface {
	// Save stores a type property if it doesn't exist.
	Save(tp TypeProperty) (TypeProperty, error)
}
