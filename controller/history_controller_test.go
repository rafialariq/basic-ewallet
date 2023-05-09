package controller

import (
	"bytes"
	"encoding/json"
	"final_project_easycash/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyData = []model.Bill{
	{
		Id:                1,
		TransactionId:     "FM001",
		SenderTypeId:      1,
		SenderId:          "082123456789",
		TypeId:            1,
		Date:              time.Date(2022, time.December, 10, 18, 22, 44, 0, time.Local),
		Amount:            80000,
		DestinationTypeId: 1,
		DestinationId:     "085712345678",
		Status:            1,
	},
	{
		Id:                2,
		TransactionId:     "FM002",
		SenderTypeId:      1,
		SenderId:          "082123456789",
		TypeId:            2,
		Amount:            45000,
		Date:              time.Date(2022, time.December, 10, 18, 22, 44, 0, time.Local),
		DestinationTypeId: 2,
		DestinationId:     "7750821758759",
		Status:            1,
	},
	{
		Id:                3,
		TransactionId:     "FM003",
		SenderTypeId:      1,
		SenderId:          "085712345678",
		TypeId:            1,
		Amount:            50000,
		Date:              time.Date(2022, time.December, 10, 18, 22, 44, 0, time.Local),
		DestinationTypeId: 1,
		DestinationId:     "082123456789",
		Status:            1,
	},
}

type historyUsecaseMock struct {
	mock.Mock
}

func (h *historyUsecaseMock) HistoryByUser(user model.User) ([]model.Bill, error) {
	args := h.Called(user)
	return args.Get(0).([]model.Bill), args.Error(1)
}

func (h *historyUsecaseMock) HistoryWithAccountFilter(user model.User, accountTypeId int) ([]model.Bill, error) {
	args := h.Called(user, accountTypeId)
	return args.Get(0).([]model.Bill), args.Error(1)
}

func (h *historyUsecaseMock) HistoryWithTypeFilter(user model.User, typeId int) ([]model.Bill, error) {
	args := h.Called(user, typeId)
	return args.Get(0).([]model.Bill), args.Error(1)
}

func (h *historyUsecaseMock) HistoryWithAmountFilter(user model.User, moreThan, lessThan float64) ([]model.Bill, error) {
	args := h.Called(user, moreThan, lessThan)
	return args.Get(0).([]model.Bill), args.Error(1)
}

type HistoryControllerTestSuite struct {
	suite.Suite
	usecaseMock *historyUsecaseMock
	routerMock  *gin.Engine
}

func (suite *HistoryControllerTestSuite) TestFindAllByUser_Success() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: "082123456789"}
	body, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history", bytes.NewBuffer(body))

	res := dummyData
	suite.usecaseMock.On("HistoryByUser", user).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	h.FindAllByUser(ctx)

	assert.Equal(suite.T(), http.StatusOK, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), res, response)
}

func (suite *HistoryControllerTestSuite) TestFindAllByUser_Failed() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: ""}
	body, _ := json.Marshal(user.Email)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history", bytes.NewBuffer(body))

	res := dummyData
	suite.usecaseMock.On("HistoryByUser", user).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	h.FindAllByUser(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.Error(suite.T(), err)
}

func (suite *HistoryControllerTestSuite) TestFindByAccountType_Success() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: "082123456789"}
	body, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history/account", bytes.NewBuffer(body))

	res := dummyData
	accountTypeId := 1
	suite.usecaseMock.On("HistoryWithAccountFilter", user, accountTypeId).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	ctx.Params = []gin.Param{{Key: "accountTypeId", Value: strconv.Itoa(accountTypeId)}}
	h.FindByAccountType(ctx)

	assert.Equal(suite.T(), http.StatusOK, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), res, response)

	suite.usecaseMock.AssertCalled(suite.T(), "HistoryWithAccountFilter", user, accountTypeId)
}

func (suite *HistoryControllerTestSuite) TestFindByAccountType_FailedBindJson() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: ""}
	body, _ := json.Marshal(user.Email)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history/account", bytes.NewBuffer(body))

	res := dummyData
	accountTypeId := 1
	suite.usecaseMock.On("HistoryWithAccountFilter", user, accountTypeId).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	ctx.Params = []gin.Param{{Key: "accountTypeId", Value: strconv.Itoa(accountTypeId)}}
	h.FindByAccountType(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.Error(suite.T(), err)

}

func (suite *HistoryControllerTestSuite) TestFindByAccountType_FailedGetParam() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: "082123456789"}
	body, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history/account", bytes.NewBuffer(body))

	res := dummyData
	suite.usecaseMock.On("HistoryWithAccountFilter", user).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	h.FindByAccountType(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.Error(suite.T(), err)
}

