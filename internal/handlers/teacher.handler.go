package handlers

import (
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Teacher struct {
	teacherService *services.Teacher
}

func NewTeacher(teacherService *services.Teacher) *Teacher {
	return &Teacher{teacherService}
}

// @Summary      Update existing teacher profile
// @Description  Admin with permission update teacher resource only
// @Tags         teacher
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "teacher id"
// @Param				 payload  body			payloads.RequestUpdateTeacher	true	"data to update teacher"
// @Success      200  		{object}  swaglib.Envelope{data=models.Teacher{subjects=[]models.Subject}}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /teachers/{id} [put]
func (h *Teacher) UpdateTeacher(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestUpdateTeacher
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}
	payload.ID = c.Param("id")

	teacher, errPayload := h.teacherService.ApplyContext(c).UpdateTeacher(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(teacher).OkJSON()
}
