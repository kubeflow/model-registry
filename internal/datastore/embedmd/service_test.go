package embedmd

import (
	"testing"

	"github.com/kubeflow/model-registry/internal/db/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbedMDConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *EmbedMDConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "valid mysql basic dsn",
			cfg: &EmbedMDConfig{
				DatabaseType: types.DatabaseTypeMySQL,
				DatabaseDSN:  "user:pass@tcp(localhost:3306)/dbname",
			},
			wantErr: false,
		},
		{
			name: "mysql with single query params allowed",
			cfg: &EmbedMDConfig{
				DatabaseType: types.DatabaseTypeMySQL,
				DatabaseDSN:  "user:pass@tcp(localhost:3306)/dbname?charset=utf8mb4",
			},
			wantErr: false,
		},
		{
			name: "mysql invalid dsn parse",
			cfg: &EmbedMDConfig{
				DatabaseType: types.DatabaseTypeMySQL,
				DatabaseDSN:  "://not-a-valid-dsn",
			},
			wantErr:     true,
			errContains: "invalid MySQL DSN",
		},
		{
			name: "mysql with two question marks should fail",
			cfg: &EmbedMDConfig{
				DatabaseType: types.DatabaseTypeMySQL,
				DatabaseDSN:  "user:pass@tcp(localhost:3306)/dbname?charset=utf8mb4?parseTime=true",
			},
			wantErr:     true,
			errContains: "invalid MySQL DSN",
		},
		{
			name: "postgres url ok",
			cfg: &EmbedMDConfig{
				DatabaseType: types.DatabaseTypePostgres,
				DatabaseDSN:  "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
			},
			wantErr: false,
		},

		{
			name: "postgres key=value ok",
			cfg: &EmbedMDConfig{
				DatabaseType: types.DatabaseTypePostgres,
				DatabaseDSN:  "host=localhost port=5432 user=user password=pass dbname=mydb sslmode=disable",
			},
			wantErr: false,
		},
		{
			name: "postgres key=value with '?' in dbname should fail",
			cfg: &EmbedMDConfig{
				DatabaseType: types.DatabaseTypePostgres,
				DatabaseDSN:  "host=localhost port=5432 user=user dbname=my?db sslmode=disable",
			},
			wantErr:     true,
			errContains: "invalid PostgreSQL DSN",
		},
		{
			name: "unsupported db type",
			cfg: &EmbedMDConfig{
				DatabaseType: "sqlite",
				DatabaseDSN:  "file:test.db",
			},
			wantErr:     true,
			errContains: "unsupported database type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
