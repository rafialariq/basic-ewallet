package controller

import (
	"errors"
	"final_project_easycash/model"
	"testing"

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

var dummyUsers = []model.User{
	{
		Id:           1,
		Username:     "Dummy Username 1",
		Password:     "Dummy Password",
		Email:        "dummy1@email.com",
		PhoneNumber:  "08111111111",
		PhotoProfile: "-",
		Balance:      100000.00,
	},
	{
		Id:           2,
		Username:     "Dummy Username 2",
		Password:     "Dummy Password 2",
		Email:        "dummy2@email.com",
		PhoneNumber:  "08222222222",
		PhotoProfile: "Dummy Photo Profile",
		Balance:      200000.00,
	},
}

type TransactionUsecaseMock struct {
	mock.Mock
}

type TransactionControllerTestSuite struct {
	suite.Suite
	routerMock             *gin.Engine
	routerGroupMock        *gin.RouterGroup
	transactionUsecaseMock *TransactionUsecaseMock
}

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (t *TransactionUsecaseMock) TransferMoney(sender string, receiver string, amount float64) error {
	args := t.Called(mock.Anything, mock.Anything, mock.AnythingOfType("float64"))
	return args.Error(0)
}

func (suite *TransactionControllerTestSuite) TestTransferMoneyToMerchant_Success() {
	dummyAmount := 10000.00
	transactionUsecaseMock := new(TransactionUsecaseMock)
	suite.transactionUsecaseMock = transactionUsecaseMock
	suite.transactionUsecaseMock.On("TransferMoney", dummyUsers[0].PhoneNumber, dummyMerchants[0].MerchantCode, dummyAmount).Return(nil)

	err := suite.transactionUsecaseMock.TransferMoney(dummyUsers[0].PhoneNumber, dummyMerchants[0].MerchantCode, dummyAmount)
	assert.Nil(suite.T(), err)
}

func (suite *TransactionControllerTestSuite) TestTransferMoneyToMerchant_Failed() {
	dummyAmount := -10000.00
	transactionUsecaseMock := new(TransactionUsecaseMock)
	suite.transactionUsecaseMock = transactionUsecaseMock
	suite.transactionUsecaseMock.On("TransferMoney", dummyUsers[0].PhoneNumber, dummyMerchants[0].MerchantCode, dummyAmount).Return(errors.New("Transfer failed"))

	err := suite.transactionUsecaseMock.TransferMoney(dummyUsers[0].PhoneNumber, dummyMerchants[0].MerchantCode, dummyAmount)
	assert.NotNil(suite.T(), err)
}

func (suite *TransactionControllerTestSuite) SetupSuite() {
	suite.transactionUsecaseMock = new(TransactionUsecaseMock)
}

func TestTransactionControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionControllerTestSuite))
}
