package usecase

import (
	"testing"

	"final_project_easycash/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyNewUser = []model.User{
	{
		Username:    "userDummy1",
		Email:       "user1@gmail.com",
		PhoneNumber: "081234567891",
		Password:    "passwordUser1",
	},
	{
		Username:    "userDummy2",
		Email:       "user2@gmail.com",
		PhoneNumber: "081234567892",
		Password:    "passwordUser2",
	},
	{
		Username:    "userDummy3",
		Email:       "user3@gmail.com",
		PhoneNumber: "081234567893",
		Password:    "passwordUser3",
	},
}

type registerRepoMock struct {
	mock.Mock
}

func (r *registerRepoMock) UserRegister(newUser *model.User) (bool, string) {
	args := r.Called(newUser)
	return args.Bool(0), args.String(1)
}

func (r *registerRepoMock) RegisterValidate(recUser *model.User) bool {
	args := r.Called(recUser)
	return args.Bool(0)
}

type RegisterUsecaseTestSuite struct {
	repoMock *registerRepoMock
	suite.Suite
}

func (suite *RegisterUsecaseTestSuite) TestUserSignUp_Success() {

	suite.repoMock.On("RegisterValidate", mock.Anything).Return(false)
	suite.repoMock.On("UserRegister", mock.Anything).Return(true, "")

	newUser := &dummyNewUser[0]
	newUser.Password = "secretPass123"

	registerUsecase := NewRegisterService(suite.repoMock)
	res, msg := registerUsecase.UserSignup(newUser)

	assert.True(suite.T(), res)
	assert.NotEmpty(suite.T(), msg)
}

func (suite *RegisterUsecaseTestSuite) TestUserSignUp_UsernameFailed() {
	newUser := dummyNewUser[0]
	newUser.Username = "ab"

	registerUsecase := NewRegisterService(suite.repoMock)
	res, msg := registerUsecase.UserSignup(&newUser)

	assert.False(suite.T(), res)
	assert.Equal(suite.T(), "your username is too short or too long", msg)
}

func (suite *RegisterUsecaseTestSuite) TestUserSignUp_InvalidPassword() {
	newUser := dummyNewUser[0]
	newUser.Password = "ab"

	suite.repoMock.On("RegisterValidate", mock.AnythingOfType("*model.User")).Return(false)

	registerUsecase := NewRegisterService(suite.repoMock)
	res, msg := registerUsecase.UserSignup(&newUser)

	assert.False(suite.T(), res)
	assert.Equal(suite.T(), "invalid password", msg)
}

func (suite *RegisterUsecaseTestSuite) TestUserSignUp_InvalidEmail() {
	newUser := dummyNewUser[0]
	newUser.Password = "secretPass123"
	newUser.Email = "dummy[]@com"

	registerUsecase := NewRegisterService(suite.repoMock)
	res, msg := registerUsecase.UserSignup(&newUser)

	assert.False(suite.T(), res)
	assert.Equal(suite.T(), "invalid email", msg)
}

func (suite *RegisterUsecaseTestSuite) TestUserSignUp_InvalidPhoneNumber() {
	newUser := dummyNewUser[0]
	newUser.Password = "secretPass123"
	newUser.PhoneNumber = "087812"

	registerUsecase := NewRegisterService(suite.repoMock)
	res, msg := registerUsecase.UserSignup(&newUser)

	assert.False(suite.T(), res)
	assert.Equal(suite.T(), "invalid phone number", msg)
}

func (suite *RegisterUsecaseTestSuite) TestUserSignUp_UserAlreadyExist() {
	newUser := dummyNewUser[0]
	newUser.Password = "secretPass123"

	suite.repoMock.On("RegisterValidate", mock.AnythingOfType("*model.User")).Return(true)

	registerUsecase := NewRegisterService(suite.repoMock)
	res, msg := registerUsecase.UserSignup(&newUser)

	assert.False(suite.T(), res)
	assert.Equal(suite.T(), "user already exist", msg)
}

// func (suite *RegisterUsecaseTestSuite) TestUserSignUp_FailedToGenerateToken() {
// 	user := &model.User{
// 		Username:    "user123",
// 		Email:       "user123@example.com",
// 		PhoneNumber: "081234567890",
// 		Password:    "password",
// 	}

// 	suite.repoMock.On("RegisterValidate", mock.Anything).Return(false)
// 	suite.repoMock.On("UserRegister", mock.Anything).Return(true, "")

// 	registerUsecase := NewRegisterService(suite.repoMock)

// 	suite..secretKey = []byte("invalid")

// 	res, msg := registerUsecase.UserSignup(user)

// 	assert.False(suite.T(), res)
// 	assert.Equal(suite.T(), "failed to generate", msg)
// }

func (suite *RegisterUsecaseTestSuite) TestUserSignUp_Failed() {
	newUser := &dummyNewUser[0]

	suite.repoMock.On("RegisterValidate", mock.AnythingOfType("*model.User")).Return(false)
	suite.repoMock.On("UserRegister", mock.AnythingOfType("*model.User")).Return(false, "failed to create user")

	registerUsecase := NewRegisterService(suite.repoMock)
	res, msg := registerUsecase.UserSignup(newUser)

	assert.False(suite.T(), res)
	assert.Equal(suite.T(), "failed to create user", msg)

}

func (suite *RegisterUsecaseTestSuite) SetupTest() {
	suite.repoMock = new(registerRepoMock)
}

func TestRegisterUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterUsecaseTestSuite))
}
