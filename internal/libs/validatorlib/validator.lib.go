package validatorlib

import (
	"errors"
	"fmt"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/slicelib"
	"slices"
	"strings"

	"github.com/chesta132/goreply/reply"
	"github.com/go-playground/validator/v10"
)

func TranslateError(err error) (error, validator.ValidationErrors) {
	valErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err, nil
	}
	messages := make([]string, 0, len(valErrs))
	requiredErrs := make([]string, 0, len(valErrs))
	fieldPairs := make(map[string]string, len(valErrs))

	for _, err := range valErrs {
		fieldPairs[err.ActualTag()] = err.Field()
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
			case "min":
				messages = append(messages, fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param()))
			case "max":
				messages = append(messages, fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param()))
			case "required_if":
				param := strings.Split(err.Param(), " ")
				fieldParam, ok := fieldPairs[param[0]]
				if !ok {
					fieldParam = param[0]
				}
				messages = append(messages, fmt.Sprintf("%s is required if value of %s is %s", err.Field(), fieldParam, param[1]))
			}
		}
	}

	return errors.New(strings.Join(messages, ", ")), valErrs
}

func ValidateStructToReply(s any, fieldPrefix ...string) *reply.ErrorPayload {
	prefix := ""
	if len(fieldPrefix) > 0 {
		prefix = fieldPrefix[0]
	}
	if err := Client.Struct(s); err != nil {
		err, valErrs := TranslateError(err)
		fields := slicelib.Unique(slicelib.Map(valErrs, func(i int, val validator.FieldError) string { return prefix + val.Field() }))

		return &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: "invalid payload",
			Details: err.Error(),
			Fields:  fields,
		}
	}
	return nil
}
