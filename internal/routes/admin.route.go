package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/models"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterAdmin(group *gin.RouterGroup) {
	adminService := services.NewAdmin(rt.rp.User(), rt.rp.Admin())
	roleSetterService := services.NewRoleSetter(
		rt.rp.User(),
		rt.rp.Admin(),
		rt.rp.Student(),
		rt.rp.Class(),
		rt.rp.Parent(),
		rt.rp.Teacher(),
		rt.rp.Subject(),
	)

	handler := handlers.NewAdmin(adminService, roleSetterService)
	mw := middlewares.NewAuth(rt.rp.User(), rt.rp.Revoked())

	group.POST("/initiate", mw.Protected(true), handler.InitiateAdmin)

	group.Use(mw.RoleProtected(models.RoleAdmin))

	group.PUT(
		"/set-role",
		mw.PermissionProtected(models.ResourceRole, []models.PermissionAction{models.ActionRead, models.ActionUpdate}),
		handler.SetRole,
	)

	rt.RegisterPermission(group.Group("/permissions"))
}
