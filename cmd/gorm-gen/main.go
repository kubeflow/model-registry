package main

import (
	"flag"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kubeflow/model-registry/internal/db/gen"
)

func main() {
	var (
		dsn    string
		dbType string
	)

	flag.StringVar(&dsn, "dsn", "", "Database connection string")
	flag.StringVar(&dbType, "db", "", "Database type (mysql or postgres)")
	flag.Parse()

	if dsn == "" || dbType == "" {
		fmt.Println("Error: --dsn and --db flags are required")
		os.Exit(1)
	}

	db, err := connectDB(dsn, dbType)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	gen.GenerateModel(db)

	fmt.Println("Successfully generated models")
}

func connectDB(dsn string, dbType string) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch dbType {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	return gorm.Open(dialector, &gorm.Config{})
} 