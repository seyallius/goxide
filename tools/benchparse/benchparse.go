// Package main. benchparse.go parses raw `go test -bench` output and regenerates
// docs/content/benchmarks.md. Benchmarks are aggregated (median of N runs),
// bucketed into workload groups, measured against a configurable baseline
// style, and rendered as one Markdown table per group. The output file is
// fully overwritten on every run, keeping the pipeline idempotent.
package main

import (
	"fmt"
	"os"
)

// ------------------------------------------- <Main> ------------------------------------------- //

// main is the CLI entrypoint. Delegates entirely to Run().
func main() {
	cfg := DefaultConfig()
	if err := run(cfg); err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}
}

// -------------------------------------- Internal Helpers -------------------------------------- //

// run executes the full benchmark parsing pipeline.
// This is the primary entry point for external callers and tests.
func run(cfg Config) error {
	rows, err := ParseRawFile(cfg.RawFile, cfg)
	if err != nil {
		return fmt.Errorf("parse raw file: %w", err)
	}
	if len(rows) == 0 {
		return fmt.Errorf("no benchmark lines found in %s", cfg.RawFile)
	}

	aggregated, runs := AggregateByMedian(rows)
	groups := GroupByWorkload(aggregated, cfg)
	maxRunCount := MaxMapValue(runs)

	if err := WriteMarkdownReport(cfg.MDFile, groups, maxRunCount, cfg); err != nil {
		return fmt.Errorf("write markdown: %w", err)
	}

	fmt.Println("✅ regenerated", cfg.MDFile)
	return nil
}
