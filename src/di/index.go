package di

import (
	"app/interfaces"
	"app/serializers"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type container struct {
	RepositoryService interfaces.IRepositoryService

	DB *gorm.DB
	Logger *zap.Logger

	Scheduler *gocron.Scheduler

	RepositorySerializer *serializers.RepositorySerializer

	RepositoryValidator interfaces.IRepositoryValidator
}

var Container = new(container)
