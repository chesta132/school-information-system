package handlers

import (
	"school-information-system/config"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/models"
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
func (h *Class) CreateClass(c *gin.Context) {
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

// @Summary      Get existing class with id
// @Description  Admin with permission read class resource or teacher only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "class id"
// @Success      200  		{object}  swaglib.Envelope{data=models.Class}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes/{id} [get]
func (h *Class) GetClass(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestGetClass
	c.ShouldBindUri(&payload)

	class, errPayload := h.classService.ApplyContext(c).GetClass(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(class).OkJSON()
}

// @Summary      Get existing classes
// @Description  Admin with permission read class resource or teacher only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param				 payload  query			payloads.RequestGetClasses	true	"config to accept classes"
// @Success      200  		{object}  swaglib.Envelope{data=[]models.Class}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes [get]
func (h *Class) GetClasses(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestGetClasses
	c.ShouldBindQuery(&payload)

	if payload.Offset < 0 {
		payload.Offset = 0
	}

	classes, errPayload := h.classService.ApplyContext(c).GetClasses(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(classes).PaginateCursor(config.LIMIT_PAGINATED_DATA, payload.Offset).OkJSON()
}

// @Summary      Update class
// @Description  Admin with permission update class resource only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "class id"
// @Param				 payload  body 			payloads.RequestUpdateClass	true	"data to update class"
// @Success      200  		{object}  swaglib.Envelope{data=models.Class}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes/{id} [put]
func (h *Class) UpdateClasss(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestUpdateClass
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}
	payload.ID = c.Param("id")

	class, errPayload := h.classService.ApplyContext(c).UpdateClass(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(class).OkJSON()
}

// @Summary      Delete class
// @Description  Admin with permission delete class resource only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "class id"
// @Success      200  		{object}  swaglib.Envelope{data=models.Id}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes/{id} [delete]
func (h *Class) DeleteClass(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestDeleteClass
	c.ShouldBindUri(&payload)

	errPayload := h.classService.ApplyContext(c).DeleteClass(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(&models.Id{ID: payload.ID}).OkJSON()
}

// ------------------------------------ //
// ------------------------------------ //
// --------------RELATION-------------- //
// ------------------------------------ //
// ------------------------------------ //

// @Summary      Get form teacher in class
// @Description  Admin with permission read class and teacher resource only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "class id"
// @Success      200  		{object}  swaglib.Envelope{data=models.User}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes/{id}/form-teacher [get]
func (h *Class) GetFormTeacher(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestGetClass
	c.ShouldBindUri(&payload)

	teacher, errPayload := h.classService.ApplyContext(c).GetFormTeacher(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(teacher).OkJSON()
}

// @Summary      Get students in class
// @Description  Admin with permission read class and student resource only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "class id"
// @Success      200  		{object}  swaglib.Envelope{data=[]models.User}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes/{id}/students [get]
func (h *Class) GetStudents(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestGetClass
	c.ShouldBindUri(&payload)

	students, errPayload := h.classService.ApplyContext(c).GetStudents(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(students).OkJSON()
}

// @Summary      Get class, form teacher, and students
// @Description  Admin with permission read class, teacher, and student resource only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "class id"
// @Success      200  		{object}  swaglib.Envelope{data=payloads.ResponseGetFullClass}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes/{id}/full [get]
func (h *Class) GetFull(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestGetClass
	c.ShouldBindUri(&payload)

	full, errPayload := h.classService.ApplyContext(c).GetFull(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(full).OkJSON()
}

// @Summary      Set form teacher of class
// @Description  Admin with permission update class and teacher resource only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "class id"
// @Param				 payload  body 			payloads.RequestSetFormTeacher	true	"data to set form teacher"
// @Success      200  		{object}  swaglib.Envelope{data=models.User}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes/{id}/form-teacher [put]
func (h *Class) SetFormTeacher(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestSetFormTeacher
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}
	c.ShouldBindUri(&payload)

	teacher, errPayload := h.classService.ApplyContext(c).SetFormTeacher(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(teacher).OkJSON()
}

// @Summary      Add students in class
// @Description  Admin with permission update class and student resource only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "class id"
// @Param				 payload  body 			payloads.RequestAddStudents	true	"data to add students"
// @Success      200  		{object}  swaglib.Envelope{data=[]models.User}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes/{id}/students [post]
func (h *Class) AddStudents(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestAddStudents
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}
	c.ShouldBindUri(&payload)

	students, errPayload := h.classService.ApplyContext(c).AddStudents(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(students).OkJSON()
}

// @Summary      Remove form teacher of class
// @Description  Admin with permission update class and teacher resource only
// @Tags         class
// @Accept       json
// @Produce      json
// @Param				 Cookie   header 		string 	false	"access_token"
// @Param				 Cookie2  header 		string 	true	"refresh_token"
// @Param 			 id				path 			string  true  "class id"
// @Success      200  		{object}  swaglib.Envelope{data=models.User}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /classes/{id}/form-teacher [delete]
func (h *Class) RemoveFormTeacher(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestGetClass
	c.ShouldBindUri(&payload)

	teacher, errPayload := h.classService.ApplyContext(c).RemoveFormTeacher(payload)
	if errPayload != nil {
		rp.Error(replylib.ErrorPayloadToArgs(errPayload)).FailJSON()
		return
	}

	rp.Success(teacher).OkJSON()
}
