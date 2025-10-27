package embedmd

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/tls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cont_mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

// Mock repository types for testing
type mockArtifactRepo struct {
	typeID int32
}

type mockContextRepo struct {
	typeID int32
}

type mockExecutionRepo struct {
	typeID int32
}

type mockOtherRepo struct {
	db *gorm.DB
}

// Mock initializer functions
func newMockArtifactRepo(db *gorm.DB, typeID int32) *mockArtifactRepo {
	return &mockArtifactRepo{typeID: typeID}
}

func newMockContextRepo(db *gorm.DB, typeID int32) *mockContextRepo {
	return &mockContextRepo{typeID: typeID}
}

func newMockExecutionRepo(db *gorm.DB, typeID int32) *mockExecutionRepo {
	return &mockExecutionRepo{typeID: typeID}
}

func newMockOtherRepo(db *gorm.DB) *mockOtherRepo {
	return &mockOtherRepo{db: db}
}

func newMockOtherRepoWithError(db *gorm.DB) (*mockOtherRepo, error) {
	return nil, errors.New("mock initialization error")
}

// Test helper to create a test database with types using a simplified MySQL setup
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	// Create a simple MySQL container without config file
	ctx := context.Background()

	mysqlContainer, err := cont_mysql.Run(ctx, "mysql:8.3",
		cont_mysql.WithDatabase("test"),
		cont_mysql.WithUsername("root"),
		cont_mysql.WithPassword("root"),
	)
	require.NoError(t, err)

	// Get connection string
	dsn, err := mysqlContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Connect to database
	dbConnector := mysql.NewMySQLDBConnector(dsn, &tls.TLSConfig{})
	db, err := dbConnector.Connect()
	require.NoError(t, err)

	// Create the required tables
	err = db.AutoMigrate(&schema.Type{})
	require.NoError(t, err)

	// Insert test types
	testTypes := []schema.Type{
		{ID: 1, Name: "TestArtifact", TypeKind: 1},
		{ID: 2, Name: "TestDoc", TypeKind: 1},
		{ID: 3, Name: "TestContext", TypeKind: 2},
		{ID: 4, Name: "TestModel", TypeKind: 2},
		{ID: 5, Name: "TestExecution", TypeKind: 3},
	}

	for _, typ := range testTypes {
		require.NoError(t, db.Create(&typ).Error)
	}

	cleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			//nolint:errcheck
			sqlDB.Close()
		}
		//nolint:errcheck
		mysqlContainer.Terminate(ctx)
	}

	return db, cleanup
}

func TestNewRepoSet_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	spec := datastore.NewSpec().
		AddArtifact("TestArtifact", datastore.NewSpecType(newMockArtifactRepo)).
		AddArtifact("TestDoc", datastore.NewSpecType(newMockArtifactRepo)).
		AddContext("TestContext", datastore.NewSpecType(newMockContextRepo)).
		AddContext("TestModel", datastore.NewSpecType(newMockContextRepo)).
		AddExecution("TestExecution", datastore.NewSpecType(newMockExecutionRepo)).
		AddOther(newMockOtherRepo)

	repoSet, err := newRepoSet(db, spec)
	require.NoError(t, err)
	assert.NotNil(t, repoSet)

	// Verify we can get repositories by type
	mockArtifact, err := repoSet.Repository(reflect.TypeOf(&mockArtifactRepo{}))
	require.NoError(t, err)
	assert.NotNil(t, mockArtifact)

	mockContext, err := repoSet.Repository(reflect.TypeOf(&mockContextRepo{}))
	require.NoError(t, err)
	assert.NotNil(t, mockContext)

	mockExecution, err := repoSet.Repository(reflect.TypeOf(&mockExecutionRepo{}))
	require.NoError(t, err)
	assert.NotNil(t, mockExecution)

	mockOther, err := repoSet.Repository(reflect.TypeOf(&mockOtherRepo{}))
	require.NoError(t, err)
	assert.NotNil(t, mockOther)

	// Verify TypeMap returns correct mappings
	typeMap := repoSet.TypeMap()
	assert.Equal(t, int32(1), typeMap["TestArtifact"])
	assert.Equal(t, int32(2), typeMap["TestDoc"])
	assert.Equal(t, int32(3), typeMap["TestContext"])
	assert.Equal(t, int32(4), typeMap["TestModel"])
	assert.Equal(t, int32(5), typeMap["TestExecution"])
}

