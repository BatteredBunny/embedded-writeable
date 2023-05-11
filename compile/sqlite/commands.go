package main

import (
	"fmt"
	"time"
)

var commands = map[string]func(state *State){
	"add":  AddCommand,
	"exit": ExitCommand,
	"help": HelpCommand,
}

func AddCommand(state *State) {
	fmt.Println("Adding row")
	if _, err := state.db.Query(
		insertDataQuery,
		"test",
		time.Now().Format("2006/01/02 15:04:05"),
	); err != nil {
		PrintError(err)
	}
}

func ExitCommand(state *State) {
	fmt.Println("Exiting!")
	state.Running = false
}
func HelpCommand(state *State) {
	fmt.Println("Commands: /add /exit /help")
}
