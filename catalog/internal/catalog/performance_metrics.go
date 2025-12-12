package catalog

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/glog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
)

// metadataJSON represents the minimal structure needed from metadata.json files
// Only the ID field is needed to look up existing models
type metadataJSON struct {
	ID              string   `json:"id"`               // Maps to model name for lookup
	OverallAccuracy *float64 `json:"overall_accuracy"` // Overall accuracy score for the model
	Size            *string  `json:"size"`             // Model parameter count (e.g., "8B params")
	TensorType      *string  `json:"tensor_type"`      // Data precision (e.g., "FP16", "INT4")
	VariantGroupID  *string  `json:"variant_group_id"` // UUID linking model variants together
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
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&raw); err != nil {
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
	// Cache of model ID -> directory path mapping to avoid repeated directory scans
	modelDirCache map[string]string
}

func NewPerformanceMetricsLoader(path []string, modelRepo dbmodels.CatalogModelRepository, metricsArtifactRepo dbmodels.CatalogMetricsArtifactRepository, typeMap map[string]int32) (*PerformanceMetricsLoader, error) {
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
	modelTypeID, exists := typeMap[service.CatalogModelTypeName]
	if !exists {
		return nil, fmt.Errorf("CatalogModel type not found in type map")
	}
	glog.V(2).Infof("Using catalog model type ID: %d", modelTypeID)

	// Get the TypeID for CatalogMetricsArtifact from the type map
	metricsArtifactTypeID, exists := typeMap[service.CatalogMetricsArtifactTypeName]
	if !exists {
		return nil, fmt.Errorf("CatalogMetricsArtifact type not found in type map")
	}
	glog.V(2).Infof("Using metrics artifact type ID: %d", metricsArtifactTypeID)

	loader := &PerformanceMetricsLoader{
		path:                  path,
		modelRepo:             modelRepo,
		metricsArtifactRepo:   metricsArtifactRepo,
		modelTypeID:           modelTypeID,
		metricsArtifactTypeID: metricsArtifactTypeID,
		modelDirCache:         make(map[string]string),
	}

	// Build the model directory cache once during initialization
	if err := loader.buildModelDirCache(); err != nil {
		return nil, fmt.Errorf("failed to build model directory cache: %v", err)
	}

	return loader, nil
}

// buildModelDirCache scans directories once and builds a cache of model ID -> directory path
func (pml *PerformanceMetricsLoader) buildModelDirCache() error {
	modelCount := 0
	for _, rootPath := range pml.path {
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

			// Read and parse metadata.json to extract the model ID
			metadataData, err := os.ReadFile(metadataPath)
			if err != nil {
				glog.Warningf("Failed to read metadata file %s: %v", metadataPath, err)
				return nil // Continue with other directories
			}

			// Parse metadata to extract the model ID for lookup
			metadata, err := parseMetadataJSON(metadataData)
			if err != nil {
				glog.Warningf("Failed to parse metadata file %s: %v", metadataPath, err)
				return nil // Continue with other directories
			}

			// Add to cache
			pml.modelDirCache[metadata.ID] = dirPath
			modelCount++
			glog.V(3).Infof("Cached model directory: %s -> %s", metadata.ID, dirPath)

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to walk directory %s: %v", rootPath, err)
		}
	}

	glog.Infof("Built model directory cache (%d models indexed)", modelCount)

	return nil
}

func (pml *PerformanceMetricsLoader) Load(ctx context.Context, record ModelProviderRecord) error {
	if pml == nil {
		return nil
	}

	attrs := record.Model.GetAttributes()
	if attrs == nil || attrs.Name == nil {
		return nil
	}

	modelName := *attrs.Name
	glog.Infof("Loading performance metrics for %s", modelName)

	// Look up the model directory in the cache
	dirPath, found := pml.modelDirCache[modelName]
	if !found {
		glog.V(2).Infof("No performance metrics directory found for model %s", modelName)
		return nil
	}

	glog.V(2).Infof("Found cached directory for model %s: %s", modelName, dirPath)

	// Process this specific model directory using the cached path
	artifactsCreated, err := processModelDirectory(dirPath, pml.modelRepo, pml.metricsArtifactRepo, pml.modelTypeID, pml.metricsArtifactTypeID)
	if err != nil {
		return fmt.Errorf("failed to process metrics for model %s: %v", modelName, err)
	}

	if artifactsCreated > 0 {
		glog.Infof("Loaded %d performance metrics artifacts for model %s", artifactsCreated, modelName)
	}

	return nil
}

