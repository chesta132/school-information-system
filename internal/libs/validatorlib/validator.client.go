package validatorlib

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Client = validator.New(validator.WithRequiredStructEnabled())

func RegisterTagName() {
	Client.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
