package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"final_project_easycash/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TransactionUsecaseMock struct {
	mock.Mock
}

type TransactionControllerTestSuite struct {
	suite.Suite
	routerMock             *gin.Engine
	routerGroupMock        *gin.RouterGroup
	transactionUsecaseMock *TransactionUsecaseMock
	userUsecaseMock        *UserUsecaseMock
}

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (u *TransactionUsecaseMock) TransferMoney(sender string, receiver string, amount float64) error {
	args := u.Called(sender, receiver, amount)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return nil
}

func (u *TransactionUsecaseMock) TopUpBalance(sender string, receiver string, amount float64) error {
	args := u.Called(sender, receiver, amount)
	if err := args.Error(0); err != nil {
		return err
	}
	return nil
}

func (u *TransactionUsecaseMock) WithdrawBalance(sender string, receiver string, amount float64) error {
	args := u.Called(sender, receiver, amount)
	if err := args.Error(0); err != nil {
		return err
	}
	return nil
}

func (u *TransactionUsecaseMock) TransferBalance(sender string, receiver string, amount float64) error {
	args := u.Called(sender, receiver, amount)
	if err := args.Error(0); err != nil {
		return err
	}
	return nil
}

func (u *TransactionUsecaseMock) SplitBill(sender string, receiver []string, amount []float64) error {
	args := u.Called(sender, receiver, amount)
	if err := args.Error(0); err != nil {
		return err
	}
	return nil
}

func (u *TransactionUsecaseMock) PayBill(receiver string, id_transaction string) error {
	args := u.Called(receiver, id_transaction)
	if err := args.Error(0); err != nil {
		return err
	}
	return nil
}

func (suite *TransactionControllerTestSuite) TestTopUpBalance_Success() {
	NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	var topUpDummy model.Bill
	topUpDummy.SenderId = dummyBanks[0].BankNumber
	topUpDummy.DestinationId = dummyUsers[0].PhoneNumber
	topUpDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(topUpDummy)
	request, err := http.NewRequest(http.MethodPost, "/menu/topup", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)

	responseWriter := httptest.NewRecorder()

	suite.transactionUsecaseMock.On("TopUpBalance", topUpDummy.SenderId, topUpDummy.DestinationId, topUpDummy.Amount).Return(nil)
	suite.routerMock.ServeHTTP(responseWriter, request)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusOK, responseWriter.Code)
	assert.Equal(suite.T(), "", actual.Error)
}

