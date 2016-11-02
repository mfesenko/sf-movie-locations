package main

import (
	"log"
	"os"

	"flag"

	"github.com/mfesenko/sf-movie-locations/config"
	"github.com/mfesenko/sf-movie-locations/datasf"
	"github.com/mfesenko/sf-movie-locations/models"
	"github.com/mfesenko/sf-movie-locations/persistence"
)

var API_KEY = os.Getenv("GOOGLE_MAPS_API_KEY")

func main() {
	configFilePath := flag.String("config", "dataloader-config.toml", "path to config file")
	config, err := config.LoadConfigFile(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	datastore, err := persistence.NewDataStore(config.Db)
	if err != nil {
		log.Fatal(err)
	}
	geocodingService, err := datasf.NewGeocodingService(API_KEY)
	if err != nil {
		log.Fatal(err)
	}
	dataSFService, err := datasf.NewDataSFService(config.Dataloader.BaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	normalizer, err := datasf.NewNormalizer(geocodingService)
	if err != nil {
		log.Fatal(err)
	}

	limit := config.Dataloader.BatchSize
	offset := 0

	movies := make([]models.Movie, 0)
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
		if len(normalizedMovies) > 0 {
			movies = append(movies, normalizedMovies...)
		}
		if err != nil {
			log.Print(err)
			continueRetrieve = false
		}
	}

	log.Printf("Read %d record(s) from datasf", recordsCount)
	log.Printf("Saving %d movie record(s) to db", len(movies))

	if len(movies) > 0 {
		datastore.Reset()
		datastore.AddMovies(movies)
	}
}
