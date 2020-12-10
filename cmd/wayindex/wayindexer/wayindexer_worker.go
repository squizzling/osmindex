package wayindexer

import (
	"bytes"
	"context"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/cmd/wayindex/wayindexer/state"
	"github.com/squizzling/osmindex/internal/pb"
	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/pool"
	"github.com/squizzling/osmindex/internal/t"
	"github.com/squizzling/osmindex/internal/wayindex"
)

func padToAlignment(b []byte, amt int) []byte {
	return append(b, bytes.Repeat([]byte{0}, amt-((len(b))%amt))...)
}

func (ni *WayIndexer) Worker(alignmentShift int) func(ctx context.Context, chIn <-chan pbf.Block, chOut chan<- pbf.Block) {
	return func(ctx context.Context, chIn <-chan pbf.Block, chOut chan<- pbf.Block) {
		pbr := &t.PBReader{
			SkipStringTable: true,
			SkipDenseNodes:  true,
			SkipRelations:   true,
		}

		locationAlignmentShift := alignmentShift
		locationAlignment := 1 << locationAlignmentShift

		sw := &state.Worker{}
		ni.TrackWorker(sw)
		defer ni.UntrackWorker(sw)

		sw.SetCurrentState(state.WorkerReading)
		for data := range chIn {
			if data.BlobType != "OSMData" {
				chOut <- data
				continue
			}
			sw.SetStateDecodeBlock(data.Index)

			var b t.Blob
			pbr.ReadBlob(data.Data.(mmap.MMap), &b)
			rawData := b.GetRawData()

			var pblk t.PrimitiveBlock
			pbr.ReadPrimitiveBlock(rawData.Buffer, &pblk)

			sw.SetCurrentState(state.WorkerProcessing)

			var ixData []wayindex.WayIndexData
			var outputLocBytes []byte

			for _, pg := range pblk.PrimitiveGroup {
				for _, w := range pg.Ways {
					ixData = append(ixData, wayindex.MakeIndexData(int64(w.ID), int64(len(outputLocBytes))>>locationAlignmentShift))

					// Filter out sequential duplicates
					outputLocations := make([]int64, len(w.Refs))
					last := int64(0)
					idx := 0
					for _, ref := range w.Refs {
						if ref != last {
							outputLocations[idx] = ref
							last = ref
							idx++
						}
					}

					encodedBytes := pb.EncodeS64PackedDeltaZero(nil, outputLocations[:idx])
					outputLocBytes = append(outputLocBytes, encodedBytes...)
					outputLocBytes = padToAlignment(outputLocBytes, locationAlignment)
				}
			}

			msg := pbf.Block{
				Index: data.Index,
				Data: indexElement{
					blocks:    ixData,
					locations: outputLocBytes,
				},
				BlobType: "IndexData",
			}

			if len(ixData) == 0 {
				// Exactly one thing must be written, because of the re-ordering, so flag
				// the block as something to be ignored.
				msg.BlobType = "Filler"
			}

			pool.ByteSlice.Put(rawData)

			sw.SetCurrentState(state.WorkerWriting)
			chOut <- msg
			sw.SetCurrentState(state.WorkerReading)
		}
	}
}
