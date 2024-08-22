package api

import (
	"encoding/json"
	"github.com/kubeflow/model-registry/ui/bff/data"
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModelRegistryHandler(t *testing.T) {
	mockK8sClient, _ := mocks.NewKubernetesClient(nil)

	testApp := App{
		kubernetesClient: mockK8sClient,
	}

	req, err := http.NewRequest(http.MethodGet, ModelRegistry, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	testApp.ModelRegistryHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var modelRegistryRes Envelope
	err = json.Unmarshal(body, &modelRegistryRes)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Convert the unmarshalled data to the expected type
	actualModelRegistry := make([]data.ModelRegistryModel, 0)
	for _, v := range modelRegistryRes["model_registry"].([]interface{}) {
		model := v.(map[string]interface{})
		actualModelRegistry = append(actualModelRegistry, data.ModelRegistryModel{Name: model["name"].(string), Description: model["description"].(string), DisplayName: model["displayName"].(string)})
	}
	modelRegistryRes["model_registry"] = actualModelRegistry

	var expected = Envelope{
		"model_registry": []data.ModelRegistryModel{
			{Name: "model-registry", Description: "Model registry description", DisplayName: "Model Registry"},
			{Name: "model-registry-dora", Description: "Model registry dora description", DisplayName: "Model Registry Dora"},
			{Name: "model-registry-bella", Description: "Model registry bella description", DisplayName: "Model Registry Bella"},
		},
	}

	assert.Equal(t, expected, modelRegistryRes)

}
