/*
Copyright © 2023 Dhiraj Bokde dhirajsb@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	grpc2 "github.com/opendatahub-io/model-registry/internal/server/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/server/graph"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InterceptorLogger(l *log.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			msg = fmt.Sprintf("DEBUG :%v", msg)
		case logging.LevelInfo:
			msg = fmt.Sprintf("INFO :%v", msg)
		case logging.LevelWarn:
			msg = fmt.Sprintf("WARN :%v", msg)
		case logging.LevelError:
			msg = fmt.Sprintf("ERROR :%v", msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
		l.Println(append([]any{"msg", msg}, fields...))
	})
}

var (
	dbFile string
	host       = "localhost"
	port   int = 8080

	// serveCmd represents the serve command
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the ml-metadata go server",
		Long: `This command launches the ml-metadata go server.

The server connects to a SQlite database. It supports options to customize the 
location of the database file and the hostname and port where it listens.'`,
		RunE: runServer,
	}
)

func runServer(cmd *cobra.Command, args []string) error {
	glog.Info("server started...")

	// Create a channel to receive signals
	signalChannel := make(chan os.Signal, 1)

	// Notify the channel on SIGINT (Ctrl+C) and SIGTERM signals
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// connect to the DB using Gorm
	db, err := NewDatabaseConnection(dbFile)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	// listen on host:port
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("server listen failed: %v", err)
	}
	m := cmux.New(listener)
	gqlListener := m.Match(cmux.HTTP1Fast())
	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))

	// gRPC server
	grpcServer := grpcListenerServer(grpc2.NewGrpcServer(db))
	// GraphQL server
	gqlServer := graphQlListenerServer(db)

	// start cmux listeners
	g := new(errgroup.Group)
	g.Go(func() error {
		glog.Info("starting gRPC server...")
		return grpcServer.Serve(grpcListener)
	})
	g.Go(func() error {
		glog.Info("starting GraphQL server...")
		return gqlServer.Serve(gqlListener)
	})
	g.Go(func() error {
		return m.Serve()
	})

	go func() {
		err = g.Wait()
		// error starting server
		if err != nil || err != http.ErrServerClosed || err != grpc.ErrServerStopped || err != cmux.ErrServerClosed {
			glog.Errorf("server listener error: %v", err)
		}
		signalChannel <- syscall.SIGINT
	}()

	// Wait for a signal
	receivedSignal := <-signalChannel
	glog.Infof("received signal: %s\n", receivedSignal)

	// Perform cleanup or other graceful shutdown actions here
	glog.Info("shutting down services...")
	grpcServer.Stop()
	_ = gqlServer.Shutdown(context.Background())

	// close DB
	glog.Info("closing DB...")
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("error accessing DB: %v", err)
	}
	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("error closing DB: %v", err)
	}
	glog.Info("shutdown!")
	return nil
}

func graphQlListenerServer(db *gorm.DB) *http.Server {
	mux := http.NewServeMux()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	return &http.Server{Handler: mux}
}

func grpcListenerServer(server proto.MetadataStoreServiceServer) *grpc.Server {
	// TODO map server options from flags
	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	lopts := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent, logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), lopts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(logger), lopts...),
		),
	}

	grpcServer := grpc.NewServer(opts...)
	// simple health check
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	proto.RegisterMetadataStoreServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	return grpcServer
}

func NewDatabaseConnection(dbFile string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serveCmd.Flags().StringVarP(&dbFile, "db-file", "d", "metadata.sqlite.db", "Sqlite DB file")
	serveCmd.Flags().StringVarP(&host, "hostname", "n", host, "Server listen hostname")
	serveCmd.Flags().IntVarP(&port, "port", "p", port, "Server listen port")
}
