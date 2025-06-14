package service

import (
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"gorm.io/gorm"
)

type TypeRepositoryImpl struct {
	db *gorm.DB
}

func NewTypeRepository(db *gorm.DB) models.TypeRepository {
	return &TypeRepositoryImpl{db: db}
}

func (r *TypeRepositoryImpl) GetAll() ([]models.Type, error) {
	var types []schema.Type

	if err := r.db.Find(&types).Error; err != nil {
		return nil, err
	}

	typesModels := make([]models.Type, len(types))

	for i, t := range types {
		typesModels[i] = &models.TypeImpl{
			ID: &t.ID,
			Attributes: &models.TypeAttributes{
				Name:        &t.Name,
				Version:     t.Version,
				TypeKind:    &t.TypeKind,
				Description: t.Description,
				InputType:   t.InputType,
				OutputType:  t.OutputType,
				ExternalID:  t.ExternalID,
			},
		}
	}

	return typesModels, nil
}
