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
	validUsername := ValidateUsername(username)

	assert.True(suite.T(), validUsername)
}

func (suite *ValidationTestSuite) TestValidateUsername_Failed() {
	username := "user"
	validUsername := ValidateUsername(username)

	assert.False(suite.T(), validUsername)
}

func (suite *ValidationTestSuite) TestValidatePhoneNumber_Success() {
	phoneNumber := "082123456789"
	validPhoneNumber := ValidatePhoneNumber(phoneNumber)

	assert.True(suite.T(), validPhoneNumber)
}

func (suite *ValidationTestSuite) TestValidatePhoneNumber_Failed() {
	phoneNumber := "0821abc"
	validPhoneNumber := ValidatePhoneNumber(phoneNumber)

	assert.False(suite.T(), validPhoneNumber)
}

func (suite *ValidationTestSuite) TestValidateEmail_Success() {
	email := "rmalariq@gmail.com"
	validEmail := ValidateEmail(email)

	assert.False(suite.T(), validEmail)
}

func (suite *ValidationTestSuite) TestValidateEmail_FailedPattern() {
	email := "rmalariq@.com"
	validEmail := ValidateEmail(email)

	assert.True(suite.T(), validEmail)
}

func (suite *ValidationTestSuite) TestValidateEmail_FailedCharacter() {
	email := "rmalariq@gmail<>.com"
	validEmail := ValidateEmail(email)

	assert.True(suite.T(), validEmail)
}

func TestRunValidationSuite(t *testing.T) {
	suite.Run(t, new(ValidationTestSuite))
}
