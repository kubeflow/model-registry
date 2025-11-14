package schema

import "github.com/lib/pq"

const TableNameArtifactPropertyOption = "artifact_property_options"

// ArtifactPropertyOption mapped from materialized view <artifact_property_options>
type ArtifactPropertyOption struct {
	TypeID           int32           `gorm:"column:type_id;not null" json:"type_id"`
	Name             string          `gorm:"column:name;not null" json:"name"`
	IsCustomProperty bool            `gorm:"column:is_custom_property;not null" json:"is_custom_property"`
	StringValue      *pq.StringArray `gorm:"column:string_value;type:text[]" json:"string_value"`
	ArrayValue       *pq.StringArray `gorm:"column:array_value;type:text[]" json:"array_value"`
	MinDoubleValue   *float64        `gorm:"column:min_double_value" json:"min_double_value"`
	MaxDoubleValue   *float64        `gorm:"column:max_double_value" json:"max_double_value"`
	MinIntValue      *int64          `gorm:"column:min_int_value" json:"min_int_value"`
	MaxIntValue      *int64          `gorm:"column:max_int_value" json:"max_int_value"`
}

// TableName ArtifactPropertyOption's table name
func (*ArtifactPropertyOption) TableName() string {
	return TableNameArtifactPropertyOption
}
