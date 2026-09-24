<!-- chaining.md explains the chain package's fluent pipeline API, -->
<!-- how it compares to nested Result calls, and when to prefer it. -->

# 🔗 Fluent Chaining

The `chain` package lets you build readable, left-to-right data pipelines
instead of deeply nested `Map`/`AndThen` calls.

## The Problem: Nested Style

```go
// Package main. nested demonstrates the hard-to-read nested approach.
package main

import "github.com/seyallius/goxide/rusty/result"

// GetUserBioNested is functionally correct but visually painful.
func GetUserBioNested(userID int) result.Result[string] {
    return result.AndThen(
        result.Map(
            findUser(userID),
            func(u User) int { return u.ProfileID },
        ),
        func(profileID int) result.Result[Profile] {
            return findProfile(profileID)
        },
    )
}
```

You read it **inside-out**. The first operation is buried in the middle.

## The Solution: Chain Style

```go
// Package main. chained demonstrates the fluent chain approach.
package main

import (
    "github.com/seyallius/goxide/rusty/chain"
    "github.com/seyallius/goxide/rusty/result"
)

// GetUserBio reads top-to-bottom like a recipe.
func GetUserBio(userID int) result.Result[string] {
    return chain.Chain(findUser(userID)).
        Map(func(u User) int { return u.ProfileID }).
        AndThen(findProfile).
        Map(func(p Profile) string { return p.Bio }).
        Unwrap()
}
```

## Core Operations

| Method         | Input → Output  | Use When                           |
| -------------- | --------------- | ---------------------------------- |
| `.Map(f)`      | `T → U`         | Transform the value (infallible)   |
| `.AndThen(f)`  | `T → Result[U]` | Next step can also fail            |
| `.MapError(f)` | `error → error` | Wrap/annotate errors               |
| `.Unwrap()`    | → `Result[T]`   | End the chain, get the Result back |

## Error Short-Circuiting

If **any** step returns `Err`, all subsequent steps are skipped:

```go
// If findUser fails, Map and AndThen never execute.
res := chain.Chain(findUser(999)).   // ← Err("user not found")
    Map(func(u User) string {        // ← SKIPPED
        return u.Name
    }).
    AndThen(validateName).           // ← SKIPPED
    Unwrap()                         // ← Contains original error

res.IsErr() // true
res.Err()   // "user not found"
```

## Chain vs Chain2

Use `Chain2` when you know the exact number of transformations at compile time.
It gives the compiler more type information:

```go
// Chain2[string, User, int] means:
//   start with Result[int], transform to User, then to string
result := chain.Chain2[string, User, int](findUser(123)).
    Map(func(u User) string { return u.Name }).
    AndThen(validateName)
```

## Traditional vs Chained: Side by Side

| Aspect              | Traditional Nesting           | Fluent Chaining            |
| ------------------- | ----------------------------- | -------------------------- |
| **Readability**     | Deeply nested, hard to follow | Linear, easy to read       |
| **Maintainability** | Difficult to modify           | Easy to add/remove steps   |
| **Debugging**       | Hard to trace through nesting | Clear step-by-step flow    |
| **Type Safety**     | Manual type tracking          | Compiler-enforced          |
| **Error Handling**  | Manual error propagation      | Automatic short-circuiting |

## When to Use Chain vs BubbleUp

| Use **Chain** when...                          | Use **BubbleUp** when...                 |
| ---------------------------------------------- | ---------------------------------------- |
| You have a data transformation pipeline        | You have sequential side-effecting steps |
| Operations are pure functions                  | Operations involve I/O, mutations        |
| You want to see the whole pipeline at a glance | You want imperative, top-to-bottom flow  |
| You're composing reusable transformations      | You're writing business logic            |

Both are valid. Mix them freely across your codebase.
