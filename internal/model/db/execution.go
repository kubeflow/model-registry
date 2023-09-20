package db

const TableNameExecution = "Execution"

// Execution mapped from table <Execution>
type Execution struct {
	ID                       int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"-"`
	TypeID                   int64   `gorm:"column:type_id;not null;uniqueIndex:UniqueExecutionTypeName,priority:1" json:"-"`
	LastKnownState           *int8   `gorm:"column:last_known_state" json:"-"`
	Name                     *string `gorm:"column:name;type:varchar(255);uniqueIndex:UniqueExecutionTypeName,priority:2" json:"-"`
	ExternalID               *string `gorm:"column:external_id;type:varchar(255);uniqueIndex:idx_execution_external_id,priority:1" json:"-"`
	CreateTimeSinceEpoch     int64   `gorm:"autoCreateTime:milli;column:create_time_since_epoch;not null;index:idx_execution_create_time_since_epoch,priority:1" json:"-"`
	LastUpdateTimeSinceEpoch int64   `gorm:"autoUpdateTime:milli;column:last_update_time_since_epoch;not null;index:idx_execution_last_update_time_since_epoch,priority:1" json:"-"`

	// relationships
	Properties    []ExecutionProperty
	ExecutionType Type          `gorm:"foreignKey:TypeID;references:ID"`
	Associations  []Association `gorm:"foreignKey:ExecutionID;references:ID"`
	Events        []Event
}

// TableName Execution's table name
func (*Execution) TableName() string {
	return TableNameExecution
}

type ExecutionState int

const (
	EXECUTION_STATE_UNKNOWN ExecutionState = iota
	NEW
	RUNNING
	COMPLETE
	FAILED
	CACHED
	CANCELED
)
