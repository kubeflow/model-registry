package repository

import (
	"errors"
	"fmt"

	"github.com/kubeflow/model-registry/internal/platform/db/entity"
	"github.com/kubeflow/model-registry/internal/platform/db/schema"
	"gorm.io/gorm"
)

type TypeRepositoryImpl struct {
	db *gorm.DB
}

func NewTypeRepository(db *gorm.DB) entity.TypeRepository {
	return &TypeRepositoryImpl{db: db}
}

func (r *TypeRepositoryImpl) GetAll() ([]entity.Type, error) {
	var types []schema.Type

	if err := r.db.Find(&types).Error; err != nil {
		return nil, err
	}

	typesModels := make([]entity.Type, len(types))

	for i, t := range types {
		typesModels[i] = &entity.TypeImpl{
			ID: &t.ID,
			Attributes: &entity.TypeAttributes{
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

func (r *TypeRepositoryImpl) Save(t entity.Type) (entity.Type, error) {
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
		if st.TypeKind != *attr.TypeKind {
			return t, fmt.Errorf("invalid type: kind is %d, cannot change to kind %d", st.TypeKind, *attr.TypeKind)
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
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

	return &entity.TypeImpl{
		ID: &st.ID,
		Attributes: &entity.TypeAttributes{
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
