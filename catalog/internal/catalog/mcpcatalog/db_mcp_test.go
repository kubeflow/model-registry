package mcpcatalog

import (
	"context"
	"errors"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/mcpcatalog/models"
	sharedmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	internalmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPropertyOptionsRepo is a minimal PropertyOptionsRepository for unit testing.
type mockPropertyOptionsRepo struct {
	contextProps  []sharedmodels.PropertyOption
	artifactProps []sharedmodels.PropertyOption
	listErr       error
}

func (m *mockPropertyOptionsRepo) Refresh(_ sharedmodels.PropertyOptionType) error { return nil }
func (m *mockPropertyOptionsRepo) List(t sharedmodels.PropertyOptionType, _ int32) ([]sharedmodels.PropertyOption, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	if t == sharedmodels.ContextPropertyOptionType {
		return m.contextProps, nil
	}
	return m.artifactProps, nil
}

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
	getResult  models.MCPServer
	getErr     error
	// capturedOptions stores the last MCPServerListOptions passed to List.
	capturedOptions models.MCPServerListOptions
	TypeID          int32 // Set to a non-zero value when testing GetFilterOptions scoping
}

func (m *mockMCPServerRepo) List(opts models.MCPServerListOptions) (*internalmodels.ListWrapper[models.MCPServer], error) {
	m.capturedOptions = opts
	return m.listResult, m.listErr
}

