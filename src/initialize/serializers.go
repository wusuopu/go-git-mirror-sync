package initialize

import (
	"app/di"
	"app/serializers"
)

func InitSerializers() {
	di.Container.RepositorySerializer = &serializers.RepositorySerializer{}
}