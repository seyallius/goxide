<!-- home.md is the landing page of the documentation site. It gives a high-level -->
<!-- overview of what Goxide is, why it exists, and where to go next. -->

# 🦀🐹 Welcome to Goxide

**Goxide** brings the elegance of Rust's `Result<T, E>` and `Option<T>` types to Go,
eliminating the verbosity of traditional `if err != nil` checks while maintaining
Go's explicit nature.

## Why Goxide?

Traditional Go error handling is explicit, but it often leads to the "arrow of death"
(deeply nested code) and makes it easy to accidentally use a zero-value when an error
actually occurred.

Goxide leverages Go generics to provide:

1. **Compiler-enforced safety** — you cannot access a value without acknowledging
   the error state.
2. **Composability** — chain operations together; errors automatically short-circuit
   the pipeline.
3. **Ergonomics** — `BubbleUp()` mimics Rust's `?` operator using Go's
   `panic`/`recover` mechanism safely, keeping business logic linear and readable.

## The Four Packages

| Package                                                                 | Purpose             | Rust Equivalent           |
| ----------------------------------------------------------------------- | ------------------- | ------------------------- |
| [`result`](https://pkg.go.dev/github.com/seyallius/goxide/rusty/result) | Fallible operations | `Result<T, E>`            |
| [`option`](https://pkg.go.dev/github.com/seyallius/goxide/rusty/option) | Optional values     | `Option<T>`               |
| [`chain`](https://pkg.go.dev/github.com/seyallius/goxide/rusty/chain)   | Fluent pipelines    | Method chaining           |
| [`types`](https://pkg.go.dev/github.com/seyallius/goxide/rusty/types)   | Functional helpers  | `Fn` traits, `identity()` |

## Quick Taste

```go
// Package main. main demonstrates the BubbleUp pattern for clean error handling.
package main

import "github.com/seyallius/goxide/rusty/result"

// ProcessOrder fetches an order, charges payment, and builds a receipt.
// Zero `if err != nil` checks. Errors short-circuit automatically.
func ProcessOrder(orderID int) (res result.Result[string]) {
    defer result.Catch(&res)

    order   := FindOrder(orderID).BubbleUp()
    payment := ProcessPayment(order).BubbleUp()
    receipt := GenerateReceipt(payment).BubbleUp()

    return result.Ok(receipt)
}
```

## Where to Go Next

- **[Installation](#installation)** — add Goxide to your project
- **[Result & Option](#core-concepts)** — understand the two core types
- **[The BubbleUp Pattern](#bubble-up)** — the headline feature
- **[Migration Guide](#migration)** — adopt Goxide gradually in existing code
