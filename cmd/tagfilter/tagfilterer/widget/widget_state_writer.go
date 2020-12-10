package widget

import (
	"fmt"

	"github.com/gdamore/tcell/views"

	"github.com/squizzling/osmindex/cmd/tagfilter/tagfilterer/state"
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
		wsw.SetText(fmt.Sprintf("Writer: waiting for block, dropped %d/%d blocks", wsw.writer.Filler(), wsw.writer.CurrentBlock()))
	case state.WriterWriting:
		wsw.SetText(fmt.Sprintf("Writer: writing block, dropped %d/%d blocks", wsw.writer.Filler(), wsw.writer.CurrentBlock()))
	}
}