// processModelDirectory processes a single model directory containing metadata.json and metric files
// Only processes metrics for models that already exist in the database
// Returns the number of artifacts created and any error encountered
func processModelDirectory(dirPath string, modelRepo dbmodels.CatalogModelRepository, metricsArtifactRepo dbmodels.CatalogMetricsArtifactRepository, modelTypeID int32, metricsArtifactTypeID int32) (int, error) {
	// Read and parse metadata.json to extract the model ID
	metadataPath := filepath.Join(dirPath, "metadata.json")
	metadataData, err := os.ReadFile(metadataPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read metadata file %s: %v", metadataPath, err)
	}

	// Parse metadata to extract the model ID for lookup
	metadata, err := parseMetadataJSON(metadataData)
	if err != nil {
		return 0, fmt.Errorf("failed to parse metadata file %s: %v", metadataPath, err)
	}

	// Check if the model already exists - only process metrics for existing models
	existingModel, err := modelRepo.GetByName(metadata.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to check for existing model: %v", err)
	}

	// Skip processing if model doesn't exist
	if existingModel == nil {
		glog.V(2).Infof("Model %s does not exist in database, skipping metrics processing", metadata.ID)
		return 0, nil
	}

	// Enrich the model with metadata before processing metrics artifacts
	if err := enrichCatalogModelFromMetadata(existingModel, metadata, modelRepo); err != nil {
		glog.Warningf("Failed to enrich model %s with metadata: %v", metadata.ID, err)
		// Continue processing - don't fail the whole operation
	}

	modelID := *existingModel.GetID()
	glog.V(2).Infof("Found existing model %s with ID %d, processing metrics", metadata.ID, modelID)

	// Use batch processing for all artifacts
	return processModelArtifactsBatch(dirPath, modelID, metadata.ID, metadata.OverallAccuracy, metricsArtifactRepo, metricsArtifactTypeID)
}

// processModelArtifactsBatch processes all metric artifacts for a model in batch
// This reduces DB overhead by parsing, checking, and inserting in optimized phases
func processModelArtifactsBatch(dirPath string, modelID int32, modelName string, overallAccuracy *float64, metricsArtifactRepo dbmodels.CatalogMetricsArtifactRepository, metricsArtifactTypeID int32) (int, error) {
	// Parse all metrics files
	var evaluationRecords []evaluationRecord
	var performanceRecords []performanceRecord

	// Parse evaluation metrics if file exists
	evaluationsPath := filepath.Join(dirPath, "evaluations.ndjson")
	if _, err := os.Stat(evaluationsPath); err == nil {
		records, err := parseEvaluationFile(evaluationsPath)
		if err != nil {
			glog.Errorf("Failed to parse evaluations file for %s: %v", modelName, err)
		} else {
			evaluationRecords = records
		}
	}

	// Parse performance metrics if file exists
	performancePath := filepath.Join(dirPath, "performance.ndjson")
	if _, err := os.Stat(performancePath); err == nil {
		records, err := parsePerformanceFile(performancePath)
		if err != nil {
			glog.Errorf("Failed to parse performance file for %s: %v", modelName, err)
		} else {
			performanceRecords = records
		}
	}

	totalRecords := len(evaluationRecords) + len(performanceRecords)
	if totalRecords == 0 {
		return 0, nil
	}

	// Bulk load all existing artifacts for this model and check in-memory
	// Single DB query to get ALL existing artifacts for this model
	existingArtifactsList, err := metricsArtifactRepo.List(dbmodels.CatalogMetricsArtifactListOptions{
		ParentResourceID: &modelID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to load existing artifacts for model: %v", err)
	}

	// Build in-memory map for O(1) lookups: external_id -> artifact
	existingArtifactsMap := make(map[string]bool, existingArtifactsList.Size)
	for _, artifact := range existingArtifactsList.Items {
		if artifact.GetAttributes() != nil && artifact.GetAttributes().ExternalID != nil {
			existingArtifactsMap[*artifact.GetAttributes().ExternalID] = true
		}
	}

	// Check which artifacts need to be created using the in-memory map
	artifactsToInsert := make([]*dbmodels.CatalogMetricsArtifactImpl, 0, totalRecords)

	// Check evaluation artifacts
	if len(evaluationRecords) > 0 {
		externalID := fmt.Sprintf("accuracy-metrics-model-%d", modelID)
		if !existingArtifactsMap[externalID] {
			artifact := createAccuracyMetricsArtifact(evaluationRecords, modelID, metricsArtifactTypeID, overallAccuracy, nil, nil)
			artifactsToInsert = append(artifactsToInsert, artifact)
		} else {
			glog.V(2).Infof("Accuracy metrics artifact already exists, skipping")
		}
	}

	// Check performance artifacts
	for _, perfRecord := range performanceRecords {
		if !existingArtifactsMap[perfRecord.ID] {
			artifact := createPerformanceArtifact(perfRecord, modelID, metricsArtifactTypeID, nil, nil)
			artifactsToInsert = append(artifactsToInsert, artifact)
		} else {
			glog.V(2).Infof("Performance artifact %s already exists, skipping", perfRecord.ID)
		}
	}

	if len(artifactsToInsert) == 0 {
		glog.V(2).Infof("All artifacts already exist for model %s, nothing to insert", modelName)
		return 0, nil
	}

	// Batch insert all new artifacts using BatchSave
	// Convert to slice of interface type for BatchSave
	artifactsToSave := make([]dbmodels.CatalogMetricsArtifact, len(artifactsToInsert))
	for i, artifact := range artifactsToInsert {
		artifactsToSave[i] = artifact
	}

	savedArtifacts, err := metricsArtifactRepo.BatchSave(artifactsToSave, &modelID)
	if err != nil {
		return 0, fmt.Errorf("failed to batch save artifacts: %v", err)
	}

	return len(savedArtifacts), nil
}

// parseEvaluationFile reads and parses an evaluations.ndjson file
func parseEvaluationFile(filePath string) ([]evaluationRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open evaluation file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	evaluationRecords := []evaluationRecord{}

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
		return nil, fmt.Errorf("error reading evaluation file: %v", err)
	}

	return evaluationRecords, nil
}

