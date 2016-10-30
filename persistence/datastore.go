package persistence

import (
	"log"

	"fmt"
	"time"

	"github.com/mfesenko/sf-movie-locations/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DataStore struct {
	session        *mgo.Session
	dbName         string
	collectionName string
}

func NewDataStore(conf config.DbConfig) (*DataStore, error) {
	dialInfo := mgo.DialInfo{
		Addrs:    []string{conf.Host},
		Database: conf.DbName,
		Username: conf.Username,
		Password: conf.Password,
		Timeout:  time.Duration(conf.Timeout) * time.Second,
	}
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		return nil, fmt.Errorf("Connection to database failed: %s", err)
	}
	return &DataStore{session, conf.DbName, conf.CollectionName}, nil
}

func (ds DataStore) getSession() *mgo.Session {
	return ds.session.Copy()
}

func (ds DataStore) collection(session *mgo.Session) *mgo.Collection {
	return session.DB(ds.dbName).C(ds.collectionName)
}

func (ds DataStore) GetMovies(title string) []Movie {
	session := ds.getSession()
	defer session.Close()

	var movies []Movie
	query := bson.M{}
	if title != "" {
		query = ds.createTitleQuery(title)
	}
	ds.collection(session).Find(query).All(&movies)
	return movies
}

func (ds DataStore) AddMovies(movies []Movie) {
	session := ds.getSession()
	defer session.Close()

	entities := make([]interface{}, len(movies))
	for index := range movies {
		entities[index] = movies[index]
		index++
	}
	err := ds.collection(session).Insert(entities...)
	if err != nil {
		log.Fatal(err)
	}
}

func (ds DataStore) AddMovie(movie Movie) {
	session := ds.getSession()
	defer session.Close()
	err := ds.collection(session).Insert(movie)
	if err != nil {
		log.Fatal(err)
	}
}

func (ds DataStore) Reset() {
	session := ds.getSession()
	defer session.Close()

	ds.collection(session).DropCollection()
	index := mgo.Index{Key: []string{"title"}}
	ds.collection(session).EnsureIndex(index)
}

func (ds DataStore) FindMovies(title string) []Movie {
	session := ds.getSession()
	defer session.Close()

	query := ds.createTitleQuery(title)
	selector := bson.M{
		"title": 1,
	}
	movies := make([]Movie, 0)
	ds.collection(session).Find(query).Select(selector).All(&movies)
	return movies
}

func (ds DataStore) createTitleQuery(title string) bson.M {
	return bson.M{
		"title": bson.RegEx{title, "i"},
	}
}
