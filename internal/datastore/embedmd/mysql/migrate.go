package mysql

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrations embed.FS

const (
	MigrationDir = "migrations"
)

type MySQLMigrator struct {
	migrator *migrate.Migrate
}

func NewMySQLMigrator(db *gorm.DB) (*MySQLMigrator, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
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
		"mysql",
		driver,
	)
	if err != nil {
		return nil, err
	}

	return &MySQLMigrator{
		migrator: m,
	}, nil
}

func (m *MySQLMigrator) Migrate() error {
	if err := m.Up(nil); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (m *MySQLMigrator) Up(steps *int) error {
	if steps == nil {
		return m.migrator.Up()
	}

	if *steps < 0 {
		return fmt.Errorf("steps cannot be negative")
	}

	return m.migrator.Steps(*steps)
}

func (m *MySQLMigrator) Down(steps *int) error {
	if steps == nil {
		return m.migrator.Down()
	}

	if *steps > 0 {
		return fmt.Errorf("steps cannot be positive")
	}

	return m.migrator.Steps(*steps)
}
