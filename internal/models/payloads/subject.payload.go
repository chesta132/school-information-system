package payloads

type RequestCreateSubject struct {
	Name string `json:"name" validate:"required"`
}
