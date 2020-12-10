package relindexerinit

import (
	"context"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/cmd/relindex/relindexerinit/state"
	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/pool"
	"github.com/squizzling/osmindex/internal/t"
)

func (ri *RelIndexerInitializer) Worker() func(ctx context.Context, chIn <-chan pbf.Block, chOut chan<- pbf.Block) {
	return func(ctx context.Context, chIn <-chan pbf.Block, chOut chan<- pbf.Block) {
		pbr := &t.PBReader{
			SkipStringTable: false,
			SkipDenseNodes:  true,
			SkipWays:        true,
			SkipRelations:   false,
		}

		sw := &state.Worker{}
		ri.TrackWorker(sw)
		defer ri.UntrackWorker(sw)

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

			var pb t.PrimitiveBlock
			pbr.ReadPrimitiveBlock(rawData.Buffer, &pb)

			sw.SetCurrentState(state.WorkerProcessing)

			var outputData []int64

			for _, pg := range pb.PrimitiveGroup {
				for _, r := range pg.Relations {
					outputData = append(outputData, int64(r.ID))
					wayCount := int64(0)
					for _, mt := range r.MemberType {
						if mt == t.RelationMemberTypeWay {
							wayCount++
						}
					}
					outputData = append(outputData, wayCount)
					for idx, wayId := range r.MemIDs {
						if r.MemberType[idx] != t.RelationMemberTypeWay {
							continue
						}

						outputData = append(outputData, -wayId)
					}
				}
			}

			msg := pbf.Block{
				Index:    data.Index,
				Data:     outputData,
				BlobType: "RelationData",
			}

			pool.ByteSlice.Put(rawData)

			sw.SetCurrentState(state.WorkerWriting)
			chOut <- msg
			sw.SetCurrentState(state.WorkerReading)
		}
	}
}
