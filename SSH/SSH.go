package ssh

import (
	"flag"
	"fmt"
	"log"
	"strings"

	config "Orchestrate/Config"
	utils "Orchestrate/Utils"

	ssh_lib "golang.org/x/crypto/ssh"
)

func HandleRunCommand(args []string) {

	/*****************
	*** Parse args ***
	******************/

	// Create a new FlagSet to parse the provided args
	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)

	// Initialize vars to hold flag values
	var flag_cmd string
	var flag_targetIP string
	var flag_targetIPs string
	var flag_target_all bool

	// Define and populate flag variables
	flagSet.StringVar(&flag_cmd, "x", "", "Command to run")
	flagSet.StringVar(&flag_targetIP, "t", "", "Specify target to run command on")
	flagSet.StringVar(&flag_targetIPs, "T", "", "Specify multiple targets to run command, seperated by a space in the quotes")
	flagSet.BoolVar(&flag_target_all, "a", false, "Command to run")

	// Parse the provided args slice
	err := flagSet.Parse(args)
	if err != nil {
		return
	}

	// Make sure that a command was given
	if flag_cmd == "" {
		fmt.Print("\n[Orchestrate]\n\tNo command was given\n\n")
		utils.PrintManPage("run")
	}

	/**********************
	*** Get the targets ***
	 **********************/

	// Load the config into the Config struct
	var conf config.Config
	config.LoadConfig(&conf)

	// Initialize targets array
	var targets []config.Target

	// Check for arg -t
	if flag_targetIP != "" {
		// Iterate through conf.Targets, checking if it has the IP specified
		for i := range conf.Targets {
			if flag_targetIP == conf.Targets[i].IP {
				targets = append(targets, conf.Targets[i])
			}
		}
	}

	// Check for arg -T
	if flag_targetIPs != "" {

		// Split the arg string into an array of IPs
		targetIPs := strings.Split(flag_targetIPs, " ")

		// Iterate through conf.Targets, checking if it has one of the IPs specified
		for i := range conf.Targets {
			for j := range targetIPs {
				if conf.Targets[i].IP == targetIPs[j] {
					targets = append(targets, conf.Targets[i])
				}
			}
		}
	}

	// Check for arg -a
	if flag_target_all {
		targets = append(targets, conf.Targets...)
	}

	// Make sure atleast 1 target was found
	if len(targets) == 0 {
		fmt.Print("\n[Orchestrate]\n\tNo targets found in Config.json that match the given IP(s)\n\n")
		return
	}

	// run the command with the targets
	runCommand(targets, flag_cmd)
}

func runCommand(targets []config.Target, cmd string) {
	for i := range targets {

		user := targets[i].User
		pass := targets[i].Pass
		ip := fmt.Sprintf("%s/22", targets[i].IP)

		// Configure the client config used to start SSH session
		ssh_config := &ssh_lib.ClientConfig{
			User: user,
			Auth: []ssh_lib.AuthMethod{
				ssh_lib.Password(pass),
			},
			HostKeyCallback: ssh_lib.InsecureIgnoreHostKey(),
		}

		// Make the SSH connection
		client, err := ssh_lib.Dial("tcp", ip, ssh_config)
		if err != nil {
			log.Fatal("Failed to dial: ", err)
		}
		defer client.Close()

		// Open a session
		session, err := client.NewSession()
		if err != nil {
			log.Fatal("Failed to create session: ", err)
		}
		defer session.Close()

		// Run the command and capture output
		output, err := session.CombinedOutput(cmd)
		if err != nil {
			log.Fatal("Failed to run command: ", err)
		}

		fmt.Printf("\n[Orchestrate] run\n%s@%s:~$%s\n%s", user, ip, cmd, string(output))
	}
}
