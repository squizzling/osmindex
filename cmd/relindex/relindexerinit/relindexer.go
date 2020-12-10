package relindexerinit

import (
	"github.com/squizzling/osmindex/cmd/relindex/relindexerinit/state"
	"github.com/squizzling/osmindex/internal/tracker"
)

type RelIndexerInitializer struct {
	tracker.WorkerTracker

	writer state.Writer
}

func (ri *RelIndexerInitializer) WriterState() interface{} {
	return &ri.writer
}
