package openapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildCombinedFilterQuery(t *testing.T) {
	testCases := []struct {
		name        string
		filterQuery string
		nameParam   string
		externalID  string
		tablePrefix string
		expected    string
	}{
		{
			name:        "empty parameters",
			filterQuery: "",
			nameParam:   "",
			externalID:  "",
			tablePrefix: "Artifact",
			expected:    "",
		},
		{
			name:        "only filter query",
			filterQuery: "state = 'LIVE'",
			nameParam:   "",
			externalID:  "",
			tablePrefix: "Artifact",
			expected:    "(state = 'LIVE')",
		},
		{
			name:        "only name parameter",
			filterQuery: "",
			nameParam:   "my-artifact",
			externalID:  "",
			tablePrefix: "Artifact",
			expected:    "Artifact.name = 'my-artifact'",
		},
		{
			name:        "only externalID parameter",
			filterQuery: "",
			nameParam:   "",
			externalID:  "ext-123",
			tablePrefix: "Artifact",
			expected:    "Artifact.external_id = 'ext-123'",
		},
		{
			name:        "name takes precedence over externalID (MLMD behavior)",
			filterQuery: "",
			nameParam:   "my-artifact",
			externalID:  "ext-123",
			tablePrefix: "Artifact",
			expected:    "Artifact.name = 'my-artifact'",
		},
		{
			name:        "filter query with name",
			filterQuery: "state = 'LIVE'",
			nameParam:   "my-artifact",
			externalID:  "",
			tablePrefix: "Artifact",
			expected:    "(state = 'LIVE') AND Artifact.name = 'my-artifact'",
		},
		{
			name:        "filter query with externalID",
			filterQuery: "state = 'LIVE'",
			nameParam:   "",
			externalID:  "ext-123",
			tablePrefix: "Artifact",
			expected:    "(state = 'LIVE') AND Artifact.external_id = 'ext-123'",
		},
		{
			name:        "filter query with both name and externalID (name wins)",
			filterQuery: "state = 'LIVE'",
			nameParam:   "my-artifact",
			externalID:  "ext-123",
			tablePrefix: "Artifact",
			expected:    "(state = 'LIVE') AND Artifact.name = 'my-artifact'",
		},
		{
			name:        "name with single quotes (escaped)",
			filterQuery: "",
			nameParam:   "my-artifact's-name",
			externalID:  "",
			tablePrefix: "Artifact",
			expected:    "Artifact.name = 'my-artifact''s-name'",
		},
		{
			name:        "externalID with single quotes (escaped)",
			filterQuery: "",
			nameParam:   "",
			externalID:  "ext-'123'",
			tablePrefix: "Artifact",
			expected:    "Artifact.external_id = 'ext-''123'''",
		},
		{
			name:        "Context table prefix with name",
			filterQuery: "",
			nameParam:   "my-model-version",
			externalID:  "",
			tablePrefix: "Context",
			expected:    "Context.name = 'my-model-version'",
		},
		{
			name:        "Context table prefix with externalID",
			filterQuery: "",
			nameParam:   "",
			externalID:  "ctx-123",
			tablePrefix: "Context",
			expected:    "Context.external_id = 'ctx-123'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := buildCombinedFilterQuery(tc.filterQuery, tc.nameParam, tc.externalID, tc.tablePrefix)
			assert.Equal(t, tc.expected, result)
		})
	}
}
