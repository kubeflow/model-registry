package api

import (
	"encoding/json"
	"github.com/kubeflow/model-registry/ui/bff/config"
	"github.com/kubeflow/model-registry/ui/bff/data"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {

	app := App{config: config.EnvConfig{
		Port: 4000,
	}}

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, HealthCheckPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	app.HealthcheckHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal("Failed to read response body")
	}

	var healthCheckRes data.HealthCheckModel
	err = json.Unmarshal(body, &healthCheckRes)
	if err != nil {
		t.Fatalf("Error unmarshalling response JSON: %v", err)
	}

	expected := data.HealthCheckModel{
		Status: "available",
		SystemInfo: data.SystemInfo{
			Version: Version,
		},
	}

	assert.Equal(t, expected, healthCheckRes)
}
