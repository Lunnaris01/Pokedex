package main 

import (
	"strings"
	"fmt"
	"bufio"
	"os"

)

func cleanInput(text string) []string{
	text = strings.ToLower(text)
	return strings.Fields(text)

}

func run_repl(){
	scanner := bufio.NewScanner(os.Stdin)
	//var full_command string
	//var main_command string
	running := true
	

	for running == true{
		fmt.Print("Pokedex > ")
		scanner.Scan()
		full_command :=  scanner.Text()
		split_commands := cleanInput(full_command)
		if len(split_commands) == 0{
			fmt.Println("empty command detected - use q to quit!")
		} else if split_commands[0] == "q" || split_commands[0] == "quit" {
			running = false
		}

		fmt.Printf("Your command was: %s\n", split_commands[0])
		
	}

}