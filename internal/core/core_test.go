package core_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/core"
	"github.com/opendatahub-io/model-registry/internal/core/mapper"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"github.com/stretchr/testify/assert"
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

func TestCreateRegisteredModel(t *testing.T) {
	conn, client, teardown := SetupTestContainer(t)
	defer teardown(t)

	// [TEST CASE]

	// create mode registry service
	service, err := core.NewModelRegistryService(conn)
	assert.Nilf(t, err, "error creating core service: %v", err)

	modelName := "PricingModel"
	externalId := "myExternalId"
	owner := "Myself"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &externalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	assert.Nilf(t, err, "error creating registered model: %v", err)
	assert.NotNilf(t, createdModel.Id, "created registered model should not have nil Id")

	byTypeAndNameResp, err := client.GetContextByTypeAndName(context.Background(), &proto.GetContextByTypeAndNameRequest{
		TypeName:    &core.RegisteredModelTypeName,
		ContextName: &modelName,
	})
	assert.Nilf(t, err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctxId := mapper.IdToString(*byTypeAndNameResp.Context.Id)
	assert.Equal(t, *createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	assert.Equal(t, modelName, *byTypeAndNameResp.Context.Name, "saved model name should match the provided one")
	assert.Equal(t, externalId, *byTypeAndNameResp.Context.ExternalId, "saved external id should match the provided one")
	assert.Equal(t, owner, byTypeAndNameResp.Context.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assert.Nilf(t, err, "error retrieving all contexts, not related to the test itself: %v", err)
	assert.Equal(t, 1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")
}

func TestGetRegisteredModelById(t *testing.T) {
	conn, _, teardown := SetupTestContainer(t)
	defer teardown(t)

	// [TEST CASE]

	// create mode registry service
	service, err := core.NewModelRegistryService(conn)
	assert.Nilf(t, err, "error creating core service: %v", err)

	modelName := "PricingModel"
	externalId := "mysupermodel"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &externalId,
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	assert.Nilf(t, err, "error creating registered model: %v", err)

	modelId, _ := mapper.IdToInt64(*createdModel.Id)
	getModelById, err := service.GetRegisteredModelById((*core.BaseResourceId)(modelId))
	assert.Nilf(t, err, "error getting registered model by id %d: %v", *modelId, err)

	assert.Equal(t, modelName, *getModelById.Name, "saved model name should match the provided one")
	assert.Equal(t, externalId, *getModelById.ExternalID, "saved external id should match the provided one")
}

// #################
// ##### Utils #####
// #################

func clearMetadataSqliteDB(wd string) error {
	if err := os.Remove(fmt.Sprintf("%s/%s", wd, sqliteFile)); err != nil {
		return fmt.Errorf("expected to clear sqlite file but didn't find: %v", err)
	}
	return nil
}

// SetupTestContainer creates a MLMD gRPC test container
// Returns
//   - gRPC client connection to the test container
//   - ml-metadata client used to double check the database
//   - teardown function
func SetupTestContainer(t *testing.T) (*grpc.ClientConn, proto.MetadataStoreServiceClient, func(t *testing.T)) {
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
