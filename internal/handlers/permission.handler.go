package handlers

import (
	"fmt"
	"school-information-system/config"
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

func (h *Permission) CreatePermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestCreatePermission
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	perm, errPayload := h.permService.ApplyContext(c).CreatePermission(payload)
	if errPayload != nil {
		rp.Error(errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}).FailJSON()
		return
	}

	// strings.Join can't process []models.PermissionAction
	permActs := slicelib.Map(perm.Actions, func(i int, act models.PermissionAction) string { return string(act) })
	rp.Success(perm).Info(fmt.Sprintf("new permission to %s %s created", strings.Join(permActs, ", "), perm.Resource)).CreatedJSON()
}

func (h *Permission) GetPermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	id := c.Param("id")

	perm, errPayload := h.permService.ApplyContext(c).GetPermission(id)
	if errPayload != nil {
		rp.Error(errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}).FailJSON()
		return
	}

	rp.Success(perm).OkJSON()
}

func (h *Permission) GetPermissions(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestGetPermissions
	c.ShouldBindQuery(&payload)

	perm, errPayload := h.permService.ApplyContext(c).GetPermissions(payload)
	if errPayload != nil {
		rp.Error(errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}).FailJSON()
		return
	}

	rp.Success(perm).PaginateCursor(config.LIMIT_PAGINATED_DATA, payload.Offset).OkJSON()
}

func (h *Permission) DeletePermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	id := c.Param("id")

	errPayload := h.permService.ApplyContext(c).DeletePermission(id)
	if errPayload != nil {
		rp.Error(errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}).FailJSON()
		return
	}

	rp.Success(map[string]string{"id": id}).OkJSON()
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

func (h *Permission) RevokePermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestRevokePermission
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	user, permission, errPayload := h.permService.ApplyContext(c).RevokePermission(payload)
	if errPayload != nil {
		rp.Error(errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}).FailJSON()
		return
	}

	// strings.Join can't process []models.PermissionAction
	permActs := slicelib.Map(permission.Actions, func(i int, act models.PermissionAction) string { return string(act) })
	rp.Success(user).Info(fmt.Sprintf("%s's no longer permitted to %s %s", user.FullName, strings.Join(permActs, ", "), permission.Resource)).OkJSON()
}
