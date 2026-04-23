package service

import (
	"os"
	"testing"

	"github.com/kubeflow/hub/internal/testutils"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.TestMainPostgresHelper(m))
}
