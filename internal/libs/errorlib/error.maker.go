package errorlib

import (
	"errors"
	"school-information-system/internal/libs/replylib"

	"github.com/chesta132/goreply/reply"
	"gorm.io/gorm"
)

func MakeUserByTargetIDNotFound(err error) *reply.ErrorPayload {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &reply.ErrorPayload{
			Code:    replylib.CodeNotFound,
			Message: "user with targeted id doesn't exist",
			Fields:  []string{"target_id"},
		}
	}
	return &reply.ErrorPayload{
		Code:    replylib.CodeServerError,
		Message: err.Error(),
	}
}
