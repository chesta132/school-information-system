package handlers

import (
	"school-information-system/config"
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

// @Summary      Get existing subject with id
// @Description  Admin with permission read subject resource or teacher only
// @Tags         subject
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "subject id"
// @Success      200  		{object}  swaglib.Envelope{data=models.Subject}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /subjects/{id} [get]
func (h *Subject) GetSubject(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestGetSubject
	c.ShouldBindUri(&payload)

	subject, errPayload := h.subjectService.ApplyContext(c).GetSubject(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(subject).OkJSON()
}

// @Summary      Get existing subjects
// @Description  Admin with permission read subject resource or teacher only
// @Tags         subject
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  query			payloads.RequestGetSubjects	true	"config to accept subjects"
// @Success      200  		{object}  swaglib.Envelope{data=[]models.Subject,meta=swaglib.Pagination}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /subjects [get]
func (h *Subject) GetSubjects(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestGetSubjects
	c.ShouldBindQuery(&payload)

	subjects, errPayload := h.subjectService.ApplyContext(c).GetSubjects(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(subjects).PaginateCursor(config.LIMIT_PAGINATED_DATA, payload.Offset).OkJSON()
}
