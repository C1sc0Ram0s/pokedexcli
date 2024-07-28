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

// Defining structs for Encountered Pokemon (explore)
type ExploreLocationAreas struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}
type PokemonEncounters struct {
	Pokemon EncounterdPokemon `json:"pokemon"`
}
type EncounterdPokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Defining structs for Pokemon (catch)
type Pokemon struct {
	Name           string  `json:"name"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Stats          []Stats `json:"stats"`
	Types          []Types `json:"types"`
	BaseExperience int     `json:"base_experience"`
}
type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effor"`
	Stat     Stat `json:"stat"`
}
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

// locationCache is globally accessible or passed
var locationCache pokecache.Cache = pokecache.NewCache(time.Minute * 2)

/*
	GetLocationAreas

args:

	nextUrl: string which defines the url for the next batch
	previousUrl: string which defines the url for the previous batch
	command: string which defines the command being used (i.e., map & mapb)
	args: slice which defines the arguments passed with the command

return:

	GetLocationAreas: struct which defines the output of the 'location-area' endpoint
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

/*
	GetExploreLocationAreas

args:

	nextUrl: string which defines the url for the next batch
	previousUrl: string which defines the url for the previous batch
	command: string which defines the command being used (i.e., map & mapb)
	args: slice which defines the arguments passed with the command

return:

	GetLocationAreas: struct which defines the output of the 'location-area' endpoint
	error: error handling
*/
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

/*
	GetPokemon

args:

	command: string which defines the command
	args: slice which defines the arguments of the command

return:

	Pokemon: Pokemon struct which contains the specified pokemon's data
	error: error handling
*/
func GetPokemon(command string, args []string) (Pokemon, error) {
	var response Pokemon
	endpoint := "https://pokeapi.co/api/v2/pokemon/" + args[1]

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
