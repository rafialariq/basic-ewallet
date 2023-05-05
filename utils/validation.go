package utils

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {

	for _, c := range email {
		if c < 31 || c > 127 || strings.ContainsAny(string(c), `()<>,;:\\"[]`) {
			return true
		}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return true
	}

	return false
}

func ValidateUsername(username string) bool {
	if len(username) < 6 || len(username) > 20 {
		return false
	}

	return true
}

func ValidatePhoneNumber(phone string) bool {
	phone = strings.ReplaceAll(phone, " ", "")

	if len(phone) < 10 || len(phone) > 14 {
		return false
	}

	var validPhone = regexp.MustCompile(`^[0-9]+$`).MatchString(phone)

	return validPhone
}
