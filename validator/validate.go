package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func TransformValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		var message string
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", err.Field())
		case "min":
			message = fmt.Sprintf("%s must be more than %s characters long", err.Field(), err.Param())
		case "max":
			message = fmt.Sprintf("%s must be less than %s character long", err.Field(), err.Param())
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", err.Field())
		case "gte":
			message = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Value())
		case "gt":
			message = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Value())
		case "lte":
			message = fmt.Sprintf("%s must be less than %s", err.Field(), err.Value())
		default:
			message = err.Error()
		}
		errors[strings.ToLower(err.Field())] = message

	}
	return errors
}
