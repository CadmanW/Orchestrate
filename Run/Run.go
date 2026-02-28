package run

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/melbahja/goph"

	config "Orchestrate/Config"
)

func HandleRunCommand(args []string) {

	/*****************
	*** Parse args ***
	******************/

	// Create a new FlagSet to parse the provided args
	flagSet := flag.NewFlagSet("run", flag.ContinueOnError)

	// Initialize vars to hold flag values
	var flag_cmd string
	var flag_sudo bool

	var flag_targetIP string
	var flag_targetIPs string
	var flag_target_all bool

	// Define and populate flag variables
	flagSet.StringVar(&flag_cmd, "x", "", "Command to run")
	flagSet.BoolVar(&flag_sudo, "s", false, "Command to run as sudo")

	flagSet.StringVar(&flag_targetIP, "t", "", "Specify target to run command on")
	flagSet.StringVar(&flag_targetIPs, "T", "", "Specify multiple targets to run command, seperated by a space in the quotes")
	flagSet.BoolVar(&flag_target_all, "a", false, "Run command on all targets in Config.json")

	// Parse the provided args slice
	err := flagSet.Parse(args)
	if err != nil {
		log.Fatalf("\n[Orchestrate] run\n\tError when parsing args: %s\n\n", err)
	}

	// Make sure that a command was given
	if flag_cmd == "" {
		log.Fatal("\n[Orchestrate] run\n\tNo command was given\n\n")
	}

	// Get the targets
	targets := config.GetTargets(flag_targetIP, flag_targetIPs, flag_target_all)

	// Make sure atleast 1 target was found
	if len(targets) == 0 {
		log.Fatal("\n[Orchestrate] run\n\tNo targets found in Config.json that match the given IP(s)\n\n")
	}

	// run the command with the targets
	for i := range targets {
		output, err := runCommandOverSSH(targets[i], flag_cmd, flag_sudo)
		if err != nil {
			log.Fatalf("\n[Orchestrate] run: Error on %s: %v\n", targets[i].IP, err)
		}
		fmt.Printf("\n[Orchestrate] run\n%s@%s:~$ %s\n%s", targets[i].User, targets[i].IP, flag_cmd, output)
	}
}

/*********************
** Helper Functions **
*********************/

func runCommandOverSSH(target config.Target, cmd string, sudo bool) (string, error) {
	// Create the connection to the target
	client, err := goph.New(target.User, target.IP, goph.Password(target.Pass))
	if err != nil {
		return "", fmt.Errorf("Error connecting to target: %s", err)
	}
	defer client.Close()

	// Check if we should run the command with sudo or not
	if sudo {
		// Create a new session
		session, err := client.NewSession()
		if err != nil {
			return "", fmt.Errorf("Error creating session: %s", err)
		}
		defer session.Close()

		// Set up stdin to pipe the password
		session.Stdin = strings.NewReader(target.Pass + "\n")

		// Run the command as sudo and capture combined output
		// -S: read password from stdin
		// -p '': empty prompt to avoid mixing with output
		sudoCmd := fmt.Sprintf("sudo -S -p '' %s", cmd)
		output, err := session.CombinedOutput(sudoCmd)
		if err != nil {
			return string(output), err
		}
		return string(output), nil
	} else {
		output, err := client.Run(cmd)
		if err != nil {
			return string(output), err
		}
		return string(output), nil
	}
}
