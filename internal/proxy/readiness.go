package proxy

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/db"
)

// ReadinessHandler is a readiness probe that requires schema_migrations.dirty to be false before allowing traffic.
func ReadinessHandler(datastore datastore.Datastore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// skip embedmd check for mlmd datastore
		if datastore.Type != "embedmd" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
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

		dbConnector, err := db.NewConnector(datastore.EmbedMD.DatabaseType, dsn, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("database new connector error: %s", err), http.StatusServiceUnavailable)
			return
		}

		db, err := dbConnector.Connect()
		if err != nil {
			http.Error(w, fmt.Sprintf("database connection error: %v", err), http.StatusServiceUnavailable)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			http.Error(w, fmt.Sprintf("database sqlDB connection error: %v", err), http.StatusServiceUnavailable)
			return
		}
		defer func() {
			if err := sqlDB.Close(); err != nil {
				glog.Errorf("error closing database: %v", err)
			}
		}()

		var result struct {
			Version int64
			Dirty   int
		}
		query := "SELECT version, dirty FROM schema_migrations ORDER BY version DESC LIMIT 1"
		if err := db.Raw(query).Scan(&result).Error; err != nil {
			http.Error(w, fmt.Sprintf("schema_migrations query error: %v", err), http.StatusServiceUnavailable)
			return
		}

		if result.Dirty != 0 {
			http.Error(w, "database schema is in dirty state", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}
