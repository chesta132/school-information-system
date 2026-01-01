package handlers

import (
	"school-information-system/config"
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

	parent, errPayload := h.parentService.ApplyContext(c).CreateParent(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(parent).Info("new parent profile created").CreatedJSON()
}

// @Summary      Get existing parent with id
// @Description  Admin with permission read parent resource only
// @Tags         parent
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "parent id"
// @Success      200  		{object}  swaglib.Envelope{data=models.Parent}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /parents/{id} [get]
func (h *Parent) GetParent(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestGetParent
	c.ShouldBindUri(&payload)

	parent, errPayload := h.parentService.ApplyContext(c).GetParent(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(parent).OkJSON()
}

// @Summary      Get existing parents
// @Description  Admin with permission read parent resource only
// @Tags         permission
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  query			payloads.RequestGetPermissions	true	"config to accept parents"
// @Success      200  		{array}  	swaglib.Envelope{data=[]models.Parent,meta=swaglib.Pagination}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /parents [get]
func (h *Parent) GetParents(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestGetParents
	c.ShouldBindQuery(&payload)

	if payload.Offset < 0 {
		payload.Offset = 0
	}

	parents, errPayload := h.parentService.ApplyContext(c).GetParents(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(parents).PaginateCursor(config.LIMIT_PAGINATED_DATA, payload.Offset).OkJSON()
}

// @Summary      Update existing parent
// @Description  Admin with permission update parent resource only
// @Tags         parent
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "parent id"
// @Param				 payload  body			payloads.RequestUpdateParent	true	"data to update parent"
// @Success      200  		{object}  swaglib.Envelope{data=models.Parent}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /parents/{id} [put]
func (h *Parent) UpdateParent(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestUpdateParent
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}
	payload.ID = c.Param("id")

	parent, errPayload := h.parentService.ApplyContext(c).UpdateParent(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(parent).OkJSON()
}
