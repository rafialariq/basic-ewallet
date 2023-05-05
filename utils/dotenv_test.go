package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DotEnvSuite struct {
	suite.Suite
	Key string
}

func (suite *DotEnvSuite) SetupTest() {
	suite.Key = "DUMMY_KEY"
}

func (suite *DotEnvSuite) CreateEnvFile() {
	file, err := os.Create(".env.test")
	if err != nil {
		suite.T().Fatalf("Failed to create test file: %v", err)
	}

	if _, err := file.Write([]byte("DUMMY_KEY=dummy_value\n")); err != nil {
		suite.T().Fatalf("Failed to write to test file: %v", err)
	}
	file.Close()
}

func (suite *DotEnvSuite) TestDotEnv_Success() {
	suite.CreateEnvFile()
	expectedValue := "dummy_value"
	actualValue := DotEnv("DUMMY_KEY", "../utils/.env.test")
	assert.Equal(suite.T(), expectedValue, actualValue)

	defer func() {
		if err := os.Remove(".env.test"); err != nil {
			if !os.IsNotExist(err) {
				suite.T().Errorf("error deleting .env.test file: %v", err)
			}
		}
	}()
}

func TestRunDotEnvSuite(t *testing.T) {
	suite.Run(t, new(DotEnvSuite))
}
