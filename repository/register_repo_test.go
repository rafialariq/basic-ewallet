package repository

import (
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegisterRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sqlx.DB
	mockSql sqlmock.Sqlmock
}

func (suite *RegisterRepositoryTestSuite) TestUserRegister_Success() {
	newUser := &dummyUsers[0]
	suite.mockSql.ExpectExec(`INSERT INTO mst_user \(username, email, phone_number, password\) VALUES \(\$1, \$2, \$3, \$4\);`).
		WithArgs(newUser.Username, newUser.Email, newUser.PhoneNumber, newUser.Password).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(`INSERT INTO mst_transaction_codes \(code\) VALUES \(\$1\)`).
		WithArgs(newUser.PhoneNumber).
		WillReturnResult(sqlmock.NewResult(0, 1))
	repo := NewRegisterRepo(suite.mockDb)
	actual, err := repo.UserRegister(newUser)

	assert.Equal(suite.T(), "user created successfully", err)
	assert.Equal(suite.T(), true, actual)
}

func (suite *RegisterRepositoryTestSuite) TestUserRegisterFirstQuery_Failed() {
	newUser := &dummyUsers[0]
	suite.mockSql.ExpectExec(`INSERT INTO mst_user \(username, email, phone_number, password\) VALUES \(\$1, \$2, \$3, \$4\);`).
		WillReturnError(errors.New("Failed"))
	repo := NewRegisterRepo(suite.mockDb)
	actual, err := repo.UserRegister(newUser)

	assert.Equal(suite.T(), "failed to create user", err)
	assert.Equal(suite.T(), false, actual)
}

func (suite *RegisterRepositoryTestSuite) TestUserRegisterSecondQuery_Failed() {
	newUser := &dummyUsers[1]
	suite.mockSql.ExpectExec(`INSERT INTO mst_user \(username, email, phone_number, password\) VALUES \(\$1, \$2, \$3, \$4\);`).
		WithArgs(newUser.Username, newUser.Email, newUser.PhoneNumber, newUser.Password).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(`INSERT INTO mst_transaction_codes \(code\) VALUES \(\$1\)`).
		WillReturnError(errors.New("Failed"))
	repo := NewRegisterRepo(suite.mockDb)
	actual, err := repo.UserRegister(newUser)

	assert.Equal(suite.T(), "failed to create user", err)
	assert.Equal(suite.T(), false, actual)
}

func (suite *RegisterRepositoryTestSuite) TestRegisterValidate_Success() {
	row := sqlmock.NewRows([]string{"username", "password"})
	row.AddRow(dummyUsers[0].Username, "$2y$10$.0qAQ9W4smgbIAbf/zwuseNnbC7.baKDY41IIFNsQcxk2UEdTPzcy")
	newUser := &dummyUsers[1]

	suite.mockSql.ExpectQuery(`SELECT username, phone_number FROM mst_user WHERE username \= \$1 OR phone_number \= \$2;`).
		WithArgs(newUser.Username, newUser.PhoneNumber).
		WillReturnRows(row)
	repo := NewRegisterRepo(suite.mockDb)
	actual := repo.RegisterValidate(newUser)

	assert.Equal(suite.T(), false, actual)
}

func (suite *RegisterRepositoryTestSuite) TestRegisterValidate_Failed() {
	row := sqlmock.NewRows([]string{"username", "password"})
	row.AddRow(dummyUsers[0].Username, "$2y$10$.0qAQ9W4smgbIAbf/zwuseNnbC7.baKDY41IIFNsQcxk2UEdTPzcy")
	newUser := &dummyUsers[0]

	suite.mockSql.ExpectQuery(`SELECT username, phone_number FROM mst_user WHERE username \= \$1 OR phone_number \= \$2;`).
		WithArgs(newUser.Username, newUser.PhoneNumber).
		WillReturnRows(row)
	repo := NewRegisterRepo(suite.mockDb)
	actual := repo.RegisterValidate(newUser)

	assert.Equal(suite.T(), true, actual)
}

func (suite *RegisterRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("An error when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	suite.mockDb = sqlxDB
	suite.mockSql = mockSql
}

func (suite *RegisterRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestRegisterRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterRepositoryTestSuite))
}
