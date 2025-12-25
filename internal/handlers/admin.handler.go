package handlers

import (
	"fmt"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
)

type Admin struct {
	adminService      *services.Admin
	roleSetterService *services.RoleSetter
}

func NewAdmin(adminService *services.Admin, roleSetterService *services.RoleSetter) *Admin {
	return &Admin{adminService, roleSetterService}
}

func (h *Admin) InitiateAdmin(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestInitiateAdmin
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}
	user, errPayload := h.adminService.ApplyContext(c).InitiateAdmin(payload)
	if errPayload != nil {
		rp.Error(errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}).FailJSON()
		return
	}

	rp.Success(user).Info(fmt.Sprintf("Admin created: %s (%s)", user.FullName, payload.StaffRole)).OkJSON()
}

func (h *Admin) SetRole(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestSetRole
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	user, errPayload := h.roleSetterService.ApplyContext(c).SetRole(payload)
	if errPayload != nil {
		rp.Error(errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}).FailJSON()
		return
	}

	rp.Success(user).Info(fmt.Sprintf("%s's role setted", user.FullName)).OkJSON()
}
