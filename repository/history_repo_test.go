package repository

import (
	"errors"
	"final_project_easycash/model"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
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

type HistoryRepoTestSuite struct {
	suite.Suite
	mockDb  *sqlx.DB
	mockSql sqlmock.Sqlmock
}

func (suite *HistoryRepoTestSuite) TestGetHistoryByUser_Success() {
	user := model.User{PhoneNumber: "082123456789"}

	rows := sqlmock.NewRows([]string{"id", "id_transaction", "sender_type_id", "sender_id", "type_id", "amount", "date", "destination_type_id", "destination_id", "status"})
	for _, v := range dummyData {
		rows.AddRow(v.Id, v.TransactionId, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.Date, v.DestinationTypeId, v.DestinationId, v.Status)
	}

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date,  destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1);"
	suite.mockSql.ExpectQuery(query).
		WithArgs(&user.PhoneNumber).WillReturnRows(rows)

	historyRepo := NewHistoryRepo(suite.mockDb)
	historyList, err := historyRepo.GetHistoryByUser(user)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), historyList, 3)
	assert.Equal(suite.T(), dummyData[0].Id, historyList[0].Id)
	assert.Equal(suite.T(), dummyData[0].Amount, historyList[0].Amount)
	assert.Equal(suite.T(), dummyData[1].Id, historyList[1].Id)
	assert.Equal(suite.T(), dummyData[1].Amount, historyList[1].Amount)
}

func (suite *HistoryRepoTestSuite) TestGetHistoryByUser_FailedNoArg() {
	user := model.User{}

	rows := sqlmock.NewRows([]string{"id", "id_transaction", "sender_type_id", "sender_id", "type_id", "amount", "date", "destination_type_id", "destination_id", "status"})
	for _, v := range dummyData {
		rows.AddRow(v.Id, v.TransactionId, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.Date, v.DestinationTypeId, v.DestinationId, v.Status)
	}

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date,  destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1);"
	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("failed"))

	historyRepo := NewHistoryRepo(suite.mockDb)

	historyList, err := historyRepo.GetHistoryByUser(user)

	assert.Nil(suite.T(), historyList)
	assert.Error(suite.T(), err)
}

// func (suite *HistoryRepoTestSuite) TestGetHistoryByUser_FailedNoData() {
// 	user := model.User{PhoneNumber: "082123456789"}

// 	rows := sqlmock.NewRows([]string{"id", "sender_type_id", "sender_id", "type_id", "destination_type_id", "destination_id"})
// 	for _, v := range dummyData {
// 		rows.AddRow(v.Id, v.SenderTypeId, v.SenderId, v.TypeId, v.DestinationTypeId, v.DestinationId)
// 	}

// 	query := "SELECT id, sender_type_id, sender_id, type_id, amount, destination_type_id, destination_id FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1);"
// 	suite.mockSql.ExpectQuery(query).
// 		WithArgs(&user.PhoneNumber).WillReturnError(errors.New("failed no data"))

// 	historyRepo := NewHistoryRepo(suite.mockDb)
// 	historyList, err := historyRepo.GetHistoryByUser(user)

// 	assert.Nil(suite.T(), historyList)
// 	assert.Error(suite.T(), err)
// }

func (suite *HistoryRepoTestSuite) TestGetHistoryWithAccountFilter_Success() {
	user := model.User{PhoneNumber: "082123456789"}
	accountTypeId := 1

	rows := sqlmock.NewRows([]string{"id", "id_transaction", "sender_type_id", "sender_id", "type_id", "amount", "date", "destination_type_id", "destination_id", "status"})
	for _, v := range dummyData {
		rows.AddRow(v.Id, v.TransactionId, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.Date, v.DestinationTypeId, v.DestinationId, v.Status)
	}

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND (sender_type_id = $2 OR destination_type_id = $2);"
	suite.mockSql.ExpectQuery(query).
		WithArgs(&user.PhoneNumber, accountTypeId).WillReturnRows(rows)

	historyRepo := NewHistoryRepo(suite.mockDb)
	historyList, err := historyRepo.GetHistoryWithAccountFilter(user, accountTypeId)

	assert.Len(suite.T(), historyList, 3)
	assert.NoError(suite.T(), err)

}

