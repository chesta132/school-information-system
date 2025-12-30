package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/models"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterSubject(group *gin.RouterGroup) {
	subjectService := services.NewSubject(rt.rp.Subject(), rt.rp.Teacher())
	handler := handlers.NewSubject(subjectService)

	mw := middlewares.NewAuth(rt.rp.User(), rt.rp.Revoked())

	group.Use(mw.RoleProtected(models.RoleAdmin, models.RoleTeacher))

	group.POST("/", mw.PermissionProtected(
		models.ResourceSubject, []models.PermissionAction{models.ActionCreate},
	), handler.CreateSubject)

	group.GET("/:id", mw.PermissionProtected(
		models.ResourceSubject,
		[]models.PermissionAction{models.ActionRead},
		middlewares.WithSkipRole(models.RoleTeacher),
	), handler.GetSubject)
	group.GET("/", mw.PermissionProtected(
		models.ResourceSubject,
		[]models.PermissionAction{models.ActionRead},
		middlewares.WithSkipRole(models.RoleTeacher),
	), handler.GetSubjects)

	group.PUT("/:id", mw.PermissionProtected(
		models.ResourceSubject,
		[]models.PermissionAction{models.ActionUpdate},
	), handler.UpdateSubject)

	group.DELETE("/:id", mw.PermissionProtected(
		models.ResourceSubject,
		[]models.PermissionAction{models.ActionDelete},
	), handler.DeleteSubject)

	group.GET("/:id/teachers", mw.PermissionProtected(
		models.ResourceSubject,
		[]models.PermissionAction{models.ActionRead},
		middlewares.WithSkipRole(models.RoleTeacher),
	), mw.PermissionProtected(
		models.ResourceTeacher,
		[]models.PermissionAction{models.ActionRead},
		middlewares.WithSkipRole(models.RoleTeacher),
	), handler.GetTeacherOfSubject)
}
