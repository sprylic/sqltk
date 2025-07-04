package sqltk

import (
	"testing"

	"github.com/sprylic/sqltk/ddl"
)

func init() {
	SetDialect(Standard())
}

func TestDropTableBuilder(t *testing.T) {
	t.Run("basic drop table", func(t *testing.T) {
		q := ddl.DropTable("users")

		sql, args, err := q.Build()
		wantSQL := "DROP TABLE users"

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

	t.Run("drop table with if exists", func(t *testing.T) {
		q := ddl.DropTable("users").
			IfExists()

		sql, args, err := q.Build()
		wantSQL := "DROP TABLE IF EXISTS users"

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

	t.Run("drop table with cascade", func(t *testing.T) {
		q := ddl.DropTable("users").
			Cascade()

		sql, args, err := q.Build()
		wantSQL := "DROP TABLE users CASCADE"

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

	t.Run("drop table with restrict", func(t *testing.T) {
		q := ddl.DropTable("users").
			Restrict()

		sql, args, err := q.Build()
		wantSQL := "DROP TABLE users RESTRICT"

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

	t.Run("drop table with if exists and cascade", func(t *testing.T) {
		q := ddl.DropTable("users").
			IfExists().
			Cascade()

		sql, args, err := q.Build()
		wantSQL := "DROP TABLE IF EXISTS users CASCADE"

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

	t.Run("drop table with dialect quoting", func(t *testing.T) {
		sql, _, err := ddl.DropTable("users").WithDialect(MySQL()).Build()
		wantSQL := "DROP TABLE `users`"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("basic drop table (postgres)", func(t *testing.T) {
		sql, args, err := ddl.DropTable("users").WithDialect(Postgres()).Build()
		wantSQL := "DROP TABLE \"users\""
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

	t.Run("drop table if exists (postgres)", func(t *testing.T) {
		sql, _, err := ddl.DropTable("users").IfExists().WithDialect(Postgres()).Build()
		wantSQL := "DROP TABLE IF EXISTS \"users\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop table cascade (postgres)", func(t *testing.T) {
		sql, _, err := ddl.DropTable("users").Cascade().WithDialect(Postgres()).Build()
		wantSQL := "DROP TABLE \"users\" CASCADE"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop table restrict (postgres)", func(t *testing.T) {
		sql, _, err := ddl.DropTable("users").Restrict().WithDialect(Postgres()).Build()
		wantSQL := "DROP TABLE \"users\" RESTRICT"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop table with dialect quoting (postgres)", func(t *testing.T) {
		sql, _, err := ddl.DropTable("users").WithDialect(Postgres()).Build()
		wantSQL := "DROP TABLE \"users\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}
