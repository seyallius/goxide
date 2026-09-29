package result_test

import (
	"os"
	"testing"
)

// TestProfile is a manual harness for generating CPU and memory profiles.
// It is skipped unless PROFILE is set, so it never runs during `go test ./...`.
//
// Usage (usually via `just profile`):
//
//	PROFILE=cpu  go test -run TestProfile -bench=BenchmarkResultWithTry -cpuprofile=cpu.prof
//	PROFILE=mem  go test -run TestProfile -bench=.                  -memprofile=mem.prof
//
// Or explicitly:
//
//	go test -run TestProfile -bench=BenchmarkResultWithTry \
//	    -cpuprofile=cpu.prof -memprofile=mem.prof ./rusty/result/
func TestProfile(t *testing.T) {
	if os.Getenv("PROFILE") == "" {
		t.Skip("set PROFILE=cpu or PROFILE=mem to run; see justfile")
	}
	// The actual work is done by the benchmarks selected via -bench=.
	// This test exists only so a -cpuprofile/-memprofile run has a target.
}
