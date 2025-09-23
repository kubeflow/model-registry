package sqlite

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrations embed.FS

const (
	MigrationDir = "migrations"
)

type SQLiteMigrator struct {
	migrator *migrate.Migrate
}

func NewSQLiteMigrator(db *gorm.DB) (*SQLiteMigrator, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}

	// Create a new source instance from the embedded files
	source, err := iofs.New(migrations, MigrationDir)
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		source,
		"sqlite3",
		driver,
	)
	if err != nil {
		return nil, err
	}

	return &SQLiteMigrator{
		migrator: m,
	}, nil
}

func (m *SQLiteMigrator) Migrate() error {
	if err := m.Up(nil); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (m *SQLiteMigrator) Up(steps *int) error {
	if steps == nil {
		return m.migrator.Up()
	}

	if *steps < 0 {
		return fmt.Errorf("steps cannot be negative")
	}

	return m.migrator.Steps(*steps)
}

func (m *SQLiteMigrator) Down(steps *int) error {
	if steps == nil {
		return m.migrator.Down()
	}

	if *steps > 0 {
		return fmt.Errorf("steps cannot be positive")
	}

	return m.migrator.Steps(*steps)
}