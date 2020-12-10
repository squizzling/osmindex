package widget

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"

	"github.com/squizzling/osmindex/internal/ui"
)

type WidgetCommonLine struct {
	*views.SimpleStyledText
	update func()
	View   views.View
}

func NewWidgetCommonLine(update func()) *WidgetCommonLine {
	return &WidgetCommonLine{
		SimpleStyledText: views.NewSimpleStyledText(),
		update:           update,
	}
}

func (wsw *WidgetCommonLine) HandleEvent(ev tcell.Event) bool {
	if _, ok := ev.(*ui.EventUpdate); ok {
		wsw.update()
	}
	return wsw.SimpleStyledText.HandleEvent(ev)
}

func (wsw *WidgetCommonLine) SetView(view views.View) {
	wsw.View = view
	wsw.SimpleStyledText.SetView(view)
}
