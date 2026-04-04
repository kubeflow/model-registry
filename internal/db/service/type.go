package service

import (
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/platform/db/repository"
	"gorm.io/gorm"
)

type TypeRepositoryImpl = repository.TypeRepositoryImpl

func NewTypeRepository(db *gorm.DB) models.TypeRepository {
	return repository.NewTypeRepository(db)
}
