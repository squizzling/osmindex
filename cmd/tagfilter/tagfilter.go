package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/squizzling/osmindex/cmd/tagfilter/tagfilterer"
	wcmd "github.com/squizzling/osmindex/cmd/tagfilter/tagfilterer/widget"
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

	tf := &tagfilterer.TagFilterer{}

	mainUI := ui.NewUI(opts.Refresh, cancel)

	pfs := pbf.ProcessFileAsync(
		ctx,
		opts.InputFile,
		tf.Worker(),
		tf.Writer(opts.OutputFile),
		opts.Workers,
	)
	mainUI.ReplaceWidgets(
		wui.NewProcessFileStateWidget("TagFilter", pfs),
		wui.NewWidgetCmd(tf, wcmd.NewWidgetStateWorker, wcmd.NewWidgetStateWriter),
	)
	pfs.Wait()
	mainUI.Stop()

	fmt.Printf("Total run time: %s\n", time.Since(startTime))
}
