# 🫧 The BubbleUp Pattern

The `BubbleUp()` method is Goxide's answer to Rust's `?` operator. It allows you to write linear, happy-path code while safely short-circuiting on errors.

## How it works

Under the hood, `BubbleUp()` triggers a controlled `panic` if the `Result` is an error. The `result.Catch()` deferred function instantly recovers this panic and populates your named return variable.

## Rules of the Road

1. **Always defer `Catch()` first**: It must be the very first line in your function.
2. **Use named returns**: Your function signature must use `(res result.Result[T])`.

```go
func GetUser(id int) (res result.Result[User]) {
    // 1. Catch MUST be deferred first
    defer result.Catch(&res)

    // 2. BubbleUp handles the early return automatically
    user := db.FindUser(id).BubbleUp()

    return result.Ok(user)
}
```
