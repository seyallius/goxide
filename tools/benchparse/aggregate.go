// Package main. aggregate.go performs pure statistical transformations on
// parsed benchmark results: median aggregation, workload grouping, and
// style-based sorting. Contains NO I/O or configuration side effects.
package main

import (
	"sort"
	"strings"
)

// ----------------------------------------- Public API ----------------------------------------- //

// AggregateByMedian collapses repeated runs into their median value.
// Returns aggregated results and a map of run counts per benchmark name.
func AggregateByMedian(rows []BenchResult) ([]BenchResult, map[string]int) {
	byName := make(map[string][]BenchResult)
	order := make([]string, 0, len(rows))

	for _, r := range rows {
		if _, exists := byName[r.RawName]; !exists {
			order = append(order, r.RawName)
		}
		byName[r.RawName] = append(byName[r.RawName], r)
	}

	runCounts := make(map[string]int, len(order))
	aggregated := make([]BenchResult, 0, len(order))

	for _, name := range order {
		rs := byName[name]
		sort.Slice(rs, func(i, j int) bool {
			return rs[i].NsPerOp < rs[j].NsPerOp
		})
		runCounts[name] = len(rs)
		aggregated = append(aggregated, rs[len(rs)/2])
	}

	return aggregated, runCounts
}

// GroupByWorkload buckets results by GroupKey, sorted by style rank then name.
func GroupByWorkload(rows []BenchResult, cfg Config) map[string][]BenchResult {
	groups := make(map[string][]BenchResult)
	for _, r := range rows {
		groups[r.GroupKey] = append(groups[r.GroupKey], r)
	}

	for k := range groups {
		sort.SliceStable(groups[k], func(i, j int) bool {
			ri, rj := groups[k][i], groups[k][j]
			rankI := StyleRank(ri.Style, cfg.Styles)
			rankJ := StyleRank(rj.Style, cfg.Styles)
			if rankI != rankJ {
				return rankI < rankJ
			}
			return ri.RawName < rj.RawName
		})
	}

	return groups
}

// MaxMapValue returns the maximum value in a map[string]int.
func MaxMapValue(m map[string]int) int {
	max := 0
	for _, v := range m {
		if v > max {
			max = v
		}
	}
	return max
}

// -------------------------------------- Internal Helpers -------------------------------------- //

// DetectStyle extracts implementation style and workload key from a base name.
// Trailing style tokens take precedence over leading ones (more specific).
func DetectStyle(base string, cfg Config) (style, key string) {
	// Check trailing styles first (higher specificity)
	for _, s := range cfg.Styles {
		if len(base) > len(s) && strings.HasSuffix(base, s) {
			stripped := strings.TrimSuffix(base, s)
			return s, applyGroupRules(trimLeadingStyle(stripped, cfg.Styles), cfg.GroupRules)
		}
	}

	// Fall back to leading styles
	for _, s := range cfg.Styles {
		if strings.HasPrefix(base, s) {
			return s, applyGroupRules(strings.TrimPrefix(base, s), cfg.GroupRules)
		}
	}

	return "other", applyGroupRules(base, cfg.GroupRules)
}

// StyleRank returns sort priority for a style (lower = higher priority).
func StyleRank(style string, styles []string) int {
	for i, known := range styles {
		if style == known {
			return i
		}
	}
	return len(styles)
}

// trimLeadingStyle removes a leading style token if present.
func trimLeadingStyle(key string, styles []string) string {
	for _, s := range styles {
		if strings.HasPrefix(key, s) {
			return strings.TrimPrefix(key, s)
		}
	}
	return key
}

// applyGroupRules finds the first matching rule for a workload key.
func applyGroupRules(key string, rules []GroupRule) string {
	for _, r := range rules {
		if strings.Contains(key, r.Contains) {
			if r.Exclude == "" || !strings.Contains(key, r.Exclude) {
				return r.Key
			}
		}
	}
	return key
}
