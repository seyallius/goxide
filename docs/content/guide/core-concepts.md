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

````

### `docs/content/guide/bubble-up.md`

```markdown
<!-- File: docs/content/guide/bubble-up.md -->
<!-- bubble-up.md is the deep-dive guide for the BubbleUp/Catch pattern, -->
<!-- which is Goxide's equivalent of Rust's ? operator. -->

# 🫧 The BubbleUp Pattern

`BubbleUp()` is Goxide's answer to Rust's `?` operator. It lets you write
linear, happy-path code while safely short-circuiting on errors.

## Before and After

**Traditional Go (5 operations = 5 error checks):**

```go
func ProcessOrder(orderID int) (Receipt, error) {
    order, err := FindOrder(orderID)
    if err != nil { return Receipt{}, err }

    customer, err := FindCustomer(order.CustomerID)
    if err != nil { return Receipt{}, err }

    payment, err := ChargePayment(customer, order.Total)
    if err != nil { return Receipt{}, err }

    receipt, err := GenerateReceipt(payment)
    if err != nil { return Receipt{}, err }

    saved, err := SaveReceipt(receipt)
    if err != nil { return Receipt{}, err }

    return saved, nil
}
````

**With BubbleUp (zero error checks):**

```go
// Package main. main demonstrates the BubbleUp pattern for clean error handling.
package main

import "github.com/seyallius/goxide/rusty/result"

// ProcessOrder fetches, charges, receipts, and saves — no if-err noise.
func ProcessOrder(orderID int) (res result.Result[Receipt]) {
    defer result.Catch(&res) // MUST be first line

    order    := FindOrder(orderID).BubbleUp()
    customer := FindCustomer(order.CustomerID).BubbleUp()
    payment  := ChargePayment(customer, order.Total).BubbleUp()
    receipt  := GenerateReceipt(payment).BubbleUp()
    saved    := SaveReceipt(receipt).BubbleUp()

    return result.Ok(saved)
}
```

## How It Works Under the Hood

```
BubbleUp() called on Ok(value)  →  returns value, execution continues
BubbleUp() called on Err(err)   →  panics with a wrapped error
                                     ↓
                              defer Catch(&res) recovers the panic
                                     ↓
                              res is populated with Err(err)
                                     ↓
                              function returns normally to caller
```

The `panic`/`recover` mechanism is **contained entirely within your function**.
The caller sees a normal `Result` return — no panics leak.

## The Three Rules

### Rule 1: Always `defer Catch()` First

```go
func DoWork() (res result.Result[int]) {
    defer result.Catch(&res) // ← MUST be the very first statement
    // ...
}
```

If you forget this, a `BubbleUp()` on an error will panic _uncaught_ and crash
your program. There is no compile-time check for this — it's a convention you
must follow.

### Rule 2: Use Named Return Values

```go
// ✅ Correct — named return
func DoWork() (res result.Result[int]) { ... }

// ❌ Wrong — Catch has nothing to write into
func DoWork() result.Result[int] { ... }
```

`Catch` needs a pointer to your return variable so it can populate it during
recovery.

### Rule 3: Don't Mix BubbleUp and Manual Checks in the Same Function

```go
// ❌ Confusing — pick one style per function
func DoWork() (res result.Result[int]) {
    defer result.Catch(&res)

    a := Step1().BubbleUp()

    b, err := Step2()       // traditional style mixed in
    if err != nil {
        return result.Err(err)
    }

    return result.Ok(a + b)
}
```

If you need traditional `(T, error)` interop, use `result.Wrap()`:

```go
// ✅ Clean — Wrap bridges the gap
func DoWork() (res result.Result[int]) {
    defer result.Catch(&res)

    a := Step1().BubbleUp()
    b := result.Wrap(Step2()).BubbleUp() // (T, error) → Result[T]

    return result.Ok(a + b)
}
```

## Performance

From the project benchmarks (i5-11400H):

```
Traditional error handling:    ~0.25 ns/op    0 allocs
BubbleUp (happy path):       ~0.30 ns/op    0 allocs
BubbleUp (error path):      ~150 ns/op     2 allocs  (panic/recover cost)
```

The happy path is nearly free. The error path pays for `panic`/`recover`, but
errors should be _exceptional_ — if you're erroring on every call, that's a
design problem, not a Goxide problem.

## When NOT to Use BubbleUp

- **Hot loops** where errors are expected (e.g., scanning for a delimiter)
- **Single-operation functions** where `if err != nil` is already one line
- **Library boundaries** where callers expect `(T, error)` — use `CatchErr` instead:

```go
// CatchErr adapts BubbleUp to traditional signatures for interface compliance.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (User, error) {
    var user User
    var err error
    defer result.CatchErr(&user, &err)

    user = ExtractUser(r).BubbleUp()
    return user, nil
}
```
