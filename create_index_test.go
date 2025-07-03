package cqb

import (
	"testing"
)

func TestCreateIndexBuilder(t *testing.T) {
	t.Run("basic create index", func(t *testing.T) {
		q := CreateIndex("idx_users_email").
			On("users").
			Columns("email").
			WithDialect(MySQL())

		sql, args, err := q.Build()
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

	t.Run("create unique index", func(t *testing.T) {
		q := CreateIndex("idx_users_email_unique").
			On("users").
			Columns("email").
			Unique().
			WithDialect(MySQL())

		sql, args, err := q.Build()
		wantSQL := "CREATE UNIQUE INDEX `idx_users_email_unique` ON `users` (`email`)"

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

	t.Run("create index with multiple columns", func(t *testing.T) {
		q := CreateIndex("idx_users_name_email").
			On("users").
			Columns("name", "email").
			WithDialect(MySQL())

		sql, args, err := q.Build()
		wantSQL := "CREATE INDEX `idx_users_name_email` ON `users` (`name`, `email`)"

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

	t.Run("create index with if not exists", func(t *testing.T) {
		q := CreateIndex("idx_users_email").
			On("users").
			Columns("email").
			IfNotExists().
			WithDialect(MySQL())

		sql, args, err := q.Build()
		wantSQL := "CREATE INDEX IF NOT EXISTS `idx_users_email` ON `users` (`email`)"

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

	t.Run("create unique index with if not exists", func(t *testing.T) {
		q := CreateIndex("idx_users_email_unique").
			On("users").
			Columns("email").
			Unique().
			IfNotExists().
			WithDialect(MySQL())

		sql, args, err := q.Build()
		wantSQL := "CREATE UNIQUE INDEX IF NOT EXISTS `idx_users_email_unique` ON `users` (`email`)"

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

func TestCreateIndexBuilder_Errors(t *testing.T) {
	t.Run("empty index name", func(t *testing.T) {
		_, _, err := CreateIndex("").Build()
		if err == nil {
			t.Errorf("expected error for empty index name, got none")
		}
	})

	t.Run("no table name", func(t *testing.T) {
		_, _, err := CreateIndex("idx_test").Build()
		if err == nil {
			t.Errorf("expected error for no table name, got none")
		}
	})

	t.Run("no columns", func(t *testing.T) {
		_, _, err := CreateIndex("idx_test").On("users").Build()
		if err == nil {
			t.Errorf("expected error for no columns, got none")
		}
	})

	t.Run("empty table name", func(t *testing.T) {
		_, _, err := CreateIndex("idx_test").On("").Columns("email").Build()
		if err == nil {
			t.Errorf("expected error for empty table name, got none")
		}
	})
}

func TestCreateIndexBuilder_Dialect(t *testing.T) {
	t.Run("MySQL dialect", func(t *testing.T) {
		q := CreateIndex("idx_users_email").
			On("users").
			Columns("email").
			WithDialect(MySQL())

		sql, args, err := q.Build()
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
		q := CreateIndex("idx_users_email").
			On("users").
			Columns("email").
			WithDialect(Postgres())

		sql, args, err := q.Build()
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

	t.Run("Standard dialect", func(t *testing.T) {
		q := CreateIndex("idx_users_email").
			On("users").
			Columns("email").
			WithDialect(Standard())

		sql, args, err := q.Build()
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
		q := CreateIndex("idx_users_email").
			On("users").
			Columns("email").
			WithDialect(Postgres())

		sql, args, err := q.Build()
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
		q := CreateIndex("idx_users_email_unique").
			On("users").
			Columns("email").
			Unique().
			WithDialect(Postgres())

		sql, args, err := q.Build()
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
		q := CreateIndex("idx_users_email").
			On("users").
			Columns("email").
			IfNotExists().
			WithDialect(Postgres())

		sql, args, err := q.Build()
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

	t.Run("create index with multiple columns (postgres)", func(t *testing.T) {
		q := CreateIndex("idx_users_name_email").
			On("users").
			Columns("name", "email").
			WithDialect(Postgres())

		sql, args, err := q.Build()
		wantSQL := "CREATE INDEX \"idx_users_name_email\" ON \"users\" (\"name\", \"email\")"

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
