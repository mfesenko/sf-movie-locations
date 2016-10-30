package datasf

import (
	"fmt"
	"log"

	"github.com/mfesenko/sf-movie-locations/persistence"
	"gopkg.in/mgo.v2/bson"
)

type Normalizer struct {
	geocodingService *GeocodingService
}

func NewNormalizer(geocodingService *GeocodingService) *Normalizer {
	return &Normalizer{geocodingService}
}

func (n Normalizer) Normalize(records []DataSFRecord) []persistence.Movie {
	recordsCount := len(records)
	moviesMap := make(map[string]*persistence.Movie)
	for i := 0; i < recordsCount; i++ {
		record := records[i]
		if record.Locations == "" {
			log.Printf("Skipping record for movie '%s' because location was empty", record.Title)
			continue
		}

		movieKey := fmt.Sprintf("%s|%s", record.Title, record.Release_year)
		movie := moviesMap[movieKey]
		location := n.geocodingService.ConvertAddressToCoordinates(record.Locations + ", San Francisco, CA")
		if movie == nil {
			movie = &persistence.Movie{
				Id:        bson.NewObjectId(),
				Title:     record.Title,
				Year:      record.Release_year,
				Actors:    buildNotEmptySlice(record.Actor_1, record.Actor_2, record.Actor_3),
				Director:  record.Director,
				Locations: make([]persistence.Location, 0),
			}
			moviesMap[movieKey] = movie
		}
		if location != nil {
			location.Description = record.Locations
			movie.Locations = append(movie.Locations, *location)
		} else {
			log.Printf("Failed to parse location for movie %s - %s", record.Title, record.Locations)
		}
	}
	moviesList := make([]persistence.Movie, len(moviesMap))
	index := 0
	for _, movie := range moviesMap {
		moviesList[index] = *movie
		index++
	}
	return moviesList
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
