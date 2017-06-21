package log

import (
	"fmt"
	"os"
)

type Logger struct {
	Verbose bool
	Exiter  Exiter
}

type Exiter func(code int)

// NewLogger creates a new logger with
func NewLogger(verbose bool, exiter func(code int)) *Logger {
	if exiter == nil {
		exiter = func(code int) {
			os.Exit(code)
		}
	}

	return &Logger{
		Verbose: verbose,
		Exiter:  exiter,
	}
}

// Printf outputs a formatted message if this logger has verbose enabled
func (logger Logger) Printf(msg string, args ...interface{}) {
	if logger.Verbose {
		logger.Output(msg, args...)
	}
}

// Error outputs a formatted error message. Errors are always shown even if verbose is disabled.
func (logger Logger) Error(err error) {
	logger.Output(err.Error())
}

// Fatal outputs a formatted error message then exits with status 1. Errors are always shown.
func (logger Logger) Fatal(err error) {
	fmt.Printf(err.Error() + "\n")
	logger.Exiter(1)
}

// Output outputs a message by-passing verbose
func (logger Logger) Output(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}
