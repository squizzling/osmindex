package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"

	"github.com/squizzling/osmindex/internal/args"
	"github.com/squizzling/osmindex/internal/system"
	"github.com/squizzling/osmindex/internal/ui"
)

type Opts struct {
	InputFile  string `short:"i" long:"input-file"  description:"Input PBF file"    required:"true"`
	OutputFile string `short:"o" long:"output-file" description:"Output PBF file"   required:"true"`
	Workers    int    `short:"w" long:"workers"     description:"Number of workers"`
	ui.UIOpts
	system.SystemOpts
}

func parseArgs(commandLine []string) *Opts {
	var opts Opts

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.LongDescription = ``
	opts.Workers = runtime.NumCPU()

	// Main arg parsing
	positional, err := parser.ParseArgs(commandLine)
	if err != nil {
		if !args.IsHelp(err) {
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

	errors = append(errors, opts.UIOpts.Validate()...)
	errors = append(errors, opts.SystemOpts.Validate()...)

	// Validation results
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
