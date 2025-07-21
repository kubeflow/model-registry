package models

var (
	Artifact_State_name = map[int32]string{
		0: "UNKNOWN",
		1: "PENDING",
		2: "LIVE",
		3: "MARKED_FOR_DELETION",
		4: "DELETED",
		5: "ABANDONED",
		6: "REFERENCE",
	}
	Artifact_State_value = map[string]int32{
		"UNKNOWN":             0,
		"PENDING":             1,
		"LIVE":                2,
		"MARKED_FOR_DELETION": 3,
		"DELETED":             4,
		"ABANDONED":           5,
		"REFERENCE":           6,
	}
)

type ArtifactListOptions struct {
	Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
	ArtifactType     *string
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
