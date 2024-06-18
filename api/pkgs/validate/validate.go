package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// map[field][]tags
type ErrFields map[string][]string

var v *validator.Validate

func getValidator() *validator.Validate {
	if v == nil {
		v = validator.New()
	}

	return v
}

func ValidateStruct(v any) error {
	return validatorErrHandling(getValidator().Struct(v))
}

func validatorErrHandling(validatorErr error) error {
	if validatorErr == nil {
		return nil
	}

	var errField = make(ErrFields)
	for _, err := range validatorErr.(validator.ValidationErrors) {
		tags := errField[err.Field()]
		tags = append(tags, err.Tag())
		errField[err.Field()] = tags
	}

	var errList []string
	for field, tags := range errField {
		errText := fmt.Sprintf("field '%s' failed validate for tag [%s]", field, strings.Join(tags, ", "))
		errList = append(errList, errText)
	}

	return errors.New(fmt.Sprintf("failed on validate: [%s]", strings.Join(errList, " - ")))

}
