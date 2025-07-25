package main

import (
	"errors"
	"fmt"
)

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("explore command requires a location name. Usage: explore <location_name>")
	}

	locationName := args[0]

	location, err := cfg.pokeapiClient.GetLocationByName(locationName)

	if err != nil {
		return err
	}
	
	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("URL: ", location.URL)
	
	if len(args) > 1 {
		fmt.Printf("Additional arguments provided: %v\n", args[1:])
	}

	areaResp, err := cfg.pokeapiClient.GetAreaDetails(location.URL)

	if err != nil {
		return err
	}

	fmt.Printf("Area: %s\n", areaResp.Name)

	for _, pokemonEncounter := range areaResp.PokemonEncounters {
		fmt.Printf("Pokemon: %s\n", pokemonEncounter.Pokemon.Name)
	}
	
	return nil
} 