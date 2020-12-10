package widget

import (
	"fmt"

	"github.com/gdamore/tcell/views"

	"github.com/squizzling/osmindex/cmd/nodeindex/indexer/state"
	"github.com/squizzling/osmindex/internal/ui"
	"github.com/squizzling/osmindex/internal/ui/widget"
)

type widgetStateWriter struct {
	*widget.WidgetCommonLine

	writer *state.Writer
}

func NewWidgetStateWriter(writer interface{}) views.Widget {
	wsw := &widgetStateWriter{
		writer: writer.(*state.Writer),
	}
	wsw.WidgetCommonLine = widget.NewWidgetCommonLine(wsw.update)
	return wsw
}

func (wsw *widgetStateWriter) update() {
	switch wsw.writer.CurrentState() {
	case state.WriterReading:
		wsw.SetText("Writer: waiting for block")
	case state.WriterWritingIndex:
		wsw.SetText(fmt.Sprintf("Writer: writing index block %d", wsw.writer.CurrentBlock()))
	case state.WriterWritingLocations:
		wsw.SetText("Writer: writing locations")
	case state.WriterCopyingLocations:
		at := wsw.writer.CopyAmount()
		max := wsw.writer.CopyTotal()
		width, _ := wsw.View.Size()
		pct := float64(at) / float64(max)
		msg := fmt.Sprintf("Writer: copying locations (%d/%d) %2.01f%%", at, max, 100*pct)
		wsw.SetMarkup(ui.MarkupAsProgress(msg, width, pct))
	case state.WriteFinalizing:
		wsw.SetText("Writer: finalizing index")
	}
}
