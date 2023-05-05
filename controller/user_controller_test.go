package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"final_project_easycash/model"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyMerchants = []model.Merchant{
	{
		Id:           1,
		MerchantCode: "Dummy Merchant Code 1",
		Name:         "Dummy Merchant Name 1",
		Amount:       20000.00,
	},
}

var dummyBanks = []model.Bank{
	{
		Id:         1,
		BankNumber: "Dummy Bank Code 1",
		Name:       "Dummy Bank Name 1",
	},
}
var dummyUsers = []model.User{
	{
		Id:           1,
		Username:     "DummyUsername1",
		Password:     "Dummy Password",
		Email:        "dummy1@email.com",
		PhoneNumber:  "081111111111",
		PhotoProfile: "Dummy Photo Profile 1",
		Balance:      100000.00,
	},
	{
		Id:           2,
		Username:     "Dummy Usename 2",
		Password:     "Dummy Password 2",
		Email:        "dummy2@email.com",
		PhoneNumber:  "082222222222",
		PhotoProfile: "Dummy Photo Profile 2",
		Balance:      200000.00,
	},
}

type UserUsecaseMock struct {
	mock.Mock
}

type UserControllerTestSuite struct {
	suite.Suite
	routerMock      *gin.Engine
	routerGroupMock *gin.RouterGroup
	usecaseMock     *UserUsecaseMock
}

func (u *UserUsecaseMock) CheckProfile(username string) (model.User, error) {
	args := u.Called(username)
	if args.Get(0) == nil {
		return model.User{}, args.Error(1)
	}
	return args.Get(0).(model.User), args.Error(1)
}

func (u *UserUsecaseMock) EditProfile(updatedUserData *model.User) error {
	args := u.Called(updatedUserData)
	if err := args.Error(0); err != nil {
		return err
	}
	return nil
}

func (u *UserUsecaseMock) EditPhotoProfile(username string, fileExt string, file *multipart.File) error {
	args := u.Called(username, fileExt, file)
	if args == nil {
		return errors.New("Failed")
	}
	return nil
}
func (u *UserUsecaseMock) UnregProfile(username string) error {
	args := u.Called(username)
	if args == nil {
		return errors.New("Failed")
	}
	return nil
}

func (suite *UserControllerTestSuite) TestCheckProfile_Success() {
	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/menu/profile/%s", dummyUsers[0].Username), nil)
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	// Mock the CheckProfile method to return a dummy user
	suite.usecaseMock.On("CheckProfile", "DummyUsername1").Return(dummyUsers[0], nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	ginContext.Params = []gin.Param{{Key: "username", Value: dummyUsers[0].Username}}
	userController.CheckProfile(ginContext)

	// Check the status code and response body
	var actual model.User
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusOK, responseWriter.Code)
	assert.Equal(suite.T(), dummyUsers[0], actual)
}

