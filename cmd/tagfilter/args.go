package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"

	"github.com/squizzling/osmindex/internal/system"
	"github.com/squizzling/osmindex/internal/ui"
)

type Opts struct {
	InputFile  string `short:"i" long:"input-file"  description:"Input PBF file"   required:"true"`
	OutputFile string `short:"o" long:"output-file" description:"Output PBF file"  required:"true"`
	Workers    int    `short:"w" long:"workers"     description:"Number of workers               "`
	ui.UIOpts
	system.SystemOpts
}

func parseArgs(args []string) *Opts {
	var opts Opts

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.LongDescription = ``
	opts.Workers = runtime.NumCPU()

	// Main arg parsing
	positional, err := parser.ParseArgs(args)
	if err != nil {
		if !isHelp(err) {
			parser.WriteHelp(os.Stderr)
			_, _ = fmt.Fprintf(os.Stderr, "\n\nerror parsing command line: %v\n", err)
			os.Exit(1)
		}
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	// Arg validation
	var errors []string
	if len(positional) != 0 { // Near as I can tell there's no way to say no positional arguments allowed.
		errors = append(errors, "no positional arguments allowed")
	}

	errors = append(errors, opts.SystemOpts.Validate()...)

	if len(errors) > 0 {
		parser.WriteHelp(os.Stderr)
		_, _ = fmt.Fprintf(os.Stderr, "\n\n")
		for _, err := range errors {
			_, _ = fmt.Fprintf(os.Stderr, "error parsing command line: %s\n", err)
		}
		os.Exit(1)
	}

	return &opts
}

// isHelp is a helper to test the error from ParseArgs() to
// determine if the help message was written. It is safe to
// call without first checking that error is nil.
func isHelp(err error) bool {
	// This was copied from https://github.com/jessevdk/go-flags/blame/master/help.go#L499, as there has not been an
	// official release yet with this code. Renamed from WriteHelp to isHelp, as flags.ErrHelp is still returned when
	// flags.HelpFlag is set, flags.PrintError is clear, and -h/--help is passed on the command line, even though the
	// help is not displayed in such a situation.
	if err == nil { // No error
		return false
	}

	flagError, ok := err.(*flags.Error)
	if !ok { // Not a go-flag error
		return false
	}

	if flagError.Type != flags.ErrHelp { // Did not print the help message
		return false
	}

	return true
}
