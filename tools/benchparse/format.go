// Package main. format.go provides shared formatting utilities for numbers,
// camelCase splitting, and percentage deltas. Used by both parse.go and
// render.go. Contains NO domain-specific logic or I/O.
package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ----------------------------------------- Public API ----------------------------------------- //

// FormatCamelCase splits CamelCase into readable space-separated words.
func FormatCamelCase(s string) string {
	s = CamelRe1.ReplaceAllString(s, "${1} ${2}")
	s = CamelRe2.ReplaceAllString(s, "${1} ${2}")
	return s
}

// FormatThousands adds comma separators to integers.
func FormatThousands(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}

	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return b.String()
}

// FormatNsOp formats ns/op values: commas for large, raw for small.
func FormatNsOp(f float64) string {
	if f >= 1000 {
		return FormatThousands(int(f))
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// CompareLabel renders percentage delta with emoji indicators.
// Deltas under 2% are considered noise.
func CompareLabel(ns, base float64) string {
	if base <= 0 {
		return "—"
	}
	delta := (ns - base) / base * 100
	switch {
	case delta <= -2:
		return fmt.Sprintf("⚡ %.1f%% faster", -delta)
	case delta >= 2:
		return fmt.Sprintf("🐢 %.1f%% slower", delta)
	default:
		return "≈ on par"
	}
}