// Package main. main.go parses raw `go test -bench` output and regenerates
// docs/content/benchmarks.md. Benchmarks are aggregated (median of N runs),
// bucketed into workload groups, measured against a configurable baseline
// style, and rendered as one Markdown table per group. The output file is
// fully overwritten on every run, keeping the pipeline idempotent.
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ----------------------------------- Types, Variables & Constants ----------------------------- //

// BenchResult holds one parsed benchmark measurement line.
type BenchResult struct {
	RawName     string
	PrettyName  string
	Style       string // implementation style token, e.g. "Result"
	GroupKey    string // workload key shared by comparable benchmarks
	Iterations  int
	NsPerOp     float64
	BytesPerOp  int
	AllocsPerOp int
}

// groupRule reroutes a family of benchmarks into a shared comparison group
// when their names do not follow the <Style><Workload> convention.
type groupRule struct {
	contains string
	key      string
}

// category is a report section: groups whose key contains `contains` are
// rendered under `title`.
type category struct {
	contains string
	title    string
}

// This block is the ONLY place you need to touch when new benchmark styles,
// families or sections appear. Everything else adapts automatically.

var (
	// styles lists the implementation style tokens found in benchmark names.
	//   - styles[0] is the BASELINE every other style is compared against.
	//   - Detection: a trailing token wins (most specific, e.g.
	//     ResultDBChainedOperationsBubbleUp -> "BubbleUp"), otherwise the
	//     leading token wins (e.g. ResultWithTry -> "Result").
	//   - Add new styles here (e.g. "Option", "Mutex", "Channel") and
	//     benchmarks like BenchmarkOptionSuccess pair up automatically.
	styles = []string{"Traditional", "Result", "Option", "BubbleUp"}

	// groupRules maps name substrings to a shared workload key. Substring
	// match on the style-stripped name, first rule wins. Use these for
	// families whose names diverge from the baseline's naming.
	groupRules = []groupRule{
		{contains: "ErrorHandling", key: "ErrorHandling"},
		{contains: "WithTry", key: "ErrorHandling"},
		{contains: "WithAndThen", key: "ErrorHandling"},
	}

	// categories defines the report sections, in display order.
	// First match wins; unmatched groups land in basicTitle.
	categories = []category{
		{contains: "DB", title: "🗄️ Database Operations"},
		{contains: "Chain", title: "🔗 Chaining & Pipelines"},
		{contains: "Error", title: "🛟 Error Handling & Recovery"},
		{contains: "Fallback", title: "🛟 Error Handling & Recovery"},
		{contains: "Catch", title: "🛟 Error Handling & Recovery"},
	}
	basicTitle = "⚡ Basic Operations"

	rawFile = "docs/benchmarks_raw.txt"
	mdFile  = "docs/content/benchmarks.md"

	// lineRe matches: BenchmarkName-12  1000  1234.5 ns/op  56 B/op  7 allocs/op
	lineRe = regexp.MustCompile(`^(Benchmark[a-zA-Z0-9_]+)(?:-\d+)?\s+(\d+)\s+([\d.eE+]+)\s+ns/op\s+(?:(\d+)\s+B/op\s+)?(?:(\d+)\s+allocs/op)?`)

	// camelRe1/camelRe2 split CamelCase names into readable words.
	camelRe1 = regexp.MustCompile("([a-z0-9])([A-Z])")
	camelRe2 = regexp.MustCompile("([A-Z]+)([A-Z][a-z])")
)

// ----------------------------------------- Public API ----------------------------------------- //

// main is the entrypoint: parse -> aggregate -> group -> compare -> write.
func main() {
	rows, err := parseRaw(rawFile)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}
	if len(rows) == 0 {
		fmt.Println("❌ no benchmark lines found in", rawFile)
		os.Exit(1)
	}

	aggregated, runs := aggregate(rows)
	groups := groupByKey(aggregated)

	if err := writeMarkdown(mdFile, groups, maxRuns(runs)); err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ regenerated", mdFile)
}

