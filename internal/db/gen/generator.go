package gen

import (
	"log"
	"strings"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func GenerateModel(db *gorm.DB) {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/db/schema",
		ModelPkgPath: "schema",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	g.WithDataTypeMap(map[string]func(column gorm.ColumnType) string{
		"bytea": func(column gorm.ColumnType) string {
			return "[]byte"
		},
	})

	g.WithOpts(
		gen.FieldGORMTag("*", func(tag field.GormTag) field.GormTag {
			if vals, ok := tag["default"]; ok {
				if len(vals) > 0 {
					val := strings.Trim(strings.TrimSpace(vals[0]), `"'`)
					// Check for NULL, 0, or empty string, which are all ways
					// a nullable default might be represented.
					if strings.ToUpper(val) == "NULL" || val == "0" || val == "" {
						tag.Remove("default")
					}
				}
			}
			return tag
		}),
	)

	tables, err := db.Migrator().GetTables()
	if err != nil {
		log.Fatalf("Failed to get tables: %v", err)
	}

	for _, tableName := range tables {
		if strings.EqualFold(tableName, "Type") {
			g.GenerateModel(tableName, gen.FieldType("type_kind", "int32"))
		} else {
			g.GenerateModel(tableName)
		}
	}

	g.Execute()
} 