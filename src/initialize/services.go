package initialize

import (
	"app/di"
	"app/interfaces"
	"app/services"
)

func InitServices() {
	di.Container.RepositoryService = interfaces.IRepositoryService(new(services.RepositoryService))
}
