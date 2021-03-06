package services

import (
	"context"
	"dgraph-osm/database/internal/connection"
	"dgraph-osm/database/internal/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgraph-io/dgo/v210/protos/api"
)

type Cities struct {
	Name        string `json:"name"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	CountryName string `json:"country_name"`
	CountryCode string `json:"country_code"`
	StateName   string `json:"state_name"`
	StateCode   string `json:"state_code"`
}

func getCities() []*models.City {
	cityFile, _ := os.Open("/home/douglasblnk/Documentos/Projects/dgraph-go/database/data/cities.json")
	defer cityFile.Close()

	citiesByteValue, _ := ioutil.ReadAll(cityFile)

	cityModel := []*Cities{}
	json.Unmarshal(citiesByteValue, &cityModel)

	var cities []*models.City

	for _, city := range cityModel {
		c := &models.City{
			Name:        city.Name,
			Latitude:    city.Latitude,
			Longitude:   city.Longitude,
			CountryName: city.CountryName,
			CountryCode: city.CountryCode,
			StateName:   city.StateName,
			StateCode:   city.StateCode,
		}

		cities = append(cities, c)
	}

	return cities
}

func InsertCities() {
	dgcon, close := connection.DgraphConnection()
	defer close()

	cities := getCities()

	json, _ := json.Marshal(&cities)

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   json,
	}

	response, err := dgcon.NewTxn().Mutate(context.Background(), mu)

	if err != nil {
		log.Fatal("err", err)
	}

	println(response)
}
