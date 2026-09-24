<!-- recovery.md covers CatchWith, Fallback, and CatchErr — the three -->
<!-- mechanisms for recovering from errors instead of just propagating them. -->

# 🛟 Error Recovery & Fallbacks

Sometimes you don't want to propagate an error — you want to **recover**.
Goxide provides three mechanisms for this, each with different granularity.

## The Three Recovery Tools

| Tool                                | Granularity                      | Use When             |
| ----------------------------------- | -------------------------------- | -------------------- |
| `Fallback(&res, defaultVal)`        | Any error → default value        | Simple defaults      |
| `CatchWith(&res, matcher, handler)` | Specific error → custom recovery | Error-specific logic |
| `CatchErr(&val, &err)`              | Adapt to `(T, error)` signatures | Interface compliance |

## 1. Fallback — Simple Default Values

`Fallback` provides a default value for **any** error. It's the simplest
recovery mechanism.

```go
// Package main. main demonstrates Fallback for default configuration values.
package main

import "github.com/seyallius/goxide/rusty/result"

// GetTimeout returns the configured timeout or a safe default.
func GetTimeout() (res result.Result[int]) {
    defer result.Catch(&res)
    defer result.Fallback(&res, 30) // Default to 30s for ANY error

    return LoadTimeoutConfig()
}
```

**Order matters:** `Fallback` must be deferred **after** `Catch`:

```go
defer result.Catch(&res)     // ← recovers the BubbleUp panic first
defer result.Fallback(&res, 30) // ← then applies the default
```

Defers execute in LIFO order, so `Fallback` runs first, then `Catch`.
If `Fallback` ran after `Catch`, the error would already be set and
`Fallback` would overwrite your specific error with a generic default.

## 2. CatchWith — Error-Specific Recovery

`CatchWith` lets you match on specific errors and recover differently:

```go
// Package main. main demonstrates CatchWith for error-specific recovery.
package main

import (
    "errors"
    "github.com/seyallius/goxide/rusty/result"
)

var ErrRateLimited = errors.New("rate limited")
var ErrNotFound    = errors.New("not found")

// FetchWithRetry handles rate limits by retrying, and not-found by defaulting.
func FetchWithRetry(url string) (res result.Result[Response]) {
    defer result.Catch(&res)

    defer result.CatchWith(&res,
        func(err error) bool { return errors.Is(err, ErrRateLimited) },
        func(err error) {
            time.Sleep(time.Second)
            res = Fetch(url) // retry once
        },
    )

    defer result.CatchWith(&res,
        func(err error) bool { return errors.Is(err, ErrNotFound) },
        func(err error) {
            res = result.Ok(Response{Status: 404, Body: "default page"})
        },
    )

    return Fetch(url)
}
```

## 3. CatchErr — Adapt to Traditional Signatures

When implementing interfaces that expect `(T, error)` returns (like
`http.Handler`, `io.Reader`, etc.), use `CatchErr`:

```go
// Package main. main demonstrates CatchErr for interface compliance.
package main

import (
    "net/http"
    "github.com/seyallius/goxide/rusty/result"
)

// HandleRequest uses BubbleUp internally but returns (User, error) externally.
func HandleRequest(w http.ResponseWriter, r *http.Request) (user User, err error) {
    defer result.CatchErr(&user, &err)

    userID := ExtractUserID(r).BubbleUp()
    user = GetUser(userID).BubbleUp()

    return user, nil
}
```

## Multi-Layer Fallback

You can stack multiple recovery layers:

```go
// Package main. main demonstrates layered fallback strategy.
package main

import "github.com/seyallius/goxide/rusty/result"

// GetFeatureFlag tries: DB → config file → hardcoded default.
func GetFeatureFlag(name string) (res result.Result[bool]) {
    defer result.Catch(&res)

    // Layer 3: ultimate fallback
    defer result.Fallback(&res, false)

    // Layer 2: try config file
    defer result.CatchWith(&res,
        func(err error) bool { return true }, // catch any remaining error
        func(err error) {
            res = LoadFlagFromFile(name)
        },
    )

    // Layer 1: try database
    return LoadFlagFromDB(name)
}
```

## Decision Flowchart

```
Error occurred
    │
    ├─ Do I need the error details? ──No──► Fallback(&res, default)
    │
    ├─ Do I need DIFFERENT recovery per error type?
    │       │
    │      Yes ──► CatchWith(&res, matcher, handler)
    │
    └─ Am I implementing a (T, error) interface?
            │
           Yes ──► CatchErr(&val, &err)
```
