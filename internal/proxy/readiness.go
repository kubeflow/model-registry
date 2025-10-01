package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kubeflow/model-registry/internal/db"
	"github.com/kubeflow/model-registry/pkg/api"
)

const (
	// Health check names
	HealthCheckDatabase      = "database"
	HealthCheckModelRegistry = "model-registry"
	HealthCheckMeta          = "meta"

	// Health check statuses
	StatusPass = "pass"
	StatusFail = "fail"

	// HTTP response messages
	responseOK   = "OK"
	responseFail = "FAIL"

	// Database schema query
	schemaMigrationsQuery = "SELECT version, dirty FROM schema_migrations ORDER BY version DESC LIMIT 1"

	// Detail keys
	detailDatastoreType                 = "datastore_type"
	detailSchemaVersion                 = "schema_version"
	detailSchemaDirty                   = "schema_dirty"
	detailRegisteredModelsAccessible    = "registered_models_accessible"
	detailRegisteredModelsCount         = "registered_models_count"
	detailArtifactsAccessible           = "artifacts_accessible"
	detailArtifactsCount                = "artifacts_count"
	detailModelVersionsAccessible       = "model_versions_accessible"
	detailModelVersionsCount            = "model_versions_count"
	detailServingEnvironmentsAccessible = "serving_environments_accessible"
	detailServingEnvironmentsCount      = "serving_environments_count"
	detailInferenceServicesAccessible   = "inference_services_accessible"
	detailInferenceServicesCount        = "inference_services_count"
	detailExperimentsAccessible         = "experiments_accessible"
	detailExperimentsCount              = "experiments_count"
	detailExperimentRunsAccessible      = "experiment_runs_accessible"
	detailExperimentRunsCount           = "experiment_runs_count"
	detailTotalResourcesChecked         = "total_resources_checked"
	detailCheckDurationMs               = "check_duration_ms"
	detailTimestamp                     = "timestamp"
)

// HealthCheck represents a single health check
type HealthCheck struct {
	Name    string
	Status  string
	Message string
	Details map[string]interface{}
}

// HealthStatus represents the overall health status
type HealthStatus struct {
	Status string                 `json:"status"`
	Checks map[string]HealthCheck `json:"checks"`
}

// HealthChecker defines the interface for health checks
type HealthChecker interface {
	Check() HealthCheck
}

// DatabaseHealthChecker checks database connectivity and schema state
type DatabaseHealthChecker struct {
}

func NewDatabaseHealthChecker() *DatabaseHealthChecker {
	return &DatabaseHealthChecker{}
}

func (d *DatabaseHealthChecker) Check() HealthCheck {
	check := HealthCheck{
		Name:    HealthCheckDatabase,
		Details: make(map[string]interface{}),
	}

	// Check database connector
	dbConnector, ok := db.GetConnector()
	if !ok {
		check.Status = StatusFail
		check.Message = "database connector not initialized"
		return check
	}

	// Test database connection
	database, err := dbConnector.Connect()
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("database connection error: %v", err)
		return check
	}

	// Check schema migration state
	var result struct {
		Version int64
		Dirty   bool
	}

	query := schemaMigrationsQuery
	if err := database.Raw(query).Scan(&result).Error; err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("schema_migrations query error: %v", err)
		return check
	}

	check.Details[detailSchemaVersion] = result.Version
	check.Details[detailSchemaDirty] = result.Dirty

	if result.Dirty {
		check.Status = StatusFail
		check.Message = "database schema is in dirty state"
		return check
	}

	check.Status = StatusPass
	check.Message = "database is healthy"
	return check
}

// ModelRegistryHealthChecker checks model registry service health
type ModelRegistryHealthChecker struct {
	service api.ModelRegistryApi
}

func NewModelRegistryHealthChecker(service api.ModelRegistryApi) *ModelRegistryHealthChecker {
	return &ModelRegistryHealthChecker{
		service: service,
	}
}

