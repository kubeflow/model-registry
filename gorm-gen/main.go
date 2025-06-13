package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// genModels is gorm/gen generated models
func genModels(g *gen.Generator, db *gorm.DB, tables []string) (err error) {
	if len(tables) == 0 {
		// Execute tasks for all tables in the database
		tables, err = db.Migrator().GetTables()
		if err != nil {
			return fmt.Errorf("GORM migrator get all tables fail: %w", err)
		}
	}

	// Execute some data table tasks
	for _, tableName := range tables {
		if tableName == "Type" {
			// Special handling for Type table to set TypeKind as int32
			g.GenerateModel(tableName, gen.FieldType("type_kind", "int32"))
		} else {
			g.GenerateModel(tableName)
		}
	}
	return nil
}

func main() {
	// Database connection configuration
	dsn := "root:root@tcp(localhost:3306)/model-registry?charset=utf8mb4&parseTime=True&loc=Local"

	// Allow DSN override via environment variable
	if envDSN := os.Getenv("GORM_GEN_DSN"); envDSN != "" {
		dsn = envDSN
	}

	// Connect to database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize the generator with configuration for models only
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../internal/db/schema",
		ModelPkgPath:      "schema",
		Mode:              0,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  false,
	})

	// Use the database connection
	g.UseDB(db)

	// Generate models for all tables using custom function
	err = genModels(g, db, nil)
	if err != nil {
		log.Fatalf("Failed to generate models: %v", err)
	}

	// Generate the code
	fmt.Println("Generating GORM models (structs only)...")
	g.Execute()
	fmt.Println("GORM models generated successfully!")
}
