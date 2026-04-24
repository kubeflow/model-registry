package basecatalog

import (
	"sort"

	sharedmodels "github.com/kubeflow/hub/catalog/internal/db/models"
	apimodels "github.com/kubeflow/hub/catalog/pkg/openapi"
	"github.com/kubeflow/hub/internal/apiutils"
)

// DbPropToAPIOption converts a database PropertyOption to an API FilterOption.
// Returns nil if the property has no values (empty set or nil range).
func DbPropToAPIOption(prop sharedmodels.PropertyOption) *apimodels.FilterOption {
	var option apimodels.FilterOption

	switch prop.ValueField() {
	case sharedmodels.StringValueField:
		if len(prop.StringValue) == 0 {
			return nil
		}
		option.Type = "string"
		sort.Strings(prop.StringValue)
		option.Values = AnySlice(prop.StringValue)

	case sharedmodels.ArrayValueField:
		if len(prop.ArrayValue) == 0 {
			return nil
		}
		option.Type = "string"
		sort.Strings(prop.ArrayValue)
		option.Values = AnySlice(prop.ArrayValue)

	case sharedmodels.IntValueField:
		if prop.MinIntValue == nil || prop.MaxIntValue == nil {
			return nil
		}

		option.Type = "number"
		option.Range = &apimodels.FilterOptionRange{
			Min: apiutils.Of(float64(*prop.MinIntValue)),
			Max: apiutils.Of(float64(*prop.MaxIntValue)),
		}

	case sharedmodels.DoubleValueField:
		if prop.MinDoubleValue == nil || prop.MaxDoubleValue == nil {
			return nil
		}

		option.Type = "number"
		option.Range = &apimodels.FilterOptionRange{
			Min: prop.MinDoubleValue,
			Max: prop.MaxDoubleValue,
		}
	}

	return &option
}

// AnySlice converts a typed slice to []any.
func AnySlice[T any](s []T) []any {
	as := make([]any, len(s))
	for i, v := range s {
		as[i] = v
	}
	return as
}

// ConvertNamedQueries converts internal named queries to the API representation,
// resolving "min"/"max" sentinel values against the provided filter options.
// Returns nil if queries is empty.
func ConvertNamedQueries(
	queries map[string]map[string]FieldFilter,
	options map[string]apimodels.FilterOption,
) *map[string]map[string]apimodels.FieldFilter {
	if len(queries) == 0 {
		return nil
	}

	apiNamedQueries := make(map[string]map[string]apimodels.FieldFilter, len(queries))
	for queryName, fieldFilters := range queries {
		apiFieldFilters := make(map[string]apimodels.FieldFilter, len(fieldFilters))
		for fieldName, filter := range fieldFilters {
			apiFieldFilters[fieldName] = apimodels.FieldFilter{
				Operator: filter.Operator,
				Value:    filter.Value,
			}
		}
		ApplyMinMax(apiFieldFilters, options)
		apiNamedQueries[queryName] = apiFieldFilters
	}

	return &apiNamedQueries
}

// ApplyMinMax resolves "min"/"max" sentinel string values in a named query
// to the actual min/max from the corresponding filter options.
func ApplyMinMax(query map[string]apimodels.FieldFilter, options map[string]apimodels.FilterOption) {
	for key, filter := range query {
		value, ok := filter.Value.(string)
		if !ok || (value != "min" && value != "max") {
			continue
		}

		option, ok := options[key]
		if !ok || option.Range == nil {
			continue
		}

		switch value {
		case "min":
			if option.Range.Min != nil {
				filter.Value = *option.Range.Min
			}
		case "max":
			if option.Range.Max != nil {
				filter.Value = *option.Range.Max
			}
		}

		query[key] = filter
	}
}
