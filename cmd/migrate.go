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
	"fmt"
	"github.com/dhirajsb/ml-metadata-go-server/model/db"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
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
		if err != nil {
			err = fmt.Errorf("db connection failed: %w", err)
			glog.Error(err)
			return err
		}
		// migrate all DB types
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
			err = fmt.Errorf("db migration failed: %w", err)
			glog.Error(err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	migrateCmd.Flags().StringVar(&dbFile, "db-file", "metadata.sqlite.db", "Sqlite DB file")
}
