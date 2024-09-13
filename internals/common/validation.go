package common

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error       bool        `json:"error"`
	FailedField string      `json:"failed_field"`
	Tag         string      `json:"tag"`
	Value       interface{} `json:"value"`
}

type XValidator struct {
	validator *validator.Validate
}

var validate = validator.New()

func (v XValidator) Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func New() *XValidator {
	return &XValidator{
		validator: validate,
	}
}

func Validate(data interface{}) error {
	myValidator := New()

	if errs := myValidator.Validate(data); len(errs) > 0 && errs[0].Error {
		return &FieldValidationError{
			Errs: errs,
		}
	}

	return nil
}

type FieldValidationError struct {
	Errs []ErrorResponse
}

func (e *FieldValidationError) Error() string {
	errMsgs := make([]string, 0)
	for _, err := range e.Errs {
		errMsgs = append(errMsgs, fmt.Sprintf(
			"[%s]: '%v' | Needs to implement '%s'",
			err.FailedField,
			err.Value,
			err.Tag,
		))
	}

	return strings.Join(errMsgs, " and ")
}
