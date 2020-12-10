package rewriter

import (
	"context"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/cmd/rewritenodes/rewriter/state"
	"github.com/squizzling/osmindex/internal/nodeindex"
	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/pool"
	"github.com/squizzling/osmindex/internal/t"
)

func (nr *NodeRewriter) Worker(ix *nodeindex.NodeSubIndex, finalPass bool) pbf.WorkFunc {
	return func(ctx context.Context, chIn <-chan pbf.Block, chOut chan<- pbf.Block) {
		pbr := t.PBReader{
			SkipStringTable: true,
			SkipDenseNodes:  true,
			SkipRelations:   true,
		}

		sw := &state.Worker{}
		nr.TrackWorker(sw)
		defer nr.UntrackWorker(sw)

		for data := range chIn {
			if data.BlobType != "OSMData" {
				chOut <- data
				continue
			}
			sw.SetStateDecoding(data.Index)

			buffer := data.Data.(mmap.MMap)

			var b t.Blob
			pbr.ReadBlob(buffer, &b)
			rawData := b.GetRawData()

			var pb t.PrimitiveBlock
			pbr.ReadPrimitiveBlock(rawData.Buffer, &pb)

			totalWays := 0
			for _, pg := range pb.PrimitiveGroup {
				totalWays += len(pg.Ways)
			}

			if totalWays == 0 {
				sw.SetCurrentState(state.WorkerWriting)
				chOut <- data
				sw.SetCurrentState(state.WorkerReading)
				pool.ByteSlice.Put(rawData)
				continue
			}

			sw.SetStateWorking(totalWays)
			hint := 0
			for _, pg := range pb.PrimitiveGroup {
				for _, w := range pg.Ways {
					sw.NextWay()
					outIndex := 0
					var found bool
					var mortonLocation int64
					for _, ref := range w.Refs {
						if ref < 0 { // Pass through IDs which are already converted to a Morton location (see wayindexer_worker.go)
							w.Refs[outIndex] = ref
							outIndex++
						} else if mortonLocation, hint, found = ix.FindNode(uint64(ref), hint); found { // Found ref
							w.Refs[outIndex] = mortonLocation
							outIndex++
						} else if !finalPass { // On non-final passes keep un-found refs as-is
							w.Refs[outIndex] = ref
							outIndex++
						} else { // On the final pass, delete un-found refs by not outputting anything
							// TODO: Maybe we keep it.  There's some strange behavior in NSW.
							w.Refs[outIndex] = ref
							outIndex++
						}
					}
					w.Refs = w.Refs[:outIndex]
				}
			}

			sw.SetCurrentState(state.WorkerWriting)
			chOut <- pbf.Block{
				Index:    data.Index,
				Data:     (&t.Blob{RawFunc: pb.Write}).Write(make([]byte, 0, 2*1048576)),
				BlobType: "OSMData",
			}
			sw.SetCurrentState(state.WorkerReading)
			pool.ByteSlice.Put(rawData)
		}
	}
}
