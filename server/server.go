package server

import (
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/mfesenko/sf-movie-locations/config"
	"github.com/mfesenko/sf-movie-locations/persistence"
)

type Server struct {
	config config.Config
}

func NewServer(configFilePath string) (*Server, error) {
	config, err := config.LoadConfigFile(configFilePath)
	if err != nil {
		return nil, err
	}
	return &Server{config}, nil
}

func (s Server) Serve() {
	router := httprouter.New()
	log.Print("Initializing data store...")
	dataStore, err := persistence.NewDataStore(s.config.Db)
	if err != nil {
		log.Fatal(err)
	}
	moviesController := NewMoviesController(dataStore)
	router.GET("/api/movieLocations", moviesController.GetAllMovieLocations)
	router.GET("/api/movies", moviesController.GetMovies)
	router.GET("/static/*filepath", s.serveStatics)
	s.serveFile("/", "/index.html", router)
	s.serveFile("/favicon.ico", "/images/favicon.ico", router)
	log.Printf("Start listening at port %d...", s.config.Server.Port)
	address := fmt.Sprintf(":%d", s.config.Server.Port)
	http.ListenAndServe(address, router)
}

func (s Server) serveStatics(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	filepath := params.ByName("filepath")
	if filepath != "/index.html" {
		fileSystemPath := s.config.Server.StaticContentPath + filepath
		file, err := os.Stat(fileSystemPath)
		if err == nil && !file.IsDir() {
			http.ServeFile(writer, request, fileSystemPath)
			return
		}
	}
	http.NotFound(writer, request)
}

func (s Server) serveFile(path string, filePath string, router *httprouter.Router) {
	staticPath := s.config.Server.StaticContentPath + filePath
	router.GET(path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.ServeFile(writer, request, staticPath)
	})
}
