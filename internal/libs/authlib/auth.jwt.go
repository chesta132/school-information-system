package authlib

import (
	"school-information-system/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string     `json:"user_id"`
	Role     string     `json:"role"`
	RotateAt *time.Time `json:"rotate_at"`
	jwt.RegisteredClaims
}

func createClaim(id, role string, expiresAt time.Duration, rotateAt *time.Duration) *jwt.Token {
	var rotate *time.Time
	if rotateAt != nil {
		r := time.Now().Add(*rotateAt)
		rotate = &r
	}
	return jwt.NewWithClaims(config.SIGN_METHOD, Claims{
		UserID:   id,
		Role:     role,
		RotateAt: rotate,
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

func CreateRefreshToken(id, role string) string {
	token := createClaim(id, role, config.REFRESH_TOKEN_EXPIRY, &config.ROTATE_REFRESH_TOKEN_AFTER)
	str, _ := token.SignedString([]byte(config.REFRESH_SECRET))
	return str
}

func CreateAccessToken(id, role string) string {
	token := createClaim(id, role, config.ACCESS_TOKEN_EXPIRY, nil)
	str, _ := token.SignedString([]byte(config.ACCESS_SECRET))
	return str
}

func ParseRefreshToken(str string) (claims Claims, ok bool) {
	token, err := jwt.ParseWithClaims(str, &claims, createKeyFunc(config.REFRESH_SECRET))
	if err != nil {
		ok = false
		return
	}
	ok = token.Valid
	return
}

func ParseAccessToken(str string) (claims Claims, ok bool) {
	token, err := jwt.ParseWithClaims(str, &claims, createKeyFunc(config.ACCESS_SECRET))
	if err != nil {
		ok = false
		return
	}
	ok = token.Valid
	return
}
