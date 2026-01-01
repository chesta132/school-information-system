package routes

import (
	"school-information-system/internal/handlers"
	"school-information-system/internal/middlewares"
	"school-information-system/internal/models"
	"school-information-system/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Route) RegisterParent(group *gin.RouterGroup) {
	parentService := services.NewParent(rt.rp.Parent())

	handler := handlers.NewParent(parentService)

	mw := middlewares.NewAuth(rt.rp.User(), rt.rp.Revoked())

	group.POST("/", mw.PermissionProtected(
		models.ResourceParent,
		[]models.PermissionAction{models.ActionCreate},
	), handler.CreateParent)

	group.GET("/:id", mw.PermissionProtected(
		models.ResourceParent,
		[]models.PermissionAction{models.ActionRead},
	), handler.GetParent)
}
