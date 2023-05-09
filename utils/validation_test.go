package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ValidationTestSuite struct {
	suite.Suite
}

func (suite *ValidationTestSuite) TestValidateUsername_Success() {
	username := "dummyUsername"
	validUsername := IsUsernameValid(username)

	assert.True(suite.T(), validUsername)
}

func (suite *ValidationTestSuite) TestValidateUsername_Failed() {
	username := "user"
	validUsername := IsUsernameValid(username)

	assert.False(suite.T(), validUsername)
}

func (suite *ValidationTestSuite) TestValidatePassword_Success() {
	password := "secretPass"
	validPassword := IsPasswordValid(password)

	assert.Equal(suite.T(), true, validPassword)
}

func (suite *ValidationTestSuite) TestValidPassword_FailedLength() {
	password := "pass"
	validPassword := IsPasswordValid(password)

	assert.Equal(suite.T(), false, validPassword)
}

func (suite *ValidationTestSuite) TestValidPassword_FailedCharacter() {
	password := "secretPass[]"
	validPassword := IsPasswordValid(password)

	assert.Equal(suite.T(), false, validPassword)
}

func (suite *ValidationTestSuite) TestValidatePhoneNumber_Success() {
	phoneNumber := "082123456789"
	validPhoneNumber := IsPhoneNumberValid(phoneNumber)

	assert.True(suite.T(), validPhoneNumber)
}

func (suite *ValidationTestSuite) TestValidatePhoneNumber_Failed() {
	phoneNumber := "0821abc"
	validPhoneNumber := IsPhoneNumberValid(phoneNumber)

	assert.False(suite.T(), validPhoneNumber)
}

func (suite *ValidationTestSuite) TestValidateEmail_Success() {
	email := "rmalariq@gmail.com"
	validEmail := IsEmailValid(email)

	assert.True(suite.T(), validEmail)
}

func (suite *ValidationTestSuite) TestValidateEmail_FailedFormat() {
	email := "rmalariq@.com"
	validEmail := IsEmailValid(email)

	assert.False(suite.T(), validEmail)
}

func (suite *ValidationTestSuite) TestValidateEmail_FailedCharacter() {
	email := "rmalariq@gmail<>.com"
	validEmail := IsEmailValid(email)

	assert.False(suite.T(), validEmail)
}

func TestRunValidationSuite(t *testing.T) {
	suite.Run(t, new(ValidationTestSuite))
}
