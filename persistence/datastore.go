package persistence

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DataStore struct {
	session        *mgo.Session
	dbName         string
	collectionName string
}

func NewDataStore(dbHost string, dbName string, collectionName string) *DataStore {
	session, err := mgo.Dial("mongodb://" + dbHost)
	if err != nil {
		log.Fatal("Connection to database failed: %s", err)
	}
	return &DataStore{session, dbName, collectionName}
}

func (ds DataStore) collection() *mgo.Collection {
	return ds.session.DB(ds.dbName).C(ds.collectionName)
}

func (ds DataStore) GetAllMovies() []Movie {
	var movies []Movie
	ds.collection().Find(bson.M{}).All(&movies)
	return movies
}

func (ds DataStore) GetMovie(id string) *Movie {
	if bson.IsObjectIdHex(id) {
		var movie Movie
		err := ds.collection().FindId(bson.ObjectIdHex(id)).One(&movie)
		if err == nil {
			return &movie
		}
	}
	return nil
}

func (ds DataStore) AddMovies(movies []Movie) {
	entities := make([]interface{}, len(movies))
	for index := range movies {
		entities[index] = movies[index]
		index++
	}
	err := ds.collection().Insert(entities...)
	if err != nil {
		log.Fatal(err)
	}
}

func (ds DataStore) AddMovie(movie Movie) {
	err := ds.collection().Insert(movie)
	if err != nil {
		log.Fatal(err)
	}
}

func (ds DataStore) Reset() {
	ds.collection().DropCollection()
	index := mgo.Index{Key: []string{"title"}}
	ds.collection().EnsureIndex(index)
}

func (ds DataStore) FindMovies(title string) []Movie {
	query := bson.M{
		"title": bson.RegEx{title, "i"},
	}
	selector := bson.M{
		"title": 1,
		"year":  1,
	}
	movies := make([]Movie, 0)
	ds.collection().Find(query).Select(selector).All(&movies)
	return movies
}
