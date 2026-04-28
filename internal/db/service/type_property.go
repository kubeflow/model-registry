package service

import (
	"github.com/kubeflow/hub/internal/db/models"
	"github.com/kubeflow/hub/internal/platform/db/repository"
	"gorm.io/gorm"
)

func NewTypePropertyRepository(db *gorm.DB) models.TypePropertyRepository {
	return repository.NewTypePropertyRepository(db)
}
