package ui

import (
	"github.com/gdamore/tcell/views"
)

type debugUI struct{}

func (d debugUI) ReplaceWidgets(widgets ...views.Widget) {}

func (d debugUI) Stop() {}
