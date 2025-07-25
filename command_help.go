package main

import "fmt"

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available commands:")
	fmt.Println("help - Prints the help menu")
	fmt.Println("exit - Exits the Pokedex")
	fmt.Println("map - Get the next page of locations")
	fmt.Println("explore <location_name> - Explore a specific location")
	return nil
}