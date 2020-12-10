package relindexerinit

import (
	"context"

	"github.com/squizzling/osmindex/cmd/relindex/relindexerinit/state"
	"github.com/squizzling/osmindex/internal/iio"
	"github.com/squizzling/osmindex/internal/pbf"
)

func (ni *RelIndexerInitializer) Output(outputFileName string) func(ctx context.Context, chIn <-chan pbf.Block) {
	return func(ctx context.Context, chIn <-chan pbf.Block) {
		outputFileIdx := iio.OsCreate(outputFileName)

		ni.writer.SetCurrentState(state.WriterReading)
		for data := range chIn {
			if data.BlobType != "RelationData" {
				continue
			}
			ni.writer.SetStateWriting(data.Index)
			iio.FWriteSliceI64(outputFileIdx, data.Data.([]int64))
			ni.writer.SetCurrentState(state.WriterReading)
		}

		iio.FClose(outputFileIdx)
	}
}
