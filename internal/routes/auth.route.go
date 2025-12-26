package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterAuth(group *gin.RouterGroup) {
	authService := services.NewAuth(rt.rp.User(), rt.rp.Revoked())

	handler := handlers.NewAuth(authService)

	group.POST("/sign-up", handler.SignUp)
	group.POST("/sign-in", handler.SignIn)
	group.POST("/sign-out", handler.SignOut)
}
