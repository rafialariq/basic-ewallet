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

type TransactionRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sqlx.DB
	mockSql sqlmock.Sqlmock
}

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

func (suite *TransactionRepositoryTestSuite) TestTransferMoney_Success() {
	sender := dummyUsers[0]
	receiver := dummyMerchants[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowMerchant := sqlmock.NewRows([]string{"merchantcode"})
	rowMerchant.AddRow(dummyMerchants[0].MerchantCode)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT merchantcode FROM mst_merchant WHERE merchantcode \= \$1`).
		WithArgs(receiver.MerchantCode).
		WillReturnRows(rowMerchant)
	suite.mockSql.ExpectExec("BEGIN;").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mockSql.ExpectExec(`INSERT INTO trx_bill \(sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
		WithArgs(1, sender.PhoneNumber, 2, amount, time.Now(), 3, receiver.MerchantCode).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(`UPDATE mst_user SET balance \= balance \- \$1 WHERE phone_number \= \$2;`).
		WithArgs(amount, sender.PhoneNumber).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(`UPDATE mst_merchant SET amount \= amount \+ \$1 WHERE merchantcode \= \$2;`).
		WithArgs(amount, receiver.MerchantCode).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec("COMMIT;").WillReturnResult(sqlmock.NewResult(0, 0))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.TransferMoney(sender.PhoneNumber, receiver.MerchantCode, amount)

	assert.Nil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestTransferMoneyCheckBalanceQuery_Failed() {
	sender := dummyUsers[0]
	receiver := dummyMerchants[0]
	amount := 15000.00

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.TransferMoney(sender.PhoneNumber, receiver.MerchantCode, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestTransferMoneyCheckBalance_Failed() {
	sender := dummyUsers[0]
	receiver := dummyMerchants[0]
	amount := 150000.00
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.TransferMoney(sender.PhoneNumber, receiver.MerchantCode, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestTransferMoneyCheckPhoneNumber_Failed() {
	sender := dummyUsers[0]
	receiver := dummyMerchants[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.TransferMoney(sender.PhoneNumber, receiver.MerchantCode, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestTransferMoneyCheckMerchant_Failed() {
	sender := dummyUsers[0]
	receiver := dummyMerchants[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT merchantcode FROM mst_merchant WHERE merchantcode \= \$1`).
		WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.TransferMoney(sender.PhoneNumber, receiver.MerchantCode, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestTransferMoneyBegin_Failed() {
	sender := dummyUsers[0]
	receiver := dummyMerchants[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowMerchant := sqlmock.NewRows([]string{"merchantcode"})
	rowMerchant.AddRow(dummyMerchants[0].MerchantCode)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT merchantcode FROM mst_merchant WHERE merchantcode \= \$1`).
		WithArgs(receiver.MerchantCode).
		WillReturnRows(rowMerchant)
	suite.mockSql.ExpectExec("BEGIN;").WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.TransferMoney(sender.PhoneNumber, receiver.MerchantCode, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestTransferMoneyInsert_Failed() {
	sender := dummyUsers[0]
	receiver := dummyMerchants[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowMerchant := sqlmock.NewRows([]string{"merchantcode"})
	rowMerchant.AddRow(dummyMerchants[0].MerchantCode)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT merchantcode FROM mst_merchant WHERE merchantcode \= \$1`).
		WithArgs(receiver.MerchantCode).
		WillReturnRows(rowMerchant)
	suite.mockSql.ExpectExec("BEGIN;").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mockSql.ExpectExec(`INSERT INTO trx_bill \(sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
		WillReturnError(errors.New("failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.TransferMoney(sender.PhoneNumber, receiver.MerchantCode, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestTransferMoneyUpdateSenderBalance_Failed() {
	sender := dummyUsers[0]
	receiver := dummyMerchants[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowMerchant := sqlmock.NewRows([]string{"merchantcode"})
	rowMerchant.AddRow(dummyMerchants[0].MerchantCode)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT merchantcode FROM mst_merchant WHERE merchantcode \= \$1`).
		WithArgs(receiver.MerchantCode).
		WillReturnRows(rowMerchant)
	suite.mockSql.ExpectExec("BEGIN;").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mockSql.ExpectExec(`INSERT INTO trx_bill \(sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
		WithArgs(1, sender.PhoneNumber, 2, amount, time.Now(), 3, receiver.MerchantCode).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(`UPDATE mst_user SET balance \= balance \- \$1 WHERE phone_number \= \$2;`).
		WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.TransferMoney(sender.PhoneNumber, receiver.MerchantCode, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestTransferMoneyUpdateReceiverBalance_Failed() {
	sender := dummyUsers[0]
	receiver := dummyMerchants[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowMerchant := sqlmock.NewRows([]string{"merchantcode"})
	rowMerchant.AddRow(dummyMerchants[0].MerchantCode)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT merchantcode FROM mst_merchant WHERE merchantcode \= \$1`).
		WithArgs(receiver.MerchantCode).
		WillReturnRows(rowMerchant)
	suite.mockSql.ExpectExec("BEGIN;").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mockSql.ExpectExec(`INSERT INTO trx_bill \(sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
		WithArgs(1, sender.PhoneNumber, 2, amount, time.Now(), 3, receiver.MerchantCode).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(`UPDATE mst_user SET balance \= balance \- \$1 WHERE phone_number \= \$2;`).
		WithArgs(amount, sender.PhoneNumber).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(`UPDATE mst_merchant SET amount \= amount \+ \$1 WHERE merchantcode \= \$2;`).
		WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.TransferMoney(sender.PhoneNumber, receiver.MerchantCode, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestTransferMoneyCommit_Failed() {
	sender := dummyUsers[0]
	receiver := dummyMerchants[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowMerchant := sqlmock.NewRows([]string{"merchantcode"})
	rowMerchant.AddRow(dummyMerchants[0].MerchantCode)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT merchantcode FROM mst_merchant WHERE merchantcode \= \$1`).
		WithArgs(receiver.MerchantCode).
		WillReturnRows(rowMerchant)
	suite.mockSql.ExpectExec("BEGIN;").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mockSql.ExpectExec(`INSERT INTO trx_bill \(sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
		WithArgs(1, sender.PhoneNumber, 2, amount, time.Now(), 3, receiver.MerchantCode).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(`UPDATE mst_user SET balance \= balance \- \$1 WHERE phone_number \= \$2;`).
		WithArgs(amount, sender.PhoneNumber).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(`UPDATE mst_merchant SET amount \= amount \+ \$1 WHERE merchantcode \= \$2;`).
		WithArgs(amount, receiver.MerchantCode).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec("COMMIT;").WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.TransferMoney(sender.PhoneNumber, receiver.MerchantCode, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestWithdrawBalance_Success() {
	sender := dummyUsers[0]
	receiver := dummyBanks[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowBank := sqlmock.NewRows([]string{"bank_number"})
	rowBank.AddRow(dummyBanks[0].BankNumber)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT bank_number FROM mst_bank WHERE bank_number \= \$1`).
		WithArgs(receiver.BankNumber).
		WillReturnRows(rowBank)
	suite.mockSql.ExpectExec("BEGIN;").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mockSql.ExpectExec(`INSERT INTO trx_bill \(sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
		WithArgs(1, sender.PhoneNumber, 3, amount, time.Now(), 2, receiver.BankNumber).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(`UPDATE mst_user SET balance \= balance \- \$1 WHERE phone_number \= \$2;`).
		WithArgs(amount, sender.PhoneNumber).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec("COMMIT;").WillReturnResult(sqlmock.NewResult(0, 0))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.WithdrawBalance(sender.PhoneNumber, receiver.BankNumber, amount)

	assert.Nil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestWithdrawBalanceBalance_Failed() {
	sender := dummyUsers[0]
	receiver := dummyBanks[0]
	amount := 15000.00

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.WithdrawBalance(sender.PhoneNumber, receiver.BankNumber, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestWithdrawPhoneNumber_Failed() {
	sender := dummyUsers[0]
	receiver := dummyBanks[0]
	amount := 15000.00
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.WithdrawBalance(sender.PhoneNumber, receiver.BankNumber, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestWithdrawBalanceInsert_Failed() {
	sender := dummyUsers[0]
	receiver := dummyBanks[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowBank := sqlmock.NewRows([]string{"bank_number"})
	rowBank.AddRow(dummyBanks[0].BankNumber)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT bank_number FROM mst_bank WHERE bank_number \= \$1`).
		WithArgs(receiver.BankNumber).
		WillReturnRows(rowBank)
	suite.mockSql.ExpectExec("BEGIN;").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mockSql.ExpectExec(`INSERT INTO trx_bill \(sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
		WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.WithdrawBalance(sender.PhoneNumber, receiver.BankNumber, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestWithdrawBalanceUpdate_Failed() {
	sender := dummyUsers[0]
	receiver := dummyBanks[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowBank := sqlmock.NewRows([]string{"bank_number"})
	rowBank.AddRow(dummyBanks[0].BankNumber)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT bank_number FROM mst_bank WHERE bank_number \= \$1`).
		WithArgs(receiver.BankNumber).
		WillReturnRows(rowBank)
	suite.mockSql.ExpectExec(`INSERT INTO trx_bill \(sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
		WithArgs(1, sender.PhoneNumber, 3, amount, time.Now(), 2, receiver.BankNumber).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec("BEGIN;").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mockSql.ExpectExec(`UPDATE mst_user SET balance \= balance \- \$1 WHERE phone_number \= \$2;`).
		WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.WithdrawBalance(sender.PhoneNumber, receiver.BankNumber, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestWithdrawBalanceBegin_Failed() {
	sender := dummyUsers[0]
	receiver := dummyBanks[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowBank := sqlmock.NewRows([]string{"bank_number"})
	rowBank.AddRow(dummyBanks[0].BankNumber)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT bank_number FROM mst_bank WHERE bank_number \= \$1`).
		WithArgs(receiver.BankNumber).
		WillReturnRows(rowBank)
	suite.mockSql.ExpectExec("BEGIN;").WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.WithdrawBalance(sender.PhoneNumber, receiver.BankNumber, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestWithdrawBalanceCommit_Failed() {
	sender := dummyUsers[0]
	receiver := dummyBanks[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)
	rowBank := sqlmock.NewRows([]string{"bank_number"})
	rowBank.AddRow(dummyBanks[0].BankNumber)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT bank_number FROM mst_bank WHERE bank_number \= \$1`).
		WithArgs(receiver.BankNumber).
		WillReturnRows(rowBank)
	suite.mockSql.ExpectExec(`INSERT INTO trx_bill \(sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
		WithArgs(1, sender.PhoneNumber, 3, amount, time.Now(), 2, receiver.BankNumber).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec("BEGIN;").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mockSql.ExpectExec(`UPDATE mst_user SET balance \= balance \- \$1 WHERE phone_number \= \$2;`).
		WithArgs(amount, sender.PhoneNumber).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec("COMMIT;").WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.WithdrawBalance(sender.PhoneNumber, receiver.BankNumber, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) TestWithdrawBalanceBankNumber_Failed() {
	sender := dummyUsers[0]
	receiver := dummyBanks[0]
	amount := 15000.00
	rowUserPhoneNumber := sqlmock.NewRows([]string{"phone_number"})
	rowUserPhoneNumber.AddRow(dummyUsers[0].PhoneNumber)
	rowUserBalance := sqlmock.NewRows([]string{"balance"})
	rowUserBalance.AddRow(dummyUsers[0].Balance)

	suite.mockSql.ExpectQuery(`SELECT balance FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserBalance)
	suite.mockSql.ExpectQuery(`SELECT phone_number FROM mst_user WHERE phone_number \= \$1`).
		WithArgs(sender.PhoneNumber).
		WillReturnRows(rowUserPhoneNumber)
	suite.mockSql.ExpectQuery(`SELECT bank_number FROM mst_bank WHERE bank_number \= \$1`).
		WillReturnError(errors.New("Failed"))
	repo := NewTransactionRepo(suite.mockDb)
	actual := repo.WithdrawBalance(sender.PhoneNumber, receiver.BankNumber, amount)

	assert.NotNil(suite.T(), actual)
}

func (suite *TransactionRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("An error when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	suite.mockDb = sqlxDB
	suite.mockSql = mockSql
}

func (suite *TransactionRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestTransactionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
