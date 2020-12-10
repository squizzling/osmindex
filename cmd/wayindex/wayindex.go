package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/squizzling/osmindex/cmd/wayindex/wayindexer"
	wi "github.com/squizzling/osmindex/cmd/wayindex/wayindexer/widget"
	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/system"
	"github.com/squizzling/osmindex/internal/ui"
	wui "github.com/squizzling/osmindex/internal/ui/widget"
)

func main() {
	opts := parseArgs(os.Args[1:])
	startTime := time.Now()

	teardown := system.Setup(&opts.SystemOpts)
	defer teardown()

	ctx, cancel := context.WithCancel(context.Background())

	mainUI := ui.NewUI(opts.Refresh, cancel)

	ni := &wayindexer.WayIndexer{}
	pfs := pbf.ProcessFileAsync(
		ctx,
		opts.InputFile,
		ni.Worker(opts.Alignment),
		ni.Output(opts.OutputFile, opts.Alignment),
		opts.Workers,
	)

	mainUI.ReplaceWidgets(
		wui.NewProcessFileStateWidget("WayIndex", pfs),
		wui.NewWidgetCmd(ni, wi.NewWidgetStateWorker, wi.NewWidgetStateWriter),
	)
	pfs.Wait()
	mainUI.Stop()
	fmt.Printf("Total run time: %s\n", time.Since(startTime))
}
