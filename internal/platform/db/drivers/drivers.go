// Package drivers registers all built-in database backends (MySQL, PostgreSQL)
// with the platform/db connector and migrator factories.
//
// Import this package for its side effects to ensure all supported database
// types are available:
//
//	import _ "github.com/kubeflow/model-registry/internal/platform/db/drivers"
package drivers

import (
	_ "github.com/kubeflow/model-registry/internal/platform/db/mysql"
	_ "github.com/kubeflow/model-registry/internal/platform/db/postgres"
)
