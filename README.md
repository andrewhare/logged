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
log := logged.New(&logged.Config{
	Serializer: logged.NewJSONSerializer(os.Stdout),
})

// Text
log := logged.New(&logged.Config{
	Serializer: logged.NewTextSerializer(os.Stdout),
})
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
log := logged.New(&logged.Config{
	Serializer: logged.NewJSONSerializer(os.Stdout),
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
log := logged.New(&logged.Config{
	Serializer: logged.NewJSONSerializer(os.Stdout),
	DebugPackages: []string{"github.com/foo/bar"},
})

// Allow all packages to write debug logs
log := logged.New(&logged.Config{
	Serializer: logged.NewJSONSerializer(os.Stdout),
	DebugPackages: []string{"*"},
})
```

# Levels

A log has three levels: error, info, and debug. Error and info are synonymous except for the level they are written at which make it possible to parse errors from the log output.

The `Log` interface is very simple:

```go
type Log interface {
	Error(err error, data map[string]string) error
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

// If there is expesive computations, use IsDebug
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

