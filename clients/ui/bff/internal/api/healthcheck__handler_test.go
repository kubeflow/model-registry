package api

import (
	"encoding/json"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {

	mockMRClient, _ := mocks.NewModelRegistryClient(nil)

	app := App{config: config.EnvConfig{
		Port: 4000,
	},
		repositories: repositories.NewRepositories(mockMRClient),
	}

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, HealthCheckPath, nil)
	assert.NoError(t, err)

	req.Header.Set(kubeflowUserId, mocks.KubeflowUserIDHeaderValue)

	app.HealthcheckHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)

	var healthCheckRes models.HealthCheckModel
	err = json.Unmarshal(body, &healthCheckRes)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := models.HealthCheckModel{
		Status: "available",
		SystemInfo: models.SystemInfo{
			Version: Version,
		},
		UserID: mocks.KubeflowUserIDHeaderValue,
	}

	assert.Equal(t, expected, healthCheckRes)
}
