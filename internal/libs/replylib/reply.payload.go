package replylib

import (
	"time"

	"github.com/chesta132/goreply/reply"
)

type Pagination struct {
	Current int  `json:"current"`  // current offset
	HasNext bool `json:"has_next"` // true if data more than replied
	Next    int  `json:"next"`     // if hasNext is false, next is 0
}

type Meta struct {
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Pagination  *Pagination       `json:"pagination,omitempty"`
	Information string            `json:"information,omitempty"`
	Tokens      map[string]string `json:"tokens,omitempty"`
	Debug       any               `json:"debug,omitempty"`
}

type Envelope struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data"`
}

func transformer(rp *reply.Reply) any {
	meta := rp.Meta()
	data := rp.Data()

	transMeta := Meta{
		Status:      meta.Status,
		Timestamp:   time.Unix(meta.Timestamp, 0),
		Information: meta.Info,
		Tokens:      meta.Tokens,
		Debug:       meta.Debug,
	}

	if meta.Pagination != nil {
		transMeta.Pagination = &Pagination{
			Current: meta.Pagination.Current,
			HasNext: meta.Pagination.HasNext,
			Next:    meta.Pagination.Next,
		}
	}

	return &Envelope{Data: data, Meta: transMeta}
}