func (m *mockMCPServerRepo) GetByID(_ int32) (models.MCPServer, error) {
	if m.getResult != nil || m.getErr != nil {
		return m.getResult, m.getErr
	}
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
func (m *mockMCPServerRepo) GetTypeID() int32 { return m.TypeID }

// mockMCPServerToolRepo is a configurable MCPServerToolRepository for unit testing.
type mockMCPServerToolRepo struct {
	listResult  *internalmodels.ListWrapper[models.MCPServerTool]
	listErr     error
	countResult int32
	countErr    error
	// countByID allows per-parent-ID count results for multi-server tests.
	// When set, CountByParentIDs uses this map instead of countResult.
	countByID map[int32]int32
	// capturedListOptions stores the last MCPServerToolListOptions passed to List.
	capturedListOptions *models.MCPServerToolListOptions
}

func (m *mockMCPServerToolRepo) List(opts models.MCPServerToolListOptions) (*internalmodels.ListWrapper[models.MCPServerTool], error) {
	m.capturedListOptions = &opts
	return m.listResult, m.listErr
}
func (m *mockMCPServerToolRepo) GetByID(_ int32) (models.MCPServerTool, error) {
	return nil, errors.New("not implemented")
}
func (m *mockMCPServerToolRepo) Save(_ models.MCPServerTool, _ *int32) (models.MCPServerTool, error) {
	return nil, errors.New("not implemented")
}
func (m *mockMCPServerToolRepo) CountByParentIDs(parentIDs []int32) (map[int32]int32, error) {
	if m.countErr != nil {
		return nil, m.countErr
	}
	result := make(map[int32]int32, len(parentIDs))
	for _, id := range parentIDs {
		if m.countByID != nil {
			result[id] = m.countByID[id]
		} else {
			result[id] = m.countResult
		}
	}
	return result, nil
}
func (m *mockMCPServerToolRepo) DeleteByParentID(_ int32) error { return errors.New("not implemented") }
func (m *mockMCPServerToolRepo) DeleteByID(_ int32) error       { return errors.New("not implemented") }

// --- helpers ---

func newTestCatalog(repo *mockMCPServerRepo, resolver NamedQueryResolver) *dbMCPCatalogImpl {
	return newTestCatalogWithToolRepo(repo, &mockMCPServerToolRepo{}, resolver)
}

func newTestCatalogWithToolRepo(repo *mockMCPServerRepo, toolRepo *mockMCPServerToolRepo, resolver NamedQueryResolver) *dbMCPCatalogImpl {
	return &dbMCPCatalogImpl{
		mcpServerRepo:     repo,
		mcpServerToolRepo: toolRepo,
		resolveNamedQuery: resolver,
	}
}

func emptyList() *internalmodels.ListWrapper[models.MCPServer] {
	return &internalmodels.ListWrapper[models.MCPServer]{Items: []models.MCPServer{}}
}

func serverID(id int32) *int32 { return &id }

func serverName(name string) *string { return &name }

func listWithServer(id int32, name string) *internalmodels.ListWrapper[models.MCPServer] {
	server := &models.MCPServerImpl{
		ID:         serverID(id),
		Attributes: &models.MCPServerAttributes{Name: serverName(name)},
	}
	return &internalmodels.ListWrapper[models.MCPServer]{Items: []models.MCPServer{server}}
}

type serverStub struct {
	id   int32
	name string
}

func listWithServers(servers ...serverStub) *internalmodels.ListWrapper[models.MCPServer] {
	items := make([]models.MCPServer, 0, len(servers))
	for _, s := range servers {
		items = append(items, &models.MCPServerImpl{
			ID:         serverID(s.id),
			Attributes: &models.MCPServerAttributes{Name: serverName(s.name)},
		})
	}
	return &internalmodels.ListWrapper[models.MCPServer]{Items: items}
}

func makeTool(name string) models.MCPServerTool {
	return &models.MCPServerToolImpl{
		Attributes: &models.MCPServerToolAttributes{Name: serverName(name)},
	}
}

func toolList(tools ...models.MCPServerTool) *internalmodels.ListWrapper[models.MCPServerTool] {
	return &internalmodels.ListWrapper[models.MCPServerTool]{Items: tools}
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

// --- GetFilterOptions tests ---

func newTestCatalogWithFilterOptions(repo *mockMCPServerRepo, propRepo *mockPropertyOptionsRepo, sources *MCPSourceCollection) *dbMCPCatalogImpl {
	return &dbMCPCatalogImpl{
		mcpServerRepo:             repo,
		mcpServerToolRepo:         &mockMCPServerToolRepo{},
		propertyOptionsRepository: propRepo,
		mcpSources:                sources,
	}
}

func TestGetFilterOptions_EmptyDB(t *testing.T) {
	repo := &mockMCPServerRepo{}
	propRepo := &mockPropertyOptionsRepo{}
	cat := newTestCatalogWithFilterOptions(repo, propRepo, nil)

	result, err := cat.GetFilterOptions(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.Filters)
	assert.Empty(t, *result.Filters)
	assert.Nil(t, result.NamedQueries)
}

func TestGetFilterOptions_FiltersFromDB(t *testing.T) {
	providerValues := []string{"OpenAI", "GitHub"}
	repo := &mockMCPServerRepo{}
	propRepo := &mockPropertyOptionsRepo{
		contextProps: []sharedmodels.PropertyOption{
			{Name: "provider", StringValue: providerValues},
			{Name: "source_id", StringValue: []string{"src1"}}, // should be skipped
			{Name: "logo", StringValue: []string{"logo.png"}},  // should be skipped
		},
	}
	cat := newTestCatalogWithFilterOptions(repo, propRepo, nil)

	result, err := cat.GetFilterOptions(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Filters)
	filters := *result.Filters
	assert.Contains(t, filters, "provider")
	assert.NotContains(t, filters, "source_id")
	assert.NotContains(t, filters, "logo")
}

func TestGetFilterOptions_NamedQueriesFromSources(t *testing.T) {
	sources := NewMCPSourceCollection()
	err := sources.MergeWithNamedQueries("test", nil, map[string]map[string]basecatalog.FieldFilter{
		"verified_only": {
			"verifiedSource": {Operator: "=", Value: true},
		},
	})
	require.NoError(t, err)

	repo := &mockMCPServerRepo{}
	propRepo := &mockPropertyOptionsRepo{}
	cat := newTestCatalogWithFilterOptions(repo, propRepo, sources)

	result, err := cat.GetFilterOptions(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.NamedQueries)
	nq := *result.NamedQueries
	require.Contains(t, nq, "verified_only")
	assert.Equal(t, true, nq["verified_only"]["verifiedSource"].Value)
}

func TestGetFilterOptions_PropertyOptionsError(t *testing.T) {
	repo := &mockMCPServerRepo{}
	propRepo := &mockPropertyOptionsRepo{listErr: errors.New("db error")}
	cat := newTestCatalogWithFilterOptions(repo, propRepo, nil)

	_, err := cat.GetFilterOptions(context.Background())
	require.Error(t, err)
}

// --- ToolCount tests ---

func TestListMCPServers_ToolCountWithoutIncludeTools(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: listWithServer(1, "test-server")}
	toolRepo := &mockMCPServerToolRepo{countResult: 5}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	result, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		IncludeTools: false,
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 1)
	assert.Equal(t, int32(5), result.Items[0].ToolCount)
	assert.Empty(t, result.Items[0].Tools)
}

