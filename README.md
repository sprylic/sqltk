# Sprylic - SQL Tool Kit

[![CI](https://github.com/sprylic/sqltk/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/sprylic/sqltk/actions/workflows/ci.yml)

A SQL builder for Go with composable query building and DDL operations.

## Warning
This project is still in development and is not yet ready for production use. 

A considerable amount of the code was generated using AI. Notably, almost all the tests. 
I tried to keep a close eye on any generated code, but there could still be some bugs or vulnerabilities.

## Example Usage

### SELECT
```go
q := sqltk.Select("id", "name").From("users").Where(sqltk.NewStringCondition("active = ?", true))
sql, args, err := q.Build()
// sql: "SELECT `id`, `name` FROM `users` WHERE active = ?" (MySQL dialect by default)
// args: [true]
```

### Aliasing and Subqueries
```go
import "github.com/sprylic/sqltk/raw"

sub := sqltk.Select(raw.Raw("COUNT(*)")).From("orders").WhereEqual("orders.user_id", "users.id")
q := sqltk.Select(sqltk.Alias(sub, "order_count")).From("users")
sql, args, err := q.Build()
// sql: "SELECT (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) AS order_count FROM `users`"
```

### Query Composition
```go
isActive := sqltk.Select().WhereEqual("status", 1)
isAdult := sqltk.Select().WhereGreaterThanOrEqual("age", 18)

q := sqltk.Select("id").From("users").Compose(isActive, isAdult)
sql, args, err := q.Build()
// sql: "SELECT `id` FROM `users` WHERE status = ? AND age >= ?"
// args: [1, 18]
```

### Condition Builder

The `ConditionBuilder` provides a composable API for building complex SQL conditions without resorting to raw SQL. Use `NewCond()` to start a condition chain, and pass it to `.Where()` or `.Having()` in any builder (`Select`, `Update`, `Delete`).

**Type-Safe Conditions:**

The library provides a `Condition` interface for type-safe condition handling. The `Where()` and `Having()` methods accept only `Condition` types:

```go
import "github.com/sprylic/sqltk/raw"

// string condition
cond := sqltk.NewStringCondition("active = ? AND age > ?", true, 18)
q := sqltk.Select("id").From("users").Where(cond)

// raw condition

cond := raw.Cond("id = 1")
q := sqltk.Select("id").From("users").Where(cond)

// ConditionBuilder
cond := sqltk.NewCond().Equal("active", true).And(sqltk.NewCond().GreaterThan("age", 18))
q := sqltk.Select("id").From("users").Where(cond)
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
- `Raw` - for raw SQL conditions (directly implements Condition)
- `*ConditionBuilder` - for fluent condition building

**Examples:**

SELECT with conditions
```go
cond := sqltk.NewCond().
    Equal("active", true).
    And(sqltk.NewCond().GreaterThan("age", 18)).
    And(sqltk.NewCond().In("status", "active", "pending"))
q := sqltk.Select("id", "name").From("users").Where(cond)
sql, args, err := q.Build()
// sql: "SELECT `id`, `name` FROM `users` WHERE active = ? AND age > ? AND status IN (?, ?)"
// args: [true, 18, "active", "pending"]
```

UPDATE with condition builder
```go
cond := sqltk.NewCond().
    Equal("active", true).
    Or(sqltk.NewCond().Equal("vip", true)).
    And(sqltk.NewCond().GreaterThan("age", 16))
q := sqltk.Update("users").Set("name", "Alice").Where(cond)
sql, args, err := q.Build()
// sql: "UPDATE `users` SET `name` = ? WHERE (active = ?) OR (vip = ?) AND age > ?"
// args: ["Alice", true, true, 16]
```

DELETE with condition builder
```go
cond := sqltk.NewCond().
    Equal("active", false).
    Or(sqltk.NewCond().IsNull("deleted_at"))
q := sqltk.Delete("users").Where(cond)
sql, args, err := q.Build()
// sql: "DELETE FROM `users` WHERE (active = ?) OR (deleted_at IS NULL)"
// args: [false]
```

**Supported condition methods:**
- `Equal`, `NotEqual`, `GreaterThan`, `LessThan`, `GreaterThanOrEqual`, `LessThanOrEqual`
- `In`, `NotIn`, `Between`, `NotBetween`, `IsNull`, `IsNotNull`, `Like`, `NotLike`
- `Exists`, `NotExists`, `Case`, `And`, `Or`
- All methods are chainable and support table-qualified columns.

**Type Safety:**
String conditions must use explicit wrappers to prevent SQL injection:
```go
// ❌ This will cause an error
q := sqltk.Select("id").From("users").Where("active = ?", 1)

// ✅ Use explicit wrapper for string conditions
q := sqltk.Select("id").From("users").Where(sqltk.NewStringCondition("active = ?", 1))

// ✅ Raw conditions work directly (implements Condition interface)
import "github.com/sprylic/sqltk/raw"

q := sqltk.Select("id").From("users").Where(raw.Cond("active = 1"))

// For simple where clauses, it's best to use one of the specific Where methods (WhereEqual, WhereNotEqual, etc.)
q := sqltk.Select("id").From("users").WhereEqual("active", 1),
```

### INSERT
```go
q := sqltk.Insert("users").Columns("id", "name").Values(1, "Alice").Values(2, "Bob")
sql, args, err := q.Build()
// sql: "INSERT INTO `users` (`id`, `name`) VALUES (?, ?), (?, ?)"
// args: [1, "Alice", 2, "Bob"]
```

### UPDATE
```go
q := sqltk.Update("users").Set("name", "Alice").WhereEqual("id", 1)
sql, args, err := q.Build()
// sql: "UPDATE `users` SET `name` = ? WHERE id = ?"
// args: ["Alice", 1]
```

### DELETE
```go
q := sqltk.Delete("users").WhereEqual("id", 1)
sql, args, err := q.Build()
// sql: "DELETE FROM `users` WHERE id = ?"
// args: [1]
```

## SQL Dialect
**MySQL is the default dialect.**
- Identifiers are quoted with backticks (`` `foo` ``) and placeholders are `?`.
- If you use Postgres or another database, **set the dialect explicitly**:
```go 
import "github.com/sprylic/sqltk/sqldialect"

sqltk.SetDialect(sqldialect.Postgres()) // for Postgres

builder := sqltk.Select("id", "name").From("users").
	WhereEqual("active", 1).
	WithDialect(sqldialect.Postgres()) // set dialect per builder.
sql, args, err := builder.Build()
```

### Warning:
Using the global dialect can be problematic when using different dialects concurrently. If you need to support a different dialect, use WithDialect on the builder instead.

## DDL Operations

### Database Operations
```go
import "github.com/sprylic/sqltk/ddl"

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
// Create table with primary key specified on column
createTable := ddl.CreateTable("users").
    AddColumn(ddl.Column("id").Type("INT").AutoIncrement().NotNull().PrimaryKey()).
    AddColumn(ddl.Column("name").Type("VARCHAR").Size(255).NotNull()).
	AddColumns(
        ddl.Column("email").Type("VARCHAR").
			Size(255).
			NotNull(),
		ddl.Column("created_at").Type("DATETIME").
			NotNull().
			Default(mysqlfunc.CurrentTimestamp()).,
        ddl.Column("updated_at").Type("DATETIME").
			NotNull().
			Default(mysqlfunc.CurrentTimestamp()).
            OnUpdate(mysqlfunc.CurrentTimestamp()), 
    ).
	Unique("idx_email", "email")
sql, _, err := createTable.Build()
// sql: "CREATE TABLE `users` (`id` INT AUTO_INCREMENT NOT NULL, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255), PRIMARY KEY (`id`), CONSTRAINT idx_email UNIQUE (`email`))"

// Create table with composite primary key (multiple columns)
createTableWithCompositePK := ddl.CreateTable("user_roles").
	AddColumns(
        ddl.Column("user_id").Type("INT").NotNull().PrimaryKey(),
        ddl.Column("role_id").Type("INT").NotNull().PrimaryKey(),
    )
sql, _, err := createTableWithCompositePK.Build()
// sql: "CREATE TABLE `user_roles` (`user_id` INT NOT NULL, `role_id` INT NOT NULL, PRIMARY KEY (`user_id`, `role_id`))"

// Alternative: specify primary key separately (legacy method)
createTableLegacy := ddl.CreateTable("users").
    AddColumns(
        ddl.Column("id").Type("INT").AutoIncrement().NotNull(),
        ddl.Column("name").Type("VARCHAR").Size(255).NotNull()
    ).PrimaryKey("id")
sql, _, err := createTableLegacy.Build()
// sql: "CREATE TABLE `users` (`id` INT AUTO_INCREMENT NOT NULL, `name` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`))"

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
import "github.com/sprylic/sqltk/sqldialect"

truncateTableCascade := ddl.TruncateTable("users").Cascade().WithDialect(sqldialect.Postgres())
sql, _, err := truncateTableCascade.Build()
// sql: "TRUNCATE TABLE \"users\" CASCADE"

// Truncate table with restart identity (PostgreSQL)
truncateTableRestart := ddl.TruncateTable("users").Restart().WithDialect(sqldialect.Postgres())
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
subq := sqltk.Select("name", "COUNT(*) as count").From("users").GroupBy("name")
createView := ddl.CreateView("user_stats").As(subq)
sql, _, err := createView.Build()
// sql: "CREATE VIEW `user_stats` AS SELECT `name`, COUNT(*) as count FROM `users` GROUP BY `name`"

// Drop view
dropView := ddl.DropView("user_stats").IfExists()
sql, _, err := dropView.Build()
// sql: "DROP VIEW IF EXISTS `user_stats`"
```

## Database Function Helpers

Helper functions provided for common database operations, making it easier to write database-specific SQL without using raw strings.

### MySQL Functions

```go
import "github.com/sprylic/sqltk/mysqlfunc"

// Date/Time functions
q := sqltk.Select(mysqlfunc.CurrentTimestamp()).From("users")
// SELECT CURRENT_TIMESTAMP FROM users

// String functions
q := sqltk.Select(
    sqltk.Alias(mysqlfunc.Concat("first_name", "' '", "last_name"), "full_name"),
).From("users")
// SELECT CONCAT(first_name, ' ', last_name) AS full_name FROM users

// Aggregate functions
q := sqltk.Select(
    sqltk.Alias(mysqlfunc.Count("*"), "total_users"),
    sqltk.Alias(mysqlfunc.Avg("age"), "avg_age"),
).From("users")
// SELECT COUNT(*) AS total_users, AVG(age) AS avg_age FROM users

// Conditional functions
q := sqltk.Select(
    sqltk.Alias(mysqlfunc.If("active", "'Active'", "'Inactive'"), "status"),
).From("users")
// SELECT IF(active, 'Active', 'Inactive') AS status FROM users
```

### PostgreSQL Functions

```go
import "github.com/sprylic/sqltk/pgfunc"

// Date/Time functions
q := sqltk.Select(pgfunc.Now()).From("users")
// SELECT now() FROM users

// String functions
q := sqltk.Select(
    sqltk.Alias(pgfunc.Concat("first_name", "' '", "last_name"), "full_name"),
).From("users")
// SELECT concat(first_name, ' ', last_name) AS full_name FROM users

// Aggregate functions
q := sqltk.Select(
    sqltk.Alias(pgfunc.Count("*"), "total_users"),
    sqltk.Alias(pgfunc.Avg("age"), "avg_age"),
).From("users")
// SELECT count(*) AS total_users, avg(age) AS avg_age FROM users

// Array functions
q := sqltk.Select(
    sqltk.Alias(pgfunc.ArrayAgg("tag"), "all_tags"),
).From("posts").GroupBy("category")
// SELECT array_agg(tag) AS all_tags FROM posts GROUP BY category
```

### Using Functions in WHERE Clauses

```go
import "github.com/sprylic/sqltk/raw"
// MySQL
q := sqltk.Select("id", "name").From("users").Where(
    raw.Cond("created_at > " + string(mysqlfunc.CurrentTimestamp())),
)
// SELECT id, name FROM users WHERE created_at > CURRENT_TIMESTAMP

// PostgreSQL
q := sqltk.Select("id", "name").From("users").Where(
    raw.Cond("created_at > " + string(pgfunc.Now())),
)
// SELECT id, name FROM users WHERE created_at > now()
```

## Available Functions

### MySQL Functions (`mysqlfunc`)

**Date/Time Functions:**
- `CurrentTimestamp()`, `CurrentDate()`, `CurrentTime()`, `Now()`, `UnixTimestamp()`
- `DateFormat(expr, format)`, `DateAdd(date, interval)`, `DateSub(date, interval)`
- `Year(date)`, `Month(date)`, `Day(date)`, `Hour(time)`, `Minute(time)`, `Second(time)`

**String Functions:**
- `Concat(args...)`, `ConcatWs(separator, args...)`
- `Upper(str)`, `Lower(str)`, `Length(str)`, `Substring(str, pos, len)`
- `Trim(str)`, `Ltrim(str)`, `Rtrim(str)`, `Replace(str, from, to)`

**Numeric Functions:**
- `Abs(num)`, `Ceiling(num)`, `Floor(num)`, `Round(num, decimals)`
- `Mod(dividend, divisor)`, `Power(base, exponent)`, `Sqrt(num)`
- `Random()`, `Pi()`, `E()`

**Aggregate Functions:**
- `Count(expr)`, `Sum(expr)`, `Avg(expr)`, `Min(expr)`, `Max(expr)`
- `GroupConcat(expr, separator)`

**Conditional Functions:**
- `If(condition, trueVal, falseVal)`, `IfNull(expr, nullVal)`, `NullIf(expr1, expr2)`

**Type Conversion:**
- `Cast(expr, asType)`, `Convert(expr, asType)`

**JSON Functions:**
- `JsonExtract(jsonDoc, path)`, `JsonUnquote(jsonVal)`, `JsonLength(jsonDoc, path)`

### PostgreSQL Functions (`pgfunc`)

**Date/Time Functions:**
- `Now()`, `CurrentTimestamp()`, `CurrentDate()`, `CurrentTime()`
- `ClockTimestamp()`, `StatementTimestamp()`, `TransactionTimestamp()`
- `Extract(field, source)`, `DatePart(field, source)`
- `DateTrunc(field, source)`, `Age(timestamp)`, `Now()`

**String Functions:**
- `Concat(args...)`, `ConcatWs(separator, args...)`
- `Upper(str)`, `Lower(str)`, `Length(str)`, `Substr(str, from, count)`
- `Trim(str)`, `Ltrim(str)`, `Rtrim(str)`, `Replace(str, from, to)`

**Numeric Functions:**
- `Abs(num)`, `Ceiling(num)`, `Floor(num)`, `Round(num, decimals)`
- `Trunc(num, decimals)`, `Mod(dividend, divisor)`, `Power(base, exponent)`
- `Sqrt(num)`, `Random()`, `Pi()`, `E()`

**Aggregate Functions:**
- `Count(expr)`, `Sum(expr)`, `Avg(expr)`, `Min(expr)`, `Max(expr)`
- `StringAgg(expr, delimiter)`, `ArrayAgg(expr)`, `JsonAgg(expr)`, `JsonbAgg(expr)`

**Conditional Functions:**
- `Coalesce(args...)`, `NullIf(expr1, expr2)`, `Greatest(args...)`, `Least(args...)`

**Type Conversion:**
- `Cast(expr, asType)`, `Convert(expr, asType)`

**JSON Functions:**
- `JsonExtract(jsonDoc, path)`, `JsonUnquote(jsonVal)`, `JsonLength(jsonDoc, path)`

## Examples

See the `examples/` directory for more detailed examples:

## Installation

```bash
go get github.com/sprylic/sqltk
```
