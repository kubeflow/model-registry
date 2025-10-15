package catalog

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/models"
)

// metadataJSON represents the minimal structure needed from metadata.json files
// Only the ID field is needed to look up existing models
type metadataJSON struct {
	ID string `json:"id"` // Maps to model name for lookup
}

// parseMetadataJSON parses JSON data into metadataJSON struct, extracting only the ID field
func parseMetadataJSON(data []byte) (metadataJSON, error) {
	var metadata metadataJSON
	if err := json.Unmarshal(data, &metadata); err != nil {
		return metadataJSON{}, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if metadata.ID == "" {
		return metadataJSON{}, fmt.Errorf("missing required 'id' field in metadata")
	}

	return metadata, nil
}

// evaluationRecord represents a single evaluation result from evaluations.ndjson
// Only minimal fields needed for association are explicitly defined
// evaluationRecords will be merged into a single accuracy-metrics artifact
type evaluationRecord struct {
	// Core fields needed to associate evaluation with model
	ModelID   string `json:"model_id"`
	Benchmark string `json:"benchmark"`

	// CustomProperties captures all other fields dynamically
	CustomProperties map[string]interface{} `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling to capture all undefined fields as CustomProperties
func (er *evaluationRecord) UnmarshalJSON(data []byte) error {
	// First unmarshal into a generic map to get all fields
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract the core fields
	if modelID, ok := raw["model_id"].(string); ok {
		er.ModelID = modelID
	}
	if benchmark, ok := raw["benchmark"].(string); ok {
		er.Benchmark = benchmark
	}

	// Initialize CustomProperties if nil
	if er.CustomProperties == nil {
		er.CustomProperties = make(map[string]interface{})
	}

	// Copy all fields to CustomProperties, including the core ones
	for key, value := range raw {
		er.CustomProperties[key] = value
	}

	return nil
}

// performanceRecord represents a single performance result from performance.ndjson
// Only minimal fields needed for association are explicitly defined
type performanceRecord struct {
	// Core fields needed to associate performance data with model
	ID      string `json:"id"`
	ModelID string `json:"model_id"`

	// CustomProperties captures remaining fields dynamically
	CustomProperties map[string]interface{} `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling to capture all undefined fields as CustomProperties
func (pr *performanceRecord) UnmarshalJSON(data []byte) error {
	// First unmarshal into a generic map to get all fields
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract the core fields
	if id, ok := raw["id"].(string); ok {
		pr.ID = id
	}
	if modelID, ok := raw["model_id"].(string); ok {
		pr.ModelID = modelID
	}

	// Initialize CustomProperties if nil
	if pr.CustomProperties == nil {
		pr.CustomProperties = make(map[string]interface{})
	}

	// Copy all fields to CustomProperties, including the core ones
	for key, value := range raw {
		pr.CustomProperties[key] = value
	}

	return nil
}

type PerformanceMetricsLoader struct {
	path                  []string
	modelRepo             dbmodels.CatalogModelRepository
	metricsArtifactRepo   dbmodels.CatalogMetricsArtifactRepository
	modelTypeID           int32
	metricsArtifactTypeID int32
}

func NewPerformanceMetricsLoader(path []string, modelRepo dbmodels.CatalogModelRepository, metricsArtifactRepo dbmodels.CatalogMetricsArtifactRepository, typeMap map[string]int64) (*PerformanceMetricsLoader, error) {
	if len(path) == 0 {
		glog.Info("No performance metrics path provided, skipping performance metrics loading")
		return nil, nil
	}

	// Check if path exists
	for _, p := range path {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			glog.Warningf("Performance metrics path %s does not exist, skipping performance metrics loading", p)
			return nil, nil
		}
	}

	glog.Infof("Loading performance metrics data from %s", path)

	// Get the TypeID for CatalogModel from the type map
	modelTypeIDInt64, exists := typeMap[service.CatalogModelTypeName]
	if !exists {
		return nil, fmt.Errorf("CatalogModel type not found in type map")
	}
	// Bounds check for int64 to int32 conversion
	if modelTypeIDInt64 > math.MaxInt32 || modelTypeIDInt64 < math.MinInt32 {
		return nil, fmt.Errorf("CatalogModel type ID %d is out of int32 range", modelTypeIDInt64)
	}
	modelTypeID := int32(modelTypeIDInt64)
	glog.V(2).Infof("Using catalog model type ID: %d", modelTypeID)

	// Get the TypeID for CatalogMetricsArtifact from the type map
	metricsArtifactTypeIDInt64, exists := typeMap[service.CatalogMetricsArtifactTypeName]
	if !exists {
		return nil, fmt.Errorf("CatalogMetricsArtifact type not found in type map")
	}
	// Bounds check for int64 to int32 conversion
	if metricsArtifactTypeIDInt64 > math.MaxInt32 || metricsArtifactTypeIDInt64 < math.MinInt32 {
		return nil, fmt.Errorf("CatalogMetricsArtifact type ID %d is out of int32 range", metricsArtifactTypeIDInt64)
	}
	metricsArtifactTypeID := int32(metricsArtifactTypeIDInt64)
	glog.V(2).Infof("Using metrics artifact type ID: %d", metricsArtifactTypeID)

	return &PerformanceMetricsLoader{
		path:                  path,
		modelRepo:             modelRepo,
		metricsArtifactRepo:   metricsArtifactRepo,
		modelTypeID:           modelTypeID,
		metricsArtifactTypeID: metricsArtifactTypeID,
	}, nil
}

func (pml *PerformanceMetricsLoader) Load(ctx context.Context, record ModelProviderRecord) error {
	if pml == nil {
		return nil
	}

	attrs := record.Model.GetAttributes()
	if attrs == nil || attrs.Name == nil {
		return nil
	}

	glog.Infof("Loading performance metrics for %s", *attrs.Name)

	// Create a type map from the stored type IDs
	typeMap := map[string]int64{
		service.CatalogModelTypeName:           int64(pml.modelTypeID),
		service.CatalogMetricsArtifactTypeName: int64(pml.metricsArtifactTypeID),
	}

	// Call the existing LoadPerformanceMetricsData function
	return LoadPerformanceMetricsData(pml.path, pml.modelRepo, pml.metricsArtifactRepo, typeMap)
}

// LoadPerformanceMetricsData loads performance metrics data from the specified directory
// into the database using the catalog model and artifact repositories.
// Only loads metrics for models that already exist in the database and where
// the metrics artifacts don't already exist.
func LoadPerformanceMetricsData(path []string, modelRepo dbmodels.CatalogModelRepository, metricsArtifactRepo dbmodels.CatalogMetricsArtifactRepository, typeMap map[string]int64) error {
	if len(path) == 0 {
		glog.Info("No performance metrics path provided, skipping performance metrics loading")
		return nil
	}

	// Check if path exists
	for _, p := range path {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			glog.Warningf("Performance metrics path %s does not exist, skipping performance metrics loading", p)
			return nil
		}
	}

	glog.Infof("Loading performance metrics data from %s", path)

	// Get the TypeID for CatalogModel from the type map
	modelTypeIDInt64, exists := typeMap[service.CatalogModelTypeName]
	if !exists {
		return fmt.Errorf("CatalogModel type not found in type map")
	}
	// Bounds check for int64 to int32 conversion
	if modelTypeIDInt64 > math.MaxInt32 || modelTypeIDInt64 < math.MinInt32 {
		return fmt.Errorf("CatalogModel type ID %d is out of int32 range", modelTypeIDInt64)
	}
	modelTypeID := int32(modelTypeIDInt64)
	glog.V(2).Infof("Using catalog model type ID: %d", modelTypeID)

	// Get the TypeID for CatalogMetricsArtifact from the type map
	metricsArtifactTypeIDInt64, exists := typeMap[service.CatalogMetricsArtifactTypeName]
	if !exists {
		return fmt.Errorf("CatalogMetricsArtifact type not found in type map")
	}
	// Bounds check for int64 to int32 conversion
	if metricsArtifactTypeIDInt64 > math.MaxInt32 || metricsArtifactTypeIDInt64 < math.MinInt32 {
		return fmt.Errorf("CatalogMetricsArtifact type ID %d is out of int32 range", metricsArtifactTypeIDInt64)
	}
	metricsArtifactTypeID := int32(metricsArtifactTypeIDInt64)
	glog.V(2).Infof("Using metrics artifact type ID: %d", metricsArtifactTypeID)

	processedCount := 0

	// Walk through the directory structure to find model directories
	for _, rootPath := range path {
		err := filepath.Walk(rootPath, func(dirPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip if not a directory
			if !info.IsDir() {
				return nil
			}

			// Check if this directory contains metadata.json
			metadataPath := filepath.Join(dirPath, "metadata.json")
			if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
				return nil // Skip directories without metadata.json
			}

			glog.Infof("Processing model directory: %s", dirPath)

			// Process the model directory
			if err := processModelDirectory(dirPath, modelRepo, metricsArtifactRepo, modelTypeID, metricsArtifactTypeID); err != nil {
				glog.Errorf("Failed to process model directory %s: %v", dirPath, err)
				// Continue processing other directories
				return nil
			}

			processedCount++
			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to walk performance metrics directory %s: %v", rootPath, err)
		}
	}

	glog.Infof("Successfully processed %d model directories", processedCount)
	return nil
}

