//go:generate ldetool --package parsers --go-string parsers.lde

package parsers

// Event is the interface used to extract an event from a log line.
type Event interface {
	Extract(line string) (bool, error)
}

// Parse a log line. Returns true if an event was successfully parsed.
func Parse(line string) (Event, bool) {
	events := []Event{
		&AWSLambdaStart{},
		&AWSLambdaReportInit{},
		&AWSLambdaReport{},
		&AWSLambdaEnd{},
		&AWSLambdaTimeout{},
	}

	for _, e := range events {
		if ok, _ := e.Extract(line); ok {
			return e, true
		}
	}

	return nil, false
}
