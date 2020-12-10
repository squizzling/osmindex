package tagfilterer

import (
	"github.com/squizzling/osmindex/cmd/tagfilter/tagfilterer/state"
	"github.com/squizzling/osmindex/internal/tracker"
)

type TagFilterer struct {
	tracker.WorkerTracker

	writer state.Writer
}

func (tf *TagFilterer) WriterState() interface{} {
	return &tf.writer
}
