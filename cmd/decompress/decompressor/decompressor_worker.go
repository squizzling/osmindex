package decompressor

import (
	"context"
	"encoding/binary"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/cmd/decompress/decompressor/state"
	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/pool"
	"github.com/squizzling/osmindex/internal/t"
)

func (decomp *Decompressor) Worker(ctx context.Context, chIn <-chan pbf.Block, chOut chan<- pbf.Block) {
	sw := &state.Worker{}
	decomp.TrackWorker(sw)
	defer decomp.UntrackWorker(sw)

	var pbr t.PBReader

	for obh := range chIn {
		// read
		sw.SetDecodingBlob(obh.Index)

		var b t.Blob
		pbr.ReadBlob(obh.Data.(mmap.MMap), &b)

		// decompress
		sw.SetCurrentState(state.WorkerDecompressingBlob)
		decompData := b.GetRawData()

		// re-encode
		sw.SetCurrentState(state.WorkerEncodingBlob)
		var newBlob t.Blob
		newBlob.Raw = decompData.Buffer
		blobBuffer := newBlob.Write(nil)

		var bh t.BlobHeader
		bh.DataSize = int32(len(blobBuffer))
		bh.Type = obh.BlobType
		outputBuffer := bh.Write(make([]byte, 4))
		binary.BigEndian.PutUint32(outputBuffer, uint32(len(outputBuffer)-4))
		outputBuffer = append(outputBuffer, blobBuffer...)

		// release
		pool.ByteSlice.Put(decompData)

		// write
		sw.SetCurrentState(state.WorkerSendingBlock)
		chOut <- pbf.Block{
			Index:    obh.Index,
			Data:     outputBuffer,
			BlobType: obh.BlobType,
		}

		// read
		sw.SetCurrentState(state.WorkerReceivingBlock)
	}
}
