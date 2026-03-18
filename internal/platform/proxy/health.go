package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	platformdb "github.com/kubeflow/model-registry/internal/platform/db"
)

const (
	HealthCheckDatabase = "database"
	HealthCheckMeta     = "meta"

	StatusPass = "pass"
	StatusFail = "fail"

	responseOK   = "OK"
	responseFail = "FAIL"

	schemaMigrationsQuery = "SELECT version, dirty FROM schema_migrations ORDER BY version DESC LIMIT 1"

	detailSchemaVersion   = "schema_version"
	detailSchemaDirty     = "schema_dirty"
	detailCheckDurationMs = "check_duration_ms"
	detailTimestamp        = "timestamp"
)

// HealthCheck represents a single health check
type HealthCheck struct {
	Name    string
	Status  string
	Message string
	Details map[string]any
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
		Details: make(map[string]any),
	}

	dbConnector := platformdb.GetConnector()
	if dbConnector == nil {
		check.Status = StatusFail
		check.Message = "database connector not initialized"
		return check
	}

	database, err := dbConnector.Connect()
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("database connection error: %v", err)
		return check
	}

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

// GeneralReadinessHandler creates a general readiness handler with configurable health checks
func GeneralReadinessHandler(additionalCheckers ...HealthChecker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		start := time.Now()

		health := HealthStatus{
			Status: StatusPass,
			Checks: make(map[string]HealthCheck),
		}

		for _, checker := range additionalCheckers {
			check := checker.Check()
			health.Checks[check.Name] = check

			if check.Status == StatusFail {
				health.Status = StatusFail
			}
		}

		duration := time.Since(start)
		if _, exists := health.Checks[HealthCheckMeta]; !exists {
			health.Checks[HealthCheckMeta] = HealthCheck{
				Name:   HealthCheckMeta,
				Status: StatusPass,
				Details: map[string]any{
					detailCheckDurationMs: duration.Milliseconds(),
					detailTimestamp:       time.Now().UTC().Format(time.RFC3339),
				},
			}
		}

		statusCode := http.StatusOK
		if health.Status == StatusFail {
			statusCode = http.StatusServiceUnavailable
		}

		if r.URL.Query().Get("format") == "json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(health)
		} else {
			w.WriteHeader(statusCode)
			if health.Status == StatusPass {
				_, _ = w.Write([]byte(responseOK))
			} else {
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
