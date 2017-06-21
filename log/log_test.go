package log

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewLogger(t *testing.T) {
	l := NewLogger(true, nil)

	if !l.Verbose {
		t.Fatalf("Unexpected value for verbose. Expected os.Stdout, got %v", l.Verbose)
	}
}

func ExampleLogger_Output() {
	verboseLogger := NewLogger(true, nil)
	silentLogger := NewLogger(false, nil)

	verboseLogger.Output("Hello, %s! This is %v.", "World", true)
	silentLogger.Output("Hello, %s! This is %v.", "World", true)

	// Output:
	// Hello, World! This is true.
	// Hello, World! This is true.
}

func ExampleLogger_Printf() {
	verboseLogger := NewLogger(true, nil)
	silentLogger := NewLogger(false, nil)

	verboseLogger.Printf("Hello, %s! This is %v.", "World", true)
	silentLogger.Printf("Hello, %s! This is %v.", "World", true)

	// Output:
	// Hello, World! This is true.
}

func ExampleLogger_Error() {
	verboseLogger := NewLogger(true, nil)
	silentLogger := NewLogger(false, nil)

	err := errors.New("This is an error, very sorry!")

	verboseLogger.Error(err)
	silentLogger.Error(err)
	// Output:
	// This is an error, very sorry!
	// This is an error, very sorry!
}

func ExampleLogger_Fatal() {
	verboseLogger := NewLogger(true, func(code int) { fmt.Printf("Return code: %v\n", code) })
	silentLogger := NewLogger(false, func(code int) { fmt.Printf("Return code: %v\n", code) })

	err := errors.New("Something has gone horribly wrong")

	verboseLogger.Fatal(err)
	silentLogger.Fatal(err)

	// Output:
	// Something has gone horribly wrong
	// Return code: 1
	// Something has gone horribly wrong
	// Return code: 1
}
