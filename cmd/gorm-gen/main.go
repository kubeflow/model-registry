package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kubeflow/model-registry/internal/db/gen"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Parse command line flags
	dsn := flag.String("dsn", "", "Database connection string")
	dbType := flag.String("db", "mysql", "Database type (mysql or postgres)")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("DSN is required")
	}

	// Connect to database
	var db *gorm.DB
	var err error

	switch *dbType {
	case "mysql":
		db, err = gorm.Open(mysql.Open(*dsn), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(*dsn), &gorm.Config{})
	default:
		log.Fatalf("Unsupported database type: %s", *dbType)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Call the central generator logic
	fmt.Println("Generating GORM models...")
	gen.GenerateModel(db)
	fmt.Println("GORM models generated successfully!")
} 