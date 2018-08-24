# monitor-go/pprof

## Usage

### Import: use import side effects

```go
import _ "github.com/meinside/monitor-go/pprof"
```

### How to see the reports

Open `http://some-address:61000/debug/pprof` in a web browser.

#### CPU Profiling

Do CPU profiling with:

```bash
$ go tool pprof http://some-address:61000/debug/pprof/profile
```

#### Memory Profiling

Do memory profiling with:

```bash
$ go tool pprof http://some-address:61000/debug/pprof/heap
```

#### Trace

Save trace data with:

```bash
$ curl http://some-address:61000/debug/pprof/trace?seconds=30 > trace.out
```

and start a tracing interface with it:

```bash
$ go tool trace -http=':8080' trace.out
```

