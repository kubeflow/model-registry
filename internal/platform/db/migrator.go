package db

import (
	"fmt"
	"slices"
	"sync"

	"gorm.io/gorm"
)

type DBMigrator interface {
	Migrate() error
	Up(steps *int) error
	Down(steps *int) error
}

type MigratorFactory func(db *gorm.DB) (DBMigrator, error)

var (
	migratorFactories   = make(map[string]MigratorFactory)
	migratorFactoriesMu sync.RWMutex
)

func RegisterMigratorFactory(dbType string, factory MigratorFactory) {
	migratorFactoriesMu.Lock()
	defer migratorFactoriesMu.Unlock()

	if _, exists := migratorFactories[dbType]; exists {
		panic(fmt.Sprintf("duplicate migrator factory for database type %q", dbType))
	}

	migratorFactories[dbType] = factory
}

func NewDBMigrator(db *gorm.DB) (DBMigrator, error) {
	dbType := db.Name()

	migratorFactoriesMu.RLock()
	factory, ok := migratorFactories[dbType]
	migratorFactoriesMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unsupported database type: %s. Registered types: %v", dbType, registeredMigratorTypes())
	}

	return factory(db)
}

func registeredMigratorTypes() []string {
	migratorFactoriesMu.RLock()
	defer migratorFactoriesMu.RUnlock()

	types := make([]string, 0, len(migratorFactories))
	for dbType := range migratorFactories {
		types = append(types, dbType)
	}
	slices.Sort(types)
	return types
}
