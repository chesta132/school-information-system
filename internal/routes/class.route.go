package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/models"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterClass(group *gin.RouterGroup) {
	classService := services.NewClass(rt.rp.Class(), rt.rp.Teacher())

	handler := handlers.NewClass(classService)

	mw := middlewares.NewAuth(rt.rp.User(), rt.rp.Revoked())

	group.POST("/", mw.PermissionProtected(
		models.ResourceClass,
		[]models.PermissionAction{models.ActionCreate},
	), handler.CreateSubject)

	group.GET("/:id", mw.PermissionProtected(
		models.ResourceClass,
		[]models.PermissionAction{models.ActionRead},
		middlewares.WithSkipRole(models.RoleTeacher),
	), handler.GetClass)
}
