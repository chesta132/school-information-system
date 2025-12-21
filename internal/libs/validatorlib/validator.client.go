package validatorlib

import "github.com/go-playground/validator/v10"

var Client = validator.New(validator.WithRequiredStructEnabled())