func (suite *HistoryControllerTestSuite) TestFindByType_Success() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: "082123456789"}
	body, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history/type", bytes.NewBuffer(body))

	res := dummyData
	typeId := 1
	suite.usecaseMock.On("HistoryWithTypeFilter", user, typeId).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	ctx.Params = []gin.Param{{Key: "typeId", Value: strconv.Itoa(typeId)}}
	h.FindByType(ctx)

	assert.Equal(suite.T(), http.StatusOK, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), res, response)

	suite.usecaseMock.AssertCalled(suite.T(), "HistoryWithTypeFilter", user, typeId)
}

func (suite *HistoryControllerTestSuite) TestFindByType_FailedBindJson() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: ""}
	body, _ := json.Marshal(user.Email)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history/type", bytes.NewBuffer(body))

	res := dummyData
	typeId := 1
	suite.usecaseMock.On("HistoryWithTypeFilter", user, typeId).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	ctx.Params = []gin.Param{{Key: "typeId", Value: strconv.Itoa(typeId)}}
	h.FindByType(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.Error(suite.T(), err)
}

func (suite *HistoryControllerTestSuite) TestFindByType_FailedGetParam() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: "082123456789"}
	body, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history/type", bytes.NewBuffer(body))

	res := dummyData
	suite.usecaseMock.On("HistoryWithAccountFilter", user).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	h.FindByType(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.Error(suite.T(), err)
}

func (suite *HistoryControllerTestSuite) TestFindByAmount_Success() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: "082123456789"}
	body, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history", bytes.NewBuffer(body))

	res := dummyData
	var moreThan float64 = 40000
	var lessThan float64 = 60000
	suite.usecaseMock.On("HistoryWithAmountFilter", user, moreThan, lessThan).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	ctx.Params = []gin.Param{{Key: "more_than", Value: strconv.FormatFloat(moreThan, 'f', 0, 64)},
		{Key: "less_than", Value: strconv.FormatFloat(lessThan, 'f', 0, 64)}}
	h.FindByAmount(ctx)

	assert.Equal(suite.T(), http.StatusOK, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), res, response)

	suite.usecaseMock.AssertCalled(suite.T(), "HistoryWithAmountFilter", user, moreThan, lessThan)
}

func (suite *HistoryControllerTestSuite) TestFindByAmount_FailedBindJson() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: ""}
	body, _ := json.Marshal(user.Email)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history", bytes.NewBuffer(body))

	res := dummyData
	var moreThan float64 = 40000
	var lessThan float64 = 60000
	suite.usecaseMock.On("HistoryWithAmountFilter", user, moreThan, lessThan).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	ctx.Params = []gin.Param{{Key: "more_than", Value: strconv.FormatFloat(moreThan, 'f', 0, 64)},
		{Key: "less_than", Value: strconv.FormatFloat(lessThan, 'f', 0, 64)}}
	h.FindByAmount(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.Error(suite.T(), err)
}

func (suite *HistoryControllerTestSuite) TestFindByAmount_FailedGetParamMoreThan() {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: "082123456789"}
	body, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history", bytes.NewBuffer(body))

	res := dummyData
	var lessThan float64 = 60000
	suite.usecaseMock.On("HistoryWithAmountFilter", user, lessThan).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	ctx.Params = []gin.Param{{Key: "less_than", Value: strconv.FormatFloat(lessThan, 'f', 0, 64)}}
	h.FindByAmount(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.Error(suite.T(), err)
}

func (suite *HistoryControllerTestSuite) TestFindByAmount_FailedGetParamLessThan() {
	// membuat context dan user mock
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	user := model.User{PhoneNumber: "082123456789"}
	body, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/history", bytes.NewBuffer(body))

	res := dummyData
	var moreThan float64 = 40000
	suite.usecaseMock.On("HistoryWithAmountFilter", user, moreThan).Return(res, nil)

	h := &HistoryController{suite.usecaseMock}
	ctx.Params = []gin.Param{{Key: "more_than", Value: strconv.FormatFloat(moreThan, 'f', 0, 64)}}
	h.FindByAmount(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, ctx.Writer.Status())
	var response []model.Bill
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	assert.Error(suite.T(), err)
}

func (suite *HistoryControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.usecaseMock = new(historyUsecaseMock)
}

func TestHistoryControllerTestSuite(t *testing.T) {
	suite.Run(t, new(HistoryControllerTestSuite))
}
