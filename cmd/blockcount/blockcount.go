package main

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/system"
)

func countBlocks(blockCount *uint64) pbf.WorkFunc {
	return func(ctx context.Context, chIn <-chan pbf.Block, chOut chan<- pbf.Block) {
		localBlockCount := uint64(0)
		for range chIn {
			localBlockCount++
		}
		atomic.AddUint64(blockCount, localBlockCount)
	}
}

func main() {
	opts := parseArgs(os.Args[1:])
	startTime := time.Now()

	teardown := system.Setup(&opts.SystemOpts)
	defer teardown()

	var blockCount uint64
	pbf.ProcessFileAsync(
		context.Background(),
		opts.InputFile,
		countBlocks(&blockCount),
		nil,
		opts.Workers,
	).Wait()

	fmt.Printf("Counted %d blocks\n", atomic.LoadUint64(&blockCount))
	fmt.Printf("Total run time: %s\n", time.Since(startTime))
}
