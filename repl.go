package main 

import (
	"strings"
	"fmt"
	"bufio"
	"os"
	"log"
	"net/http"
	"io"
	"encoding/json"
	"errors"
	"time"
	"pokedex/internal"
)

var cliCommands map[string]cliCommand

func registerCliCommands(){
	cliCommands = map[string] cliCommand{
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},

		"exit": {
			name: "exit",
			description: "Close the Pokedex",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "View Next 20 Locations",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "View Previous 20 Locations",
			callback: commandMapB,
		},
	}
}




type cliCommandConfig struct {
    nextURL     string
    previousURL string
	cache *pokecache.Cache
}

type cliCommand struct {
	name string
	description string
	callback func(*cliCommandConfig) error
}

type locationAreaResponse struct {
    Count    int     `json:"count"`
    Next     *string `json:"next"`
    Previous *string `json:"previous"`
    Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}


func cleanInput(text string) []string{
	text = strings.ToLower(text)
	return strings.Fields(text)

}

func runRepl(){
	registerCliCommands()
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &cliCommandConfig{
		cache: pokecache.NewCache(1*time.Minute),
	}
	running := true
	
	for running == true{
		fmt.Print("Pokedex > ")
		scanner.Scan()
		fullCommand :=  scanner.Text()
		splitCommands := cleanInput(fullCommand)
		if len(splitCommands) == 0{
			fmt.Println("empty command detected - try again!")
		} 
		mainCommand := splitCommands[0]

		cliCommand, ok := cliCommands[mainCommand]
		if ok {
			cliCommand.callback(cfg)
		} else {
			fmt.Println("Unknown command")
		}
		
	}

}

func commandExit(cfg *cliCommandConfig) error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("os.Exit failed!")
}

func commandHelp(cfg *cliCommandConfig) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, helptext := range cliCommands {
		fmt.Printf("%s: %s\n", helptext.name, helptext.description )
	}
	return nil
}

func commandMap(cfg *cliCommandConfig) error {
	if cfg.nextURL == "" {
		cfg.nextURL = "https://pokeapi.co/api/v2/location-area"
	}
	err := requestMap(cfg.nextURL, cfg)
	return err
}

func commandMapB(cfg *cliCommandConfig) error {
	if cfg.previousURL == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	err := requestMap(cfg.previousURL, cfg)
	return err

}

func requestMap(url string, cfg *cliCommandConfig) error {
	data_body, ok := cfg.cache.Get(url)
	if !ok{
		res, err := http.Get(url)
		if err != nil {
			fmt.Errorf("Encoutered error when requesting data!")
			return err
		}
		data_body, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if err!= nil{
			log.Fatal(err)
		}
		if res.StatusCode > 299 {
			fmt.Errorf("Response failed with status code %d.", res.StatusCode)
			return errors.New(fmt.Sprintf("Failed with responce Code %d", res.StatusCode))
		}
		cfg.cache.Add(url,data_body)
	}
	locArea := locationAreaResponse{}
	err := json.Unmarshal(data_body,&locArea)
	if err != nil {
		fmt.Errorf("Failed to Unmarshal body!")
		return err
	}
	for _,result := range locArea.Results{
		fmt.Println(result.Name)
	}
	if locArea.Next != nil{
		cfg.nextURL = *locArea.Next
	} else {
		cfg.nextURL = ""
	}
	if locArea.Previous != nil {
		cfg.previousURL = *locArea.Previous
	} else {
		cfg.previousURL = ""
	}
	return nil

}
