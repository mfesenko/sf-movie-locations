package main

import (
	"net/http"

	"log"

	"fmt"

	"github.com/julienschmidt/httprouter"
	"github.com/mfesenko/sf-movie-locations/controllers"
	"github.com/mfesenko/sf-movie-locations/persistence"
)

const (
	DB_HOST            = "localhost"
	DB_NAME            = "sf-movie-locations"
	DB_COLLECTION_NAME = "movies"
	PORT               = 2222
)

func main() {
	router := httprouter.New()
	log.Print("Initializing data store...")
	dataStore := persistence.NewDataStore(DB_HOST, DB_NAME, DB_COLLECTION_NAME)
	moviesController := controllers.NewMoviesController(dataStore)
	router.GET("/api/movieLocations", moviesController.GetAllMovieLocations)
	router.GET("/api/movieLocations/:id", moviesController.GetMovieLocations)
	router.GET("/api/movies/:title", moviesController.GetMovies)
	log.Printf("Start listening at port %d...", PORT)
	address := fmt.Sprintf("localhost:%d", PORT)
	http.ListenAndServe(address, router)
}
