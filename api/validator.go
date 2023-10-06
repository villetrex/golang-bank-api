package api

import "github.com/go-playground/validator/v10"

var validCurency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is supported
		return util.IsValidCurrency(currency)
	}
	return false
}
