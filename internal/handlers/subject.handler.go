package handlers

import (
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Subject struct {
	subjectService *services.Subject
}

func NewSubject(subjectService *services.Subject) *Subject {
	return &Subject{subjectService}
}

// @Summary      Create new subject
// @Description  Admin with permission create subject resource only
// @Tags         subject
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  body 			payloads.RequestCreateSubject	true	"data of new subject"
// @Success      201  		{object}  swaglib.Envelope{data=models.Subject,meta=swaglib.Info}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /subjects [post]
func (h *Subject) CreateSubject(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestCreateSubject
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	subject, errPayload := h.subjectService.ApplyContext(c).CreateSubject(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(subject).Info("new subject " + subject.Name + " created").CreatedJSON()
}
