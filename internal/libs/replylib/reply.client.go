package replylib

import "github.com/chesta132/goreply/reply"

var Client = reply.NewClient(reply.Client{
	CodeAliases: CodeAliases,
})