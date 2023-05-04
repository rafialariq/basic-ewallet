package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHashing(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to generate hashed password")
	}

	return string(hashedPassword)
}
