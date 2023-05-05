package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IoSuite struct {
	suite.Suite
	fileBasePath string
	fileName     string
	fileContent  string
}

func (suite *IoSuite) SetupTest() {
	suite.fileBasePath = ""
	suite.fileName = "test-file"
	suite.fileContent = "test-content"
}

func (suite *IoSuite) CreateMultipartFile(fileContents []byte, fileName string) (multipart.File, *multipart.FileHeader, error) {
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

func (suite *IoSuite) TestSaveToLocalFile() {
	// Create a temporary file for testing purposes
	tempFile, err := ioutil.TempFile(suite.fileBasePath, suite.fileName)
	if err != nil {
		suite.T().Fatalf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Create a sample multipart.File for testing purposes
	fileContent := []byte(suite.fileContent)
	file, _, _ := suite.CreateMultipartFile(fileContent, suite.fileName)

	// Call SaveToLocalFile with the temporary file and sample multipart.File
	err = SaveToLocalFile(tempFile.Name(), &file)

	// Assert that there were no errors
	assert.NoError(suite.T(), err)

	// Read the saved file to verify its contents
	savedFileContent, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		suite.T().Fatalf("error reading saved file: %v", err)
	}

	// Assert that the saved file contents match the sample file contents
	assert.Equal(suite.T(), suite.fileContent, string(savedFileContent))
}

func TestRunIoSuite(t *testing.T) {
	suite.Run(t, new(IoSuite))
}
