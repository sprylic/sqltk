package ddl

import (
	"github.com/sprylic/sqltk"
	"testing"
)

func TestDropViewBuilder(t *testing.T) {
	t.Run("basic drop view", func(t *testing.T) {
		q := DropView("active_users")

		sql, args, err := q.WithDialect(sqltk.NoQuoteIdent()).Build()
		wantSQL := "DROP VIEW active_users"

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

	t.Run("drop view with if exists", func(t *testing.T) {
		q := DropView("user_stats").
			IfExists()

		sql, args, err := q.WithDialect(sqltk.NoQuoteIdent()).Build()
		wantSQL := "DROP VIEW IF EXISTS user_stats"

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

	t.Run("drop view with cascade", func(t *testing.T) {
		q := DropView("complex_view").
			Cascade()

		sql, args, err := q.WithDialect(sqltk.NoQuoteIdent()).Build()
		wantSQL := "DROP VIEW complex_view CASCADE"

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

	t.Run("drop view with restrict", func(t *testing.T) {
		q := DropView("important_view").
			Restrict()

		sql, args, err := q.WithDialect(sqltk.NoQuoteIdent()).Build()
		wantSQL := "DROP VIEW important_view RESTRICT"

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

	t.Run("drop view with if exists and cascade", func(t *testing.T) {
		q := DropView("temp_view").
			IfExists().
			Cascade()

		sql, args, err := q.WithDialect(sqltk.NoQuoteIdent()).Build()
		wantSQL := "DROP VIEW IF EXISTS temp_view CASCADE"

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

	t.Run("drop view with if exists and restrict", func(t *testing.T) {
		q := DropView("protected_view").
			IfExists().
			Restrict()

		sql, args, err := q.WithDialect(sqltk.NoQuoteIdent()).Build()
		wantSQL := "DROP VIEW IF EXISTS protected_view RESTRICT"

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
}

func TestDropViewBuilder_Errors(t *testing.T) {
	t.Run("empty view name", func(t *testing.T) {
		q := DropView("")
		_, _, err := q.WithDialect(sqltk.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for empty view name, got none")
		}
	})

	t.Run("cascade and restrict together", func(t *testing.T) {
		q := DropView("test_view").
			Cascade().
			Restrict()
		_, _, err := q.WithDialect(sqltk.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for using both CASCADE and RESTRICT, got none")
		}
	})

	t.Run("restrict and cascade together", func(t *testing.T) {
		q := DropView("test_view").
			Restrict().
			Cascade()
		_, _, err := q.WithDialect(sqltk.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for using both RESTRICT and CASCADE, got none")
		}
	})
}

func TestDropViewBuilder_Dialect(t *testing.T) {
	t.Run("MySQL dialect", func(t *testing.T) {
		q := DropView("active_users")

		sql, args, err := q.WithDialect(sqltk.MySQL()).Build()
		wantSQL := "DROP VIEW `active_users`"

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

	t.Run("Postgres dialect", func(t *testing.T) {
		q := DropView("active_users")

		sql, args, err := q.WithDialect(sqltk.Postgres()).Build()
		wantSQL := "DROP VIEW \"active_users\""

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

	t.Run("Postgres dialect with if exists", func(t *testing.T) {
		q := DropView("user_stats").
			IfExists()

		sql, args, err := q.WithDialect(sqltk.Postgres()).Build()
		wantSQL := "DROP VIEW IF EXISTS \"user_stats\""

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

	t.Run("Postgres dialect with cascade", func(t *testing.T) {
		q := DropView("complex_view").
			Cascade()

		sql, args, err := q.WithDialect(sqltk.Postgres()).Build()
		wantSQL := "DROP VIEW \"complex_view\" CASCADE"

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

	t.Run("Postgres dialect with restrict", func(t *testing.T) {
		q := DropView("important_view").
			Restrict()

		sql, args, err := q.WithDialect(sqltk.Postgres()).Build()
		wantSQL := "DROP VIEW \"important_view\" RESTRICT"

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
}
