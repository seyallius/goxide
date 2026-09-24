<!-- cookbook.md provides standalone, runnable examples for quick reference. -->

# 📖 Quick Cookbook

Welcome to the Goxide Cookbook! These are standalone, runnable examples.
You can copy any of these blocks into a `main.go` file, run `go run main.go`, and see Goxide in action immediately.

## 1. Basic Error Handling with BubbleUp

The most common use case: replacing `if err != nil` with clean, linear code.

```go
// Package main. main demonstrates basic BubbleUp error handling.
package main

import (
    "errors"
    "fmt"
    "github.com/seyallius/goxide/rusty/result"
)

func divide(a, b int) result.Result[int] {
    if b == 0 {
        return result.Err[int](errors.New("division by zero"))
    }
    return result.Ok(a / b)
}

func complexMath(a, b, c int) (res result.Result[int]) {
    defer result.Catch(&res)

    // If divide(a, b) fails, it immediately returns the error.
    // No nested if-statements!
    step1 := divide(a, b).BubbleUp()
    step2 := divide(step1, c).BubbleUp()

    return result.Ok(step2 * 10)
}

func main() {
    // Success case
    res1 := complexMath(100, 2, 5)
    fmt.Printf("Success: %d\n", res1.Unwrap()) // Output: Success: 100

    // Error case
    res2 := complexMath(100, 0, 5)
    fmt.Printf("Error: %v\n", res2.Err())     // Output: Error: division by zero
}
```

## 2. Safe Configuration with Option

Use `Option` when a value might be missing, but it's not an "error".

```go
// Package main. main demonstrates Option for configuration defaults.
package main

import (
    "fmt"
    "github.com/seyallius/goxide/rusty/option"
)

type Config struct {
    Port    option.Option[int]
    Timeout option.Option[int]
}

func main() {
    // Simulate loading config where some values are missing
    cfg := Config{
        Port:    option.Some(8080),
        Timeout: option.None[int](), // Missing!
    }

    // UnwrapOr provides safe defaults without nil-checks
    port := cfg.Port.UnwrapOr(3000)
    timeout := cfg.Timeout.UnwrapOr(30)

    fmt.Printf("Server starting on port %d with %ds timeout\n", port, timeout)
    // Output: Server starting on port 8080 with 30s timeout
}
```

## 3. Data Transformation Pipelines with Chain

When you need to transform data through multiple steps, `chain` reads left-to-right.

```go
// Package main. main demonstrates fluent method chaining.
package main

import (
    "fmt"
    "strings"
    "github.com/seyallius/goxide/rusty/chain"
    "github.com/seyallius/goxide/rusty/result"
)

func validate(s string) result.Result[string] {
    if len(s) < 3 {
        return result.Err[string](fmt.Errorf("too short"))
    }
    return result.Ok(s)
}

func main() {
    input := "  hello world  "

    res := chain.Chain(result.Ok(input)).
        Map(strings.TrimSpace).
        Map(strings.ToUpper).
        AndThen(validate).
        Unwrap()

    fmt.Println(res.Unwrap()) // Output: HELLO WORLD

    // If validation fails, the chain short-circuits
    badInput := "hi"
    badRes := chain.Chain(result.Ok(badInput)).
        Map(strings.TrimSpace).
        AndThen(validate).
        Unwrap()

    fmt.Println(badRes.Err()) // Output: too short
}
```
