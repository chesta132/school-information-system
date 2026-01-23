package payloads

type RequestUpdateStudent struct {
	ID        string   `uri:"id" validate:"required,uuid4"`
	NISN      string   `example:"0091913711"`
	ParentIDs []string `json:"parent_ids" validate:"omitempty,min=2,max=2"` // replace current parents
}
