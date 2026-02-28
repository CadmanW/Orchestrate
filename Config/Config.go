package config

import (
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
		log.Fatalf("\n[Orchestrate] config\n\tError when parsing args: %s\n\n", err)
	}

	// ADD TARGET
	if flag_add != "" {
		err := AddTarget(flag_add)
		if err != nil {
			log.Fatalf("\n[Orchestrate] config\n\tError adding target to Config.json: %s\n\n", err)
		}
	}

	// REMOVE TARGET
	if flag_remove != "" {
		err := RemoveTarget(flag_remove)
		if err != nil {
			log.Fatalf("\n[Orchestrate] config\n\tError removing target from Config.json: %s\n\n", err)
		}
	}
}

/*
* Adds a target to the config
 */
func AddTarget(targetArg string) error {
	// Normalize the target argument so we can split it up
	// user:pass@ip --> user:pass:ip
	normalized := strings.ReplaceAll(targetArg, "@", ":")

	// Split up the target argument
	// "user:pass:ip" --> ["user", "pass", "ip"]
	parts := strings.Split(normalized, ":")

	// Make sure user:pass@ip format was followed
	if len(parts) != 3 {
		return fmt.Errorf("Invalid target format: %s ; Use: username:password@127.0.0.1", targetArg)
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
		"\n[Orchestrate] config\n\tAdded target to config:\n\tUser: %s\n\tPass: %s\n\tIP:   %s\n\n",
		user, pass, ip,
	)

	// return nil on success
	return nil
}

/*
* Removes a target from the config by IP
 */
func RemoveTarget(targetIP string) error {
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
				"\n[Orchestrate] config\n\tRemoved target %s from Config.json\n\n",
				targetIP,
			)
			// return nil on success
			return nil
		}
	}

	// If nothing was found, specified IP does not match any targets in Config.json
	return fmt.Errorf("IP: %s was not found in Config.json", targetIP)
}

/*
* loads the config into a Config struct passed into the function as a pointer
 */
func LoadConfig(config *Config) error {
	// Read the file
	configFile, err := os.ReadFile("Config.json")
	if err != nil {
		return fmt.Errorf("Error reading Config.json: %s", err)
	}

	// Parse the JSON bytes into a GO struct (Config)
	err = json.Unmarshal(configFile, config)
	if err != nil {
		return fmt.Errorf("Error unmarshaling Config.json: %s", err)
	}

	// Return nil on success
	return nil
}

/*
* Writes the config to Config.json
 */
func Write_config(config Config) error {
	// Parse the config struct into a JSON byte slice
	jsonData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("Error writing to Config.json: %s", err)
	}

	// Write new config to Config.json
	err = os.WriteFile("Config.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to Config.json: %s", err)
	}

	// Return nil on success
	return nil
}

/*
* Returns targets based on the flags provided
 */
func GetTargets(t string, T string, a bool) []Target {
	var targets []Target
	// Load the config into the Config struct
	var conf Config
	LoadConfig(&conf)

	// Check for arg -t
	if t != "" {
		// Iterate through conf.Targets, checking if it has the IP specified
		for i := range conf.Targets {
			if t == conf.Targets[i].IP {
				targets = append(targets, conf.Targets[i])
			}
		}
	}

	// Check for arg -T
	if T != "" {

		// Split the arg string into an array of IPs
		targetIPs := strings.Split(T, " ")

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
	if a {
		targets = append(targets, conf.Targets...)
	}

	return targets
}
