package datastore

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/mlmdtypes"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	// maxGRPCRetryAttempts is the maximum number of attempts to retry GRPC requests to the MLMD service.
	maxGRPCRetryAttempts = 25 // 25 attempts with incremental backoff (1s, 2s, 3s, ..., 25s) it's ~5 minutes
)

var (
	ErrMLMDConnectionStart = errors.New("error dialing connection to mlmd service")
	ErrMLMDTypeCreation    = errors.New("error creating MLMD types")
	ErrMLMDCoreCreation    = errors.New("error creating core service")
	ErrMLMDConnectionClose = errors.New("error closing connection to mlmd service")
)

type MLMDConfig struct {
	Hostname string
	Port     int
}

func (c *MLMDConfig) Validate() error {
	if c.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be in the range 1-65535")
	}

	return nil
}

type MLMDService struct {
	*MLMDConfig
	gRPCConnection *grpc.ClientConn
}

func NewMLMDService(cfg *MLMDConfig) *MLMDService {
	return &MLMDService{
		MLMDConfig: cfg,
	}
}

func (s *MLMDService) Connect() (api.ModelRegistryApi, error) {
	uri := fmt.Sprintf("%s:%d", s.Hostname, s.Port)

	glog.Infof("Connecting to MLMD service at %s..", uri)

	conn, err := grpc.NewClient(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%w %s: %w", ErrMLMDConnectionStart, uri, err)
	}

	s.gRPCConnection = conn

	mlmdTypeNamesConfig := mlmdtypes.NewMLMDTypeNamesConfigFromDefaults()

	// Backoff and retry GRPC requests to the MLMD service, until the service
	// becomes available or the maximum number of attempts is reached.
	for i := range maxGRPCRetryAttempts {
		_, err = mlmdtypes.CreateMLMDTypes(conn, mlmdTypeNamesConfig)
		if err == nil {
			break
		}

		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Unavailable {
			return nil, fmt.Errorf("%w: %w", ErrMLMDTypeCreation, err)
		}

		glog.Warningf("Retrying connection to MLMD service (attempt %d/%d): %v", i+1, maxGRPCRetryAttempts, err)

		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMLMDTypeCreation, err)
	}

	service, err := core.NewModelRegistryService(conn, mlmdTypeNamesConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMLMDCoreCreation, err)
	}

	glog.Infof("Successfully connected to MLMD service")

	return service, nil
}

func (s *MLMDService) Teardown() error {
	glog.Info("Closing connection to MLMD service")

	if s.gRPCConnection == nil {
		return nil
	}

	if err := s.gRPCConnection.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrMLMDConnectionClose, err)
	}

	return nil
}
