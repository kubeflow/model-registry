package db

const TableNameEvent = "Event"

// Event mapped from table <Event>
type Event struct {
	ID                     int64 `gorm:"column:id;not null;primaryKey;autoIncrement:true" json:"-"`
	ArtifactID             int64 `gorm:"column:artifact_id;not null;uniqueIndex:UniqueEvent,priority:1" json:"-"`
	ExecutionID            int64 `gorm:"column:execution_id;not null;uniqueIndex:UniqueEvent,priority:2;index:idx_event_execution_id,priority:1" json:"-"`
	Type                   int8  `gorm:"column:type;not null;uniqueIndex:UniqueEvent,priority:3" json:"-"`
	MillisecondsSinceEpoch int64 `gorm:"autoCreateTime:milli;column:milliseconds_since_epoch;not null" json:"-"`

	// relationships
	PathSteps []EventPath
	Artifact  Artifact
	Execution Execution
}

// TableName Event's table name
func (*Event) TableName() string {
	return TableNameEvent
}

type EventType int

// Events distinguish between an artifact that is written by the execution
// (possibly as a cache), versus artifacts that are part of the declared
// output of the Execution. For more information on what DECLARED_ means,
// see the comment on the message.
const (
	EVENT_TYPE_UNKNOWN EventType = iota
	DECLARED_OUTPUT
	DECLARED_INPUT
	INPUT
	OUTPUT
	INTERNAL_INPUT
	INTERNAL_OUTPUT
	PENDING_OUTPUT
)
