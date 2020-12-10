package indexer

import (
	"github.com/squizzling/osmindex/cmd/nodeindex/indexer/state"
	"github.com/squizzling/osmindex/internal/tracker"
)

type indexRange struct {
	first uint64
	last  uint64
}

type elementBatch struct {
	ranges    []indexRange
	locations []int64
}

type NodeIndexer struct {
	tracker.WorkerTracker
	writer state.Writer
}

func (ni *NodeIndexer) WriterState() interface{} {
	return &ni.writer
}
