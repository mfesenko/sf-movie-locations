package datasf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const requestOffset = 10
const requestLimit = 5

func TestNewDataSFServiceSuccess(t *testing.T) {
	validUrls := []string{
		"http://example.com",
		"https://data.sfgov.org/resource/wwmu-gmzc.json",
	}
	assert := assert.New(t)
	for _, url := range validUrls {
		service, err := NewDataSFService(url)
		assert.Nil(err)
		assert.NotNil(service)
		assert.Equal(url, service.baseUrl)
	}
}

func TestNewDataSFServiceError(t *testing.T) {
	invalidUrls := []string{
		"",
		"asdf",
		"asdf.com",
	}
	assert := assert.New(t)
	for _, url := range invalidUrls {
		service, err := NewDataSFService(url)
		assert.Nil(service)
		assert.EqualError(err, fmt.Sprintf("Invalid base url: %s", url))
	}
}

func TestDataSFService_RetrieveRecordsSuccess(t *testing.T) {
	assert := assert.New(t)
	var expectedRecords []DataSFRecord
	recordsData := readTestDataFile(assert, "datasfservice_records.json")
	err := json.Unmarshal(recordsData, &expectedRecords)
	assert.Nil(err)

	responseData := readTestDataFile(assert, "datasfservice_response.json")
	server := httptest.NewServer(createSuccessfulHandler(string(responseData)))
	defer server.Close()

	service := createDataSFService(assert, server.URL)
	records, err := service.RetrieveRecords(requestOffset, requestLimit)
	assert.Nil(err)
	assert.Equal(expectedRecords, records)
}

func TestDataSFService_RetrieveRecordsRequestError(t *testing.T) {
	assert := assert.New(t)
	url := "http://127.0.0.1:12345"
	service := createDataSFService(assert, url)
	validateRequestFailure(assert, url, service)

}

func TestDataSFService_RetrieveRecordsRequestFailure(t *testing.T) {
	handler := func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(500)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	assert := assert.New(t)
	service := createDataSFService(assert, server.URL)
	validateRequestFailure(assert, server.URL, service)
}

func TestDataSFService_RetrieveRecordsNotJsonContent(t *testing.T) {
	tests := []string{"", "test"}
	assert := assert.New(t)
	for _, test := range tests {
		testWithNotJsonContent(assert, test)
	}
}

func testWithNotJsonContent(assert *assert.Assertions, content string) {
	server := httptest.NewServer(createSuccessfulHandler(content))
	defer server.Close()

	service := createDataSFService(assert, server.URL)

	records, err := service.RetrieveRecords(requestOffset, requestLimit)
	assert.EqualError(err, "datasf: Failed to deserialize records from json response")
	assert.Equal(0, len(records))
}

func createSuccessfulHandler(content string) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if isRequestParamValid(request, "$offset", requestOffset) &&
			isRequestParamValid(request, "$limit", requestLimit) {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(200)
			fmt.Fprint(writer, content)
		} else {
			writer.WriteHeader(404)
		}
	})
}

func isRequestParamValid(request *http.Request, paramName string, expectedValue int) bool {
	paramValue := request.URL.Query().Get(paramName)
	return paramValue == fmt.Sprintf("%d", expectedValue)
}

func createDataSFService(assert *assert.Assertions, url string) *DataSFService {
	service, err := NewDataSFService(url)
	assert.Nil(err)
	assert.NotNil(service)
	return service
}

func validateRequestFailure(assert *assert.Assertions, url string, service *DataSFService) {
	records, err := service.RetrieveRecords(requestOffset, requestLimit)
	expectedError := fmt.Sprintf("datasf: Request to %s was not successful", url)
	assert.EqualError(err, expectedError)
	assert.Equal(0, len(records))
}
