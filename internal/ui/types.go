package ui

import (
	"github.com/gdamore/tcell/views"
)

type UI interface {
	ReplaceWidgets(widgets ...views.Widget)
	Stop()
}
