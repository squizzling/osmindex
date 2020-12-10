package wayindexer

import (
	"context"
	"io"

	"github.com/squizzling/osmindex/cmd/wayindex/wayindexer/state"
	"github.com/squizzling/osmindex/internal/icontext"
	"github.com/squizzling/osmindex/internal/iio"
	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/wayindex"
)

func (ni *WayIndexer) Output(outputFileName string, alignmentShift int) func(ctx context.Context, chIn <-chan pbf.Block) {
	return func(ctx context.Context, chIn <-chan pbf.Block) {
		locationAlignmentShift := alignmentShift

		outputFileIdx := iio.OsCreate(outputFileName)
		outputFileLoc := iio.OsCreate(outputFileName + ".loc.tmp")

		// Reserve 16 bytes.  Don't need nearly this much for the header, but it keeps things in cache lines
		iio.FWriteSliceI64(outputFileIdx, []int64{0, 0})

		var outputData []uint64
		locationLength := int64(0)

		indexBlockCount := int64(0)

		ni.writer.SetCurrentState(state.WriterReading)
		for data := range chIn {
			if data.BlobType != "IndexData" {
				continue
			}
			outputData = outputData[:0]
			element := data.Data.(indexElement)
			ni.writer.SetStateWritingIndex(uint64(len(element.blocks)))

			for _, b := range element.blocks {
				// On input, Offset will be chunks from the start of element.locations
				// On output, Offset needs to be chunks from the start of the file
				globalData := wayindex.MakeIndexData(
					b.WayID(),
					b.Offset()+locationLength,
				)
				outputData = append(outputData, uint64(globalData))
			}
			iio.FWriteSliceU64(outputFileIdx, outputData)
			iio.FWrite(outputFileLoc, element.locations)
			locationLength += int64(len(element.locations)) >> locationAlignmentShift
			indexBlockCount += int64(len(element.blocks))
			ni.writer.SetCurrentState(state.WriterReading)
		}

		// Append location file to end of index file
		// TODO: Put location file at the start
		ni.writer.SetStateCopyLocations(locationLength << locationAlignmentShift)
		iio.FSeek(outputFileLoc, 0, io.SeekStart)
		buf := make([]byte, 1048576)
		for {
			if icontext.IsCancelled(ctx) {
				// TODO: Check if this should be done less frequently.
				break
			}
			// TODO: Make an iio.FCopyWithContext
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
		// TODO: Version
		iio.FWriteSliceI64(outputFileIdx, []int64{indexBlockCount, int64(alignmentShift)})
		iio.FClose(outputFileIdx)
	}
}
