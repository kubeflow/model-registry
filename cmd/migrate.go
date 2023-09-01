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
	"github.com/dhirajsb/ml-metadata-go-server/model/db"
	"github.com/dhirajsb/ml-metadata-go-server/server"
	"github.com/dhirajsb/ml-metadata-go-server/server/library"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate ml-metadata DB to latest schema",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// connect to DB
		dbConn, err := NewDatabaseConnection(dbFile)
		defer func() {
			// close DB connection on exit
			db, err2 := dbConn.DB()
			if err2 != nil {
				err2 = db.Close()
				if err2 != nil {
					glog.Warningf("error closing DB connection: %v", err2)
				}
			}
		}()
		if err != nil {
			return fmt.Errorf("db connection failed: %w", err)
		}
		// migrate all DB types
		// TODO add support for more elaborate Gorm migrations
		err = dbConn.AutoMigrate(
			db.Artifact{},
			db.ArtifactProperty{},
			db.Association{},
			db.Attribution{},
			db.Context{},
			db.ContextProperty{},
			db.Event{},
			db.EventPath{},
			db.Execution{},
			db.ExecutionProperty{},
			db.ParentContext{},
			db.ParentType{},
			db.Type{},
			db.TypeProperty{},
		)
		if err != nil {
			return fmt.Errorf("db migration failed: %w", err)
		}

		// load metadata type libraries
		for _, dir := range libraryDirs {
			abs, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("error getting absolute library path for %s: %w", dir, err)
			}
			_, err = os.Stat(abs)
			if err != nil {
				return fmt.Errorf("error opening library path for %s: %w", abs, err)
			}
			err = filepath.WalkDir(abs, func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					glog.Warningf("error reading library path %s: %v", path, err)
					return filepath.SkipDir
				}
				if entry.IsDir() || !isYamlFile(path) {
					return nil
				}

				bytes, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read library file %s: %w", path, err)
				}
				var lib library.MetadataLibrary
				err = yaml.Unmarshal(bytes, &lib)
				grpcServer := server.NewGrpcServer(dbConn)
				typesRequest := proto.PutTypesRequest{}
				for _, ar := range lib.ArtifactTypes {
					typesRequest.ArtifactTypes = append(typesRequest.ArtifactTypes, &proto.ArtifactType{
						Name:        ar.Name,
						Version:     ar.Version,
						Description: ar.Description,
						ExternalId:  ar.ExternalId,
						Properties:  ar.Properties,
					})
				}
				for _, ar := range lib.ContextTypes {
					typesRequest.ContextTypes = append(typesRequest.ContextTypes, &proto.ContextType{
						Name:        ar.Name,
						Version:     ar.Version,
						Description: ar.Description,
						ExternalId:  ar.ExternalId,
						Properties:  ar.Properties,
					})
				}
				for _, ar := range lib.ExecutionTypes {
					typesRequest.ExecutionTypes = append(typesRequest.ExecutionTypes, &proto.ExecutionType{
						Name:        ar.Name,
						Version:     ar.Version,
						Description: ar.Description,
						ExternalId:  ar.ExternalId,
						Properties:  ar.Properties,
					})
				}
				response, err := grpcServer.PutTypes(context.Background(), &typesRequest)
				if err != nil {
					return fmt.Errorf("failed to add library from file %s: %w", path, err)
				}
				glog.Infof("added %d artifacts, %d contexts and %d execution types from library file %s",
					len(response.ArtifactTypeIds), len(response.ContextTypeIds), len(response.ExecutionTypeIds), path)
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to read library directory %s: %w", abs, err)
			}
		}
		return nil
	},
}

func isYamlFile(path string) bool {
	lowerPath := strings.ToLower(filepath.Ext(path))
	return strings.HasSuffix(lowerPath, ".yaml") || strings.HasSuffix(lowerPath, ".yml")
}

var libraryDirs []string

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	migrateCmd.Flags().StringVarP(&dbFile, "db-file", "d", "metadata.sqlite.db", "Sqlite DB file")
	migrateCmd.Flags().StringSliceVarP(&libraryDirs, "metadata-library-dir", "m", libraryDirs, "Built-in metadata types library directories containing yaml files")
}
