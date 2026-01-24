package handlers

import (
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Student struct {
	studentService *services.Student
}

func NewStudent(studentService *services.Student) *Student {
	return &Student{studentService}
}

// @Summary      Update existing student profile
// @Description  Admin with permission update student resource only
// @Tags         student
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "student id"
// @Param				 payload  body			payloads.RequestUpdateStudent	true	"data to update student"
// @Success      200  		{object}  swaglib.Envelope{data=models.Student{parents=[]models.Parent}}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /students/{id} [put]
func (h *Student) UpdateStudent(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestUpdateStudent
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}
	payload.ID = c.Param("id")

	student, errPayload := h.studentService.ApplyContext(c).UpdateStudent(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(student).OkJSON()
}
