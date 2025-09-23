package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

var (
	dbType string
	dsn    string
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

	// Custom ModelOpt to clean up primary key tags
	modelOpt := gen.FieldGORMTag("*", func(tag field.GormTag) field.GormTag {
		// Remove "not null" from primary key fields since it's redundant
		if _, hasPrimaryKey := tag["primaryKey"]; hasPrimaryKey {
			tag.Remove("not null")
		}
		return tag
	})

	// Execute some data table tasks
	for _, tableName := range tables {
		if tableName == "Type" {
			// Special handling for Type table
			if dbType == "sqlite" {
				g.GenerateModel(tableName, 
					gen.FieldType("type_kind", "int32"), 
					gen.FieldType("id", "int32"),
					gen.FieldGORMTag("id", func(tag field.GormTag) field.GormTag {
						tag.Set("column", "id")
						tag.Set("primaryKey", "")
						tag.Set("autoIncrement", "true")
						return tag
					}),
					modelOpt)
			} else {
				g.GenerateModel(tableName, 
					gen.FieldType("type_kind", "int32"), 
					gen.FieldGORMTag("id", func(tag field.GormTag) field.GormTag {
						tag.Set("column", "id")
						tag.Set("primaryKey", "")
						tag.Set("autoIncrement", "true")
						return tag
					}),
					modelOpt)
			}
		} else if tableName == "schema_migrations" && dbType == "sqlite" {
			// Special handling for schema_migrations table in SQLite to match MySQL/PostgreSQL types
			g.GenerateModel(tableName, 
				gen.FieldType("version", "int64"),
				gen.FieldType("dirty", "bool"), 
				gen.FieldGORMTag("version", func(tag field.GormTag) field.GormTag {
					tag.Set("column", "version")
					tag.Set("primaryKey", "")
					return tag
				}),
				gen.FieldGORMTag("dirty", func(tag field.GormTag) field.GormTag {
					tag.Set("column", "dirty")
					tag.Set("not null", "")
					return tag
				}),
				modelOpt)
		} else if isPropertyTable(tableName) {
			// Special handling for property tables to ensure composite primary keys
			var opts []gen.ModelOpt
			if dbType == "sqlite" && hasIDField(tableName) {
				opts = append(opts, gen.FieldType("id", "int32"))
				opts = append(opts, gen.FieldGORMTag("id", func(tag field.GormTag) field.GormTag {
					tag.Set("column", "id")
					tag.Set("primaryKey", "")
					tag.Set("autoIncrement", "true")
					return tag
				}))
			}
			
			// Fix composite primary key fields
			opts = append(opts, gen.FieldGORMTag("name", func(tag field.GormTag) field.GormTag {
				tag.Set("column", "name")
				tag.Set("primaryKey", "")
				tag.Remove("not null") // Remove not null since primaryKey implies it
				return tag
			}))
			opts = append(opts, gen.FieldGORMTag("is_custom_property", func(tag field.GormTag) field.GormTag {
				tag.Set("column", "is_custom_property")
				tag.Set("primaryKey", "")
				tag.Remove("not null") // Remove not null since primaryKey implies it
				return tag
			}))
			opts = append(opts, modelOpt)
			g.GenerateModel(tableName, opts...)
		} else if hasIDField(tableName) {
			// Tables with ID fields - ensure autoIncrement is set
			var opts []gen.ModelOpt
			if dbType == "sqlite" {
				opts = append(opts, gen.FieldType("id", "int32"))
			}
			opts = append(opts, gen.FieldGORMTag("id", func(tag field.GormTag) field.GormTag {
				tag.Set("column", "id")
				tag.Set("primaryKey", "")
				tag.Set("autoIncrement", "true")
				return tag
			}))
			opts = append(opts, modelOpt)
			g.GenerateModel(tableName, opts...)
		} else {
			g.GenerateModel(tableName, modelOpt)
		}
	}
	return nil
}