// processModelDirectory processes a single model directory containing metadata.json and metric files
// Only processes metrics for models that already exist in the database
func processModelDirectory(dirPath string, modelRepo dbmodels.CatalogModelRepository, metricsArtifactRepo dbmodels.CatalogMetricsArtifactRepository, modelTypeID int32, metricsArtifactTypeID int32) error {
	// Read and parse metadata.json to extract the model ID
	metadataPath := filepath.Join(dirPath, "metadata.json")
	metadataData, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata file %s: %v", metadataPath, err)
	}

	// Parse metadata to extract the model ID for lookup
	metadata, err := parseMetadataJSON(metadataData)
	if err != nil {
		return fmt.Errorf("failed to parse metadata file %s: %v", metadataPath, err)
	}

	// Check if the model already exists - only process metrics for existing models
	existingModel, err := findExistingModel(metadata, modelRepo)
	if err != nil {
		return fmt.Errorf("failed to check for existing model: %v", err)
	}

	// Skip processing if model doesn't exist
	if existingModel == nil {
		glog.V(2).Infof("Model %s does not exist in database, skipping metrics processing", metadata.ID)
		return nil
	}

	glog.V(2).Infof("Found existing model %s with ID %d, processing metrics", metadata.ID, *existingModel.GetID())

	// Process evaluation metrics if evaluations.ndjson exists and create separate metric artifacts
	evaluationsPath := filepath.Join(dirPath, "evaluations.ndjson")
	if _, err := os.Stat(evaluationsPath); err == nil {
		if err := processEvaluationArtifacts(evaluationsPath, *existingModel.GetID(), metricsArtifactRepo, metricsArtifactTypeID); err != nil {
			glog.Errorf("Failed to process evaluation artifacts for %s: %v", metadata.ID, err)
		}
	}

	// Process performance metrics if performance.ndjson exists
	performancePath := filepath.Join(dirPath, "performance.ndjson")
	if _, err := os.Stat(performancePath); err == nil {
		if err := processPerformanceArtifacts(performancePath, *existingModel.GetID(), metricsArtifactRepo, metricsArtifactTypeID); err != nil {
			glog.Errorf("Failed to process performance artifacts for %s: %v", metadata.ID, err)
		}
	}

	return nil
}

