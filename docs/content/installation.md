<!-- installation.md covers how to add Goxide to a Go project, verify the -->
<!-- installation, and understand the module layout. -->

# 📥 Installation

## Requirements

- **Go 1.25+** (Goxide uses generics extensively)
- No external dependencies — Goxide is a pure-Go library

## Add to Your Project

```bash
go get github.com/seyallius/goxide
```

This adds the module to your `go.mod`:

```
require github.com/seyallius/goxide v0.14.0
```

## Import Individual Packages

You only import what you need. Each package is independent:

```go
// Package main. main shows importing individual Goxide packages.
package main

import (
    "github.com/seyallius/goxide/rusty/result"  // Result[T] + BubbleUp + Catch
    "github.com/seyallius/goxide/rusty/option"  // Option[T] for optional values
    "github.com/seyallius/goxide/rusty/chain"   // Fluent method chaining
    "github.com/seyallius/goxide/rusty/types"   // Id, Return, Compose helpers
)
```

## Module Layout

```
github.com/seyallius/goxide/
├── rusty/
│   ├── result/     ← Result[T], BubbleUp(), Catch(), CatchWith(), Fallback()
│   ├── option/     ← Option[T], Some(), None(), UnwrapOr()
│   ├── chain/      ← Chain(), Map(), AndThen(), Unwrap()
│   ├── types/      ← Id(), Return(), Compose()
│   └── examples/   ← Worked examples (database, HTTP, validation)
├── reflect/        ← Type-safe struct reflection utilities
└── internal/
    └── tests/      ← Integration test helpers
```

## Verify Installation

Create a quick smoke test:

```go
// Package main. main verifies that Goxide is installed correctly.
package main

import (
    "fmt"
    "github.com/seyallius/goxide/rusty/result"
)

func main() {
    res := result.Ok(42)
    fmt.Println(res.Unwrap()) // 42
}
```

```bash
go run main.go
# Output: 42
```

## Upgrading

```bash
go get -u github.com/seyallius/goxide@latest
```

> **Note:** When v2 is released, the module path will become
> `github.com/seyallius/goxide/v2` per Go's semantic import versioning rules.
> All import paths will need the `/v2` suffix at that point.
