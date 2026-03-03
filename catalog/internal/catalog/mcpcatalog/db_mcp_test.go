package mcpcatalog

import (
	"context"
	"errors"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/mcpcatalog/models"
	internalmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mergeFilterQueries tests ---

func TestMergeFilterQueries(t *testing.T) {
	t.Run("both empty returns empty string", func(t *testing.T) {
		result := mergeFilterQueries("", "")
		assert.Equal(t, "", result)
	})

	t.Run("only first non-empty returns first", func(t *testing.T) {
		result := mergeFilterQueries("a = 1", "")
		assert.Equal(t, "a = 1", result)
	})

	t.Run("only second non-empty returns second", func(t *testing.T) {
		result := mergeFilterQueries("", "b = 2")
		assert.Equal(t, "b = 2", result)
	})

	t.Run("both non-empty wraps each in parens and joins with AND", func(t *testing.T) {
		result := mergeFilterQueries("a = 1", "b = 2")
		assert.Equal(t, "(a = 1) AND (b = 2)", result)
	})
}

// --- mock repository ---

// mockMCPServerRepo is a minimal MCPServerRepository for unit testing.
type mockMCPServerRepo struct {
	listResult *internalmodels.ListWrapper[models.MCPServer]
	listErr    error
	// capturedOptions stores the last MCPServerListOptions passed to List.
	capturedOptions models.MCPServerListOptions
}

func (m *mockMCPServerRepo) List(opts models.MCPServerListOptions) (*internalmodels.ListWrapper[models.MCPServer], error) {
	m.capturedOptions = opts
	return m.listResult, m.listErr
}

func (m *mockMCPServerRepo) GetByID(_ int32) (models.MCPServer, error) {
	return nil, errors.New("not implemented")
}

func (m *mockMCPServerRepo) GetByNameAndVersion(_ string, _ string) (models.MCPServer, error) {
	return nil, errors.New("not implemented")
}

func (m *mockMCPServerRepo) Save(_ models.MCPServer) (models.MCPServer, error) {
	return nil, errors.New("not implemented")
}

func (m *mockMCPServerRepo) DeleteBySource(_ string) error { return errors.New("not implemented") }
func (m *mockMCPServerRepo) DeleteByID(_ int32) error      { return errors.New("not implemented") }
func (m *mockMCPServerRepo) GetDistinctSourceIDs() ([]string, error) {
	return nil, errors.New("not implemented")
}

// mockMCPServerToolRepo is a minimal MCPServerToolRepository that always returns empty.
type mockMCPServerToolRepo struct{}

func (m *mockMCPServerToolRepo) List(_ models.MCPServerToolListOptions) ([]models.MCPServerTool, error) {
	return nil, nil
}
func (m *mockMCPServerToolRepo) GetByID(_ int32) (models.MCPServerTool, error) {
	return nil, errors.New("not implemented")
}
func (m *mockMCPServerToolRepo) Save(_ models.MCPServerTool, _ *int32) (models.MCPServerTool, error) {
	return nil, errors.New("not implemented")
}
func (m *mockMCPServerToolRepo) DeleteByParentID(_ int32) error { return errors.New("not implemented") }
func (m *mockMCPServerToolRepo) DeleteByID(_ int32) error       { return errors.New("not implemented") }

// --- ListMCPServers named query tests ---

func newTestCatalog(repo *mockMCPServerRepo, resolver NamedQueryResolver) *dbMCPCatalogImpl {
	return &dbMCPCatalogImpl{
		mcpServerRepo:     repo,
		mcpServerToolRepo: &mockMCPServerToolRepo{},
		resolveNamedQuery: resolver,
	}
}

func emptyList() *internalmodels.ListWrapper[models.MCPServer] {
	return &internalmodels.ListWrapper[models.MCPServer]{Items: []models.MCPServer{}}
}

func TestListMCPServers_NoNamedQuery(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: emptyList()}
	cat := newTestCatalog(repo, nil)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		FilterQuery: "provider = 'OpenAI'",
	})
	require.NoError(t, err)
	assert.Equal(t, "provider = 'OpenAI'", *repo.capturedOptions.FilterQuery)
}

func TestListMCPServers_NamedQueryResolved(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: emptyList()}
	resolver := func(name string) (map[string]basecatalog.FieldFilter, bool) {
		if name == "verified_only" {
			return map[string]basecatalog.FieldFilter{
				"verified": {Operator: "=", Value: true},
			}, true
		}
		return nil, false
	}
	cat := newTestCatalog(repo, resolver)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		NamedQuery: "verified_only",
	})
	require.NoError(t, err)
	assert.Equal(t, "verified = true", *repo.capturedOptions.FilterQuery)
}

func TestListMCPServers_NamedQueryMergedWithFilterQuery(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: emptyList()}
	resolver := func(name string) (map[string]basecatalog.FieldFilter, bool) {
		if name == "active" {
			return map[string]basecatalog.FieldFilter{
				"status": {Operator: "=", Value: "active"},
			}, true
		}
		return nil, false
	}
	cat := newTestCatalog(repo, resolver)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		NamedQuery:  "active",
		FilterQuery: "provider = 'OpenAI'",
	})
	require.NoError(t, err)
	// user filterQuery comes first, named query resolved second, both wrapped in parens
	assert.Equal(t, "(provider = 'OpenAI') AND (status = 'active')", *repo.capturedOptions.FilterQuery)
}

func TestListMCPServers_NoResolverReturnsError(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: emptyList()}
	cat := newTestCatalog(repo, nil)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		NamedQuery: "any_query",
	})
	require.Error(t, err)
	assert.True(t, errors.Is(err, api.ErrBadRequest))
}

func TestListMCPServers_UnknownNamedQueryReturnsError(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: emptyList()}
	resolver := func(_ string) (map[string]basecatalog.FieldFilter, bool) {
		return nil, false
	}
	cat := newTestCatalog(repo, resolver)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		NamedQuery: "nonexistent",
	})
	require.Error(t, err)
	assert.True(t, errors.Is(err, api.ErrBadRequest))
}