// findExistingModel checks if a model with the given metadata already exists in the database
func findExistingModel(metadata metadataJSON, modelRepo dbmodels.CatalogModelRepository) (dbmodels.CatalogModel, error) {
	// Check if a model with this Name already exists
	existingModels, err := modelRepo.List(dbmodels.CatalogModelListOptions{
		Name: &metadata.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing model: %v", err)
	}

	// Return the existing model if found, nil otherwise
	if existingModels.Size > 0 {
		existingModel := existingModels.Items[0]
		return existingModel, nil
	}

	return nil, nil
}

// processEvaluationArtifacts processes evaluation metrics from evaluations.ndjson and creates a single accuracy-metrics artifact
// Only creates the artifact if it doesn't already exist in the database
func processEvaluationArtifacts(filePath string, modelID int32, metricsArtifactRepo dbmodels.CatalogMetricsArtifactRepository, metricsArtifactTypeID int32) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open evaluation file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	evaluationRecords := []evaluationRecord{}

	// Read all evaluation records
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var evalRecord evaluationRecord
		if err := json.Unmarshal([]byte(line), &evalRecord); err != nil {
			glog.Errorf("Failed to parse evaluation record: %v", err)
			continue
		}

		evaluationRecords = append(evaluationRecords, evalRecord)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading evaluation file: %v", err)
	}

	if len(evaluationRecords) == 0 {
		glog.V(2).Infof("No evaluations found in %s", filePath)
		return nil
	}

	// Use a consistent external_id for the accuracy metrics artifact
	externalID := fmt.Sprintf("accuracy-metrics-model-%d", modelID)

	// Check if artifact with this external_id already exists
	existingArtifacts, err := metricsArtifactRepo.List(dbmodels.CatalogMetricsArtifactListOptions{
		ExternalID: &externalID,
	})
	if err != nil {
		return fmt.Errorf("failed to check for existing accuracy metrics artifact: %v", err)
	}

	// Skip creating artifact if it already exists
	if existingArtifacts.Size > 0 {
		existingArtifact := existingArtifacts.Items[0]
		glog.V(2).Infof("Accuracy metrics artifact already exists with ID %d, skipping", *existingArtifact.GetID())
		return nil
	}

	glog.V(2).Infof("Creating new accuracy metrics artifact")

	// Create artifact (no existing ID since we're only creating new ones)
	artifact := createAccuracyMetricsArtifact(evaluationRecords, modelID, metricsArtifactTypeID, nil, nil)

	// Save artifact to database
	if _, err := metricsArtifactRepo.Save(artifact, &modelID); err != nil {
		return fmt.Errorf("failed to save accuracy metrics artifact: %v", err)
	}

	glog.V(2).Infof("Processed accuracy metrics artifact with %d evaluations from %s", len(evaluationRecords), filePath)
	return nil
}

