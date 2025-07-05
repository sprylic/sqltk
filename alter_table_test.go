package sqltk

import (
	"testing"

	"github.com/sprylic/sqltk/ddl"
)

func TestAlterTableBuilder(t *testing.T) {
	t.Run("add column", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").AddColumn(ddl.Column("age").Type("INT")).
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD COLUMN age INT"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop column", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").DropColumn("old_field").
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users DROP COLUMN old_field"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("rename column", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").RenameColumn("username", "user_name").
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users RENAME COLUMN username TO user_name"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("rename table", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").RenameTable("accounts").
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users RENAME TO accounts"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("modify column", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			ModifyColumn(ddl.Column("age").Type("BIGINT").NotNull()).
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users MODIFY COLUMN age BIGINT NOT NULL"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add constraint", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddConstraint(ddl.NewConstraint().Unique("idx_email", "email")).
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD CONSTRAINT idx_email UNIQUE (email)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop constraint", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			DropConstraint("idx_email").
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users DROP CONSTRAINT idx_email"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add index", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddIndex("idx_name", "name").
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD INDEX idx_name (name)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add multi-column index", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddIndex("idx_name_email", "name", "email").
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD INDEX idx_name_email (name, email)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop index", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			DropIndex("idx_name").
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users DROP INDEX idx_name"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("complex alter with all operations", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddColumn(ddl.Column("age").Type("INT")).
			ModifyColumn(ddl.Column("name").Type("VARCHAR").Size(100).NotNull()).
			DropColumn("old_field").
			AddConstraint(ddl.NewConstraint().Check("chk_age", "age >= 0")).
			AddIndex("idx_age", "age").
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD COLUMN age INT, MODIFY COLUMN name VARCHAR(100) NOT NULL, DROP COLUMN old_field, ADD CONSTRAINT chk_age CHECK (age >= 0), ADD INDEX idx_age (age)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add column (postgres)", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").AddColumn(ddl.Column("age").Type("INT")).
			WithDialect(Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" ADD COLUMN \"age\" INT"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop column (postgres)", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").DropColumn("old_field").
			WithDialect(Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" DROP COLUMN \"old_field\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("rename column (postgres)", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").RenameColumn("username", "user_name").
			WithDialect(Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" RENAME COLUMN \"username\" TO \"user_name\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("rename table (postgres)", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").RenameTable("accounts").
			WithDialect(Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" RENAME TO \"accounts\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("modify column (postgres)", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			ModifyColumn(ddl.Column("age").Type("BIGINT").NotNull()).
			WithDialect(Postgres()).Build()
		// Postgres uses ALTER COLUMN ... TYPE ...
		wantSQL := "ALTER TABLE \"users\" MODIFY COLUMN \"age\" BIGINT NOT NULL"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add constraint (postgres)", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddConstraint(ddl.NewConstraint().Unique("idx_email", "email")).
			WithDialect(Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" ADD CONSTRAINT \"idx_email\" UNIQUE (\"email\")"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop constraint (postgres)", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			DropConstraint("idx_email").
			WithDialect(Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" DROP CONSTRAINT \"idx_email\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add index (postgres)", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddIndex("idx_name", "name").
			WithDialect(Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" ADD INDEX \"idx_name\" (\"name\")"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop index (postgres)", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			DropIndex("idx_name").
			WithDialect(Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" DROP INDEX \"idx_name\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}