func (suite *UserControllerTestSuite) TestCheckProfileMissingClaims_Failed() {
	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/menu/profile/%s", dummyUsers[0].Username), nil)
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	// Mock the CheckProfile method to return a dummy user
	suite.usecaseMock.On("CheckProfile", "DummyUsername1").Return(dummyUsers[0], nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	userController.CheckProfile(ginContext)

	// Check the status code and response body
	var actual model.User
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.Equal(suite.T(), model.User{}, actual)
}

func (suite *UserControllerTestSuite) TestCheckProfileMissingUsername_Failed() {
	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/menu/profile/%s", dummyUsers[0].Username), nil)
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	// Mock the CheckProfile method to return a dummy user
	suite.usecaseMock.On("CheckProfile", dummyUsers[1].Username).Return(dummyUsers[0], nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{})
	ginContext.Params = []gin.Param{{Key: "username", Value: dummyUsers[1].Username}}
	userController.CheckProfile(ginContext)

	// Check the status code and response body
	var actual model.User
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.Equal(suite.T(), model.User{}, actual)
}

func (suite *UserControllerTestSuite) TestCheckProfileMissmatchedUsername_Failed() {
	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/menu/profile/%s", dummyUsers[0].Username), nil)
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	// Mock the CheckProfile method to return a dummy user
	suite.usecaseMock.On("CheckProfile", dummyUsers[1].Username).Return(dummyUsers[0], nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	ginContext.Params = []gin.Param{{Key: "username", Value: dummyUsers[1].Username}}
	userController.CheckProfile(ginContext)

	// Check the status code and response body
	var actual model.User
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.Equal(suite.T(), model.User{}, actual)
}

func (suite *UserControllerTestSuite) TestEditProfile_Success() {
	updatedUserData := dummyUsers[0]
	jsonData, _ := json.Marshal(updatedUserData)

	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodPost, "/menu/profile/edit", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	suite.usecaseMock.On("EditProfile", &dummyUsers[0]).Return(nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	userController.EditProfile(ginContext)

	// Check the status code and response body
	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusOK, responseWriter.Code)
	assert.Equal(suite.T(), "", actual.Error)
}

func (suite *UserControllerTestSuite) TestEditProfileMissingClaims_Failed() {
	updatedUserData := dummyUsers[0]
	jsonData, _ := json.Marshal(updatedUserData)

	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodPost, "/menu/profile/edit", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	suite.usecaseMock.On("EditProfile", &dummyUsers[0]).Return(nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	userController.EditProfile(ginContext)

	// Check the status code and response body
	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *UserControllerTestSuite) TestEditProfileEmptyBody_Failed() {
	updatedUserData := "dummyUsers[0]"
	jsonData, _ := json.Marshal(updatedUserData)

	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodPost, "/menu/profile/edit", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	suite.usecaseMock.On("EditProfile", &dummyUsers[0]).Return(nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	userController.EditProfile(ginContext)

	// Check the status code and response body
	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusBadRequest, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *UserControllerTestSuite) TestEditProfileMissingUsername_Failed() {
	updatedUserData := dummyUsers[0]
	jsonData, _ := json.Marshal(updatedUserData)

	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodPost, "/menu/profile/edit", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()
	suite.usecaseMock.On("EditProfile", &dummyUsers[0]).Return(nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{})
	userController.EditProfile(ginContext)

	// Check the status code and response body
	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *UserControllerTestSuite) TestEditProfileMismatchedUsername_Failed() {
	updatedUserData := dummyUsers[1]
	jsonData, _ := json.Marshal(updatedUserData)

	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodPost, "/menu/profile/edit", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()
	suite.usecaseMock.On("EditProfile", &dummyUsers[0]).Return(nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	userController.EditProfile(ginContext)

	// Check the status code and response body
	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *UserControllerTestSuite) TestUnregProfile_Success() {
	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/menu/profile/%s", dummyUsers[0].Username), nil)
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	// Mock the CheckProfile method to return a dummy user
	suite.usecaseMock.On("UnregProfile", dummyUsers[0].Username).Return(nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	ginContext.Params = []gin.Param{{Key: "username", Value: dummyUsers[0].Username}}
	userController.UnregProfile(ginContext)

	// Check the status code and response body
	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusOK, responseWriter.Code)
	assert.Equal(suite.T(), "", actual.Error)
}

func (suite *UserControllerTestSuite) TestUnregProfileMissingClaims_Failed() {
	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/menu/profile/%s", dummyUsers[0].Username), nil)
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	// Mock the CheckProfile method to return a dummy user
	suite.usecaseMock.On("UnregProfile", dummyUsers[0].Username).Return(nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Params = []gin.Param{{Key: "username", Value: dummyUsers[0].Username}}
	userController.UnregProfile(ginContext)

	// Check the status code and response body
	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *UserControllerTestSuite) TestUnregProfileMissingUsername_Failed() {
	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/menu/profile/%s", dummyUsers[0].Username), nil)
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	// Mock the CheckProfile method to return a dummy user
	suite.usecaseMock.On("UnregProfile", dummyUsers[0].Username).Return(nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{})
	ginContext.Params = []gin.Param{{Key: "username", Value: dummyUsers[0].Username}}
	userController.UnregProfile(ginContext)

	// Check the status code and response body
	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *UserControllerTestSuite) TestUnregProfileMismatchedUsername_Failed() {
	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/menu/profile/%s", dummyUsers[0].Username), nil)
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	// Mock the CheckProfile method to return a dummy user
	suite.usecaseMock.On("UnregProfile", dummyUsers[1].Username).Return(nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	ginContext.Params = []gin.Param{{Key: "username", Value: dummyUsers[1].Username}}
	userController.UnregProfile(ginContext)

	// Check the status code and response body
	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *UserControllerTestSuite) TestEditPhotoProfile_Success() {
	// Create a new user controller and router
	userController := NewUserController(suite.routerGroupMock, suite.usecaseMock)

	// Create a new HTTP request with the token in the header
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/menu/profile/photo/%s", dummyUsers[0].Username), nil)
	suite.Require().NoError(err)

	// Create a new recorder to capture the response
	responseWriter := httptest.NewRecorder()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := writer.CreateFormFile("photo", "photo.jpg")
	suite.Require().NoError(err)

	fileContent := []byte("file content")
	_, err = file.Write(fileContent)
	suite.Require().NoError(err)

	writer.Close()

	request.Body = ioutil.NopCloser(body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.usecaseMock.On("EditPhotoProfile", dummyUsers[0].Username, "jpg", mock.AnythingOfType("*multipart.File")).Return(nil)

	// Call the handler function with the context
	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	ginContext.Params = []gin.Param{{Key: "username", Value: dummyUsers[0].Username}}
	userController.EditPhotoProfile(ginContext)

	// Check the status code and response body
	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusOK, responseWriter.Code)
	assert.Equal(suite.T(), "", actual.Error)
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.routerGroupMock = suite.routerMock.Group("/menu")
	suite.usecaseMock = new(UserUsecaseMock)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
