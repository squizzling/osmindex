package widget

import (
	"fmt"

	"github.com/gdamore/tcell/views"

	"github.com/squizzling/osmindex/cmd/rewritenodes/rewriter/state"
	"github.com/squizzling/osmindex/internal/ui"
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
	case state.WorkerWorking:
		at := wsw.worker.CurrentWay()
		max := wsw.worker.WayCount()
		msg := fmt.Sprintf("Worker: working on block %d, way %d/%d",
			wsw.worker.CurrentBlock(),
			at,
			max,
		)
		if wsw.View == nil {
			wsw.SetText(msg)
		} else {
			w, _ := wsw.View.Size()
			pct := float64(at) / float64(max)
			wsw.SetMarkup(ui.MarkupAsProgress(msg, w, pct))
		}
	case state.WorkerWriting:
		wsw.SetText(fmt.Sprintf("Worker: sending block %d", wsw.worker.CurrentBlock()))
	}
}
