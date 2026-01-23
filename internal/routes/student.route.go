package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/models"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterStudent(group *gin.RouterGroup) {
	studentService := services.NewStudent(rt.rp.Student(), rt.rp.Parent())
	handler := handlers.NewStudent(studentService)

	mw := middlewares.NewAuth(rt.rp.User(), rt.rp.Revoked())

	group.PUT("/:id", mw.PermissionProtected(
		models.ResourceStudent,
		[]models.PermissionAction{models.ActionUpdate},
	), handler.UpdateStudent)
}
