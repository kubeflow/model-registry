package models

const DataSetType = "dataset-artifact"

type DataSetListOptions struct {
	Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
}

type DataSetAttributes struct {
	Name                     *string
	URI                      *string
	State                    *string
	ArtifactType             *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type DataSet interface {
	Entity[DataSetAttributes]
}

type DataSetImpl = BaseEntity[DataSetAttributes]

type DataSetRepository interface {
	GetByID(id int32) (DataSet, error)
	List(listOptions DataSetListOptions) (*ListWrapper[DataSet], error)
	Save(dataSet DataSet, parentResourceID *int32) (DataSet, error)
}
