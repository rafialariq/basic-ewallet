package usecase

import (
	"errors"
	"final_project_easycash/model"
	"testing"

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
		Username:     "Dummy Username 1",
		Password:     "DummyPassword",
		Email:        "dummy1@email.com",
		PhoneNumber:  "08111111111",
		PhotoProfile: "-",
		Balance:      100000.00,
	},
	{
		Id:           2,
		Username:     "DummyUsername 2",
		Password:     "DummyPassword2",
		Email:        "dummy2@email.com",
		PhoneNumber:  "08222222222",
		PhotoProfile: "Dummy Photo Profile",
		Balance:      200000.00,
	},
}

type transRepoMock struct {
	mock.Mock
}

type TransactionUsecaseTestSuite struct {
	repoMock *transRepoMock
	suite.Suite
}

func (t *transRepoMock) TransferMoney(sender string, receiver string, amount float64) error {
	args := t.Called(sender, receiver, amount)
	if args == nil {
		return errors.New("Failed")
	}
	return nil
}

func (t *transRepoMock) TopUpBalance(sender string, receiver string, amount float64) error {
	args := t.Called(sender, receiver, amount)
	if args == nil {
		return errors.New("Failed")
	}
	return nil
}

func (t *transRepoMock) WithdrawBalance(sender string, receiver string, amount float64) error {
	args := t.Called(sender, receiver, amount)
	if args == nil {
		return errors.New("Failed")
	}
	return nil
}

func (t *transRepoMock) TransferBalance(sender string, receiver string, amount float64) error {
	args := t.Called(sender, receiver, amount)
	if args == nil {
		return errors.New("Failed")
	}
	return nil
}

func (t *transRepoMock) SplitBill(sender string, receiver []string, amount []float64) error {
	args := t.Called(sender, receiver, amount)
	if args == nil {
		return errors.New("Failed")
	}
	return nil
}

func (t *transRepoMock) PayBill(receiver string, id_transaction string) error {
	args := t.Called(receiver, id_transaction)
	if args == nil {
		return errors.New("Failed")
	}
	return nil
}

func (suite *TransactionUsecaseTestSuite) TestTopUpBalance_Success() {
	dummyAmount := 20000.00
	dummyAmountAfterAdmin := 19000.00
	transactionUsecase := NewTransactionUsecase(suite.repoMock)
	suite.repoMock.On("TopUpBalance", dummyBanks[0].BankNumber, dummyUsers[0].PhoneNumber, dummyAmountAfterAdmin).Return(nil)
	err := transactionUsecase.TopUpBalance(dummyBanks[0].BankNumber, dummyUsers[0].PhoneNumber, dummyAmount)
	assert.Nil(suite.T(), err)
}

func (suite *TransactionUsecaseTestSuite) TestTopUpBalance_Failed() {
	dummyAmount := -20000.00
	dummyAmountAfterAdmin := 19000.00
	transactionUsecase := NewTransactionUsecase(suite.repoMock)
	suite.repoMock.On("TopUpBalance", dummyBanks[0].BankNumber, dummyUsers[0].PhoneNumber, dummyAmountAfterAdmin).Return(nil)
	err := transactionUsecase.TopUpBalance(dummyBanks[0].BankNumber, dummyUsers[0].PhoneNumber, dummyAmount)
	assert.NotNil(suite.T(), err)
}

func (suite *TransactionUsecaseTestSuite) TestWithdrawBalance_Success() {
	dummyAmount := 20000.00
	dummyAmountAfterAdmin := 22500.00
	transactionUsecase := NewTransactionUsecase(suite.repoMock)
	suite.repoMock.On("WithdrawBalance", dummyUsers[0].PhoneNumber, dummyBanks[0].BankNumber, dummyAmountAfterAdmin).Return(nil)
	err := transactionUsecase.WithdrawBalance(dummyUsers[0].PhoneNumber, dummyBanks[0].BankNumber, dummyAmount)
	assert.Nil(suite.T(), err)
}

func (suite *TransactionUsecaseTestSuite) TestWithdrawBalance_Failed() {
	dummyAmount := -20000.00
	dummyAmountAfterAdmin := 22500
	transactionUsecase := NewTransactionUsecase(suite.repoMock)
	suite.repoMock.On("WithdrawBalance", dummyBanks[0].BankNumber, dummyUsers[0].PhoneNumber, dummyAmountAfterAdmin).Return(nil)
	err := transactionUsecase.WithdrawBalance(dummyBanks[0].BankNumber, dummyUsers[0].PhoneNumber, dummyAmount)
	assert.NotNil(suite.T(), err)
}

func (suite *TransactionUsecaseTestSuite) TestTransferBalance_Success() {
	dummyAmount := 20000.00
	transactionUsecase := NewTransactionUsecase(suite.repoMock)
	suite.repoMock.On("TransferBalance", dummyUsers[0].PhoneNumber, dummyUsers[1].PhoneNumber, dummyAmount).Return(nil)
	err := transactionUsecase.TransferBalance(dummyUsers[0].PhoneNumber, dummyUsers[1].PhoneNumber, dummyAmount)
	assert.Nil(suite.T(), err)
}

func (suite *TransactionUsecaseTestSuite) TestTransferBalance_Failed() {
	dummyAmount := -20000.00
	transactionUsecase := NewTransactionUsecase(suite.repoMock)
	suite.repoMock.On("TransferBalance", dummyUsers[0].PhoneNumber, dummyUsers[1].PhoneNumber, dummyAmount).Return(nil)
	err := transactionUsecase.TransferBalance(dummyUsers[0].PhoneNumber, dummyUsers[1].PhoneNumber, dummyAmount)
	assert.NotNil(suite.T(), err)
}

func (suite *TransactionUsecaseTestSuite) SetupTest() {
	suite.repoMock = new(transRepoMock)
}

func TestTransactionUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionUsecaseTestSuite))
}
