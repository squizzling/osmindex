package ui

import (
	"context"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type normalUI struct {
	screen tcell.Screen
	cancel context.CancelFunc
	app    *views.Application

	appReady chan struct{}
	view     views.View
	*views.BoxLayout

	shutdown chan struct{}
	wg       sync.WaitGroup
}

func NewUI(updateInterval time.Duration, cancel context.CancelFunc) UI {
	if updateInterval == 0 {
		return &debugUI{}
	}
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	app := &views.Application{}
	u := &normalUI{
		app:       app,
		screen:    screen,
		cancel:    cancel,
		BoxLayout: views.NewBoxLayout(views.Vertical),
		appReady:  make(chan struct{}),
		shutdown:  make(chan struct{}),
	}

	app.SetScreen(screen)
	app.SetRootWidget(u)
	app.Start()

	// app.Start() will call SetView() on the root widget
	// as the last action before entering the event loop.
	<-u.appReady

	u.wg.Add(1)
	go func() {
		ticker := time.NewTicker(updateInterval)
		for {
			select {
			case <-ticker.C:
				screen.PostEventWait(NewEventUpdate())
			case <-u.shutdown:
				ticker.Stop()
				u.wg.Done()
				return
			}
		}
	}()

	return u
}

func (u *normalUI) SetView(v views.View) {
	if u.view == nil {
		close(u.appReady)
	}
	u.view = v
	u.BoxLayout.SetView(v)
}

func (u *normalUI) HandleEvent(e tcell.Event) bool {
	switch e := e.(type) {
	case *tcell.EventKey:
		if e.Modifiers() != tcell.ModCtrl {
			break
		}
		switch e.Key() {
		case tcell.KeyCtrlA:
			panic("abort")
		case tcell.KeyCtrlC:
			u.cancel()
			return true
		case tcell.KeyCtrlD:
			_ = pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
			return true
		case tcell.KeyCtrlL:
			u.app.Refresh()
			return true
		}
	}
	return u.BoxLayout.HandleEvent(e)

}

func (u *normalUI) ReplaceWidgets(widgets ...views.Widget) {
	u.app.PostFunc(func() {
		for _, w := range u.Widgets() {
			u.RemoveWidget(w)
		}
		for i, w := range widgets {
			if i == len(widgets)-1 {
				u.AddWidget(w, 1.0)
			} else {
				u.AddWidget(w, 0)
			}
		}
	})
}

func (u *normalUI) Stop() {
	close(u.shutdown)
	u.wg.Wait()
	u.app.Quit()
	if err := u.app.Wait(); err != nil {
		panic(err)
	}
}
