package main

import (
	"errors"
	"fmt"
)

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("explore command requires a location name. Usage: explore <location_name>")
	}

	locationName := args[0]

	areaResp, err := cfg.pokeapiClient.GetAreaDetails(locationName)

	if err != nil {
		return err
	}
	
	fmt.Printf("Exploring %s...\n", areaResp.Name)

	if len(areaResp.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found in this area")
		return nil
	}

	fmt.Println("Pokemon in this area:")

	for _, pokemonEncounter := range areaResp.PokemonEncounters {
		fmt.Printf("- %s\n", pokemonEncounter.Pokemon.Name)
	}
	
	return nil
} 