package main

import (
	"flag"

	"log"

	"github.com/mfesenko/sf-movie-locations/server"
)

func main() {
	configFilePath := flag.String("config", "webapp-config.toml", "path to config file")
	flag.Parse()
	server, err := server.NewServer(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	server.Serve()
}
