package pprof

// Import this package as side-effect
// (import _ "github.com/meinside/monitor-go/pprof")
//
// for running pprof server locally.

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
