package models

const DocArtifactType = "doc-artifact"

type DocArtifactListOptions struct {
	Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
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
	Entity[DocArtifactAttributes]
}

type DocArtifactImpl = BaseEntity[DocArtifactAttributes]

type DocArtifactRepository interface {
	GetByID(id int32) (DocArtifact, error)
	List(listOptions DocArtifactListOptions) (*ListWrapper[DocArtifact], error)
	Save(docArtifact DocArtifact, parentResourceID *int32) (DocArtifact, error)
}
