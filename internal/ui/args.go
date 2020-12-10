package ui

import (
	"time"
)

type UIOpts struct {
	Refresh time.Duration `short:"r" long:"refresh" description:"UI refresh interval (0 to disable)" default:"100ms"`
}

func (opts *UIOpts) Validate() []string {
	return nil
}
