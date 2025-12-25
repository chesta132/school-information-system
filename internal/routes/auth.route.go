package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/repos"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterAuth(group *gin.RouterGroup) {
	userRepo := repos.NewUser(rt.db)
	revokedRepo := repos.NewRevoked(rt.db)

	authService := services.NewAuth(userRepo, revokedRepo)

	handler := handlers.NewAuth(authService)

	group.POST("/sign-up", handler.SignUp)
	group.POST("/sign-in", handler.SignIn)
	group.POST("/sign-out", handler.SignOut)
}
