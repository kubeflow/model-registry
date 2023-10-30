package testutils

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	useProvider      = testcontainers.ProviderDefault // or explicit to testcontainers.ProviderPodman if needed
	mlmdImage        = "gcr.io/tfx-oss-public/ml_metadata_store_server:1.14.0"
	sqliteFile       = "metadata.sqlite.db"
	testConfigFolder = "test/config/ml-metadata"
)

func clearMetadataSqliteDB(wd string) error {
	if err := os.Remove(fmt.Sprintf("%s/%s", wd, sqliteFile)); err != nil {
		return fmt.Errorf("expected to clear sqlite file but didn't find: %v", err)
	}
	return nil
}

// SetupMLMDTestContainer creates a MLMD gRPC test container
// Returns
//   - gRPC client connection to the test container
//   - ml-metadata client used to double check the database
//   - teardown function
func SetupMLMDTestContainer(t *testing.T) (*grpc.ClientConn, proto.MetadataStoreServiceClient, func(t *testing.T)) {
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("error getting working directory: %v", err)
	}
	wd = fmt.Sprintf("%s/../../%s", wd, testConfigFolder)
	t.Logf("using working directory: %s", wd)

	req := testcontainers.ContainerRequest{
		Image:        mlmdImage,
		ExposedPorts: []string{"8080/tcp"},
		Env: map[string]string{
			"METADATA_STORE_SERVER_CONFIG_FILE": "/tmp/shared/conn_config.pb",
		},
		Mounts: testcontainers.ContainerMounts{
			testcontainers.ContainerMount{
				Source: testcontainers.GenericBindMountSource{
					HostPath: wd,
				},
				Target: "/tmp/shared",
			},
		},
		WaitingFor: wait.ForLog("Server listening on"),
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
	t.Log("MLMD test container setup at: ", mlmdAddr)

	// setup grpc connection
	conn, err := grpc.DialContext(
		context.Background(),
		mlmdAddr,
		grpc.WithReturnConnectionError(),
		grpc.WithBlock(),
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
		if err := clearMetadataSqliteDB(wd); err != nil {
			t.Error(err)
		}
	}
}
