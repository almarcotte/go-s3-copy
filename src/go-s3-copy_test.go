package src

import (
	"github.com/gnumast/go-s3-copy/config"
	"github.com/gnumast/go-s3-copy/log"
	"testing"
)

func TestNewWatcher(t *testing.T) {
	c := &config.Parsed{
		Verbose: true,
		Global: config.Global{
			Delay:  5,
			Bucket: "my-bucket",
		},
	}

	l := log.NewLogger(true, nil)
	w := NewWatcher(c, l)

	if !w.Logger.Verbose {
		t.Fatalf("Couldn't attach a logger to the Watcher")
	}

	if w.Config.Global.Delay != 5 {
		t.Fatalf("Couldn't attach a configuration to the Watcher")
	}

	if len(w.Queue) != 0 {
		t.Fatalf("Couldn't initialize the upload queue for this Watcher")
	}
}
