package service

import (
	"gorm.io/gorm"
)

type artifactHandler struct {
	db *gorm.DB
}
