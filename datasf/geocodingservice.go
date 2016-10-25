package datasf

import (
	"log"

	"github.com/mfesenko/sf-movie-locations/persistence"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type GeocodingService struct {
	client           *maps.Client
	coordinatesCache map[string]*persistence.Location
}

func NewGeocodingService(apiKey string) *GeocodingService {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create Google Maps client: %s", err)
	}
	return &GeocodingService{client: client, coordinatesCache: make(map[string]*persistence.Location)}
}

func (gs *GeocodingService) ConvertAddressToCoordinates(address string) *persistence.Location {
	coordinates := gs.coordinatesCache[address]
	if coordinates == nil {
		coordinates = gs.requestCoordinates(address)
		gs.coordinatesCache[address] = coordinates
	}
	return coordinates
}

func (gs GeocodingService) requestCoordinates(address string) *persistence.Location {
	placeId := gs.getPlaceId(address)
	var response []maps.GeocodingResult
	var err error
	if placeId == nil {
		response, err = gs.getAddressCoordinates(address)

	} else {
		response, err = gs.getPlaceCoordinates(*placeId)
	}

	if err != nil {
		log.Printf("Failed to get geo coordinates for address '%s' with error: %s", address, err)
		return nil
	}
	location := response[0].Geometry.Location
	return &persistence.Location{location.Lat, location.Lng}
}

func (gs GeocodingService) getPlaceId(address string) *string {
	request := &maps.PlaceAutocompleteRequest{Input: address}
	response, err := gs.client.PlaceAutocomplete(context.Background(), request)
	if err != nil {
		return nil
	}
	return &response.Predictions[0].PlaceID
}

func (gs GeocodingService) getPlaceCoordinates(placeId string) ([]maps.GeocodingResult, error) {
	request := &maps.GeocodingRequest{PlaceID: placeId}
	return gs.client.ReverseGeocode(context.Background(), request)
}

func (gs GeocodingService) getAddressCoordinates(address string) ([]maps.GeocodingResult, error) {
	request := &maps.GeocodingRequest{Address: address}
	return gs.client.Geocode(context.Background(), request)
}
