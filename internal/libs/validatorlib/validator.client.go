package validatorlib

import (
	"reflect"
	"school-information-system/internal/models"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Client = validator.New(validator.WithRequiredStructEnabled())

func registFirstTag(tags map[string]string) validator.TagNameFunc {
	return func(fld reflect.StructField) string {
		for tag, splitter := range tags {
			name, _, _ := strings.Cut(fld.Tag.Get(tag), splitter)
			if name != "-" && name != "" {
				return name
			}
		}
		return ""
	}
}

func registEnumValidation[E ~string](enum []E) validator.Func {
	return func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.String {
			return false
		}

		value := fl.Field().String()
		if !slices.Contains(enum, E(value)) {
			return false
		}

		return true
	}
}

func init() {
	// register tags
	Client.RegisterTagNameFunc(registFirstTag(map[string]string{
		"json": ",",
		"uri":  ",",
		"form": ",",
	}))

	// register enum user role tag
	Client.RegisterValidation("user_role", registEnumValidation(models.UserRoles))

	// register enum user gender tag
	Client.RegisterValidation("user_gender", registEnumValidation(models.UserGenders))

	// register enum perm action tag
	Client.RegisterValidation("permission_action", registEnumValidation(models.PermissionActions))

	// register enum perm resource tag
	Client.RegisterValidation("permission_resource", registEnumValidation(models.PermissionResources))
}
