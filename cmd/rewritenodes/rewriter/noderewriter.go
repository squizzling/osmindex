package rewriter

import (
	"github.com/squizzling/osmindex/cmd/rewritenodes/rewriter/state"
	"github.com/squizzling/osmindex/internal/tracker"
)

type NodeRewriter struct {
	tracker.WorkerTracker

	writer state.Writer
}

func (nr *NodeRewriter) WriterState() interface{} {
	return &nr.writer
}