func (suite *HistoryRepoTestSuite) TestGetHistoryWithAccountFilter_FailedNoArgs() {
	user := model.User{}
	accountTypeId := 1

	rows := sqlmock.NewRows([]string{"id", "id_transaction", "sender_type_id", "sender_id", "type_id", "amount", "date", "destination_type_id", "destination_id", "status"})
	for _, v := range dummyData {
		rows.AddRow(v.Id, v.TransactionId, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.Date, v.DestinationTypeId, v.DestinationId, v.Status)

	}

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND (sender_type_id = $2 OR destination_type_id = $2);"
	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("Failed"))

	historyRepo := NewHistoryRepo(suite.mockDb)
	historyList, err := historyRepo.GetHistoryWithAccountFilter(user, accountTypeId)

	assert.Nil(suite.T(), historyList)
	assert.Error(suite.T(), err)

}

// func (suite *HistoryRepoTestSuite) TestGetHistoryWithAccountFilter_FailedNoData() {
// 	user := model.User{}
// 	accountTypeId := 1

// 	rows := sqlmock.NewRows([]string{"id", "sender_type_id", "sender_id", "type_id", "amount", "destination_type_id", "destination_id"})
// 	for _, v := range dummyData {
// 		rows.AddRow(v.Id, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.DestinationTypeId, v.DestinationId)
// 	}

// 	query := "SELECT id, sender_type_id, sender_id, type_id, amount, destination_type_id, destination_id FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND (sender_type_id = $2 OR destination_type_id = $2);"
// 	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("Failed"))

// 	historyRepo := NewHistoryRepo(suite.mockDb)
// 	historyList, err := historyRepo.GetHistoryWithAccountFilter(user, accountTypeId)

// 	assert.Nil(suite.T(), historyList)
// 	assert.Error(suite.T(), err)

// }

func (suite *HistoryRepoTestSuite) TestGetHistoryWithTypeFilter_Success() {
	user := model.User{PhoneNumber: "082123456789"}
	typeId := 1

	rows := sqlmock.NewRows([]string{"id", "id_transaction", "sender_type_id", "sender_id", "type_id", "amount", "date", "destination_type_id", "destination_id", "status"})
	for _, v := range dummyData {
		rows.AddRow(v.Id, v.TransactionId, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.Date, v.DestinationTypeId, v.DestinationId, v.Status)
	}

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND type_id = $2;"
	suite.mockSql.ExpectQuery(query).
		WithArgs(&user.PhoneNumber, &typeId).WillReturnRows(rows)

	historyRepo := NewHistoryRepo(suite.mockDb)
	historyList, err := historyRepo.GetHistoryWithTypeFilter(user, typeId)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), historyList, 3)

}

func (suite *HistoryRepoTestSuite) TestGetHistoryWithTypeFilter_FailedNoArg() {
	user := model.User{}
	typeId := 1

	rows := sqlmock.NewRows([]string{"id", "id_transaction", "sender_type_id", "sender_id", "type_id", "amount", "date", "destination_type_id", "destination_id", "status"})
	for _, v := range dummyData {
		rows.AddRow(v.Id, v.TransactionId, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.Date, v.DestinationTypeId, v.DestinationId, v.Status)
	}

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND type_id = $2;"
	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("Failed"))

	historyRepo := NewHistoryRepo(suite.mockDb)
	historyList, err := historyRepo.GetHistoryWithTypeFilter(user, typeId)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), historyList)

}

// func (suite *HistoryRepoTestSuite) TestGetHistoryWithTypeFilter_FailedNoData() {
// 	user := model.User{PhoneNumber: "082123456789"}
// 	typeId := "1"

// 	rows := sqlmock.NewRows([]string{"id", "sender_type_id", "sender_id", "type_id", "amount", "destination_type_id", "destination_id"})
// 	for _, v := range dummyData {
// 		rows.AddRow(v.Id, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.DestinationTypeId, v.DestinationId)
// 	}

