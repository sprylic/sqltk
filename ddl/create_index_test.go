package ddl

import (
	"github.com/sprylic/sqltk/sqldialect"
	"testing"
)

func TestCreateIndexBuilder(t *testing.T) {
	t.Run("basic index", func(t *testing.T) {
		sql, _, err := CreateIndex("idx_users_email", "users").
			Columns("email").WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE INDEX idx_users_email ON users (email)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("unique index", func(t *testing.T) {
		sql, _, err := CreateIndex("idx_users_email_unique", "users").
			Unique().Columns("email").WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE UNIQUE INDEX idx_users_email_unique ON users (email)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("multi-column index", func(t *testing.T) {
		sql, _, err := CreateIndex("idx_users_name_email", "users").
			Columns("name", "email").WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE INDEX idx_users_name_email ON users (name, email)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		sql, _, err := CreateIndex("idx_users_email", "users").IfNotExists().
			Columns("email").WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE INDEX IF NOT EXISTS idx_users_email ON users (email)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("unique if not exists", func(t *testing.T) {
		sql, _, err := CreateIndex("idx_users_email_unique", "users").Unique().
			IfNotExists().Columns("email").WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users (email)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}

func TestCreateIndexBuilder_Errors(t *testing.T) {
	t.Run("empty index name", func(t *testing.T) {
		_, _, err := CreateIndex("", "users").Columns("email").
			WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for empty index name, got none")
		}
	})

	t.Run("no table name", func(t *testing.T) {
		_, _, err := CreateIndex("idx_test", "").Columns("email").
			WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for no table name, got none")
		}
	})

	t.Run("no columns", func(t *testing.T) {
		_, _, err := CreateIndex("idx_test", "users").Columns().
			WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for no columns, got none")
		}
	})

	t.Run("empty table name", func(t *testing.T) {
		_, _, err := CreateIndex("idx_test", "").Columns("email").
			WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for empty table name, got none")
		}
	})
}

func TestCreateIndexBuilder_Dialect(t *testing.T) {
	t.Run("MySQL dialect", func(t *testing.T) {
		sql, args, err := CreateIndex("idx_users_email", "users").
			Columns("email").WithDialect(sqldialect.MySQL()).Build()
		wantSQL := "CREATE INDEX `idx_users_email` ON `users` (`email`)"
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
		sql, args, err := CreateIndex("idx_users_email", "users").
			Columns("email").WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE INDEX \"idx_users_email\" ON \"users\" (\"email\")"
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

	t.Run("NoQuoteIdent dialect", func(t *testing.T) {
		sql, args, err := CreateIndex("idx_users_email", "users").
			Columns("email").WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE INDEX idx_users_email ON users (email)"
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

func TestCreateIndexBuilder_Postgres(t *testing.T) {
	t.Run("basic create index (postgres)", func(t *testing.T) {
		sql, args, err := CreateIndex("idx_users_email", "users").
			Columns("email").WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE INDEX \"idx_users_email\" ON \"users\" (\"email\")"
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

	t.Run("create unique index (postgres)", func(t *testing.T) {
		sql, args, err := CreateIndex("idx_users_email_unique", "users").
			Unique().Columns("email").WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE UNIQUE INDEX \"idx_users_email_unique\" ON \"users\" (\"email\")"
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

	t.Run("create index with if not exists (postgres)", func(t *testing.T) {
		sql, args, err := CreateIndex("idx_users_email", "users").IfNotExists().
			Columns("email").WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE INDEX IF NOT EXISTS \"idx_users_email\" ON \"users\" (\"email\")"
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

	t.Run("create unique index with if not exists (postgres)", func(t *testing.T) {
		sql, args, err := CreateIndex("idx_users_email_unique", "users").Unique().
			IfNotExists().Columns("email").WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE UNIQUE INDEX IF NOT EXISTS \"idx_users_email_unique\" ON \"users\" (\"email\")"
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
