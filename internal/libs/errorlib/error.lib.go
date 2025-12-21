package errorlib

import "errors"

var (
	ErrNotActivated = errors.New("account not activated yet, please wait for admin approval")
)
