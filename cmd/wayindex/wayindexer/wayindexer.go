package wayindexer

import (
	"github.com/squizzling/osmindex/cmd/wayindex/wayindexer/state"
	"github.com/squizzling/osmindex/internal/tracker"
	"github.com/squizzling/osmindex/internal/wayindex"
)

type indexElement struct {
	blocks    []wayindex.WayIndexData
	locations []byte
}

type WayIndexer struct {
	tracker.WorkerTracker

	writer state.Writer
}

func (ni *WayIndexer) WriterState() interface{} {
	return &ni.writer
}
