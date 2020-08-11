package parsers_test

import (
	"testing"

	"github.com/tj/assert"

	"github.com/apex/parsers"
)

// Test parsing AWS Lambda events.
func TestParse_lambda(t *testing.T) {
	t.Run("START", func(t *testing.T) {
		s := "START RequestId: f7172574-5884-44d9-95f4-7438fb83e9b0 Version: 26"
		v, _ := parsers.Parse(s)
		e := v.(*parsers.AWSLambdaStart)

		assert.Equal(t, "f7172574-5884-44d9-95f4-7438fb83e9b0", e.RequestID)
		assert.Equal(t, 26, e.Version)
	})

	t.Run("END", func(t *testing.T) {
		s := "END RequestId: f7172574-5884-44d9-95f4-7438fb83e9b0"
		v, _ := parsers.Parse(s)
		e := v.(*parsers.AWSLambdaEnd)

		assert.Equal(t, "f7172574-5884-44d9-95f4-7438fb83e9b0", e.RequestID)
	})

	t.Run("REPORT", func(t *testing.T) {
		s := "REPORT RequestId: 136f2f48-069e-4808-8d73-b31c4d97e146\tDuration: 7.80 ms\tBilled Duration: 100 ms\tMemory Size: 512 MB\tMax Memory Used: 115 MB\t\n"
		v, _ := parsers.Parse(s)
		e := v.(*parsers.AWSLambdaReport)

		assert.Equal(t, "136f2f48-069e-4808-8d73-b31c4d97e146", e.RequestID)
		assert.Equal(t, 7.8, e.Duration)
		assert.Equal(t, 100.0, e.BilledDuration)
		assert.Equal(t, 512, e.MemorySize)
		assert.Equal(t, 115, e.MaxMemoryUsed)
	})
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
