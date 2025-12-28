// for docs only lib
package swaglib

import "time"

type status string

const (
	success status = "SUCCESS"
	error   status = "ERROR"
)

type Pagination struct {
	Pagination struct {
		Current int  `json:"current"`  // current offset
		HasNext bool `json:"has_next"` // true if data more than replied
		Next    int  `json:"next"`     // if hasNext is false, next is 0
	} `json:"pagination"`
}

type Token struct {
	Tokens map[string]string `json:"tokens"`
}

type Info struct {
	Information string `json:"information" example:"resource successfully updated"` // information message
}

type Meta struct {
	Status    status    `json:"status"`
	Timestamp time.Time `json:"timestamp" example:"2006-01-02T15:04:05Z07:00"` // in UTC
	Debug     string    `json:"debug" example:"inconsistent value"`            // please dont process debug fields because it inconsistently
}

type Envelope struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data" extensions:"x-null-if-{}"`
}
