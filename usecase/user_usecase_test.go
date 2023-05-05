package usecase

import (
	"encoding/base64"
	"errors"
	"final_project_easycash/model"
	"io/ioutil"
	"mime/multipart"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type userRepoMock struct {
	mock.Mock
}

type fileRepoMock struct {
	mock.Mock
}

type utilsMock struct {
	mock.Mock
}

type UserUsecaseTestSuite struct {
	utilsMock    *utilsMock
	fileRepoMock *fileRepoMock
	userRepoMock *userRepoMock
	suite.Suite
}

func (u *userRepoMock) GetUserById(username string) (model.User, error) {
	args := u.Called(username)
	if args == nil {
		return model.User{}, errors.New("Failed")
	}
	return args.Get(0).(model.User), nil
}
func (u *userRepoMock) UpdateUserById(updatedUserData *model.User) error {
	args := u.Called(updatedUserData)
	if args[0] != nil {
		return errors.New("Failed")
	}
	return nil
}
func (u *userRepoMock) UpdatePhotoProfile(username string, filePath string) error {
	args := u.Called(username, filePath)
	if args != nil {
		return args.Error(0)
	}
	return nil
}
func (u *userRepoMock) DeleteUserById(username string) error {
	args := u.Called(username)
	if args != nil {
		return args.Error(0)
	}
	return nil
}

func (f *fileRepoMock) Save(fileName string, file *multipart.File) (string, error) {
	args := f.Called(fileName, file)
	if args != nil {
		return "", args.Error(4)
	}
	return "Dummy File Path", nil
}

func (u *utilsMock) ValidatePhoneNumber(phone string) bool {
	args := u.Called(phone)
	if args == nil {
		return false
	}
	return true
}

func (u *utilsMock) ValidateEmail(email string) bool {
	args := u.Called(email)
	if args == nil {
		return true
	}
	return false
}

func (suite *UserUsecaseTestSuite) TestCheckProfile_Success() {
	userUsecase := NewUserUsecase(suite.userRepoMock, suite.fileRepoMock)
	suite.userRepoMock.On("GetUserById", dummyUsers[0].Username).Return(dummyUsers[0], nil)
	user, err := userUsecase.CheckProfile(dummyUsers[0].Username)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyUsers[0], user)
}

func (suite *UserUsecaseTestSuite) TestCheckProfile_EncodePhotoProfile_Success() {
	// Create a temporary file with some content
	file, err := ioutil.TempFile("", "test-image.*")
	require.NoError(suite.T(), err)
	defer os.Remove(file.Name())

	_, err = file.WriteString("test content")
	require.NoError(suite.T(), err)

	// Create a user with a photo profile pointing to the temporary file
	user := model.User{
		Username:     "testuser",
		PhotoProfile: file.Name(),
	}

	// Create a user repo mock that returns the user
	suite.userRepoMock.On("GetUserById", "testuser").Return(user, nil)

	// Create a user usecase with the user repo mock
	usecase := NewUserUsecase(suite.userRepoMock, nil)

	// Call the CheckProfile function
	result, err := usecase.CheckProfile("testuser")
	require.NoError(suite.T(), err)

	// Verify that the photo profile was encoded correctly
	expectedEncoded := base64.StdEncoding.EncodeToString([]byte("test content"))
	assert.Equal(suite.T(), expectedEncoded, result.PhotoProfile)
}

func (suite *UserUsecaseTestSuite) TestCheckProfile_EncodePhotoProfile_Failed() {
	user := model.User{
		Username:     "testuser",
		PhotoProfile: "non-existent-file.jpg",
	}
	suite.userRepoMock.On("GetUserById", "testuser").Return(user, nil)

	usecase := NewUserUsecase(suite.userRepoMock, nil)
	result, err := usecase.CheckProfile("testuser")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.User{}, result)
}

func (suite *UserUsecaseTestSuite) TestEditProfile_Success() {
	userUsecase := NewUserUsecase(suite.userRepoMock, suite.fileRepoMock)
	suite.utilsMock.On("ValidateEmail", &dummyUsers[0].Email).Return(false)
	suite.utilsMock.On("ValidatePhoneNumber", &dummyUsers[0].PhoneNumber).Return(true)
	suite.userRepoMock.On("UpdateUserById", &dummyUsers[0]).Return(nil)
	err := userUsecase.EditProfile(&dummyUsers[0])
	assert.Nil(suite.T(), err)
}

// func (suite *UserUsecaseTestSuite) TestEditPhotoProfile_Success() {
// 	userUsecase := NewUserUsecase(suite.userRepoMock, suite.fileRepoMock)
// 	dummyFileExt := "jpg"
// 	dummyFileName := "user_Dummy Username 1.jpg"
// 	multipartFile := &multipart.FileHeader{
// 		Filename: "test-image.jpg",
// 		Size:     int64(len("test content")),
// 	}
// 	openedFile, err := multipartFile.Open()
// 	require.NoError(suite.T(), err)
// 	defer openedFile.Close()

// 	suite.fileRepoMock.On("Save", dummyFileName, &openedFile).Return(nil)
// 	suite.userRepoMock.On("UpdatePhotoProfile", dummyUsers[0].Username, dummyFileName).Return(nil)
// 	res := userUsecase.EditPhotoProfile(dummyUsers[0].Username, dummyFileExt, &openedFile)
// 	assert.Nil(suite.T(), res)
// }

func (suite *UserUsecaseTestSuite) TestUnregProfile_Success() {
	userUsecase := NewUserUsecase(suite.userRepoMock, suite.fileRepoMock)
	suite.userRepoMock.On("DeleteUserById", dummyUsers[0].Username).Return(nil)
	err := userUsecase.UnregProfile(dummyUsers[0].Username)
	assert.Nil(suite.T(), err)
}

func (suite *UserUsecaseTestSuite) TestUnregProfile_Failed() {
	userUsecase := NewUserUsecase(suite.userRepoMock, suite.fileRepoMock)
	suite.userRepoMock.On("DeleteUserById", dummyUsers[0].Username).Return(errors.New("Failed"))
	err := userUsecase.UnregProfile(dummyUsers[0].Username)
	assert.NotNil(suite.T(), err)
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.userRepoMock = new(userRepoMock)
	suite.fileRepoMock = new(fileRepoMock)
	suite.utilsMock = new(utilsMock)
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

// fileInfo, err := file.Stat()
// require.NoError(suite.T(), err)

// fileHeader := make([]byte, 512)
// _, err = file.Read(fileHeader)
// require.NoError(suite.T(), err)
// file.Seek(0, 0)

// multipartFile.Filename = file.Name()
// multipartFile.Size = fileInfo.Size()
