package controller

import (
	"github.com/go-playground/validator/v10"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
)

var validate = validator.New()

type ValidationError struct {
	Error  string      `json:"error"`
	Field  string      `json:"field"`
	Reason string      `json:"reason"`
	Value  interface{} `json:"value"`
}

func processValidationError(err error) ValidationError {
	validationErr := ValidationError{}
	for _, err := range err.(validator.ValidationErrors) {
		validationErr.Error = constant.ErrMessageValidationFailed
		validationErr.Field = err.Field()
		validationErr.Reason = err.Tag()
		validationErr.Value = err.Value()

		// Just return one validation error at once
		break
	}
	return validationErr
}
