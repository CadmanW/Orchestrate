package upload

import (
	config "Orchestrate/Config"
	"flag"
	"fmt"
	"log"

	"github.com/melbahja/goph"
)

func HandleUploadCommand(args []string) {
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
	flagSet.StringVar(&flag_file, "f", "", "File to upload")
	flagSet.StringVar(&flag_directory, "F", "", "Directory to upload")
	flagSet.StringVar(&flag_destination, "d", "", "Destination to upload files to")

	flagSet.StringVar(&flag_targetIP, "t", "", "Specify target to upload file(s) to")
	flagSet.StringVar(&flag_targetIPs, "T", "", "Specify multiple targets to upload file(s) to, seperated by a space in the quotes")
	flagSet.BoolVar(&flag_target_all, "a", false, "Upload file(s) to all targets in Config.json")

	// Parse the provided args slice
	err := flagSet.Parse(args)
	if err != nil {
		log.Fatalf("\n[Orchestrate] upload\n\tError when parsing args: %s\n\n", err)
	}

	// Get the targets
	targets := config.GetTargets(flag_targetIP, flag_targetIPs, flag_target_all)

	// Make sure atleast 1 target was found
	if len(targets) == 0 {
		log.Fatal("\n[Orchestrate] upload\n\tNo targets found in Config.json that match the given IP(s)\n\n")
		return
	}

	// Make sure a destination was provided
	if flag_destination == "" {
		log.Fatal("\n[Orchestrate] upload\n\tNo destination was provided")
	}

	// Loop over the targets, uploading file(s) to each, depending on the flag provided
	// Also make sure file(s) were provided
	for i := range targets {
		if flag_file != "" {
			err = uploadFile(targets[i], flag_file, flag_destination)
		} else if flag_directory != "" {
			err = uploadDirectory(targets[i], flag_directory, flag_destination)
		} else {
			log.Fatal("\n[Orchestrate] upload\n\tNo file or directory was provided")
		}
	}

	// Check for upload errors
	if err != nil {
		log.Fatalf("\n[Orchestrate] upload \n\tError uploading file(s): %s", err)
	}
}

// Helper function to upload a file
func uploadFile(target config.Target, filePath string, destinationPath string) error {
	// Create the connection to the target
	client, err := goph.New(target.User, target.IP, goph.Password(target.Pass))
	if err != nil {
		return fmt.Errorf("Error creating connection to target: %s ; %s", target.IP, err)
	}
	defer client.Close()

	err = client.Upload(filePath, destinationPath)
	if err != nil {
		return fmt.Errorf("Error uploading file: %s", err)
	}

	// return nil on success
	fmt.Printf("[Orchestrate] upload\n\tUploaded file %s to %s at %s successfully\n\n", filePath, target.IP, destinationPath)
	return nil
}

// TODO Helper function to upload a directory
func uploadDirectory(target config.Target, directoryPath string, destinationPath string) error {
	// Create the connection to the target
	client, err := goph.New(target.User, target.IP, goph.Password(target.Pass))
	if err != nil {
		return fmt.Errorf("Error creating connection to target: %s ; %s", target.IP, err)
	}
	defer client.Close()

	//TODO upload the directory

	// Return nil on success
	return nil
}
