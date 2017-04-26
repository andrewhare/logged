# logged

A simple Go logging package inspired by [Dave Cheney's discussion on logging](https://dave.cheney.net/2015/11/05/lets-talk-about-logging).

GoDoc: https://godoc.org/github.com/andrewhare/logged  
Go Report Card: https://goreportcard.com/report/github.com/andrewhare/logged  
CI: https://circleci.com/gh/andrewhare/logged

# Installation

```
$ go get github.com/andrewhare/logged
```

# Create a log

Logged is designed to be very simple and only expose the bare minimum functionality. 

### Serializer

The log can be configured with either a JSON or text serializer that write to an `io.Writer`:

```go
// JSON
log := logged.New(logged.NewJSONSerializer(os.Stdout))

// Text
log := logged.New(logged.NewTextSerializer(os.Stdout))
```

If you want to write a custom format, you only need to implement the `Serializer` interface:

```go
type Serializer interface {
	Write(e *Entry) error
}
```

### Default data

The log is designed to write metadata as a set of key/value pairs. This means that each log write can be optionally accompanied by a `map[string]string` that will be written with the message. Often times, there is certain data that should be written out with every write by default. You can specify the default data when you create the log:

```go
log := logged.NewOpts(logged.NewJSONSerializer(os.Stdout), logged.Opts{
	Defaults: map[string]string{
		"app_name": "fldsmdfr",
		"version":  "1.0.0",
	},
})
```

### Debug packages

Logged allows you to specific only certain packages to write debug logs. By passing debug packages you implicitly enable debug logging. Passing nothing or an empty slice disables debug logging:

```go
// Only allow github.com/foo/bar to write debug logs
log := logged.New(logged.NewJSONSerializer(os.Stdout), logged.Opts{
	DebugPackages: []string{"github.com/foo/bar"},
})

// Allow all packages to write debug logs
log := logged.New(logged.NewJSONSerializer(os.Stdout), logged.Opts{
	DebugPackages: []string{"*"},
})
```

### Spawning a sub-log

A log can be used to spawn a sub-log with all default metadata merged together:

```go
subLog := log.New(map[string]string{
	"default_for_subLog": "1234",
})
```

Now `subLog` can be used and the original defaults will be written out with every entry in addition to the defaults passed to `log.New`.

# Levels

A log has two levels: info and debug. Info should be used for actionable entries, debug should be used for everything else and can be disblaed on a per-package level.

The `Log` interface is very simple:

```go
type Log interface {
	Info(message string, data map[string]string) error
	Debug(message string, data map[string]string) error
	IsDebug() bool
}
```

You don't need to guard calls to `Debug` with `IsDebug` - rather use this function to guard against expensive computations in your code for data you want to pass to `Debug`:

```go
// You don't need to guard this since nothing is being computed for the data - this call
// is a no-op unless DebugPackages allows for it
log.Debug("something happened", map[string]string{
	"cheap_data": "abc123",
})

// If there are expensive computations, use IsDebug
if log.IsDebug() {
	log.Debug("something else happened", map[string]string{
		"expensive_data": someExpensiveFunc(),
	})
}
```

If you don't need to pass data, just pass `nil`:

```go
log.Info("an event", nil)
```

# Why another logging package?

Logged is designed to be extremely fast:

* A custom JSON serializer is provided that uses no reflection. 
* Key/value data must be a string which means that the serializer doesn't have to use reflection to write to the underlying writer.

Logged is designed to be minimalist:

* Only three log levels are provided.
* The `Log` interface provides only four functions - there is no extensive DSL or API to learn.

