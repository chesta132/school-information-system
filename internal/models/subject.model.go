package models

type Subject struct {
	Id
	Name string `json:"name"`

	Timestamp
}
