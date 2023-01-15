package app

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(validate *validator.Validate, input interface{}) (*ErrorResponse, error) {
	err := validate.Struct(input)
	errorRes := NewErrorResponse()

	if err == nil {
		return errorRes, nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())

		switch err.Tag() {
		case "required":
			errorRes.Add(field, fmt.Sprintf("The %s field is required.", strings.ReplaceAll(field, "_", " ")))
		case "gte":
			errorRes.Add(field, fmt.Sprintf("The %s field should a minimum length of %s character(s).", strings.ReplaceAll(field, "_", " "), err.Param()))
		}
	}

	fmt.Println(errorRes)

	return errorRes, err
}
