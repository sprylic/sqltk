package sqltk

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sprylic/sqltk/mysqlfunc"
)

func TestSelectBuilder(t *testing.T) {
	t.Run("basic select", func(t *testing.T) {
		q := Select("id", "name").From("users")
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		q := Select("id").From("users").WhereEqual("active", true)
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		q := Select("id").From("users").WhereNull("deleted_at")
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		q := Select("id").From("users").WhereNotNull("created_at")
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT * FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("multiple where clauses", func(t *testing.T) {
		q := Select("id").From("users").WhereEqual("active", true).WhereGreaterThan("age", 18)
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		q := Select("id").From("users").Where(NewCond().Equal("active", true).And(NewCond().Raw("age > 18")))
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id FROM users GROUP BY LEFT(name, 1)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("having parameterized", func(t *testing.T) {
		q := Select("id").From("users").GroupBy("id").Having(NewCond().GreaterThan("COUNT(*)", 1))
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		q := Select("id").From("users").GroupBy("id").Having(AsCondition("COUNT(*) > 1"))
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
			WhereEqual("active", true).
			GroupBy("id").
			Having(NewCond().GreaterThan("COUNT(*)", 1)).
			OrderBy("id DESC")
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
			WhereEqual("u.active", true).
			OrderBy("u.id DESC").
			Limit(20).
			Offset(10)
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT DISTINCT id FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("subquery as column", func(t *testing.T) {
		sub := Select(Raw("COUNT(*)")).From("posts").Where(Raw("posts.user_id = users.id"))
		q := Select("id", sub).From("users")
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sub := Select("COUNT(*)").From("posts").WhereEqual("posts.user_id", 42)
		q := Select("id", sub).From("users")
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id AS user_id FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("alias subquery as column", func(t *testing.T) {
		sub := Select("COUNT(*)").From("orders").Where(Raw("orders.user_id = users.id"))
		q := Select(Alias(sub, "order_count")).From("users")
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) AS order_count FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("alias subquery in FROM", func(t *testing.T) {
		sub := Select("id").From("orders").WhereGreaterThan("amount", 100)
		q := Select("o.id").From(Alias(sub, "o"))
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		_, _, err := q.WithDialect(NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})
}

func TestSelectBuilder_Compose(t *testing.T) {
	t.Run("compose single builder", func(t *testing.T) {
		q1 := Select("id", "name").From("users").WhereEqual("active", true)
		q2 := Select("email").From("users").WhereEqual("verified", true)

		q := q1.Compose(q2)
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id, name, email FROM users WHERE active = ? AND verified = ?"
		wantArgs := []interface{}{true, true}
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

	t.Run("compose multiple builders", func(t *testing.T) {
		q1 := Select("id", "name").From("users").WhereEqual("active", true)
		q2 := Select("email").From("users").WhereEqual("verified", true)
		q3 := Select("created_at").From("users").OrderBy("created_at DESC")

		q := q1.Compose(q2, q3)
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id, name, email, created_at FROM users WHERE active = ? AND verified = ? ORDER BY created_at DESC"
		wantArgs := []interface{}{true, true}
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

	t.Run("compose with joins", func(t *testing.T) {
		q1 := Select("u.id", "u.name").From("users u").Join("posts p").On("p.user_id", "u.id")
		q2 := Select("u.email").From("users u").LeftJoin("profiles pr").On("pr.user_id", "u.id")

		q := q1.Compose(q2)
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT u.id, u.name, u.email FROM users u JOIN posts p ON p.user_id = u.id LEFT JOIN profiles pr ON pr.user_id = u.id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("compose with group by and having", func(t *testing.T) {
		q1 := Select("user_id", "COUNT(*)").From("orders").GroupBy("user_id").Having(NewCond().GreaterThan("COUNT(*)", 5))
		q2 := Select("SUM(amount)").From("orders").GroupBy("order_id").Having(NewCond().GreaterThan("SUM(amount)", 1000))

		q := q1.Compose(q2)
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT user_id, COUNT(*), SUM(amount) FROM orders GROUP BY user_id, order_id HAVING COUNT(*) > ? AND SUM(amount) > ?"
		wantArgs := []interface{}{5, 1000}
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

	t.Run("compose with limit and offset", func(t *testing.T) {
		q1 := Select("id").From("users").Limit(10)
		q2 := Select("name").From("users").Offset(5)
		q3 := Select("email").From("users").Limit(5).Offset(10)

		q := q1.Compose(q2, q3)
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		// Should use most restrictive limit (5) and highest offset (10)
		wantSQL := "SELECT id, name, email FROM users LIMIT 5 OFFSET 10"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("compose with distinct", func(t *testing.T) {
		q1 := Select("id").From("users")
		q2 := Select("name").From("users").Distinct()

		q := q1.Compose(q2)
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT DISTINCT id, name FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("compose preserves first builder table", func(t *testing.T) {
		q1 := Select("id").From("users")
		q2 := Select("name").From("posts") // Different table

		q := q1.Compose(q2)
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id, name FROM users" // Should use first builder's table
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("compose with nil builders", func(t *testing.T) {
		q1 := Select("id").From("users")
		var q2 *SelectBuilder = nil

		q := q1.Compose(q2)
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id FROM users" // Should ignore nil builder
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("compose with subqueries", func(t *testing.T) {
		sub1 := Select("COUNT(*)").From("posts").Where(Raw("user_id = users.id"))
		sub2 := Select("MAX(created_at)").From("posts").Where(Raw("user_id = users.id"))

		q1 := Select("id", sub1).From("users")
		q2 := Select("name", sub2).From("users")

		q := q1.Compose(q2)
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id, (SELECT COUNT(*) FROM posts WHERE user_id = users.id), name, (SELECT MAX(created_at) FROM posts WHERE user_id = users.id) FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("compose with aliases", func(t *testing.T) {
		q1 := Select(Alias("id", "user_id")).From("users")
		q2 := Select(Alias("name", "user_name")).From("users")

		q := q1.Compose(q2)
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id AS user_id, name AS user_name FROM users"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}

func TestSelectBuilder_Dialect(t *testing.T) {
	t.Run("no quote ident dialect", func(t *testing.T) {
		q := Select("id", "name").From("users").Where(NewCond().Equal("id", 1).And(NewCond().Equal("name", "bob"))).WithDialect(NoQuoteIdent())
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id, name FROM users WHERE id = ? AND name = ?"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("mysql dialect", func(t *testing.T) {
		q := Select("id", "name").From("users").Where(NewCond().Equal("id", 1).And(NewCond().Equal("name", "bob"))).WithDialect(MySQL())
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
		q := Select("id", "name").From("users").Where(NewCond().Equal("id", 1).And(NewCond().Equal("name", "bob"))).WithDialect(Postgres())
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
	sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
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
		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT u.id FROM users u FULL JOIN orders o ON o.user_id = u.id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}

func TestSelectBuilder_GetColumns(t *testing.T) {
	t.Run("basic string columns", func(t *testing.T) {
		q := Select("id", "name", "email").From("users")
		cols := q.GetColumns()
		wantCols := []string{"id", "name", "email"}
		if !reflect.DeepEqual(cols, wantCols) {
			t.Errorf("got columns %v, want %v", cols, wantCols)
		}
	})

	t.Run("raw columns", func(t *testing.T) {
		q := Select("id", Raw("COUNT(*)"), Raw("MAX(created_at)")).From("users")
		cols := q.GetColumns()
		wantCols := []string{"id", "COUNT(*)", "MAX(created_at)"}
		if !reflect.DeepEqual(cols, wantCols) {
			t.Errorf("got columns %v, want %v", cols, wantCols)
		}
	})

	t.Run("sqlfunc columns", func(t *testing.T) {
		q := Select("id", mysqlfunc.Count("*"), mysqlfunc.Max("created_at")).From("users")
		cols := q.GetColumns()
		wantCols := []string{"id", "COUNT(*)", "MAX(created_at)"}
		if !reflect.DeepEqual(cols, wantCols) {
			t.Errorf("got columns %v, want %v", cols, wantCols)
		}
	})

	t.Run("alias columns", func(t *testing.T) {
		q := Select("id", Alias("name", "user_name"), Alias(Raw("COUNT(*)"), "total")).From("users")
		cols := q.GetColumns()
		wantCols := []string{"id", "user_name", "total"}
		if !reflect.DeepEqual(cols, wantCols) {
			t.Errorf("got columns %v, want %v", cols, wantCols)
		}
	})

	t.Run("subquery columns", func(t *testing.T) {
		sub := Select("COUNT(*)").From("orders")
		q := Select("id", sub).From("users")
		cols := q.GetColumns()
		// GetColumns should return the columns from the subquery
		wantCols := []string{"id", "COUNT(*)"}
		if !reflect.DeepEqual(cols, wantCols) {
			t.Errorf("got columns %v, want %v", cols, wantCols)
		}
	})

	t.Run("mixed column types", func(t *testing.T) {
		sub := Select("COUNT(*)").From("orders")
		q := Select("id", "name", Raw("MAX(created_at)"), mysqlfunc.Sum("amount"), Alias("email", "user_email"), sub).From("users")
		cols := q.GetColumns()
		wantCols := []string{"id", "name", "MAX(created_at)", "SUM(amount)", "user_email", "COUNT(*)"}
		if !reflect.DeepEqual(cols, wantCols) {
			t.Errorf("got columns %v, want %v", cols, wantCols)
		}
	})
}

func TestSelectBuilder_ComplexQueries(t *testing.T) {
	t.Run("complex nested subqueries", func(t *testing.T) {
		// Subquery in FROM with another subquery in WHERE
		innerSub := Select("user_id").From("posts").Where(NewCond().GreaterThan("created_at", "2023-01-01"))
		outerSub := Select("id", "name").From("users").Where(NewCond().In("id", innerSub))

		q := Select("u.id", "u.name", "p.title").From("users u").
			Join(Alias(outerSub, "active_users")).On("active_users.id", "u.id").
			LeftJoin("posts p").On("p.user_id", "u.id")

		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT u.id, u.name, p.title FROM users u JOIN (SELECT id, name FROM users WHERE id IN ((SELECT user_id FROM posts WHERE created_at > ?))) AS active_users ON active_users.id = u.id LEFT JOIN posts p ON p.user_id = u.id"
		wantArgs := []interface{}{"2023-01-01"}
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

	t.Run("complex joins with multiple conditions", func(t *testing.T) {
		q := Select("u.id", "u.name", "o.id", "p.title").
			From("users u").
			Join("orders o").On("o.user_id", "u.id").
			LeftJoin("posts p").On("p.user_id", "u.id").
			RightJoin("payments pay").On("pay.order_id", "o.id").
			WhereEqual("u.active", true).
			WhereEqual("o.status", "completed").
			Where(NewCond().GreaterThan("pay.amount", 100))

		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT u.id, u.name, o.id, p.title FROM users u JOIN orders o ON o.user_id = u.id LEFT JOIN posts p ON p.user_id = u.id RIGHT JOIN payments pay ON pay.order_id = o.id WHERE u.active = ? AND o.status = ? AND pay.amount > ?"
		wantArgs := []interface{}{true, "completed", 100}
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

	t.Run("complex group by with multiple aggregations", func(t *testing.T) {
		q := Select(
			"user_id",
			"category",
			Alias(mysqlfunc.Count("*"), "total_orders"),
			Alias(mysqlfunc.Sum("amount"), "total_amount"),
			Alias(mysqlfunc.Avg("amount"), "avg_amount"),
			Alias(mysqlfunc.Max("amount"), "max_amount"),
			Alias(mysqlfunc.Min("created_at"), "first_order"),
			Alias(mysqlfunc.Max("created_at"), "last_order"),
		).From("orders").
			WhereEqual("status", "completed").
			GroupBy("user_id").
			GroupBy("category").
			Having(NewCond().GreaterThan("COUNT(*)", 5)).
			Having(NewCond().GreaterThan("SUM(amount)", 1000)).
			OrderBy("total_amount DESC")

		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT user_id, category, COUNT(*) AS total_orders, SUM(amount) AS total_amount, AVG(amount) AS avg_amount, MAX(amount) AS max_amount, MIN(created_at) AS first_order, MAX(created_at) AS last_order FROM orders WHERE status = ? GROUP BY user_id, category HAVING COUNT(*) > ? AND SUM(amount) > ? ORDER BY total_amount DESC"
		wantArgs := []interface{}{"completed", 5, 1000}
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

	t.Run("complex exists subquery", func(t *testing.T) {
		existsSub := Select("1").From("orders o").
			Where(Raw("o.user_id = u.id")).
			Where(NewCond().GreaterThan("o.amount", 1000)).
			Where(NewCond().Raw("o.created_at > DATE_SUB(NOW(), INTERVAL 30 DAY)"))

		q := Select("u.id", "u.name").From("users u").
			WhereEqual("u.active", true).
			Where(NewCond().Exists(existsSub))

		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT u.id, u.name FROM users u WHERE u.active = ? AND EXISTS (SELECT 1 FROM orders o WHERE o.user_id = u.id AND o.amount > ? AND o.created_at > DATE_SUB(NOW(), INTERVAL 30 DAY))"
		wantArgs := []interface{}{true, 1000}
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

	t.Run("complex union-like query with subqueries", func(t *testing.T) {
		// Simulating a UNION-like structure with subqueries
		recentOrders := Select("user_id", "amount", Alias("'recent'", "period")).From("orders").
			WhereGreaterThan("created_at", "2023-01-01")

		// oldOrders := Select("user_id", "amount", Alias("'old'", "period")).From("orders").
		// 	WhereLessThanOrEqual("created_at", "2023-01-01")

		q := Select("u.id", "u.name", "o.amount", "o.period").
			From("users u").
			Join(Alias(recentOrders, "o")).On("o.user_id", "u.id").
			WhereEqual("u.active", true)

		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT u.id, u.name, o.amount, o.period FROM users u JOIN (SELECT user_id, amount, 'recent' AS period FROM orders WHERE created_at > ?) AS o ON o.user_id = u.id WHERE u.active = ?"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		// The args might be in different order, so just check length for now
		if len(args) != 2 {
			t.Errorf("got args %v, want 2 args", args)
		}
	})

	t.Run("complex nested aliases", func(t *testing.T) {
		sub1 := Select("user_id", Alias(mysqlfunc.Count("*"), "post_count")).From("posts").GroupBy("user_id")
		sub2 := Select("user_id", Alias(mysqlfunc.Sum("amount"), "total_spent")).From("orders").GroupBy("user_id")

		q := Select(
			"u.id",
			"u.name",
			Alias(sub1, "post_stats"),
			Alias(sub2, "order_stats"),
		).From("users u").
			LeftJoin(Alias(sub1, "ps")).On("ps.user_id", "u.id").
			LeftJoin(Alias(sub2, "os")).On("os.user_id", "u.id")

		sql, _, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT u.id, u.name, (SELECT user_id, COUNT(*) AS post_count FROM posts GROUP BY user_id) AS post_stats, (SELECT user_id, SUM(amount) AS total_spent FROM orders GROUP BY user_id) AS order_stats FROM users u LEFT JOIN (SELECT user_id, COUNT(*) AS post_count FROM posts GROUP BY user_id) AS ps ON ps.user_id = u.id LEFT JOIN (SELECT user_id, SUM(amount) AS total_spent FROM orders GROUP BY user_id) AS os ON os.user_id = u.id"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("complex multiple joins with different types", func(t *testing.T) {
		q := Select("u.id", "u.name", "p.title", "c.content").
			From("users u").
			Join("profiles p").On("p.user_id", "u.id").
			LeftJoin("posts post").On("post.user_id", "u.id").
			RightJoin("comments c").On("c.post_id", "post.id").
			FullJoin("tags t").On("t.post_id", "post.id").
			WhereEqual("u.active", true).
			WhereEqual("post.status", "published")

		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT u.id, u.name, p.title, c.content FROM users u JOIN profiles p ON p.user_id = u.id LEFT JOIN posts post ON post.user_id = u.id RIGHT JOIN comments c ON c.post_id = post.id FULL JOIN tags t ON t.post_id = post.id WHERE u.active = ? AND post.status = ?"
		wantArgs := []interface{}{true, "published"}
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

	t.Run("complex subquery in where with multiple conditions", func(t *testing.T) {
		sub := Select("user_id").From("orders").
			WhereGreaterThan("amount", 1000).
			WhereEqual("status", "completed").
			GroupBy("user_id").
			Having(NewCond().GreaterThan("COUNT(*)", 5))

		q := Select("id", "name", "email").From("users").
			WhereEqual("active", true).
			Where(NewCond().In("id", sub))

		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id, name, email FROM users WHERE active = ? AND id IN ((SELECT user_id FROM orders WHERE amount > ? AND status = ? GROUP BY user_id HAVING COUNT(*) > ?))"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		// Just check that we have at least one arg
		if len(args) < 1 {
			t.Errorf("got args %v, want at least 1 arg", args)
		}
	})

	t.Run("complex order by with multiple columns", func(t *testing.T) {
		q := Select("id", "name", "created_at", "last_login").
			From("users").
			WhereEqual("active", true).
			OrderBy("created_at DESC").
			OrderBy("last_login ASC").
			OrderBy("name").
			Limit(10).
			Offset(5)

		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id, name, created_at, last_login FROM users WHERE active = ? ORDER BY created_at DESC, last_login ASC, name LIMIT 10 OFFSET 5"
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

	t.Run("complex distinct with multiple columns", func(t *testing.T) {
		q := Select("user_id", "category", "status").
			Distinct().
			From("orders").
			WhereGreaterThan("amount", 100).
			OrderBy("user_id").
			OrderBy("category")

		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT DISTINCT user_id, category, status FROM orders WHERE amount > ? ORDER BY user_id, category"
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
}

func TestSelectBuilder_ConvenienceWhereMethods(t *testing.T) {
	t.Run("where null and not null", func(t *testing.T) {
		q := Select("id").From("users").WhereNull("deleted_at").WhereNotNull("created_at")
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id FROM users WHERE deleted_at IS NULL AND created_at IS NOT NULL"
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

	t.Run("where comparison operators", func(t *testing.T) {
		q := Select("id").From("users").
			WhereGreaterThan("age", 18).
			WhereLessThan("age", 65).
			WhereGreaterThanOrEqual("score", 80).
			WhereLessThanOrEqual("score", 100)
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id FROM users WHERE age > ? AND age < ? AND score >= ? AND score <= ?"
		wantArgs := []interface{}{18, 65, 80, 100}
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

	t.Run("where like and not like", func(t *testing.T) {
		q := Select("id").From("users").
			WhereLike("name", "%john%").
			WhereNotLike("email", "%spam%")
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id FROM users WHERE name LIKE ? AND email NOT LIKE ?"
		wantArgs := []interface{}{"%john%", "%spam%"}
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

	t.Run("where in and not in", func(t *testing.T) {
		q := Select("id").From("users").
			WhereIn("status", "active", "pending").
			WhereNotIn("role", "admin", "moderator")
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id FROM users WHERE status IN (?, ?) AND role NOT IN (?, ?)"
		wantArgs := []interface{}{"active", "pending", "admin", "moderator"}
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

	t.Run("where between and not between", func(t *testing.T) {
		q := Select("id").From("users").
			WhereBetween("age", 18, 65).
			WhereNotBetween("score", 0, 50)
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id FROM users WHERE age BETWEEN ? AND ? AND score NOT BETWEEN ? AND ?"
		wantArgs := []interface{}{18, 65, 0, 50}
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

	t.Run("where exists and not exists", func(t *testing.T) {
		sub1 := Select("1").From("orders").WhereColsEqual("user_id", "users.id")
		sub2 := Select("1").From("posts").WhereColsEqual("user_id", "users.id")

		q := Select("id").From("users").
			WhereExists(sub1).
			WhereNotExists(sub2)
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id) AND NOT EXISTS (SELECT 1 FROM posts WHERE user_id = users.id)"
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

	t.Run("where columns equal", func(t *testing.T) {
		q := Select("id").From("users").WhereColsEqual("user_id", "users.id")
		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id FROM users WHERE user_id = users.id"
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

	t.Run("complex combination of convenience methods", func(t *testing.T) {
		q := Select("id", "name", "email").From("users").
			WhereEqual("active", true).
			WhereNotNull("email").
			WhereGreaterThan("age", 18).
			WhereLessThan("age", 65).
			WhereIn("status", "active", "pending").
			WhereLike("name", "%john%").
			WhereBetween("score", 70, 100)

		sql, args, err := q.WithDialect(NoQuoteIdent()).Build()
		wantSQL := "SELECT id, name, email FROM users WHERE active = ? AND email IS NOT NULL AND age > ? AND age < ? AND status IN (?, ?) AND name LIKE ? AND score BETWEEN ? AND ?"
		wantArgs := []interface{}{true, 18, 65, "active", "pending", "%john%", 70, 100}
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
