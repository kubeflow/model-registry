package models

const ModelArtifactType = "model-artifact"

type ModelArtifactListOptions struct {
	Pagination
	Name           *string
	ExternalID     *string
	ModelVersionID *int32
}

type ModelArtifactAttributes struct {
	Name                     *string
	URI                      *string
	State                    *string
	ArtifactType             *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type ModelArtifact interface {
	GetID() *int32
	GetTypeID() *int32
	GetAttributes() *ModelArtifactAttributes
	GetProperties() *[]Properties
	GetCustomProperties() *[]Properties
}

type ModelArtifactImpl struct {
	ID               *int32
	TypeID           *int32
	Attributes       *ModelArtifactAttributes
	Properties       *[]Properties
	CustomProperties *[]Properties
}

func (r *ModelArtifactImpl) GetID() *int32 {
	return r.ID
}

func (r *ModelArtifactImpl) GetTypeID() *int32 {
	return r.TypeID
}

func (r *ModelArtifactImpl) GetAttributes() *ModelArtifactAttributes {
	return r.Attributes
}

func (r *ModelArtifactImpl) GetProperties() *[]Properties {
	return r.Properties
}

func (r *ModelArtifactImpl) GetCustomProperties() *[]Properties {
	return r.CustomProperties
}

type ModelArtifactRepository interface {
	GetByID(id int32) (ModelArtifact, error)
	List(listOptions ModelArtifactListOptions) (*ListWrapper[ModelArtifact], error)
	Save(modelArtifact ModelArtifact, modelVersionID *int32) (ModelArtifact, error)
}
