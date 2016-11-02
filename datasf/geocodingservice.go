package datasf

import (
	"fmt"

	"github.com/mfesenko/sf-movie-locations/models"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type GeocodingService interface {
	ConvertAddressToCoordinates(address string) (models.Location, error)
}

type googleGeocodingService struct {
	client           *maps.Client
	coordinatesCache map[string]models.Location
}

func NewGeocodingService(apiKey string) (GeocodingService, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("datasf: Failed to create Google Maps client: %s", err)
	}
	return &googleGeocodingService{
		client:           client,
		coordinatesCache: make(map[string]models.Location),
	}, nil
}

func (gs *googleGeocodingService) ConvertAddressToCoordinates(address string) (models.Location, error) {
	coordinates, ok := gs.coordinatesCache[address]
	if ok == false {
		var err error
		coordinates, err = gs.requestCoordinates(address)
		if err != nil {
			return models.Location{}, err
		}
		gs.coordinatesCache[address] = coordinates
	}
	return coordinates, nil
}

func (gs googleGeocodingService) requestCoordinates(address string) (models.Location, error) {
	placeId, err := gs.getPlaceId(address)
	var response []maps.GeocodingResult
	if err != nil {
		response, err = gs.getAddressCoordinates(address)

	} else {
		response, err = gs.getPlaceCoordinates(placeId)
	}

	if err != nil {
		return models.Location{},
			fmt.Errorf("Failed to get geo coordinates for address '%s' with error: %s", address, err)
	}
	location := response[0].Geometry.Location
	return models.Location{
		Latitude:  location.Lat,
		Longitude: location.Lng,
	}, nil
}

func (gs googleGeocodingService) getPlaceId(address string) (string, error) {
	request := &maps.PlaceAutocompleteRequest{Input: address}
	response, err := gs.client.PlaceAutocomplete(context.Background(), request)
	if err != nil {
		return "", err
	}
	return response.Predictions[0].PlaceID, nil
}

func (gs googleGeocodingService) getPlaceCoordinates(placeId string) ([]maps.GeocodingResult, error) {
	request := &maps.GeocodingRequest{PlaceID: placeId}
	return gs.client.ReverseGeocode(context.Background(), request)
}

func (gs googleGeocodingService) getAddressCoordinates(address string) ([]maps.GeocodingResult, error) {
	request := &maps.GeocodingRequest{Address: address}
	return gs.client.Geocode(context.Background(), request)
}
