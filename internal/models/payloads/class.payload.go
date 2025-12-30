package payloads

type RequestCreateClass struct {
	Grade         int    `json:"grade" validate:"required,min=1,max=12" example:"10"`
	Major         string `json:"major" validate:"required" example:"TJKT"`
	ClassNumber   int    `json:"class_number" validate:"required,min=1" example:"3"`
	FormTeacherID string `json:"form_teacher_id" validate:"required,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
}

type RequestGetClass struct {
	ID string `uri:"id" validate:"required,uuid4"`
}

type RequestGetClasses struct {
	Offset      int    `form:"offset" example:"10"`
	Grade       int    `form:"grade" example:"10"`
	Major       string `form:"major" example:"TJKT"`
	ClassNumber int    `form:"class_number" example:"3"`
}
