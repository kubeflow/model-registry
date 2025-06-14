package models

const DocArtifactType = "doc-artifact"

type DocArtifactListOptions struct {
	Pagination
	Name           *string
	ExternalID     *string
	ModelVersionID *int32
}

type DocArtifactAttributes struct {
	Name                     *string
	URI                      *string
	State                    *string
	ArtifactType             *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type DocArtifact interface {
	GetID() *int32
	GetTypeID() *int32
	GetAttributes() *DocArtifactAttributes
	GetProperties() *[]Properties
	GetCustomProperties() *[]Properties
}

type DocArtifactImpl struct {
	ID               *int32
	TypeID           *int32
	Attributes       *DocArtifactAttributes
	Properties       *[]Properties
	CustomProperties *[]Properties
}

func (r *DocArtifactImpl) GetID() *int32 {
	return r.ID
}

func (r *DocArtifactImpl) GetTypeID() *int32 {
	return r.TypeID
}

func (r *DocArtifactImpl) GetAttributes() *DocArtifactAttributes {
	return r.Attributes
}

func (r *DocArtifactImpl) GetProperties() *[]Properties {
	return r.Properties
}

func (r *DocArtifactImpl) GetCustomProperties() *[]Properties {
	return r.CustomProperties
}

type DocArtifactRepository interface {
	GetByID(id int32) (DocArtifact, error)
	List(listOptions DocArtifactListOptions) (*ListWrapper[DocArtifact], error)
	Save(docArtifact DocArtifact, modelVersionID *int32) (DocArtifact, error)
}
