package repl 

import (
	"strings"
	"fmt"
	"bufio"
	"os"
	"time"
	"pokedex/internal/commands"
	"pokedex/internal/models"
	"pokedex/internal/cache"
)

type CliCommand = models.CliCommand



func cleanInput(text string) []string{
	text = strings.ToLower(text)
	return strings.Fields(text)

}

func RunRepl(){
	commands.RegisterCliCommands()
	commands.RegisterPokedex()
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &models.CliCommandConfig{
		Cache: pokecache.NewCache(1*time.Minute),
	}
	running := true
	
	for running == true{
		fmt.Print("Pokedex > ")
		scanner.Scan()
		fullCommand :=  scanner.Text()
		splitCommands := cleanInput(fullCommand)
		if len(splitCommands) == 0{
			fmt.Println("empty command detected - try again!")
			continue
		} 
		mainCommand := splitCommands[0]
		cfg.Arguments = splitCommands[1:]
		cliCommand, ok := commands.CliCommands[mainCommand]
		if ok {
			err := cliCommand.Callback(cfg)
			if err != nil{
				fmt.Println("Encountered an error when executing the command")
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
		
	}

}
