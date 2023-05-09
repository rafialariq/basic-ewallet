package repository

import (
	"log"
	"testing"

	"final_project_easycash/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegisterRepoTestSuite struct {
	suite.Suite
	mockDb  *sqlx.DB
	mockSql sqlmock.Sqlmock
}

var dummyNewUser = []model.User{
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

func (suite *RegisterRepoTestSuite) TestUserRegister_Success() {
	newUser := dummyNewUser[0]
	suite.mockSql.ExpectExec("INSERT INTO mst_user").WithArgs(newUser.Username, newUser.Email, newUser.PhoneNumber, newUser.Password).WillReturnResult(sqlmock.NewResult(1, 1))
	registerRepo := NewRegisterRepo(suite.mockDb)
	user, res := registerRepo.UserRegister(&newUser)

	assert.True(suite.T(), user)
	assert.Equal(suite.T(), "user created successfully", res)
}

func (suite *RegisterRepoTestSuite) TestUserRegister_Failed() {
	newUser := dummyNewUser[0]
	suite.mockSql.ExpectExec("INSERT INTO mst_user")
	registerRepo := NewRegisterRepo(suite.mockDb)
	user, res := registerRepo.UserRegister(&newUser)

	assert.False(suite.T(), user)
	assert.Equal(suite.T(), "failed to create user", res)

}

func (suite *RegisterRepoTestSuite) TestRegisterValidate_UsernameFound() {
	recUser := &dummyNewUser[1]

	rows := sqlmock.NewRows([]string{"username", "phone_number"}).AddRow("userDummy2", "081234567892")
	query := "SELECT username, phone_number FROM mst_user WHERE username = \\$1 OR phone_number = \\$2;"
	suite.mockSql.ExpectQuery(query).WithArgs(recUser.Username, recUser.PhoneNumber).WillReturnRows(rows)

	registerRepo := NewRegisterRepo(suite.mockDb)
	result := registerRepo.RegisterValidate(recUser)

	assert.True(suite.T(), result)
	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}

func (suite *RegisterRepoTestSuite) TestRegisterValidate_PhoneNumberFound() {
	recUser := &dummyNewUser[2]

	rows := sqlmock.NewRows([]string{"username", "phone_number"}).AddRow("userDummy3", "081234567893")
	query := "SELECT username, phone_number FROM mst_user WHERE username = \\$1 OR phone_number = \\$2;"
	suite.mockSql.ExpectQuery(query).WithArgs(recUser.Username, recUser.PhoneNumber).WillReturnRows(rows)

	registerRepo := NewRegisterRepo(suite.mockDb)
	result := registerRepo.RegisterValidate(recUser)

	assert.True(suite.T(), result)
	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}

func (suite *RegisterRepoTestSuite) TestRegisterValidate_NotFound() {
	recUser := &dummyNewUser[0]

	rows := sqlmock.NewRows([]string{"userDummy1", "081234567891"})
	query := "SELECT username, phone_number FROM mst_user WHERE username = \\$1 OR phone_number = \\$2;"
	suite.mockSql.ExpectQuery(query).WithArgs(recUser.Username, recUser.PhoneNumber).WillReturnRows(rows)

	registerRepo := NewRegisterRepo(suite.mockDb)
	result := registerRepo.RegisterValidate(recUser)

	assert.False(suite.T(), result)
	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}

func (suite *RegisterRepoTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("An error when opening a stub database connection", err)
	}

	db := sqlx.NewDb(mockDb, "postgres")

	suite.mockDb = db
	suite.mockSql = mockSql
}

func (suite *RegisterRepoTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestRegisterRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterRepoTestSuite))
}
