<!-- option-patterns.md covers the Option[T] type: when to use it, -->
<!-- how it replaces nil pointers, and common access patterns. -->

# 🎯 Option Patterns

`Option[T]` represents a value that might not exist. It replaces the
Go idiom of returning `*T` (pointer) or `(T, bool)` (comma-ok).

## The Problem with Nil Pointers

```go
// Package main. nilpointer demonstrates the danger of *T for optional values.
package main

// GetBio returns nil if the user hasn't set a bio.
// The caller MUST remember to check for nil, or panic.
func GetBio(userID int) *string {
    user := db.FindUser(userID)
    if user.Bio == "" {
        return nil // ← caller might forget to check this
    }
    return &user.Bio
}

// Somewhere else...
bio := GetBio(42)
fmt.Println(*bio) // 💥 nil pointer dereference if bio not set
```

## The Option Solution

```go
// Package main. optsafe demonstrates safe optional value handling.
package main

import "github.com/seyallius/goxide/rusty/option"

// GetBio makes the optionality EXPLICIT in the signature.
func GetBio(userID int) option.Option[string] {
    user := db.FindUser(userID)
    if user.Bio == "" {
        return option.None[string]()
    }
    return option.Some(user.Bio)
}

// The caller is FORCED to handle both cases.
bio := GetBio(42)
fmt.Println(bio.UnwrapOr("No bio set")) // Safe. Always works.
```

## Common Access Patterns

### ✅ DO: Use `UnwrapOr` for Safe Defaults

```go
name := GetUserName(id).UnwrapOr("Guest")
timeout := GetConfig("timeout").UnwrapOr(30 * time.Second)
```

### ✅ DO: Use `Map` for Transformations

```go
// Transform the value inside the Option without unwrapping.
upperName := GetUserName(id).Map(strings.ToUpper)
// Some("ALICE") or None — no manual checking
```

### ✅ DO: Use `FlatMap` for Optional Operations

```go
// FlatMap chains operations that themselves return Option.
bio := GetUser(id).FlatMap(func(u User) option.Option[string] {
    return GetBio(u.ID)
})
```

### ✅ DO: Use `Some()` for Go-Idiomatic Access

```go
if val, ok := opt.Some(); ok {
    fmt.Println("Got:", val)
} else {
    fmt.Println("Nothing here")
}
```

### ❌ DON'T: Use `Unwrap()` in Production Code

```go
// ❌ Panics if None — only use in tests or when logically impossible
val := opt.Unwrap()
```

### ❌ DON'T: Use Option for Error Handling

```go
// ❌ Wrong — you lose the error information
func FindUser(id int) option.Option[User] {
    user, err := db.Query(id)
    if err != nil {
        return option.None[User]() // Was it a timeout? Bad query? Gone.
    }
    return option.Some(user)
}

// ✅ Right — use Result when failure has a reason
func FindUser(id int) result.Result[User] {
    return result.Wrap(db.Query(id))
}
```

## Migration from Pointers

**Before:**

```go
type Config struct {
    Timeout *int    // nil means "not set"
    Retries *int
}
```

**After:**

```go
type Config struct {
    Timeout option.Option[int]
    Retries option.Option[int]
}

// Usage
timeout := cfg.Timeout.UnwrapOr(30)
```

## Integration with Result

Convert freely between the two:

```go
// Result → Option (discard error details)
opt := result.Value()  // Some(v) if Ok, None if Err

// Option → Result (attach an error to None)
res := opt.ToResult(ErrNotFound)  // Ok(v) if Some, Err(ErrNotFound) if None
```
