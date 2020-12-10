package decompressor

import (
	"context"

	"github.com/squizzling/osmindex/cmd/decompress/decompressor/state"
	"github.com/squizzling/osmindex/internal/iio"
	"github.com/squizzling/osmindex/internal/pbf"
)

func (decomp *Decompressor) Output(outputFileName string) func(ctx context.Context, chIn <-chan pbf.Block) {
	return func(ctx context.Context, chIn <-chan pbf.Block) {
		outputFile := iio.OsCreate(outputFileName)
		decomp.writer.SetCurrentState(state.WriterReceivingBlock)

		for obo := range chIn {
			decomp.writer.SetWritingBlock(obo.Index)
			iio.FWrite(outputFile, obo.Data.([]byte))
			decomp.writer.SetCurrentState(state.WriterReceivingBlock)
		}

		iio.FClose(outputFile)
	}
}
