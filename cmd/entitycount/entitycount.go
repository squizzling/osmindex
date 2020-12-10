package main

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/pool"
	"github.com/squizzling/osmindex/internal/system"
	"github.com/squizzling/osmindex/internal/t"
)

func atomicSetMaximum(new uint64, max *uint64) {
	for old := atomic.LoadUint64(max); new > old && !atomic.CompareAndSwapUint64(max, old, new); old = atomic.LoadUint64(max) {
		/* do nothing */
	}
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func countEntities(blockCount, nodeCount, wayCount, relationCount, wayRefCount, maxNode, maxWay, maxRelation, maxWayRefCount *uint64) pbf.WorkFunc {
	return func(ctx context.Context, chIn <-chan pbf.Block, chOut chan<- pbf.Block) {
		localBlockCount := uint64(0)
		localNodeCount := uint64(0)
		localWayCount := uint64(0)
		localRelationCount := uint64(0)
		localWayRefCount := uint64(0)
		localMaxNode := uint64(0)
		localMaxWay := uint64(0)
		localMaxRelation := uint64(0)
		localMaxWayRefCount := uint64(0)

		pbr := &t.PBReader{
			SkipStringTable: true,
		}

		for data := range chIn {
			localBlockCount++
			if data.BlobType != "OSMData" {
				continue
			}

			buffer := data.Data.(mmap.MMap)

			var b t.Blob
			pbr.ReadBlob(buffer, &b)
			rawData := b.GetRawData()

			var pb t.PrimitiveBlock
			pbr.ReadPrimitiveBlock(rawData.Buffer, &pb)

			for _, pg := range pb.PrimitiveGroup {
				if pg.Dense != nil {
					localNodeCount += uint64(len(pg.Dense.Id))
					for _, nid := range pg.Dense.Id {
						localMaxNode = max(uint64(nid), localMaxNode)
					}
				}
				localWayCount += uint64(len(pg.Ways))
				localRelationCount += uint64(len(pg.Relations))
				for _, w := range pg.Ways {
					localWayRefCount += uint64(len(w.Refs))
					localMaxWay = max(uint64(w.ID), localMaxWay)
					localMaxWayRefCount = max(localMaxWayRefCount, uint64(len(w.Refs)))
				}
				for _, r := range pg.Relations {
					localMaxRelation = max(uint64(r.ID), localMaxRelation)
				}
			}
			pool.ByteSlice.Put(rawData)
		}

		atomic.AddUint64(blockCount, localBlockCount)
		atomic.AddUint64(nodeCount, localNodeCount)
		atomic.AddUint64(wayCount, localWayCount)
		atomic.AddUint64(relationCount, localRelationCount)
		atomic.AddUint64(wayRefCount, localWayRefCount)
		atomicSetMaximum(localMaxNode, maxNode)
		atomicSetMaximum(localMaxWay, maxWay)
		atomicSetMaximum(localMaxRelation, maxRelation)
		atomicSetMaximum(localMaxWayRefCount, maxWayRefCount)
	}
}

func main() {
	opts := parseArgs(os.Args[1:])

	teardown := system.Setup(&opts.SystemOpts)
	defer teardown()

	var blockCount, nodeCount, wayCount, relationCount, wayRefCount uint64
	var maxNode, maxWay, maxRelation uint64
	var maxWayRefCount uint64
	pbf.ProcessFileAsync(
		context.Background(),
		opts.InputFile,
		countEntities(&blockCount, &nodeCount, &wayCount, &relationCount, &wayRefCount, &maxNode, &maxWay, &maxRelation, &maxWayRefCount),
		nil,
		opts.Workers,
	).Wait()

	fmt.Printf("Counted %d blocks\n", blockCount)
	fmt.Printf("Counted %d nodes, %d max\n", nodeCount, maxNode)
	fmt.Printf("Counted %d ways, %d max\n", wayCount, maxWay)
	fmt.Printf("Counted %d ways refs, %d maximum length\n", wayRefCount, maxWayRefCount)
	fmt.Printf("Counted %d relations, %d max\n", relationCount, maxRelation)
}
