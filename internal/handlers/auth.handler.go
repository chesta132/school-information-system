package handlers

import (
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models/payload"
	"school-information-system/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	authService *services.Auth
}

func NewAuth(authService *services.Auth) *Auth {
	return &Auth{authService}
}

func (h *Auth) SignUp(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payload.RequestSignUp
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	user, cookies, errPayload := h.authService.ApplyContext(c).SignUp(payload)
	if errPayload != nil {
		rp.Error(errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}).FailJSON()
		return
	}

	rp.SetCookies(cookies...).Success(user).OkJSON()
}

func (h *Auth) SignIn(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payload.RequestSignIn
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	user, cookies, errPayload := h.authService.ApplyContext(c).SignIn(payload)
	if errPayload != nil {
		rp.Error(errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}).FailJSON()
		return
	}

	rp.SetCookies(cookies...).Success(user).OkJSON()
}

func (h *Auth) SignOut(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	cookies := h.authService.ApplyContext(c).SignOut()
	rp.SetCookies(cookies...).Success(nil).OkJSON()
}