func TestListMCPServers_CountByParentIDError(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: listWithServer(1, "test-server")}
	toolRepo := &mockMCPServerToolRepo{countErr: errors.New("db count error")}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		IncludeTools: false,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error counting tools")
}

func TestGetMCPServer_ToolCountWithoutIncludeTools(t *testing.T) {
	serverEntity := &models.MCPServerImpl{
		ID:         serverID(1),
		Attributes: &models.MCPServerAttributes{Name: serverName("test-server")},
	}
	repo := &mockMCPServerRepo{getResult: serverEntity}
	toolRepo := &mockMCPServerToolRepo{countResult: 7}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	result, err := cat.GetMCPServer(context.Background(), "1", false, 0)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, int32(7), result.ToolCount)
	assert.Empty(t, result.Tools)
}

func TestGetMCPServer_CountByParentIDError(t *testing.T) {
	serverEntity := &models.MCPServerImpl{
		ID:         serverID(1),
		Attributes: &models.MCPServerAttributes{Name: serverName("test-server")},
	}
	repo := &mockMCPServerRepo{getResult: serverEntity}
	toolRepo := &mockMCPServerToolRepo{countErr: errors.New("db count error")}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	_, err := cat.GetMCPServer(context.Background(), "1", false, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error counting tools")
}

func TestListMCPServers_ToolCountZero(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: listWithServer(1, "no-tools-server")}
	toolRepo := &mockMCPServerToolRepo{countResult: 0}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	result, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		IncludeTools: false,
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 1)
	assert.Equal(t, int32(0), result.Items[0].ToolCount)
	assert.Empty(t, result.Items[0].Tools)
}

func TestListMCPServers_MultipleServersWithDifferentToolCounts(t *testing.T) {
	repo := &mockMCPServerRepo{
		listResult: listWithServers(
			serverStub{1, "server-a"},
			serverStub{2, "server-b"},
			serverStub{3, "server-c"},
		),
	}
	toolRepo := &mockMCPServerToolRepo{
		countByID: map[int32]int32{1: 3, 2: 0, 3: 10},
	}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	result, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		IncludeTools: false,
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 3)
	assert.Equal(t, int32(3), result.Items[0].ToolCount)
	assert.Equal(t, int32(0), result.Items[1].ToolCount)
	assert.Equal(t, int32(10), result.Items[2].ToolCount)
}

func TestListMCPServers_IncludeToolsUsesAccurateCount(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: listWithServer(1, "test-server")}
	toolRepo := &mockMCPServerToolRepo{
		countResult: 5,
		listResult:  toolList(makeTool("tool-1"), makeTool("tool-2")),
	}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	result, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		IncludeTools: true,
		ToolLimit:    2,
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 1)
	// toolCount should reflect total (5), not len(returned tools) (2)
	assert.Equal(t, int32(5), result.Items[0].ToolCount)
	assert.Len(t, result.Items[0].Tools, 2)
}

func TestListMCPServers_ToolLimitPassedAsPageSize(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: listWithServer(1, "test-server")}
	toolRepo := &mockMCPServerToolRepo{
		countResult: 5,
		listResult:  toolList(makeTool("tool-1")),
	}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		IncludeTools: true,
		ToolLimit:    3,
	})
	require.NoError(t, err)
	require.NotNil(t, toolRepo.capturedListOptions, "List should have been called on tool repo")
	require.NotNil(t, toolRepo.capturedListOptions.PageSize, "PageSize should be set when ToolLimit > 0")
	assert.Equal(t, int32(3), *toolRepo.capturedListOptions.PageSize)
}

func TestListMCPServers_ZeroToolLimitDoesNotSetPageSize(t *testing.T) {
	repo := &mockMCPServerRepo{listResult: listWithServer(1, "test-server")}
	toolRepo := &mockMCPServerToolRepo{
		countResult: 5,
		listResult:  toolList(makeTool("tool-1"), makeTool("tool-2")),
	}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		IncludeTools: true,
		ToolLimit:    0,
	})
	require.NoError(t, err)
	require.NotNil(t, toolRepo.capturedListOptions, "List should have been called on tool repo")
	assert.Nil(t, toolRepo.capturedListOptions.PageSize, "PageSize should be nil when ToolLimit is 0")
}

