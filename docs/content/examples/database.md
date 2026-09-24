<!-- database.md is a worked example showing Goxide in a database -->
<!-- repository layer with transactions, fallbacks, and error wrapping. -->

# 🗄️ Example: Database Operations

This example shows a user repository that uses Goxide for clean error handling,
transaction management, and fallback strategies.

## Domain Types

```go
// Package main. main demonstrates Goxide in a database repository layer.
package main

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "github.com/seyallius/goxide/rusty/result"
)

// User represents a database entity.
type User struct {
    ID        int
    Email     string
    Name      string
    CreatedAt time.Time
}

// Sentinel errors for the repository layer.
var (
    ErrUserNotFound  = errors.New("user not found")
    ErrEmailExists   = errors.New("email already exists")
    ErrDBConnection  = errors.New("database connection failed")
)
```

## Repository with Result[T]

```go
// UserRepository wraps database operations in Result[T].
type UserRepository struct {
    db *sql.DB
}

// NewUserRepository creates a repository with the given database connection.
func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

// FindByID retrieves a user by primary key.
func (r *UserRepository) FindByID(ctx context.Context, id int) result.Result[User] {
    var u User
    err := r.db.QueryRowContext(ctx,
        "SELECT id, email, name, created_at FROM users WHERE id = $1", id,
    ).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt)

    if errors.Is(err, sql.ErrNoRows) {
        return result.Err(ErrUserNotFound)
    }
    return result.Wrap(u, err)
}

// Create inserts a new user and returns the generated ID.
func (r *UserRepository) Create(ctx context.Context, email, name string) result.Result[int] {
    var id int
    err := r.db.QueryRowContext(ctx,
        "INSERT INTO users (email, name, created_at) VALUES ($1, $2, $3) RETURNING id",
        email, name, time.Now(),
    ).Scan(&id)

    return result.Wrap(id, err)
}
```

## Transaction Pattern

```go
// CreateUserWithProfile creates a user and profile atomically.
// If either step fails, the transaction is rolled back automatically.
func (r *UserRepository) CreateUserWithProfile(
    ctx context.Context, email, name, bio string,
) (res result.Result[User]) {
    defer result.Catch(&res)

    // Begin transaction
    tx := result.Wrap(r.db.BeginTx(ctx, nil)).BubbleUp()

    // Ensure rollback on any error via a deferred check
    defer func() {
        if res.IsErr() {
            _ = tx.Rollback()
        }
    }()

    // Step 1: Insert user
    var userID int
    err := tx.QueryRowContext(ctx,
        "INSERT INTO users (email, name, created_at) VALUES ($1, $2, $3) RETURNING id",
        email, name, time.Now(),
    ).Scan(&userID)
    result.Wrap(userID, err).BubbleUp()

    // Step 2: Insert profile
    _, err = tx.ExecContext(ctx,
        "INSERT INTO profiles (user_id, bio) VALUES ($1, $2)",
        userID, bio,
    )
    result.Wrap(struct{}{}, err).BubbleUp()

    // Commit
    result.Wrap(struct{}{}, tx.Commit()).BubbleUp()

    return result.Ok(User{ID: userID, Email: email, Name: name})
}
```

## Fallback: Cache → Database

```go
// GetUser tries the cache first, falls back to the database.
func (s *UserService) GetUser(ctx context.Context, id int) (res result.Result[User]) {
    defer result.Catch(&res)

    // Try cache — treat miss as non-fatal
    defer result.CatchWith(&res,
        func(err error) bool { return errors.Is(err, ErrCacheMiss) },
        func(err error) {
            // Cache miss → fall through to database
            res = s.repo.FindByID(ctx, id)
        },
    )

    return s.cache.Get(ctx, id)
}
```

## Traditional Interop

```go
// LegacyHandler implements an old interface that expects (User, error).
func (s *UserService) LegacyHandler(id int) (User, error) {
    var user User
    var err error
    defer result.CatchErr(&user, &err)

    user = s.GetUser(context.Background(), id).BubbleUp()
    return user, nil
}
```
