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

	req, err := http.NewRequest(http.MethodGet, ModelRegistryListPath, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	testApp.ModelRegistryHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var actual ModelRegistryListEnvelope
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	var expected = ModelRegistryListEnvelope{
		Data: []data.ModelRegistryModel{
			{Name: "model-registry", Description: "Model registry description", DisplayName: "Model Registry"},
			{Name: "model-registry-dora", Description: "Model registry dora description", DisplayName: "Model Registry Dora"},
			{Name: "model-registry-bella", Description: "Model registry bella description", DisplayName: "Model Registry Bella"},
		},
	}

	assert.Equal(t, expected, actual)

}
