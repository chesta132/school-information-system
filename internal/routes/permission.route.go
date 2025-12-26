package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/models"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterPermission(group *gin.RouterGroup) {
	permService := services.NewPermission(rt.rp.User(), rt.rp.Permission())

	handler := handlers.NewPermission(permService)

	mw := middlewares.NewAuth(rt.rp.User(), rt.rp.Revoked())

	// create
	group.POST("/", mw.PermissionProtected(
		models.ResourcePermission,
		[]models.PermissionAction{models.ActionCreate},
	), handler.CreatePermission)
	// get 1
	group.GET("/:id", mw.PermissionProtected(
		models.ResourcePermission,
		[]models.PermissionAction{models.ActionRead},
	), handler.GetPermission)
	// get many
	group.GET("/", mw.PermissionProtected(
		models.ResourcePermission,
		[]models.PermissionAction{models.ActionRead},
	), handler.GetPermissions)
	// update
	group.PUT("/:id", mw.PermissionProtected(
		models.ResourcePermission,
		[]models.PermissionAction{models.ActionUpdate},
	), handler.UpdatePermission)
	// delete
	group.DELETE("/:id", mw.PermissionProtected(
		models.ResourcePermission,
		[]models.PermissionAction{models.ActionDelete},
	), handler.DeletePermission)

	// grant & revoke
	group.PUT("/grant", mw.PermissionProtected(
		models.ResourcePermission,
		[]models.PermissionAction{models.ActionUpdate},
	), handler.GrantPermission)
	group.DELETE("/revoke", mw.PermissionProtected(
		models.ResourcePermission,
		[]models.PermissionAction{models.ActionUpdate},
	), handler.RevokePermission)
}
