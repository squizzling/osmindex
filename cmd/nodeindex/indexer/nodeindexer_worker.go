package indexer

import (
	"context"
	"math"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/cmd/nodeindex/indexer/state"
	"github.com/squizzling/osmindex/internal/morton"
	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/pool"
	"github.com/squizzling/osmindex/internal/t"
)

const (
	notAssigned = uint64(math.MaxUint64)
)

func (ni *NodeIndexer) Worker(ctx context.Context, chIn <-chan pbf.Block, chOut chan<- pbf.Block) {
	pbr := &t.PBReader{
		SkipStringTable: true,
		SkipWays:        true,
		SkipRelations:   true,
	}

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

		var pb t.PrimitiveBlock
		pbr.ReadPrimitiveBlock(rawData.Buffer, &pb)

		sw.SetCurrentState(state.WorkerProcessing)

		var indexBlocks []indexRange
		var outputLoc []int64

		firstId := notAssigned
		lastId := notAssigned

		for _, pg := range pb.PrimitiveGroup {
			if pg.Dense == nil {
				continue
			}
			for idx, id := range pg.Dense.Id {
				uid := uint64(id)
				if uid == lastId+1 { // extend current block
					lastId = uid
				} else if uid == lastId+2 {
					outputLoc = append(outputLoc, 0)
					lastId = uid
				} else {
					if firstId != notAssigned {
						indexBlocks = append(indexBlocks, indexRange{
							first: firstId,
							last:  lastId,
						})
					}

					firstId = uid
					lastId = uid
				}

				lat := int32((pb.LatOffset + (int64(pb.Granularity) * pg.Dense.Lat[idx])) / 100)
				lon := int32((pb.LonOffset + (int64(pb.Granularity) * pg.Dense.Lon[idx])) / 100)
				outputLoc = append(outputLoc, int64(morton.Encode(lon, lat|math.MinInt32)))
			}
		}

		if firstId != notAssigned {
			indexBlocks = append(indexBlocks, indexRange{
				first: firstId,
				last:  lastId,
			})
		}

		msg := pbf.Block{
			Index: data.Index,
			Data: elementBatch{
				ranges:    indexBlocks,
				locations: outputLoc,
			},
			BlobType: "IndexData",
		}

		if len(indexBlocks) == 0 {
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
