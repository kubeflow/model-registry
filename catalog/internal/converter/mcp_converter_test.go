package converter_test

import (
	"testing"

	"github.com/kubeflow/hub/catalog/internal/catalog/mcpcatalog/models"
	"github.com/kubeflow/hub/catalog/internal/converter"
	"github.com/kubeflow/hub/internal/apiutils"
	dbmodels "github.com/kubeflow/hub/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertDbMCPToolToOpenapi_StripsQualifiedPrefix(t *testing.T) {
	tests := []struct {
		name         string
		storedName   string
		expectedName string
	}{
		{
			name:         "strips server@version prefix",
			storedName:   "weather-api@1.0.0:get_current_weather",
			expectedName: "get_current_weather",
		},
		{
			name:         "strips server-only prefix (no version)",
			storedName:   "myserver:my_tool",
			expectedName: "my_tool",
		},
		{
			name:         "no prefix passes through unchanged",
			storedName:   "plain_tool_name",
			expectedName: "plain_tool_name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			accessType := "read_only"
			tool := &models.MCPServerToolImpl{
				Attributes: &models.MCPServerToolAttributes{
					Name: apiutils.Of(tc.storedName),
				},
				Properties: &[]dbmodels.Properties{
					{Name: "accessType", StringValue: apiutils.Of(accessType)},
				},
			}

			result := converter.ConvertDbMCPToolToOpenapi(tool)
			require.NotNil(t, result)
			assert.Equal(t, tc.expectedName, result.Name)
		})
	}
}
