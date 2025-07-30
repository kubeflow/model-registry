package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/db"
	"github.com/kubeflow/model-registry/pkg/api"
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
	datastore datastore.Datastore
}

func (d *DatabaseHealthChecker) Check() HealthCheck {
	check := HealthCheck{
		Name:    "database",
		Details: make(map[string]interface{}),
	}

	// Skip embedmd check for mlmd datastore
	if d.datastore.Type != "embedmd" {
		check.Status = "pass"
		check.Message = "MLMD datastore - skipping database check"
		check.Details["datastore_type"] = d.datastore.Type
		return check
	}

	// Check DSN configuration
	dsn := d.datastore.EmbedMD.DatabaseDSN
	if dsn == "" {
		check.Status = "fail"
		check.Message = "database DSN not configured"
		return check
	}

	// Check database connector
	dbConnector, ok := db.GetConnector()
	if !ok {
		check.Status = "fail"
		check.Message = "database connector not initialized"
		return check
	}

	// Test database connection
	database, err := dbConnector.Connect()
	if err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("database connection error: %v", err)
		return check
	}

	// Check schema migration state
	var result struct {
		Version int64
		Dirty   bool
	}

	query := "SELECT version, dirty FROM schema_migrations ORDER BY version DESC LIMIT 1"
	if err := database.Raw(query).Scan(&result).Error; err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("schema_migrations query error: %v", err)
		return check
	}

	check.Details["schema_version"] = result.Version
	check.Details["schema_dirty"] = result.Dirty

	if result.Dirty {
		check.Status = "fail"
		check.Message = "database schema is in dirty state"
		return check
	}

	check.Status = "pass"
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
		Name:    "model-registry",
		Details: make(map[string]interface{}),
	}

	if m.service == nil {
		check.Status = "fail"
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
		check.Status = "fail"
		check.Message = fmt.Sprintf("failed to list registered models: %v", err)
		return check
	}

	check.Details["registered_models_accessible"] = true
	check.Details["registered_models_count"] = models.Size

	// Test artifacts listing (all artifacts, not tied to specific model version)
	artifacts, err := m.service.GetArtifacts(listOptions, nil)
	if err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("failed to list artifacts: %v", err)
		return check
	}

	check.Details["artifacts_accessible"] = true
	check.Details["artifacts_count"] = artifacts.Size

	// Test model versions listing
	versions, err := m.service.GetModelVersions(listOptions, nil)
	if err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("failed to list model versions: %v", err)
		return check
	}

	check.Details["model_versions_accessible"] = true
	check.Details["model_versions_count"] = versions.Size

	// Test serving environments listing
	servingEnvs, err := m.service.GetServingEnvironments(listOptions)
	if err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("failed to list serving environments: %v", err)
		return check
	}

	check.Details["serving_environments_accessible"] = true
	check.Details["serving_environments_count"] = servingEnvs.Size

	// Test inference services listing (all services, not tied to specific serving environment or runtime)
	inferenceServices, err := m.service.GetInferenceServices(listOptions, nil, nil)
	if err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("failed to list inference services: %v", err)
		return check
	}

	check.Details["inference_services_accessible"] = true
	check.Details["inference_services_count"] = inferenceServices.Size

	check.Status = "pass"
	check.Message = "model registry service is healthy"
	check.Details["total_resources_checked"] = 5

	return check
}

// GeneralReadinessHandler creates a general readiness handler with configurable health checks
func GeneralReadinessHandler(datastore datastore.Datastore, additionalCheckers ...HealthChecker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		start := time.Now()

		// Initialize health status
		health := HealthStatus{
			Status: "pass",
			Checks: make(map[string]HealthCheck),
		}

		// Always include database health check
		dbChecker := &DatabaseHealthChecker{datastore: datastore}
		dbCheck := dbChecker.Check()
		health.Checks["database"] = dbCheck

		if dbCheck.Status == "fail" {
			health.Status = "fail"
		}

		// Run additional health checks
		for _, checker := range additionalCheckers {
			check := checker.Check()
			health.Checks[check.Name] = check

			if check.Status == "fail" {
				health.Status = "fail"
			}
		}

		// Add timing information
		duration := time.Since(start)
		if _, exists := health.Checks["meta"]; !exists {
			health.Checks["meta"] = HealthCheck{
				Name:   "meta",
				Status: "pass",
				Details: map[string]interface{}{
					"check_duration_ms": duration.Milliseconds(),
					"timestamp":         time.Now().UTC().Format(time.RFC3339),
				},
			}
		}

		// Set response status
		statusCode := http.StatusOK
		if health.Status == "fail" {
			statusCode = http.StatusServiceUnavailable
		}

		// Return JSON response for detailed health info, or simple OK for basic checks
		if r.URL.Query().Get("format") == "json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(health)
		} else {
			w.WriteHeader(statusCode)
			if health.Status == "pass" {
				_, _ = w.Write([]byte("OK"))
			} else {
				_, _ = w.Write([]byte("FAIL"))
			}
		}
	})
}

// ReadinessHandler is a readiness probe that requires schema_migrations.dirty to be false before allowing traffic.
func ReadinessHandler(datastore datastore.Datastore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// skip embedmd check for mlmd datastore
		if datastore.Type != "embedmd" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		dsn := datastore.EmbedMD.DatabaseDSN
		if dsn == "" {
			http.Error(w, "database DSN not configured", http.StatusServiceUnavailable)
			return
		}

		dbConnector, ok := db.GetConnector()
		if !ok {
			http.Error(w, "database connector not initialized", http.StatusServiceUnavailable)
			return
		}

		database, err := dbConnector.Connect()
		if err != nil {
			http.Error(w, fmt.Sprintf("database connection error: %v", err), http.StatusServiceUnavailable)
			return
		}

		var result struct {
			Version int64
			Dirty   bool
		}

		query := "SELECT version, dirty FROM schema_migrations ORDER BY version DESC LIMIT 1"
		if err := database.Raw(query).Scan(&result).Error; err != nil {
			http.Error(w, fmt.Sprintf("schema_migrations query error: %v", err), http.StatusServiceUnavailable)
			return
		}

		if result.Dirty {
			http.Error(w, "database schema is in dirty state", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}

func ReadyzHandler(datastore datastore.Datastore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}
