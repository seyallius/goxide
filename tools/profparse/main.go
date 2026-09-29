// Command profparse converts `go tool pprof -top` text dumps from
// docs/profiles/*.txt into docs/content/profiling.md.
//
// File naming convention:
//
//	cpu_full.txt                             full-suite CPU profile
//	cpu_<BenchmarkName>.txt                  CPU profile of one benchmark
//	mem_space_full.txt                       full-suite alloc_space
//	mem_space_<BenchmarkName>.txt            alloc_space of one benchmark
//	mem_objects_full.txt                     full-suite alloc_objects
//	mem_objects_<BenchmarkName>.txt          alloc_objects of one benchmark
//
// Run via `just profparse` (or `just perf-docs`).
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// --------------------------------- Types, Constants & Variables ------------------------------- //

// pprofRow matches one data row of `go tool pprof -top` output.
var pprofRow = regexp.MustCompile(
	`^\s*([\d.]+)([a-zA-Z]*)\s+([\d.]+%)\s+([\d.]+%)\s+([\d.]+)([a-zA-Z]*)\s+([\d.]+%)\s+(.+?)\s*$`,
)

// filePattern matches cpu_full.txt / mem_space_BenchmarkResultWithTry.txt etc.
var filePattern = regexp.MustCompile(
	`^(cpu|mem_space|mem_objects)_(.+)\.txt$`,
)

type section struct {
	// Kind is one of "cpu", "mem_space", "mem_objects".
	Kind string
	// Target is "full" or a benchmark name.
	Target   string
	Title    string
	Source   string
	Metadata []string
	Rows     []pprofRowData
}

type pprofRowData struct {
	Flat, FlatPct, SumPct, Cum, CumPct, Name string
}

// humanKind maps the filename prefix to a display label.
func humanKind(kind string) string {
	switch kind {
	case "cpu":
		return "CPU Profile"
	case "mem_space":
		return "Memory — alloc_space (bytes)"
	case "mem_objects":
		return "Memory — alloc_objects (count)"
	}
	return kind
}

// ------------------------------------------- <Main> ------------------------------------------- //

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// -------------------------------------- Internal Helpers -------------------------------------- //

func run() error {
	entries, err := os.ReadDir("docs/profiles")
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("docs/profiles/ does not exist; run `just profile-cpu` / `just profile-mem` first")
		}
		return err
	}

	var sections []section
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		m := filePattern.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		kind, target := m[1], m[2]
		src := filepath.Join("docs/profiles", e.Name())
		data, err := os.ReadFile(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skipping %s: %v\n", src, err)
			continue
		}
		s := parse(kind, target, src, string(data))
		sections = append(sections, s)
	}

	if len(sections) == 0 {
		return fmt.Errorf("no profile sources matched docs/profiles/{cpu,mem_space,mem_objects}_*.txt")
	}

	// Ordering: full-suite sections first, then per-benchmark.
	// Within each group, keep a stable order (cpu, mem_space, mem_objects)
	// and then alphabetical by target.
	sort.SliceStable(sections, func(i, j int) bool {
		a, b := sections[i], sections[j]
		aFull, bFull := a.Target == "full", b.Target == "full"
		if aFull != bFull {
			return aFull // full first
		}
		if a.Kind != b.Kind {
			return kindOrder(a.Kind) < kindOrder(b.Kind)
		}
		return a.Target < b.Target
	})

	hw, _ := os.ReadFile("docs/hardware.txt")
	md := render(sections, strings.TrimRight(string(hw), "\n"))

	if err := os.MkdirAll("docs/content", 0o755); err != nil {
		return err
	}
	out := "docs/content/profiling.md"
	if err := os.WriteFile(out, []byte(md), 0o644); err != nil {
		return err
	}
	fmt.Printf("✅ wrote %s (%d sections)\n", out, len(sections))
	return nil
}

func kindOrder(kind string) int {
	switch kind {
	case "cpu":
		return 0
	case "mem_space":
		return 1
	case "mem_objects":
		return 2
	}
	return 99
}

func parse(kind, target, source, raw string) section {
	s := section{
		Kind:   kind,
		Target: target,
		Source: source,
	}
	if target == "full" {
		s.Title = humanKind(kind) + " — full suite"
	} else {
		s.Title = humanKind(kind) + " — " + target
	}

	inTable := false
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "flat") &&
			strings.Contains(line, "cum%") {
			inTable = true
			continue
		}

		if !inTable {
			if idx := strings.Index(line, ":"); idx > 0 && idx < 30 {
				switch strings.TrimSpace(line[:idx]) {
				case "File", "Type", "Time", "Duration":
					s.Metadata = append(s.Metadata, line)
					continue
				}
			}
			if strings.HasPrefix(line, "Showing ") || strings.HasPrefix(line, "Dropped ") {
				s.Metadata = append(s.Metadata, line)
			}
			continue
		}

		m := pprofRow.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		s.Rows = append(s.Rows, pprofRowData{
			Flat:    m[1] + m[2],
			FlatPct: m[3],
			SumPct:  m[4],
			Cum:     m[5] + m[6],
			CumPct:  m[7],
			Name:    m[8],
		})
	}
	return s
}

func render(sections []section, hardware string) string {
	var b strings.Builder

	b.WriteString("<!-- profiling.md is auto-generated by `just profparse`. DO NOT EDIT. -->\n\n")
	b.WriteString("# 🔬 Profiling Reports\n\n")
	b.WriteString("Generated from `go tool pprof -top` via `just profile-docs` ")
	b.WriteString("(full suite) or `just profile-cpu <name>` / `just profile-mem <name>` ")
	b.WriteString("(per-benchmark). Each table lists the **top 20** functions for that ")
	b.WriteString("metric. `flat` = time/bytes spent **in** the function; ")
	b.WriteString("`cum` = total including callees.\n\n")

	if hardware != "" {
		b.WriteString("## 🖥️ Test Environment\n\n")
		b.WriteString("```\n")
		b.WriteString(hardware)
		b.WriteString("\n```\n\n")
	}

	b.WriteString(fmt.Sprintf("_Generated: %s_\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Group: full suite first, then per-benchmark.
	writeGroup := func(group string, pred func(section) bool) {
		var members []section
		for _, s := range sections {
			if pred(s) {
				members = append(members, s)
			}
		}
		if len(members) == 0 {
			return
		}
		b.WriteString("## " + group + "\n\n")
		for _, s := range members {
			writeSection(&b, s)
		}
	}

	writeGroup("Full Suite", func(s section) bool { return s.Target == "full" })
	writeGroup("Per-Benchmark", func(s section) bool { return s.Target != "full" })

	return b.String()
}

func writeSection(b *strings.Builder, s section) {
	b.WriteString("### " + s.Title + "\n\n")

	if len(s.Metadata) > 0 {
		b.WriteString("<details><summary>Run metadata</summary>\n\n```\n")
		for _, m := range s.Metadata {
			b.WriteString(m + "\n")
		}
		b.WriteString("```\n\n</details>\n\n")
	}

	if len(s.Rows) == 0 {
		b.WriteString("_No rows parsed from the pprof output._\n\n")
		return
	}

	b.WriteString("| Function | flat | flat% | cum | cum% |\n")
	b.WriteString("|---|---:|---:|---:|---:|\n")
	for _, r := range s.Rows {
		b.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %s |\n",
			r.Name, r.Flat, r.FlatPct, r.Cum, r.CumPct))
	}
	b.WriteString("\n")
}
