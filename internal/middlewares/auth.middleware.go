package middlewares

import (
	"school-information-system/config"
	"school-information-system/internal/libs/authlib"
	"school-information-system/internal/libs/replylib"
	"time"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	// add blacklist system later
}

func (*Auth) Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		rp := replylib.Client.New(adapter.AdaptGin(c))

		// read & validate access token
		accessCookie, err := c.Request.Cookie(config.ACCESS_TOKEN_KEY)
		if err == nil {
			access, err := authlib.ParseAccessToken(accessCookie.Value)
			if err == nil {
				c.Set("userID", access.UserID)
				c.Set("role", access.Role)
				c.Next()
				return
			}
		}

		// read & validate refresh token
		refreshCookie, err := c.Request.Cookie(config.REFRESH_TOKEN_KEY)
		if err != nil {
			rp.Error(replylib.CodeUnauthorized, "no refresh token provided").FailJSON()
			return
		}

		refresh, err := authlib.ParseRefreshToken(refreshCookie.Value)
		if err != nil {
			rp.Error(replylib.CodeUnauthorized, err.Error()).FailJSON()
			return
		}

		// update access token
		rememberMe := authlib.IsCookieRememberMe(*refreshCookie)
		ac := authlib.CreateAccessCookie(refresh.UserID, refresh.Role, rememberMe)
		rp.SetCookies(ac)

		// rotate refresh token
		if refresh.RotateAt.Before(time.Now()) {
			rc := authlib.CreateRefreshCookie(refresh.UserID, refresh.Role, rememberMe)
			rp.SetCookies(rc)
		}

		c.Set("userID", refresh.UserID)
		c.Set("role", refresh.Role)
		c.Next()
	}
}
