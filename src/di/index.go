package di

import (
	"app/interfaces"

	"gorm.io/gorm"
)

type container struct {
	RepositoryService interfaces.IRepositoryService
	DB *gorm.DB
}

var Container = new(container)
