package indexer

import (
	"context"
	"io"

	"github.com/squizzling/osmindex/cmd/nodeindex/indexer/state"
	"github.com/squizzling/osmindex/internal/icontext"
	"github.com/squizzling/osmindex/internal/iio"
	"github.com/squizzling/osmindex/internal/nodeindex"
	"github.com/squizzling/osmindex/internal/pbf"
)

func (ni *NodeIndexer) Output(outputFileName string) func(ctx context.Context, chIn <-chan pbf.Block) {
	return func(ctx context.Context, chIn <-chan pbf.Block) {
		outputFileIdx := iio.OsCreate(outputFileName)
		outputFileLoc := iio.OsCreate(outputFileName + ".loc.tmp")

		// Reserve 16 bytes, to be filled in at the end
		iio.FWriteSliceI64(outputFileIdx, []int64{0, 0})

		var outputData []uint64
		lastLocationIndex := uint64(0)
		indexElementCount := uint64(0)

		ni.writer.SetCurrentState(state.WriterReading)
		for data := range chIn {
			if data.BlobType != "IndexData" {
				continue
			}
			outputData = outputData[:0]
			element := data.Data.(elementBatch)
			ni.writer.SetStateWritingIndex(uint64(len(element.ranges)))
			for _, indexBlock := range element.ranges {
				itemCount := (indexBlock.last - indexBlock.first) + 1
				outputData = append(
					outputData,
					indexBlock.first,
					uint64(nodeindex.NewCountOffset(itemCount, lastLocationIndex)),
				)
				indexElementCount++
				lastLocationIndex += itemCount
			}
			iio.FWriteSliceU64(outputFileIdx, outputData)
			ni.writer.SetCurrentState(state.WriterWritingLocations)
			iio.FWriteSliceI64(outputFileLoc, element.locations)
			ni.writer.SetCurrentState(state.WriterReading)
		}

		// Append location file to end of index file
		ni.writer.SetStateCopyLocations(8 * lastLocationIndex)
		iio.FSeek(outputFileLoc, 0, io.SeekStart)
		buf := make([]byte, 1048576)
		for {
			if icontext.IsCancelled(ctx) {
				break
			}
			n := iio.FCopy(outputFileIdx, outputFileLoc, buf)
			if n == 0 {
				break
			}
			ni.writer.AdvanceCopy(uint64(n))
		}
		iio.FClose(outputFileLoc)
		iio.OsRemove(outputFileName + ".loc.tmp")

		// Put the index block count at the start of the index file
		ni.writer.SetCurrentState(state.WriteFinalizing)
		iio.FSeek(outputFileIdx, 0, io.SeekStart)
		iio.FWriteSliceU64(outputFileIdx, []uint64{nodeindex.Version, indexElementCount})
		iio.FClose(outputFileIdx)
	}
}
