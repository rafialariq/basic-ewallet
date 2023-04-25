package utils

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	// Validasi format email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}

	// Validasi karakter khusus
	for _, c := range email {
		if c <= 31 || c >= 127 || strings.ContainsAny(string(c), `()<>@,;:\\".[]`) {
			return false
		}
	}

	// Validasi domain
	parts := strings.Split(email, "@")
	domain := parts[len(parts)-1]
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") || strings.Count(domain, ".") < 1 {
		return false
	}

	return true
}

func ValidateUsername(username string) bool {
	if len(username) < 6 || len(username) > 20 {
		return false
	}

	return true
}

func ValidatePhoneNumber(phone string) bool {
	// Menghilangkan spasi pada nomor telepon
	phone = strings.ReplaceAll(phone, " ", "")

	if len(phone) < 10 || len(phone) > 14 {
		return false
	}

	// Validasi nomor telepon hanya terdiri dari karakter angka
	var validPhone = regexp.MustCompile(`^[0-9]+$`).MatchString(phone)

	return validPhone
}
