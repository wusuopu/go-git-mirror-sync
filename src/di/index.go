package di

import (
	"app/interfaces"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type container struct {
	RepositoryService interfaces.IRepositoryService
	DB *gorm.DB
	Logger *zap.Logger
	Scheduler *gocron.Scheduler
}

var Container = new(container)
