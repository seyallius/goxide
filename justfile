# Use the reusable benchmark/profile tool. Install it with:
#   go install github.com/seyallius/benchbook/cmd/benchbook@latest
set shell := ["bash", "-c"]

benchbook := env_var_or_default("BENCHBOOK", "benchbook")

default:
    @just --list

[group('Development')]
treeclip dir="":
    treeclip run {{ dir }} -f -t -v -c --stats

[group('Development')]
serve-docs port="8000":
    @echo "📚 Serving docs at http://localhost:{{ port }}"
    @echo "   (Press Ctrl+C to stop)"
    cd docs && python -m http.server {{ port }}

[group('Development')]
test:
    go test ./...

[group('Development')]
test-cover:
    go test -cover ./...

# Record benchmark samples and rebuild the generated performance pages.
[group('Profiling')]
bench count="6":
    {{benchbook}} run --count {{count}}

# Capture CPU and heap profiles for the package configured in benchbook.toml.
[group('Profiling')]
profile bench="":
    {{benchbook}} profile{{ if bench != "" { " --bench " + bench } else { "" } }}

# Compatibility aliases. benchbook records all configured profile kinds together.
[group('Profiling')]
profile-cpu bench="":
    {{benchbook}} profile{{ if bench != "" { " --bench " + bench } else { "" } }}

[group('Profiling')]
profile-mem bench="":
    {{benchbook}} profile{{ if bench != "" { " --bench " + bench } else { "" } }}

# Regenerate all performance documentation through the reusable tool.
[group('Profiling')]
perf-docs:
    {{benchbook}} run
    {{benchbook}} profile
    {{benchbook}} build

[group('Git')]
amend:
    git commit -a --amend

[group('Git')]
empty:
    git commit --allow-empty

[group('Git')]
rebase n="3":
    git rebase -i HEAD~{{ n }}

[group('Git')]
[linux]
diff-cp:
    git diff HEAD | xclip -selection clipboard
