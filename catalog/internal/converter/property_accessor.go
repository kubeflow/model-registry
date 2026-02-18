package converter

import (
	"encoding/json"

	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
)

// PropertyAccessor provides O(1) access to database properties via map lookup.
// This improves performance when multiple properties need to be accessed from
// the same properties slice, avoiding O(n²) complexity from repeated linear scans.
//
// The accessor is safe for concurrent read operations after construction.
type PropertyAccessor struct {
	propMap map[string]dbmodels.Properties
}

// NewPropertyAccessor creates a property accessor from a properties slice.
// Returns a nil-safe accessor even if props is nil.
//
// Time complexity: O(n) where n is the number of properties.
// Subsequent lookups via Get* methods are O(1).
func NewPropertyAccessor(props *[]dbmodels.Properties) *PropertyAccessor {
	propMap := make(map[string]dbmodels.Properties)
	if props != nil {
		for _, prop := range *props {
			propMap[prop.Name] = prop
		}
	}
	return &PropertyAccessor{propMap: propMap}
}

// GetString retrieves a string property value.
// Returns empty string if the property doesn't exist or has no value.
func (pa *PropertyAccessor) GetString(name string) string {
	if prop, exists := pa.propMap[name]; exists && prop.StringValue != nil {
		return *prop.StringValue
	}
	return ""
}

// GetStringPtr retrieves a string property as a pointer.
// Returns nil if the property doesn't exist, has no value, or is an empty string.
func (pa *PropertyAccessor) GetStringPtr(name string) *string {
	val := pa.GetString(name)
	if val == "" {
		return nil
	}
	return &val
}

// GetBoolPtr retrieves a boolean property as a pointer.
// Returns nil if the property doesn't exist or has no value.
func (pa *PropertyAccessor) GetBoolPtr(name string) *bool {
	if prop, exists := pa.propMap[name]; exists && prop.BoolValue != nil {
		return prop.BoolValue
	}
	return nil
}

// GetInt retrieves an int64 property value (converts from int32 storage).
// Returns 0 if the property doesn't exist or has no value.
func (pa *PropertyAccessor) GetInt(name string) int64 {
	if prop, exists := pa.propMap[name]; exists && prop.IntValue != nil {
		return int64(*prop.IntValue)
	}
	return 0
}

// GetStringArray retrieves a JSON-encoded string array property.
// Returns nil if the property doesn't exist, is empty, or cannot be unmarshaled.
func (pa *PropertyAccessor) GetStringArray(name string) []string {
	jsonStr := pa.GetString(name)
	if jsonStr == "" {
		return nil
	}
	var result []string
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		// Silently return nil on unmarshal error (consistent with existing behavior)
		return nil
	}
	return result
}

// HasAny returns true if any of the specified property names exist in the property map.
// This is useful for checking if a group of optional properties has at least one value.
func (pa *PropertyAccessor) HasAny(names ...string) bool {
	for _, name := range names {
		if _, exists := pa.propMap[name]; exists {
			return true
		}
	}
	return false
}
