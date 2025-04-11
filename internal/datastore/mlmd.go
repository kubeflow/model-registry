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
	// maxGRPCRetryAttempts is the maximum number of attempts to retry GRPC requests to the MLMD server.
	maxGRPCRetryAttempts = 1 // 25 attempts with incremental backoff (1s, 2s, 3s, ..., 25s) it's ~5 minutes
)

var (
	ErrMLMDConnectionStart = errors.New("error dialing connection to mlmd server")
	ErrMLMDTypeCreation    = errors.New("error creating MLMD types")
	ErrMLMDCoreCreation    = errors.New("error creating core service")
	ErrMLMDConnectionClose = errors.New("error closing connection to mlmd server")
)

type MLMDService struct {
	Hostname       string
	Port           int
	GRPCConnection *grpc.ClientConn
}

func NewMLMDService(hostname string, port int) *MLMDService {
	return &MLMDService{
		Hostname: hostname,
		Port:     port,
	}
}

func (s *MLMDService) Build() (api.ModelRegistryApi, error) {
	mlmdAddr := fmt.Sprintf("%s:%d", s.Hostname, s.Port)

	glog.Infof("connecting to MLMD server %s..", mlmdAddr)

	conn, err := grpc.NewClient(mlmdAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%w %s: %w", ErrMLMDConnectionStart, mlmdAddr, err)
	}

	s.GRPCConnection = conn

	mlmdTypeNamesConfig := mlmdtypes.NewMLMDTypeNamesConfigFromDefaults()

	// Backoff and retry GRPC requests to the MLMD server, until the server
	// becomes available or the maximum number of attempts is reached.
	for i := 0; i < maxGRPCRetryAttempts; i++ {
		_, err = mlmdtypes.CreateMLMDTypes(conn, mlmdTypeNamesConfig)
		if err == nil {
			break
		}

		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Unavailable {
			return nil, fmt.Errorf("%w: %w", ErrMLMDTypeCreation, err)
		}

		time.Sleep(time.Duration(i+1) * time.Second)
	}

	service, err := core.NewModelRegistryService(conn, mlmdTypeNamesConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMLMDCoreCreation, err)
	}

	return service, nil
}

func (s *MLMDService) Teardown() error {
	if err := s.GRPCConnection.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrMLMDConnectionClose, err)
	}

	return nil
}
