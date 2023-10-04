package service

import (
	"fmt"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/model/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func migrateDatabase(dbConn *gorm.DB) error {
	// using only needed RDBMS type for the scope under test
	err := dbConn.AutoMigrate(
		db.Type{},
		db.TypeProperty{},
		// TODO: add as needed.
	)
	if err != nil {
		return fmt.Errorf("db migration failed: %w", err)
	}
	return nil
}

func setup() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = migrateDatabase(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Bare minimal test of PutArtifactType with a given Name, and Get.
func TestInsertTypeThenReadAllType(t *testing.T) {
	dbc, err := setup()
	if err != nil {
		t.Errorf("Should expect DB connection: %v", err)
	}
	defer func() {
		dbi, err := dbc.DB()
		if err != nil {
			t.Errorf("Test need to clear sqlite DB for the next one, but errored: %v", err)
		}
		dbi.Close()
	}()
	dal := NewDBService(dbc)

	artifactName := "John Doe"
	newType := db.Type{
		Name:     artifactName,
		TypeKind: int8(db.ARTIFACT_TYPE),
	}

	at, err := dal.InsertType(newType)
	if err != nil {
		t.Errorf("Should create ArtifactType: %v", err)
	}
	if at.ID < 0 {
		t.Errorf("Should have ID for ArtifactType: %v", at.ID)
	}
	if at.Name != artifactName {
		t.Errorf("Should have Name for ArtifactType per constant: %v", at.Name)
	}

	ats, err2 := dal.ReadAllType(newType)
	if err2 != nil {
		t.Errorf("Should get ArtifactType: %v", err2)
	}
	if len(ats) != 1 { // TODO if temp file is okay, this is superfluos
		t.Errorf("The test is running under different assumption")
	}
	at0 := ats[0]
	t.Logf("at0: %v", at0)
	if at0.ID != at.ID {
		t.Errorf("Should have same ID")
	}
	if at0.Name != at.Name {
		t.Errorf("Should have same Name")
	}

}

func TestReadAllType(t *testing.T) {
	dbc, err := setup()
	if err != nil {
		t.Errorf("Should expect DB connection: %v", err)
	}
	defer func() {
		dbi, err := dbc.DB()
		if err != nil {
			t.Errorf("Test need to clear sqlite DB for the next one, but errored: %v", err)
		}
		dbi.Close()
	}()
	dal := NewDBService(dbc)

	fixVersion := "version"

	if _, err := dal.InsertType(db.Type{Name: "at0", Version: &fixVersion, TypeKind: int8(db.ARTIFACT_TYPE)}); err != nil {
		t.Errorf("Should create ArtifactType: %v", err)
	}
	if _, err := dal.InsertType(db.Type{Name: "at1", Version: &fixVersion, TypeKind: int8(db.ARTIFACT_TYPE)}); err != nil {
		t.Errorf("Should create ArtifactType: %v", err)
	}

	results, err := dal.ReadAllType(db.Type{Version: &fixVersion})
	t.Logf("results: %v", results)
	if err != nil {
		t.Errorf("Should get ArtifactTypes: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Should have retrieved 2 artifactTypes")
	}
}

func TestUpsertType(t *testing.T) {
	dbc, err := setup()
	if err != nil {
		t.Errorf("Should expect DB connection: %v", err)
	}
	defer func() {
		dbi, err := dbc.DB()
		if err != nil {
			t.Errorf("Test need to clear sqlite DB for the next one, but errored: %v", err)
		}
		dbi.Close()
	}()
	dal := NewDBService(dbc)

	artifactName := "John Doe"
	v0 := "v0"
	v1 := "v1"
	if _, err := dal.InsertType(db.Type{Name: artifactName, Version: &v0, TypeKind: int8(db.ARTIFACT_TYPE)}); err != nil {
		t.Errorf("Should Insert ArtifactType: %v", err)
	}
	if res, err := dal.InsertType(db.Type{Name: artifactName, Version: &v0, TypeKind: int8(db.ARTIFACT_TYPE)}); err == nil {
		t.Errorf("Subsequent Insert must have failed: %v", res)
	}
	if _, err := dal.UpsertType(db.Type{Name: artifactName, Version: &v1, TypeKind: int8(db.ARTIFACT_TYPE)}); err != nil {
		t.Errorf("Should Upsert ArtifactType: %v", err)
	}
}
