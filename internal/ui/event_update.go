package ui

import (
	"time"

	"github.com/gdamore/tcell"
)

type EventUpdate struct {
	t time.Time
}

func NewEventUpdate() tcell.Event {
	return &EventUpdate{
		t: time.Now(),
	}
}

func (eu *EventUpdate) When() time.Time {
	return eu.t
}
