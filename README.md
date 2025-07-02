# Composable Query Builder

A SQL query builder for Go.

## Goals
- **Lightweight**: Minimal abstractions, no heavy ORM features.
- **Thread Safe**: Safe for concurrent use.
- **Database Agnostic**: Works with any database that implements Go's `database/sql` interface.
- **Simple API**: Make common queries (SELECT, INSERT, UPDATE, DELETE) easy, but allow for custom/raw SQL.

## Setting the SQL Dialect
Set the dialect globally for your application:

```go
import "github.com/sprylic/cqb"

sq.SetDialect(sq.Postgres()) // or sq.MySQL(), sq.Standard()
```

## Example Usage

### SELECT
```go
q := sq.Select("id", "name").From("users").Where("active = ?", true)
sql, args, err := q.Build()
// sql: "SELECT id, name FROM users WHERE active = ?" (Standard dialect)
// args: [true]
```

### Aliasing and Subqueries
```go
sub := sq.Select("COUNT(*)").From("orders").Where("orders.user_id = users.id")
q := sq.Select(sq.Alias(sub, "order_count")).From("users")
sql, args, err := q.Build()
// sql: "SELECT (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) AS order_count FROM users"
```

### Query Composition
```go
isActive := func(b *sq.SelectBuilder) *sq.SelectBuilder {
    return b.Where("active = ?", true)
}
isAdult := func(b *sq.SelectBuilder) *sq.SelectBuilder {
    return b.Where("age >= ?", 18)
}
q := sq.Select("id").From("users").Compose(isActive, isAdult)
sql, args, err := q.Build()
// sql: "SELECT id FROM users WHERE active = ? AND age >= ?"
// args: [true, 18]
```

### INSERT
```go
q := sq.Insert("users").Columns("id", "name").Values(1, "Alice").Values(2, "Bob")
sql, args, err := q.Build()
// sql: "INSERT INTO users (id, name) VALUES (?, ?), (?, ?)"
// args: [1, "Alice", 2, "Bob"]
```

### UPDATE
```go
q := sq.Update("users").Set("name", "Alice").Where("id = ?", 1)
sql, args, err := q.Build()
// sql: "UPDATE users SET name = ? WHERE id = ?"
// args: ["Alice", 1]
```

### DELETE
```go
q := sq.Delete("users").Where("id = ?", 1)
sql, args, err := q.Build()
// sql: "DELETE FROM users WHERE id = ?"
// args: [1]
```

## SQL Dialect Examples

#### MySQL
```go
sq.SetDialect(sq.MySQL())
q := sq.Select("id", "name").From("users").Where("id = ?", 1)
sql, args, err := q.Build()
// sql: "SELECT `id`, `name` FROM `users` WHERE id = ?"
```

#### Postgres
```go
sq.SetDialect(sq.Postgres())
q := sq.Select("id", "name").From("users").Where("id = ? AND name = ?", 1, "bob")
sql, args, err := q.Build()
// sql: "SELECT \"id\", \"name\" FROM \"users\" WHERE id = $1 AND name = $2"
```

## Status
Work in progress. 