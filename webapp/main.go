package main

import (
	"flag"

	"github.com/mfesenko/sf-movie-locations/server"
)

func main() {
	configFilePath := flag.String("config", "config.toml", "path to config file")
	flag.Parse()
	server := server.NewServer(*configFilePath)
	server.Serve()
}
