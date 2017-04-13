package main

import (
	"os"

	"github.com/andrewhare/logged"
)

func main() {
	log := logged.New(&logged.Config{
		Serializer: logged.NewJSONSerializer(os.Stdout),
		Defaults: logged.Data{
			"app_name": "fldsmdfr",
			"version":  "1.2.3.4",
		},
	})

	// Log just a message
	log.Info("an info message", nil)

	// Log a message with extended data
	log.Info("an info message with data", logged.Data{
		"some_number": "111",
		"some_string": "abc",
	})

	// Guard all debug statements to prevent expensive
	// computation from running
	if log.IsDebug() {
		log.Debug("a debug message", nil)
		log.Debug("a debug message with data", logged.Data{
			"some_date": "Tue Apr 11 11:47:48 EDT 2017",
		})
	}
}
