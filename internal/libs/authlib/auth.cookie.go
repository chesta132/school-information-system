package authlib

import (
	"net/http"
	"school-information-system/config"
	"time"
)

// set expires<0 to invalidate cookie
func ToCookie(name string, str string, expires time.Duration) (cookie http.Cookie) {
	cookie = http.Cookie{
		Name:     name,
		Value:    str,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   config.IsEnvProd(),
		HttpOnly: config.IsEnvDev(),
	}
	if expires > 0 {
		cookie.Expires = time.Now().Add(expires)
	} else if expires < 0 {
		cookie.MaxAge = -1
		cookie.Expires = time.Unix(0, 0)
	}
	return
}

func Invalidate(name string) http.Cookie {
	return ToCookie(name, "", -1)
}

func CreateRefreshCookie(id, role string, rememberMe bool) http.Cookie {
	expires := time.Duration(0)
	str := CreateRefreshToken(id, role)
	if rememberMe {
		expires = config.REFRESH_TOKEN_EXPIRY
	}

	return ToCookie(config.REFRESH_TOKEN_KEY, str, expires)
}

func CreateAccessCookie(id, role string, rememberMe bool) http.Cookie {
	expires := time.Duration(0)
	str := CreateAccessToken(id, role)
	if rememberMe {
		expires = config.ACCESS_TOKEN_EXPIRY
	}

	return ToCookie(config.ACCESS_TOKEN_KEY, str, expires)
}
