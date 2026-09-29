#set shell := ["bash", "-c"]

# ------------------------------------------------------------------------------
# Variables
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# Default
# ------------------------------------------------------------------------------

# Default target: List available commands
default:
    @just --list

# ------------------------------------------------------------------------------
# Development
# ------------------------------------------------------------------------------

# Run treeclip with default flags.
[group('Development')]
treeclip dir="":
    treeclip run {{ dir }} -f -t -v -c --stats

# Serve the documentation site locally for testing
[group('Development')]
serve-docs port="8000":
    @echo "📚 Serving docs at http://localhost:{{ port }}"
    @echo "   (Press Ctrl+C to stop)"
    # We cd into docs/ so that index.html is at the root URL
    cd docs && python -m http.server {{ port }}

# Run tests
[group('Development')]
test:
    go test ./...

# Run tests with coverage
[group('Development')]
test-cover:
    go test -cover ./...

# ═══════════════════════════════════════════════════════════════════════════
# 🔬 Profiling
# ═══════════════════════════════════════════════════════════════════════════

# Capture hardware info (CPU, RAM, OS) used in benchmark and profile docs.
[group('Profiling')]
hwinfo:
    @mkdir -p docs
    go run ./tools/sysinfo/ | tee docs/hardware.txt

# Render docs/benchmarks_raw.txt into docs/content/benchmarks.md
[group('Profiling')]
benchparse:
    go run ./tools/benchparse/

# Render docs/profiles/*.txt into docs/content/profiling.md
[group('Profiling')]
profparse:
    go run ./tools/profparse/

# Run all benchmarks, save raw output, and regenerate docs/content/benchmarks.md
[group('Profiling')]
bench count="1":
    @echo "🏃 Running benchmarks..."
    @mkdir -p docs
    @just hwinfo
    go test -count={{ count }} -bench=. -benchmem ./rusty/... | tee docs/benchmarks_raw.txt || true
    @echo "📊 Parsing benchmarks into Markdown..."
    @just benchparse
    @echo "✅ Done! Updated docs/content/benchmarks.md"
    @echo "💡 Don't forget to commit the updated benchmarks.md!"

# CPU profile. Pass a benchmark name to profile just that one.
#
#   just profile-cpu                       → docs/profiles/cpu_full.txt
#   just profile-cpu BenchmarkResultWithTry → docs/profiles/cpu_BenchmarkResultWithTry.txt
[group('Profiling')]
profile-cpu bench="":
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p docs/profiles
    if [ -z "{{ bench }}" ]; then
        echo "🔬 CPU profiling FULL result benchmark suite..."
        suffix="full"
        filter="."
    else
        echo "🔬 CPU profiling {{ bench }}..."
        suffix="{{ bench }}"
        filter="{{ bench }}"
    fi
    PROFILE=cpu go test -run TestProfile -bench="$filter" \
        -benchmem -cpuprofile=cpu.prof ./rusty/result/
    out="docs/profiles/cpu_${suffix}.txt"
    go tool pprof -top -nodecount=20 cpu.prof | tee "$out"
    echo ""
    echo "✅ wrote $out"
    echo "Interactive:  go tool pprof -http=:8080 cpu.prof"

# Memory profile. Pass a benchmark name to profile just that one.
#
#   just profile-mem                       → docs/profiles/mem_{space,objects}_full.txt
#   just profile-mem BenchmarkResultWithTry → docs/profiles/mem_{space,objects}_BenchmarkResultWithTry.txt
[group('Profiling')]
profile-mem bench="":
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p docs/profiles
    if [ -z "{{ bench }}" ]; then
        echo "🔬 Memory profiling FULL result benchmark suite..."
        suffix="full"
        filter="."
    else
        echo "🔬 Memory profiling {{ bench }}..."
        suffix="{{ bench }}"
        filter="{{ bench }}"
    fi
    PROFILE=mem go test -run TestProfile -bench="$filter" -benchmem \
        -memprofile=mem.prof ./rusty/result/
    space="docs/profiles/mem_space_${suffix}.txt"
    objects="docs/profiles/mem_objects_${suffix}.txt"
    echo ""
    echo "Top allocators by bytes (alloc_space):"
    go tool pprof -top -alloc_space -nodecount=20 mem.prof | tee "$space"
    echo ""
    echo "Top allocators by count (alloc_objects):"
    go tool pprof -top -alloc_objects -nodecount=20 mem.prof | tee "$objects"
    echo ""
    echo "✅ wrote $space"
    echo "✅ wrote $objects"
    echo "Interactive:  go tool pprof -http=:8080 mem.prof"

