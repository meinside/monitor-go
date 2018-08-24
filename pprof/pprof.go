package pprof

// Import this package for running pprof server locally.

import (
	"time"

	"github.com/meinside/monitor-go"
)

const (
	pprofHTTPPort = 61000
)

func init() {
	mon := monitor.New(monitor.MonitorNone, 1*time.Hour, pprofHTTPPort, nil)
	mon.SetVerbose(true)

	mon.Begin()
}