// parsePerformanceFile reads and parses a performance.ndjson file
func parsePerformanceFile(filePath string) ([]performanceRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open performance file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	performanceRecords := []performanceRecord{}

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

		performanceRecords = append(performanceRecords, perfRecord)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading performance file: %v", err)
	}

	return performanceRecords, nil
}

// createAccuracyMetricsArtifact creates a single metrics artifact from all evaluation records
func createAccuracyMetricsArtifact(evalRecords []evaluationRecord, modelID int32, typeID int32, overallAccuracy *float64, existingID *int32, existingCreateTime *int64) *dbmodels.CatalogMetricsArtifactImpl {
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
		delete(evalRecord.CustomProperties, "updated_at")
		delete(evalRecord.CustomProperties, "created_at")
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

	// Add overall_average custom property from metadata.json overall_accuracy field
	if overallAccuracy != nil {
		customProperties = append(customProperties, models.Properties{
			Name:        "overall_average",
			DoubleValue: overallAccuracy,
		})
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
		if createdAtNum, ok := perfRecord.CustomProperties["created_at"].(json.Number); ok {
			createdAt, err := createdAtNum.Int64()
			if err == nil {
				createTime = &createdAt
			} else {
				glog.Warningf("%s: invalid created_at value: %v", artifactName, err)
			}
		}
	}
	if createTime == nil {
		createTime = apiutils.Of(time.Now().UnixMilli())
	}

	if updatedAtNum, ok := perfRecord.CustomProperties["updated_at"].(json.Number); ok {
		updatedAt, err := updatedAtNum.Int64()
		if err == nil {
			updateTime = &updatedAt
		} else {
			glog.Warningf("%s: invalid updated_at value: %v", artifactName, err)
		}
	}
	if updateTime == nil {
		updateTime = apiutils.Of(time.Now().UnixMilli())
	}
	delete(perfRecord.CustomProperties, "updated_at")
	delete(perfRecord.CustomProperties, "created_at")

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
			prop.SetInt64Value(v)
		case int:
			intVal := int32(v)
			prop.IntValue = &intVal
		case bool:
			prop.BoolValue = &v
		case json.Number:
			if n, err := v.Int64(); err == nil {
				prop.SetInt64Value(n)
			} else if f, err := v.Float64(); err == nil {
				prop.DoubleValue = &f
			} else {
				// This shouldn't happen, but convert it to a string if it does.
				strVal := v.String()
				prop.StringValue = &strVal
			}
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

// enrichCatalogModelFromMetadata updates CatalogModel with additional fields from metadata.json
func enrichCatalogModelFromMetadata(existingModel dbmodels.CatalogModel, metadata metadataJSON, modelRepo dbmodels.CatalogModelRepository) error {
	// Build custom properties to add/update
	var customProperties []models.Properties

	if metadata.Size != nil && *metadata.Size != "" {
		customProperties = append(customProperties, models.Properties{
			Name:             "size",
			StringValue:      metadata.Size,
			IsCustomProperty: true,
		})
	}

	if metadata.TensorType != nil && *metadata.TensorType != "" {
		customProperties = append(customProperties, models.Properties{
			Name:             "tensor_type",
			StringValue:      metadata.TensorType,
			IsCustomProperty: true,
		})
	}

	if metadata.VariantGroupID != nil && *metadata.VariantGroupID != "" {
		customProperties = append(customProperties, models.Properties{
			Name:             "variant_group_id",
			StringValue:      metadata.VariantGroupID,
			IsCustomProperty: true,
		})
	}

	if len(customProperties) == 0 {
		return nil // Nothing to update
	}

	// Add the new custom properties to existing model
	existingCustomProperties := existingModel.GetCustomProperties()
	if existingCustomProperties == nil {
		existingCustomProperties = &customProperties
	} else {
		*existingCustomProperties = append(*existingCustomProperties, customProperties...)
	}

	// Save the updated model
	_, err := modelRepo.Save(existingModel)
	if err != nil {
		return fmt.Errorf("failed to save enriched model: %v", err)
	}

	glog.V(2).Infof("Enriched model %s with %d custom properties", *existingModel.GetAttributes().Name, len(customProperties))
	return nil
}
