package decompressor

import (
	"github.com/squizzling/osmindex/cmd/decompress/decompressor/state"
	"github.com/squizzling/osmindex/internal/tracker"
)

type Decompressor struct {
	tracker.WorkerTracker
	writer state.Writer
}

func (decomp *Decompressor) WriterState() interface{} {
	return &decomp.writer
}
