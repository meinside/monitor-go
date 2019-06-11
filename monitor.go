package monitor

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

// Option is a type for setting monitoring options
type Option int

// constants
const (
	// Monitoring options
	MonitorNone       Option = 0
	MonitorGoroutines Option = 1
	MonitorMemory     Option = 1 << 1

	// default monitoring option
	DefaultMonitorOption Option = MonitorGoroutines | MonitorMemory

	// default monitoring interval seconds
	DefaultMonitorInterval = 10 * time.Second
)

var _stdout = log.New(os.Stdout, "", log.LstdFlags)

// return string value of given Option
func (o Option) String() string {
	switch o {
	case MonitorNone:
		return "None"
	case MonitorGoroutines:
		return "GoRoutines"
	case MonitorMemory:
		return "Memory"
	default:
		return "UnknownOption"
	}
}

// Monitor struct
type Monitor struct {
	// values for monitoring
	option   Option
	interval time.Duration
	callback func(stats map[Option]string)

	// for pprof
	httpPort   int
	httpServer *http.Server

	// for stopping
	stopChan chan interface{}

	// verbose flag
	verbose bool
}

// Default returns a Monitor with default settings.
func Default(callback func(stats map[Option]string)) *Monitor {
	return New(DefaultMonitorOption, DefaultMonitorInterval, 0, callback)
}

// New generates a Monitor with given settings.
//
// `stats` will be empty if `option` is `MonitorNone`.
func New(option Option, interval time.Duration, port int, callback func(stats map[Option]string)) *Monitor {
	return &Monitor{
		option:     option,
		interval:   interval,
		callback:   callback,
		httpPort:   port,
		httpServer: nil,
		stopChan:   make(chan interface{}),
		verbose:    false,
	}
}

// SetOption sets the option.
func (m *Monitor) SetOption(option Option) {
	m.option = option
}

// SetInterval sets the interval.
func (m *Monitor) SetInterval(interval time.Duration) {
	m.interval = interval
}

// SetHTTPPort sets the HTTP port.
func (m *Monitor) SetHTTPPort(port int) {
	m.httpPort = port
}

// SetVerbose sets verbose flag.
func (m *Monitor) SetVerbose(verbose bool) {
	m.verbose = verbose
}

// Begin starts monitoring.
func (m *Monitor) Begin() {
	timer := time.NewTicker(m.interval)

	if m.callback != nil {
		go func() {
			// begin monitoring
			m.verboseLog("start monitoring...")

			for {
				select {
				case <-timer.C:
					// time interval
					m.callback(m.stat())
				case <-m.stopChan:
					// stop monitoring
					m.verboseLog("received from stop channel, stopping...")
					break
				}
			}
		}()
	}

	if m.httpPort > 0 {
		go func() {
			addr := fmt.Sprintf(":%d", m.httpPort)

			// additional pprof handlers
			mux := http.NewServeMux()
			mux.HandleFunc("/debug/pprof/", pprof.Index)
			mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

			m.httpServer = &http.Server{
				Addr:              addr,
				Handler:           mux,
				WriteTimeout:      120 * time.Second,
				ReadTimeout:       10 * time.Second,
				ReadHeaderTimeout: 10 * time.Second,
				IdleTimeout:       300 * time.Second,
			}

			// begin http server
			m.verboseLog("start pprof http server on http://localhost%s/debug/pprof", addr)

			if err := m.httpServer.ListenAndServe(); err != nil {
				m.verboseLog("pprof http server stopped: %s", err)

				m.httpServer = nil
			}
		}()
	}
}

// Stop finishes monitoring.
func (m *Monitor) Stop() {
	// stop monitoring
	if m.callback != nil {
		m.verboseLog("stopping monitoring...")

		m.stopChan <- struct{}{}
	}

	// stop http server
	if m.httpPort > 0 && m.httpServer != nil {
		m.verboseLog("stopping pprof http server...")

		m.httpServer.Shutdown(context.Background())
	}
}

// CurrentStat fetches current stat.
func (m *Monitor) CurrentStat() map[Option]string {
	return m.stat()
}

// print verbose log message
func (m *Monitor) verboseLog(format string, args ...interface{}) {
	if m.verbose {
		_stdout.Printf(format, args...)
	}
}
