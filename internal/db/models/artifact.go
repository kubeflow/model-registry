package models

import (
	"github.com/kubeflow/model-registry/internal/db/constants"
	"github.com/kubeflow/model-registry/internal/db/filter"
)

var (
	// Use centralized state mappings from constants package
	Artifact_State_name  = constants.ArtifactStateNames
	Artifact_State_value = constants.ArtifactStateMapping
)

type ArtifactListOptions struct {
	Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
	ArtifactType     *string
}

// GetRestEntityType implements the FilterApplier interface
// This enables advanced filtering support for artifacts
func (a *ArtifactListOptions) GetRestEntityType() filter.RestEntityType {
	// Determine the appropriate REST entity type based on artifact type
	if a.ArtifactType != nil {
		switch *a.ArtifactType {
		case "model-artifact":
			return filter.RestEntityModelArtifact
		case "doc-artifact":
			return filter.RestEntityDocArtifact
		case "dataset-artifact":
			return filter.RestEntityDataSet
		case "metric":
			return filter.RestEntityMetric
		case "parameter":
			return filter.RestEntityParameter
		}
	}
	// Default to ModelArtifact if no specific type is provided
	// This allows filtering on common artifact properties
	return filter.RestEntityModelArtifact
}

type Artifact struct {
	ModelArtifact *ModelArtifact
	DocArtifact   *DocArtifact
	DataSet       *DataSet
	Metric        *Metric
	Parameter     *Parameter
}

type ArtifactRepository interface {
	GetByID(id int32) (Artifact, error)
	List(listOptions ArtifactListOptions) (*ListWrapper[Artifact], error)
}
