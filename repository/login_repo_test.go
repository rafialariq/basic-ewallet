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

type LoginRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sqlx.DB
	mockSql sqlmock.Sqlmock
}

func (suite *LoginRepositoryTestSuite) TestFindUser_Success() {
	row := sqlmock.NewRows([]string{"username", "password"})
	row.AddRow(dummyUsers[0].Username, "$2y$10$.0qAQ9W4smgbIAbf/zwuseNnbC7.baKDY41IIFNsQcxk2UEdTPzcy")

	suite.mockSql.ExpectQuery(`SELECT username, password FROM mst_user WHERE username = \$1`).WithArgs(dummyUsers[0].Username).WillReturnRows(row)
	repo := NewLoginRepo(suite.mockDb)

	actual, err := repo.FindUser(dummyUsers[0])

	assert.Equal(suite.T(), "successfully login", err)
	assert.Equal(suite.T(), true, actual)
}

func (suite *LoginRepositoryTestSuite) TestFindUserInvalidPassword_Failed() {
	row := sqlmock.NewRows([]string{"username", "password"})
	row.AddRow(dummyUsers[1].Username, "$2y$10$.0qAQ9W4smgbIAbf/zwuseNnbC7.baKDY41IIFNsQcxk2UEdTPzcy")

	suite.mockSql.ExpectQuery(`SELECT username, password FROM mst_user WHERE username = \$1`).WithArgs(dummyUsers[1].Username).WillReturnRows(row)
	repo := NewLoginRepo(suite.mockDb)

	actual, err := repo.FindUser(dummyUsers[1])

	assert.Equal(suite.T(), "invalid password", err)
	assert.Equal(suite.T(), false, actual)
}

func (suite *LoginRepositoryTestSuite) TestFindUserUserNotFound_Failed() {
	row := sqlmock.NewRows([]string{"username", "password"})
	row.AddRow(dummyUsers[0].Username, "$2y$10$.0qAQ9W4smgbIAbf/zwuseNnbC7.baKDY41IIFNsQcxk2UEdTPzcy")

	suite.mockSql.ExpectQuery(`SELECT username, password FROM mst_user WHERE username = \$1`).WithArgs(dummyUsers[1].Username).WillReturnError(errors.New("failed"))
	repo := NewLoginRepo(suite.mockDb)

	actual, err := repo.FindUser(dummyUsers[1])

	assert.Equal(suite.T(), "user not found", err)
	assert.Equal(suite.T(), false, actual)
}

func (suite *LoginRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("An error when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	suite.mockDb = sqlxDB
	suite.mockSql = mockSql
}

func (suite *LoginRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestLoginRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LoginRepositoryTestSuite))
}
