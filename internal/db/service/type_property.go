package service

import (
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/platform/db/repository"
	"gorm.io/gorm"
)

func NewTypePropertyRepository(db *gorm.DB) models.TypePropertyRepository {
	return repository.NewTypePropertyRepository(db)
}
