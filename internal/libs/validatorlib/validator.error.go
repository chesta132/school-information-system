package validatorlib

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

func TranslateError(err error) (error, validator.ValidationErrors) {
	valErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err, nil
	}
	messages := make([]string, 0, len(valErrs))
	requiredErrs := make([]string, 0, len(valErrs))

	for _, err := range valErrs {
		if err.Tag() == "required" {
			messages = append(messages, fmt.Sprintf("%s is required", err.Field()))
			requiredErrs = append(requiredErrs, err.Field())
		}

		if !slices.Contains(requiredErrs, err.Field()) {
			switch err.Tag() {
			case "email":
				messages = append(messages, fmt.Sprintf("%s is not valid email", err.Field()))
			case "oneof":
				messages = append(messages, fmt.Sprintf("%s is not a valid enum of [%s]", err.Field(), err.Param()))
			}
		}
	}

	return errors.New(strings.Join(messages, ", ")), valErrs
}