func TestGetMCPServer_ToolLimitSetsPageSize(t *testing.T) {
	serverEntity := &models.MCPServerImpl{
		ID:         serverID(1),
		Attributes: &models.MCPServerAttributes{Name: serverName("test-server")},
	}
	repo := &mockMCPServerRepo{getResult: serverEntity}
	toolRepo := &mockMCPServerToolRepo{
		countResult: 5,
		listResult:  toolList(makeTool("tool-1")),
	}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	_, err := cat.GetMCPServer(context.Background(), "1", true, 3)
	require.NoError(t, err)
	require.NotNil(t, toolRepo.capturedListOptions, "List should have been called on tool repo")
	require.NotNil(t, toolRepo.capturedListOptions.PageSize, "PageSize should be set when toolLimit > 0")
	assert.Equal(t, int32(3), *toolRepo.capturedListOptions.PageSize)
}

func TestGetMCPServer_ZeroToolLimitOmitsPageSize(t *testing.T) {
	serverEntity := &models.MCPServerImpl{
		ID:         serverID(1),
		Attributes: &models.MCPServerAttributes{Name: serverName("test-server")},
	}
	repo := &mockMCPServerRepo{getResult: serverEntity}
	toolRepo := &mockMCPServerToolRepo{
		countResult: 5,
		listResult:  toolList(makeTool("tool-1"), makeTool("tool-2")),
	}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	_, err := cat.GetMCPServer(context.Background(), "1", true, 0)
	require.NoError(t, err)
	require.NotNil(t, toolRepo.capturedListOptions, "List should have been called on tool repo")
	assert.Nil(t, toolRepo.capturedListOptions.PageSize, "PageSize should be nil when toolLimit is 0")
}

func TestGetMCPServerTool_StripsQualifiedPrefix(t *testing.T) {
	serverEntity := &models.MCPServerImpl{
		ID:         serverID(88),
		Attributes: &models.MCPServerAttributes{Name: serverName("dynatrace-mcp")},
	}
	repo := &mockMCPServerRepo{getResult: serverEntity}
	toolRepo := &mockMCPServerToolRepo{
		listResult: toolList(makeTool("dynatrace-mcp@1.6.1:list_problems")),
	}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	result, err := cat.GetMCPServerTool(context.Background(), "88", "list_problems")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "list_problems", result.Name)
}

func TestGetMCPServerTool_NotFound(t *testing.T) {
	serverEntity := &models.MCPServerImpl{
		ID:         serverID(88),
		Attributes: &models.MCPServerAttributes{Name: serverName("dynatrace-mcp")},
	}
	repo := &mockMCPServerRepo{getResult: serverEntity}
	// DB-level filter returns empty when tool name doesn't match.
	toolRepo := &mockMCPServerToolRepo{
		listResult: toolList(),
	}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	_, err := cat.GetMCPServerTool(context.Background(), "88", "nonexistent_tool")
	require.Error(t, err)
	assert.ErrorIs(t, err, api.ErrNotFound)
}

func TestGetMCPServerTool_PassesToolNameFilter(t *testing.T) {
	serverEntity := &models.MCPServerImpl{
		ID:         serverID(88),
		Attributes: &models.MCPServerAttributes{Name: serverName("dynatrace-mcp")},
	}
	repo := &mockMCPServerRepo{getResult: serverEntity}
	toolRepo := &mockMCPServerToolRepo{
		listResult: toolList(makeTool("dynatrace-mcp@1.6.1:list_problems")),
	}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	_, err := cat.GetMCPServerTool(context.Background(), "88", "list_problems")
	require.NoError(t, err)
	require.NotNil(t, toolRepo.capturedListOptions, "List should have been called")
	require.NotNil(t, toolRepo.capturedListOptions.ToolName, "ToolName should be set on list options")
	assert.Equal(t, "list_problems", *toolRepo.capturedListOptions.ToolName)
}

// --- Pagination default tests ---

func TestListMCPServers_EmptyOrderByUsesDefault(t *testing.T) {
	// When OrderBy and SortOrder are not specified (empty strings from the API layer),
	// the pagination must still apply default ordering ("id" / "ASC") so that
	// cursor-based pagination works correctly on subsequent pages.
	repo := &mockMCPServerRepo{listResult: emptyList()}
	cat := newTestCatalog(repo, nil)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		FilterQuery: "license='Apache 2.0'",
		PageSize:    10,
		// OrderBy and SortOrder are zero values (empty strings)
	})
	require.NoError(t, err)

	pagination := repo.capturedOptions.Pagination
	// The pagination must resolve to the defaults ("id" / "ASC").
	// Either OrderBy is nil (so GetOrderBy falls through to DefaultOrderBy),
	// or it is set to the explicit default value — but it must NOT be a
	// pointer to an empty string, which would skip the ORDER BY clause.
	orderBy := pagination.GetOrderBy()
	sortOrder := pagination.GetSortOrder()
	assert.Equal(t, "id", orderBy, "expected default orderBy 'id' when param is empty")
	assert.Equal(t, "ASC", sortOrder, "expected default sortOrder 'ASC' when param is empty")
}

