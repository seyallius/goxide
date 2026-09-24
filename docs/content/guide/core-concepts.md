<!-- core-concepts.md explains the philosophy behind Result and Option, -->
<!-- when to use each, and how they relate to traditional Go patterns. -->

# 🧠 Core Concepts: Result & Option

## The Problem with `(T, error)`

In traditional Go, a function returns `(T, error)`. Every caller must manually
inspect the error, branch, and return early:

```go
// Package main. traditional demonstrates standard Go error handling fatigue.
package main

// ProcessUserTraditional shows the verbose nature of standard error handling.
func ProcessUserTraditional(id int) (string, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return "", err // Manual short-circuit
    }

    profile, err := db.FindProfile(user.ID)
    if err != nil {
        return "", err // Manual short-circuit
    }

    return profile.Bio, nil
}
```

Two problems:

1. **The arrow of death** — each `if err != nil` adds a level of nesting.
2. **Zero-value traps** — nothing stops you from using `user` after an error.
   The compiler won't warn you.

## Railway-Oriented Programming

Goxide wraps operations in a `Result[T]`. Think of it like a train track:

```
     ┌─────────┐      ┌─────────┐     ┌─────────┐
───►│ FindUser ├────►│ Process ├────►│  Save   ├───► Ok(receipt)
     └─────────┘      └─────────┘     └─────────┘
         │               │               │
         ▼               ▼               ▼
    ┌─────────────────────────────────────────┐
    │              Error track                │───► Err(error)
    └─────────────────────────────────────────┘
```

If any operation fails, the train switches to the **Error track** and all
subsequent stations are skipped automatically. No manual `if err != nil`.

## Result vs Option

|                | `result.Result[T]`                          | `option.Option[T]`                                 |
| -------------- | ------------------------------------------- | -------------------------------------------------- |
| **Use when**   | Operation can **fail** and you care **why** | Value might be **absent**, but that's not an error |
| **Examples**   | DB timeout, invalid email, network error    | Cache miss, unset bio, optional config             |
| **Rust equiv** | `Result<T, E>`                              | `Option<T>`                                        |
| **States**     | `Ok(value)` or `Err(error)`                 | `Some(value)` or `None`                            |

### Rule of Thumb

> If you'd write `if err != nil`, use **Result**.
> If you'd write `if x == nil`, use **Option**.

## Converting Between Them

You can discard error details when you only care about presence:

```go
// Package main. conversion demonstrates moving from Result to Option.
package main

import (
    "github.com/seyallius/goxide/rusty/option"
    "github.com/seyallius/goxide/rusty/result"
)

// GetCachedUser treats a failed cache lookup as simply "None".
func GetCachedUser(id int) option.Option[User] {
    cacheResult := FetchFromCache(id)
    return cacheResult.Value() // Some(user) if Ok, None if Err
}
```

And you can promote an `Option` to a `Result` when absence _is_ an error:

```go
// GetRequiredUser converts None into a meaningful error.
func GetRequiredUser(id int) result.Result[User] {
    opt := cache.GetUser(id)
    return opt.ToResult(ErrUserNotFound)
}
```

## The Three Styles

Goxide doesn't force you into one paradigm. Pick per-function:

| Style           | When                                      | Example                                    |
| --------------- | ----------------------------------------- | ------------------------------------------ |
| **BubbleUp**    | Business logic with many sequential steps | `defer result.Catch(&res)` + `.BubbleUp()` |
| **Chain**       | Data transformation pipelines             | `chain.Chain(...).Map(...).AndThen(...)`   |
| **Traditional** | Performance-critical loops, simple cases  | `if err != nil { return }`                 |

You can mix all three in the same codebase. That's by design.
