package config

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	APP_NAME = "school-information-system"

	// query

	MAX_CREATE_BATCH = 100 // max data per batch when inserting multiple records

	// auth token (jwt)

	REFRESH_TOKEN_KEY = "refresh_token"
	ACCESS_TOKEN_KEY  = "access_token"
)

var (
	// auth token (jwt)

	SIGN_METHOD                              = jwt.SigningMethodHS256
	REFRESH_TOKEN_EXPIRY       time.Duration = (time.Hour * 24 * 7 * 2) + (time.Hour * 24 * 3) // 2 weeks 3 days
	ACCESS_TOKEN_EXPIRY        time.Duration = time.Minute * 5                                 // 5 minutes
	ROTATE_REFRESH_TOKEN_AFTER time.Duration = time.Hour * 24 * 7 * 2                          // 2 weeks
)
