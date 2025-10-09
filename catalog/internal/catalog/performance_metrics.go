package catalog

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/models"
)

// metadataJSON represents the structure of metadata.json files
// Core fields are those defined by CatalogModel (BaseModel + BaseResource)
// All other fields are captured dynamically in CustomProperties
type metadataJSON struct {
	// Core fields defined by CatalogModel/BaseResource
	ID          string   `json:"id"`            // Maps to name and externalId
	Description string   `json:"description"`   // From BaseModel
	Readme      string   `json:"readme"`        // From BaseModel
	Maturity    string   `json:"maturity"`      // From BaseModel
	Languages   []string `json:"languages"`     // Maps to "language" from BaseModel
	Tasks       []string `json:"tasks"`         // From BaseModel
	Provider    string   `json:"provider_name"` // Maps to "provider" from BaseModel
	Logo        string   `json:"logo"`          // From BaseModel
	License     string   `json:"license"`       // From BaseModel
	LicenseLink string   `json:"license_link"`  // Maps to "licenseLink" from BaseModel
	LibraryName string   `json:"library_name"`  // Maps to "libraryName" from BaseModel
	CreatedAt   int64    `json:"created_at"`    // Maps to createTimeSinceEpoch from BaseResource
	UpdatedAt   int64    `json:"updated_at"`    // Maps to lastUpdateTimeSinceEpoch from BaseResource

	// CustomProperties captures any additional fields not defined above
	CustomProperties map[string]interface{} `json:"-"`
}

// parseMetadataJSON parses JSON data into metadataJSON struct, capturing unknown fields as custom properties
func parseMetadataJSON(data []byte) (metadataJSON, error) {
	// First, unmarshal into a map to capture all fields
	var allFields map[string]interface{}
	if err := json.Unmarshal(data, &allFields); err != nil {
		return metadataJSON{}, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Define the known fields that should not be included in custom properties
	knownFields := map[string]bool{
		"id":            true,
		"description":   true,
		"readme":        true,
		"maturity":      true,
		"languages":     true,
		"tasks":         true,
		"provider_name": true,
		"logo":          true,
		"license":       true,
		"license_link":  true,
		"library_name":  true,
		"created_at":    true,
		"updated_at":    true,
	}

	// Extract custom properties (fields not in the known list)
	customProperties := make(map[string]interface{})
	for key, value := range allFields {
		if !knownFields[key] {
			customProperties[key] = value
		}
	}

	// Now unmarshal into the structured type for the known fields
	var metadata metadataJSON
	if err := json.Unmarshal(data, &metadata); err != nil {
		return metadataJSON{}, fmt.Errorf("failed to unmarshal into structured type: %v", err)
	}

	// Add the custom properties
	metadata.CustomProperties = customProperties

	return metadata, nil
}

// LoadPerformanceMetricsData loads performance metrics data from the specified directory
// into the database using the catalog model and artifact repositories.
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
func processModelDirectory(dirPath string, modelRepo dbmodels.CatalogModelRepository, metricsArtifactRepo dbmodels.CatalogMetricsArtifactRepository, modelTypeID int32, metricsArtifactTypeID int32) error {
	// Read and parse metadata.json
	metadataPath := filepath.Join(dirPath, "metadata.json")
	metadataData, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata file %s: %v", metadataPath, err)
	}

	// Parse metadata with dynamic field capture
	metadata, err := parseMetadataJSON(metadataData)
	if err != nil {
		return fmt.Errorf("failed to parse metadata file %s: %v", metadataPath, err)
	}

	// Create and save the catalog model
	_, err = createAndSaveModel(metadata, dirPath, modelTypeID, modelRepo)
	if err != nil {
		return fmt.Errorf("failed to create/save model: %v", err)
	}

	return nil
}

