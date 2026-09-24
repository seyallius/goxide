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

# Run all benchmarks, save raw output, and generate documentation Markdown
[group('Development')]
bench count="1":
    @echo "🏃 Running benchmarks..."
    @mkdir -p docs
    go test -count={{ count }} -bench=. -benchmem ./rusty/... | tee docs/benchmarks_raw.txt || true
    @echo "📊 Parsing benchmarks into Markdown..."
    go run ./tools/benchparse/main.go
    @echo "✅ Done! Updated docs/content/benchmarks.md"
    @echo "💡 Don't forget to commit the updated benchmarks.md!"

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
