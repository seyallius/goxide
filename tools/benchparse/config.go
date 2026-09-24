// Package main. config.go defines benchmark parsing configuration, domain types,
// and default values. This is the single source of truth for styles, grouping
// rules, and report categories. Change ONLY this file when adding new benchmarks.
package main

import "regexp"

// ----------------------------------- Types, Variables & Constants ----------------------------- //

// BenchResult holds one parsed benchmark measurement line.
type BenchResult struct {
	RawName     string
	PrettyName  string
	Style       string // implementation style token, e.g. "Result "
	GroupKey    string // workload key shared by comparable benchmarks
	Iterations  int
	NsPerOp     float64
	BytesPerOp  int
	AllocsPerOp int
}

// GroupRule maps name substrings to a shared workload key.
// First matching rule wins; Exclude prevents false positives.
type GroupRule struct {
	Contains string
	Exclude  string // if non-empty, names containing this are NOT matched
	Key      string
}

// Category defines a report section for grouping workload keys.
type Category struct {
	Contains string
	Title    string
}

// Config holds all tunable parameters for the benchmark parser.
type Config struct {
	RawFile    string
	MDFile     string
	Styles     []string // Styles[0] is always the baseline
	GroupRules []GroupRule
	Categories []Category
	BasicTitle string
}

// LineRe matches: BenchmarkName-12  1000  1234.5 ns/op  56 B/op  7 allocs/op
var LineRe = regexp.MustCompile(`^(Benchmark[a-zA-Z0-9_]+)(?:-\d+)?\s+(\d+)\s+([\d.eE+]+)\s+ns/op\s+(?:(\d+)\s+B/op\s+)?(?:(\d+)\s+allocs/op)?`)

// CamelRe1/CamelRe2 split CamelCase names into readable words.
var (
	CamelRe1 = regexp.MustCompile("([a-z0-9])([A-Z])")
	CamelRe2 = regexp.MustCompile("([A-Z]+)([A-Z][a-z])")
)

// ----------------------------------------- Public API ----------------------------------------- //

// DefaultConfig returns the standard configuration for this project.
// This is the ONLY place to touch when adding new styles or categories.
func DefaultConfig() Config {
	return Config{
		RawFile: "docs/benchmarks_raw.txt",
		MDFile:  "docs/content/benchmarks.md",
		Styles:  []string{"Traditional ", "Result ", "Option ", "BubbleUp"},
		GroupRules: []GroupRule{
			{Contains: "ErrorHandling", Exclude: "DB", Key: "ErrorHandlingCPU"},
			{Contains: "WithTry", Exclude: "DB", Key: "ErrorHandlingCPU"},
			{Contains: "WithAndThen", Exclude: "DB", Key: "ErrorHandlingCPU"},
		},
		Categories: []Category{
			{Contains: "DB", Title: "🗄️ Database Operations"},
			{Contains: "Chain", Title: "🔗 Chaining & Pipelines"},
			{Contains: "Error", Title: "🛟 Error Handling & Recovery"},
			{Contains: "Fallback", Title: "🛟 Error Handling & Recovery"},
			{Contains: "Catch", Title: "🛟 Error Handling & Recovery"},
		},
		BasicTitle: "⚡ Basic Operations",
	}
}
