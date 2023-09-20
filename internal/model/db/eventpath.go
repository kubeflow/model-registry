package db

const TableNameEventPath = "EventPath"

// EventPath mapped from table <EventPath>
type EventPath struct {
	EventID     int64   `gorm:"column:event_id;not null;index:idx_eventpath_event_id,priority:1" json:"-"`
	IsIndexStep bool    `gorm:"column:is_index_step;not null" json:"-"`
	StepIndex   *int    `gorm:"column:step_index" json:"-"`
	StepKey     *string `gorm:"column:step_key" json:"-"`
}

// TableName EventPath's table name
func (*EventPath) TableName() string {
	return TableNameEventPath
}
