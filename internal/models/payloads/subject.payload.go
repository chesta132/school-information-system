package payloads

type RequestCreateSubject struct {
	Name string `json:"name" validate:"required"`
}

type RequestGetSubject struct {
	ID string `uri:"id" json:"id" validate:"required,uuid4"`
}

type RequestGetSubjects struct {
	Offset int    `form:"offset" example:"10"`
	Query  string `form:"q" example:"infor"`
}

type RequestUpdateSubject struct {
	ID   string `uri:"id" json:"id" validate:"required,uuid4"`
	Name string `json:"name" validate:"required"`
}