func TestListMCPServers_EmptyOrderByPaginationNotNilEmptyString(t *testing.T) {
	// Regression: when OrderBy/SortOrder are empty, they must NOT be set as
	// pointers to empty strings. The downstream paginator skips ORDER BY when
	// both are empty, but still applies the cursor WHERE clause, causing
	// non-deterministic results and empty second pages.
	repo := &mockMCPServerRepo{listResult: emptyList()}
	cat := newTestCatalog(repo, nil)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		FilterQuery: "license='Apache 2.0'",
		PageSize:    1,
		// OrderBy and SortOrder deliberately left as zero values
	})
	require.NoError(t, err)

	pagination := repo.capturedOptions.Pagination

	// If OrderBy is non-nil, it must not point to an empty string.
	if pagination.OrderBy != nil {
		assert.NotEmpty(t, *pagination.OrderBy, "OrderBy must not be a pointer to an empty string")
	}
	// If SortOrder is non-nil, it must not point to an empty string.
	if pagination.SortOrder != nil {
		assert.NotEmpty(t, *pagination.SortOrder, "SortOrder must not be a pointer to an empty string")
	}
}

func TestListMCPServers_ExplicitOrderByIsPreserved(t *testing.T) {
	// When OrderBy/SortOrder are explicitly provided, they should be passed through.
	repo := &mockMCPServerRepo{listResult: emptyList()}
	cat := newTestCatalog(repo, nil)

	_, err := cat.ListMCPServers(context.Background(), ListMCPServersParams{
		FilterQuery: "license='Apache 2.0'",
		PageSize:    10,
		OrderBy:     "CREATE_TIME",
		SortOrder:   "DESC",
	})
	require.NoError(t, err)

	pagination := repo.capturedOptions.Pagination
	assert.Equal(t, "CREATE_TIME", pagination.GetOrderBy())
	assert.Equal(t, "DESC", pagination.GetSortOrder())
}

func TestListMCPServerTools_EmptyOrderByUsesDefault(t *testing.T) {
	// Same bug applies to ListMCPServerTools — empty OrderBy/SortOrder must
	// resolve to defaults for correct cursor-based pagination.
	serverEntity := &models.MCPServerImpl{
		ID:         serverID(1),
		Attributes: &models.MCPServerAttributes{Name: serverName("test-server")},
	}
	repo := &mockMCPServerRepo{getResult: serverEntity}
	toolRepo := &mockMCPServerToolRepo{listResult: toolList()}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	_, err := cat.ListMCPServerTools(context.Background(), "1", ListMCPServerToolsParams{
		PageSize: 10,
		// OrderBy and SortOrder are zero values (empty strings)
	})
	require.NoError(t, err)

	require.NotNil(t, toolRepo.capturedListOptions)
	pagination := toolRepo.capturedListOptions.Pagination
	orderBy := pagination.GetOrderBy()
	sortOrder := pagination.GetSortOrder()
	assert.Equal(t, "id", orderBy, "expected default orderBy 'id' when param is empty")
	assert.Equal(t, "ASC", sortOrder, "expected default sortOrder 'ASC' when param is empty")
}

func TestGetMCPServer_IncludeToolsUsesAccurateCount(t *testing.T) {
	serverEntity := &models.MCPServerImpl{
		ID:         serverID(1),
		Attributes: &models.MCPServerAttributes{Name: serverName("test-server")},
	}
	repo := &mockMCPServerRepo{getResult: serverEntity}
	toolRepo := &mockMCPServerToolRepo{
		countResult: 10,
		listResult:  toolList(makeTool("tool-1")),
	}
	cat := newTestCatalogWithToolRepo(repo, toolRepo, nil)

	result, err := cat.GetMCPServer(context.Background(), "1", true, 0)
	require.NoError(t, err)
	require.NotNil(t, result)
	// toolCount should reflect total (10), not len(returned tools) (1)
	assert.Equal(t, int32(10), result.ToolCount)
	assert.Len(t, result.Tools, 1)
}
