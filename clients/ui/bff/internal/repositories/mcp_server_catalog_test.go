package repositories

import (
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllMcpServers_Success(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}
	responseJSON := `{"size": 2, "pageSize": 10, "nextPageToken": "", "items": [{"id": 1, "name": "Server 1", "toolCount": 5}, {"id": 2, "name": "Server 2", "toolCount": 3}]}`
	mockClient.On("GET", "/mcp_servers").Return([]byte(responseJSON), nil)

	repo := CatalogSources{}
	pageValues := url.Values{}
	result, err := repo.GetAllMcpServers(mockClient, pageValues)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int32(2), result.Size)
	assert.Equal(t, int32(10), result.PageSize)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, 1, result.Items[0].ID)
	assert.Equal(t, "Server 1", result.Items[0].Name)
	mockClient.AssertExpectations(t)
}

func TestGetAllMcpServers_WithQueryParams(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}
	responseJSON := `{"size": 1, "pageSize": 10, "nextPageToken": "", "items": [{"id": 1, "name": "MCP Server", "toolCount": 2}]}`
	mockClient.On("GET", mock.MatchedBy(func(path string) bool {
		return strings.HasPrefix(path, "/mcp_servers?") &&
			strings.Contains(path, "name=test")
	})).Return([]byte(responseJSON), nil)

	repo := CatalogSources{}
	pageValues := url.Values{}
	pageValues.Set("name", "test")
	result, err := repo.GetAllMcpServers(mockClient, pageValues)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "MCP Server", result.Items[0].Name)
	mockClient.AssertExpectations(t)
}

func TestGetAllMcpServers_ClientError(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}
	expectedErr := &httpclient.HTTPError{
		StatusCode: 500,
		ErrorResponse: httpclient.ErrorResponse{
			Code:    "500",
			Message: "internal server error",
		},
	}
	mockClient.On("GET", "/mcp_servers").Return([]byte(nil), expectedErr)

	repo := CatalogSources{}
	result, err := repo.GetAllMcpServers(mockClient, url.Values{})

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	mockClient.AssertExpectations(t)
}

func TestGetAllMcpServers_InvalidJSON(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}
	mockClient.On("GET", "/mcp_servers").Return([]byte("not valid json"), nil)

	repo := CatalogSources{}
	result, err := repo.GetAllMcpServers(mockClient, url.Values{})

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decoding")
	mockClient.AssertExpectations(t)
}

func TestGetMcpServersFilter_Success(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}
	responseJSON := `{"filters": {}, "namedQueries": {}}`
	mockClient.On("GET", "/mcp_servers/filter_options").Return([]byte(responseJSON), nil)

	repo := CatalogSources{}
	result, err := repo.GetMcpServersFilter(mockClient)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockClient.AssertExpectations(t)
}

func TestGetMcpServersFilter_ClientError(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}
	expectedErr := errors.New("network failure")
	mockClient.On("GET", "/mcp_servers/filter_options").Return([]byte(nil), expectedErr)

	repo := CatalogSources{}
	result, err := repo.GetMcpServersFilter(mockClient)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fetching")
	mockClient.AssertExpectations(t)
}

func TestGetMcpServer_Success(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}
	responseJSON := `{"id": 1, "name": "Test MCP Server", "toolCount": 5}`
	mockClient.On("GET", "/mcp_servers/server-1").Return([]byte(responseJSON), nil)

	repo := CatalogSources{}
	result, err := repo.GetMcpServer(mockClient, "server-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Test MCP Server", result.Name)
	mockClient.AssertExpectations(t)
}

func TestGetMcpServer_ClientError(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}
	expectedErr := &httpclient.HTTPError{StatusCode: 404}
	mockClient.On("GET", "/mcp_servers/missing").Return([]byte(nil), expectedErr)

	repo := CatalogSources{}
	result, err := repo.GetMcpServer(mockClient, "missing")

	assert.Nil(t, result)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

func TestGetMcpServersTools_Success(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}
	responseJSON := `{"size": 2, "pageSize": 10, "nextPageToken": "", "items": []}`
	mockClient.On("GET", "/mcp_servers/server-1/tools").Return([]byte(responseJSON), nil)

	repo := CatalogSources{}
	result, err := repo.GetMcpServersTools(mockClient, "server-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int32(2), result.Size)
	assert.Len(t, result.Items, 0)
	mockClient.AssertExpectations(t)
}

func TestGetMcpServersTools_ClientError(t *testing.T) {
	mockClient := &mocks.MockHTTPClient{}
	expectedErr := errors.New("fetch failed")
	mockClient.On("GET", "/mcp_servers/bad/tools").Return([]byte(nil), expectedErr)

	repo := CatalogSources{}
	result, err := repo.GetMcpServersTools(mockClient, "bad")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fetching")
	mockClient.AssertExpectations(t)
}
