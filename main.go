// Package main is the entry point of the application
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gnumast/gossip/config"
	"github.com/gnumast/gossip/watcher"
	"os"
)

const (
	ENV_CONFIG = "GOSSIP_CONFIG"
	Version    = "0.0.1"
)

func main() {
	file, err := getConfigFileLocation()

	if err != nil {
		usage()
		os.Exit(1)
	}

	config, err = config.LoadConfig(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// usage displays the usage statement
func usage() {
	fmt.Printf("gossip %v\n", Version)
	fmt.Println("usage: gossip [-h] [-config file]")
	fmt.Println(" -h: displays this")
	fmt.Println(" -config file: configuration file in json")
}

// getConfigFileLocation reads the arguments from the command line and returns the location of the config file
func getConfigFileLocation() (file string, err error) {
	file = *flag.String("config", "", "path to configuration file")
	flag.Parse()

	if file != "" {
		return
	}

	file = os.Getenv(ENV_CONFIG)

	if file != "" {
		return
	}

	err = errors.New(fmt.Sprintf("Configuration file is missing. Use -config or set $%v", ENV_CONFIG))

	return
}
