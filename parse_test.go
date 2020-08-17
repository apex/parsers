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
