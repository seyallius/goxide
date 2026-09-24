// Package main. parse.go reads raw benchmark output files and converts each
// valid line into a BenchResult. Handles file I/O, regex matching, and
// style/key detection. No aggregation or rendering logic belongs here.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ----------------------------------------- Public API ----------------------------------------- //

// ParseRawFile reads and parses all valid benchmark lines from a file.
func ParseRawFile(path string, cfg Config) ([]BenchResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var results []BenchResult
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		result, ok := parseLine(scanner.Text(), cfg)
		if !ok {
			continue
		}
		results = append(results, result)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning %s: %w", path, err)
	}
	return results, nil
}

// -------------------------------------- Internal Helpers -------------------------------------- //

// parseLine attempts to parse a single benchmark output line.
func parseLine(line string, cfg Config) (BenchResult, bool) {
	matches := LineRe.FindStringSubmatch(line)
	if matches == nil {
		return BenchResult{}, false
	}

	nsPerOp, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return BenchResult{}, false
	}

	iterations, _ := strconv.Atoi(matches[2])
	bytesPerOp, _ := atoiOrZero(matches[4])
	allocsPerOp, _ := atoiOrZero(matches[5])

	baseName := strings.TrimPrefix(matches[1], "Benchmark")
	style, groupKey := DetectStyle(baseName, cfg)

	return BenchResult{
		RawName:     matches[1],
		PrettyName:  FormatCamelCase(baseName),
		Style:       style,
		GroupKey:    groupKey,
		Iterations:  iterations,
		NsPerOp:     nsPerOp,
		BytesPerOp:  bytesPerOp,
		AllocsPerOp: allocsPerOp,
	}, true
}

// atoiOrZero parses an integer, returning 0 for empty optional groups.
func atoiOrZero(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}
