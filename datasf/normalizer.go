package datasf

import (
	"fmt"
	"log"

	"errors"
	"strings"

	"github.com/mfesenko/sf-movie-locations/models"
	"gopkg.in/mgo.v2/bson"
)

type Normalizer struct {
	geocodingService *GeocodingService
}

func NewNormalizer(geocodingService *GeocodingService) *Normalizer {
	return &Normalizer{geocodingService}
}

func (n Normalizer) Normalize(records []DataSFRecord) ([]models.Movie, error) {
	moviesMap := make(map[string]*models.Movie)
	var resultError error = nil
	for _, record := range records {
		if record.Locations == "" {
			log.Printf("Skipping record for movie '%s' because location was empty", record.Title)
			continue
		}
		movieKey := fmt.Sprintf("%s|%s", record.Title, record.Release_year)
		movie := moviesMap[movieKey]
		location, err := n.geocodingService.ConvertAddressToCoordinates(record.Locations + ", San Francisco, CA")
		if err != nil {
			log.Printf("Failed to parse location for movie %s with error: %s", record.Title, err)
			if strings.Contains(err.Error(), "OVER_QUERY_LIMIT") {
				resultError = errors.New("Query limit was exceeded")
				break
			}
		} else {
			if movie == nil {
				movie = &models.Movie{
					Id:        bson.NewObjectId(),
					Title:     record.Title,
					Year:      record.Release_year,
					Actors:    buildNotEmptySlice(record.Actor_1, record.Actor_2, record.Actor_3),
					Director:  record.Director,
					Locations: make([]models.Location, 0),
				}
				moviesMap[movieKey] = movie
			}
			location.Description = record.Locations
			movie.Locations = append(movie.Locations, location)
		}
	}
	moviesList := make([]models.Movie, len(moviesMap))
	index := 0
	for _, movie := range moviesMap {
		moviesList[index] = *movie
		index++
	}
	return moviesList, resultError
}

func buildNotEmptySlice(values ...string) []string {
	result := make([]string, 0)
	for _, value := range values {
		if len(value) > 0 {
			result = append(result, value)
		}
	}
	return result
}
