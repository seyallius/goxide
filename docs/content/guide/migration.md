<!-- migration.md is the step-by-step guide for adopting Goxide in an -->
<!-- existing codebase without a big-bang rewrite. -->

# 🔄 Migration Guide

You don't need to rewrite your codebase. Goxide is designed for **gradual
adoption** — you can migrate one function at a time.

## Strategy: Wrap, Don't Rewrite

The key insight: `result.Wrap()` bridges traditional `(T, error)` functions
into `Result[T]` without changing the original function.

```go
// Package main. main demonstrates the Wrap bridge pattern.
package main

import "github.com/seyallius/goxide/rusty/result"

// These are existing library functions you CANNOT change:
//   func db.FindUser(id int) (User, error)
//   func db.FindProfile(uid int) (Profile, error)

// Step 1: Create Result-wrapped versions (one line each)
var findUser    = result.WrapFunc1(db.FindUser)
var findProfile = result.WrapFunc1(db.FindProfile)

// Step 2: Use them with BubbleUp in new code
func GetUserBio(id int) (res result.Result[string]) {
    defer result.Catch(&res)

    user    := findUser(id).BubbleUp()
    profile := findProfile(user.ID).BubbleUp()

    return result.Ok(profile.Bio)
}

// Step 3: Old callers still work — the original functions are untouched
func LegacyHandler(id int) (string, error) {
    user, err := db.FindUser(id) // ← still works, unchanged
    if err != nil {
        return "", err
    }
    return user.Name, nil
}
```

## Phase 1: New Code Uses Goxide

Write all **new** functions with `Result[T]` and `BubbleUp()`. Don't touch
existing code yet.

```go
// Package main. main demonstrates new code using Goxide patterns.
package main

import "github.com/seyallius/goxide/rusty/result"

// NewFeature is brand-new code — use Goxide from the start.
func NewFeature(input Input) (res result.Result[Output]) {
    defer result.Catch(&res)

    validated := Validate(input).BubbleUp()
    processed := Process(validated).BubbleUp()

    return result.Ok(processed)
}
```

## Phase 2: Boundary Adaptation with CatchErr

When new Goxide code needs to satisfy old interfaces:

```go
// Package main. main demonstrates CatchErr at the boundary.
package main

import "github.com/seyallius/goxide/rusty/result"

// OldInterface expects traditional (T, error) returns.
type OldInterface interface {
    DoWork(input string) (Result, error)
}

// NewImpl uses Goxide internally but complies with OldInterface externally.
type NewImpl struct{}

func (n *NewImpl) DoWork(input string) (out Result, err error) {
    defer result.CatchErr(&out, &err)

    parsed := Parse(input).BubbleUp()
    out = Compute(parsed).BubbleUp()

    return out, nil
}
```

## Phase 3: Migrate Existing Functions (Optional)

When you're ready, convert old functions one at a time:

**Before:**

```go
func GetUserData(id int) (UserData, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return UserData{}, err
    }

    profile, err := db.FindProfile(user.ID)
    if err != nil {
        return UserData{}, err
    }

    return ProcessData(user, profile), nil
}
```

**After:**

```go
func GetUserData(id int) (res result.Result[UserData]) {
    defer result.Catch(&res)

    user    := db.FindUser(id).BubbleUp()    // Wrap + BubbleUp
    profile := db.FindProfile(user.ID).BubbleUp()

    return result.Ok(ProcessData(user, profile))
}
```

> **Tip:** You can keep the old function signature and add a new one alongside
> it during the transition period. Deprecate the old one with a comment.

## Phase 4: Introduce Chaining for Pipelines

Once individual functions return `Result[T]`, you can compose them with chains:

```go
// Package main. main demonstrates composing migrated functions with chain.
package main

import (
    "github.com/seyallius/goxide/rusty/chain"
    "github.com/seyallius/goxide/rusty/result"
)

// FullPipeline composes three independently-migrated functions.
func FullPipeline(userID int) result.Result[Report] {
    return chain.Chain(GetUserData(userID)).
        Map(func(ud UserData) Summary { return Summarize(ud) }).
        AndThen(GenerateReport).
        Unwrap()
}
```

## Checklist

- [ ] Add `go get github.com/seyallius/goxide` to your project
- [ ] Write new functions with `Result[T]` + `BubbleUp()`
- [ ] Use `result.Wrap()` / `result.WrapFunc1()` to bridge existing functions
- [ ] Use `CatchErr` at interface boundaries
- [ ] Migrate existing functions opportunistically (not big-bang)
- [ ] Introduce `chain` for multi-step pipelines once functions return `Result[T]`
- [ ] Add `option.Option[T]` for nullable fields in new structs
