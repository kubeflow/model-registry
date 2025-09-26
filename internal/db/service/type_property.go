package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"golang.org/x/exp/constraints"
	"gorm.io/gorm"
)

type typePropertyRepositoryImpl struct {
	db *gorm.DB
}

func NewTypePropertyRepository(db *gorm.DB) models.TypePropertyRepository {
	return &typePropertyRepositoryImpl{db: db}
}

func (r *typePropertyRepositoryImpl) Save(tp models.TypeProperty) (models.TypeProperty, error) {
	var stp schema.TypeProperty
	err := r.db.Where("type_id=? AND name=?", tp.GetTypeID(), tp.GetName()).First(&stp).Error
	if err == nil {
		oldType := intPointerString(stp.DataType)
		newType := intPointerString(tp.GetDataType())
		if oldType != newType {
			return tp, fmt.Errorf("invalid property type: data type is %s, cannot change to %s", oldType, newType)
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		stp.TypeID = tp.GetTypeID()
		stp.Name = tp.GetName()
		stp.DataType = tp.GetDataType()

		if err := r.db.Create(&stp).Error; err != nil {
			return tp, err
		}
	}

	return &models.TypePropertyImpl{
		TypeID:   stp.TypeID,
		Name:     stp.Name,
		DataType: stp.DataType,
	}, nil
}

func intPointerString[T constraints.Integer](v *T) string {
	if v == nil {
		return "<nil>"
	}
	return strconv.Itoa(int(*v))
}
