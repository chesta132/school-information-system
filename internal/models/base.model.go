package models

type ID struct {
	ID string `gorm:"default:gen_random_uuid()" json:"id"`
}
