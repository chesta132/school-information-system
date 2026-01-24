package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/models"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterTeacher(group *gin.RouterGroup) {
	teacherService := services.NewTeacher(rt.rp.Teacher(), rt.rp.Subject())
	handler := handlers.NewTeacher(teacherService)

	mw := middlewares.NewAuth(rt.rp.User(), rt.rp.Revoked())

	group.PUT("/:id", mw.PermissionProtected(
		models.ResourceTeacher,
		[]models.PermissionAction{models.ActionUpdate},
	), handler.UpdateTeacher)
}
