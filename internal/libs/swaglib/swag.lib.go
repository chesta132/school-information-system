package swaglib

type status string

const (
	success status = "SUCCESS"
	error   status = "ERROR"
)

type Pagination struct {
	Pagination struct {
		Current int  `json:"current"` // current offset
		HasNext bool `json:"hasNext"` // true if data more than replied
		Next    int  `json:"next"`    // if hasNext is false, next is 0
	} `json:"pagination"`
}

type Token struct {
	Token map[string]string `json:"token"`
}

type Info struct {
	Information string `json:"information" example:"resource successfully updated"` // information message
}

type Meta struct {
	Status status `json:"status"`
}

type Envelope struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data" extensions:"x-null-if-{}"`
}
