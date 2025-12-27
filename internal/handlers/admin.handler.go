package handlers

import (
	"fmt"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Admin struct {
	adminService      *services.Admin
	roleSetterService *services.RoleSetter
}

func NewAdmin(adminService *services.Admin, roleSetterService *services.RoleSetter) *Admin {
	return &Admin{adminService, roleSetterService}
}

// @Summary      Initiates new admin
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  body			payloads.RequestInitiateAdmin	true	"data of initiated admin"
// @Success      200  		{object}  swaglib.Envelope{data=models.User{admin_profile=models.Admin{permissions=[]models.Permission}},meta=swaglib.Info}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /admins/initiate [post]
func (h *Admin) InitiateAdmin(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestInitiateAdmin
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}
	user, errPayload := h.adminService.ApplyContext(c).InitiateAdmin(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(user).Info(fmt.Sprintf("Admin created: %s (%s)", user.FullName, payload.StaffRole)).OkJSON()
}

// @Summary      Set another user's role
// @Description  Admin with permission update role resource only
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  body			payloads.RequestSetRole	true	"data of targeted user's role, only insert data that match with target_role"
// @Success      200  		{object}  swaglib.Envelope{data=models.User{admin_profile=models.Admin,student_profile=models.Student{class=models.Class,parents=[]models.Parent},teacher_profile=models.Teacher{subjects=[]models.Subject}},meta=swaglib.Info} "*_data is match with role"
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /admins/set-role [put]
func (h *Admin) SetRole(c *gin.Context) {
	rp := replylib.Client.New(adapter.AdaptGin(c))
	var payload payloads.RequestSetRole
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	user, errPayload := h.roleSetterService.ApplyContext(c).SetRole(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(user).Info(fmt.Sprintf("%s's role setted", user.FullName)).OkJSON()
}
