package replylib

import (
	"github.com/chesta132/goreply/reply"
)

var (
	// phone number errors

	ErrInvalidPhone = reply.ErrorPayload{
		Code:    CodeBadRequest,
		Message: "invalid payload",
		Fields:  reply.FieldsError{"phone": "invalid phone number"},
	}
	ErrPhoneRegistered = reply.ErrorPayload{
		Code:    CodeConflict,
		Message: "invalid payload",
		Fields:  reply.FieldsError{"phone": "phone number already registered"},
	}

	// email errors

	ErrEmailRegistered = reply.ErrorPayload{
		Code:    CodeConflict,
		Message: "invalid payload",
		Fields:  reply.FieldsError{"email": "email already registered"},
	}
	ErrEmailNotRegistered = reply.ErrorPayload{
		Code:    CodeNotFound,
		Message: "invalid payload",
		Fields:  reply.FieldsError{"email": "email not registered yet"},
	}

	// password errors

	ErrIncorrectPassword = reply.ErrorPayload{
		Code:    CodeUnauthorized,
		Message: "invalid payload",
		Fields:  reply.FieldsError{"password": "password is incorrect"},
	}
	ErrIncorrectKey = reply.ErrorPayload{
		Code:    CodeUnauthorized,
		Message: "invalid payload",
		Fields:  reply.FieldsError{"key": "invalid key"},
	}

	// presence/absence error

	ErrAdminExist = reply.ErrorPayload{
		Code:    CodeConflict,
		Message: "admin is already exist",
	}
	ErrPermissionNameExist = reply.ErrorPayload{
		Code:    CodeConflict,
		Message: "invalid payload",
		Fields:  reply.FieldsError{"name": "another permission with this name already registered"},
	}
	ErrPermissionImmutable = reply.ErrorPayload{
		Code:    CodeUnprocessableEntity,
		Message: "this permission is immutable",
	}
)

func ErrorPayloadToArgs(errPayload *reply.ErrorPayload) (string, string, reply.ErrorOption, reply.ErrorOption) {
	return errPayload.Code, errPayload.Message, reply.WithDetails(errPayload.Details), reply.WithFields(errPayload.Fields)
}
