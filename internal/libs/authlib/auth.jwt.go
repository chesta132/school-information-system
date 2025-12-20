package authlib

import (
	"errors"
	"school-information-system/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     string     `json:"user_id"`
	Role       string     `json:"role"`
	RotateAt   *time.Time `json:"rotate_at"`
	RememberMe bool       `json:"remember_me"`
	jwt.RegisteredClaims
}

func createClaim(id, role string, rememberMe bool, expiresAt time.Duration, rotateAt *time.Duration) *jwt.Token {
	var rotate *time.Time

	if rotateAt != nil {
		r := time.Now().Add(*rotateAt)
		rotate = &r

		if !rememberMe {
			// max 1 day expiry for not remember me refresh token
			maxSessionExpiry := time.Hour * 24
			if expiresAt > maxSessionExpiry {
				expiresAt = maxSessionExpiry
			}
		}
	}

	return jwt.NewWithClaims(config.SIGN_METHOD, Claims{
		UserID:     id,
		Role:       role,
		RotateAt:   rotate,
		RememberMe: rememberMe,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.APP_NAME,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		},
	})
}

func createKeyFunc(secret string) jwt.Keyfunc {
	return func(t *jwt.Token) (any, error) {
		if t.Method != config.SIGN_METHOD {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(secret), nil
	}
}

func CreateRefreshToken(id, role string, rememberMe bool) string {
	token := createClaim(id, role, rememberMe, config.REFRESH_TOKEN_EXPIRY, &config.ROTATE_REFRESH_TOKEN_AFTER)
	str, _ := token.SignedString([]byte(config.REFRESH_SECRET))
	return str
}

func CreateAccessToken(id, role string, rememberMe bool) string {
	token := createClaim(id, role, rememberMe, config.ACCESS_TOKEN_EXPIRY, nil)
	str, _ := token.SignedString([]byte(config.ACCESS_SECRET))
	return str
}

func ParseRefreshToken(str string) (claims Claims, err error) {
	token, err := jwt.ParseWithClaims(str, &claims, createKeyFunc(config.REFRESH_SECRET))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = errors.New("refresh token is expired")
		}
		return
	}
	if !token.Valid {
		err = errors.New("invalid refresh token")
	}
	return
}

func ParseAccessToken(str string) (claims Claims, err error) {
	token, err := jwt.ParseWithClaims(str, &claims, createKeyFunc(config.ACCESS_SECRET))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = errors.New("access token is expired")
		}
		return
	}
	if !token.Valid {
		err = errors.New("invalid access token")
	}
	return
}
