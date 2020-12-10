package widget

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/views"

	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/ui"
)

type PassProgress struct {
	*WidgetCommonLine

	pfs *pbf.ProcessFileState

	pass  uint64
	count uint64
}

func NewPassProgress(pass, count uint64, pfs *pbf.ProcessFileState) views.Widget {
	pp := &PassProgress{
		pass:  pass,
		count: count,
		pfs:   pfs,
	}
	pp.WidgetCommonLine = NewWidgetCommonLine(pp.update)
	return pp
}

func (wpp *PassProgress) update() {
	msg := fmt.Sprintf("Pass %d/%d", wpp.pass+1, wpp.count)

	pfsPct := float64(wpp.pfs.At()) / float64(wpp.pfs.Max())
	if pfsPct < 0 {
		pfsPct = 0
	} else if pfsPct > 1 {
		pfsPct = 1
	}
	pct := float64(wpp.pass) / float64(wpp.count)
	if !math.IsNaN(pfsPct) {
		pct += pfsPct / float64(wpp.count)
	}

	if wpp.View != nil {
		w, _ := wpp.View.Size()
		wpp.SetMarkup(ui.MarkupAsProgress(msg, w, pct))
	} else {
		wpp.SetText(msg)
	}
}
