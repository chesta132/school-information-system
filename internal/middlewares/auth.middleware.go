package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"school-information-system/config"
	"school-information-system/internal/libs/authlib"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models"
	"slices"
	"strings"
	"time"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	// add blacklist system later
}

func (*Auth) protected(c *gin.Context) (claims authlib.Claims, newAccessCookie, newRefreshCookie *http.Cookie, err error) {
	accessCookie, err := c.Request.Cookie(config.ACCESS_TOKEN_KEY)
	if err == nil {
		claims, err = authlib.ParseAccessToken(accessCookie.Value)
		if err == nil {
			if claims.Role == string(models.RoleUnsetted) {
				err = errorlib.ErrNotActivated
			}
			return
		}
	}

	// read & validate refresh token
	refreshCookie, err := c.Request.Cookie(config.REFRESH_TOKEN_KEY)
	if err != nil {
		err = errors.New("no refresh token provided")
		return
	}

	claims, err = authlib.ParseRefreshToken(refreshCookie.Value)
	if err != nil {
		return
	}
	if claims.Role == string(models.RoleUnsetted) {
		err = errorlib.ErrNotActivated
		return
	}

	// update access token
	na := authlib.CreateAccessCookie(claims.UserID, claims.Role, claims.RememberMe)
	newAccessCookie = &na

	// rotate refresh token
	if claims.RotateAt.Before(time.Now()) {
		nr := authlib.CreateRefreshCookie(claims.UserID, claims.Role, claims.RememberMe)
		newRefreshCookie = &nr
	}
	return
}

func (*Auth) applyInternalProtectedReturn(c *gin.Context, rp *reply.Reply, claims authlib.Claims, newAccessCookie, newRefreshCookie *http.Cookie) {
	if newAccessCookie != nil {
		rp.SetCookies(*newAccessCookie)
	}
	if newRefreshCookie != nil {
		rp.SetCookies(*newRefreshCookie)
	}
	c.Set("userID", claims.UserID)
	c.Set("role", claims.Role)
}

func (mw *Auth) Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		rp := replylib.Client.New(adapter.AdaptGin(c))
		claims, newAccessCookie, newRefreshCookie, err := mw.protected(c)
		if err != nil {
			rp.Error(replylib.CodeUnauthorized, err.Error()).FailJSON()
			c.Abort()
			return
		}
		mw.applyInternalProtectedReturn(c, rp, claims, newAccessCookie, newRefreshCookie)
		c.Next()
	}
}

func (mw *Auth) RoleProtected(roles ...models.UserRole) gin.HandlerFunc {
	strRoles := make([]string, len(roles))
	for i, r := range roles {
		strRoles[i] = string(r)
	}

	return func(c *gin.Context) {
		rp := replylib.Client.New(adapter.AdaptGin(c))

		// get role
		var role string
		if r, exists := c.Get("role"); exists {
			role, _ = r.(string)
		} else {
			claims, newAccessCookie, newRefreshCookie, err := mw.protected(c)
			if err != nil {
				// if unsetted role allowed and claims role is unsetted then next
				if errors.Is(err, errorlib.ErrNotActivated) && slices.Contains(strRoles, string(models.RoleUnsetted)) {
					mw.applyInternalProtectedReturn(c, rp, claims, newAccessCookie, newRefreshCookie)
					c.Next()
					return
				}
				rp.Error(replylib.CodeUnauthorized, err.Error()).FailJSON()
				c.Abort()
				return
			}
			mw.applyInternalProtectedReturn(c, rp, claims, newAccessCookie, newRefreshCookie)
			role = claims.Role
		}

		// protect unsetted role and send better error message
		if role == string(models.RoleUnsetted) && !slices.Contains(strRoles, string(models.RoleUnsetted)) {
			rp.Error(replylib.CodeForbidden, errorlib.ErrNotActivated.Error()).FailJSON()
			c.Abort()
			return
		}

		if slices.Contains(strRoles, role) {
			c.Next()
			return
		}

		rp.Error(replylib.CodeForbidden, fmt.Sprintf("invalid role, only %s can access this resource", strings.Join(strRoles, ", "))).FailJSON()
		c.Abort()
	}
}
