package validatorlib

import (
	"errors"
	"fmt"
	"reflect"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/slicelib"
	"slices"
	"strings"

	"github.com/chesta132/goreply/reply"
	"github.com/go-playground/validator/v10"
)

// extract prefix map recursively
func extractPrefixMap(s any) map[string]string {
	prefixMap := make(map[string]string)
	buildPrefixMap(reflect.ValueOf(s), "", prefixMap)
	return prefixMap
}

func buildPrefixMap(val reflect.Value, parentPrefix string, prefixMap map[string]string) {
	// handle pointer
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	// only process struct
	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if !field.IsExported() {
			continue
		}

		fieldVal := val.Field(i)

		// get field name from json tag
		jsonTag := field.Tag.Get("json")
		fieldName, _, _ := strings.Cut(jsonTag, ",")
		if fieldName == "" || fieldName == "-" {
			fieldName = field.Name
		}

		// check if field has prefix tag
		prefix := field.Tag.Get("prefix")

		// if has prefix, nested fields inside need this prefix
		if prefix != "" {
			// recurse into nested struct
			if fieldVal.Kind() == reflect.Ptr || fieldVal.Kind() == reflect.Struct {
				buildPrefixMap(fieldVal, prefix, prefixMap)
			}
		} else if parentPrefix != "" {
			// if no prefix but has parent prefix, apply to this field
			prefixMap[fieldName] = parentPrefix
		}

		// still recurse for nested struct without prefix tag
		if prefix == "" && (fieldVal.Kind() == reflect.Ptr || fieldVal.Kind() == reflect.Struct) {
			buildPrefixMap(fieldVal, parentPrefix, prefixMap)
		}
	}
}

func TranslateError(err error, prefixMap map[string]string) (error, validator.ValidationErrors) {
	valErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err, nil
	}
	messages := make([]string, 0, len(valErrs))
	requiredErrs := make([]string, 0, len(valErrs))

	for _, err := range valErrs {
		// get prefix from map if exists
		fieldName := err.Field()
		if prefix, ok := prefixMap[fieldName]; ok {
			fieldName = prefix + fieldName
		}

		// handle required error
		if err.Tag() == "required" {
			messages = append(messages, fmt.Sprintf("%s is required", fieldName))
			requiredErrs = append(requiredErrs, fieldName)
		}

		// handle other validation errors
		if !slices.Contains(requiredErrs, fieldName) {
			switch err.Tag() {
			case "email":
				messages = append(messages, fmt.Sprintf("%s is not valid email", fieldName))
			case "oneof":
				messages = append(messages, fmt.Sprintf("%s is not a valid enum of [%s]", fieldName, err.Param()))
			case "min", "max":
				suffix := ""
				switch err.Kind() {
				case reflect.Slice, reflect.Array:
					suffix = " items"
				case reflect.String:
					suffix = " characters"
				}

				operator := "at least"
				if err.Tag() == "max" {
					operator = "at most"
				}
				messages = append(messages, fmt.Sprintf("%s must be %s %s%s", fieldName, operator, err.Param(), suffix))
			case "required_if":
				targetField, targetValue, _ := strings.Cut(err.Param(), " ")
				messages = append(messages, fmt.Sprintf("%s is required if value of %s is %s", fieldName, targetField, targetValue))
			}
		}
	}

	return errors.New(strings.Join(messages, ", ")), valErrs
}

func ValidateStructToReply(s any) *reply.ErrorPayload {
	if err := Client.Struct(s); err != nil {
		// extract prefix map from struct tags
		prefixMap := extractPrefixMap(s)

		// translate validator errors with prefix
		err, valErrs := TranslateError(err, prefixMap)

		// collect field names with prefix applied
		fields := slicelib.Unique(slicelib.Map(valErrs, func(i int, val validator.FieldError) string {
			fieldName := val.Field()
			if prefix, ok := prefixMap[fieldName]; ok {
				return prefix + fieldName
			}
			return fieldName
		}))

		return &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: "invalid payload",
			Details: err.Error(),
			Fields:  fields,
		}
	}
	return nil
}
