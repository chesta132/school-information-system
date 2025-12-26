package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/models"
	"school-information-system/internal/repos"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterPermission(group *gin.RouterGroup) {
	userRepo := repos.NewUser(rt.db)
	permRepo := repos.NewPermission(rt.db)
	revokedRepo := repos.NewRevoked(rt.db)

	permService := services.NewPermission(userRepo, permRepo)

	handler := handlers.NewPermission(permService)

	mw := middlewares.NewAuth(userRepo, revokedRepo)

	group.POST("/grant", mw.PermissionProtected(
		models.ResourcePermission,
		[]models.PermissionAction{models.ActionRead, models.ActionUpdate},
	), handler.GrantPermission)
}
