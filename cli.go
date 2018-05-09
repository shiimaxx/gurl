package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (c *CLI) Run(args []string) int {
	var (
		output  string
		version bool
	)

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(c.outStream)

	flags.StringVar(&output, "output", "./", "output file")
	flags.StringVar(&output, "o", "./", "output file(Short)")
	flags.BoolVar(&version, "version", false, "print version information")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if version {
		fmt.Fprintf(c.outStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	if len(flags.Args()) < 1 {
		fmt.Fprintln(c.errStream, "missing arguments")
		return ExitCodeError
	}

	filePath := flags.Args()[0]
	if _, err := os.Stat(filePath); os.IsExist(err) {
		fmt.Fprintf(c.errStream, "%s: already exits\n", filePath)
		return ExitCodeError
	}

	return ExitCodeOK
}
