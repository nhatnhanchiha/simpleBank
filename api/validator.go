package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/nhatnhanchiha/simpleBank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		util.IsSupportedCurrency(currency)
	}

	return false
}
