package datasf

import (
	"fmt"
	"github.com/mfesenko/sf-movie-locations/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGeocodingServiceWithEmptyApiKey(t *testing.T) {
	service, err := NewGeocodingService("")
	assert := assert.New(t)
	assert.Nil(service)
	assert.EqualError(err, "datasf: Failed to create Google Maps client: maps: API Key or Maps for Work credentials missing")
}

func TestNewGeocodingService(t *testing.T) {
	assert := assert.New(t)
	service := createGeocodingService(assert)
	assert.NotNil(service.client)
}

func TestGeocodingService_ConvertAddressToCoordinatesForPlace(t *testing.T) {
	assert := assert.New(t)
	service := createGeocodingService(assert)

	address := "Golden Gate Bridge, San Francisco, CA"
	expectedLocation := models.Location{
		Latitude:  37.8199286,
		Longitude: -122.4782551,
	}

	location, err := service.ConvertAddressToCoordinates(address)
	assert.Nil(err)
	assert.Equal(expectedLocation, location)
	assert.Equal(1, len(service.coordinatesCache))

	location, err = service.ConvertAddressToCoordinates(address)
	assert.Nil(err)
	assert.Equal(expectedLocation, location)
	assert.Equal(1, len(service.coordinatesCache))
}

func TestGeocodingService_ConvertAddressToCoordinatesForAddress(t *testing.T) {
	assert := assert.New(t)
	service := createGeocodingService(assert)

	address := "Cesar Chavez & Mission Street (Mission), San Francisco, CA"
	expectedLocation := models.Location{
		Latitude:  37.7482293,
		Longitude: -122.4182139,
	}
	location, err := service.ConvertAddressToCoordinates(address)
	assert.Nil(err)
	assert.Equal(expectedLocation, location)
	assert.Equal(1, len(service.coordinatesCache))
}

func TestGeocodingService_ConvertAddressToCoordinatesForInvalidAddress(t *testing.T) {
	assert := assert.New(t)
	service := createGeocodingService(assert)

	address := "Pier 50- end of the pier, San Francisco, CA"
	expectedErrorMessage := fmt.Sprintf("Failed to get geo coordinates for address '%s' with error: maps: ZERO_RESULTS - ",
		address)
	expectedLocation := models.Location{}

	location, err := service.ConvertAddressToCoordinates(address)
	assert.EqualError(err, expectedErrorMessage)
	assert.Equal(expectedLocation, location)
	assert.Equal(0, len(service.coordinatesCache))
}

func createGeocodingService(assert *assert.Assertions) *GeocodingService {
	service, err := NewGeocodingService("AIzaSyBDII7q46NRV97keR-R14tbS_N7xgsLzE8")
	assert.NotNil(service)
	assert.Nil(err)
	return service
}
