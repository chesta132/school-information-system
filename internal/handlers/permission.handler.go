package handlers

import (
	"fmt"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/slicelib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/services"
	"strings"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
)

type Permission struct {
	permService *services.Permission
}

func NewPermission(permService *services.Permission) *Permission {
	return &Permission{permService}
}

func (h *Permission) GrantPermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestGrantPermission
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	user, permission, errPayload := h.permService.ApplyContext(c).GrantPermission(payload)
	if errPayload != nil {
		rp.Error(errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}).FailJSON()
		return
	}

	// strings.Join can't process []models.PermissionAction
	permActs := slicelib.Map(permission.Actions, func(i int, act models.PermissionAction) string { return string(act) })
	rp.Success(user).Info(fmt.Sprintf("%s's permitted to %s %s", user.FullName, strings.Join(permActs, ", "), permission.Resource)).OkJSON()
}
