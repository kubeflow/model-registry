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
	Entity[ModelArtifactAttributes]
}

type ModelArtifactImpl = BaseEntity[ModelArtifactAttributes]

type ModelArtifactRepository interface {
	GetByID(id int32) (ModelArtifact, error)
	List(listOptions ModelArtifactListOptions) (*ListWrapper[ModelArtifact], error)
	Save(modelArtifact ModelArtifact, modelVersionID *int32) (ModelArtifact, error)
}
