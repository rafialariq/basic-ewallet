package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHashingSuite struct {
	suite.Suite
	password string
}

func (suite *PasswordHashingSuite) SetupTest() {
	// Inisialisasi password
	suite.password = "password123"
}

func (suite *PasswordHashingSuite) TestPasswordHashing_Success() {
	// Panggil fungsi PasswordHashing dengan password yang sudah diinisialisasi
	hashedPassword := PasswordHashing(suite.password)

	// Verifikasi hasil hash tidak kosong
	assert.NotEmpty(suite.T(), hashedPassword)

	// Verifikasi hasil hash sama dengan password yang di-hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(suite.password))
	assert.Nil(suite.T(), err, "Hashed password doesn't match original password")
}

// func (suite *PasswordHashingSuite) TestPasswordHashing_Failed() {
// 	password := "password123"

// 	_, err := bcrypt.Cost([]byte{})
// 	suite.NoError(err)

// 	// Test when bcrypt.GenerateFromPassword returns an error
// 	hashedPassword := PasswordHashing(password + "%invalid_chars")
// 	assert.Equal(suite.T(), "", hashedPassword)
// 	log.Println("Failed to generate hashed password")
// }

func TestRunPasswordHashingSuite(t *testing.T) {
	suite.Run(t, new(PasswordHashingSuite))
}
