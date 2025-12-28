package handlers

import (
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	authService *services.Auth
}

func NewAuth(authService *services.Auth) *Auth {
	return &Auth{authService}
}

// @Summary      Creates new account and sign in
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param				 payload  body	payloads.RequestSignUp	true	"data of new account"
// @Success      201  		{object}  swaglib.Envelope{data=models.User}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /auth/sign-up [post]
func (h *Auth) SignUp(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestSignUp
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	user, cookies, errPayload := h.authService.ApplyContext(c).SignUp(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.SetCookies(cookies...).Success(user).CreatedJSON()
}

// @Summary      Sign in to registered account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param				 payload  body			payloads.RequestSignIn	true	"data of registered account"
// @Success      200  		{object}  swaglib.Envelope{data=models.User}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /auth/sign-in [post]
func (h *Auth) SignIn(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestSignIn
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	user, cookies, errPayload := h.authService.ApplyContext(c).SignIn(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.SetCookies(cookies...).Success(user).OkJSON()
}

// @Summary      Sign out from session
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  		{object}  swaglib.Envelope{data=nil}
// @Router       /auth/sign-out [post]
func (h *Auth) SignOut(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	cookies := h.authService.ApplyContext(c).SignOut()
	rp.SetCookies(cookies...).Success(nil).OkJSON()
}
