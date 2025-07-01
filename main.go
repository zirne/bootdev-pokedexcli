package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/zirne/pokedexcli/internal/pokeapi"
)

var cmdMap = map[string]cliCommand{}
var api = pokeapi.NewClient()
var pokedex = make(map[string]Pokemon)

type Pokemon struct {
	Caught   bool
	PokeData pokeapi.Pokemon
}

func cleanInput(text string) []string {
	r := make([]string, 0)
	words := strings.Fields(text)
	for _, word := range words {
		r = append(r, strings.ToLower(word))
	}
	return r
}

type config struct {
	LastCall string
	Next     string
	Previous any
}

type cliCommand struct {
	name        string
	description string
	callback    func(C *config, param string) error
}

func commandExit(c *config, param string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, param string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range cmdMap {
		fmt.Println(fmt.Sprintf("%v: %v", cmd.name, cmd.description))
	}
	return nil
}

func commandPokedex(c *config, param string) error {
	fmt.Println("Your Pokédex:")
	for k, _ := range pokedex {
		fmt.Println(" - " + k)
	}

	return nil
}

func commandExplore(c *config, param string) error {
	body, err := api.Get("https://pokeapi.co/api/v2/location-area/" + param + "/")
	if err != nil {
		return err
	}
	obj := pokeapi.LocationArea{}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokémon:")
	for _, o := range obj.PokemonEncounters {
		fmt.Println(" - " + o.Pokemon.Name)
	}
	return nil
}

func commandInspect(c *config, param string) error {
	if param == "" {
		return fmt.Errorf("A Pokémon name must be provided")
	}
	if _, ok := pokedex[strings.ToLower(param)]; !ok {
		return fmt.Errorf("You haven't attempted to catch this Pokémon")
	}
	o := pokedex[strings.ToLower(param)]
	if !o.Caught {
		return fmt.Errorf("You haven't caught this Pokémon")
	}
	fmt.Println("Name: " + o.PokeData.Name)
	fmt.Println(fmt.Sprintf("Height: %v", o.PokeData.Height))
	fmt.Println(fmt.Sprintf("Weight: %v", o.PokeData.Weight))
	// fmt.Println(fmt.Sprintf("Stats: %v", o.PokeData.Stats))
	fmt.Println("Stats:")
	for _, stat := range o.PokeData.Stats {
		fmt.Println(fmt.Sprintf(" -%v: %v", stat.Stat.Name, stat.BaseStat))
	}
	fmt.Println("Types:")
	for _, t := range o.PokeData.Types {
		fmt.Println(fmt.Sprintf(" - %v", t.Type.Name))
	}
	return nil
}

func commandCatch(c *config, param string) error {
	if param == "" {
		return fmt.Errorf("A Pokémon name must be provided")
	}
	body, err := api.Get("https://pokeapi.co/api/v2/pokemon/" + param + "/")
	if err != nil {
		return err
	}
	obj := pokeapi.Pokemon{}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return err
	}
	// Add info to pokedex if not present
	if _, ok := pokedex[strings.ToLower(param)]; !ok {
		pokedex[strings.ToLower(param)] = Pokemon{false, pokeapi.Pokemon{}}
	}
	// Attempt catch
	fmt.Println("Throwing a Pokeball at " + param + "...")
	chance := 512 / rand.Intn(obj.BaseExperience)
	if chance >= 10 {
		pokedex[strings.ToLower(param)] = Pokemon{true, obj}
		fmt.Println(param + " was caught!")
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Println(param + " escaped!")
	}
	return nil
}

func commandMap(c *config, params string) error {
	var url = "https://pokeapi.co/api/v2/location-area/"
	if c.LastCall == "map" && c.Next != "" {
		url = c.Next
	}
	body, err := api.Get(url)
	if err != nil {
		return err
	}
	obj := pokeapi.LocationAreas{}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return err
	}
	for _, o := range obj.Results {
		fmt.Println(o.Name)
	}
	c.Next = obj.Next
	c.Previous = obj.Previous
	return nil
}

func commandMapb(c *config, param string) error {
	if c.LastCall == "map" && c.Previous != nil {
		url := c.Previous.(string)
		if c.LastCall == "map" && c.Next != "" {
			url = c.Next
		}
		body, err := api.Get(url)
		obj := pokeapi.LocationAreas{}
		err = json.Unmarshal(body, &obj)
		if err != nil {
			return err
		}
		for _, o := range obj.Results {
			fmt.Println(o.Name)
		}
		c.Next = obj.Next
		c.Previous = obj.Previous
		return nil

	}
	fmt.Println("you're on the first page")
	return nil
}

func setData() { // Får inte heta init av nån anledning
	cmdMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Lists next map locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists previous map locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Lists Pokémon in the given area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a given Pokémon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a given Pokémon (if you have caught it)",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all your caught Pokémon",
			callback:    commandPokedex,
		},
	}
}

func main() {
	setData()
	scanner := bufio.NewScanner(os.Stdin)
	cfg := config{}
	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		clean := cleanInput(input)
		if len(clean) > 0 {
			cmd := clean[0]
			if v, ok := cmdMap[cmd]; ok {
				param := ""
				if len(clean) > 1 {
					param = clean[1]
				}
				err := v.callback(&cfg, param)
				if err != nil {
					fmt.Println(fmt.Errorf("Error: %v", err))
				}
				cfg.LastCall = cmd
			} else {
				fmt.Println("Unknown command")
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}
