package widget

import (
	"fmt"

	"github.com/gdamore/tcell/views"

	"github.com/squizzling/osmindex/cmd/tagfilter/tagfilterer/state"
	"github.com/squizzling/osmindex/internal/ui/widget"
)

type widgetStateWorker struct {
	*widget.WidgetCommonLine

	worker *state.Worker
}

func NewWidgetStateWorker(worker interface{}) views.Widget {
	wsw := &widgetStateWorker{
		worker: worker.(*state.Worker),
	}
	wsw.WidgetCommonLine = widget.NewWidgetCommonLine(wsw.update)
	return wsw
}

func (wsw *widgetStateWorker) update() {
	switch wsw.worker.CurrentState() {
	case state.WorkerReading:
		wsw.SetText("Worker: waiting for block")
	case state.WorkerDecoding:
		wsw.SetText(fmt.Sprintf("Worker: decoding block %d, %d kept, %d dropped", wsw.worker.CurrentBlock(), wsw.worker.Kept(), wsw.worker.Dropped()))
	case state.WorkerWorking:
		wsw.SetText(fmt.Sprintf("Worker: processing block %d, %d kept, %d dropped", wsw.worker.CurrentBlock(), wsw.worker.Kept(), wsw.worker.Dropped()))
	case state.WorkerEncoding:
		wsw.SetText(fmt.Sprintf("Worker: encoding block %d, %d kept, %d dropped", wsw.worker.CurrentBlock(), wsw.worker.Kept(), wsw.worker.Dropped()))
	case state.WorkerWriting:
		wsw.SetText(fmt.Sprintf("Worker: sending block %d, %d kept, %d dropped", wsw.worker.CurrentBlock(), wsw.worker.Kept(), wsw.worker.Dropped()))
	}
}
