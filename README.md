# monitor-go

A go library for monitoring application's stats periodically.

## supported stats

- [X] Memory usage
- [X] Number of goroutines
- [ ] ?

## how to get

```bash
$ go get -u github.com/meinside/monitor-go
```

## usage

### example

```go
package main

import (
	"fmt"
	"time"

	"github.com/meinside/monitor-go"
)

func logStat(stats map[monitor.Option]string) {
	fmt.Println("---- stat ----")

	for _, v := range stats {
		fmt.Printf("%s\n", v)
	}
}

func main() {
	mon := monitor.Default(logStat)
	mon.SetInterval(3 * time.Second)
	mon.SetOption(monitor.MonitorMemory)
	mon.SetVerbose(true)

	fmt.Println("> starting monitoring...")

	mon.Begin()

	// do something which takes some time
	finished := make(chan struct{})
	go func(ch chan struct{}) {
		fmt.Println("> starting task...")

		time.Sleep(10 * time.Second)

		ch <- struct{}{}

		fmt.Println("> job finished...")
	}(finished)

	// wait for it
	<-finished

	mon.Stop()

	fmt.Println("> monitor finished.")
}
```

## license

MIT

