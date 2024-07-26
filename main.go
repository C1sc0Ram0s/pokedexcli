package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/C1sc0Ram0s/pokedexcli/internal/pokeapi"
)

// cliName is the name used in the repl prompts
var cliName string = "pokedex"

// printPrompt displays the repl prompt at the start of each loop
func printPrompt() {
	fmt.Print(cliName, "> ")
}

type config struct {
	nextUrl     string
	previousUrl string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println()
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the names of 20 location areas in the Pokemon world")
	fmt.Println("mapb: Displays the names of the previous 20 location areas in the Pokemon world")
	fmt.Println()
	return nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Exiting the Pokedex. Goodbye!")
	os.Exit(0)
	return nil // This line will never be reached, keeps function signiture consistent
}

func commandMap(cfg *config, args []string) error {
	command := args[0]
	result, err := pokeapi.GetLocationAreas(cfg.nextUrl, cfg.previousUrl, command, args)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	cfg.nextUrl = result.Next
	cfg.previousUrl = result.Previous
	for _, location := range result.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(cfg *config, args []string) error {
	command := args[0]
	result, err := pokeapi.GetLocationAreas(cfg.nextUrl, cfg.previousUrl, command, args)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	cfg.nextUrl = result.Next
	cfg.previousUrl = result.Previous
	for _, location := range result.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(cfg *config, args []string) error {
	command := args[0]
	if len(args) < 1 {
		return errors.New("invalid number of arguments")
	} else {
		fmt.Printf("Exploring %s...\n", args[1])
		result, err := pokeapi.GetExploreLocationAreas(command, args)
		if err != nil {
			fmt.Println("Error: ", err)
			return nil
		}
		fmt.Println("Found Pokemon:")
		for _, pokemon := range result.PokemonEncounters {
			fmt.Println(" - ", pokemon.Pokemon.Name)
		}

	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	//command := args[0]
	return nil
}

func main() {
	// CLI Commands
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb (map back)",
			description: "Displays the names of the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays a list of all the Pokemon in a given area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon and adds them to the user Pokedex",
			callback:    commandCatch,
		},
	}

	// Begin REPL loop
	reader := bufio.NewScanner(os.Stdin)
	cfg := config{}
	printPrompt()
	for reader.Scan() {
		input := reader.Text()
		args := strings.Split(input, " ") //The first argument will always be the command
		commandInput := args[0]
		if command, exists := commands[commandInput]; exists {
			command.callback(&cfg, args)
		} else {
			fmt.Printf("Unkown command: %s\n\n", input)
		}
		printPrompt()
	}
}