func (m *ModelRegistryHealthChecker) Check() HealthCheck {
	check := HealthCheck{
		Name:    HealthCheckModelRegistry,
		Details: make(map[string]interface{}),
	}

	if m.service == nil {
		check.Status = StatusFail
		check.Message = "model registry service not available"
		return check
	}

	// Test basic listing operation with minimal page size
	listOptions := api.ListOptions{
		PageSize: func() *int32 { i := int32(1); return &i }(),
	}

	// Test registered models listing
	models, err := m.service.GetRegisteredModels(listOptions)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("failed to list registered models: %v", err)
		return check
	}

	check.Details[detailRegisteredModelsAccessible] = true
	check.Details[detailRegisteredModelsCount] = models.Size

	// Test artifacts listing (all artifacts, not tied to specific model version)
	artifacts, err := m.service.GetArtifacts("", listOptions, nil)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("failed to list artifacts: %v", err)
		return check
	}

	check.Details[detailArtifactsAccessible] = true
	check.Details[detailArtifactsCount] = artifacts.Size

	// Test model versions listing
	versions, err := m.service.GetModelVersions(listOptions, nil)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("failed to list model versions: %v", err)
		return check
	}

	check.Details[detailModelVersionsAccessible] = true
	check.Details[detailModelVersionsCount] = versions.Size

	// Test serving environments listing
	servingEnvs, err := m.service.GetServingEnvironments(listOptions)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("failed to list serving environments: %v", err)
		return check
	}

	check.Details[detailServingEnvironmentsAccessible] = true
	check.Details[detailServingEnvironmentsCount] = servingEnvs.Size

	// Test inference services listing (all services, not tied to specific serving environment or runtime)
	inferenceServices, err := m.service.GetInferenceServices(listOptions, nil, nil)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("failed to list inference services: %v", err)
		return check
	}

	check.Details[detailInferenceServicesAccessible] = true
	check.Details[detailInferenceServicesCount] = inferenceServices.Size

	// Test experiments listing
	experiments, err := m.service.GetExperiments(listOptions)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("failed to list experiments: %v", err)
		return check
	}

	check.Details[detailExperimentsAccessible] = true
	check.Details[detailExperimentsCount] = experiments.Size

	// Test experiment runs listing
	experimentRuns, err := m.service.GetExperimentRuns(listOptions, nil)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("failed to list experiment runs: %v", err)
		return check
	}

	check.Details[detailExperimentRunsAccessible] = true
	check.Details[detailExperimentRunsCount] = experimentRuns.Size

	check.Status = StatusPass
	check.Message = "model registry service is healthy"
	check.Details[detailTotalResourcesChecked] = 5

	return check
}

// GeneralReadinessHandler creates a general readiness handler with configurable health checks
func GeneralReadinessHandler(additionalCheckers ...HealthChecker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		start := time.Now()

		// Initialize health status
		health := HealthStatus{
			Status: StatusPass,
			Checks: make(map[string]HealthCheck),
		}

		// Run health checks
		for _, checker := range additionalCheckers {
			check := checker.Check()
			health.Checks[check.Name] = check

			if check.Status == StatusFail {
				health.Status = StatusFail
			}
		}

		// Add timing information
		duration := time.Since(start)
		if _, exists := health.Checks[HealthCheckMeta]; !exists {
			health.Checks[HealthCheckMeta] = HealthCheck{
				Name:   HealthCheckMeta,
				Status: StatusPass,
				Details: map[string]interface{}{
					detailCheckDurationMs: duration.Milliseconds(),
					detailTimestamp:       time.Now().UTC().Format(time.RFC3339),
				},
			}
		}

		// Set response status
		statusCode := http.StatusOK
		if health.Status == StatusFail {
			statusCode = http.StatusServiceUnavailable
		}

		// Return JSON response for detailed health info, or simple OK for basic checks
		if r.URL.Query().Get("format") == "json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(health)
		} else {
			w.WriteHeader(statusCode)
			if health.Status == StatusPass {
				_, _ = w.Write([]byte(responseOK))
			} else {
				// Return the first failed check's error message
				for _, check := range health.Checks {
					if check.Status == StatusFail {
						_, _ = w.Write([]byte(check.Message))
						return
					}
				}
				_, _ = w.Write([]byte(responseFail))
			}
		}
	})
}
