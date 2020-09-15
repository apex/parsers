package parsers_test

import (
	"testing"

	"github.com/tj/assert"

	"github.com/apex/parsers"
)

var cases = []struct {
	Label  string
	Input  string
	Output parsers.Event
}{
	{
		Label: "Lambda start",
		Input: "START RequestId: f7172574-5884-44d9-95f4-7438fb83e9b0 Version: 26",
		Output: &parsers.AWSLambdaStart{
			RequestID: "f7172574-5884-44d9-95f4-7438fb83e9b0",
			Version:   "26",
		},
	},
	{
		Label: "Lambda start with $LATEST",
		Input: "START RequestId: f7172574-5884-44d9-95f4-7438fb83e9b0 Version: $LATEST",
		Output: &parsers.AWSLambdaStart{
			RequestID: "f7172574-5884-44d9-95f4-7438fb83e9b0",
			Version:   "$LATEST",
		},
	},
	{
		Label: "Lambda end",
		Input: "END RequestId: f7172574-5884-44d9-95f4-7438fb83e9b0",
		Output: &parsers.AWSLambdaEnd{
			RequestID: "f7172574-5884-44d9-95f4-7438fb83e9b0",
		},
	},
	{
		Label: "Lambda report",
		Input: "REPORT RequestId: 136f2f48-069e-4808-8d73-b31c4d97e146\tDuration: 7.80 ms\tBilled Duration: 100 ms\tMemory Size: 512 MB\tMax Memory Used: 115 MB\t",
		Output: &parsers.AWSLambdaReport{
			RequestID:      "136f2f48-069e-4808-8d73-b31c4d97e146",
			Duration:       7.8,
			BilledDuration: 100,
			MemorySize:     512,
			MaxMemoryUsed:  115,
		},
	},
	{
		Label: "Lambda report with init duration",
		Input: "REPORT RequestId: 136f2f48-069e-4808-8d73-b31c4d97e146\tDuration: 7.80 ms\tBilled Duration: 100 ms\tMemory Size: 512 MB\tMax Memory Used: 115 MB\tInit Duration: 185.62 ms\t",
		Output: &parsers.AWSLambdaReportInit{
			RequestID:      "136f2f48-069e-4808-8d73-b31c4d97e146",
			Duration:       7.8,
			BilledDuration: 100,
			InitDuration:   185.62,
			MemorySize:     512,
			MaxMemoryUsed:  115,
		},
	},
	{
		Label: "Lambda timeout",
		Input: "2020-08-19T09:20:47.075Z 8173dbda-4443-4bcd-8d4c-33704efa0f05 Task timed out after 30.03 seconds",
		Output: &parsers.AWSLambdaTimeout{
			Timestamp: "2020-08-19T09:20:47.075Z",
			RequestID: "8173dbda-4443-4bcd-8d4c-33704efa0f05",
			Duration:  30.03,
		},
	},
	{
		Label: "Heroku syslog",
		Input: "<45>1 2020-08-28T10:38:06.285004+00:00 host app api - Some random message here",
		Output: &parsers.Syslog{
			Priority:      45,
			SyslogVersion: 1,
			Timestamp:     "2020-08-28T10:38:06.285004+00:00",
			Hostname:      "host",
			Appname:       "app",
			ProcID:        "api",
			MsgID:         "-",
			Message:       "Some random message here",
		},
	},
	{
		Label:  "Unmatched",
		Input:  `{ "some": "json" }`,
		Output: nil,
	},
	{
		Label:  "Empty",
		Input:  ``,
		Output: nil,
	},
}

// Test parsing.
func TestParse(t *testing.T) {
	for _, c := range cases {
		t.Run(c.Label, func(t *testing.T) {
			v, _ := parsers.Parse(c.Input)
			assert.Equal(t, c.Output, v)
		})
	}
}

