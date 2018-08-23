# monitor-go

A go library for monitoring application's stats (periodically).

## Supported stats and functionalities

- [X] Memory usage
- [X] Number of goroutines
- [X] pprof monitoring through http server
- [ ] ?

## How to get

```bash
$ go get -u github.com/meinside/monitor-go
```

## Usages

### Example

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

### Easy pprof monitoring

Run your code with `github.com/meinside/monitor-go/pprof` package imported as side-effect:

```go
package main

import (
	_ "github.com/meinside/monitor-go/pprof"
)

func main() {
	// ... do your business here
}
```

then visit `http://HOST_NAME:61000/debug/pprof/` for pprof stats.

## license

MIT

