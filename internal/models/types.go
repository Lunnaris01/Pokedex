package models


import (
	"pokedex/internal/cache"
)


type CliCommandConfig struct {
    NextURL     string
    PreviousURL string
	Arguments []string
	Cache *pokecache.Cache
}


type CliCommand struct {
	Name string
	Description string
	Callback func(*CliCommandConfig) error
}

type LocationAreaResponse struct {
    Count    int     `json:"count"`
    Next     *string `json:"next"`
    Previous *string `json:"previous"`
    Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}

type LocationAreaDetailsResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Game_index int `json:"game_index"`
	EncouterMethodRates  [] struct {
		EncounterMethod struct {
			Name string `json:"name"`
			ULR string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`

	} `json:"encounter_method_rates"`


	Location struct {
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name string `json:"name"`
		Language struct{
			Name string `json:"name"`
			URL string `json:"url"`
		} `json:"language"`

	} `json:"names"`

	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`

}

type PokemonEncounter struct {
	Pokemon struct{
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"pokemon"`
	VersionDetails []struct {
		Version struct {
			Name string `json:"name"`
			URL string `json:"url"`
		} `json:"version"`
		MaxChance int `json:"max_chance"`
		Encounter_details []struct{
			MinLevel int `json:"min_level"`
			MaxLevel int `json:"max_level"`
			ConditionValues []struct{
				Name string `json:"name"`
				URL string `json:"url"`
			} `json:"condition_values"`
			Chance int `json:"chance"`
			Method struct {
				Name string `json:"name"`
				URL string `json:"url"`
			} `json:"method"`
		} `json:"encounter_details"`
	} `json:"version_details"`
}


type PokemonSimplified struct{
	Name string `json:"name"`
	BaseExp int `json:"base_experience"`

}