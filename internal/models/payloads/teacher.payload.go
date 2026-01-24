package payloads

type RequestUpdateTeacher struct {
	ID         string   `uri:"id" validate:"required,uuid4"`
	SubjectIDs []string `json:"subject_ids" validate:"omitempty,min=1"` // replace current subjects
	NUPTK      string   `example:"1234567890123456"`
}
