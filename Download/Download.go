package download

import (
	config "Orchestrate/Config"
	"flag"
	"log"
)

func HandleDownloadCommand(args []string) {
	// Create a new FlagSet to parse the provided args
	flagSet := flag.NewFlagSet("upload", flag.ContinueOnError)

	// Initialize vars to hold flag values
	var flag_file string
	var flag_directory string
	var flag_destination string

	var flag_targetIP string
	var flag_targetIPs string
	var flag_target_all bool

	// Define and populate flag variables
	flagSet.StringVar(&flag_file, "f", "", "File to download")
	flagSet.StringVar(&flag_directory, "F", "", "Directory to download")
	flagSet.StringVar(&flag_destination, "d", "", "Destination to download files to")

	flagSet.StringVar(&flag_targetIP, "t", "", "Specify target to download file(s) from")
	flagSet.StringVar(&flag_targetIPs, "T", "", "Specify multiple targets to download file(s) from, seperated by a space in the quotes")
	flagSet.BoolVar(&flag_target_all, "a", false, "download file(s) from all targets in Config.json")

	// Parse the provided args slice
	err := flagSet.Parse(args)
	if err != nil {
		log.Fatalf("\n[Orchestrate] upload\n\tError when parsing args: %s\n\n", err)
	}

	// Get the targets
	targets := config.GetTargets(flag_targetIP, flag_targetIPs, flag_target_all)

	// Make sure atleast 1 target was found
	if len(targets) == 0 {
		log.Fatal("\n[Orchestrate] download\n\tNo targets found in Config.json that match the given IP(s)\n\n")
		return
	}

	// Make sure a destination was provided
	if flag_destination == "" {
		log.Fatal("\n[Orchestrate] download\n\tNo destination was provided")
	}

	// Loop over all artgets found, download
}

func downloadFile() error {

	return nil
}

func downloadDirectory() error {

	return nil
}
