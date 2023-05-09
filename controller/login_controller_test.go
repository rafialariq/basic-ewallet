package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"final_project_easycash/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type loginUsecaseMock struct {
	mock.Mock
}

func (l *loginUsecaseMock) UserLogin(user model.User) (bool, string) {
	args := l.Called(user)
	return args.Bool(0), args.String(1)
}

type LoginControllerTestSuite struct {
	suite.Suite
	usecaseMock *loginUsecaseMock
	routerMock  *gin.Engine
}

func (suite *LoginControllerTestSuite) TestLoginHandler_Success() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{
		Username: "dummyUser1",
		Password: "secretPass1",
	}
	body, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))

	suite.usecaseMock.On("UserLogin", user).Return(true, "")

	l := &LoginController{suite.usecaseMock}
	l.LoginHandler(ctx)

	assert.Equal(suite.T(), http.StatusOK, ctx.Writer.Status())
}

func (suite *LoginControllerTestSuite) TestLoginHandler_FailedBindJson() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{}
	body, _ := json.Marshal(user.PhotoProfile)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))

	suite.usecaseMock.On("UserLogin").Return(true, "")

	l := &LoginController{suite.usecaseMock}
	l.LoginHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
}

func (suite *LoginControllerTestSuite) TestLoginHandler_FailedUnauthorized() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{
		Username: "dummyUser1",
		Password: "secretPass1",
	}
	body, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))

	suite.usecaseMock.On("UserLogin", user).Return(false, "")

	l := &LoginController{suite.usecaseMock}
	l.LoginHandler(ctx)

	assert.Equal(suite.T(), http.StatusUnauthorized, ctx.Writer.Status())
}

func (suite *LoginControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.usecaseMock = new(loginUsecaseMock)
}

func TestLoginControllerTestSuite(t *testing.T) {
	suite.Run(t, new(LoginControllerTestSuite))
}
