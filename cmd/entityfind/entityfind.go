package main

import (
	"context"
	"fmt"
	"os"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/pool"
	"github.com/squizzling/osmindex/internal/system"
	"github.com/squizzling/osmindex/internal/t"
)

func isIn(ids []uint64, id uint64) bool {
	for _, id2 := range ids {
		if id2 == id {
			return true
		}
	}
	return false
}

func findEntity(ids []uint64) pbf.WorkFunc {
	return func(ctx context.Context, chIn <-chan pbf.Block, chOut chan<- pbf.Block) {
		pbr := &t.PBReader{}

		for data := range chIn {
			if data.BlobType != "OSMData" {
				continue
			}

			buffer := data.Data.(mmap.MMap)

			var b t.Blob
			pbr.ReadBlob(buffer, &b)
			rawData := b.GetRawData()

			var pb t.PrimitiveBlock
			pbr.ReadPrimitiveBlock(rawData.Buffer, &pb)

			ns := pb.ToNodes()
			for _, n := range ns {
				if isIn(ids, uint64(n.Id)) {
					fmt.Printf("%#v\n", n)
				}
			}

			ws := pb.ToWays()
			for _, w := range ws {
				if isIn(ids, uint64(w.Id)) {
					fmt.Printf("%#v\n", w)
				}
			}

			rs := pb.ToRelations()
			for _, r := range rs {
				if isIn(ids, uint64(r.Id)) {
					fmt.Printf("%#v\n", r)
				}
			}

			pool.ByteSlice.Put(rawData)
		}
	}
}

func main() {
	opts := parseArgs(os.Args[1:])

	teardown := system.Setup(&opts.SystemOpts)
	defer teardown()

	pbf.ProcessFileAsync(
		context.Background(),
		opts.InputFile,
		findEntity(opts.EntityID),
		nil,
		opts.Workers,
	).Wait()
}
