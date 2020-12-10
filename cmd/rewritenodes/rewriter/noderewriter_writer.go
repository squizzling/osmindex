package rewriter

import (
	"context"
	"encoding/binary"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/cmd/rewritenodes/rewriter/state"
	"github.com/squizzling/osmindex/internal/iio"
	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/t"
)

func (nr *NodeRewriter) Writer(outputFilename string) pbf.WriteFunc {
	return func(ctx context.Context, chIn <-chan pbf.Block) {
		outputFile := iio.OsCreate(outputFilename)

		blobHeaderBuf := make([]byte, 0, 64)
		for data := range chIn {
			nr.writer.SetStateWriting(data.Index)
			blobBuf, ok := data.Data.([]byte)
			if !ok {
				blobBuf = data.Data.(mmap.MMap)
			}

			// prepare blob header
			var bh t.BlobHeader
			bh.DataSize = int32(len(blobBuf))
			bh.Type = data.BlobType
			blobHeaderBuf = blobHeaderBuf[:4]
			blobHeaderBuf = (&bh).Write(blobHeaderBuf)

			// prepare blob header header
			binary.BigEndian.PutUint32(blobHeaderBuf, uint32(len(blobHeaderBuf)-4))

			iio.FWrite(outputFile, blobHeaderBuf)
			iio.FWrite(outputFile, blobBuf)
			nr.writer.SetCurrentState(state.WriterReading)
		}

		iio.FClose(outputFile)
	}
}
