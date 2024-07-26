package pokeapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/C1sc0Ram0s/pokedexcli/internal/pokecache"
)

// Defining structs for LocationAreas
type LocationAreas struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Results `json:"results"`
}
type Results struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Defining structs for Explore
type ExploreLocationAreas struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}
type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}
type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// locationCache is globally accessible or passed
var locationCache pokecache.Cache = pokecache.NewCache(time.Minute * 2)

/*
	LocationAreas

args:

	nextUrl: string which defines the url for the next batch
	previousUrl: string which defines the url for the previous batch
	command: string which defines the command being used (i.e., map & mapb)

return:

	LocationAreas: struct which defines the output of the 'location-area' endpoint
	error: error handling
*/
func GetLocationAreas(nextUrl, previousUrl, command string, args []string) (LocationAreas, error) {
	var response LocationAreas
	var endpoint string

	if command == "map" { //commandMap - Returns next batch
		if nextUrl == "" {
			endpoint = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
		} else {
			endpoint = nextUrl
		}
	} else if command == "mapb" { //commandMapb - Returns previous batch
		if previousUrl == "" {
			return response, errors.New("there are no more previous locations")
		} else {
			endpoint = previousUrl
		}
	} else {
		return response, errors.New("command using LocationAreas endpoint was not found")
	}

	// Check if the endpoint is cached
	if value, exists := locationCache.Get(endpoint); exists {
		json.Unmarshal(value, &response)
		return response, nil
	} else {
		res, err := http.Get(endpoint)
		if err != nil {
			return response, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return response, errors.New("unexpected status code")
		}

		// Decode API response
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return response, err
		}

		// Cache result
		data, err := json.Marshal(response)
		if err == nil {
			locationCache.Add(endpoint, data)
		}
	}

	return response, nil
}

func GetExploreLocationAreas(command string, args []string) (ExploreLocationAreas, error) {
	var response ExploreLocationAreas
	endpoint := "https://pokeapi.co/api/v2/location-area/" + args[1]

	// Check if the endpoint is cached
	if value, exists := locationCache.Get(endpoint); exists {
		json.Unmarshal(value, &response)
		return response, nil
	} else {
		res, err := http.Get(endpoint)
		if err != nil {
			return response, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return response, errors.New("unexpected status code")
		}

		// Decode API response
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return response, err
		}

		// Cache result
		data, err := json.Marshal(response)
		if err == nil {
			locationCache.Add(endpoint, data)
		}
	}

	return response, nil

}