func (suite *TransactionControllerTestSuite) TestTopUpBalanceInvalidJSON_Failed() {
	NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/topup", nil)
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()
	suite.routerMock.ServeHTTP(responseWriter, request)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusBadRequest, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestTopUpInvalidNumber_Failed() {
	NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	var topUpDummy model.Bill
	topUpDummy.SenderId = dummyBanks[0].BankNumber
	topUpDummy.DestinationId = dummyUsers[0].PhoneNumber
	topUpDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(topUpDummy)
	request, err := http.NewRequest(http.MethodPost, "/menu/topup", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)

	responseWriter := httptest.NewRecorder()

	suite.transactionUsecaseMock.On("TopUpBalance", topUpDummy.SenderId, topUpDummy.DestinationId, topUpDummy.Amount).Return(errors.New("Receiver number not found"))
	suite.routerMock.ServeHTTP(responseWriter, request)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusBadRequest, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestTopUpError_Failed() {
	NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	var topUpDummy model.Bill
	topUpDummy.SenderId = dummyBanks[0].BankNumber
	topUpDummy.DestinationId = dummyUsers[0].PhoneNumber
	topUpDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(topUpDummy)
	request, err := http.NewRequest(http.MethodPost, "/menu/topup", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)

	responseWriter := httptest.NewRecorder()

	suite.transactionUsecaseMock.On("TopUpBalance", topUpDummy.SenderId, topUpDummy.DestinationId, topUpDummy.Amount).Return(errors.New("Failed"))
	suite.routerMock.ServeHTTP(responseWriter, request)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusInternalServerError, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestWithdrawBalance_Success() {
	var withdrawDummy model.Bill
	withdrawDummy.SenderId = dummyUsers[0].PhoneNumber
	withdrawDummy.DestinationId = dummyBanks[0].BankNumber
	withdrawDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(withdrawDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/bank", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()
	suite.transactionUsecaseMock.On("WithdrawBalance", withdrawDummy.SenderId, withdrawDummy.DestinationId, withdrawDummy.Amount).Return(nil)
	suite.userUsecaseMock.On("CheckProfile", dummyUsers[0].Username).Return(dummyUsers[0], nil)

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	transactionController.WithdrawBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusOK, responseWriter.Code)
	assert.Equal(suite.T(), "", actual.Error)
}

func (suite *TransactionControllerTestSuite) TestWithdrawMissingClaims_Failed() {
	var withdrawDummy model.Bill
	withdrawDummy.SenderId = dummyUsers[0].PhoneNumber
	withdrawDummy.DestinationId = dummyBanks[0].BankNumber
	withdrawDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(withdrawDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/bank", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{})
	transactionController.WithdrawBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestWithdrawMissingUsername_Failed() {
	var withdrawDummy model.Bill
	withdrawDummy.SenderId = dummyUsers[0].PhoneNumber
	withdrawDummy.DestinationId = dummyBanks[0].BankNumber
	withdrawDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(withdrawDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/bank", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	transactionController.WithdrawBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestWithdrawMismatchedUsername_Failed() {
	var withdrawDummy model.Bill
	withdrawDummy.SenderId = dummyUsers[1].PhoneNumber
	withdrawDummy.DestinationId = dummyBanks[0].BankNumber
	withdrawDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(withdrawDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/bank", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()
	suite.userUsecaseMock.On("CheckProfile", dummyUsers[0].Username).Return(dummyUsers[0], nil)

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	transactionController.WithdrawBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestWithdrawErrorCheckProfile_Failed() {
	var withdrawDummy model.Bill
	withdrawDummy.SenderId = dummyUsers[1].PhoneNumber
	withdrawDummy.DestinationId = dummyBanks[0].BankNumber
	withdrawDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(withdrawDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/bank", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()
	suite.userUsecaseMock.On("CheckProfile", dummyUsers[0].Username).Return(model.User{}, errors.New("Failed"))

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	transactionController.WithdrawBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusBadRequest, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestWithdrawBalanceNumberNotFound_Failed() {
	var withdrawDummy model.Bill
	withdrawDummy.SenderId = dummyUsers[0].PhoneNumber
	withdrawDummy.DestinationId = dummyBanks[0].BankNumber
	withdrawDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(withdrawDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/bank", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()
	suite.transactionUsecaseMock.On("WithdrawBalance", withdrawDummy.SenderId, withdrawDummy.DestinationId, withdrawDummy.Amount).Return(errors.New("Receiver number not found"))
	suite.userUsecaseMock.On("CheckProfile", dummyUsers[0].Username).Return(dummyUsers[0], nil)

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	transactionController.WithdrawBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusBadRequest, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestWithdrawBalanceErrorUsecase_Failed() {
	var withdrawDummy model.Bill
	withdrawDummy.SenderId = dummyUsers[0].PhoneNumber
	withdrawDummy.DestinationId = dummyBanks[0].BankNumber
	withdrawDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(withdrawDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/bank", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()
	suite.transactionUsecaseMock.On("WithdrawBalance", withdrawDummy.SenderId, withdrawDummy.DestinationId, withdrawDummy.Amount).Return(errors.New("Failed"))
	suite.userUsecaseMock.On("CheckProfile", dummyUsers[0].Username).Return(dummyUsers[0], nil)

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	transactionController.WithdrawBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusInternalServerError, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestTransferBalance_Success() {
	var transferDummy model.Bill
	transferDummy.SenderId = dummyUsers[0].PhoneNumber
	transferDummy.DestinationId = dummyUsers[1].PhoneNumber
	transferDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(transferDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/user", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()
	suite.transactionUsecaseMock.On("TransferBalance", transferDummy.SenderId, transferDummy.DestinationId, transferDummy.Amount).Return(nil)
	suite.userUsecaseMock.On("CheckProfile", dummyUsers[0].Username).Return(dummyUsers[0], nil)

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	transactionController.TransferBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusOK, responseWriter.Code)
	assert.Equal(suite.T(), "", actual.Error)
}

func (suite *TransactionControllerTestSuite) TestTransferBalanceMissingClaims_Failed() {

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/user", nil)
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	transactionController.TransferBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestTransferBalanceMissingUsername_Failed() {
	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/user", nil)
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{})
	transactionController.TransferBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestTransferBalanceMismatchedUsername_Failed() {
	var transferDummy model.Bill
	transferDummy.SenderId = dummyUsers[1].PhoneNumber
	transferDummy.DestinationId = dummyUsers[0].PhoneNumber
	transferDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(transferDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/user", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	suite.userUsecaseMock.On("CheckProfile", dummyUsers[0].Username).Return(dummyUsers[0], nil)
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	transactionController.TransferBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestTransferBalanceErrorCheckProfile_Failed() {
	var transferDummy model.Bill
	transferDummy.SenderId = dummyUsers[0].PhoneNumber
	transferDummy.DestinationId = dummyUsers[1].PhoneNumber
	transferDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(transferDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/user", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	suite.userUsecaseMock.On("CheckProfile", dummyUsers[0].Username).Return(model.User{}, errors.New("Failed"))
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	transactionController.TransferBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusBadRequest, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestTransferBalanceNumberNotFound_Failed() {
	var transferDummy model.Bill
	transferDummy.SenderId = dummyUsers[0].PhoneNumber
	transferDummy.DestinationId = dummyUsers[1].PhoneNumber
	transferDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(transferDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/user", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()
	suite.transactionUsecaseMock.On("TransferBalance", transferDummy.SenderId, transferDummy.DestinationId, transferDummy.Amount).Return(errors.New("Receiver number not found"))
	suite.userUsecaseMock.On("CheckProfile", dummyUsers[0].Username).Return(dummyUsers[0], nil)

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	transactionController.TransferBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusBadRequest, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) TestTransferBalanceErrorUsecase_Failed() {
	var transferDummy model.Bill
	transferDummy.SenderId = dummyUsers[0].PhoneNumber
	transferDummy.DestinationId = dummyUsers[1].PhoneNumber
	transferDummy.Amount = 10000.00
	jsonData, _ := json.Marshal(transferDummy)

	transactionController := NewTransactionController(suite.routerGroupMock, suite.transactionUsecaseMock, suite.userUsecaseMock)
	request, err := http.NewRequest(http.MethodPost, "/menu/transfer/user", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	responseWriter := httptest.NewRecorder()
	suite.transactionUsecaseMock.On("TransferBalance", transferDummy.SenderId, transferDummy.DestinationId, transferDummy.Amount).Return(errors.New("Failed"))
	suite.userUsecaseMock.On("CheckProfile", dummyUsers[0].Username).Return(dummyUsers[0], nil)

	ginContext, _ := gin.CreateTestContext(responseWriter)
	ginContext.Request = request
	ginContext.Set("claims", jwt.MapClaims{"username": dummyUsers[0].Username})
	transactionController.TransferBalance(ginContext)

	var actual Response
	response := responseWriter.Body.String()
	json.Unmarshal([]byte(response), &actual)

	assert.Equal(suite.T(), http.StatusInternalServerError, responseWriter.Code)
	assert.NotNil(suite.T(), actual.Error)
}

func (suite *TransactionControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.routerGroupMock = suite.routerMock.Group("/menu")
	suite.transactionUsecaseMock = new(TransactionUsecaseMock)
	suite.userUsecaseMock = new(UserUsecaseMock)
}

func TestTransactionControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionControllerTestSuite))
}
