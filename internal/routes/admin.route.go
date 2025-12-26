package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/models"
	"school-information-system/internal/repos"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterAdmin(group *gin.RouterGroup) {
	userRepo := repos.NewUser(rt.db)
	adminRepo := repos.NewAdmin(rt.db)
	revokedRepo := repos.NewRevoked(rt.db)
	// role setter repos
	studentRepo := repos.NewStudent(rt.db)
	classRepo := repos.NewClass(rt.db)
	parentRepo := repos.NewParent(rt.db)
	teacherRepo := repos.NewTeacher(rt.db)
	subjectRepo := repos.NewSubject(rt.db)

	adminService := services.NewAdmin(userRepo, adminRepo)
	roleSetterService := services.NewRoleSetter(userRepo, adminRepo, studentRepo, classRepo, parentRepo, teacherRepo, subjectRepo)

	handler := handlers.NewAdmin(adminService, roleSetterService)
	mw := middlewares.NewAuth(userRepo, revokedRepo)

	group.POST("/initiate", mw.Protected(true), handler.InitiateAdmin)

	group.Use(mw.RoleProtected(models.RoleAdmin))

	group.POST(
		"/set-role",
		mw.PermissionProtected(models.ResourceRole, []models.PermissionAction{models.ActionRead, models.ActionUpdate}),
		handler.SetRole,
	)

	rt.RegisterPermission(group.Group("/permissions"))
}
