package models

type Subject struct {
	Id
	Name string `gorm:"unique;not null" json:"name" example:"informatika"`

	Timestamp
}
