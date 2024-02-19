package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Property string `json:"Property"`
	Tag      string `json:"Tag"`
	Value    string `json:"Value"`
	Message  string `json:"Message"`
}

func GetValidationError(err error) *[]ValidationError {
	var validationErrors []ValidationError
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, e := range err.(validator.ValidationErrors) {
			var el ValidationError
			el.Message = e.StructField()
			el.Property = e.Field()
			el.Tag = e.Tag()
			el.Value = e.Param()
			validationErrors = append(validationErrors, el)
		}
		return &validationErrors
	}
	return nil
}
