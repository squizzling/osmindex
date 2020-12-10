package widget

import (
	"fmt"

	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/ui"
)

type ProcessFileStateWidget struct {
	*WidgetCommonLine
	action string

	ps *pbf.ProcessFileState
}

func NewProcessFileStateWidget(action string, ps *pbf.ProcessFileState) *ProcessFileStateWidget {
	pfsw := &ProcessFileStateWidget{
		action: action,
		ps:     ps,
	}
	pfsw.WidgetCommonLine = NewWidgetCommonLine(pfsw.update)
	return pfsw
}

func (pfsw *ProcessFileStateWidget) update() {
	prefix := fmt.Sprintf("[%s] Processing PBF %s: ", pfsw.action, pfsw.ps.Filename())
	switch s := pfsw.ps.CurrentState(); s {
	case pbf.ProcessFileStateInitializing:
		pfsw.SetText(prefix + "initializing")
	case pbf.ProcessFileStateLoadingBlock, pbf.ProcessFileStateSendingBlock:
		action := "loading block"
		if s == pbf.ProcessFileStateSendingBlock {
			action = "sending block"
		}
		at := pfsw.ps.At()
		max := pfsw.ps.Max()
		pct := float64(at) / float64(max)

		msg := fmt.Sprintf("%s %d/%d %2.01f%%", action, at, max, 100*pct)
		if pfsw.View != nil {
			w, _ := pfsw.View.Size()
			pfsw.SetMarkup(ui.MarkupAsProgress(prefix+msg, w, pct))
		} else {
			pfsw.SetText(prefix + msg)
		}
	case pbf.ProcessFileStateShuttingDown:
		pfsw.SetText(prefix + "shutting down")
	}
}
