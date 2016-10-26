package server

import (
	"fmt"
	"log"
	"net/http"

	"io/ioutil"

	"github.com/julienschmidt/httprouter"
	"github.com/mfesenko/sf-movie-locations/persistence"
)

type Server struct {
	config Config
}

func NewServer(configFilePath string) *Server {
	config := LoadConfigFile(configFilePath)
	return &Server{config}
}

func (s Server) Serve() {
	router := httprouter.New()
	log.Print("Initializing data store...")
	dataStore := persistence.NewDataStore(s.config.Db.Host, s.config.Db.DbName, s.config.Db.CollectionName)
	moviesController := NewMoviesController(dataStore)
	router.GET("/api/movieLocations", moviesController.GetAllMovieLocations)
	router.GET("/api/movieLocations/:id", moviesController.GetMovieLocations)
	router.GET("/api/movies/:title", moviesController.GetMovies)
	router.ServeFiles("/static/*filepath", http.Dir(s.config.Server.StaticContentPath))
	router.GET("/", s.GetIndex)
	log.Printf("Start listening at port %d...", s.config.Server.Port)
	address := fmt.Sprintf("localhost:%d", s.config.Server.Port)
	http.ListenAndServe(address, router)
}

func (s Server) GetIndex(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	path := s.config.Server.StaticContentPath + "/index.html"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		writer.WriteHeader(404)
	} else {
		writer.WriteHeader(200)
		writer.Write(content)
	}
}
