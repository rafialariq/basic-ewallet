package utils

import (
	"regexp"
	"strings"
)

func IsUsernameValid(username string) bool {
	if len(username) < 6 || len(username) > 20 {
		return false
	}

	return true
}

func IsPasswordValid(password string) bool {

	if len(password) <= 8 || len(password) >= 20 {
		return false
	} else if strings.ContainsAny(password, ` ^*+=-_()<>,;:\\"[]`) {
		return false
	}

	return true
}

func IsEmailValid(email string) bool {
	for _, c := range email {
		if c < 31 || c > 127 || strings.ContainsAny(string(c), `()<>,;:\\"[]`) {
			return false
		}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsPhoneNumberValid(phoneNumber string) bool {
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")

	if len(phoneNumber) < 10 || len(phoneNumber) > 14 {
		return false
	}

	var validPhone = regexp.MustCompile(`^[0-9]+$`).MatchString(phoneNumber)

	return validPhone
}
