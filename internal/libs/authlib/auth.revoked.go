package authlib

import "school-information-system/internal/models"

const (
	ReasonUserSignOut = "user sign out"
)

var revokedMessages = map[string]string{
	ReasonUserSignOut: "user already signed out",
}

func MessageOfRevoke(revoked models.Revoked) string {
	if msg, ok := revokedMessages[revoked.Reason]; ok {
		return msg
	}
	return revoked.Reason
}
