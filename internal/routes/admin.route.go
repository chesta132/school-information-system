package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/repos"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterAdmin(group *gin.RouterGroup) {
	userRepo := repos.NewUser(rt.db)
	adminRepo := repos.NewAdmin(rt.db)
	revokedRepo := repos.NewRevoked(rt.db)

	adminService := services.NewAdmin(userRepo, adminRepo)

	handler := handlers.NewAdmin(adminService)
	mw := middlewares.NewAuth(revokedRepo)

	group.POST("/initiate", mw.Protected(true), handler.InitiateAdmin)
}
