package errorlib

import "errors"

var (
	ErrInvalidEnv        = errors.New("[ENVIRONMENT] invalid environment, must be 'development' or 'production'")
	ErrNotActivated      = errors.New("account not activated yet, please wait for admin approval")
	ErrInvalidRole       = errors.New("you can not access this resource")
	ErrTargetInvalidRole = errors.New("targeted user does not have a permitted role to access the resource")
)
