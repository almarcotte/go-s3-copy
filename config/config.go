package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Path struct {
	Root      string        `json:"root"`
	Bucket    string        `json:"bucket"`
	Recursive bool          `json:"recursive"`
	Delete    bool          `json:"delete"`
	Delay     time.Duration `json:"delay"`
}

type Global struct {
	Bucket string        `json:"bucket"`
	Delay  time.Duration `json:"delay"`
}

type Credentials struct {
	Access string `json:"access"`
	Secret string `json:"secret"`
	Region string `json:"region"`
}

type Parsed struct {
	Paths       []Path      `json:"paths"`
	Global      Global      `json:"global"`
	Credentials Credentials `json:"credentials"`
	Verbose     bool        `json:"verbose"`
}

var (
	DelayTooShortError       = errors.New("Delay should be at least")
	AwsSecretKeyMissingError = errors.New("AWS secret key is not provided in your config file and $AWS_SECRET is not set.")
	AwsAccessKeyMissingError = errors.New("AWS access key is not provided in your config file and $AWS_ACCESS is not set.")
	AwsRegionMissingError    = errors.New("AWS Region is missing")
	NoPathToWatchError       = errors.New("There must be at least one path to watch")
	MissingRootPathError     = errors.New("Missing `root` for at least one path")
	NoBucketDefinedError     = errors.New("A path doesn't have an explicit bucket and none is set globally")
)

// LoadConfig reads the JSON file and creates a Parsed. Global settings are applied to all paths if none are returned.
func LoadConfig(filename string) (config *Parsed, err error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return
	}

	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return
	}

	mergeGlobals(config)

	if err = validate(*config); err != nil {
		return
	}

	return
}

func validate(config Parsed) error {
	// Paths
	if len(config.Paths) == 0 {
		return NoPathToWatchError
	}

	for _, path := range config.Paths {
		if path.Root == "" {
			return MissingRootPathError
		}

		if path.Bucket == "" {
			return NoBucketDefinedError
		}

		if path.Delay == 0 {
			return DelayTooShortError
		}
	}

	// Credentials
	if config.Credentials.Secret == "" && os.Getenv("AWS_SECRET") == "" {
		return AwsSecretKeyMissingError
	}

	if config.Credentials.Access == "" && os.Getenv("AWS_ACCESS") == "" {
		return AwsAccessKeyMissingError
	}

	if config.Credentials.Region == "" && os.Getenv("AWS_REGION") == "" {
		return AwsRegionMissingError
	}

	return nil
}

// MergeGlobals takes the global parameters for Delay and Bucket and applies them to Paths that don't have
// an explicit value
func mergeGlobals(config *Parsed) {
	for i := range config.Paths {
		current := &config.Paths[i]

		if current.Delay == 0 {
			config.Paths[i].Delay = config.Global.Delay
		}

		if current.Bucket == "" {
			current.Bucket = config.Global.Bucket
		}
	}
}
