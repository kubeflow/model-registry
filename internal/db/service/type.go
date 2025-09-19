package service

import (
	"errors"
	"fmt"

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

func (r *TypeRepositoryImpl) Save(t models.Type) (models.Type, error) {
	attr := t.GetAttributes()
	if attr == nil {
		return t, errors.New("invalid type: missing attributes")
	}
	if attr.Name == nil {
		return t, errors.New("invalid type: missing name")
	}
	if attr.TypeKind == nil {
		return t, errors.New("invalid type: missing kind")
	}

	var st schema.Type
	err := r.db.Where("name = ?", *attr.Name).First(&st).Error

	if err == nil {
		// Record already exists. We don't support updates, but we can return the full details.

		// Catch this case in particular.
		if st.TypeKind != *attr.TypeKind {
			return t, fmt.Errorf("invalid type: kind is %d, cannot change to kind %d", st.TypeKind, *attr.TypeKind)
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Record doesn't exist, so we'll create it.
		st = schema.Type{
			Name:        *attr.Name,
			Version:     attr.Version,
			TypeKind:    *attr.TypeKind,
			Description: attr.Description,
			InputType:   attr.InputType,
			OutputType:  attr.OutputType,
			ExternalID:  attr.ExternalID,
		}

		if err := r.db.Create(&st).Error; err != nil {
			return t, err
		}
	} else {
		return t, err
	}

	return &models.TypeImpl{
		ID: &st.ID,
		Attributes: &models.TypeAttributes{
			Name:        &st.Name,
			Version:     st.Version,
			TypeKind:    &st.TypeKind,
			Description: st.Description,
			InputType:   st.InputType,
			OutputType:  st.OutputType,
			ExternalID:  st.ExternalID,
		},
	}, nil
}
