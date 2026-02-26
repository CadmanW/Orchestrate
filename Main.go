package main

import (
	config "Orchestrate/Config"
	ssh "Orchestrate/SSH"
	utils "Orchestrate/Utils"
	"os"
)

func main() {

	/*********************
		COMMAND PARSING
	**********************/
	// If no command is provided, print the manual
	if len(os.Args) == 1 {
		utils.PrintManPage("orchestrate")
		return
	}

	// Subcommand dispatch
	switch os.Args[1] {
	case "config":
		config.HandleConfigCommand(os.Args[2:])
	case "run":
		ssh.HandleRunCommand(os.Args[2:])
	default:
		// If command is unknown, print the manual
		utils.PrintManPage("orchestrate")
	}
}