var herokuCases = []struct {
	Label  string
	Input  string
	Output parsers.Event
}{
	{
		Label: "Heroku deployment",
		Input: "Deploy 059375fe by user tj@apex.sh",
		Output: &parsers.HerokuDeploy{
			Commit: "059375fe",
			User:   "tj@apex.sh",
		},
	},
	{
		Label: "Heroku release",
		Input: "Release v16 created by user tj@apex.sh",
		Output: &parsers.HerokuRelease{
			Version: "v16",
			User:    "tj@apex.sh",
		},
	},
	{
		Label: "Heroku rollback",
		Input: "Rollback to v11 by user tj@apex.sh",
		Output: &parsers.HerokuRollback{
			Version: "v11",
			User:    "tj@apex.sh",
		},
	},
	{
		Label: "Heroku build start",
		Input: "Build started by user tj@apex.sh",
		Output: &parsers.HerokuBuild{
			User: "tj@apex.sh",
		},
	},
	{
		Label: "Heroku state change",
		Input: "State changed from starting to crashed",
		Output: &parsers.HerokuStateChange{
			From: "starting",
			To:   "crashed",
		},
	},
	{
		Label: "Heroku process exit",
		Input: "Process exited with status 143",
		Output: &parsers.HerokuProcessExit{
			Status: 143,
		},
	},
	{
		Label: "Heroku starting process",
		Input: "Starting process with command `node index.js`",
		Output: &parsers.HerokuProcessStart{
			Command: "node index.js",
		},
	},
	{
		Label: "Heroku listening",
		Input: "Listening on 55766",
		Output: &parsers.HerokuProcessListening{
			Port: 55766,
		},
	},
	{
		Label: "Heroku set env var",
		Input: "Set FOO config vars by user tj@apex.sh",
		Output: &parsers.HerokuConfigSet{
			Variables: "FOO",
			User:      "tj@apex.sh",
		},
	},
	{
		Label: "Heroku set env vars",
		Input: "Set FOO, BAR config vars by user tj@apex.sh",
		Output: &parsers.HerokuConfigSet{
			Variables: "FOO, BAR",
			User:      "tj@apex.sh",
		},
	},
	{
		Label: "Heroku remove env vars",
		Input: "Remove FOO config vars by user tj@apex.sh",
		Output: &parsers.HerokuConfigRemove{
			Variables: "FOO",
			User:      "tj@apex.sh",
		},
	},
	{
		Label: "Heroku scale 0 free",
		Input: "Scaled to web@0:Free by user tj@apex.sh",
		Output: &parsers.HerokuScale{
			Dynos: "web@0:Free",
			User:  "tj@apex.sh",
		},
	},
	{
		Label: "Heroku scale 1 free",
		Input: "Scaled to web@1:Free by user tj@apex.sh",
		Output: &parsers.HerokuScale{
			Dynos: "web@1:Free",
			User:  "tj@apex.sh",
		},
	},
	{
		Label: "Heroku scale multiple",
		Input: "Scaled to web@1:Free worker@0:Free by user tj@apex.sh",
		Output: &parsers.HerokuScale{
			Dynos: "web@1:Free worker@0:Free",
			User:  "tj@apex.sh",
		},
	},
}

// Test parsing Heroku messages.
func TestParseHeroku(t *testing.T) {
	for _, c := range herokuCases {
		t.Run(c.Label, func(t *testing.T) {
			v, _ := parsers.ParseHeroku(c.Input)
			assert.Equal(t, c.Output, v)
		})
	}
}

// Benchmark parsing.
func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := "REPORT RequestId: 136f2f48-069e-4808-8d73-b31c4d97e146\tDuration: 7.80 ms\tBilled Duration: 100 ms\tMemory Size: 512 MB\tMax Memory Used: 115 MB\t\n"
		_, ok := parsers.Parse(s)
		if !ok {
			b.Fatal("failed parsing")
		}
	}
}

// Benchmark parsing.
func BenchmarkParseHeroku(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := "Set FOO, BAR config vars by user tj@apex.sh"
		_, ok := parsers.ParseHeroku(s)
		if !ok {
			b.Fatal("failed parsing")
		}
	}
}
