package models

import "gopkg.in/mgo.v2/bson"

type (
	Location struct {
		Latitude    float64 `json:"lat"`
		Longitude   float64 `json:"lng"`
		Description string  `json:"description"`
	}

	Movie struct {
		Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Title     string        `json:"title"`
		Year      string        `json:"year"`
		Director  string        `json:"director,omitempty"`
		Actors    []string      `json:"actors,omitempty"`
		Locations []Location    `json:"locations,omitempty"`
	}
)
