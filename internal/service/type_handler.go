package service

import (
	"context"
	"fmt"

	"github.com/opendatahub-io/model-registry/internal/model/db"
	"gorm.io/gorm"
)

type typeHandler struct {
	db *gorm.DB
}

func (h *typeHandler) InsertType(i db.Type) (r *db.Type, err error) {
	ctx, _ := Begin(context.Background(), h.db)
	defer handleTransaction(ctx, &err)

	result := h.db.Create(&i)
	if result.Error != nil {
		return nil, result.Error
	}
	return &i, nil
}

func (h *typeHandler) UpsertType(i db.Type) (r *db.Type, err error) {
	ctx, _ := Begin(context.Background(), h.db)
	defer handleTransaction(ctx, &err)

	if err := h.db.Where("name = ?", i.Name).Assign(i).FirstOrCreate(&i).Error; err != nil {
		err = fmt.Errorf("error creating type %s: %v", i.Name, err)
		return nil, err
	}
	return &i, nil
}

func (h *typeHandler) ReadType(i db.Type) (*db.Type, error) {
	var results []*db.Type
	rx := h.db.Find(&results)
	if rx.Error != nil {
		return nil, rx.Error
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("found more than one Type(s): %v", len(results))
	}
	return results[0], nil
}

func (h *typeHandler) ReadAllType(query interface{}, args ...interface{}) ([]*db.Type, error) {
	var results []*db.Type
	rx := h.db.Where(query, args).Find(&results)
	if rx.Error != nil {
		return nil, rx.Error
	}
	return results, nil
}

func (h *typeHandler) UpdateType(i db.Type) (result *db.Type, err error) {
	panic("unimplemented")
}

func (h *typeHandler) DeleteType(i db.Type) (r *db.Type, err error) {
	panic("unimplemented")
}
