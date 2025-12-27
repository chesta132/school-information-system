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
	"github.com/gin-gonic/gin"
)

type Permission struct {
	permService *services.Permission
}

func NewPermission(permService *services.Permission) *Permission {
	return &Permission{permService}
}

// @Summary      Create new permission combination
// @Description  Admin with permission create permission resource only
// @Tags         permission
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  body 			payloads.RequestCreatePermission	true	"data of new permission"
// @Success      201  		{object}  swaglib.Envelope{data=models.Permission,meta=swaglib.Info}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /permissions [post]
func (h *Permission) CreatePermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestCreatePermission
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	perm, errPayload := h.permService.ApplyContext(c).CreatePermission(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	// strings.Join can't process []models.PermissionAction
	permActs := slicelib.Map(perm.Actions, func(i int, act models.PermissionAction) string { return string(act) })
	rp.Success(perm).Info(fmt.Sprintf("new permission to %s %s created", strings.Join(permActs, ", "), perm.Resource)).CreatedJSON()
}

// @Summary      Get existing permission with id
// @Description  Admin with permission read permission resource only
// @Tags         permission
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "permission id"
// @Success      200  		{object}  swaglib.Envelope{data=models.Permission}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /permissions/{id} [get]
func (h *Permission) GetPermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	id := c.Param("id")

	perm, errPayload := h.permService.ApplyContext(c).GetPermission(id)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(perm).OkJSON()
}

// @Summary      Get existing permissions
// @Description  Admin with permission read permission resource only
// @Tags         permission
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  query			payloads.RequestGetPermissions	true	"config to accept permissions"
// @Success      200  		{array}  swaglib.Envelope{data=models.Permission,meta=swaglib.Pagination}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /permissions [get]
func (h *Permission) GetPermissions(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestGetPermissions
	c.ShouldBindQuery(&payload)

	perm, errPayload := h.permService.ApplyContext(c).GetPermissions(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(perm).PaginateCursor(config.LIMIT_PAGINATED_DATA, payload.Offset).OkJSON()
}

// @Summary      Update existing permission
// @Description  Admin with permission update permission resource only. Can not update permission seeds
// @Tags         permission
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "permission id"
// @Param				 payload  body			payloads.RequestUpdatePermission	true	"data to update permission"
// @Success      200  		{object}  swaglib.Envelope{data=models.Permission}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /permissions/{id} [put]
func (h *Permission) UpdatePermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestUpdatePermission
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}
	payload.ID = c.Param("id")

	perm, errPayload := h.permService.ApplyContext(c).UpdatePermission(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(perm).OkJSON()
}

// @Summary      Delete existing permission
// @Description  Admin with permission delete permission resource only. Can not deleet granted permission
// @Tags         permission
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "permission id"
// @Success      200  		{object}  swaglib.Envelope{data=models.Id}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /permissions/{id} [delete]
func (h *Permission) DeletePermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	id := c.Param("id")

	errPayload := h.permService.ApplyContext(c).DeletePermission(id)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(map[string]string{"id": id}).OkJSON()
}

// @Summary      Grant existing permission to another admin
// @Description  Admin with permission update permission resource only. Response granted user
// @Tags         permission
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  body			payloads.RequestGrantPermission	true	"data to grant permission"
// @Success      200  		{object}  swaglib.Envelope{data=models.User{admin_profile=models.Admin{permissions=[]models.Permission}},meta=swaglib.Info}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /permissions/grant [put]
func (h *Permission) GrantPermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestGrantPermission
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	user, permission, errPayload := h.permService.ApplyContext(c).GrantPermission(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	// strings.Join can't process []models.PermissionAction
	permActs := slicelib.Map(permission.Actions, func(i int, act models.PermissionAction) string { return string(act) })
	rp.Success(user).Info(fmt.Sprintf("%s's permitted to %s %s", user.FullName, strings.Join(permActs, ", "), permission.Resource)).OkJSON()
}

// @Summary      Revoke existing permission of another admin
// @Description  Admin with permission update permission resource only. Response granted user
// @Tags         permission
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  body			payloads.RequestRevokePermission	true	"data to revoke permission"
// @Success      200  		{object}  swaglib.Envelope{data=models.User{admin_profile=models.Admin{permissions=[]models.Permission}},meta=swaglib.Info}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /permissions/revoke [put]
func (h *Permission) RevokePermission(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestRevokePermission
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	user, permission, errPayload := h.permService.ApplyContext(c).RevokePermission(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	// strings.Join can't process []models.PermissionAction
	permActs := slicelib.Map(permission.Actions, func(i int, act models.PermissionAction) string { return string(act) })
	rp.Success(user).Info(fmt.Sprintf("%s's no longer permitted to %s %s", user.FullName, strings.Join(permActs, ", "), permission.Resource)).OkJSON()
}
