package persistence

import (
	"testing"

	"github.com/facebookgo/mgotest"
	"github.com/mfesenko/sf-movie-locations/config"
	"github.com/mfesenko/sf-movie-locations/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const testDbName = "testdb"
const testCollectionName = "testcollection"

type GetMoviesTest struct {
	title          string
	expectedMovies []models.Movie
}

func TestNewDataStoreFailure(t *testing.T) {
	dbConfig := createDbConfig("test:6666")
	datastore, err := NewDataStore(dbConfig)
	assert := assert.New(t)
	assert.EqualError(err, "Connection to database failed: no reachable servers")
	assert.Nil(datastore)
}

func TestNewDataStoreSuccess(t *testing.T) {
	mongo := mgotest.NewStartedServer(t)
	defer mongo.Stop()

	assert := assert.New(t)
	datastore := createDataStore(mongo, assert)
	assert.NotNil(datastore.session)
	assert.Equal(testDbName, datastore.dbName)
	assert.Equal(testCollectionName, datastore.collectionName)
}

func TestDataStore_FindMovies(t *testing.T) {
	testFindMovies(t, nil)
}

func TestDataStore_FindMoviesWithTitleOnly(t *testing.T) {
	testFindMovies(t, func(movie models.Movie) models.Movie {
		return models.Movie{
			Id:    movie.Id,
			Title: movie.Title,
		}
	}, "title")
}

func testFindMovies(t *testing.T, transformer func(models.Movie) models.Movie, fields ...string) {
	mongo := mgotest.NewStartedServer(t)
	defer mongo.Stop()

	movies := initDbCollection(mongo)
	assert := assert.New(t)
	datastore := createDataStore(mongo, assert)

	var expectedMovies []models.Movie
	if transformer == nil {
		expectedMovies = movies
	} else {
		expectedMovies = make([]models.Movie, len(movies))
		for index, movie := range movies {
			expectedMovies[index] = transformer(movie)
		}
	}

	var tests = []GetMoviesTest{
		{"test", []models.Movie{}},
		{"wonDer", expectedMovies[1:2]},
		{"Woman", expectedMovies},
		{"", expectedMovies},
	}

	for _, test := range tests {
		resultMovies := datastore.FindMovies(test.title, fields...)
		assert.Equal(test.expectedMovies, resultMovies)
	}
}

func TestDataStore_AddMovie(t *testing.T) {
	mongo := mgotest.NewStartedServer(t)
	session := mongo.Session().Copy()
	defer mongo.Stop()
	defer session.Close()

	assert := assert.New(t)
	datastore := createDataStore(mongo, assert)

	validateEntities(session, assert)
	movie := getTestMovies()[0]
	datastore.AddMovie(movie)
	validateEntities(session, assert, movie)
}

func TestDataStore_AddMovies(t *testing.T) {
	mongo := mgotest.NewStartedServer(t)
	session := mongo.Session().Copy()
	defer mongo.Stop()
	defer session.Close()

	assert := assert.New(t)
	datastore := createDataStore(mongo, assert)

	validateEntities(session, assert)
	movies := getTestMovies()
	datastore.AddMovies(movies)
	validateEntities(session, assert, movies...)
}

func TestDataStore_Reset(t *testing.T) {
	mongo := mgotest.NewStartedServer(t)
	session := mongo.Session().Copy()
	defer mongo.Stop()
	defer session.Close()

	movies := initDbCollection(mongo)
	assert := assert.New(t)
	validateEntities(session, assert, movies...)
	datastore := createDataStore(mongo, assert)
	datastore.Reset()
	validateEntities(session, assert)
}

func createDbConfig(url string) config.DbConfig {
	return config.DbConfig{
		Host:           url,
		DbName:         testDbName,
		CollectionName: testCollectionName,
		Timeout:        3,
	}
}

func createDataStore(mongo *mgotest.Server, assert *assert.Assertions) *DataStore {
	dbConfig := createDbConfig(mongo.URL())
	datastore, err := NewDataStore(dbConfig)
	assert.Nil(err)
	assert.NotNil(datastore)
	return datastore
}

func initDbCollection(mongo *mgotest.Server) []models.Movie {
	session := mongo.Session().Copy()
	defer session.Close()
	movies := getTestMovies()
	for _, movie := range movies {
		session.DB(testDbName).C(testCollectionName).Insert(movie)
	}
	return movies
}

func getTestMovies() []models.Movie {
	return []models.Movie{
		{
			Id:       bson.NewObjectId(),
			Title:    "Pretty woman",
			Year:     "1990",
			Director: "Garry Marshall",
			Actors:   []string{"Richard Gere", "Julia Roberts", "Jason Alexander"},
			Locations: []models.Location{
				{
					Latitude:    22,
					Longitude:   45,
					Description: "location1",
				},
			},
		},
		{
			Id:       bson.NewObjectId(),
			Title:    "Wonder woman",
			Year:     "2017",
			Director: "Patty Jenkins",
			Actors:   []string{"Gal Gadot", "Chris Pine", "Robin Wright"},
			Locations: []models.Location{
				{
					Latitude:    122,
					Longitude:   145,
					Description: "location2",
				},
			},
		},
	}
}

func validateEntities(session *mgo.Session, assert *assert.Assertions, expectedMovies ...models.Movie) {
	var movies []models.Movie
	err := session.DB(testDbName).C(testCollectionName).Find(bson.M{}).All(&movies)
	assert.Nil(err)
	assert.Equal(expectedMovies, movies)
}
