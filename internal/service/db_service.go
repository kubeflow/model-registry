package service

import (
	"gorm.io/gorm"

	"github.com/opendatahub-io/model-registry/internal/model/db"
)

var _ DBService = dbServiceHandler{}
var _ DBService = (*dbServiceHandler)(nil)

func NewDBService(db *gorm.DB) DBService {
	return &dbServiceHandler{
		typeHandler:     &typeHandler{db: db},
		artifactHandler: &artifactHandler{db: db},
	}
}

type DBService interface {
	InsertType(db.Type) (*db.Type, error)
	UpsertType(db.Type) (*db.Type, error)
	ReadType(db.Type) (*db.Type, error)
	// Get-like function to use a signature similar to the gorm `Where` func
	ReadAllType(query interface{}, args ...interface{}) ([]*db.Type, error)
	UpdateType(db.Type) (*db.Type, error)
	DeleteType(db.Type) (*db.Type, error)

	// InsertEEE(db.EEE) (*db.EEE, error)
	// UpsertEEE(db.EEE) (*db.EEE, error)
	// ReadEEE(db.EEE) (*db.EEE, error)
	// ReadAllEEE(query interface{}, args ...interface{}) ([]*db.EEE, error)
	// UpdateEEE(db.EEE) (*db.EEE, error)
	// DeleteEEE(db.EEE) (*db.EEE, error)
}

type dbServiceHandler struct {
	*typeHandler
	*artifactHandler
}
