package commands

import(
	"fmt"
	"os"
	"net/http"
	"io"
	"encoding/json"
	"log"
	"errors"
	"pokedex/internal/models"
	"math/rand"
	"strings"
)

type CliCommand = models.CliCommand

var CliCommands map[string]CliCommand
var Pokedex map[string]models.PokemonSimplified

func RegisterPokedex(){
	Pokedex = map[string]models.PokemonSimplified{}
}

func RegisterCliCommands(){
	CliCommands = map[string] CliCommand{
		"help": {
			Name: "help",
			Description: "Displays a help message",
			Callback: CommandHelp,
		},

		"exit": {
			Name: "exit",
			Description: "Close the Pokedex",
			Callback: CommandExit,
		},
		"map": {
			Name: "map",
			Description: "View Next 20 Locations",
			Callback: CommandMap,
		},
		"mapb": {
			Name: "mapb",
			Description: "View Previous 20 Locations",
			Callback: CommandMapB,
		},
		"explore": {
			Name: "explore <REGION>",
			Description: "Lists pokemon available in REGION",
			Callback: CommandExplore,
		},
		"catch": {
			Name: "catch <POKEMON>",
			Description: "Attempt to catch a POKEMON",
			Callback: CommandCatch,
		},
	}
}


func CommandExit(cfg *models.CliCommandConfig) error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("os.Exit failed!")
}

func CommandHelp(cfg *models.CliCommandConfig) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, helptext := range CliCommands {
		fmt.Printf("%s: %s\n", helptext.Name, helptext.Description )
	}
	return nil
}

func CommandMap(cfg *models.CliCommandConfig) error {
	if cfg.NextURL == "" {
		cfg.NextURL = "https://pokeapi.co/api/v2/location-area"
	}
	err := requestMap(cfg.NextURL, cfg)
	return err
}

func CommandMapB(cfg *models.CliCommandConfig) error {
	if cfg.PreviousURL == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	err := requestMap(cfg.PreviousURL, cfg)
	return err

}

func requestMap(url string, cfg *models.CliCommandConfig) error {
	data_body, ok := cfg.Cache.Get(url)
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
		cfg.Cache.Add(url,data_body)
	}
	locArea := models.LocationAreaResponse{}
	err := json.Unmarshal(data_body,&locArea)
	if err != nil {
		fmt.Errorf("Failed to Unmarshal body!")
		return err
	}
	for _,result := range locArea.Results{
		fmt.Println(result.Name)
	}
	if locArea.Next != nil{
		cfg.NextURL = *locArea.Next
	} else {
		cfg.NextURL = ""
	}
	if locArea.Previous != nil {
		cfg.PreviousURL = *locArea.Previous
	} else {
		cfg.PreviousURL = ""
	}
	return nil

}

func CommandExplore(cfg *models.CliCommandConfig) error {
	if len(cfg.Arguments) == 0 {
		return errors.New("Missing region for explore (Usage: explore <REGION>)")
	}
	url := "https://pokeapi.co/api/v2/location-area/" + cfg.Arguments[0]
	data_body, ok := cfg.Cache.Get(url)
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
		cfg.Cache.Add(url,data_body)
	}
	locAreaDetails := models.LocationAreaDetailsResponse{}
	err := json.Unmarshal(data_body,&locAreaDetails)
	if err != nil {
		fmt.Errorf("Failed to Unmarshal body!")
		return err
	}
	fmt.Println("Exploring " + cfg.Arguments[0] + "...")
	if len(locAreaDetails.PokemonEncounters)>0{
		fmt.Println("Found Pokemon:")
		for _,result := range locAreaDetails.PokemonEncounters{
			fmt.Println(result.Pokemon.Name)
		}
	}

	return nil
}

func CommandCatch(cfg *models.CliCommandConfig) error {
	if len(cfg.Arguments) == 0 {
		return errors.New("Are you trying to catch air? (Usage: catch <POKEMON NAME>)")
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + cfg.Arguments[0]
	data_body, ok := cfg.Cache.Get(url)
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
		cfg.Cache.Add(url,data_body)
	}
	pokedetails := models.PokemonSimplified{}
	err := json.Unmarshal(data_body,&pokedetails)
	if err != nil {
		fmt.Errorf("Failed to Unmarshal body!")
		return err
	}
	thrown_ball := "Pokeball"
	if len(cfg.Arguments)>1{
		thrown_ball = cfg.Arguments[1]
	}
	caught := catch_pokemon(pokedetails,thrown_ball) 
	if caught {
		fmt.Println(cfg.Arguments[0] + " was caught!")
		register_to_pokedex(pokedetails)
	} else {
		fmt.Println(cfg.Arguments[0] + " escaped!")
	}
	return nil

}

func catch_pokemon(pokedetails models.PokemonSimplified, ballType string) bool {
	var throw_quality int
	if ballType == "superball"{
		throw_quality = 100 +rand.Intn(500)
	} else if ballType == "hyperball" {
		throw_quality = 100 + rand.Intn(700)
	} else if ballType == "masterball" {
		throw_quality = 9001
	} else {
		ballType = "Pokeball"
		throw_quality = rand.Intn(500)
	}
	fmt.Println("Throwing a " +strings.Title(ballType) + " at " + pokedetails.Name + "...")

	return throw_quality>=pokedetails.BaseExp
}

func register_to_pokedex(pokemon models.PokemonSimplified) {
	fmt.Println("Trying to register " + pokemon.Name + " to the Pokedex...")
	_, ok := Pokedex[pokemon.Name]
	if ok{
		fmt.Println("Pokemon is already registered! Try catching new Pokemon to compelte your Pokedex!")
		return
	} else {
	Pokedex[pokemon.Name] = pokemon
	fmt.Printf("A new Pokemon! We now have %d Pokemons in our Pokedex!\n",len(Pokedex))
	}
}