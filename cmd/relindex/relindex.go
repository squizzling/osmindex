package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/squizzling/osmindex/cmd/relindex/relindexerinit"
	wi "github.com/squizzling/osmindex/cmd/relindex/relindexerinit/widget"
	"github.com/squizzling/osmindex/internal/iio"
	"github.com/squizzling/osmindex/internal/iunsafe"
	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/system"
	"github.com/squizzling/osmindex/internal/ui"
	wui "github.com/squizzling/osmindex/internal/ui/widget"
	"github.com/squizzling/osmindex/internal/wayindex"
)

var (
	_ = wi.NewWidgetStateWriter
	_ = wui.NewWidgetCmd
)

func makePassFilename(fn string, pass, count int) string {
	return fmt.Sprintf("%s.pass_%03d_%03d", fn, pass, count)
}

// rewrite is both a description (it rewrites a file) and a requirement (the function needs rewriting)
func rewrite(outputFile, inputFile string, wi *wayindex.WayIndex, pass, count int) {
	swi := wi.SubIndexWithScaledPass(pass, count)

	fIn := iio.OsOpen(inputFile)
	fOut := iio.OsCreate(outputFile)

	m := iio.MMapMap(fIn)
	u := iunsafe.ByteSliceAsInt64Slice(m)

	outputData := make([]int64, 0, 1048576)

	for i := 0; i < len(u); {
		relId := u[i]
		i++
		outputData = append(outputData, relId)

		wayCount := int(u[i])
		i++
		outputData = append(outputData, int64(wayCount))

		for j := 0; j < wayCount; j++ {
			wayId := u[i]
			i++

			if wayId < 0 {
				locations, _, ok := swi.FindLocations(-wayId, 0)
				if !ok {
					if pass+1 < count {
						outputData = append(outputData, wayId)
					} else {
						// if we're in the last pass, and we can't find a way,
						// keep the way but output a 0 length.
						outputData = append(outputData, -wayId, 0)
					}
				} else {
					outputData = append(outputData, -wayId)
					outputData = append(outputData, int64(len(locations)))
					outputData = append(outputData, locations...)
				}
			} else {
				outputData = append(outputData, wayId)
				locCount := int(u[i])
				i++
				outputData = append(outputData, int64(locCount))
				outputData = append(outputData, u[i:i+locCount]...)
				i += locCount
			}

		}

		if len(outputData) > 1000000 { // Try to keep it under 1MB
			iio.FWriteSliceI64(fOut, outputData)
			outputData = outputData[:0]
		}
	}

	if len(outputData) > 0 {
		iio.FWriteSliceI64(fOut, outputData)
	}

	iio.MMapUnmap(m)
	iio.FClose(fOut)
	iio.FClose(fIn)
}

func main() {
	opts := parseArgs(os.Args[1:])
	startTime := time.Now()

	filePass := func(p int) string {
		return makePassFilename(opts.OutputFile, p, opts.Passes)
	}

	teardown := system.Setup(&opts.SystemOpts)
	defer teardown()

	ctx, cancel := context.WithCancel(context.Background())

	mainUI := ui.NewUI(opts.Refresh, cancel)
	var passTimes []time.Duration

	ni := &relindexerinit.RelIndexerInitializer{}
	pfs := pbf.ProcessFileAsync(
		ctx,
		opts.InputFile,
		ni.Worker(),
		ni.Output(filePass(0)),
		opts.Workers,
	)
	mainUI.ReplaceWidgets(
		wui.NewProcessFileStateWidget("RelIndexInit", pfs),
		wui.NewWidgetCmd(ni, wi.NewWidgetStateWorker, wi.NewWidgetStateWriter),
	)
	pfs.Wait()
	mainUI.Stop()

	ix := wayindex.NewWayIndex(opts.WayIndex, nil)
	for i := 0; i < opts.Passes; i++ {
		fmt.Printf("Running pass %d/%d\n", i+1, opts.Passes)
		passStartTime := time.Now()

		inputFile := filePass(i)
		outputFile := opts.OutputFile
		if i+1 < opts.Passes {
			outputFile = filePass(i + 1)
		}

		rewrite(outputFile, inputFile, ix, i, opts.Passes)
		iio.OsRemoveTry(filePass(i))

		passTimes = append(passTimes, time.Since(passStartTime))
	}


	ix.Close()
	ix.PrintStats()
	for idx, passRunTime := range passTimes {
		fmt.Printf("Pass %d run time: %s\n", idx, passRunTime)
	}

	fmt.Printf("Total run time: %s\n", time.Since(startTime))
}