func TestNewRepoSet_MissingType(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	spec := datastore.NewSpec().
		AddArtifact("NonExistentType", datastore.NewSpecType(newMockArtifactRepo))

	repoSet, err := newRepoSet(db, spec)
	assert.Error(t, err)
	assert.Nil(t, repoSet)
	assert.Contains(t, err.Error(), "required type 'NonExistentType' not found in database")
}

func TestNewRepoSet_InitializerError(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	spec := datastore.NewSpec().AddOther(newMockOtherRepoWithError)

	repoSet, err := newRepoSet(db, spec)
	assert.Error(t, err)
	assert.Nil(t, repoSet)
	assert.Contains(t, err.Error(), "mock initialization error")
}

// Define a simple interface and implementation for interface testing
type TestInterface interface {
	TestMethod() string
}

type testImpl struct{}

func (ti *testImpl) TestMethod() string {
	return "test"
}

func TestRepoSetImpl_Repository_InterfaceMatch(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create initializer that returns the implementation
	newTestImpl := func(db *gorm.DB) *testImpl {
		return &testImpl{}
	}

	spec := datastore.NewSpec().AddOther(newTestImpl)

	repoSet, err := newRepoSet(db, spec)
	require.NoError(t, err)

	// Should be able to get repository by interface type
	repo, err := repoSet.Repository(reflect.TypeOf((*TestInterface)(nil)).Elem())
	require.NoError(t, err)
	assert.NotNil(t, repo)

	// Verify it's the correct implementation
	impl, ok := repo.(*testImpl)
	assert.True(t, ok)
	assert.Equal(t, "test", impl.TestMethod())
}

func TestRepoSetImpl_Repository_UnknownType(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	spec := datastore.NewSpec().
		AddArtifact("TestArtifact", datastore.NewSpecType(newMockArtifactRepo))

	repoSet, err := newRepoSet(db, spec)
	require.NoError(t, err)

	// Try to get a repository type that doesn't exist
	type unknownType struct{}
	repo, err := repoSet.Repository(reflect.TypeOf(&unknownType{}))
	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "unknown repository type")
}

func TestRepoSetImpl_Call_InvalidInitializer(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	rs := &repoSetImpl{
		db: db,
	}

	args := map[reflect.Type]any{
		reflect.TypeOf(db): db,
	}

	// Test with non-function
	_, err := rs.call("not a function", args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "initializer is not a function")

	// Test with function that has no return values
	noReturnFunc := func() {}
	_, err = rs.call(noReturnFunc, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "initializer has no return value")

	// Test with function that has too many return values
	tooManyReturnsFunc := func() (int, string, error) {
		return 0, "", nil
	}
	_, err = rs.call(tooManyReturnsFunc, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "more than 2 return values")

	// Test with missing argument type
	needsIntFunc := func(i int) string {
		return "test"
	}
	_, err = rs.call(needsIntFunc, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no initializer argument for type")
}

func TestRepoSetImpl_Call_ValidInitializers(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	rs := &repoSetImpl{
		db: db,
	}

	args := map[reflect.Type]any{
		reflect.TypeOf(db):       db,
		reflect.TypeOf(int32(0)): int32(42),
	}

	// Test function with one return value
	oneReturnFunc := func(db *gorm.DB) *mockOtherRepo {
		return &mockOtherRepo{db: db}
	}
	result, err := rs.call(oneReturnFunc, args)
	require.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo, ok := result.(*mockOtherRepo)
	assert.True(t, ok)
	assert.Equal(t, db, mockRepo.db)

	// Test function with two return values (success case)
	twoReturnSuccessFunc := func(db *gorm.DB, id int32) (*mockArtifactRepo, error) {
		return &mockArtifactRepo{typeID: id}, nil
	}
	result, err = rs.call(twoReturnSuccessFunc, args)
	require.NoError(t, err)
	assert.NotNil(t, result)
	mockArtifact, ok := result.(*mockArtifactRepo)
	assert.True(t, ok)
	assert.Equal(t, int32(42), mockArtifact.typeID)

	// Test function with two return values (error case)
	twoReturnErrorFunc := func(db *gorm.DB) (*mockOtherRepo, error) {
		return nil, errors.New("initialization failed")
	}
	result, err = rs.call(twoReturnErrorFunc, args)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "initialization failed")
}

