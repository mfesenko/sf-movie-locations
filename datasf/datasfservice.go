package datasf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type DataSFService struct {
	baseUrl string
}

type DataSFRecord struct {
	Actor_1      string `json:"actor_1"`
	Actor_2      string `json:"actor_2"`
	Actor_3      string `json:"actor_3"`
	Director     string `json:"director"`
	Locations    string `json:"locations"`
	Release_year string `json:"release_year"`
	Title        string `json:"title"`
}

const DATASF_URL_TEMPLATE = "%s?$offset=%d&$limit=%d"

func NewDataSFService(baseUrl string) (*DataSFService, error) {
	_, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("Invalid base url: %s", baseUrl)
	}
	return &DataSFService{baseUrl: baseUrl}, nil
}

func (s DataSFService) RetrieveRecords(offset int, limit int) ([]DataSFRecord, error) {
	url := fmt.Sprintf(DATASF_URL_TEMPLATE, s.baseUrl, offset, limit)
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		return []DataSFRecord{}, fmt.Errorf("datasf: Request to %s was not successful", s.baseUrl)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []DataSFRecord{}, errors.New("datasf: Failed to read response content")
	}
	var records []DataSFRecord
	err = json.Unmarshal(body, &records)
	if err != nil {
		return []DataSFRecord{}, errors.New("datasf: Failed to deserialize records from json response")
	}
	return records, nil
}