// hasIDField returns true if the table has an auto-increment ID primary key field
func hasIDField(tableName string) bool {
	tablesWithIDField := []string{
		"Artifact", "Association", "Attribution", "Context", 
		"Event", "Execution", "Type",
	}
	
	for _, table := range tablesWithIDField {
		if table == tableName {
			return true
		}
	}
	return false
}

// isPropertyTable returns true if the table is a property table with composite primary keys
func isPropertyTable(tableName string) bool {
	propertyTables := []string{
		"ArtifactProperty", "ContextProperty", "ExecutionProperty", "TypeProperty",
	}
	
	for _, table := range propertyTables {
		if table == tableName {
			return true
		}
	}
	return false
}

// getDialector returns the appropriate GORM dialector based on database type and DSN
func getDialector(dbType, dsn string) (gorm.Dialector, error) {
	switch dbType {
	case "mysql":
		return mysql.Open(dsn), nil
	case "postgres", "postgresql":
		return postgres.Open(dsn), nil
	case "sqlite":
		return sqlite.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s. Supported types: mysql, postgres, sqlite, sqlserver", dbType)
	}
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gorm-gen",
	Short: "GORM code generator for model-registry database schemas",
	Long: `GORM code generator for model-registry database schemas.

This tool generates GORM model structs from database tables for the model-registry project.
It supports multiple database types including MySQL, PostgreSQL, SQLite, and SQL Server.

The generated models are placed in the ../internal/db/schema directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGenerate()
	},
}

func runGenerate() error {
	// Allow environment variable overrides
	if envDBType := os.Getenv("GORM_GEN_DB_TYPE"); envDBType != "" {
		dbType = envDBType
	}

	if envDSN := os.Getenv("GORM_GEN_DSN"); envDSN != "" {
		dsn = envDSN
	}

	// Use default DSN if not provided
	if dsn == "" {
		return fmt.Errorf("Please provide a DSN using --dsn flag or GORM_GEN_DSN environment variable for %s database", dbType)
	}

	fmt.Printf("Connecting to %s database...\n", dbType)

	// Get the appropriate dialector
	dialector, err := getDialector(dbType, dsn)
	if err != nil {
		return fmt.Errorf("failed to get database dialector: %w", err)
	}

	// Connect to database
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Database connection successful!")

	// Initialize the generator with configuration for models only
	config := gen.Config{
		OutPath:           "../internal/db/schema",
		ModelPkgPath:      "schema",
		Mode:              0,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  false,
	}

	if dbType == "postgres" || dbType == "postgresql" {
		dataTypeMap := map[string]func(gorm.ColumnType) string{
			"bytea": func(columnType gorm.ColumnType) string {
				return "[]byte"
			},
		}
		config.WithDataTypeMap(dataTypeMap)
	}
	g := gen.NewGenerator(config)

	// Use the database connection
	g.UseDB(db)

	// Generate models for all tables using custom function
	err = genModels(g, db, nil)
	if err != nil {
		return fmt.Errorf("failed to generate models: %w", err)
	}

	// Generate the code
	fmt.Printf("Generating GORM models for %s database...\n", dbType)
	g.Execute()
	fmt.Println("GORM models generated successfully!")

	return nil
}

func init() {
	// Define flags
	rootCmd.Flags().StringVar(&dbType, "db-type", "mysql", "Database type (mysql, postgres, sqlite, sqlserver)")
	rootCmd.Flags().StringVar(&dsn, "dsn", "", "Database connection string (DSN). If not provided, uses default for the database type")

	// Add examples to the help
	rootCmd.Example = `  # Generate models for MySQL (default)
  gorm-gen --db-type=mysql --dsn="user:pass@tcp(localhost:3306)/dbname"

  # Generate models for PostgreSQL
  gorm-gen --db-type=postgres --dsn="host=localhost user=postgres dbname=mydb"

  # Generate models for SQLite
  gorm-gen --db-type=sqlite --dsn="./database.db"

  # Use environment variables
  export GORM_GEN_DB_TYPE=postgres
  export GORM_GEN_DSN="host=localhost user=postgres dbname=mydb"
  gorm-gen`
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
