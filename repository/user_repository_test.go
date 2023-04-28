package repository

import (
	"errors"
	"final_project_easycash/model"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sqlx.DB
	mockSql sqlmock.Sqlmock
}

func (suite *UserRepositoryTestSuite) TestUserGetUserById_Success() {
	row := sqlmock.NewRows([]string{"id", "username", "password", "email", "phone_number", "photo_profile", "balance"})
	row.AddRow(dummyUsers[0].Id, dummyUsers[0].Username, dummyUsers[0].Password, dummyUsers[0].Email, dummyUsers[0].PhoneNumber, dummyUsers[0].PhotoProfile, dummyUsers[0].Balance)

	suite.mockSql.ExpectQuery("SELECT (.*) FROM mst_user").WillReturnRows(row)
	repo := NewUserRepo(suite.mockDb)

	actual, err := repo.GetUserById(dummyUsers[0].Username)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, actual.Id)
}

func (suite *UserRepositoryTestSuite) TestUserGetUserById_Failed() {
	suite.mockSql.ExpectQuery("SELECT (.*) FROM mst_user").WillReturnError(errors.New("Failed"))
	repo := NewUserRepo(suite.mockDb)

	expected := model.User{}

	actual, err := repo.GetUserById(dummyUsers[0].Email)

	assert.Equal(suite.T(), expected, actual)
	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestUserUpdateUserById_Success() {
	updatedUserData := &dummyUsers[0]
	suite.mockSql.ExpectExec(`UPDATE mst_user SET password = \$1, email = \$2, phone_number = \$3 WHERE username = \$4`).WithArgs(updatedUserData.Password, updatedUserData.Email, updatedUserData.PhoneNumber, updatedUserData.Username).WillReturnResult(sqlmock.NewResult(0, 1))
	repo := NewUserRepo(suite.mockDb)

	err := repo.UpdateUserById(updatedUserData)

	assert.Nil(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestUserUpdateUserById_Failed() {
	updatedUserData := &dummyUsers[0]
	suite.mockSql.ExpectExec(`UPDATE mst_user SET password = \$1, email = \$2, phone_number = \$3 WHERE username = \$4`).WillReturnError(errors.New("Failed"))
	repo := NewUserRepo(suite.mockDb)

	err := repo.UpdateUserById(updatedUserData)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errors.New("Failed"), err)
}

func (suite *UserRepositoryTestSuite) TestUserUpdatePhotoProfile_Success() {
	updatedPhotoProfile := dummyUsers[0]
	suite.mockSql.ExpectExec(`UPDATE mst_user SET photo_profile = \$1 WHERE username = \$2`).WithArgs(updatedPhotoProfile.PhotoProfile, updatedPhotoProfile.Username).WillReturnResult(sqlmock.NewResult(0, 1))
	repo := NewUserRepo(suite.mockDb)

	err := repo.UpdatePhotoProfile(updatedPhotoProfile.Username, updatedPhotoProfile.PhotoProfile)

	assert.Nil(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestUserUpdatePhotoProfile_Failed() {
	updatedPhotoProfile := dummyUsers[0]
	suite.mockSql.ExpectExec(`UPDATE mst_user SET photo_profile = \$1 WHERE username = \$2`).WillReturnError(errors.New("Failed"))
	repo := NewUserRepo(suite.mockDb)

	err := repo.UpdatePhotoProfile(updatedPhotoProfile.Username, updatedPhotoProfile.PhotoProfile)

	assert.NotNil(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestUserDeleteUserById_Success() {
	deletedUser := dummyUsers[0]
	suite.mockSql.ExpectExec(`DELETE FROM mst_user WHERE username = \$1`).WithArgs(deletedUser.Username).WillReturnResult(sqlmock.NewResult(0, 1))
	repo := NewUserRepo(suite.mockDb)

	err := repo.DeleteUserById(deletedUser.Username)

	assert.Nil(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestUserDeleteUserById_Failed() {
	deletedUser := dummyUsers[0]
	suite.mockSql.ExpectExec(`DELETE FROM mst_user WHERE username = \$1`).WillReturnError(errors.New("Failed"))
	repo := NewUserRepo(suite.mockDb)

	err := repo.DeleteUserById(deletedUser.Username)

	assert.NotNil(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("An error when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	suite.mockDb = sqlxDB
	suite.mockSql = mockSql
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestCustomerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

//go test ./... -coverprofile=coverage.out
//go tool cover -html=coverage
