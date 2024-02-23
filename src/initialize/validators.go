package initialize

import (
	"app/di"
	"app/validators"
)

func InitValidators() {
	di.Container.RepositoryValidator = validators.NewRepositoryValidator()
}