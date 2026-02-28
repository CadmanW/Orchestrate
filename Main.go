package main

import (
	config "Orchestrate/Config"
	run "Orchestrate/Run"
	upload "Orchestrate/Upload"
	"fmt"
	"os"
)

func main() {

	/*********************
		COMMAND PARSING
	**********************/
	// If no command is provided, print the manual
	if len(os.Args) == 1 {
		printManPage()
		return
	}

	// Subcommand dispatch
	switch os.Args[1] {
	case "config":
		config.HandleConfigCommand(os.Args[2:])
	case "run":
		run.HandleRunCommand(os.Args[2:])
	case "upload":
		upload.HandleUploadCommand(os.Args[2:])
	default:
		// If command is unknown, print the manual
		printManPage()
	}
}

/*
* Prints the man page to the console
 */
func printManPage() error {
	// Read man page
	data, err := os.ReadFile("man.txt")
	if err != nil {
		return fmt.Errorf("Error reading man page: ", err)
	}

	// Output man page
	fmt.Println(string(data))

	// Return nil on success
	return nil
}