// 	query := "SELECT id, sender_type_id, sender_id, type_id, amount, destination_type_id, destination_id FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND type_id = $2;"
// 	suite.mockSql.ExpectQuery(query).
// 		WithArgs(&user.PhoneNumber, &typeId).WillReturnRows(rows)

// 	historyRepo := NewHistoryRepo(suite.mockDb)
// 	historyList, err := historyRepo.GetHistoryWithTypeFilter(user, typeId)

// 	assert.NoError(suite.T(), err)
// 	assert.Len(suite.T(), historyList, 3)

// }

func (suite *HistoryRepoTestSuite) TestGetHistoryWithAmountFilter_Success() {
	user := model.User{PhoneNumber: "082123456789"}
	var moreThan float64 = 50000
	var lessThan float64 = 100000

	rows := sqlmock.NewRows([]string{"id", "id_transaction", "sender_type_id", "sender_id", "type_id", "amount", "date", "destination_type_id", "destination_id", "status"})
	for _, v := range dummyData {
		rows.AddRow(v.Id, v.TransactionId, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.Date, v.DestinationTypeId, v.DestinationId, v.Status)
	}

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND amount >= $2 AND amount <= $3;"
	suite.mockSql.ExpectQuery(query).
		WithArgs(&user.PhoneNumber, &moreThan, &lessThan).WillReturnRows(rows)

	historyRepo := NewHistoryRepo(suite.mockDb)
	historyList, err := historyRepo.GetHistoryWithAmountFilter(user, moreThan, lessThan)

	assert.Len(suite.T(), historyList, 3)
	assert.NoError(suite.T(), err)

}

func (suite *HistoryRepoTestSuite) TestGetHistoryWithAmountFilter_FailedNoArg() {
	user := model.User{}
	var moreThan float64 = 50000
	var lessThan float64 = 100000

	rows := sqlmock.NewRows([]string{"id", "id_transaction", "sender_type_id", "sender_id", "type_id", "amount", "date", "destination_type_id", "destination_id", "status"})
	for _, v := range dummyData {
		rows.AddRow(v.Id, v.TransactionId, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.Date, v.DestinationTypeId, v.DestinationId, v.Status)
	}

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND amount >= $2 AND amount <= $3;"
	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("Failed"))

	historyRepo := NewHistoryRepo(suite.mockDb)
	historyList, err := historyRepo.GetHistoryWithAmountFilter(user, moreThan, lessThan)

	assert.Nil(suite.T(), historyList)
	assert.Error(suite.T(), err)

}

// func (suite *HistoryRepoTestSuite) TestGetHistoryWithAmountFilter_FailedNoData() {
// 	user := model.User{}
// 	var moreThan float64 = 50000
// 	var lessThan float64 = 100000

// 	rows := sqlmock.NewRows([]string{"id", "sender_type_id", "sender_id", "type_id", "amount", "destination_type_id", "destination_id"})
// 	for _, v := range dummyData {
// 		rows.AddRow(v.Id, v.SenderTypeId, v.SenderId, v.TypeId, v.Amount, v.DestinationTypeId, v.DestinationId)
// 	}

// 	query := "SELECT id, sender_type_id, sender_id, type_id, amount, destination_type_id, destination_id FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND amount >= $2 AND amount <= $3;"
// 	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("Failed"))

// 	historyRepo := NewHistoryRepo(suite.mockDb)
// 	historyList, err := historyRepo.GetHistoryWithAmountFilter(user, moreThan, lessThan)

// 	assert.Nil(suite.T(), historyList)
// 	assert.Error(suite.T(), err)

// }

func (suite *HistoryRepoTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalln("An error when opening a stub database connection", err)
	}

	db := sqlx.NewDb(mockDb, "postgres")

	suite.mockDb = db
	suite.mockSql = mockSql
}

func (suite *HistoryRepoTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestHistoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(HistoryRepoTestSuite))
}