func TestMakeTypeMap(t *testing.T) {
	specMap := map[string]*datastore.SpecType{
		"type1": datastore.NewSpecType("func1"),
		"type2": datastore.NewSpecType("func2"),
		"type3": datastore.NewSpecType("func3"),
	}

	nameIDMap := map[string]int32{
		"type1": 10,
		"type2": 20,
		"type3": 30,
	}

	// Test ArtifactTypeMap
	artifactMap := makeTypeMap[datastore.ArtifactTypeMap](specMap, nameIDMap)
	assert.Equal(t, datastore.ArtifactTypeMap{"type1": 10, "type2": 20, "type3": 30}, artifactMap)

	// Test ContextTypeMap
	contextMap := makeTypeMap[datastore.ContextTypeMap](specMap, nameIDMap)
	assert.Equal(t, datastore.ContextTypeMap{"type1": 10, "type2": 20, "type3": 30}, contextMap)

	// Test ExecutionTypeMap
	executionMap := makeTypeMap[datastore.ExecutionTypeMap](specMap, nameIDMap)
	assert.Equal(t, datastore.ExecutionTypeMap{"type1": 10, "type2": 20, "type3": 30}, executionMap)
}

func TestRepoSetImpl_TypeMapCloning(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	spec := datastore.NewSpec().
		AddArtifact("TestArtifact", datastore.NewSpecType(newMockArtifactRepo))

	repoSet, err := newRepoSet(db, spec)
	require.NoError(t, err)

	rs := repoSet.(*repoSetImpl)

	// Get the type map
	typeMap1 := rs.TypeMap()
	typeMap2 := rs.TypeMap()

	// Verify they have the same values
	assert.Equal(t, typeMap1, typeMap2)

	// Verify they are different objects (cloned)
	typeMap1["TestModification"] = 999
	assert.NotEqual(t, typeMap1, typeMap2)
	assert.NotContains(t, typeMap2, "TestModification")
}

// Integration test with real service repositories
func TestRepoSetImpl_WithRealRepositories(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create artifact types that match what the real repositories expect
	realArtifactTypes := []schema.Type{
		{ID: 10, Name: "model-artifact", TypeKind: 1},
		{ID: 11, Name: "doc-artifact", TypeKind: 1},
	}

	for _, at := range realArtifactTypes {
		require.NoError(t, db.Create(&at).Error)
	}

	spec := datastore.NewSpec().
		AddArtifact("model-artifact", datastore.NewSpecType(service.NewModelArtifactRepository)).
		AddArtifact("doc-artifact", datastore.NewSpecType(service.NewDocArtifactRepository)).
		AddOther(service.NewArtifactRepository)

	repoSet, err := newRepoSet(db, spec)
	require.NoError(t, err)
	assert.NotNil(t, repoSet)

	// Verify we can get the real repositories
	modelRepo, err := repoSet.Repository(reflect.TypeOf((*models.ModelArtifactRepository)(nil)).Elem())
	require.NoError(t, err)
	assert.NotNil(t, modelRepo)

	docRepo, err := repoSet.Repository(reflect.TypeOf((*models.DocArtifactRepository)(nil)).Elem())
	require.NoError(t, err)
	assert.NotNil(t, docRepo)

	artifactRepo, err := repoSet.Repository(reflect.TypeOf((*models.ArtifactRepository)(nil)).Elem())
	require.NoError(t, err)
	assert.NotNil(t, artifactRepo)
}
