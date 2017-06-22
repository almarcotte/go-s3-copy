package log

import (
	"fmt"
	"os"
)

type StdOut struct {
	Verbose bool
	Exiter  Exiter
}

type Exiter func(code int)

// NewLogger creates a new logger with
func NewLogger(verbose bool, exiter func(code int)) StdOut {
	if exiter == nil {
		exiter = func(code int) {
			os.Exit(code)
		}
	}

	return StdOut{
		Verbose: verbose,
		Exiter:  exiter,
	}
}

// Printf outputs a formatted message if this logger has verbose enabled
func (logger StdOut) Printf(msg string, args ...interface{}) {
	if logger.Verbose {
		logger.Output(msg, args...)
	}
}

// Error outputs a formatted error message. Errors are always shown even if verbose is disabled.
func (logger StdOut) Error(err error) {
	logger.Output(err.Error())
}

// Fatal outputs a formatted error message then exits with status 1. Errors are always shown.
func (logger StdOut) Fatal(err error) {
	fmt.Printf(err.Error() + "\n")
	logger.Exiter(1)
}

// Output outputs a message by-passing verbose
func (logger StdOut) Output(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}
