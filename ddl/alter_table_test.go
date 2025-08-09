package ddl

import (
	"testing"

	"github.com/sprylic/sqltk/shared"
)

func TestAlterTableBuilder(t *testing.T) {
	t.Run("add column", func(t *testing.T) {
		sql, _, err := AlterTable("users").AddColumn(Column("age").Type("INT")).
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD COLUMN age INT"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop column", func(t *testing.T) {
		sql, _, err := AlterTable("users").DropColumn("old_field").
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users DROP COLUMN old_field"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("rename column", func(t *testing.T) {
		sql, _, err := AlterTable("users").RenameColumn("username", "user_name").
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users RENAME COLUMN username TO user_name"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("rename table", func(t *testing.T) {
		sql, _, err := AlterTable("users").RenameTable("accounts").
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users RENAME TO accounts"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("modify column", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			ModifyColumn(Column("age").Type("BIGINT").NotNull()).
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users MODIFY COLUMN age BIGINT NOT NULL"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add constraint", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			AddConstraint(NewConstraint().Unique("idx_email", "email")).
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD CONSTRAINT idx_email UNIQUE (email)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop constraint", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			DropConstraint("idx_email").
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users DROP CONSTRAINT idx_email"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add index", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			AddIndex("idx_name", "name").
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD INDEX idx_name (name)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add multi-column index", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			AddIndex("idx_name_email", "name", "email").
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD INDEX idx_name_email (name, email)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop index", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			DropIndex("idx_name").
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users DROP INDEX idx_name"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("complex alter with all operations", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			AddColumn(Column("age").Type("INT")).
			ModifyColumn(Column("name").Type("VARCHAR").Size(100).NotNull()).
			DropColumn("old_field").
			AddConstraint(NewConstraint().Check("chk_age", "age >= 0")).
			AddIndex("idx_age", "age").
			WithDialect(shared.NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD COLUMN age INT, MODIFY COLUMN name VARCHAR(100) NOT NULL, DROP COLUMN old_field, ADD CONSTRAINT chk_age CHECK (age >= 0), ADD INDEX idx_age (age)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add column (postgres)", func(t *testing.T) {
		sql, _, err := AlterTable("users").AddColumn(Column("age").Type("INT")).
			WithDialect(shared.Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" ADD COLUMN \"age\" INT"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop column (postgres)", func(t *testing.T) {
		sql, _, err := AlterTable("users").DropColumn("old_field").
			WithDialect(shared.Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" DROP COLUMN \"old_field\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("rename column (postgres)", func(t *testing.T) {
		sql, _, err := AlterTable("users").RenameColumn("username", "user_name").
			WithDialect(shared.Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" RENAME COLUMN \"username\" TO \"user_name\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("rename table (postgres)", func(t *testing.T) {
		sql, _, err := AlterTable("users").RenameTable("accounts").
			WithDialect(shared.Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" RENAME TO \"accounts\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("modify column (postgres)", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			ModifyColumn(Column("age").Type("BIGINT").NotNull()).
			WithDialect(shared.Postgres()).Build()
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
		sql, _, err := AlterTable("users").
			AddConstraint(NewConstraint().Unique("idx_email", "email")).
			WithDialect(shared.Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" ADD CONSTRAINT \"idx_email\" UNIQUE (\"email\")"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop constraint (postgres)", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			DropConstraint("idx_email").
			WithDialect(shared.Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" DROP CONSTRAINT \"idx_email\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add index (postgres)", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			AddIndex("idx_name", "name").
			WithDialect(shared.Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" ADD INDEX \"idx_name\" (\"name\")"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("drop index (postgres)", func(t *testing.T) {
		sql, _, err := AlterTable("users").
			DropIndex("idx_name").
			WithDialect(shared.Postgres()).Build()
		wantSQL := "ALTER TABLE \"users\" DROP INDEX \"idx_name\""
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}
