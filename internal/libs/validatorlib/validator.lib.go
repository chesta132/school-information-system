package validatorlib

import (
	"fmt"
	"reflect"
	"school-information-system/internal/libs/replylib"
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

func TranslateError(err error, prefixMap map[string]string) reply.FieldsError {
	valErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	fields := make(reply.FieldsError, len(valErrs))
	insertedErr := make(map[string]struct{}, len(valErrs))

	for _, err := range valErrs {
		// get prefix from map if exists
		fieldName := err.Field()
		if prefix, ok := prefixMap[fieldName]; ok {
			fieldName = prefix + fieldName
		}

		// handle required error
		if err.Tag() == "required" {
			fields[fieldName] = fmt.Sprintf("%s is required", fieldName)
			insertedErr[fieldName] = struct{}{}
		}

		// handle other validation errors
		if _, ok := insertedErr[fieldName]; !ok {
			switch err.Tag() {
			case "email":
				fields[fieldName] = fmt.Sprintf("%s is not valid email", fieldName)
				insertedErr[fieldName] = struct{}{}
			case "oneof":
				fields[fieldName] = fmt.Sprintf("%s is not a valid enum of [%s]", fieldName, err.Param())
				insertedErr[fieldName] = struct{}{}
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
				fields[fieldName] = fmt.Sprintf("%s must be %s %s%s", fieldName, operator, err.Param(), suffix)
				insertedErr[fieldName] = struct{}{}
			case "required_if":
				targetField, targetValue, _ := strings.Cut(err.Param(), " ")
				fields[fieldName] = fmt.Sprintf("%s is required if value of %s is %s", fieldName, targetField, targetValue)
				insertedErr[fieldName] = struct{}{}
			case "required_without":
				fields[fieldName] = fmt.Sprintf("%s required if %s empty", fieldName, err.Param())
				insertedErr[fieldName] = struct{}{}
			}
		}
	}

	return fields
}

func ValidateStructToReply(s any) *reply.ErrorPayload {
	if err := Client.Struct(s); err != nil {
		// extract prefix map from struct tags
		prefixMap := extractPrefixMap(s)

		// translate validator errors with prefix
		fields := TranslateError(err, prefixMap)

		return &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: "invalid payload",
			Fields:  fields,
		}
	}
	return nil
}
