<!-- http-handlers.md shows how to use Goxide in HTTP handler functions, -->
<!-- including JSON validation, error responses, and middleware patterns. -->

# 🌐 Example: HTTP Handlers

This example shows Goxide in a REST API handler with JSON parsing,
validation, and structured error responses.

## Handler with BubbleUp + CatchErr

```go
// Package main. main demonstrates Goxide in HTTP handlers.
package main

import (
    "encoding/json"
    "net/http"

    "github.com/seyallius/goxide/rusty/result"
)

// CreateUserRequest is the expected JSON body.
type CreateUserRequest struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

// CreateUserResponse is the JSON response.
type CreateUserResponse struct {
    ID    int    `json:"id"`
    Email string `json:"email"`
}

// CreateUserHandler handles POST /users.
// It uses BubbleUp internally but returns (response, error) for the router.
func CreateUserHandler(repo *UserRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var resp CreateUserResponse
        var err error
        defer result.CatchErr(&resp, &err)

        // Step 1: Parse JSON body
        var req CreateUserRequest
        decodeErr := json.NewDecoder(r.Body).Decode(&req)
        result.Wrap(req, decodeErr).BubbleUp()

        // Step 2: Validate
        result.Wrap(validateEmail(req.Email)).BubbleUp()

        // Step 3: Create user in database
        id := repo.Create(r.Context(), req.Email, req.Name).BubbleUp()

        // Step 4: Build response
        resp = CreateUserResponse{ID: id, Email: req.Email}

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        _ = json.NewEncoder(w).Encode(resp)
    }
}
```

## Validation with Result

```go
// Package main. main demonstrates validation functions returning Result.
package main

import (
    "errors"
    "net/mail"

    "github.com/seyallius/goxide/rusty/result"
)

var ErrInvalidEmail = errors.New("invalid email format")
var ErrNameTooShort = errors.New("name must be at least 2 characters")

// validateEmail returns Ok(email) or Err with a descriptive message.
func validateEmail(email string) result.Result[string] {
    if _, err := mail.ParseAddress(email); err != nil {
        return result.Err(ErrInvalidEmail)
    }
    return result.Ok(email)
}

// validateName checks minimum length.
func validateName(name string) result.Result[string] {
    if len(name) < 2 {
        return result.Err(ErrNameTooShort)
    }
    return result.Ok(name)
}

// validateRequest runs all validations, collecting all errors.
func validateRequest(req CreateUserRequest) (res result.Result[CreateUserRequest]) {
    defer result.Catch(&res)

    validateEmail(req.Email).BubbleUp()
    validateName(req.Name).BubbleUp()

    return result.Ok(req)
}
```

## Chaining in Handlers

```go
// ProcessUserPipeline chains parse → validate → enrich → respond.
func ProcessUserPipeline(body []byte) result.Result[CreateUserResponse] {
    return chain.Chain(parseJSON(body)).
        AndThen(validateRequest).
        Map(enrichWithDefaults).
        Map(toResponse).
        Unwrap()
}
```

## Error Response Middleware

```go
// ErrorResponse converts a Result error into an HTTP error response.
func ErrorResponse(w http.ResponseWriter, res result.Result[any]) {
    if res.IsOk() {
        return
    }

    err := res.Err()
    status := http.StatusInternalServerError

    switch {
    case errors.Is(err, ErrInvalidEmail), errors.Is(err, ErrNameTooShort):
        status = http.StatusBadRequest
    case errors.Is(err, ErrUserNotFound):
        status = http.StatusNotFound
    case errors.Is(err, ErrEmailExists):
        status = http.StatusConflict
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(map[string]string{
        "error": err.Error(),
    })
}
```
