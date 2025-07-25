package ddl

import (
	"github.com/sprylic/sqltk"
	"testing"
)

func init() {
	sqltk.SetDialect(sqltk.NoQuoteIdent())
}

func TestTruncateTableBuilder(t *testing.T) {
	t.Run("basic truncate table", func(t *testing.T) {
		q := TruncateTable("users")

		sql, args, err := q.Build()
		wantSQL := "TRUNCATE TABLE users"

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

	t.Run("truncate multiple tables", func(t *testing.T) {
		q := TruncateTable("users", "orders", "products")

		sql, args, err := q.Build()
		wantSQL := "TRUNCATE TABLE users, orders, products"

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

	t.Run("truncate table with cascade", func(t *testing.T) {
		q := TruncateTable("users").
			Cascade()

		sql, args, err := q.Build()
		wantSQL := "TRUNCATE TABLE users CASCADE"

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

	t.Run("truncate table with restrict", func(t *testing.T) {
		q := TruncateTable("users").
			Restrict()

		sql, args, err := q.Build()
		wantSQL := "TRUNCATE TABLE users RESTRICT"

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

	t.Run("truncate table with dialect quoting", func(t *testing.T) {
		sql, _, err := TruncateTable("users").WithDialect(sqltk.MySQL()).Build()
		wantSQL := "TRUNCATE TABLE `users`"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("basic truncate table (postgres)", func(t *testing.T) {
		sql, args, err := TruncateTable("users").WithDialect(sqltk.Postgres()).Build()
		wantSQL := "TRUNCATE TABLE \"users\""
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

	t.Run("truncate table cascade (postgres)", func(t *testing.T) {
		sql, _, err := TruncateTable("users").Cascade().WithDialect(sqltk.Postgres()).Build()
		wantSQL := "TRUNCATE TABLE \"users\" CASCADE"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("truncate table restrict (postgres)", func(t *testing.T) {
		sql, _, err := TruncateTable("users").Restrict().WithDialect(sqltk.Postgres()).Build()
		wantSQL := "TRUNCATE TABLE \"users\" RESTRICT"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("truncate table restart identity (postgres)", func(t *testing.T) {
		sql, _, err := TruncateTable("users").Restart().WithDialect(sqltk.Postgres()).Build()
		wantSQL := "TRUNCATE TABLE \"users\" RESTART IDENTITY"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("truncate table continue identity (postgres)", func(t *testing.T) {
		sql, _, err := TruncateTable("users").Continue().WithDialect(sqltk.Postgres()).Build()
		wantSQL := "TRUNCATE TABLE \"users\" CONTINUE IDENTITY"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("truncate table restart identity cascade (postgres)", func(t *testing.T) {
		sql, _, err := TruncateTable("users").Restart().Cascade().WithDialect(sqltk.Postgres()).Build()
		wantSQL := "TRUNCATE TABLE \"users\" RESTART IDENTITY CASCADE"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("truncate table continue identity restrict (postgres)", func(t *testing.T) {
		sql, _, err := TruncateTable("users").Continue().Restrict().WithDialect(sqltk.Postgres()).Build()
		wantSQL := "TRUNCATE TABLE \"users\" CONTINUE IDENTITY RESTRICT"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("truncate multiple tables with cascade (postgres)", func(t *testing.T) {
		sql, _, err := TruncateTable("users", "orders").Cascade().WithDialect(sqltk.Postgres()).Build()
		wantSQL := "TRUNCATE TABLE \"users\", \"orders\" CASCADE"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("truncate multiple tables restart identity cascade (postgres)", func(t *testing.T) {
		sql, _, err := TruncateTable("users", "orders").Restart().Cascade().WithDialect(sqltk.Postgres()).Build()
		wantSQL := "TRUNCATE TABLE \"users\", \"orders\" RESTART IDENTITY CASCADE"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("error: no table names", func(t *testing.T) {
		q := TruncateTable()
		_, _, err := q.Build()
		if err == nil {
			t.Fatal("expected error for no table names")
		}
		if err.Error() != "at least one table name is required" {
			t.Errorf("got error %q, want %q", err.Error(), "at least one table name is required")
		}
	})

	t.Run("error: empty table name", func(t *testing.T) {
		q := TruncateTable("")
		_, _, err := q.Build()
		if err == nil {
			t.Fatal("expected error for empty table name")
		}
		if err.Error() != "table name cannot be empty" {
			t.Errorf("got error %q, want %q", err.Error(), "table name cannot be empty")
		}
	})

	t.Run("error: empty table name in multiple", func(t *testing.T) {
		q := TruncateTable("users", "", "orders")
		_, _, err := q.Build()
		if err == nil {
			t.Fatal("expected error for empty table name")
		}
		if err.Error() != "table name cannot be empty" {
			t.Errorf("got error %q, want %q", err.Error(), "table name cannot be empty")
		}
	})

	t.Run("cascade and restrict are mutually exclusive", func(t *testing.T) {
		q := TruncateTable("users").Cascade().Restrict()
		sql, _, err := q.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should only have RESTRICT since it was called last
		wantSQL := "TRUNCATE TABLE users RESTRICT"
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("restart and continue are mutually exclusive", func(t *testing.T) {
		q := TruncateTable("users").Restart().Continue().WithDialect(sqltk.Postgres())
		sql, _, err := q.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should only have CONTINUE since it was called last
		wantSQL := "TRUNCATE TABLE \"users\" CONTINUE IDENTITY"
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("restart and continue are ignored for non-postgres dialects", func(t *testing.T) {
		q := TruncateTable("users").Restart().Continue().WithDialect(sqltk.MySQL())
		sql, _, err := q.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should not include RESTART/CONTINUE for MySQL
		wantSQL := "TRUNCATE TABLE `users`"
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}
