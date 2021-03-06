package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mfesenko/sf-movie-locations/persistence"
)

type MoviesController struct {
	datastore *persistence.DataStore
}

func NewMoviesController(datastore *persistence.DataStore) *MoviesController {
	return &MoviesController{datastore}
}

func (mc MoviesController) GetAllMovieLocations(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	title := request.URL.Query().Get("title")
	movies := mc.datastore.FindMovies(title)
	writeJsonResponse(&writer, movies)
}

func (mc MoviesController) GetMovies(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	title := request.URL.Query().Get("title")
	movies := mc.datastore.FindMovies(title, "title")
	writeJsonResponse(&writer, movies)
}

func writeJsonResponse(writer *http.ResponseWriter, response interface{}) {
	responseWriter := *writer
	responseJson, err := json.Marshal(response)
	if err == nil {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(200)
		fmt.Fprintf(responseWriter, "%s", responseJson)
	} else {
		responseWriter.WriteHeader(500)
		fmt.Fprintf(responseWriter, "Failed to serialize response data")
	}
}
