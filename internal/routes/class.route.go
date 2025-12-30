package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/models"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterClass(group *gin.RouterGroup) {
	classService := services.NewClass(rt.rp.Class(), rt.rp.User(), rt.rp.Teacher(), rt.rp.Student())

	handler := handlers.NewClass(classService)

	mw := middlewares.NewAuth(rt.rp.User(), rt.rp.Revoked())

	group.POST("/", mw.PermissionProtected(
		models.ResourceClass,
		[]models.PermissionAction{models.ActionCreate},
	), handler.CreateClass)

	group.GET("/:id", mw.PermissionProtected(
		models.ResourceClass,
		[]models.PermissionAction{models.ActionRead},
		middlewares.WithSkipRole(models.RoleTeacher),
	), handler.GetClass)
	group.GET("/", mw.PermissionProtected(
		models.ResourceClass,
		[]models.PermissionAction{models.ActionRead},
		middlewares.WithSkipRole(models.RoleTeacher),
	), handler.GetClasses)

	group.PUT("/:id", mw.PermissionProtected(
		models.ResourceClass,
		[]models.PermissionAction{models.ActionUpdate},
	), handler.UpdateClasss)

	group.DELETE("/:id", mw.PermissionProtected(
		models.ResourceClass,
		[]models.PermissionAction{models.ActionDelete},
	), handler.DeleteClass)

	// RELATION

	group.GET("/:id/form-teacher", mw.PermissionProtected(
		models.ResourceClass,
		[]models.PermissionAction{models.ActionRead},
	), mw.PermissionProtected(
		models.ResourceTeacher,
		[]models.PermissionAction{models.ActionRead},
	), handler.GetFormTeacher)
}
