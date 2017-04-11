# logged

A ridiculously simple Go logging package that has a minimalist interface, fast JSON output, and follows the thought process of [Dave Cheney's discussion on logging](https://dave.cheney.net/2015/11/05/lets-talk-about-logging) as close as possible.

# Installation

```
$ go get github.com/andrewhare/logged
```

# Usage

```go
package main

import (
	"os"

	"github.com/andrewhare/logged"
)

func main() {
	log := logged.New(&logged.Config{
		Writer: os.Stdout,
		Defaults: logged.Data{
			"app_name": "fldsmdfr",
			"version":  "1.2.3.4",
		},
	})

	// Log just a message
	log.Info("an info message")

	// Log a message with extended data
	log.InfoEx("an info message with data", logged.Data{
		"some_number": "111",
		"some_string": "abc",
	})

	// Guard all debug statements to prevent expensive
	// computation from running
	if log.IsDebug() {
		log.Debug("a debug message")
		log.DebugEx("a debug message with data", logged.Data{
			"some_date": "Tue Apr 11 11:47:48 EDT 2017",
		})
	}
}
```
