package monitor

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	// for pprof
	_ "net/http/pprof"
)

// Option is a type for setting monitoring options
type Option int

// Monitoring options
const (
	MonitorGoroutines Option = 1
	MonitorMemory     Option = 1 << 1

	DefaultMonitorInterval = 10 * time.Second

	DefaultMonitorOption Option = MonitorGoroutines | MonitorMemory
)

// Monitor struct
type Monitor struct {
	option   Option
	interval time.Duration
	httpPort int
	callback func(stats map[Option]string)

	httpServer *http.Server
	stopChan   chan interface{}
	verbose    bool
}

// Default returns a Monitor with default settings
func Default(callback func(stats map[Option]string)) *Monitor {
	return New(DefaultMonitorOption, DefaultMonitorInterval, 0, callback)
}

// New generates a Monitor with given settings
func New(option Option, interval time.Duration, port int, callback func(stats map[Option]string)) *Monitor {
	return &Monitor{
		interval: interval,
		option:   option,
		httpPort: port,
		callback: callback,
		stopChan: make(chan interface{}),
		verbose:  false,
	}
}

// SetOption sets the option
func (m *Monitor) SetOption(option Option) {
	m.option = option
}

// SetInterval sets the interval
func (m *Monitor) SetInterval(interval time.Duration) {
	m.interval = interval
}

// SetHTTPPort sets the HTTP port
func (m *Monitor) SetHTTPPort(port int) {
	m.httpPort = port
}

// SetVerbose sets verbose flag
func (m *Monitor) SetVerbose(verbose bool) {
	m.verbose = verbose
}

// Begin starts monitoring
func (m *Monitor) Begin() {
	timer := time.NewTicker(m.interval)

	go func() {
		// begin monitoring
		m.verboseLog("Start monitoring...")

		for {
			select {
			case <-timer.C:
				// time interval
				m.callback(m.stat())
			case <-m.stopChan:
				// stop monitoring
				m.verboseLog("Received from stop channel, stopping...")
				break
			}
		}
	}()

	if m.httpPort > 0 {
		go func() {
			addr := fmt.Sprintf(":%d", m.httpPort)
			m.httpServer = &http.Server{
				Addr:              addr,
				WriteTimeout:      10 * time.Second,
				ReadTimeout:       10 * time.Second,
				ReadHeaderTimeout: 10 * time.Second,
				IdleTimeout:       60 * time.Second,
			}

			// begin http server
			m.verboseLog(fmt.Sprintf("Start HTTP server... (http://HOST_NAME%s/debug/pprof)", addr))

			if err := m.httpServer.ListenAndServe(); err != nil {
				m.verboseLog(fmt.Sprintf("HTTP server stopping... (%s)", err))

				m.httpServer = nil
			}
		}()
	}
}

// Stop finishes monitoring
func (m *Monitor) Stop() {
	// stop monitoring
	m.verboseLog("Stop monitoring...")
	m.stopChan <- struct{}{}

	if m.httpPort > 0 && m.httpServer != nil {
		// stop http server
		m.verboseLog("Stop HTTP server...")
		m.httpServer.Shutdown(context.Background())
	}
}

// CurrentStat fetches current stat
func (m *Monitor) CurrentStat() map[Option]string {
	return m.stat()
}

// print verbose log message
func (m *Monitor) verboseLog(str string) {
	if m.verbose {
		log.Printf("[monitor] %s", str)
	}
}
