package postgres

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrations embed.FS

const (
	MigrationDir = "migrations"
)

type PostgresMigrator struct {
	migrator *migrate.Migrate
}

func NewPostgresMigrator(db *gorm.DB) (*PostgresMigrator, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
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
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}


	return &PostgresMigrator{
		migrator: m,
	}, nil
}

func (m *PostgresMigrator) Migrate() error {
	if err := m.Up(nil); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (m *PostgresMigrator) Up(steps *int) error {
	if steps == nil {
		return m.migrator.Up()
	}

	if *steps < 0 {
		return fmt.Errorf("steps cannot be negative")
	}

	return m.migrator.Steps(*steps)
}

func (m *PostgresMigrator) Down(steps *int) error {
	if steps == nil {
		return m.migrator.Down()
	}

	if *steps > 0 {
		return fmt.Errorf("steps cannot be positive")
	}

	return m.migrator.Steps(*steps)
} 