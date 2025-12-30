package handlers

import (
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Class struct {
	classService *services.Class
}

func NewClass(classService *services.Class) *Class {
	return &Class{classService}
}

// @Summary      Create new class
// @Description  Admin with permission create class resource only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  body 			payloads.RequestCreateClass	true	"data of new class"
// @Success      201  		{object}  swaglib.Envelope{data=models.Class,meta=swaglib.Info}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes [post]
func (h *Class) CreateSubject(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestCreateClass
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	class, errPayload := h.classService.ApplyContext(c).CreateClass(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(class).Info("new class " + class.Name + " created").CreatedJSON()
}
