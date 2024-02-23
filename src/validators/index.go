package validators

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)


func validateValueByRule(v *validator.Validate, value interface{}, rule string, fieldName string, trans *ut.Translator) validator.ValidationErrorsTranslations {
	err := v.Var(value, rule)
	if err == nil {
		return nil
	}
	ret := err.(validator.ValidationErrors).Translate(*trans)
	ret[fieldName] = ret[""]
	delete(ret, "")
	return ret
}