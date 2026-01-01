package handlers

import (
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Parent struct {
	parentService *services.Parent
}

func NewParent(parentService *services.Parent) *Parent {
	return &Parent{parentService}
}

// @Summary      Create new parent profile
// @Description  Admin with permission create parent resource only
// @Tags         parent
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  body 			payloads.RequestCreateParent	true	"data of new parent"
// @Success      201  		{object}  swaglib.Envelope{data=models.Parent,meta=swaglib.Info}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /parents [post]
func (h *Parent) CreateParent(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestCreateParent
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	perm, errPayload := h.parentService.ApplyContext(c).CreateParent(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(perm).Info("new parent profile created").CreatedJSON()
}
