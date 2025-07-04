package stk

import (
	"reflect"
	"strings"
	"testing"
)

func TestSelectBuilder(t *testing.T) {
	t.Run("basic select", func(t *testing.T) {
		q := Select("id", "name").From("users")
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id, name FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if len(args) != 0 {
			t.Errorf("got args %v, want none", args)
		}
	})

	t.Run("select with where", func(t *testing.T) {
		q := Select("id").From("users").Where(NewStringCondition("active = ?", true))
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE active = ?"
		wantArgs := []interface{}{true}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("where equal", func(t *testing.T) {
		q := Select("id").From("users").WhereEqual("active", true)
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE active = ?"
		wantArgs := []interface{}{true}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("where not equal", func(t *testing.T) {
		q := Select("id").From("users").WhereNotEqual("active", false)
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE active != ?"
		wantArgs := []interface{}{false}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("where column is null", func(t *testing.T) {
		q := Select("id").From("users").WhereEqual("deleted_at", nil)
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE deleted_at IS NULL"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if len(args) != 0 {
			t.Errorf("got args %v, want none", args)
		}
	})

	t.Run("where column is not null", func(t *testing.T) {
		q := Select("id").From("users").WhereNotEqual("created_at", nil)
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE created_at IS NOT NULL"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if len(args) != 0 {
			t.Errorf("got args %v, want none", args)
		}
	})

	t.Run("select all columns", func(t *testing.T) {
		q := Select().From("users")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT * FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("multiple where clauses", func(t *testing.T) {
		q := Select("id").From("users").Where(NewStringCondition("active = ?", true)).Where(NewStringCondition("age > ?", 18))
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE active = ? AND age > ?"
		wantArgs := []interface{}{true, 18}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})
}

func TestSelectBuilder_RawWhere(t *testing.T) {
	t.Run("raw where only", func(t *testing.T) {
		q := Select("id").From("users").Where(AsCondition(Raw("age > 18")))
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE age > 18"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if len(args) != 0 {
			t.Errorf("got args %v, want none", args)
		}
	})

	t.Run("mixed parameterized and raw", func(t *testing.T) {
		q := Select("id").From("users").Where(NewStringCondition("active = ?", true)).Where(AsCondition(Raw("age > 18")))
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE active = ? AND age > 18"
		wantArgs := []interface{}{true}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("error on invalid type", func(t *testing.T) {
		// This test demonstrates that the compiler will catch invalid types
		// We can't test this at runtime since it's a compile-time error
		t.Skip("This is now a compile-time error, not a runtime error")
	})
}

func TestSelectBuilder_GroupBy_Having_OrderBy(t *testing.T) {
	t.Run("group by column", func(t *testing.T) {
		q := Select("id").AddField("COUNT(*)").From("users").GroupBy("id")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id, COUNT(*) FROM users GROUP BY id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("group by raw", func(t *testing.T) {
		q := Select("id").From("users").GroupBy(Raw("LEFT(name, 1)"))
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users GROUP BY LEFT(name, 1)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("having parameterized", func(t *testing.T) {
		q := Select("id").From("users").GroupBy("id").Having(NewStringCondition("COUNT(*) > ?", 1))
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users GROUP BY id HAVING COUNT(*) > ?"
		wantArgs := []interface{}{1}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("having raw", func(t *testing.T) {
		q := Select("id").From("users").GroupBy("id").Having(AsCondition(Raw("COUNT(*) > 1")))
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users GROUP BY id HAVING COUNT(*) > 1"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if len(args) != 0 {
			t.Errorf("got args %v, want none", args)
		}
	})

	t.Run("order by column", func(t *testing.T) {
		q := Select("id").From("users").OrderBy("id DESC")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users ORDER BY id DESC"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("order by raw", func(t *testing.T) {
		q := Select("id").From("users").OrderBy(Raw("RANDOM()"))
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users ORDER BY RANDOM()"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("full query with all clauses", func(t *testing.T) {
		q := Select("id").From("users").
			Where(NewStringCondition("active = ?", true)).
			GroupBy("id").
			Having(NewStringCondition("COUNT(*) > ?", 1)).
			OrderBy("id DESC")
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE active = ? GROUP BY id HAVING COUNT(*) > ? ORDER BY id DESC"
		wantArgs := []interface{}{true, 1}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("error on invalid group by type", func(t *testing.T) {
		_, _, err := Select("id").From("users").GroupBy(123).Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("error on invalid having type", func(t *testing.T) {
		// This test demonstrates that the compiler will catch invalid types
		// We can't test this at runtime since it's a compile-time error
		t.Skip("This is now a compile-time error, not a runtime error")
	})

	t.Run("error on invalid order by type", func(t *testing.T) {
		_, _, err := Select("id").From("users").OrderBy(123).Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})
}

func TestSelectBuilder_Join_Limit_Offset(t *testing.T) {
	t.Run("inner join fluent", func(t *testing.T) {
		q := Select("u.id", "p.id").From("users u").Join("posts p").On("p.user_id", "u.id")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT u.id, p.id FROM users u JOIN posts p ON p.user_id = u.id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("left join fluent", func(t *testing.T) {
		q := Select("u.id", "p.id").From("users u").LeftJoin("posts p").On("p.user_id", "u.id")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT u.id, p.id FROM users u LEFT JOIN posts p ON p.user_id = u.id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("multiple joins fluent", func(t *testing.T) {
		q := Select("u.id", "p.id", "c.id").From("users u").
			Join("posts p").On("p.user_id", "u.id").
			LeftJoin("comments c").On("c.post_id", "p.id")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT u.id, p.id, c.id FROM users u JOIN posts p ON p.user_id = u.id LEFT JOIN comments c ON c.post_id = p.id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("limit and offset", func(t *testing.T) {
		q := Select("id").From("users").Limit(10).Offset(5)
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users LIMIT 10 OFFSET 5"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("limit only", func(t *testing.T) {
		q := Select("id").From("users").Limit(1)
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users LIMIT 1"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("offset only", func(t *testing.T) {
		q := Select("id").From("users").Offset(2)
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users OFFSET 2"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("full query with join, limit, offset", func(t *testing.T) {
		q := Select("u.id", "p.id").From("users u").
			Join("posts p").On("p.user_id", "u.id").
			Where(NewStringCondition("u.active = ?", true)).
			OrderBy("u.id DESC").
			Limit(20).
			Offset(10)
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT u.id, p.id FROM users u JOIN posts p ON p.user_id = u.id WHERE u.active = ? ORDER BY u.id DESC LIMIT 20 OFFSET 10"
		wantArgs := []interface{}{true}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("join with alias expr", func(t *testing.T) {
		q := Select("u.id", "p.id").From(Alias("users", "u")).Join(Alias("posts", "p")).On("p.user_id", "u.id")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT u.id, p.id FROM users AS u JOIN posts AS p ON p.user_id = u.id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("join error propagation", func(t *testing.T) {
		// Test invalid join table type
		q := Select("id").From("users").Join(123).On("user_id", "id")
		_, _, err := q.Build()
		if err == nil {
			t.Error("expected error for invalid join table type")
		}
		if !strings.Contains(err.Error(), "join: table must be string, Raw, *SelectBuilder, or AliasExpr") {
			t.Errorf("expected specific error message, got: %v", err)
		}
	})
}

func TestSelectBuilder_Distinct_Subquery(t *testing.T) {
	t.Run("distinct", func(t *testing.T) {
		q := Select("id").Distinct().From("users")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT DISTINCT id FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("subquery as column", func(t *testing.T) {
		sub := Select("COUNT(*)").From("posts").Where(Raw("posts.user_id = users.id"))
		q := Select("id", sub).From("users")
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id, (SELECT COUNT(*) FROM posts WHERE posts.user_id = users.id) FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if len(args) != 0 {
			t.Errorf("got args %v, want none", args)
		}
	})

	t.Run("subquery as column with args", func(t *testing.T) {
		sub := Select("COUNT(*)").From("posts").Where(NewStringCondition("posts.user_id = ?", 42))
		q := Select("id", sub).From("users")
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id, (SELECT COUNT(*) FROM posts WHERE posts.user_id = ?) FROM users"
		wantArgs := []interface{}{42}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("subquery in FROM", func(t *testing.T) {
		sub := Select("id").From("posts").Where(NewStringCondition("published = ?", true))
		q := Select("id").From(sub)
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM (SELECT id FROM posts WHERE published = ?)"
		wantArgs := []interface{}{true}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("error on invalid select column type", func(t *testing.T) {
		_, _, err := Select(123).From("users").Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("error on invalid from type", func(t *testing.T) {
		_, _, err := Select("id").From(123).Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})
}

func TestSelectBuilder_Alias(t *testing.T) {
	t.Run("alias column", func(t *testing.T) {
		q := Select(Alias("id", "user_id")).From("users")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id AS user_id FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("alias subquery as column", func(t *testing.T) {
		sub := Select("COUNT(*)").From("orders").Where(NewStringCondition("orders.user_id = users.id"))
		q := Select(Alias(sub, "order_count")).From("users")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) AS order_count FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("alias subquery in FROM", func(t *testing.T) {
		sub := Select("id").From("orders").Where(NewStringCondition("amount > ?", 100))
		q := Select("o.id").From(Alias(sub, "o"))
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT o.id FROM (SELECT id FROM orders WHERE amount > ?) AS o"
		wantArgs := []interface{}{100}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("alias raw in FROM", func(t *testing.T) {
		q := Select("u.id").From(Alias(Raw("users u"), "u_alias"))
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT u.id FROM users u AS u_alias"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("error on invalid alias expr type", func(t *testing.T) {
		q := Select(Alias(123, "bad")).From("users")
		_, _, err := q.WithDialect(Standard()).Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})
}

func TestSelectBuilder_Compose(t *testing.T) {
	isActive := func(b *SelectBuilder) *SelectBuilder {
		return b.Where(NewStringCondition("active = ?", true))
	}
	isAdult := func(b *SelectBuilder) *SelectBuilder {
		return b.Where(NewStringCondition("age >= ?", 18))
	}

	t.Run("compose single fragment", func(t *testing.T) {
		q := Select("id").From("users").Compose(isActive)
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE active = ?"
		wantArgs := []interface{}{true}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("compose multiple fragments", func(t *testing.T) {
		q := Select("id").From("users").Compose(isActive, isAdult)
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE active = ? AND age >= ?"
		wantArgs := []interface{}{true, 18}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("compose order matters", func(t *testing.T) {
		first := func(b *SelectBuilder) *SelectBuilder { return b.Where(NewStringCondition("x = ?", 1)) }
		second := func(b *SelectBuilder) *SelectBuilder { return b.Where(NewStringCondition("y = ?", 2)) }
		q := Select("id").From("users").Compose(first, second)
		sql, args, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id FROM users WHERE x = ? AND y = ?"
		wantArgs := []interface{}{1, 2}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("compose propagates error", func(t *testing.T) {
		// This test demonstrates that the compiler will catch invalid types
		// We can't test this at runtime since it's a compile-time error
		t.Skip("This is now a compile-time error, not a runtime error")
	})
}

func TestSelectBuilder_Dialect(t *testing.T) {
	t.Run("standard dialect", func(t *testing.T) {
		q := Select("id", "name").From("users").Where(NewStringCondition("id = ? AND name = ?", 1, "bob")).WithDialect(Standard())
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT id, name FROM users WHERE id = ? AND name = ?"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("mysql dialect", func(t *testing.T) {
		q := Select("id", "name").From("users").Where(NewStringCondition("id = ? AND name = ?", 1, "bob")).WithDialect(MySQL())
		sql, _, err := q.WithDialect(MySQL()).Build()
		wantSQL := "SELECT `id`, `name` FROM `users` WHERE id = ? AND name = ?"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("postgres dialect", func(t *testing.T) {
		q := Select("id", "name").From("users").Where(NewStringCondition("id = ? AND name = ?", 1, "bob")).WithDialect(Postgres())
		sql, _, err := q.WithDialect(Postgres()).Build()
		wantSQL := "SELECT \"id\", \"name\" FROM \"users\" WHERE id = $1 AND name = $2"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}

func TestSelectBuilder_FluentJoinOn(t *testing.T) {
	q := Select("u.id").From("users u").Join("orders o").On("o.user_id", "u.id")
	sql, _, err := q.WithDialect(Standard()).Build()
	wantSQL := "SELECT u.id FROM users u JOIN orders o ON o.user_id = u.id"
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sql != wantSQL {
		t.Errorf("got SQL %q, want %q", sql, wantSQL)
	}
}

func TestSelectBuilder_FluentJoinTypes(t *testing.T) {
	t.Run("left join", func(t *testing.T) {
		q := Select("u.id").From("users u").LeftJoin("orders o").On("o.user_id", "u.id")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT u.id FROM users u LEFT JOIN orders o ON o.user_id = u.id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("right join", func(t *testing.T) {
		q := Select("u.id").From("users u").RightJoin("orders o").On("o.user_id", "u.id")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT u.id FROM users u RIGHT JOIN orders o ON o.user_id = u.id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("full join", func(t *testing.T) {
		q := Select("u.id").From("users u").FullJoin("orders o").On("o.user_id", "u.id")
		sql, _, err := q.WithDialect(Standard()).Build()
		wantSQL := "SELECT u.id FROM users u FULL JOIN orders o ON o.user_id = u.id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}
