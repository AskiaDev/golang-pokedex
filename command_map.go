package main

import "fmt"



func commandMap(cfg *config, args ...string) error {
	fmt.Println("Getting next page of locations")
	
	locationResp, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationURL)
	
	if err != nil {
		return err
	}

	cfg.nextLocationURL = &locationResp.Next
	cfg.prevLocationURL = &locationResp.Previous

	for _, location := range locationResp.Results {
		fmt.Println(location.Name)
	} 

	return nil
}