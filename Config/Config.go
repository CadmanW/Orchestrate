package config

import (
	utils "Orchestrate/Utils"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Target struct {
	User string `json:"user"`
	Pass string `json:"password"`
	IP   string `json:"ip"`
}

type Config struct {
	Targets []Target `json:"targets"`
}

/*
* Used to interact with config using args provided
 */
func HandleConfigCommand(args []string) {

	// Create a new FlagSet to parse the provided args
	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)

	// Initialize vars to hold flags values
	var flag_add string
	var flag_remove string

	// Define and populate flag variables
	flagSet.StringVar(&flag_add, "a", "", "Add a target to Config.json in the form of user:pass@ip")
	flagSet.StringVar(&flag_remove, "r", "", "Removes a target from Config.json by IP")

	// Parse the provided args slice
	err := flagSet.Parse(args)
	if err != nil {
		return
	}

	// ADD TARGET
	if flag_add != "" {
		AddTarget(flag_add)
		return
	}

	// REMOVE TARGET
	if flag_remove != "" {
		RemoveTarget(flag_remove)
		return
	}

	// No valid option provided
	utils.PrintManPage("config")
}

/*
* Adds a target to the config
 */
func AddTarget(targetArg string) {
	// Normalize the target argument so we can split it up
	// user:pass@ip --> user:pass:ip
	normalized := strings.ReplaceAll(targetArg, "@", ":")

	// Split up the target argument
	// "user:pass:ip" --> ["user", "pass", "ip"]
	parts := strings.Split(normalized, ":")

	// Make sure user:pass@ip format was followed
	if len(parts) != 3 {
		fmt.Printf(
			"\n[Orchestrate]\n\tInvalid target format: %s\n\tUse: username:password@127.0.0.1\n\n",
			targetArg,
		)
		return
	}

	// Get use pass and ip from the split up target argument
	user, pass, ip := parts[0], parts[1], parts[2]

	// load the current config into the Config struct
	var config Config
	LoadConfig(&config)

	// Add the new target to the Config struct
	config.Targets = append(config.Targets, Target{user, pass, ip})

	// Write the new Config struct to Config.json
	Write_config(config)

	fmt.Printf(
		"\n[Orchestrate]\n\tAdded target to config:\n\tUser: %s\n\tPass: %s\n\tIP:   %s\n\n",
		user, pass, ip,
	)
}

/*
* Removes a target from the config by IP
 */
func RemoveTarget(targetIP string) {
	// Load the current config into the Config struct
	var config Config
	LoadConfig(&config)

	// Loop throgh all targets, checking if IP matches target's IP to be removed
	for i := range config.Targets {
		if config.Targets[i].IP == targetIP {
			// Remove the target from the config
			config.Targets = append(config.Targets[:i], config.Targets[i+1:]...)
			// Write the new config to Config.json
			Write_config(config)
			fmt.Printf(
				"\n[Orchestrate]\n\tRemoved target %s from Config.json\n\n",
				targetIP,
			)
			return
		}
	}

	// If nothing was found, specified IP does not match any targets in Config.json
	fmt.Printf(
		"\n[Orchestrate]\n\tIP: %s was not found in Config.json\n\n",
		targetIP,
	)
}

/*
* loads the config into a Config struct passed into the function as a pointer
 */
func LoadConfig(config *Config) {
	// Read the file
	configFile, err := os.ReadFile("Config.json")
	if err != nil {
		log.Fatal("Error reading Config.json:", err)
	}

	// Parse the JSON bytes into a GO struct (Config)
	err = json.Unmarshal(configFile, config)
	if err != nil {
		log.Fatal("Error parsing Config.json:", err)
	}
}

func Write_config(config Config) {
	// Parse the config struct into a JSON byte slice
	jsonData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatal("Error writing to Config.json:", err)
	}

	// Write new config to Config.json
	err = os.WriteFile("Config.json", jsonData, 0644)
	if err != nil {
		log.Fatal("Error writing to Config.json:", err)
	}
}
