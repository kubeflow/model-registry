package service

import (
	"fmt"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/schema"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var _ models.PropertyOptionsRepository = (*PropertyOptionsRepositoryImpl)(nil)

type PropertyOptionsRepositoryImpl struct {
	db *gorm.DB
}

func NewPropertyOptionsRepository(db *gorm.DB) models.PropertyOptionsRepository {
	return &PropertyOptionsRepositoryImpl{
		db: db,
	}
}

func (r *PropertyOptionsRepositoryImpl) Refresh(t models.PropertyOptionType) error {
	if r.db.Name() != "postgres" {
		return nil
	}

	var viewName string
	switch t {
	case models.ContextPropertyOptionType:
		viewName = schema.TableNameContextPropertyOption
	case models.ArtifactPropertyOptionType:
		viewName = schema.TableNameArtifactPropertyOption
	default:
		return fmt.Errorf("invalid property option type: %d", t)
	}

	sql := fmt.Sprintf("REFRESH MATERIALIZED VIEW %s", viewName)
	if err := r.db.Exec(sql).Error; err != nil {
		return fmt.Errorf("error refreshing materialized view %s: %w", viewName, err)
	}

	return nil
}

func (r *PropertyOptionsRepositoryImpl) List(t models.PropertyOptionType, typeID int32) ([]models.PropertyOption, error) {
	if r.db.Name() != "postgres" {
		return []models.PropertyOption{}, nil
	}

	switch t {
	case models.ContextPropertyOptionType:
		return r.listContextPropertyOptions(typeID)
	case models.ArtifactPropertyOptionType:
		return r.listArtifactPropertyOptions(typeID)
	default:
		return nil, fmt.Errorf("invalid property option type: %d", t)
	}
}

func (r *PropertyOptionsRepositoryImpl) listContextPropertyOptions(typeID int32) ([]models.PropertyOption, error) {
	q := r.db
	if typeID > 0 {
		q = q.Where("type_id = ?", typeID)
	}
	q = q.Order("name")

	var contextOptions []schema.ContextPropertyOption
	if err := q.Find(&contextOptions).Error; err != nil {
		return nil, fmt.Errorf("error querying context property options: %w", err)
	}

	return convertSchemaToPropertyOptions(contextOptions), nil
}

func (r *PropertyOptionsRepositoryImpl) listArtifactPropertyOptions(typeID int32) ([]models.PropertyOption, error) {
	q := r.db
	if typeID > 0 {
		q = q.Where("type_id = ?", typeID)
	}
	q = q.Order("name")

	var artifactOptions []schema.ArtifactPropertyOption
	if err := q.Find(&artifactOptions).Error; err != nil {
		return nil, fmt.Errorf("error querying artifact property options: %w", err)
	}

	return convertSchemaToPropertyOptions(artifactOptions), nil
}

// Helper function to convert schema types to models.PropertyOption
// This works for both ContextPropertyOption and ArtifactPropertyOption since they have identical structure
func convertSchemaToPropertyOptions[T interface {
	schema.ContextPropertyOption | schema.ArtifactPropertyOption
}](options []T) []models.PropertyOption {
	result := make([]models.PropertyOption, len(options))

	for i, option := range options {
		var stringValue, arrayValue []string

		// Convert pq.StringArray to []string for StringValue
		if stringVal := getStringValue(option); stringVal != nil {
			stringValue = []string(*stringVal)
		}

		// Convert pq.StringArray to []string for ArrayValue
		if arrVal := getArrayValue(option); arrVal != nil {
			arrayValue = []string(*arrVal)
		}

		result[i] = models.PropertyOption{
			TypeID:           getTypeID(option),
			Name:             getName(option),
			IsCustomProperty: getIsCustomProperty(option),
			StringValue:      stringValue,
			ArrayValue:       arrayValue,
			MinDoubleValue:   getMinDoubleValue(option),
			MaxDoubleValue:   getMaxDoubleValue(option),
			MinIntValue:      getMinIntValue(option),
			MaxIntValue:      getMaxIntValue(option),
		}
	}

	return result
}

// Helper functions to extract fields from schema types
func getStringValue[T interface {
	schema.ContextPropertyOption | schema.ArtifactPropertyOption
}](option T) *pq.StringArray {
	switch v := any(option).(type) {
	case schema.ContextPropertyOption:
		return v.StringValue
	case schema.ArtifactPropertyOption:
		return v.StringValue
	default:
		return nil
	}
}

func getArrayValue[T interface {
	schema.ContextPropertyOption | schema.ArtifactPropertyOption
}](option T) *pq.StringArray {
	switch v := any(option).(type) {
	case schema.ContextPropertyOption:
		return v.ArrayValue
	case schema.ArtifactPropertyOption:
		return v.ArrayValue
	default:
		return nil
	}
}

func getTypeID[T interface {
	schema.ContextPropertyOption | schema.ArtifactPropertyOption
}](option T) int32 {
	switch v := any(option).(type) {
	case schema.ContextPropertyOption:
		return v.TypeID
	case schema.ArtifactPropertyOption:
		return v.TypeID
	default:
		return 0
	}
}

func getName[T interface {
	schema.ContextPropertyOption | schema.ArtifactPropertyOption
}](option T) string {
	switch v := any(option).(type) {
	case schema.ContextPropertyOption:
		return v.Name
	case schema.ArtifactPropertyOption:
		return v.Name
	default:
		return ""
	}
}

func getIsCustomProperty[T interface {
	schema.ContextPropertyOption | schema.ArtifactPropertyOption
}](option T) bool {
	switch v := any(option).(type) {
	case schema.ContextPropertyOption:
		return v.IsCustomProperty
	case schema.ArtifactPropertyOption:
		return v.IsCustomProperty
	default:
		return false
	}
}

func getMinDoubleValue[T interface {
	schema.ContextPropertyOption | schema.ArtifactPropertyOption
}](option T) *float64 {
	switch v := any(option).(type) {
	case schema.ContextPropertyOption:
		return v.MinDoubleValue
	case schema.ArtifactPropertyOption:
		return v.MinDoubleValue
	default:
		return nil
	}
}

func getMaxDoubleValue[T interface {
	schema.ContextPropertyOption | schema.ArtifactPropertyOption
}](option T) *float64 {
	switch v := any(option).(type) {
	case schema.ContextPropertyOption:
		return v.MaxDoubleValue
	case schema.ArtifactPropertyOption:
		return v.MaxDoubleValue
	default:
		return nil
	}
}

func getMinIntValue[T interface {
	schema.ContextPropertyOption | schema.ArtifactPropertyOption
}](option T) *int64 {
	switch v := any(option).(type) {
	case schema.ContextPropertyOption:
		return v.MinIntValue
	case schema.ArtifactPropertyOption:
		return v.MinIntValue
	default:
		return nil
	}
}

func getMaxIntValue[T interface {
	schema.ContextPropertyOption | schema.ArtifactPropertyOption
}](option T) *int64 {
	switch v := any(option).(type) {
	case schema.ContextPropertyOption:
		return v.MaxIntValue
	case schema.ArtifactPropertyOption:
		return v.MaxIntValue
	default:
		return nil
	}
}
