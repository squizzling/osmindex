package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/squizzling/osmindex/cmd/decompress/decompressor"
	wcmd "github.com/squizzling/osmindex/cmd/decompress/decompressor/widget"
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

	decomp := &decompressor.Decompressor{}
	pfs := pbf.ProcessFileAsync(
		ctx,
		opts.InputFile,
		decomp.Worker,
		decomp.Output(opts.OutputFile),
		opts.Workers,
	)

	mainUI.ReplaceWidgets(
		wui.NewProcessFileStateWidget("Decompress", pfs),
		wui.NewWidgetCmd(decomp, wcmd.NewWidgetStateWorker, wcmd.NewWidgetStateWriter),
	)
	pfs.Wait()
	mainUI.Stop()
	fmt.Printf("Total run time: %s\n", time.Since(startTime))
}
