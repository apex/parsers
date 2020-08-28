//go:generate ldetool --package parsers --go-string parsers.lde

package parsers

// Event is the interface used to extract an event from a log line.
type Event interface {
	Extract(line string) (bool, error)
}

// Parse a log line from any source. Returns true if an event was successfully parsed.
func Parse(line string) (Event, bool) {
	events := []Event{
		&AWSLambdaStart{},
		&AWSLambdaReportInit{},
		&AWSLambdaReport{},
		&AWSLambdaEnd{},
		&AWSLambdaTimeout{},
		&Syslog{},
	}

	for _, e := range events {
		if ok, _ := e.Extract(line); ok {
			return e, true
		}
	}

	return nil, false
}

// ParseLambda parses a log line from AWS Lambda. Returns true if an event was successfully parsed.
func ParseLambda(line string) (Event, bool) {
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

// ParseHeroku parses a log line from Heroku. Returns true if an event was successfully parsed. You should
// first parse the syslog line from Heroku using Syslog, and then ParseHeroku() for the platform specifi
// message.
func ParseHeroku(line string) (Event, bool) {
	events := []Event{
		&HerokuDeploy{},
		&HerokuRollback{},
		&HerokuBuild{},
		&HerokuRelease{},
		&HerokuProcessExit{},
		&HerokuProcessStart{},
		&HerokuStateChange{},
		&HerokuProcessListening{},
		&HerokuConfigSet{},
		&HerokuConfigRemove{},
	}

	for _, e := range events {
		if ok, _ := e.Extract(line); ok {
			return e, true
		}
	}

	return nil, false
}
