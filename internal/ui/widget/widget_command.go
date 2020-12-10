package widget

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"

	"github.com/squizzling/osmindex/internal/imath"
	"github.com/squizzling/osmindex/internal/tracker"
	"github.com/squizzling/osmindex/internal/ui"
)

type widgetCommand struct {
	*views.BoxLayout
	writerText views.Widget

	wt *tracker.WidgetWorkers
}

func NewWidgetCmd(sp tracker.StateProvider, newWorkerWidget, newWriterWidget tracker.NewWidgetFunc) views.Widget {
	wni := &widgetCommand{
		BoxLayout:  views.NewBoxLayout(views.Vertical),
		writerText: newWriterWidget(sp.WriterState()),
		wt:         tracker.NewWidgetWorkers(sp, newWorkerWidget),
	}
	wni.BoxLayout.AddWidget(wni.wt.WorkerBox, 0)
	wni.BoxLayout.AddWidget(wni.writerText, 0)
	return wni
}

func (wc *widgetCommand) Size() (int, int) {
	workerWidth, workerHeight := wc.wt.WorkerBox.Size()
	writerWidth, writerHeight := wc.writerText.Size()
	return imath.MaxInt(workerWidth, writerWidth), workerHeight + writerHeight
}

func (wc *widgetCommand) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *ui.EventUpdate:
		wc.wt.RefreshStates()
	}
	wc.wt.WorkerBox.HandleEvent(ev)
	return wc.BoxLayout.HandleEvent(ev)
}
