package di

import (
	"app/interfaces"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type container struct {
	RepositoryService interfaces.IRepositoryService
	DB *gorm.DB
	Logger *zap.Logger
}

var Container = new(container)
