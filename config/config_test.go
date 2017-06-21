package config

import (
	"io/ioutil"
	"os"
	"testing"
)

const testContent = `
{
  "paths": [
    {
      "root": "/home/alex/photos",
      "delete": true,
      "delay": 2
    },
    {
      "root": "/home/alex/music",
      "recursive": true,
      "bucket": "my-music-backup"
    }
  ],
  "global": {
    "delete": false,
    "recursive": true,
    "bucket": "default-bucket",
    "delay": 5
  },
  "credentials": {
    "access": "AWS_ACCESS_KEY",
    "secret": "AWS_SECRET_KEY"
  }
}
`

func TestFileIsMissing(t *testing.T) {
	_, err := LoadConfig("testing.json")

	if err == nil {
		t.Fatal("Expecting `LoadConfig` to return an error, nil returned")
	}
}

func TestJsonFormatIssue(t *testing.T) {
	content := []byte("This is not real JSON")
	filename := "config.json"

	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		t.Fatalf("Error writing content to %s: %v", filename, err)
	}

	defer os.Remove(filename)

	_, err := LoadConfig(filename)

	if err == nil {
		t.Fatalf("Expecting `LoadConfig` to return an error on invalid JSON, `nil` returned")
	}
}

// TestLoadConfig tests both loading the file (i.e. reading from disk) but also that the returned Config struct has the
// expected values including the global params being properly applied
func TestLoadConfig(t *testing.T) {
	content := []byte(testContent)
	filename := "config.json"

	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		t.Fatalf("Error writing content to %s: %v", filename, err)
	}

	defer os.Remove(filename)

	config, err := LoadConfig("config.json")

	if err != nil {
		t.Fatalf("Error while loading sample config: %v", err)
	}

	assert := func(expected interface{}, actual interface{}, msg string) {
		if expected != actual {
			t.Fatalf(msg+" - got: %+v, expected: %+v", actual, expected)
		}
	}

	assert(true, config.Global.Recursive, "Global.Recursive")
	assert(false, config.Global.Delete, "Global.Delete")
	assert("default-bucket", config.Global.Bucket, "Global.Bucket")

	assert("AWS_SECRET_KEY", config.Credentials.Secret, "Credentials.Secret")
	assert("AWS_ACCESS_KEY", config.Credentials.Access, "Credentials.Access")

	assert(2, len(config.Paths), "len(Paths)")

	assert(true, config.Paths[0].Delete, "Path[0].Delete")
	assert(false, config.Paths[1].Delete, "Path[1].Delete")

	assert(true, config.Paths[0].Recursive, "Path[0].Recursive")
	assert(true, config.Paths[1].Recursive, "Path[1].Recursive")
}
