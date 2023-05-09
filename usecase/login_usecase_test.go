package usecase

import (
	"log"
	"testing"
	"time"

	"final_project_easycash/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyUser = []model.User{
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

type loginRepoMock struct {
	mock.Mock
}

func (l *loginRepoMock) FindUser(recUser model.User) (bool, string) {
	args := l.Called(recUser)
	return args.Bool(0), args.String(1)
}

type LoginUsecaseTestSuite struct {
	repoMock *loginRepoMock
	suite.Suite
}

func (suite *LoginUsecaseTestSuite) TestUserLogin_Success() {

	suite.repoMock.On("FindUser", dummyUser[0]).Return(true, "successfully login")

	expectedToken := jwt.New(jwt.SigningMethodHS256)
	claims := expectedToken.Claims.(jwt.MapClaims)
	claims["username"] = dummyUser[0].Username
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	expectedTokenString, err := expectedToken.SignedString([]byte("secretkey"))
	if err != nil {
		log.Println(err)
	}

	loginUsecase := NewLoginService(suite.repoMock)
	success, token := loginUsecase.UserLogin(dummyUser[0])
	assert.True(suite.T(), success)
	assert.Equal(suite.T(), expectedTokenString, token)

}

func (suite *LoginUsecaseTestSuite) TestUserLogin_Failed() {

	suite.repoMock.On("FindUser", dummyUser[0]).Return(false, "invalid password")

	loginUsecase := NewLoginService(suite.repoMock)
	success, res := loginUsecase.UserLogin(dummyUser[0])
	assert.False(suite.T(), success)
	assert.Equal(suite.T(), "invalid password", res)
}

func (suite *LoginUsecaseTestSuite) TestUserLogin_FailedGenerateToken() {

	suite.repoMock.On("FindUser", dummyUser[0]).Return(false, "failed to generate token")

	loginUsecase := NewLoginService(suite.repoMock)
	success, res := loginUsecase.UserLogin(dummyUser[0])
	assert.False(suite.T(), success)
	assert.Equal(suite.T(), "failed to generate token", res)
}

func (suite *LoginUsecaseTestSuite) SetupTest() {
	suite.repoMock = new(loginRepoMock)
}

func TestLoginUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(LoginUsecaseTestSuite))
}
