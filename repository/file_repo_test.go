package repository

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type utilsMock struct {
	mock.Mock
}

func (u *utilsMock) SaveToLocalFile(fileLocation string, file *multipart.File) error {
	args := u.Called(fileLocation, file)
	if args != nil {
		return args.Error(0)
	}
	return nil
}

type FileRepositoryTestSuite struct {
	utilsMock *utilsMock
	suite.Suite
}

func (suite *FileRepositoryTestSuite) TestSave_Success() {
	tempDir, err := ioutil.TempDir("", "file-repo-test")
	suite.Require().NoError(err)

	defer func() {
		err := os.RemoveAll(tempDir)
		suite.Require().NoError(err)
	}()

	dummyFileName := "Dummy File Name"
	dummyFileLocation := filepath.Join(tempDir, dummyFileName)
	fileContent := []byte("file content")
	file, _, err := createMultipartFile(fileContent, dummyFileName)
	suite.Require().NoError(err)

	repo := NewFileRepository(tempDir)
	suite.utilsMock.On("SaveToLocalFile", dummyFileLocation, file).Return(nil)
	actual, err := repo.Save(dummyFileName, &file)
	assert.Equal(suite.T(), dummyFileLocation, actual)
	assert.Nil(suite.T(), err)
}

func (suite *FileRepositoryTestSuite) TestSave_Failed() {
	dummyFilePath := "Dummy File Path"
	dummyFileName := "Dummy File Name"
	dummyFileLocation := filepath.Join(dummyFilePath, dummyFileName)
	fileContent := []byte("file content")
	file, _, err := createMultipartFile(fileContent, dummyFileName)
	suite.Require().NoError(err)

	repo := NewFileRepository(dummyFilePath)
	suite.utilsMock.On("SaveToLocalFile", dummyFileLocation, file).Return(errors.New("Failed"))
	actual, err := repo.Save(dummyFileName, &file)
	assert.Equal(suite.T(), "", actual)
	assert.NotNil(suite.T(), err)
}

func (suite *FileRepositoryTestSuite) SetupTest() {
	suite.utilsMock = new(utilsMock)
}

func TestFileRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(FileRepositoryTestSuite))
}

func createMultipartFile(fileContents []byte, fileName string) (multipart.File, *multipart.FileHeader, error) {
	file := bytes.NewReader(fileContents)
	fileHeader := &multipart.FileHeader{
		Filename: fileName,
	}
	formData := &bytes.Buffer{}
	writer := multipart.NewWriter(formData)
	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}
	writer.Close()

	fileBytes := bytes.NewReader(formData.Bytes())
	request := &http.Request{
		Method: "POST",
		URL:    &url.URL{},
		Header: map[string][]string{
			"Content-Type": {writer.FormDataContentType()},
		},
		Body: ioutil.NopCloser(fileBytes),
	}
	return request.FormFile("file")
}
