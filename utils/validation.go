package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func IsUsernameValid(username string) bool {
	minUsername, _ := strconv.Atoi(DotEnv("MIN_UNAME", ".env"))
	maxUsername, _ := strconv.Atoi(DotEnv("MAX_UNAME", ".env"))

	if len(username) < minUsername || len(username) > maxUsername {
		return false
	}

	return true
}

func IsPasswordValid(password string) bool {
	minPassword, _ := strconv.Atoi(DotEnv("MIN_PASS", ".env"))
	maxPassword, _ := strconv.Atoi(DotEnv("MAX_PASS", ".env"))

	if len(password) <= minPassword || len(password) >= maxPassword {
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
	minPhoneNum, _ := strconv.Atoi(DotEnv("MIN_PHONE_NUM", ".env"))
	maxPhoneNum, _ := strconv.Atoi(DotEnv("MAX_PHONE_NUM", ".env"))
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")

	if len(phoneNumber) < minPhoneNum || len(phoneNumber) > maxPhoneNum {
		return false
	}

	var validPhone = regexp.MustCompile(`^[0-9]+$`).MatchString(phoneNumber)

	return validPhone
}
