package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadinessHandler_NonEmbedMD(t *testing.T) {
	ds := datastore.Datastore{
		Type: "mlmd",
	}
	handler := ReadinessHandler(ds)

	req, err := http.NewRequest("GET", "/readyz/isDirty", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "ok", rr.Body.String())
}
