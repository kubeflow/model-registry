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
	"github.com/dhirajsb/ml-metadata-go-server/ml_metadata/proto"
	"github.com/dhirajsb/ml-metadata-go-server/server"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"sync"
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
	dbFile   string
	grpcHost     = "localhost"
	grpcPort int = 8080

	// serveCmd represents the serve command
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("serve called")

			// connect to the DB using Gorm
			db, err := NewDatabaseConnection(dbFile)
			if err != nil {
				log.Fatalf("db connection failed: %v", err)
			}

			// serve the grpc server
			listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcHost, grpcPort))
			if err != nil {
				log.Fatalf("grpc listen failed: %v", err)
			}
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
			proto.RegisterMetadataStoreServiceServer(grpcServer, server.NewGrpcServer(db))

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				glog.Info("Starting grpc server...")
				err = grpcServer.Serve(listen)
				if err != nil {
					log.Fatalf("grpc serving failed: %v", err)
				}
				wg.Done()
			}()

			// TODO serve the GraphQL server

			// wait for servers to finish
			wg.Wait()

			return nil
		},
	}
)

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
	serveCmd.Flags().StringVar(&dbFile, "db-file", "metadata.sqlite.db", "Sqlite DB file")
	serveCmd.Flags().StringVar(&grpcHost, "grpc-host", grpcHost, "gRPC listen hostname")
	serveCmd.Flags().IntVar(&grpcPort, "grpc-port", grpcPort, "gRPC listen port")
}
