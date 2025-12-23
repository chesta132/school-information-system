package errorlib

import "errors"

var (
	ErrInvalidEnv = errors.New("[ENVIRONMENT] invalid environment, must be 'development' or 'production'")
	ErrNotActivated = errors.New("account not activated yet, please wait for admin approval")
)