// createAndSaveModel creates and saves a catalog model from metadata, skipping if model already exists
func createAndSaveModel(metadata metadataJSON, dirPath string, modelTypeID int32, modelRepo dbmodels.CatalogModelRepository) (dbmodels.CatalogModel, error) {
	// Check if a model with this Name already exists
	existingModels, err := modelRepo.List(dbmodels.CatalogModelListOptions{
		Name: &metadata.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing model: %v", err)
	}

	// If model already exists, skip creating it
	if existingModels.Size > 0 {
		existingModel := existingModels.Items[0]
		glog.V(2).Infof("Model %s already exists with ID %d, skipping", metadata.ID, *existingModel.GetID())
		return existingModel, nil
	}

	glog.V(2).Infof("Creating new model %s", metadata.ID)

	// Create the catalog model for new model (no existing ID)
	dbModel := createDBModelFromMetadata(metadata, dirPath, modelTypeID)

	// Save the model
	savedModel, err := modelRepo.Save(dbModel)
	if err != nil {
		return nil, fmt.Errorf("failed to save new model to database: %v", err)
	}

	glog.V(2).Infof("Saved model %s with ID %d", metadata.ID, *savedModel.GetID())
	return savedModel, nil
}

// createDBModelFromMetadata converts metadata to a database model, mapping properties and custom properties
func createDBModelFromMetadata(metadata metadataJSON, dirPath string, typeID int32) *dbmodels.CatalogModelImpl {
	// Create properties for core CatalogModel fields
	properties := []models.Properties{
		{
			Name:        "source_id",
			StringValue: stringPtr("performance-metrics"),
		},
	}

	// Add core fields if they're not empty
	if metadata.Description != "" {
		properties = append(properties, models.Properties{
			Name:        "description",
			StringValue: &metadata.Description,
		})
	}

	if metadata.Readme != "" {
		properties = append(properties, models.Properties{
			Name:        "readme",
			StringValue: &metadata.Readme,
		})
	}

	if metadata.Maturity != "" {
		properties = append(properties, models.Properties{
			Name:        "maturity",
			StringValue: &metadata.Maturity,
		})
	}

	if metadata.Provider != "" {
		properties = append(properties, models.Properties{
			Name:        "provider",
			StringValue: &metadata.Provider,
		})
	}

	if metadata.Logo != "" {
		properties = append(properties, models.Properties{
			Name:        "logo",
			StringValue: &metadata.Logo,
		})
	}

	if metadata.License != "" {
		properties = append(properties, models.Properties{
			Name:        "license",
			StringValue: &metadata.License,
		})
	}

	if metadata.LicenseLink != "" {
		properties = append(properties, models.Properties{
			Name:        "license_link",
			StringValue: &metadata.LicenseLink,
		})
	}

	if metadata.LibraryName != "" {
		properties = append(properties, models.Properties{
			Name:        "library_name",
			StringValue: &metadata.LibraryName,
		})
	}

	// Add language as JSON array
	if len(metadata.Languages) > 0 {
		languageJSON, _ := json.Marshal(metadata.Languages)
		languageStr := string(languageJSON)
		properties = append(properties, models.Properties{
			Name:        "language",
			StringValue: &languageStr,
		})
	}

	// Add tasks as JSON array
	if len(metadata.Tasks) > 0 {
		tasksJSON, _ := json.Marshal(metadata.Tasks)
		tasksStr := string(tasksJSON)
		properties = append(properties, models.Properties{
			Name:        "tasks",
			StringValue: &tasksStr,
		})
	}

	// Create custom properties from the dynamically captured fields
	customProperties := []models.Properties{}

	// Add all custom properties from the CustomProperties map
	for key, value := range metadata.CustomProperties {
		// Handle different value types
		switch v := value.(type) {
		case string:
			if v != "" {
				customProperties = append(customProperties, models.Properties{
					Name:        key,
					StringValue: &v,
				})
			}
		case []interface{}, map[string]interface{}:
			// For arrays and objects, marshal to JSON string
			jsonBytes, err := json.Marshal(v)
			if err == nil {
				jsonStr := string(jsonBytes)
				customProperties = append(customProperties, models.Properties{
					Name:        key,
					StringValue: &jsonStr,
				})
			}
		case float64:
			// Numbers come as float64 from JSON
			strValue := fmt.Sprintf("%v", v)
			customProperties = append(customProperties, models.Properties{
				Name:        key,
				StringValue: &strValue,
			})
		case int64:
			strValue := fmt.Sprintf("%d", v)
			customProperties = append(customProperties, models.Properties{
				Name:        key,
				StringValue: &strValue,
			})
		case bool:
			strValue := fmt.Sprintf("%t", v)
			customProperties = append(customProperties, models.Properties{
				Name:        key,
				StringValue: &strValue,
			})
		default:
			// For any other type, try to marshal to JSON
			jsonBytes, err := json.Marshal(v)
			if err == nil {
				jsonStr := string(jsonBytes)
				customProperties = append(customProperties, models.Properties{
					Name:        key,
					StringValue: &jsonStr,
				})
			}
		}
	}

	// Add metadata file path for reference
	metadataFilePath := filepath.Join(dirPath, "metadata.json")
	customProperties = append(customProperties, models.Properties{
		Name:        "metadata_file_path",
		StringValue: &metadataFilePath,
	})

	// Create the model with the provided TypeID from the type map
	model := &dbmodels.CatalogModelImpl{
		TypeID: &typeID,
		Attributes: &dbmodels.CatalogModelAttributes{
			Name:                     &metadata.ID,
			ExternalID:               &metadata.ID,
			CreateTimeSinceEpoch:     &metadata.CreatedAt,
			LastUpdateTimeSinceEpoch: &metadata.UpdatedAt,
		},
		Properties:       &properties,
		CustomProperties: &customProperties,
	}

	return model
}

// stringPtr is a helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}
