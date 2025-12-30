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
	"school-information-system/internal/repos"
	"slices"
	"strings"
	"time"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Auth struct {
	userRepo    *repos.User
	revokedRepo *repos.Revoked
}

func NewAuth(userRepo *repos.User, revokedRepo *repos.Revoked) *Auth {
	return &Auth{userRepo, revokedRepo}
}

func (mw *Auth) protected(c *gin.Context) (claims authlib.Claims, newAccessCookie, newRefreshCookie *http.Cookie, err error) {
	ctx := c.Request.Context()
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

	// check if token is revoked
	if revoked, revErr := mw.revokedRepo.GetFirst(ctx, "token = ?", refreshCookie.Value); !errors.Is(revErr, gorm.ErrRecordNotFound) {
		if revErr != nil {
			err = revErr
			return
		}
		err = errors.New(authlib.MessageOfRevoke(revoked))
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

// ensureAuthenticated to ensure protected middleware run, return false if not authenticated and aborted
func (mw *Auth) ensureAuthenticated(c *gin.Context, rp *reply.Reply, allowUnsetted bool) bool {
	// if prev middleware is protecting return true
	if _, exists := c.Get("userID"); exists {
		return true
	}

	// validate token
	claims, newAccessCookie, newRefreshCookie, err := mw.protected(c)
	if err != nil {
		if errors.Is(err, errorlib.ErrNotActivated) && allowUnsetted {
			mw.applyInternalProtectedReturn(c, rp, claims, newAccessCookie, newRefreshCookie)
			return true
		}

		if errors.Is(err, errorlib.ErrNotActivated) {
			rp.Error(replylib.CodeForbidden, err.Error()).FailJSON()
		} else {
			rp.Error(replylib.CodeUnauthorized, err.Error()).FailJSON()
		}
		c.Abort()
		return false
	}

	mw.applyInternalProtectedReturn(c, rp, claims, newAccessCookie, newRefreshCookie)
	return true
}

// Protected middleware basic auth
func (mw *Auth) Protected(allowUnsettedRole ...bool) gin.HandlerFunc {
	allowUnsetted := false
	if len(allowUnsettedRole) > 0 {
		allowUnsetted = allowUnsettedRole[0]
	}
	return func(c *gin.Context) {
		rp := replylib.Client.Use(adapter.AdaptGin(c))
		if mw.ensureAuthenticated(c, rp, allowUnsetted) {
			c.Next()
		}
	}
}

// RoleProtected protects with role validation
func (mw *Auth) RoleProtected(roles ...models.UserRole) gin.HandlerFunc {
	strRoles := make([]string, len(roles))
	for i, r := range roles {
		strRoles[i] = string(r)
	}

	return func(c *gin.Context) {
		rp := replylib.Client.Use(adapter.AdaptGin(c))

		// make sure user is authenticated
		allowUnsetted := slices.Contains(strRoles, string(models.RoleUnsetted))
		if !mw.ensureAuthenticated(c, rp, allowUnsetted) {
			return
		}

		roleInterface, _ := c.Get("role")
		role, _ := roleInterface.(string)

		// protect unsetted role
		if role == string(models.RoleUnsetted) && !allowUnsetted {
			rp.Error(replylib.CodeForbidden, errorlib.ErrNotActivated.Error()).FailJSON()
			c.Abort()
			return
		}

		// protect roles
		if !slices.Contains(strRoles, role) {
			rp.Error(replylib.CodeForbidden, fmt.Sprintf("invalid role, only %s can access this resource", strings.Join(strRoles, ", "))).FailJSON()
			c.Abort()
			return
		}

		c.Next()
	}
}

type PermissionProtectedOpt struct {
	skipOnRole []models.UserRole
}

type PermissionProtectedOptFunc func(*PermissionProtectedOpt)

func WithSkipRole(roles ...models.UserRole) PermissionProtectedOptFunc {
	return func(ppo *PermissionProtectedOpt) {
		ppo.skipOnRole = roles
	}
}

// PermissionProtected protects with permission validation
func (mw *Auth) PermissionProtected(resource models.PermissionResource, actions []models.PermissionAction, opts ...PermissionProtectedOptFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		rp := replylib.Client.Use(adapter.AdaptGin(c))
		ctx := c.Request.Context()

		option := new(PermissionProtectedOpt)
		for _, opt := range opts {
			opt(option)
		}

		// make sure user is authenticated
		if !mw.ensureAuthenticated(c, rp, false) {
			return
		}

		userIDInterface, _ := c.Get("userID")
		userID, _ := userIDInterface.(string)
		roleInterface, _ := c.Get("role")
		role, _ := roleInterface.(string)

		if len(option.skipOnRole) > 0 && slices.Contains(option.skipOnRole, models.UserRole(role)) {
			c.Next()
			return
		}

		// validate role is admin
		if role != string(models.RoleAdmin) {
			rp.Error(replylib.CodeForbidden, "invalid role, only admin can access this resource").FailJSON()
			c.Abort()
			return
		}

		// get and set user with permissions
		user, err := mw.userRepo.GetFirstWithPreload(ctx, []string{"AdminProfile.Permissions"}, "id = ? AND role = ?", userID, models.RoleAdmin)
		if err != nil {
			errPayload := errorlib.MakeNotFound(err, "your user profile not found", nil)
			rp.Error(errPayload.Code, errPayload.Message).FailJSON()
			c.Abort()
			return
		}
		if user.AdminProfile == nil || len(user.AdminProfile.Permissions) == 0 {
			rp.Error(replylib.CodeConflict, "your admin profile or permission not registered").FailJSON()
			c.Abort()
			return
		}
		c.Set("user", user)

		// user's action tracks
		existingAction := make(map[models.PermissionAction]struct{}, len(actions))

		for _, perm := range user.AdminProfile.Permissions {
			if perm.Resource == resource {
				for _, act := range perm.Actions {
					existingAction[act] = struct{}{}
				}
			}
		}

		// validate is permitted
		for _, action := range actions {
			if _, ok := existingAction[action]; !ok {
				rp.Error(replylib.CodeForbidden, fmt.Sprintf("missing permission: %s.%s", resource, action)).FailJSON()
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
