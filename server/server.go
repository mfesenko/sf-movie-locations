package server

import (
	"fmt"
	"log"
	"net/http"

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
	router.GET("/api/movies/", moviesController.GetMovies)
	router.ServeFiles("/static/*filepath", http.Dir(s.config.Server.StaticContentPath))
	s.serveFile("/", "/index.html", router)
	s.serveFile("/favicon.ico", "/images/favicon.ico", router)
	log.Printf("Start listening at port %d...", s.config.Server.Port)
	address := fmt.Sprintf(":%d", s.config.Server.Port)
	http.ListenAndServe(address, router)
}

func (s Server) serveFile(path string, filePath string, router *httprouter.Router) {
	staticPath := s.config.Server.StaticContentPath + filePath
	router.GET(path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.ServeFile(writer, request, staticPath)
	})
}
