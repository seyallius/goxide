# Repository Guidelines

## Project Structure & Module Organization

This is a Go 1.25 module (`github.com/seyallius/goxide`) that provides Rust-inspired `Result`, `Option`, and chaining APIs.

- `rusty/chain`, `rusty/option`, `rusty/result`, and `rusty/types`: library packages and their tests/benchmarks.
- `rusty/examples`: runnable usage examples.
- `internal/tests`: shared database and test setup helpers.
- `tools/benchparse`, `tools/profparse`, and `tools/sysinfo`: benchmark, profile, and hardware-report tooling.
- `docs/content`: user guides and generated benchmark/profile pages; `docs/`: documentation-site assets and raw measurements.
- `testdata`: fixtures used by tests.

## Build, Test, and Development Commands

Run `go test ./...` for the full test suite and `go test -cover ./...` for coverage. Run `go test -bench=. ./rusty/... -benchmem` for benchmarks. `just test` and `just test-cover` provide the standard shortcuts.

Use `just serve-docs` to serve the documentation site locally. Regenerate performance documentation with `just perf-docs`; narrower commands include `just benchparse`, `just profparse`, `just profile-cpu`, and `just profile-mem`. These commands update files under `docs/` and `docs/content/`.

## Coding Style & Naming Conventions

Format Go changes with `gofmt` (for example, `gofmt -w rusty/result/result.go`). Follow idiomatic Go naming: exported identifiers use PascalCase, local variables use concise lowerCamelCase, and test files end in `_test.go`. Keep package boundaries focused and preserve the existing generic API style. No separate linter is configured; `go test ./...` is the baseline validation.

## Testing Guidelines

Add behavioral tests beside the package they cover, using descriptive `Test...` names. Add `Benchmark...` functions when measuring performance-sensitive paths. Database benchmarks use the shared helpers in `internal/tests/testkit.go`. Run targeted tests with commands such as `go test ./rusty/result -run TestName` before running the full suite.

## Commit & Pull Request Guidelines

Use short, imperative, scoped messages following the repository’s history, such as `feat(benchparse): ...`, `fix(result): ...`, or `docs(benchmarks): ...`. Pull requests should explain the behavior change, list validation commands, and include updated generated documentation when benchmarks or profiles change. Include screenshots for documentation-site or visual changes when useful.
