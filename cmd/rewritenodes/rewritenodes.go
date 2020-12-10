package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/squizzling/osmindex/cmd/rewritenodes/rewriter"
	wcmd "github.com/squizzling/osmindex/cmd/rewritenodes/rewriter/widget"
	"github.com/squizzling/osmindex/internal/iio"
	"github.com/squizzling/osmindex/internal/system"
	"github.com/squizzling/osmindex/internal/ui"
	wui "github.com/squizzling/osmindex/internal/ui/widget"

	"github.com/squizzling/osmindex/internal/nodeindex"
	"github.com/squizzling/osmindex/internal/pbf"
)

func makePassFilename(fn string, pass, count uint64) string {
	return fmt.Sprintf("%s.pass_%03d_%03d", fn, pass, count)
}

func main() {
	opts := parseArgs(os.Args[1:])

	startTime := time.Now()

	teardown := system.Setup(&opts.SystemOpts)
	defer teardown()

	ctx, cancel := context.WithCancel(context.Background())

	ix := nodeindex.NewNodeIndex(opts.IndexFile, opts.Timing)

	passes := uint64(1)

	if opts.Memory != nil {
		memTarget := *opts.Memory * 1048576
		memRequired := ix.LocationCount() * 8
		passes = (memRequired / memTarget) + 1
	}

	nr := &rewriter.NodeRewriter{}

	mainUI := ui.NewUI(opts.Refresh, cancel)
	var passTimes []time.Duration

	for i := uint64(0); i < passes; i++ {
		passStartTime := time.Now()

		if i > 1 {
			iio.OsRemoveTry(makePassFilename(opts.OutputFile, i-1, passes))
		}

		inputFile := opts.InputFile
		if i > 0 {
			inputFile = makePassFilename(opts.OutputFile, i, passes)
		}

		outputFile := opts.OutputFile
		if i+1 < passes {
			outputFile = makePassFilename(opts.OutputFile, i+1, passes)
		}

		pfs := pbf.ProcessFileAsync(
			ctx,
			inputFile,
			nr.Worker(ix.SubIndex(i, passes), i == passes-1),
			nr.Writer(outputFile),
			opts.Workers,
		)
		mainUI.ReplaceWidgets(
			wui.NewPassProgress(i, passes, pfs),
			wui.NewProcessFileStateWidget("RewriteNodes", pfs),
			wui.NewWidgetCmd(nr, wcmd.NewWidgetStateWorker, wcmd.NewWidgetStateWriter),
		)
		pfs.Wait()
		passTimes = append(passTimes, time.Since(passStartTime))
	}

	if passes > 1 {
		iio.OsRemove(makePassFilename(opts.OutputFile, passes-1, passes))
	}

	mainUI.Stop()

	ix.PrintStats()

	ix.Close()
	for idx, passRunTime := range passTimes {
		fmt.Printf("Pass %d run time: %s\n", idx, passRunTime)
	}
	fmt.Printf("Total run time: %s\n", time.Since(startTime))
}
