package api

import (
	"encoding/json"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/data"
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModelRegistryHandler(t *testing.T) {
	mockK8sClient := new(mocks.KubernetesClientMock)
	mockK8sClient.On("FetchServiceNamesByComponent", "model-registry-server").Return([]string{"model-registry-dora", "model-registry-bella"}, nil)

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
	fmt.Println(string(body))
	var modelRegistryRes Envelope
	err = json.Unmarshal(body, &modelRegistryRes)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Convert the unmarshalled data to the expected type
	actualModelRegistry := make([]data.ModelRegistryModel, 0)
	for _, v := range modelRegistryRes["model_registry"].([]interface{}) {
		model := v.(map[string]interface{})
		actualModelRegistry = append(actualModelRegistry, data.ModelRegistryModel{Name: model["name"].(string)})
	}
	modelRegistryRes["model_registry"] = actualModelRegistry

	var expected = Envelope{
		"model_registry": []data.ModelRegistryModel{
			{Name: "model-registry-dora"},
			{Name: "model-registry-bella"},
		},
	}

	assert.Equal(t, expected, modelRegistryRes)

	mockK8sClient.AssertExpectations(t)
}
