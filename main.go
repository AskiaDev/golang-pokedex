package main

import (
	"strings"
	"time"

	"github.com/AskiaDev/go-pokedex/internal/pokeapi"
	"github.com/AskiaDev/go-pokedex/internal/pokecache"
)


type CliCommand struct {
	name string
	description string
	callback func(*config, ...string) error
}

func cleanInput(input string) []string {
	words := strings.Fields(input)

	return words
}

func getCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"help": {
			name:        "help",
			description: "Prints the help menu",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMap,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific location",
			callback:    commandExplore,
		},
		// "mapb": {
		// 	name:        "mapb",
		// 	description: "Get the previous page of locations",
		// 	callback:    commandMapb,
		// },
	}
}



func main() {
	cache := pokecache.NewCache(5 * time.Minute)
	pokeClient := pokeapi.NewClient(5*time.Second, cache)
	cfg := &config{
		pokeapiClient: pokeClient,
	}
	startRepl(cfg)
}