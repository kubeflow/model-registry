package catalog

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
	mrmodels "github.com/kubeflow/model-registry/internal/db/models"
)

func TestLoadCatalogSources(t *testing.T) {
	type args struct {
		catalogsPath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "test-catalog-sources",
			args:    args{catalogsPath: "testdata/test-catalog-sources.yaml"},
			want:    []string{"catalog1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock services
			services := service.NewServices(
				&MockCatalogModelRepository{},
				&MockCatalogArtifactRepository{},
				&MockCatalogModelArtifactRepository{},
				&MockCatalogMetricsArtifactRepository{},
				&MockPropertyOptionsRepository{},
			)
			loader := NewLoader(services, []string{tt.args.catalogsPath})
			err := loader.Start(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLoader().Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotKeys := make([]string, 0, len(loader.Sources.All()))
			for k := range loader.Sources.All() {
				gotKeys = append(gotKeys, k)
			}
			sort.Strings(gotKeys)
			if !reflect.DeepEqual(gotKeys, tt.want) {
				t.Errorf("NewLoader().Start() got = %v, want %v", gotKeys, tt.want)
			}
		})
	}
}

func TestLoadCatalogSourcesEnabledDisabled(t *testing.T) {
	trueValue := true
	type args struct {
		catalogsPath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]apimodels.CatalogSource
		wantErr bool
	}{
		{
			name: "test-catalog-sources-enabled-and-disabled",
			args: args{catalogsPath: "testdata/test-catalog-sources.yaml"},
			want: map[string]apimodels.CatalogSource{
				"catalog1": {
					Id:      "catalog1",
					Name:    "Catalog 1",
					Enabled: &trueValue,
					Labels:  []string{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock services
			services := service.NewServices(
				&MockCatalogModelRepository{},
				&MockCatalogArtifactRepository{},
				&MockCatalogModelArtifactRepository{},
				&MockCatalogMetricsArtifactRepository{},
				&MockPropertyOptionsRepository{},
			)
			loader := NewLoader(services, []string{tt.args.catalogsPath})
			err := loader.Start(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLoader().Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(loader.Sources.All(), tt.want) {
				t.Errorf("NewLoader().Start() got metadata = %#v, want %#v", loader.Sources.All(), tt.want)
			}
		})
	}
}

func TestLabelsValidation(t *testing.T) {
	// Create mock services
	services := service.NewServices(
		&MockCatalogModelRepository{},
		&MockCatalogArtifactRepository{},
		&MockCatalogModelArtifactRepository{},
		&MockCatalogMetricsArtifactRepository{},
		&MockPropertyOptionsRepository{},
	)

	tests := []struct {
		name    string
		config  *sourceConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid labels with name field",
			config: &sourceConfig{
				Catalogs: []Source{},
				Labels: []map[string]any{
					{"name": "labelNameOne", "displayName": "Label Name One"},
					{"name": "labelNameTwo", "displayName": "Label Name Two"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid label missing name field",
			config: &sourceConfig{
				Catalogs: []Source{},
				Labels: []map[string]any{
					{"name": "labelNameOne", "displayName": "Label Name One"},
					{"displayName": "Label Name Two"}, // Missing "name"
				},
			},
			wantErr: true,
			errMsg:  "invalid label at index 1: missing required 'name' field",
		},
		{
			name: "invalid label with empty name",
			config: &sourceConfig{
				Catalogs: []Source{},
				Labels: []map[string]any{
					{"name": "", "displayName": "Empty Name"},
				},
			},
			wantErr: true,
			errMsg:  "invalid label at index 0: missing required 'name' field",
		},
		{
			name: "duplicate label names within same origin",
			config: &sourceConfig{
				Catalogs: []Source{},
				Labels: []map[string]any{
					{"name": "labelNameOne", "displayName": "Label Name One 1"},
					{"name": "labelNameTwo", "displayName": "Label Name Two"},
					{"name": "labelNameOne", "displayName": "Label Name One 2"},
				},
			},
			wantErr: true,
			errMsg:  "duplicate label name 'labelNameOne' within the same origin",
		},
		{
			name: "nil labels should not error",
			config: &sourceConfig{
				Catalogs: []Source{},
				Labels:   nil,
			},
			wantErr: false,
		},
		{
			name: "empty labels array should not error",
			config: &sourceConfig{
				Catalogs: []Source{},
				Labels:   []map[string]any{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewLoader(services, []string{})
			err := loader.updateLabels("test-path", tt.config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("updateLabels() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("updateLabels() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("updateLabels() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestCatalogSourceLabelsDefaultToEmptySlice(t *testing.T) {
	type args struct {
		catalogsPath string
	}
	tests := []struct {
		name string
		args args
		want func(sources map[string]apimodels.CatalogSource) bool
	}{
		{
			name: "labels-default-to-empty-slice",
			args: args{catalogsPath: "testdata/test-catalog-sources.yaml"},
			want: func(sources map[string]apimodels.CatalogSource) bool {
				// Verify that all loaded catalog sources have labels defaulting to empty slice
				for _, source := range sources {
					if source.Labels == nil {
						return false // Labels should not be nil
					}
					if len(source.Labels) != 0 {
						return false // Labels should be empty slice, not nil and not containing elements
					}
				}
				return len(sources) > 0 // Ensure we actually loaded some sources to test
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock services
			services := service.NewServices(
				&MockCatalogModelRepository{},
				&MockCatalogArtifactRepository{},
				&MockCatalogModelArtifactRepository{},
				&MockCatalogMetricsArtifactRepository{},
				&MockPropertyOptionsRepository{},
			)
			loader := NewLoader(services, []string{tt.args.catalogsPath})
			err := loader.Start(context.Background())
			if err != nil {
				t.Errorf("NewLoader().Start() error = %v", err)
				return
			}

			sources := loader.Sources.All()
			if !tt.want(sources) {
				t.Errorf("Labels validation failed for sources: %#v", sources)
			}

			// Explicitly verify each source has empty labels slice
			for id, source := range sources {
				if source.Labels == nil {
					t.Errorf("Source %s has nil Labels, expected empty slice", id)
				} else if len(source.Labels) != 0 {
					t.Errorf("Source %s has non-empty Labels %v, expected empty slice", id, source.Labels)
				}
			}
		})
	}
}

func TestLoadCatalogSourcesWithMockRepositories(t *testing.T) {
	// Create mock repositories with tracking capabilities
	mockModelRepo := &MockCatalogModelRepository{}
	mockArtifactRepo := &MockCatalogArtifactRepository{}
	mockModelArtifactRepo := &MockCatalogModelArtifactRepository{}
	mockMetricsArtifactRepo := &MockCatalogMetricsArtifactRepository{}

	services := service.NewServices(
		mockModelRepo,
		mockArtifactRepo,
		mockModelArtifactRepo,
		mockMetricsArtifactRepo,
		&MockPropertyOptionsRepository{},
	)

	// Register a test provider that will create some test data
	testProviderName := "test-provider"
	RegisterModelProvider(testProviderName, func(ctx context.Context, source *Source, reldir string) (<-chan ModelProviderRecord, error) {
		ch := make(chan ModelProviderRecord, 1)

		// Create a test model
		modelName := "test-model"
		model := &dbmodels.CatalogModelImpl{
			Attributes: &dbmodels.CatalogModelAttributes{
				Name: &modelName,
			},
		}

		// Create test artifacts
		modelArtifactName := "model-artifact"
		metricsArtifactName := "metrics-artifact"

		artifacts := []dbmodels.CatalogArtifact{
			{
				CatalogModelArtifact: &dbmodels.CatalogModelArtifactImpl{
					Attributes: &dbmodels.CatalogModelArtifactAttributes{
						Name: &modelArtifactName,
					},
				},
			},
			{
				CatalogMetricsArtifact: &dbmodels.CatalogMetricsArtifactImpl{
					Attributes: &dbmodels.CatalogMetricsArtifactAttributes{
						Name: &metricsArtifactName,
					},
				},
			},
		}

		ch <- ModelProviderRecord{
			Model:     model,
			Artifacts: artifacts,
		}
		close(ch)

		return ch, nil
	})

	// Create test config content (use in-memory instead of file)
	testConfig := &sourceConfig{
		Catalogs: []Source{
			{
				CatalogSource: apimodels.CatalogSource{
					Id:      "test-catalog",
					Name:    "Test Catalog",
					Enabled: apiutils.Of(true),
				},
				Type: testProviderName,
				Properties: map[string]any{
					"test": "property",
				},
			},
		},
	}

	// Create a loader and test the database update directly
	l := NewLoader(services, []string{})
	ctx := context.Background()

	err := l.updateDatabase(ctx, "test-path", testConfig)
	if err != nil {
		t.Fatalf("updateDatabase() error = %v", err)
	}

	// Wait a bit for the goroutine to process
	time.Sleep(100 * time.Millisecond)

	// Verify that the model was saved
	if len(mockModelRepo.SavedModels) != 1 {
		t.Errorf("Expected 1 model to be saved, got %d", len(mockModelRepo.SavedModels))
	}

	if len(mockModelRepo.SavedModels) > 0 {
		savedModel := mockModelRepo.SavedModels[0]
		if savedModel.GetAttributes() == nil || savedModel.GetAttributes().Name == nil {
			t.Error("Saved model should have attributes with name")
		} else if *savedModel.GetAttributes().Name != "test-model" {
			t.Errorf("Expected model name 'test-model', got '%s'", *savedModel.GetAttributes().Name)
		}
	}

	// Verify that artifacts were saved
	if len(mockModelArtifactRepo.SavedArtifacts) != 1 {
		t.Errorf("Expected 1 model artifact to be saved, got %d", len(mockModelArtifactRepo.SavedArtifacts))
	}

	if len(mockMetricsArtifactRepo.SavedMetrics) != 1 {
		t.Errorf("Expected 1 metrics artifact to be saved, got %d", len(mockMetricsArtifactRepo.SavedMetrics))
	}
}

func TestLoadCatalogSourcesWithRepositoryErrors(t *testing.T) {
	// Create a mock repository that fails on save
	mockModelRepo := &MockCatalogModelRepositoryWithErrors{
		shouldFailSave: true,
	}
	mockArtifactRepo := &MockCatalogArtifactRepository{}
	mockModelArtifactRepo := &MockCatalogModelArtifactRepository{}
	mockMetricsArtifactRepo := &MockCatalogMetricsArtifactRepository{}

	services := service.NewServices(
		mockModelRepo,
		mockArtifactRepo,
		mockModelArtifactRepo,
		mockMetricsArtifactRepo,
		&MockPropertyOptionsRepository{},
	)

	// Register a test provider
	testProviderName := "test-error-provider"
	RegisterModelProvider(testProviderName, func(ctx context.Context, source *Source, reldir string) (<-chan ModelProviderRecord, error) {
		ch := make(chan ModelProviderRecord, 1)

		modelName := "test-model"
		model := &dbmodels.CatalogModelImpl{
			Attributes: &dbmodels.CatalogModelAttributes{
				Name: &modelName,
			},
		}

		ch <- ModelProviderRecord{
			Model:     model,
			Artifacts: []dbmodels.CatalogArtifact{},
		}
		close(ch)

		return ch, nil
	})

	testConfig := &sourceConfig{
		Catalogs: []Source{
			{
				CatalogSource: apimodels.CatalogSource{
					Id:      "test-catalog",
					Name:    "Test Catalog",
					Enabled: apiutils.Of(true),
				},
				Type: testProviderName,
			},
		},
	}

	l := NewLoader(services, []string{})
	ctx := context.Background()

	// This should not return an error even if repository operations fail
	// (errors are logged but don't stop the loading process)
	err := l.updateDatabase(ctx, "test-path", testConfig)
	if err != nil {
		t.Fatalf("updateDatabase() should not fail even with repository errors, got error = %v", err)
	}

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Verify that no models were saved due to the error
	if len(mockModelRepo.SavedModels) != 0 {
		t.Errorf("Expected 0 models to be saved due to error, got %d", len(mockModelRepo.SavedModels))
	}
}

func TestMockRepositoryBehavior(t *testing.T) {
	mockRepo := &MockCatalogModelRepository{}

	// Test Save operation
	modelName := "test-model"
	model := &dbmodels.CatalogModelImpl{
		Attributes: &dbmodels.CatalogModelAttributes{
			Name: &modelName,
		},
	}

	savedModel, err := mockRepo.Save(model)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if savedModel.GetID() == nil {
		t.Error("Saved model should have an ID")
	}

	if *savedModel.GetID() != 1 {
		t.Errorf("Expected ID 1, got %d", *savedModel.GetID())
	}

	// Test GetByID operation
	retrievedModel, err := mockRepo.GetByID(1)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if retrievedModel.GetAttributes().Name == nil || *retrievedModel.GetAttributes().Name != modelName {
		t.Errorf("Retrieved model name mismatch, expected %s", modelName)
	}

	// Test GetByName operation
	retrievedModel, err = mockRepo.GetByName(modelName)
	if err != nil {
		t.Fatalf("GetByName() error = %v", err)
	}

	if retrievedModel.GetID() == nil || *retrievedModel.GetID() != 1 {
		t.Error("Retrieved model should have ID 1")
	}

	// Test List operation
	listWrapper, err := mockRepo.List(dbmodels.CatalogModelListOptions{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(listWrapper.Items) != 1 {
		t.Errorf("Expected 1 item in list, got %d", len(listWrapper.Items))
	}

	// Test not found scenarios
	_, err = mockRepo.GetByID(999)
	if err == nil {
		t.Error("GetByID() should return error for non-existent ID")
	}

	_, err = mockRepo.GetByName("non-existent")
	if err == nil {
		t.Error("GetByName() should return error for non-existent name")
	}
}

// MockCatalogModelRepositoryWithErrors is a mock that can simulate errors
type MockCatalogModelRepositoryWithErrors struct {
	MockCatalogModelRepository
	shouldFailSave bool
}

func (m *MockCatalogModelRepositoryWithErrors) Save(model dbmodels.CatalogModel) (dbmodels.CatalogModel, error) {
	if m.shouldFailSave {
		return nil, fmt.Errorf("simulated save error")
	}
	return m.MockCatalogModelRepository.Save(model)
}

// MockCatalogModelRepository mocks the CatalogModelRepository interface.
type MockCatalogModelRepository struct {
	SavedModels []dbmodels.CatalogModel
	NextID      int32
}

func (m *MockCatalogModelRepository) GetByID(id int32) (dbmodels.CatalogModel, error) {
	for _, model := range m.SavedModels {
		if model.GetID() != nil && *model.GetID() == id {
			return model, nil
		}
	}
	return nil, &MockNotFoundError{Entity: "CatalogModel", ID: id}
}

func (m *MockCatalogModelRepository) List(listOptions dbmodels.CatalogModelListOptions) (*mrmodels.ListWrapper[dbmodels.CatalogModel], error) {
	return &mrmodels.ListWrapper[dbmodels.CatalogModel]{
		Items:    m.SavedModels,
		PageSize: int32(len(m.SavedModels)),
		Size:     int32(len(m.SavedModels)),
	}, nil
}

func (m *MockCatalogModelRepository) GetByName(name string) (dbmodels.CatalogModel, error) {
	for _, model := range m.SavedModels {
		if model.GetAttributes() != nil && model.GetAttributes().Name != nil && *model.GetAttributes().Name == name {
			return model, nil
		}
	}
	return nil, &MockNotFoundError{Entity: "CatalogModel", ID: 0}
}

func (m *MockCatalogModelRepository) Save(model dbmodels.CatalogModel) (dbmodels.CatalogModel, error) {
	m.NextID++
	id := m.NextID

	// Create a new model with assigned ID
	savedModel := &dbmodels.CatalogModelImpl{
		ID:               &id,
		TypeID:           model.GetTypeID(),
		Attributes:       model.GetAttributes(),
		Properties:       model.GetProperties(),
		CustomProperties: model.GetCustomProperties(),
	}

	m.SavedModels = append(m.SavedModels, savedModel)
	return savedModel, nil
}

// MockCatalogModelArtifactRepository mocks the CatalogModelArtifactRepository interface.
type MockCatalogModelArtifactRepository struct {
	SavedArtifacts []dbmodels.CatalogModelArtifact
	NextID         int32
}

func (m *MockCatalogModelArtifactRepository) GetByID(id int32) (dbmodels.CatalogModelArtifact, error) {
	for _, artifact := range m.SavedArtifacts {
		if artifact.GetID() != nil && *artifact.GetID() == id {
			return artifact, nil
		}
	}
	return nil, &MockNotFoundError{Entity: "CatalogModelArtifact", ID: id}
}

func (m *MockCatalogModelArtifactRepository) List(listOptions dbmodels.CatalogModelArtifactListOptions) (*mrmodels.ListWrapper[dbmodels.CatalogModelArtifact], error) {
	return &mrmodels.ListWrapper[dbmodels.CatalogModelArtifact]{
		Items:    m.SavedArtifacts,
		PageSize: int32(len(m.SavedArtifacts)),
		Size:     int32(len(m.SavedArtifacts)),
	}, nil
}

func (m *MockCatalogModelArtifactRepository) Save(modelArtifact dbmodels.CatalogModelArtifact, parentResourceID *int32) (dbmodels.CatalogModelArtifact, error) {
	m.NextID++
	id := m.NextID

	// Create a new artifact with assigned ID
	savedArtifact := &dbmodels.CatalogModelArtifactImpl{
		ID:               &id,
		TypeID:           modelArtifact.GetTypeID(),
		Attributes:       modelArtifact.GetAttributes(),
		Properties:       modelArtifact.GetProperties(),
		CustomProperties: modelArtifact.GetCustomProperties(),
	}

	m.SavedArtifacts = append(m.SavedArtifacts, savedArtifact)
	return savedArtifact, nil
}

// MockCatalogMetricsArtifactRepository mocks the CatalogMetricsArtifactRepository interface.
type MockCatalogMetricsArtifactRepository struct {
	SavedMetrics []dbmodels.CatalogMetricsArtifact
	NextID       int32
}

func (m *MockCatalogMetricsArtifactRepository) GetByID(id int32) (dbmodels.CatalogMetricsArtifact, error) {
	for _, metrics := range m.SavedMetrics {
		if metrics.GetID() != nil && *metrics.GetID() == id {
			return metrics, nil
		}
	}
	return nil, &MockNotFoundError{Entity: "CatalogMetricsArtifact", ID: id}
}

func (m *MockCatalogMetricsArtifactRepository) List(listOptions dbmodels.CatalogMetricsArtifactListOptions) (*mrmodels.ListWrapper[dbmodels.CatalogMetricsArtifact], error) {
	return &mrmodels.ListWrapper[dbmodels.CatalogMetricsArtifact]{
		Items:    m.SavedMetrics,
		PageSize: int32(len(m.SavedMetrics)),
		Size:     int32(len(m.SavedMetrics)),
	}, nil
}

func (m *MockCatalogMetricsArtifactRepository) Save(metricsArtifact dbmodels.CatalogMetricsArtifact, parentResourceID *int32) (dbmodels.CatalogMetricsArtifact, error) {
	m.NextID++
	id := m.NextID

	// Create a new metrics artifact with assigned ID
	savedMetrics := &dbmodels.CatalogMetricsArtifactImpl{
		ID:               &id,
		TypeID:           metricsArtifact.GetTypeID(),
		Attributes:       metricsArtifact.GetAttributes(),
		Properties:       metricsArtifact.GetProperties(),
		CustomProperties: metricsArtifact.GetCustomProperties(),
	}

	m.SavedMetrics = append(m.SavedMetrics, savedMetrics)
	return savedMetrics, nil
}

func (m *MockCatalogMetricsArtifactRepository) BatchSave(metricsArtifacts []dbmodels.CatalogMetricsArtifact, parentResourceID *int32) ([]dbmodels.CatalogMetricsArtifact, error) {
	savedArtifacts := make([]dbmodels.CatalogMetricsArtifact, len(metricsArtifacts))

	for i, metricsArtifact := range metricsArtifacts {
		m.NextID++
		id := m.NextID

		// Create a new metrics artifact with assigned ID
		savedMetrics := &dbmodels.CatalogMetricsArtifactImpl{
			ID:               &id,
			TypeID:           metricsArtifact.GetTypeID(),
			Attributes:       metricsArtifact.GetAttributes(),
			Properties:       metricsArtifact.GetProperties(),
			CustomProperties: metricsArtifact.GetCustomProperties(),
		}

		m.SavedMetrics = append(m.SavedMetrics, savedMetrics)
		savedArtifacts[i] = savedMetrics
	}

	return savedArtifacts, nil
}

// MockCatalogArtifactRepository mocks the CatalogArtifactRepository interface.
type MockCatalogArtifactRepository struct {
	SavedArtifacts []dbmodels.CatalogArtifact
	NextID         int32
}

func (m *MockCatalogArtifactRepository) GetByID(id int32) (dbmodels.CatalogArtifact, error) {
	for _, artifact := range m.SavedArtifacts {
		// Check both model and metrics artifacts for the ID
		if artifact.CatalogModelArtifact != nil && artifact.CatalogModelArtifact.GetID() != nil && *artifact.CatalogModelArtifact.GetID() == id {
			return artifact, nil
		}
		if artifact.CatalogMetricsArtifact != nil && artifact.CatalogMetricsArtifact.GetID() != nil && *artifact.CatalogMetricsArtifact.GetID() == id {
			return artifact, nil
		}
	}
	return dbmodels.CatalogArtifact{}, &MockNotFoundError{Entity: "CatalogArtifact", ID: id}
}

func (m *MockCatalogArtifactRepository) List(listOptions dbmodels.CatalogArtifactListOptions) (*mrmodels.ListWrapper[dbmodels.CatalogArtifact], error) {
	return &mrmodels.ListWrapper[dbmodels.CatalogArtifact]{
		Items:    m.SavedArtifacts,
		PageSize: int32(len(m.SavedArtifacts)),
		Size:     int32(len(m.SavedArtifacts)),
	}, nil
}

func (m *MockCatalogArtifactRepository) DeleteByParentID(artifactType string, parentResourceID int32) error {
	// Simple mock implementation - could be enhanced to actually filter and delete
	return nil
}

// MockNotFoundError represents an error when an entity is not found.
type MockNotFoundError struct {
	Entity string
	ID     int32
}

func (e *MockNotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Entity, e.ID)
}

// MockPropertyOptionsRepository mocks the PropertyOptionsRepository interface.
type MockPropertyOptionsRepository struct {
	RefreshCalls []dbmodels.PropertyOptionType
	ListCalls    []struct {
		Type   dbmodels.PropertyOptionType
		TypeID int32
	}
	MockOptions map[dbmodels.PropertyOptionType]map[int32][]dbmodels.PropertyOption
}

func NewMockPropertyOptionsRepository() *MockPropertyOptionsRepository {
	return &MockPropertyOptionsRepository{
		RefreshCalls: make([]dbmodels.PropertyOptionType, 0),
		ListCalls: make([]struct {
			Type   dbmodels.PropertyOptionType
			TypeID int32
		}, 0),
		MockOptions: make(map[dbmodels.PropertyOptionType]map[int32][]dbmodels.PropertyOption),
	}
}

func (m *MockPropertyOptionsRepository) Refresh(t dbmodels.PropertyOptionType) error {
	m.RefreshCalls = append(m.RefreshCalls, t)
	return nil
}

func (m *MockPropertyOptionsRepository) List(t dbmodels.PropertyOptionType, typeID int32) ([]dbmodels.PropertyOption, error) {
	m.ListCalls = append(m.ListCalls, struct {
		Type   dbmodels.PropertyOptionType
		TypeID int32
	}{Type: t, TypeID: typeID})

	if typeMap, exists := m.MockOptions[t]; exists {
		if options, exists := typeMap[typeID]; exists {
			return options, nil
		}
	}

	// Return empty slice by default
	return []dbmodels.PropertyOption{}, nil
}

// SetMockOptions allows tests to set up mock data for specific types and typeIDs.
func (m *MockPropertyOptionsRepository) SetMockOptions(t dbmodels.PropertyOptionType, typeID int32, options []dbmodels.PropertyOption) {
	if m.MockOptions[t] == nil {
		m.MockOptions[t] = make(map[int32][]dbmodels.PropertyOption)
	}
	m.MockOptions[t][typeID] = options
}
