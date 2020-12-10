package widget

import (
	"fmt"

	"github.com/gdamore/tcell/views"

	"github.com/squizzling/osmindex/cmd/decompress/decompressor/state"
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
	case state.WriterReceivingBlock:
		wsw.SetText("Writer: waiting for block")
	case state.WriterWritingData:
		wsw.SetText(fmt.Sprintf("Writer: writing block %d", wsw.writer.CurrentBlock()))
	}

}
