package repository

import (
	"database/sql"
	"log"
	"testing"

	"final_project_easycash/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var dummyUser = []model.User{
	{
		Username:    "userDummy1",
		Email:       "user1@gmail.com",
		PhoneNumber: "081234567891",
		Password:    "passwordUser1",
	},
	{
		Username:    "userDummy2",
		Email:       "user2@gmail.com",
		PhoneNumber: "081234567892",
		Password:    "passwordUser2",
	},
	{
		Username:    "userDummy3",
		Email:       "user3@gmail.com",
		PhoneNumber: "081234567893",
		Password:    "passwordUser3",
	},
}

type LoginRepoTestSuite struct {
	suite.Suite
	mockDb  *sqlx.DB
	mockSql sqlmock.Sqlmock
}

func (suite *LoginRepoTestSuite) TestFindUser_Success() {
	recUser := dummyUser[0]
	UserInDb := model.User{
		Username: "userDummy1",
		Password: "$2a$10$6wvkxozhPmUsP0sr8XciNOVPQM7XUZBYt1DeOfLI/4XRkM4YCkNiG", // hashed "passwordUser1"
	}

	rows := sqlmock.NewRows([]string{"username", "password"}).
		AddRow(UserInDb.Username, UserInDb.Password)
	suite.mockSql.ExpectQuery("SELECT username, password FROM mst_user WHERE username = (.+)").
		WithArgs(recUser.Username).
		WillReturnRows(rows)

	loginRepo := NewLoginRepo(suite.mockDb)
	result, message := loginRepo.FindUser(recUser)
	assert.True(suite.T(), result)
	assert.Equal(suite.T(), "successfully login", message)
}

func (suite *LoginRepoTestSuite) TestFindUserFailUserNotFound() {
	recUser := dummyUser[0]

	suite.mockSql.ExpectQuery("SELECT username, password FROM mst_user WHERE username = (.+)").
		WithArgs(recUser.Username).
		WillReturnError(sql.ErrNoRows)

	loginRepo := NewLoginRepo(suite.mockDb)
	result, message := loginRepo.FindUser(recUser)
	assert.False(suite.T(), result)
	assert.Equal(suite.T(), "user not found", message)
}

func (suite *LoginRepoTestSuite) TestFindUserFailInvalidPassword() {
	recUser := dummyUser[0]
	resUser := model.User{
		Username: "userDummy1",
		Password: "$2a$10$6wvkxozhPmUsP0sr8XciNOVPQM7XUZBYt1DeOfLI/4XRkM4YCkNiG", // hashed "passwordUser1"
	}

	rows := sqlmock.NewRows([]string{"username", "password"}).
		AddRow(resUser.Username, resUser.Password)
	suite.mockSql.ExpectQuery("SELECT username, password FROM mst_user WHERE username = (.+)").
		WithArgs(recUser.Username).
		WillReturnRows(rows)

	recUser.Password = "passwordUser2"

	loginRepo := NewLoginRepo(suite.mockDb)
	result, message := loginRepo.FindUser(recUser)
	assert.False(suite.T(), result)
	assert.Equal(suite.T(), "invalid password", message)
}

func (suite *LoginRepoTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("An error when opening a stub database connection", err)
	}

	db := sqlx.NewDb(mockDb, "postgres")

	suite.mockDb = db
	suite.mockSql = mockSql
}

func (suite *LoginRepoTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestLoginRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LoginRepoTestSuite))
}
