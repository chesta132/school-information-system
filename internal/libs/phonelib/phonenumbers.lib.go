package phonelib

import (
	"school-information-system/config"

	"github.com/nyaruka/phonenumbers"
)

func IsValidNumber(number string) bool {
	num, err := phonenumbers.Parse(number, config.APP_REGION)
	if err != nil {
		return false
	}

	return phonenumbers.IsValidNumber(num)
}

func FormatNumber(number string) (formattedNumber string, valid bool) {
	num, err := phonenumbers.Parse(number, config.APP_REGION)
	if err != nil {
		return "", false
	}
	if !phonenumbers.IsValidNumber(num) {
		return "", false
	}

	return phonenumbers.Format(num, phonenumbers.E164), true
}
