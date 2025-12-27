package models

type Subject struct {
	Id
	Name string `json:"name" example:"informatika"`

	Timestamp
}
