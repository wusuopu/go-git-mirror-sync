package validators

import (
	"app/di"
	"app/models"
	"app/utils/helper"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/valyala/fastjson"
)

type RepositoryValidator struct {
	validate *validator.Validate
	trans *ut.Translator
}

func NewRepositoryValidator() *RepositoryValidator {
	validate := validator.New()
	validate.RegisterValidation("uniqName", func(fl validator.FieldLevel) bool {
		// 检查该 Name 对应的 repository 是否存在
		entity := models.Repository{}
		results := di.Container.DB.First(&entity, "name = ?", fl.Field().String())
		return results.RowsAffected == 0
	})

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTranslation("uniqName", trans, func(ut ut.Translator) error {
		return ut.Add("uniqName", "{0} {1} Name 不能重复", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("uniqName", fe.Field(), fe.Tag())
		return t
	})

	v := new(RepositoryValidator)
	v.validate = validate
	v.trans = &trans
	return v
}

func (v *RepositoryValidator) ValidateCreate(value *fastjson.Value) (*models.Repository, validator.ValidationErrorsTranslations) {
	var err validator.ValidationErrorsTranslations = nil

	Name := helper.GetJSONString(value, "Name")
	err = validateValueByRule(v.validate, Name, "required,uniqName", "Name", v.trans)
	if err != nil {
		return nil, err
	}

	Alias := helper.GetJSONString(value, "Alias")
	// err = v.validate.Var(Alias, "required")
	err = validateValueByRule(v.validate, Alias, "required", "Alias", v.trans)
	if err != nil {
		return nil, err
	}

	Url := helper.GetJSONString(value, "Url")
	// err = v.validate.Var(Url, "required")
	err = validateValueByRule(v.validate, Url, "required", "Url", v.trans)
	if err != nil {
		return nil, err
	}

	AuthType := helper.GetJSONString(value, "AuthType")
	err = validateValueByRule(v.validate, AuthType, "required", "AuthType", v.trans)
	if err != nil {
		return nil, err
	}

	Username := helper.GetJSONString(value, "Username")
	Password := helper.GetJSONString(value, "Password")
	SSHKey := helper.GetJSONString(value, "SSHKey")
	entity := models.Repository{
		Name: Name,
		Alias: Alias,
		Url: Url,
		AuthType: AuthType,
		Username: &Username,
		Password: &Password,
		SSHKey: &SSHKey,
	}
	
	return &entity, nil
}

func (v *RepositoryValidator) ValidateUpdate(value *fastjson.Value) (*models.Repository, validator.ValidationErrorsTranslations) {
	return nil, nil
}