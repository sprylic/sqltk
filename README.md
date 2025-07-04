# SQL Tool Kit

A SQL toolkit for Go that provides composable query building and DDL operations.

[![CI](https://github.com/sprylic/stk/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/sprylic/stk/actions/workflows/ci.yml)

## Goals
- **Thread Safe**: Safe for concurrent use.
- **Database Agnostic**: Works with any database that implements Go's `database/sql` interface.
- **Simple API**: Make common queries (SELECT, INSERT, UPDATE, DELETE) and DDL operations easy, but allow for custom/raw SQL.

## Default SQL Dialect
**MySQL is the default dialect.**
- Identifiers are quoted with backticks (`` `foo` ``) and placeholders are `?`.
- If you use Postgres or another database, **set the dialect explicitly**:
  ```go
  stk.SetDialect(stk.Postgres()) // for Postgres
  stk.SetDialect(stk.Standard()) // for no quoting (legacy/ANSI)
  
  q := stk.Select("id", "name").From("users").Where("active = ?", true) // set dialect per builder.
  sql, args, err := q.WithDialect(stk.PostGres()).Build()
  ```

### ⚠️ Setting the dialect globally can cause issues when using different dialects concurrently. If you need to support a different dialect, use WithDialect on the builder instead.

## Setting the SQL Dialect
Set the dialect globally for your application:

```go
import "github.com/sprylic/stk"

stk.SetDialect(stk.Postgres()) // or stk.MySQL(), stk.Standard()
```

## Example Usage

### SELECT
```go
q := stk.Select("id", "name").From("users").Where("active = ?", true)
sql, args, err := q.Build()
// sql: "SELECT `id`, `name` FROM `users` WHERE active = ?" (MySQL dialect by default)
// args: [true]
```

### Aliasing and Subqueries
```go
sub := stk.Select("COUNT(*)").From("orders").Where("orders.user_id = users.id")
q := stk.Select(stk.Alias(sub, "order_count")).From("users")
sql, args, err := q.Build()
// sql: "SELECT (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) AS order_count FROM `users`"
```

### Query Composition
```go
isActive := func(b *stk.SelectBuilder) *stk.SelectBuilder {
    return b.Where("active = ?", true)
}
isAdult := func(b *stk.SelectBuilder) *stk.SelectBuilder {
    return b.Where("age >= ?", 18)
}
q := stk.Select("id").From("users").Compose(isActive, isAdult)
sql, args, err := q.Build()
// sql: "SELECT `id` FROM `users` WHERE active = ? AND age >= ?"
// args: [true, 18]
```

### Condition Builder (NewCond)

The `ConditionBuilder` provides a fluent, composable API for building complex SQL conditions without resorting to raw SQL. Use `NewCond()` to start a condition chain, and pass it to `.Where()` or `.Having()` in any builder (`Select`, `Update`, `Delete`).

**Type-Safe Conditions:**

The library now provides a `Condition` interface for type-safe condition handling. This prevents unsafe queries and makes the API more explicit. The `Where()` and `Having()` methods now accept only `Condition` types:

```go
// Type-safe string condition
cond := stk.NewStringCondition("active = ? AND age > ?", true, 18)
q := stk.Select("id").From("users").Where(cond)

// Type-safe raw condition (requires AsCondition wrapper)
cond := stk.AsCondition(stk.Raw("id = 1"))
q := stk.Select("id").From("users").Where(cond)

// ConditionBuilder (implements Condition interface)
cond := stk.NewCond().Equal("active", true).And(stk.NewCond().GreaterThan("age", 18))
q := stk.Select("id").From("users").Where(cond)

// Compile-time type safety - these will not compile:
// q.Where("active = ?", true)           // Error: string doesn't implement Condition
// q.Where(123)                          // Error: int doesn't implement Condition
// q.Where(stk.Raw("id = 1"))            // Error: Raw doesn't implement Condition
```

**Interface Design:**

The `Condition` interface provides a clean, type-safe way to handle SQL conditions:

```go
type Condition interface {
    BuildCondition() (string, []interface{}, error)
}
```

All condition types implement this interface:
- `*StringCondition` - for parameterized string conditions
- `*RawCondition` - for raw SQL conditions  
- `*ConditionBuilder` - for fluent condition building
- `AsCondition(Raw)` - helper to convert Raw to Condition

**Examples:**

```go
// SELECT with complex condition
cond := stk.NewCond().
    Equal("active", true).
    And(stk.NewCond().GreaterThan("age", 18)).
    And(stk.NewCond().In("status", "active", "pending"))
q := stk.Select("id", "name").From("users").Where(cond)
sql, args, err := q.Build()
// sql: "SELECT `id`, `name` FROM `users` WHERE active = ? AND age > ? AND status IN (?, ?)"
// args: [true, 18, "active", "pending"]

// UPDATE with condition builder
cond := stk.NewCond().
    Equal("active", true).
    Or(stk.NewCond().Equal("vip", true)).
    And(stk.NewCond().GreaterThan("age", 16))
q := stk.Update("users").Set("name", "Alice").Where(cond)
sql, args, err := q.Build()
// sql: "UPDATE `users` SET `name` = ? WHERE (active = ?) OR (vip = ?) AND age > ?"
// args: ["Alice", true, true, 16]

// DELETE with condition builder
cond := stk.NewCond().
    Equal("active", false).
    Or(stk.NewCond().IsNull("deleted_at"))
q := stk.Delete("users").Where(cond)
sql, args, err := q.Build()
// sql: "DELETE FROM `users` WHERE (active = ?) OR (deleted_at IS NULL)"
// args: [false]
```

**Supported condition methods:**
- `Equal`, `NotEqual`, `GreaterThan`, `LessThan`, `GreaterThanOrEqual`, `LessThanOrEqual`
- `In`, `NotIn`, `Between`, `NotBetween`, `IsNull`, `IsNotNull`, `Like`, `NotLike`
- `Exists`, `NotExists`, `Case`, `And`, `Or`
- All methods are chainable and support table-qualified columns and dialect quoting.

**Type Safety:**
String conditions must now use explicit wrappers to prevent SQL injection:
```go
// ❌ This will now cause an error
q := stk.Select("id").From("users").Where("active = ?", true)

// ✅ Use explicit wrapper for string conditions
q := stk.Select("id").From("users").Where(stk.NewStringCondition("active = ?", true))

// ✅ Raw conditions still work (explicitly marked as raw)
q := stk.Select("id").From("users").Where(stk.Raw("id = 1"))
```

### INSERT
```go
q := stk.Insert("users").Columns("id", "name").Values(1, "Alice").Values(2, "Bob")
sql, args, err := q.Build()
// sql: "INSERT INTO `users` (`id`, `name`) VALUES (?, ?), (?, ?)"
// args: [1, "Alice", 2, "Bob"]
```

### UPDATE
```go
q := stk.Update("users").Set("name", "Alice").Where("id = ?", 1)
sql, args, err := q.Build()
// sql: "UPDATE `users` SET `name` = ? WHERE id = ?"
// args: ["Alice", 1]
```

### DELETE
```go
q := stk.Delete("users").Where("id = ?", 1)
sql, args, err := q.Build()
// sql: "DELETE FROM `users` WHERE id = ?"
// args: [1]
```

## DDL Operations

### Database Operations
```go
import "github.com/sprylic/stk/ddl"

// Create database
createDB := ddl.CreateDatabase("myapp_db").
    IfNotExists().
    Charset("utf8mb4").
    Collation("utf8mb4_unicode_ci")
sql, _, err := createDB.Build()
// sql: "CREATE DATABASE IF NOT EXISTS `myapp_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"

// Drop database
dropDB := ddl.DropDatabase("myapp_db").IfExists()
sql, _, err := dropDB.Build()
// sql: "DROP DATABASE IF EXISTS `myapp_db`"
```

### Schema Operations (PostgreSQL)
```go
// Create schema
createSchema := ddl.CreateSchema("myapp_schema").
    IfNotExists().
    Authorization("myapp_user")
sql, _, err := createSchema.Build()
// sql: "CREATE SCHEMA IF NOT EXISTS \"myapp_schema\" AUTHORIZATION \"myapp_user\""

// Drop schema
dropSchema := ddl.DropSchema("myapp_schema").
    IfExists().
    Cascade()
sql, _, err := dropSchema.Build()
// sql: "DROP SCHEMA IF EXISTS \"myapp_schema\" CASCADE"
```

### Table Operations
```go
// Create table
createTable := ddl.CreateTable("users").
    AddColumn(ddl.Column("id").Type("INT").AutoIncrement().NotNull()).
    AddColumn(ddl.Column("name").Type("VARCHAR").Size(255).NotNull()).
    AddColumn(ddl.Column("email").Type("VARCHAR").Size(255)).
    PrimaryKey("id").
    Unique("idx_email", "email")
sql, _, err := createTable.Build()
// sql: "CREATE TABLE `users` (`id` INT AUTO_INCREMENT NOT NULL, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255), PRIMARY KEY (`id`), CONSTRAINT idx_email UNIQUE (`email`))"

// Alter table
alterTable := ddl.AlterTable("users").
    AddColumn(ddl.Column("age").Type("INT")).
    AddIndex("idx_age", "age")
sql, _, err := alterTable.Build()
// sql: "ALTER TABLE `users` ADD COLUMN `age` INT, ADD INDEX idx_age (`age`)"

// Drop table
dropTable := ddl.DropTable("users").IfExists()
sql, _, err := dropTable.Build()
// sql: "DROP TABLE IF EXISTS `users`"

// Truncate table
truncateTable := ddl.TruncateTable("users")
sql, _, err := truncateTable.Build()
// sql: "TRUNCATE TABLE `users`"

// Truncate table with cascade (PostgreSQL)
truncateTableCascade := ddl.TruncateTable("users").Cascade().WithDialect(ddl.Postgres())
sql, _, err := truncateTableCascade.Build()
// sql: "TRUNCATE TABLE \"users\" CASCADE"

// Truncate table with restart identity (PostgreSQL)
truncateTableRestart := ddl.TruncateTable("users").Restart().WithDialect(ddl.Postgres())
sql, _, err := truncateTableRestart.Build()
// sql: "TRUNCATE TABLE \"users\" RESTART IDENTITY"
```

### Index Operations
```go
// Create index
createIndex := ddl.CreateIndex("idx_users_name", "users").Columns("name")
sql, _, err := createIndex.Build()
// sql: "CREATE INDEX `idx_users_name` ON `users` (`name`)"

// Drop index
dropIndex := ddl.DropIndex("idx_users_name", "users")
sql, _, err := dropIndex.Build()
// sql: "DROP INDEX `idx_users_name` ON `users`"
```

### View Operations
```go
// Create view
subq := stk.Select("name", "COUNT(*) as count").From("users").GroupBy("name")
createView := ddl.CreateView("user_stats").As(subq)
sql, _, err := createView.Build()
// sql: "CREATE VIEW `user_stats` AS SELECT `name`, COUNT(*) as count FROM `users` GROUP BY `name`"

// Drop view
dropView := ddl.DropView("user_stats").IfExists()
sql, _, err := dropView.Build()
// sql: "DROP VIEW IF EXISTS `user_stats`"
```

## SQL Dialect Examples

#### MySQL (default)
```go
// No need to set dialect for MySQL, it's the default
q := stk.Select("id", "name").From("users").Where("id = ?", 1)
sql, args, err := q.Build()
// sql: "SELECT `id`, `name` FROM `users` WHERE id = ?"
```

#### Postgres
```go
stk.SetDialect(stk.Postgres())
q := stk.Select("id", "name").From("users").Where("id = ? AND name = ?", 1, "bob")
sql, args, err := q.Build()
// sql: "SELECT \"id\", \"name\" FROM \"users\" WHERE id = $1 AND name = $2"
```

## Status
Work in progress. 