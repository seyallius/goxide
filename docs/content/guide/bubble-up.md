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
```

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