// processPerformanceArtifacts processes performance metrics from performance.ndjson
// Only creates artifacts that don't already exist in the database
func processPerformanceArtifacts(filePath string, modelID int32, metricsArtifactRepo dbmodels.CatalogMetricsArtifactRepository, metricsArtifactTypeID int32) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open performance file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	artifactCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var perfRecord performanceRecord
		if err := json.Unmarshal([]byte(line), &perfRecord); err != nil {
			glog.Errorf("Failed to parse performance record: %v", err)
			continue
		}

		// Check if artifact with this external_id already exists
		listOptions := dbmodels.CatalogMetricsArtifactListOptions{
			ExternalID: &perfRecord.ID,
		}
		existingArtifacts, err := metricsArtifactRepo.List(listOptions)
		if err != nil {
			glog.Errorf("Failed to check for existing performance artifact: %v", err)
			continue
		}

		// Skip creating artifact if it already exists
		if existingArtifacts.Size > 0 {
			existingArtifact := existingArtifacts.Items[0]
			glog.V(2).Infof("Performance artifact %s already exists with ID %d, skipping", perfRecord.ID, *existingArtifact.GetID())
			continue
		}

		glog.V(2).Infof("Creating new performance artifact %s", perfRecord.ID)

		// Create artifact (no existing ID since we're only creating new ones)
		artifact := createPerformanceArtifact(perfRecord, modelID, metricsArtifactTypeID, nil, nil)

		// Save artifact to database
		if _, err := metricsArtifactRepo.Save(artifact, &modelID); err != nil {
			glog.Errorf("Failed to save performance artifact: %v", err)
			continue
		}

		artifactCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading performance file: %v", err)
	}

	glog.V(2).Infof("Processed %d performance artifacts from %s", artifactCount, filePath)
	return nil
}