# Re-render the alloc_objects view from an existing mem.prof (no test re-run).
[group('Profiling')]
profile-allocs bench="":
    #!/usr/bin/env bash
    set -euo pipefail
    test -f mem.prof || { echo "mem.prof not found; run 'just profile-mem' first"; exit 1; }
    mkdir -p docs/profiles
    suffix="${1:-full}"
    if [ -z "{{ bench }}" ]; then suffix="full"; else suffix="{{ bench }}"; fi
    out="docs/profiles/mem_objects_${suffix}.txt"
    go tool pprof -top -alloc_objects -nodecount=20 mem.prof | tee "$out"
    echo "✅ wrote $out"

# Regenerate profiling.md from whatever lives in docs/profiles/.
# Does NOT re-run any test — pure re-render. Safe to call any time.
[group('Profiling')]
profiling-md:
    @just profparse
    @echo "✅ Updated docs/content/profiling.md"

# Full-suite profiles only → regenerates profiling.md with the "full" sections.
# Per-benchmark files (if any) are left alone and still appear in the page.
[group('Profiling')]
profile-docs:
    @just hwinfo
    @just profile-cpu
    @just profile-mem
    @just profparse
    @echo "✅ Done! Updated docs/content/profiling.md"
    @echo "💡 Don't forget to commit the updated profiling.md!"

# Interactive pprof web UI for the CPU profile (requires profile-cpu first).
[group('Profiling')]
profile-web-cpu:
    go tool pprof -http=:8080 cpu.prof

# Interactive pprof web UI for the memory profile (requires profile-mem first).
[group('Profiling')]
profile-web-mem:
    go tool pprof -http=:8080 mem.prof

# Clean up profile artifacts.
[group('Profiling')]
profile-clean:
    -rm -f cpu.prof mem.prof *.prof
    @echo "🧹 profiles cleaned (docs/profiles/ untouched)"

# Regenerate ALL performance documentation (hardware, benchmarks, profiles).
# This is the only command you need to remember.
[group('Profiling')]
perf-docs:
    @echo "══════════════════════════════════════════════════════════════"
    @echo "📊 Regenerating performance documentation"
    @echo "══════════════════════════════════════════════════════════════"

    @echo ""
    @echo "──────────────────────────────────────────────────────────────"
    @echo "🖥️  [1/5] Capturing hardware info"
    @echo "──────────────────────────────────────────────────────────────"
    @just hwinfo

    @echo ""
    @echo "──────────────────────────────────────────────────────────────"
    @echo "🏃 [2/5] Running benchmarks"
    @echo "──────────────────────────────────────────────────────────────"
    @just bench

    @echo ""
    @echo "──────────────────────────────────────────────────────────────"
    @echo "🔬 [3/5] CPU profiling"
    @echo "──────────────────────────────────────────────────────────────"
    @just profile-cpu

    @echo ""
    @echo "──────────────────────────────────────────────────────────────"
    @echo "🧠 [4/5] Memory profiling"
    @echo "──────────────────────────────────────────────────────────────"
    @just profile-mem

    @echo ""
    @echo "──────────────────────────────────────────────────────────────"
    @echo "📝 [5/5] Rendering profiling.md"
    @echo "──────────────────────────────────────────────────────────────"
    @just profparse

    @echo ""
    @echo "══════════════════════════════════════════════════════════════"
    @echo "✅ Done."
    @echo "══════════════════════════════════════════════════════════════"
    @echo "   Updated:"
    @echo "     • docs/hardware.txt"
    @echo "     • docs/benchmarks_raw.txt"
    @echo "     • docs/profiles/cpu_full.txt"
    @echo "     • docs/profiles/mem_{space,objects}_full.txt"
    @echo "     • docs/content/benchmarks.md"
    @echo "     • docs/content/profiling.md"
    @echo "💡 Per-benchmark profiles (cpu_<name>.txt etc.) are NOT regenerated."
    @echo "   Run 'just profile-cpu <name>' to (re)generate one."

# ------------------------------------------------------------------------------
# Code Quality
# ------------------------------------------------------------------------------


# ------------------------------------------------------------------------------
# Git
# ------------------------------------------------------------------------------

# Commit staged changes with amend.
[group('Git')]
amend:
    git commit -a --amend

[group('Git')]
empty:
    git commit --allow-empty

# Rebase current branch to the specified number of commits (Usage: just rebase 5)
[group('Git')]
rebase n="3":
    git rebase -i HEAD~{{ n }}

[group('Git')]
[linux]
diff-cp:
    git diff HEAD | xclip -selection clipboard

[group('Git')]
[windows]
diff-cp:
    git diff HEAD | /c/Windows/System32/clip.exe

[group('Git')]
today:
    git log --since="today 00:00:00" --until="today 23:59:59" --oneline
