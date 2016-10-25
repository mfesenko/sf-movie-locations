package datasf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DataSFService struct {
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

const DATASF_URL_TEMPLATE = "https://data.sfgov.org/resource/wwmu-gmzc.json?$offset=%d&$limit=%d"

func NewDataSFService() *DataSFService {
	return &DataSFService{}
}

func (s DataSFService) RetrieveRecords(offset int, limit int) []DataSFRecord {
	url := fmt.Sprintf(DATASF_URL_TEMPLATE, offset, limit)
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		log.Fatal("Request to data.sfgov.org was not successful")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Failed to read response content")
	}
	var records []DataSFRecord
	err = json.Unmarshal(body, &records)
	if err != nil {
		log.Fatal("Failed to deserialize records from json response")
	}
	return records
}