// createAccuracyMetricsArtifact creates a single metrics artifact from all evaluation records
func createAccuracyMetricsArtifact(evalRecords []evaluationRecord, modelID int32, typeID int32, existingID *int32, existingCreateTime *int64) *dbmodels.CatalogMetricsArtifactImpl {
	artifactName := fmt.Sprintf("accuracy-metrics-model-%d", modelID)
	externalID := fmt.Sprintf("accuracy-metrics-model-%d", modelID)

	// Use existing create time if provided, otherwise find from evaluation records
	createTime := existingCreateTime
	var updateTime *int64

	for _, evalRecord := range evalRecords {
		if existingCreateTime == nil {
			if createdAtFloat, ok := evalRecord.CustomProperties["created_at"].(float64); ok {
				createdAt := int64(createdAtFloat)
				if createTime == nil || createdAt < *createTime {
					createTime = &createdAt
				}
			}
		}
		if updatedAtFloat, ok := evalRecord.CustomProperties["updated_at"].(float64); ok {
			updatedAt := int64(updatedAtFloat)
			if updateTime == nil || updatedAt > *updateTime {
				updateTime = &updatedAt
			}
		}
	}

	// Properties can be empty or contain general metadata
	properties := []models.Properties{}

	// Create custom properties - simple mapping of benchmark_name to score_value
	customProperties := []models.Properties{}

	for _, evalRecord := range evalRecords {
		// Add the benchmark score as a named property (e.g., "aime24": 63.3333)
		if score, ok := evalRecord.CustomProperties["score"].(float64); ok {
			customProperties = append(customProperties, models.Properties{
				Name:        evalRecord.Benchmark,
				DoubleValue: &score,
			})
		}
	}

	// Create the metrics artifact with metricsType set to accuracy-metrics
	metricsArtifact := &dbmodels.CatalogMetricsArtifactImpl{
		ID:     existingID, // Use existing ID if updating
		TypeID: &typeID,
		Attributes: &dbmodels.CatalogMetricsArtifactAttributes{
			Name:                     &artifactName,
			ExternalID:               &externalID,
			CreateTimeSinceEpoch:     createTime,
			LastUpdateTimeSinceEpoch: updateTime,
			MetricsType:              dbmodels.MetricsTypeAccuracy,
		},
		Properties:       &properties,
		CustomProperties: &customProperties,
	}

	return metricsArtifact
}

// createPerformanceArtifact creates a metrics artifact from performance record
func createPerformanceArtifact(perfRecord performanceRecord, modelID int32, typeID int32, existingID *int32, existingCreateTime *int64) *dbmodels.CatalogMetricsArtifactImpl {
	// Create artifact name (must be unique per artifact)
	artifactName := fmt.Sprintf("performance-%s", perfRecord.ID)

	// Use existing create time if provided, otherwise extract from custom properties
	createTime := existingCreateTime
	var updateTime *int64

	if existingCreateTime == nil {
		if createdAtFloat, ok := perfRecord.CustomProperties["created_at"].(float64); ok {
			createdAt := int64(createdAtFloat)
			createTime = &createdAt
		}
	}
	if updatedAtFloat, ok := perfRecord.CustomProperties["updated_at"].(float64); ok {
		updatedAt := int64(updatedAtFloat)
		updateTime = &updatedAt
	}

	// Properties can be empty - all data goes in custom properties
	properties := []models.Properties{}

	// Create custom properties - simple mapping of all performance data
	customProperties := []models.Properties{}

	// Add all fields from the performance record as custom properties
	for key, value := range perfRecord.CustomProperties {
		prop := models.Properties{Name: key}

		// Handle different value types
		switch v := value.(type) {
		case string:
			prop.StringValue = &v
		case float64:
			prop.DoubleValue = &v
		case int64:
			intVal := int32(v)
			prop.IntValue = &intVal
		case int:
			intVal := int32(v)
			prop.IntValue = &intVal
		case bool:
			prop.BoolValue = &v
		default:
			// Convert other types to string representation
			strVal := fmt.Sprintf("%v", v)
			prop.StringValue = &strVal
		}

		customProperties = append(customProperties, prop)
	}

	// Create the metrics artifact with metricsType set to performance-metrics
	metricsArtifact := &dbmodels.CatalogMetricsArtifactImpl{
		ID:     existingID, // Use existing ID if updating
		TypeID: &typeID,
		Attributes: &dbmodels.CatalogMetricsArtifactAttributes{
			Name:                     &artifactName,
			ExternalID:               &perfRecord.ID,
			CreateTimeSinceEpoch:     createTime,
			LastUpdateTimeSinceEpoch: updateTime,
			MetricsType:              dbmodels.MetricsTypePerformance,
		},
		Properties:       &properties,
		CustomProperties: &customProperties,
	}

	return metricsArtifact
}
