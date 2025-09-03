package models

import "github.com/kubeflow/model-registry/internal/db/filter"

const ModelArtifactType = "model-artifact"

type ModelArtifactListOptions struct {
	Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
}

// GetRestEntityType implements the FilterApplier interface
func (m *ModelArtifactListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityModelArtifact
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
	Save(modelArtifact ModelArtifact, parentResourceID *int32) (ModelArtifact, error)
}
