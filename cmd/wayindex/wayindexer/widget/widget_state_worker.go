package widget

import (
	"fmt"

	"github.com/gdamore/tcell/views"

	"github.com/squizzling/osmindex/cmd/wayindex/wayindexer/state"
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
		wsw.SetText(fmt.Sprintf("Worker: decoding block %d", wsw.worker.CurrentBlock()))
	case state.WorkerProcessing:
		wsw.SetText(fmt.Sprintf("Worker: processing block %d", wsw.worker.CurrentBlock()))
	case state.WorkerWriting:
		wsw.SetText(fmt.Sprintf("Worker: sending block %d", wsw.worker.CurrentBlock()))
	}
}
