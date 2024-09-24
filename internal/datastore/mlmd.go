package datastore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/mlmdtypes"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*30)

	defer cancel()

	mlmdAddr := fmt.Sprintf("%s:%d", s.Hostname, s.Port)

	glog.Infof("connecting to MLMD server %s..", mlmdAddr)

	conn, err := grpc.DialContext( // nolint:staticcheck
		ctxTimeout,
		mlmdAddr,
		grpc.WithReturnConnectionError(), // nolint:staticcheck
		grpc.WithBlock(),                 // nolint:staticcheck
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w %s: %w", ErrMLMDConnectionStart, mlmdAddr, err)
	}

	s.GRPCConnection = conn

	glog.Infof("connected to MLMD server")

	mlmdTypeNamesConfig := mlmdtypes.NewMLMDTypeNamesConfigFromDefaults()

	if _, err = mlmdtypes.CreateMLMDTypes(conn, mlmdTypeNamesConfig); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMLMDTypeCreation, err)
	}

	service, err := core.NewModelRegistryService(conn, mlmdTypeNamesConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMLMDCoreCreation, err)
	}

	return service, nil
}

func (s *MLMDService) Teardown() error {
	glog.Infof("closing connection to MLMD server")

	if err := s.GRPCConnection.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrMLMDConnectionClose, err)
	}

	return nil
}
