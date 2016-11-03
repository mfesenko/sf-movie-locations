# sf-movie-locations

## Problem
This project is a service that shows a map with locations where movies have been filmed in San Francisco. Original data is available on [DataSF](https://data.sfgov.org/Culture-and-Recreation/Film-Locations-in-San-Francisco/yitu-d5am).

I've decided to use Go for backend, MongoDB as storage, jQuery for frontend. Basically it's my first attempt to write something more that 'hello world' in Go. 

Live version deployed on Amazon EC2 instance can be found here: http://sfmovielocations.online/

## Achitecture
This projects consists of two main parts:
* webapp
* dataloader
 
**Webapp** is a simple single page application which shows on a map and allows user to filter locations on the by movie name. When user start typing movie name - list of movies with names that match the query will show up. User can either select a movie from autocomplete dropdown, either filter locations base on his input.
When a user clicks on location marked at the map - list of movies that were filmed there will be shown.
At the moment webapp is not big, it has only two api endpoints and static content.

**Dataloader** is a tool for fetching data from [original dataset](https://data.sfgov.org/Culture-and-Recreation/Film-Locations-in-San-Francisco/yitu-d5am). It loads data from SODA API, transforms it and saves to MongoDB.
This tool is also written using Go. 
Google Maps Web Service APIs is used to retrieve geographical coordinates by place description or address. For communication with this api client from Google is used. 
Since original dataset doesn't update too often - at the moment re-fetching of data is scheduled once a week.

## Build and execution
To run this application locally you will need:
* Go 1.7
* MongoDB 2.6+
To properly build the code you need to load all external dependencies. I've used [govendor](https://github.com/kardianos/govendor) for this purpose.
Also there's a **docker-compose** setup for both webapp and dataloader, using which you can avoid installing Go in your environment.

## Running webapp

Using docker-compose:
```
$ git clone git@github.com:mfesenko/sf-movie-locations.git
$ cd sf-movie-locations/docker
$ docker-compose up webapp
```

Building from scratch:
```
$ mkdir sf-movie-locations
$ cd sf-movie-locations
$ export GOPATH=`pwd`
$ go get github.com/kardianos/govendor
$ go build github.com/kardianos/govendor
$ go get github.com/mfesenko/sf-movie-locations
$ cd src/github.com/mfesenko/sf-movie-locations
$ $GOPATH/govendor init
$ $GOPATH/govendor fetch +out
$ go run webapp/main.go 
```

Database address and credentials, port to run server, path folder with static content are specified in config file. By default  **webapp-config.toml** file is used. 

Running with custom config file:
```
$ go run webapp/main.go -config custom-webaap-config.toml
```

## Running dataloader
Using docker-composer:
```
$ git clone git@github.com:mfesenko/sf-movie-locations.git
$ cd sf-movie-locations/docker
$ docker-compose up dataloader
```

Building from scratch:
```
$ mkdir sf-movie-locations
$ cd sf-movie-locations
$ export GOPATH=`pwd`
$ go get github.com/kardianos/govendor
$ go build github.com/kardianos/govendor
$ go get github.com/mfesenko/sf-movie-locations
$ cd src/github.com/mfesenko/sf-movie-locations
$ $GOPATH/govendor init
$ $GOPATH/govendor fetch +out
$ export GOOGLE_MAPS_API_KEY="<YOUR_GOOGLE_MAPS_API_KEY>"
$ go run dataloader/main.go 
```

Database address and credentials, link to api for retrieving of original dataset are specified in config file. By default  **dataloader-config.toml** file is used. 
Running with custom config file:
```
$ go run dataloader/main.go -config custom-dataloader-config.toml
```


## Persistence
Only one collection is used to store all information about movies. Name of database and collection can be configured using config file.
Example of movie document:
```
{
	"_id" : ObjectId("581a8be8cd8b220095fca335"),
	"title" : "A Night Full of Rain",
	"year" : "1978",
	"director" : "Lina Wertmuller",
	"actors" : [
		"Candice Bergen",
		"Giancarlo Gianni"
	],
	"locations" : [
		{
			"latitude" : 37.7489402,
			"longitude" : -122.3928481,
			"description" : "Embarcadero Freeway"
		}
	]
}
```

## Database initialization
Default settings assume that specific users are present in the database. 
User **dataloader** has write access to the database, user **webapp** hasa only read access.
There's sample script for creation of default users which can be executed like this:
```
$ mongo < db/createusers.js
```
 
## Testing
To execute tests you will need Go installed. Currently not all packages are completely covered with tests, but it will be fixed in the nearest future. 
For example, to run tests for package **config**:
 ```
 $ cd config
 $ go test
 ```
 
## TODO
There's always a room for improvements :) So here's the list of future improvements for this project:
1. Backend
* improve logging and error handling
* complete tests for all modules
* run MongoDB in a docker container
* convertation of location from original dataset to geographical coordinates is failing for some records
* currently dataloader is always re-creating database from scratch, some sort of incremental fetch can be implemented
2. Frontend
* fix issues with UI on mobile browsers an IE and do more testing on various environments
* think about grouping points on a map when there's a lot of them in some area and split them when user zooms in
* don't load all locations with one request, load them partially
3. General
* integrate with some other datasource to retrieve more information about movie, f.e. plot description, poster of the movie
* locations filtering not only by movie name
* currently dataloader is scheduled for execution in crontab once a week, can make it configurable option and control it from code  