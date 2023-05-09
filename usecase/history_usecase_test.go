package usecase

import (
	"final_project_easycash/model"
	"testing"
	"time"

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
		Date:              time.Now(),
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
		Date:              time.Now(),
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
		Date:              time.Now(),
		DestinationTypeId: 1,
		DestinationId:     "082123456789",
		Status:            1,
	},
}

type HistoryRepoMock struct {
	mock.Mock
}

func (h *HistoryRepoMock) GetHistoryByUser(user model.User) ([]model.Bill, error) {
	args := h.Called(&user)
	return args.Get(0).([]model.Bill), args.Error(1)
}

func (h *HistoryRepoMock) GetHistoryWithAccountFilter(user model.User, accountTypeId int) ([]model.Bill, error) {
	args := h.Called(&user, &accountTypeId)
	return args.Get(0).([]model.Bill), args.Error(1)
}

func (h *HistoryRepoMock) GetHistoryWithTypeFilter(user model.User, typeId int) ([]model.Bill, error) {
	args := h.Called(&user, &typeId)
	return args.Get(0).([]model.Bill), args.Error(1)
}

func (h *HistoryRepoMock) GetHistoryWithAmountFilter(user model.User, moreThan, lessThan float64) ([]model.Bill, error) {
	args := h.Called(&user, &moreThan, &lessThan)
	return args.Get(0).([]model.Bill), args.Error(1)
}

type HistoryUsecaseTestSuite struct {
	repoMock *HistoryRepoMock
	suite.Suite
}

func (suite *HistoryUsecaseTestSuite) TestHistoryByUser_Success() {
	user := &model.User{PhoneNumber: "082123456789"}

	suite.repoMock.On("GetHistoryByUser", user).Return(dummyData, nil)

	historyUsecase := NewHistoryUsecase(suite.repoMock)
	historyList, err := historyUsecase.HistoryByUser(*user)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), historyList, 3)
	assert.Equal(suite.T(), dummyData, historyList)

}

func (suite *HistoryUsecaseTestSuite) TestHistoryWithAccountFilter_Success() {
	user := &model.User{PhoneNumber: "082123456789"}
	accountTypeId := 1

	suite.repoMock.On("GetHistoryWithAccountFilter", user, &accountTypeId).Return(dummyData, nil)

	historyUsecase := NewHistoryUsecase(suite.repoMock)
	historyList, err := historyUsecase.HistoryWithAccountFilter(*user, accountTypeId)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), historyList, 3)
	assert.Equal(suite.T(), dummyData, historyList)
}

func (suite *HistoryUsecaseTestSuite) TestHistoryWithTypeFilter_Success() {
	user := &model.User{PhoneNumber: "082123456789"}
	typeId := 1

	suite.repoMock.On("GetHistoryWithTypeFilter", user, &typeId).Return(dummyData, nil)

	historyUsecase := NewHistoryUsecase(suite.repoMock)
	historyList, err := historyUsecase.HistoryWithTypeFilter(*user, typeId)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), historyList, 3)
	assert.Equal(suite.T(), dummyData, historyList)
}

func (suite *HistoryUsecaseTestSuite) TestHistoryWithAmountFilter_Success() {
	user := &model.User{PhoneNumber: "082123456789"}
	var moreThan float64 = 50000
	var lessThan float64 = 100000

	suite.repoMock.On("GetHistoryWithAmountFilter", user, &moreThan, &lessThan).Return(dummyData, nil)

	historyUsecase := NewHistoryUsecase(suite.repoMock)
	historyList, err := historyUsecase.HistoryWithAmountFilter(*user, moreThan, lessThan)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), historyList, 3)
	assert.Equal(suite.T(), dummyData, historyList)
}

func (suite *HistoryUsecaseTestSuite) SetupTest() {
	suite.repoMock = new(HistoryRepoMock)
}

func TestHistoryUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(HistoryUsecaseTestSuite))
}
