package db

const TableNameArtifact = "Artifact"

// Artifact mapped from table <Artifact>
type Artifact struct {
	ID                       int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"-"`
	TypeID                   int64   `gorm:"column:type_id;not null;uniqueIndex:UniqueArtifactTypeName,priority:1" json:"-"`
	URI                      *string `gorm:"column:uri;type:text;index:idx_artifact_uri,priority:1" json:"-"`
	State                    *int8   `gorm:"column:state" json:"-"`
	Name                     *string `gorm:"column:name;type:varchar(255);uniqueIndex:UniqueArtifactTypeName,priority:2" json:"-"`
	ExternalID               *string `gorm:"column:external_id;type:varchar(255);uniqueIndex:idx_artifact_external_id,priority:1" json:"-"`
	CreateTimeSinceEpoch     int64   `gorm:"autoCreateTime:milli;column:create_time_since_epoch;not null;index:idx_artifact_create_time_since_epoch,priority:1" json:"-"`
	LastUpdateTimeSinceEpoch int64   `gorm:"autoUpdateTime:milli;column:last_update_time_since_epoch;not null;index:idx_artifact_last_update_time_since_epoch,priority:1" json:"-"`

	// relationships
	Properties   []ArtifactProperty
	ArtifactType Type          `gorm:"foreignKey:TypeID;references:ID"`
	Attributions []Attribution `gorm:"foreignKey:ArtifactID;references:ID"`
	Events       []Event
}

// TableName Artifact's table name
func (*Artifact) TableName() string {
	return TableNameArtifact
}

type ArtifactState int

const (
	UNKNOWN ArtifactState = iota
	// PENDING A state indicating that the artifact may exist.
	PENDING
	// LIVE A state indicating that the artifact should exist, unless something
	// external to the system deletes it.
	LIVE
	// MARKED_FOR_DELETION A state indicating that the artifact should be deleted.
	MARKED_FOR_DELETION
	// DELETED A state indicating that the artifact has been deleted.
	DELETED
	// ABANDONED A state indicating that the artifact has been abandoned, which may be
	// due to a failed or cancelled execution.
	ABANDONED
	// REFERENCE A state indicating that the artifact is a reference artifact. At
	// execution start time, the orchestrator produces an output artifact for
	// each output key with state PENDING. However, for an intermediate
	// artifact, this first artifact's state will be REFERENCE. Intermediate
	// artifacts emitted during a component's execution will copy the REFERENCE
	// artifact's attributes. At the end of an execution, the artifact state
	// should remain REFERENCE instead of being changed to LIVE.
	REFERENCE
)
