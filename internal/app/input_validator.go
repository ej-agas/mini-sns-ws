package app

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(validate *validator.Validate, input interface{}) (*ValidationErrors, error) {
	err := validate.Struct(input)
	errorRes := NewValidationErrors()

	if err == nil {
		return errorRes, nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())

		switch err.Tag() {
		case "required":
			errorRes.Add(field, fmt.Sprintf("The %s field is required.", strings.ReplaceAll(field, "_", " ")))
		case "gte":
			errorRes.Add(field, fmt.Sprintf("The %s field should have a minimum length of %s character(s).", strings.ReplaceAll(field, "_", " "), err.Param()))
		case "alpha":
			errorRes.Add(field, fmt.Sprintf("The %s field should only contain alpha characters.", strings.ReplaceAll(field, "_", " ")))
		case "alphanum":
			errorRes.Add(field, fmt.Sprintf("The %s field should only contain alphanumeric characters.", strings.ReplaceAll(field, "_", " ")))
		case "email":
			errorRes.Add(field, fmt.Sprintf("The %s field should be a valid email address.", field))
		case "alphanumunicode":
			errorRes.Add(field, fmt.Sprintf("The %s field should only contain alphanumeric unicode characters.", strings.ReplaceAll(field, "_", " ")))
		case "ascii":
			errorRes.Add(field, fmt.Sprintf("The %s field should only contain ascii characters.", strings.ReplaceAll(field, "_", " ")))
		}
	}

	return errorRes, err
}
