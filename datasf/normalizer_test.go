package datasf

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/mfesenko/sf-movie-locations/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedGeocodingService struct {
	mock.Mock
}

const validAddress = "valid address, San Francisco, CA"
const queryLimitAddress = "query limit address, San Francisco, CA"

var expectedLocation = models.Location{
	Latitude:  37.15,
	Longitude: -122.56,
}

func (gs *MockedGeocodingService) ConvertAddressToCoordinates(address string) (models.Location, error) {
	gs.Called(address)
	if address == validAddress {
		return expectedLocation, nil
	}
	if address == queryLimitAddress {
		return models.Location{}, errors.New(" OVER_QUERY_LIMIT ")
	}
	return models.Location{}, errors.New("invalid location")
}

func TestNewNormalizerWithNilGeocodingService(t *testing.T) {
	assert := assert.New(t)
	normalizer, err := NewNormalizer(nil)
	assert.Nil(normalizer)
	assert.EqualError(err, "Geocoding service is required")
}

func TestNewNormalizer(t *testing.T) {
	assert := assert.New(t)
	geocodingService := &MockedGeocodingService{}
	normalizer := createNormalizer(assert, geocodingService)
	assert.Equal(geocodingService, normalizer.geocodingService)
}

func TestNormalizer_Normalize(t *testing.T) {
	assert := assert.New(t)
	geocodingService := &MockedGeocodingService{}
	normalizer := createNormalizer(assert, geocodingService)

	geocodingService.On("ConvertAddressToCoordinates", "valid address, San Francisco, CA").Return(
		expectedLocation, nil)
	geocodingService.On("ConvertAddressToCoordinates", "invalid address, San Francisco, CA").Return(
		models.Location{}, errors.New("invalid locations"))
	records := getRecords(assert, "normalize_success_records.json")
	movies, err := normalizer.Normalize(records)
	assert.Nil(err)
	expectedMovies := getExpectedMovies(assert, "normalize_success_movies.json")
	validateMovies(assert, expectedMovies, movies)
	geocodingService.AssertExpectations(t)
}

func TestNormalizer_NormalizeWithQueryLimitError(t *testing.T) {
	assert := assert.New(t)
	geocodingService := &MockedGeocodingService{}
	normalizer := createNormalizer(assert, geocodingService)

	geocodingService.On("ConvertAddressToCoordinates", "valid address, San Francisco, CA").Return(
		expectedLocation, nil)
	geocodingService.On("ConvertAddressToCoordinates", "query limit address, San Francisco, CA").Return(
		models.Location{}, errors.New(" OVER_QUERY_LIMIT "))

	records := getRecords(assert, "normalize_query_limit_error_records.json")
	movies, err := normalizer.Normalize(records)
	assert.EqualError(err, "Query limit was exceeded")
	expectedMovies := getExpectedMovies(assert, "normalize_query_limit_error_movies.json")
	validateMovies(assert, expectedMovies, movies)
	geocodingService.AssertExpectations(t)
}

func createNormalizer(assert *assert.Assertions, geocodingService GeocodingService) *Normalizer {
	normalizer, err := NewNormalizer(geocodingService)
	assert.NotNil(normalizer)
	assert.Nil(err)
	return normalizer
}

func validateMovies(assert *assert.Assertions, expectedMovies []models.Movie, movies []models.Movie) {
	assert.Equal(len(expectedMovies), len(movies))
	for index, expectedMovie := range expectedMovies {
		movie := movies[index]
		assert.Equal(expectedMovie.Title, movie.Title)
		assert.Equal(expectedMovie.Year, movie.Year)
		assert.Equal(expectedMovie.Director, movie.Director)
		assert.Equal(expectedMovie.Actors, movie.Actors)
		assert.Equal(expectedMovie.Locations, movie.Locations)
	}
}

func getRecords(assert *assert.Assertions, fileName string) []DataSFRecord {
	data := readTestDataFile(assert, fileName)
	var records []DataSFRecord
	json.Unmarshal(data, &records)
	return records
}

func getExpectedMovies(assert *assert.Assertions, fileName string) []models.Movie {
	data := readTestDataFile(assert, fileName)
	var movies []models.Movie
	json.Unmarshal(data, &movies)
	return movies
}
