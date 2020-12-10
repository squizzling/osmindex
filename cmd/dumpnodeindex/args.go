package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/squizzling/osmindex/internal/args"
	"github.com/squizzling/osmindex/internal/system"
)

type Opts struct {
	InputFile      string  `short:"i" long:"input-file"      description:"Input NODEIDX file"                    required:"true"`
	PrintLocations bool    `short:"p" long:"print-locations" description:"Print locations in each index element"                `
	SearchIndex    bool    `short:"s" long:"search-index"    description:"Perform search for each index"                        `
	SubIndexIndex  *uint64 `          long:"sub-index-index" description:"Zero based"                                           `
	SubIndexCount  *uint64 `          long:"sub-index-count"                                                                    `
	StartElement   *int    `          long:"start-element"                                                                      `
	EndElement     *int    `          long:"end-element"                                                                        `
	system.SystemOpts
}

func parseArgs(commandLine []string) *Opts {
	var opts Opts

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.LongDescription = `
Prints each index element stored in the index, optionally printing each location, and
performing a search of the full index for each ID, and looking for inconsistencies.

This is mainly a scratch application for debugging the index.`

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

	errors = append(errors, opts.SystemOpts.Validate()...)

	if (opts.SubIndexCount != nil) != (opts.SubIndexIndex != nil) {
		errors = append(errors, "sub-index-count and sub-index-index must be specified together")
	}

	if opts.SubIndexCount != nil && opts.SubIndexIndex != nil && *opts.SubIndexIndex >= *opts.SubIndexCount {
		errors = append(errors, "sub-index-index must be below sub-index-count")
	}

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
