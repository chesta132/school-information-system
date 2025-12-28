package replylib

import (
	"school-information-system/config"

	"github.com/chesta132/goreply/reply"
)

var Client = reply.NewClient(reply.Client{
	CodeAliases: CodeAliases,
	Transformer: transformer,
	DebugMode:   config.IsEnvDev(),
})
