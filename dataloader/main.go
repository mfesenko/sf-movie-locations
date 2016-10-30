package main

import (
	"log"
	"os"

	"github.com/mfesenko/sf-movie-locations/datasf"
	"github.com/mfesenko/sf-movie-locations/persistence"
)

var API_KEY = os.Getenv("GOOGLE_MAPS_API_KEY")

func main() {
	datastore := persistence.NewDataStore("localhost", "sf-movie-locations", "movies")
	geocodingService, err := datasf.NewGeocodingService(API_KEY)
	if err != nil {
		log.Fatal(err)
	}
	dataSFService := datasf.NewDataSFService()
	normalizer := datasf.NewNormalizer(geocodingService)

	limit := 1000
	offset := 0

	movies := make([]persistence.Movie, 0)
	recordsCount := 0
	continueRetrieve := true
	for continueRetrieve {
		log.Printf("Retrieving next %d record(s) starting from %d", limit, offset)
		records, err := dataSFService.RetrieveRecords(offset, limit)
		if err != nil {
			log.Fatal(err)
		}
		retrievedRecordsCount := len(records)
		recordsCount += retrievedRecordsCount
		continueRetrieve = retrievedRecordsCount > 0
		offset += limit

		normalizedMovies, err := normalizer.Normalize(records)
		movies = append(movies, normalizedMovies...)
		if err != nil {
			log.Print(err)
			continueRetrieve = false
		}
	}

	log.Printf("Read %d record(s) from datasf", recordsCount)
	log.Printf("Saving %d movie record(s) to db", len(movies))

	datastore.Reset()
	datastore.AddMovies(movies)
}
