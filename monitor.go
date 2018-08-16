package monitor

import (
	"log"
	"time"
)

const (
	// DefaultMonitorInterval is the default interval seconds for monitoring
	DefaultMonitorInterval = 10 * time.Second
)

// Option is a type for setting monitoring options
type Option int

// Monitoring options
const (
	DefaultMonitorOption Option = MonitorGoroutines | MonitorMemory

	MonitorGoroutines Option = 1
	MonitorMemory     Option = 1 << 1
)

// Monitor struct
type Monitor struct {
	interval time.Duration
	option   Option
	callback func(stats map[Option]string)

	stopChan chan interface{}
	verbose  bool
}

// Default returns a Monitor with default settings
func Default(callback func(stats map[Option]string)) *Monitor {
	return New(DefaultMonitorOption, DefaultMonitorInterval, callback)
}

// New generates a Monitor with given settings
func New(option Option, interval time.Duration, callback func(stats map[Option]string)) *Monitor {
	return &Monitor{
		interval: interval,
		option:   option,
		callback: callback,
		stopChan: make(chan interface{}),
		verbose:  false,
	}
}

// SetInterval sets the interval
func (m *Monitor) SetInterval(interval time.Duration) {
	m.interval = interval
}

// SetOption sets the option
func (m *Monitor) SetOption(option Option) {
	m.option = option
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
}

// Stop finishes monitoring
func (m *Monitor) Stop() {
	m.stopChan <- struct{}{}
}

// print verbose log message
func (m *Monitor) verboseLog(str string) {
	if m.verbose {
		log.Printf("[monitor] %s", str)
	}
}
