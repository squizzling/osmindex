package tracker

// TODO: This package is a bit of a mess

import (
	"github.com/gdamore/tcell/views"
)

type StateProvider interface {
	Version() uint64
	WriterState() interface{}
	WorkerStates() []interface{}
}

// WidgetWorkers is a Widget which holds and maintains a list of Widgets for each worker exposed by the StateProvider
type WidgetWorkers struct {
	provider    StateProvider
	new         func(interface{}) views.Widget
	states      map[interface{}]views.Widget
	lastVersion uint64

	WorkerBox *views.BoxLayout
}

type NewWidgetFunc func(state interface{}) views.Widget

func NewWidgetWorkers(sp StateProvider, new NewWidgetFunc) *WidgetWorkers {
	return &WidgetWorkers{
		provider:  sp,
		new:       new,
		WorkerBox: views.NewBoxLayout(views.Vertical),
	}
}

func (wt *WidgetWorkers) RefreshStates() {
	if currentVersion := wt.provider.Version(); wt.lastVersion == currentVersion {
		return
	} else {
		wt.lastVersion = currentVersion
	}

	oldWidgets := wt.states
	for _, oldWidget := range oldWidgets {
		wt.WorkerBox.RemoveWidget(oldWidget)
	}

	newWidgets := make(map[interface{}]views.Widget)
	for _, workerState := range wt.provider.WorkerStates() {
		var workerWidget views.Widget
		var ok bool
		if workerWidget, ok = oldWidgets[workerState]; !ok {
			workerWidget = wt.new(workerState)
		}
		newWidgets[workerState] = workerWidget
		wt.WorkerBox.AddWidget(workerWidget, 0)
	}
	wt.states = newWidgets
}
