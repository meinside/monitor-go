package monitor

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// collect stats for current option
func (m *Monitor) stat() map[Option]string {
	stats := make(map[Option]string)

	if m.option&MonitorGoroutines > 0 {
		stats[MonitorGoroutines] = strconv.Itoa(collectNumGoroutines())
	}

	if m.option&MonitorMemory > 0 {
		stats[MonitorMemory] = strings.Join(collectMemoryUsage(), "\n")
	}

	// TODO - add more...

	return stats
}

// collect memory usage
func collectMemoryUsage() []string {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)

	return []string{
		fmt.Sprintf("Alloc: %s", numToBytes(stat.Alloc)),
		fmt.Sprintf("TotalAlloc: %s", numToBytes(stat.TotalAlloc)),
		fmt.Sprintf("NumGC: %s", numToHumanReadable(uint64(stat.NumGC), "")),
		fmt.Sprintf("Mallocs: %s", numToHumanReadable(stat.Mallocs, "")),
		fmt.Sprintf("Frees: %s", numToHumanReadable(stat.Frees, "")),
		fmt.Sprintf("NumLiveObjects: %s", numToHumanReadable(stat.Mallocs-stat.Frees, "")),

		// TODO - add more...
	}
}

// collect the number of goroutines
func collectNumGoroutines() int {
	return runtime.NumGoroutine()
}

// make given number of bytes human readable
func numToBytes(bytes uint64) string {
	return numToHumanReadable(bytes, "B")
}

// make given number human readable
func numToHumanReadable(num uint64, unit string) string {
	if num >= 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f T%s", float64(num)/(1024*1024*1024*1024), unit)
	} else if num >= 1024*1024*1024 {
		return fmt.Sprintf("%.2f G%s", float64(num)/(1024*1024*1024), unit)
	} else if num >= 1024*1024 {
		return fmt.Sprintf("%.2f M%s", float64(num)/(1024*1024), unit)
	} else if num >= 1024 {
		return fmt.Sprintf("%.2f K%s", float64(num)/(1024), unit)
	}

	return fmt.Sprintf("%d %s", num, unit)
}
