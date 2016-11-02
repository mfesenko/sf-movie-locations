package datasf

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stretchr/testify/assert"
)

const testDataDir = "testdata/"

func readTestDataFile(assert *assert.Assertions, fileName string) []byte {
	filePath := testDataDir + fileName
	_, err := os.Stat(filePath)
	assert.Nil(err, fmt.Errorf("Test data file '%s' doesn't exist", filePath))
	data, err := ioutil.ReadFile(filePath)
	assert.Nil(err, fmt.Errorf("Failed to read test data file '%s', error: %s", filePath, err))
	return data
}
