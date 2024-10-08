package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultProvider  = testcontainers.ProviderDefault // or explicit to testcontainers.ProviderPodman if needed
	mlmdImage        = "gcr.io/tfx-oss-public/ml_metadata_store_server:1.14.0"
	sqliteFile       = "metadata.sqlite.db"
	testConfigFolder = "test/config/ml-metadata"
)

func ClearMetadataSqliteDB() error {
	wd, err := getTestConfigWorkingDir()
	if err != nil {
		return err
	}

	if err := os.Remove(fmt.Sprintf("%s/%s", wd, sqliteFile)); err != nil {
		return fmt.Errorf("expected to clear sqlite file but didn't find: %v", err)
	}
	return nil
}

func fileExists(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err == nil {
		return !info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func getTestConfigWorkingDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/../../%s", wd, testConfigFolder), nil
}

// SetupMLMetadataTestContainer setup a test container for MLMD server exposing gRPC interface
// Returns:
//   - The test container gRPC address <host>:<port>
//   - The teardown function to close and teardown the test container
func SetupMLMetadataTestContainer(t *testing.T) (*grpc.ClientConn, proto.MetadataStoreServiceClient, func(t *testing.T)) {
	ctx := context.Background()
	wd, err := getTestConfigWorkingDir()
	if err != nil {
		t.Errorf("error getting working directory: %v", err)
	}
	// t.Logf("using working directory: %s", wd)

	// when unhandled panics or other hard failures, could leave the DB in the directory
	// here we make sure it's not existing already, and that it was really cleanup by previous runs
	sqlitePath := fmt.Sprintf("%s/%s", wd, sqliteFile)
	exists, err := fileExists(sqlitePath)
	if err != nil {
		t.Errorf("error looking up for SQLite path: %v", err)
	}
	if exists {
		t.Errorf("SQLite should not exists: %v", sqlitePath)
		panic("halting immediately, SQLite should not exists: " + sqlitePath)
	}

	req := testcontainers.ContainerRequest{
		Image:        mlmdImage,
		ExposedPorts: []string{"8080/tcp"},
		Env: map[string]string{
			"METADATA_STORE_SERVER_CONFIG_FILE": "/tmp/shared/conn_config.pb",
		},
		Mounts: testcontainers.ContainerMounts{
			testcontainers.ContainerMount{
				Source: testcontainers.GenericBindMountSource{ // nolint keep deprecated method to avoid depending directly to docker api exposed by testcontainers' HostConfigModifier
					HostPath: wd,
				},
				Target: "/tmp/shared",
			},
		},
		WaitingFor: wait.ForLog("Server listening on"),
		Privileged: true,
	}

	useProvider := defaultProvider
	if useProvider == testcontainers.ProviderDefault { // user did not override
		if tryDetectPodmanRunning() {
			t.Log("Podman running detected! Will instruct TestContainers to use Podman explicitly.") // see https://github.com/testcontainers/testcontainers-go/issues/2264#issuecomment-2062092012
			useProvider = testcontainers.ProviderPodman
		}
	}
	mlmdgrpc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ProviderType:     useProvider,
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Errorf("error setting up mlmd grpc container: %v", err)
	}

	mappedHost, err := mlmdgrpc.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	mappedPort, err := mlmdgrpc.MappedPort(ctx, "8080")
	if err != nil {
		t.Error(err)
	}

	mlmdAddr := fmt.Sprintf("%s:%s", mappedHost, mappedPort.Port())
	t.Log("MLMD test container running at: ", mlmdAddr)

	// setup grpc connection
	conn, err := grpc.DialContext( // nolint:staticcheck
		context.Background(),
		mlmdAddr,
		grpc.WithReturnConnectionError(), // nolint:staticcheck
		grpc.WithBlock(),                 // nolint:staticcheck
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Errorf("error dialing connection to mlmd server %s: %v", mlmdAddr, err)
	}

	mlmdClient := proto.NewMetadataStoreServiceClient(conn)

	return conn, mlmdClient, func(t *testing.T) {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
		if err := mlmdgrpc.Terminate(ctx); err != nil {
			t.Error(err)
		}
	}
}

// simple utility function with heuristic to detect if Podman is running
func tryDetectPodmanRunning() bool {
	cmd := exec.Command("podman", "machine", "info", "--format", "json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}
	output := out.Bytes()
	type MachineInfo struct {
		Host struct {
			MachineState string `json:"MachineState"`
		} `json:"Host"`
	}
	var info MachineInfo
	err = json.Unmarshal(output, &info)
	if err != nil {
		return false
	}
	return info.Host.MachineState == "Running" || strings.Contains(os.Getenv("DOCKER_HOST"), "podman.sock")
}
