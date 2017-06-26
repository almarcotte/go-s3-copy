// Package main is the entry point of the application
package main

import (
	"flag"
	"fmt"
	"github.com/gnumast/go-s3-copy/config"
	"github.com/gnumast/go-s3-copy/log"
	"github.com/gnumast/go-s3-copy/src"
	"os"
)

const (
	ENV_CONFIG = "GO_S3_COPY_CONFIG"
	Version    = "0.0.1"
)

func main() {
	var file string

	file = os.Getenv(ENV_CONFIG)                                      // Default to the environment variable
	flag.StringVar(&file, "config", "", "path to configuration file") // overwrite if passed explicitly

	flag.Parse()

	if file == "" {
		usage()
		os.Exit(1)
	}

	parsed, err := config.LoadConfig(file)
	logger := log.NewLogger(parsed.Verbose, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	watcher := src.NewWatcher(parsed, logger)

	if err = watcher.Start(); err != nil {
		return
	}
}

// usage displays the usage statement
func usage() {
	fmt.Printf("go-s3-copy %v\n", Version)
	fmt.Println("usage: go-s3-copy [-h] [-config file]")
	fmt.Println(" -h: displays this")
	fmt.Println(" -config file: configuration file in json")
}
