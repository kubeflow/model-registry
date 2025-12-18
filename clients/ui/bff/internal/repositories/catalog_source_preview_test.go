package repositories

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCatalogSourcePreview_Success(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}

	// Mock successful response
	responseJSON := `{"size": 10, "pageSize": 10, "nextPageToken": "", "items": []}`
	mockClient.On("POSTWithContentType", mock.Anything, mock.Anything, mock.Anything).
		Return([]byte(responseJSON), nil)

	repo := CatalogSourcePreview{}
	payload := models.CatalogSourcePreviewRequest{
		Type:           "yaml",
		IncludedModels: []string{},
		ExcludedModels: []string{},
		Properties: map[string]interface{}{
			"yaml": "models:\n  - name: test",
		},
	}

	result, err := repo.CreateCatalogSourcePreview(mockClient, payload, url.Values{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int32(10), result.Size)
	mockClient.AssertExpectations(t)
}

func TestCreateCatalogSourcePreview_Unauthorized(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}

	// Mock 401 error
	httpErr := &httpclient.HTTPError{
		StatusCode: http.StatusUnauthorized,
		ErrorResponse: httpclient.ErrorResponse{
			Code:    "401",
			Message: "Invalid or missing API key",
		},
	}
	mockClient.On("POSTWithContentType", mock.Anything, mock.Anything, mock.Anything).
		Return([]byte{}, httpErr)

	repo := CatalogSourcePreview{}
	payload := models.CatalogSourcePreviewRequest{
		Type:           "huggingface",
		IncludedModels: []string{"*"},
		ExcludedModels: []string{},
		Properties: map[string]interface{}{
			"allowedOrganization": "test-org",
		},
	}

	result, err := repo.CreateCatalogSourcePreview(mockClient, payload, url.Values{})

	assert.Nil(t, result)
	assert.Error(t, err)

	// Verify it's an HTTPError with correct status
	var httpError *httpclient.HTTPError
	assert.ErrorAs(t, err, &httpError)
	assert.Equal(t, http.StatusUnauthorized, httpError.StatusCode)
	mockClient.AssertExpectations(t)
}

func TestCreateCatalogSourcePreview_Forbidden(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}

	// Mock 403 error
	httpErr := &httpclient.HTTPError{
		StatusCode: http.StatusForbidden,
		ErrorResponse: httpclient.ErrorResponse{
			Code:    "403",
			Message: "Access denied to organization",
		},
	}
	mockClient.On("POSTWithContentType", mock.Anything, mock.Anything, mock.Anything).
		Return([]byte{}, httpErr)

	repo := CatalogSourcePreview{}
	payload := models.CatalogSourcePreviewRequest{
		Type:           "huggingface",
		IncludedModels: []string{"*"},
		ExcludedModels: []string{},
		Properties: map[string]interface{}{
			"allowedOrganization": "forbidden-org",
		},
	}

	result, err := repo.CreateCatalogSourcePreview(mockClient, payload, url.Values{})

	assert.Nil(t, result)
	assert.Error(t, err)

	// Verify it's an HTTPError with correct status
	var httpError *httpclient.HTTPError
	assert.ErrorAs(t, err, &httpError)
	assert.Equal(t, http.StatusForbidden, httpError.StatusCode)
	mockClient.AssertExpectations(t)
}

func TestCreateCatalogSourcePreview_BadRequest(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}

	// Mock 400 error
	httpErr := &httpclient.HTTPError{
		StatusCode: http.StatusBadRequest,
		ErrorResponse: httpclient.ErrorResponse{
			Code:    "400",
			Message: "Invalid source type",
		},
	}
	mockClient.On("POSTWithContentType", mock.Anything, mock.Anything, mock.Anything).
		Return([]byte{}, httpErr)

	repo := CatalogSourcePreview{}
	payload := models.CatalogSourcePreviewRequest{
		Type:           "invalid_type",
		IncludedModels: []string{},
		ExcludedModels: []string{},
		Properties:     map[string]interface{}{},
	}

	result, err := repo.CreateCatalogSourcePreview(mockClient, payload, url.Values{})

	assert.Nil(t, result)
	assert.Error(t, err)

	var httpError *httpclient.HTTPError
	assert.ErrorAs(t, err, &httpError)
	assert.Equal(t, http.StatusBadRequest, httpError.StatusCode)
	mockClient.AssertExpectations(t)
}
