package interfaces

import (
	"app/models"

	"github.com/go-playground/validator/v10"
	"github.com/valyala/fastjson"
)


type IRepositoryValidator interface {
	ValidateCreate(value *fastjson.Value) (*models.Repository, validator.ValidationErrorsTranslations)
	ValidateUpdate(value *fastjson.Value) (*models.Repository, validator.ValidationErrorsTranslations)
}