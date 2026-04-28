package service

import (
	"github.com/kubeflow/hub/internal/db/models"
	"github.com/kubeflow/hub/internal/platform/db/repository"
	"gorm.io/gorm"
)

type TypeRepositoryImpl = repository.TypeRepositoryImpl

func NewTypeRepository(db *gorm.DB) models.TypeRepository {
	return repository.NewTypeRepository(db)
}
