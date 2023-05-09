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

type registerUsecaseMock struct {
	mock.Mock
}

func (r *registerUsecaseMock) UserSignup(newUser *model.User) (bool, string) {
	args := r.Called(newUser)
	return args.Bool(0), args.String(1)
}

type RegisterControllerTestSuite struct {
	suite.Suite
	usecaseMock *registerUsecaseMock
	routerMock  *gin.Engine
}

type responsePattern struct {
	Msg      string `json:"message"`
	JwtToken string `json:"toke"`
}

func (suite *RegisterControllerTestSuite) TestRegisterHandler_Success() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	newUser := &model.User{
		Username:    "dummyUser1",
		Password:    "secretPass1",
		Email:       "dummy@gmail.com",
		PhoneNumber: "082123456789",
	}
	body, _ := json.Marshal(newUser)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))

	suite.usecaseMock.On("UserSignup", newUser).Return(true, "")

	r := &RegisterController{suite.usecaseMock}
	r.RegisterHandler(ctx)

	assert.Equal(suite.T(), http.StatusCreated, ctx.Writer.Status())
	var response responsePattern
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "user created successfully", response.Msg)
}

func (suite *RegisterControllerTestSuite) TestRegisterHandler_FailedBindJson() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	newUser := &model.User{}
	body, _ := json.Marshal(newUser.PhotoProfile)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))

	suite.usecaseMock.On("UserSignup").Return(true, "")

	r := &RegisterController{suite.usecaseMock}
	r.RegisterHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
}

func (suite *RegisterControllerTestSuite) TestRegisterHandler_FailedSignup() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	newUser := &model.User{
		Username:    "dummyUser1",
		Password:    "secretPass1",
		Email:       "dummy@gmail.com",
		PhoneNumber: "082123456789",
	}
	body, _ := json.Marshal(newUser)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))

	suite.usecaseMock.On("UserSignup", newUser).Return(false, "")

	r := &RegisterController{suite.usecaseMock}
	r.RegisterHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
}

func (suite *RegisterControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.usecaseMock = new(registerUsecaseMock)
}

func TestRegisterControllerTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterControllerTestSuite))
}