// -------------------------------------- Internal Helpers -------------------------------------- //

// parseRaw reads the raw benchmark file and returns every matching line.
func parseRaw(path string) ([]BenchResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []BenchResult
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		m := lineRe.FindStringSubmatch(sc.Text())
		if m == nil {
			continue
		}
		ns, _ := strconv.ParseFloat(m[3], 64)
		iters, _ := strconv.Atoi(m[2])
		b, _ := atoiOrZero(m[4])
		a, _ := atoiOrZero(m[5])
		base := strings.TrimPrefix(m[1], "Benchmark")
		style, key := detectStyle(base)
		out = append(out, BenchResult{
			RawName: m[1], PrettyName: formatName(base),
			Style: style, GroupKey: key,
			Iterations: iters, NsPerOp: ns, BytesPerOp: b, AllocsPerOp: a,
		})
	}
	return out, sc.Err()
}

// atoiOrZero parses an integer, returning 0 for empty optional groups.
func atoiOrZero(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

// detectStyle splits a benchmark base name into its implementation style and
// workload key. A trailing style token wins over a leading one because it is
// the more specific modifier: ResultDBChainedOperationsBubbleUp is a BubbleUp
// benchmark, not merely a Result one.
func detectStyle(base string) (style, key string) {
	for _, s := range styles {
		if len(base) > len(s) && strings.HasSuffix(base, s) {
			return s, applyGroupRules(trimStylePrefix(strings.TrimSuffix(base, s)))
		}
	}
	for _, s := range styles {
		if strings.HasPrefix(base, s) {
			return s, applyGroupRules(strings.TrimPrefix(base, s))
		}
	}
	return "other", applyGroupRules(base)
}

// trimStylePrefix removes a leading style token from a workload key.
func trimStylePrefix(key string) string {
	for _, s := range styles {
		if strings.HasPrefix(key, s) {
			return strings.TrimPrefix(key, s)
		}
	}
	return key
}

// applyGroupRules reroutes name families into shared workload keys.
func applyGroupRules(key string) string {
	for _, r := range groupRules {
		if strings.Contains(key, r.contains) {
			return r.key
		}
	}
	return key
}

// aggregate collapses repeated runs of the same benchmark into its median
// row (middle element after sorting by ns/op). Returns per-name run counts.
func aggregate(rows []BenchResult) ([]BenchResult, map[string]int) {
	byName := map[string][]BenchResult{}
	order := []string{}
	for _, r := range rows {
		if _, ok := byName[r.RawName]; !ok {
			order = append(order, r.RawName)
		}
		byName[r.RawName] = append(byName[r.RawName], r)
	}

	runs := map[string]int{}
	out := make([]BenchResult, 0, len(order))
	for _, name := range order {
		rs := byName[name]
		sort.Slice(rs, func(i, j int) bool { return rs[i].NsPerOp < rs[j].NsPerOp })
		runs[name] = len(rs)
		out = append(out, rs[len(rs)/2])
	}
	return out, runs
}

// maxRuns returns the highest run count seen for any benchmark.
func maxRuns(runs map[string]int) int {
	m := 0
	for _, n := range runs {
		if n > m {
			m = n
		}
	}
	return m
}

// groupByKey buckets aggregated rows by workload key, baseline first.
func groupByKey(rows []BenchResult) map[string][]BenchResult {
	g := map[string][]BenchResult{}
	for _, r := range rows {
		g[r.GroupKey] = append(g[r.GroupKey], r)
	}
	for k := range g {
		sort.SliceStable(g[k], func(i, j int) bool {
			ri, rj := g[k][i], g[k][j]
			if ri.Style != rj.Style {
				return styleRank(ri.Style) < styleRank(rj.Style)
			}
			return ri.RawName < rj.RawName
		})
	}
	return g
}

// styleRank orders styles within a group: baseline first, then the rest in
// configuration order, unknown styles last.
func styleRank(s string) int {
	for i, known := range styles {
		if s == known {
			return i
		}
	}
	return len(styles)
}

// categoryFor assigns a workload key to a report section.
func categoryFor(key string) string {
	for _, c := range categories {
		if strings.Contains(key, c.contains) {
			return c.title
		}
	}
	return basicTitle
}

// sectionTitles returns the report sections in display order.
func sectionTitles() []string {
	titles := []string{basicTitle}
	seen := map[string]bool{basicTitle: true}
	for _, c := range categories {
		if !seen[c.title] {
			seen[c.title] = true
			titles = append(titles, c.title)
		}
	}
	return titles
}

// compareLabel renders the percentage delta of a row against its baseline.
// Deltas under 2% are reported as "on par" because benchmark noise on a
// shared CPU routinely exceeds that.
func compareLabel(ns, base float64) string {
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

// formatName turns CamelCase benchmark suffixes into readable words.
func formatName(base string) string {
	s := camelRe1.ReplaceAllString(base, "${1} ${2}")
	s = camelRe2.ReplaceAllString(s, "${1} ${2}")
	return s
}

// formatInt inserts thousands separators for readability.
func formatInt(n int) string {
	s := strconv.Itoa(n)
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return b.String()
}

// formatNs renders ns/op with commas for large values, raw for small ones.
func formatNs(f float64) string {
	if f >= 1000 {
		return formatInt(int(f))
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// writeMarkdown regenerates the benchmarks page. Each workload group gets its
// own heading and its own complete table — a blank line inside a GFM table
// terminates it, which is what broke the previous layout. The file is fully
// overwritten, so the pipeline stays idempotent no matter how often it runs.
func writeMarkdown(path string, groups map[string][]BenchResult, runs int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()
	p := func(format string, a ...any) { fmt.Fprintf(w, format, a...) }

	p("<!-- benchmarks.md is auto-generated by `just bench`. DO NOT EDIT. -->\n\n")
	p("# ⚡ Performance Benchmarks\n\n")
	p("Generated from `go test -bench=. -benchmem` via `just bench`.\n")
	if runs > 1 {
		p("Values shown are the **median of %d runs** per benchmark.\n", runs)
	}
	p("Lower `ns/op`, `B/op` and `allocs/op` are better.\n\n")
	p("Each table compares one workload across implementation styles; the\n")
	p("**vs %s** column is the delta against the baseline style.\n", styles[0])
	p("Deltas under 2%% are reported as *on par* (benchmark noise).\n\n")

	for _, title := range sectionTitles() {
		keys := []string{}
		for k := range groups {
			if categoryFor(k) == title {
				keys = append(keys, k)
			}
		}
		if len(keys) == 0 {
			continue
		}
		sort.Strings(keys)

		p("## %s\n\n", title)
		for _, k := range keys {
			rows := groups[k]
			base, hasBase := 0.0, false
			for _, r := range rows {
				if r.Style == styles[0] {
					base, hasBase = r.NsPerOp, true
					break
				}
			}

			p("### %s\n\n", formatName(k))
			p("| Benchmark | Iterations | ns/op | B/op | allocs/op | vs %s |\n", styles[0])
			p("|---|---:|---:|---:|---:|---|\n")
			for _, r := range rows {
				cmp := "—"
				switch {
				case r.Style == styles[0]:
					cmp = "baseline"
				case hasBase:
					cmp = compareLabel(r.NsPerOp, base)
				}
				p("| `%s` | %s | %s | %s | %s | %s |\n",
					r.PrettyName, formatInt(r.Iterations), formatNs(r.NsPerOp),
					formatInt(r.BytesPerOp), formatInt(r.AllocsPerOp), cmp)
			}
			p("\n")
		}
	}
	return nil
}
