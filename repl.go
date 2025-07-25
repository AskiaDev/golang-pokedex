package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AskiaDev/go-pokedex/internal/pokeapi"
)


type config struct {
	pokeapiClient pokeapi.Client
	nextLocationURL *string
	prevLocationURL *string
}

func startRepl(cfg *config){
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		reader.Scan()
		words := cleanInput(reader.Text())

		if len(words) == 0 {
			continue
		}

		if words[0] == "q" {
			break
		}

		command := words[0]
		args := words[1:] 

		cliCommand, exists := getCommands()[command]

		if exists {
			err := cliCommand.callback(cfg, args...)  // Pass args to callback
			if err != nil {
				fmt.Println("Error executing command:", err)
			}
		} else {
			fmt.Println("Invalid command")
			continue
		}

		fmt.Println("Your command is ", cliCommand.name)
		fmt.Println("Your command is ", cliCommand.description)
	}
}