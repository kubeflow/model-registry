package catalog

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

// LoadPerformanceMetricsData loads performance metrics data from the specified directory
// into the database using the catalog model and artifact repositories.
func LoadPerformanceMetricsData(path []string, modelRepo dbmodels.CatalogModelRepository, metricsArtifactRepo dbmodels.CatalogMetricsArtifactRepository, typeMap map[string]int64) (map[string]*model.CatalogModel, error) {
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

	loadedModels := make(map[string]*model.CatalogModel)
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
			// if err := processModelDirectory(dirPath, modelRepo, metricsArtifactRepo, modelTypeID, metricsArtifactTypeID, loadedModels); err != nil {
			// 	glog.Errorf("Failed to process model directory %s: %v", dirPath, err)
			// 	// Continue processing other directories
			// 	return nil
			// }

			processedCount++
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("failed to walk performance metrics directory %s: %v", rootPath, err)
		}
	}

	glog.Infof("Successfully processed %d model directories and loaded %d models into database", processedCount, len(loadedModels))
	return loadedModels, nil
}
