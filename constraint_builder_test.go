package sqltk

import (
	"testing"

	"github.com/sprylic/sqltk/ddl"
)

func TestConstraintBuilder(t *testing.T) {
	t.Run("NewPrimaryKey", func(t *testing.T) {
		constraint := ddl.NewConstraint().PrimaryKey("id").Build()
		if constraint.Type != ddl.PrimaryKeyType {
			t.Errorf("expected PrimaryKeyType, got %s", constraint.Type)
		}
		if len(constraint.Columns) != 1 || constraint.Columns[0] != "id" {
			t.Errorf("expected columns [id], got %v", constraint.Columns)
		}
	})

	t.Run("NewPrimaryKey multiple columns", func(t *testing.T) {
		constraint := ddl.NewConstraint().PrimaryKey("id", "tenant_id").Build()
		if constraint.Type != ddl.PrimaryKeyType {
			t.Errorf("expected PrimaryKeyType, got %s", constraint.Type)
		}
		if len(constraint.Columns) != 2 || constraint.Columns[0] != "id" || constraint.Columns[1] != "tenant_id" {
			t.Errorf("expected columns [id tenant_id], got %v", constraint.Columns)
		}
	})

	t.Run("NewUnique", func(t *testing.T) {
		constraint := ddl.NewConstraint().Unique("idx_email", "email").Build()
		if constraint.Type != ddl.UniqueType {
			t.Errorf("expected UniqueType, got %s", constraint.Type)
		}
		if constraint.Name != "idx_email" {
			t.Errorf("expected name idx_email, got %s", constraint.Name)
		}
		if len(constraint.Columns) != 1 || constraint.Columns[0] != "email" {
			t.Errorf("expected columns [email], got %v", constraint.Columns)
		}
	})

	t.Run("NewCheck", func(t *testing.T) {
		constraint := ddl.NewConstraint().Check("chk_age", "age >= 0").Build()
		if constraint.Type != ddl.CheckType {
			t.Errorf("expected CheckType, got %s", constraint.Type)
		}
		if constraint.Name != "chk_age" {
			t.Errorf("expected name chk_age, got %s", constraint.Name)
		}
		if constraint.CheckExpr != "age >= 0" {
			t.Errorf("expected check expression 'age >= 0', got %s", constraint.CheckExpr)
		}
	})

	t.Run("NewForeignKey", func(t *testing.T) {
		constraint := ddl.NewConstraint().ForeignKey("fk_user_role", "role_id").Build()
		if constraint.Type != ddl.ForeignKeyType {
			t.Errorf("expected ForeignKeyType, got %s", constraint.Type)
		}
		if constraint.Name != "fk_user_role" {
			t.Errorf("expected name fk_user_role, got %s", constraint.Name)
		}
		if len(constraint.Columns) != 1 || constraint.Columns[0] != "role_id" {
			t.Errorf("expected columns [role_id], got %v", constraint.Columns)
		}
	})

	t.Run("NewIndex", func(t *testing.T) {
		constraint := ddl.NewConstraint().Index("idx_name_email", "name", "email").Build()
		if constraint.Type != ddl.IndexType {
			t.Errorf("expected IndexType, got %s", constraint.Type)
		}
		if constraint.Name != "idx_name_email" {
			t.Errorf("expected name idx_name_email, got %s", constraint.Name)
		}
		if len(constraint.Columns) != 2 || constraint.Columns[0] != "name" || constraint.Columns[1] != "email" {
			t.Errorf("expected columns [name email], got %v", constraint.Columns)
		}
	})

	t.Run("NewRawConstraint", func(t *testing.T) {
		constraint := ddl.NewConstraint().Raw("chk_custom", "custom_expression").Build()
		if constraint.Type != ddl.CheckType {
			t.Errorf("expected CheckType, got %s", constraint.Type)
		}
		if constraint.Name != "chk_custom" {
			t.Errorf("expected name chk_custom, got %s", constraint.Name)
		}
		if constraint.CheckExpr != "custom_expression" {
			t.Errorf("expected check expression 'custom_expression', got %s", constraint.CheckExpr)
		}
	})

	t.Run("WithColumns", func(t *testing.T) {
		constraint := ddl.NewConstraint().Unique("idx_test", "col1").WithColumns("col1", "col2").Build()
		if len(constraint.Columns) != 2 || constraint.Columns[0] != "col1" || constraint.Columns[1] != "col2" {
			t.Errorf("expected columns [col1 col2], got %v", constraint.Columns)
		}
	})

	t.Run("WithCheckExpr", func(t *testing.T) {
		constraint := ddl.NewConstraint().Check("chk_test", "value > 0").WithCheckExpr("value > 0").Build()
		if constraint.CheckExpr != "value > 0" {
			t.Errorf("expected check expression 'value > 0', got %s", constraint.CheckExpr)
		}
	})

	t.Run("WithReference", func(t *testing.T) {
		constraint := ddl.NewConstraint().ForeignKey("fk_test", "user_id").
			WithReference("users", "id").
			Build()
		if constraint.Reference == nil {
			t.Fatal("expected reference to be set")
		}
		if constraint.Reference.Table != "users" {
			t.Errorf("expected table 'users', got %s", constraint.Reference.Table)
		}
		if len(constraint.Reference.Columns) != 1 || constraint.Reference.Columns[0] != "id" {
			t.Errorf("expected reference columns [id], got %v", constraint.Reference.Columns)
		}
	})

	t.Run("WithOnDelete", func(t *testing.T) {
		constraint := ddl.NewConstraint().ForeignKey("fk_test", "user_id").
			WithOnDelete("CASCADE").
			Build()
		if constraint.Reference == nil {
			t.Fatal("expected reference to be set")
		}
		if constraint.Reference.OnDelete != "CASCADE" {
			t.Errorf("expected ON DELETE CASCADE, got %s", constraint.Reference.OnDelete)
		}
	})

	t.Run("WithOnUpdate", func(t *testing.T) {
		constraint := ddl.NewConstraint().ForeignKey("fk_test", "user_id").
			WithOnUpdate("CASCADE").
			Build()
		if constraint.Reference == nil {
			t.Fatal("expected reference to be set")
		}
		if constraint.Reference.OnUpdate != "CASCADE" {
			t.Errorf("expected ON UPDATE CASCADE, got %s", constraint.Reference.OnUpdate)
		}
	})

	t.Run("chained methods", func(t *testing.T) {
		constraint := ddl.NewConstraint().ForeignKey("fk_user_role", "role_id").
			WithReference("roles", "id").
			WithOnDelete("CASCADE").
			WithOnUpdate("CASCADE").
			Build()

		if constraint.Type != ddl.ForeignKeyType {
			t.Errorf("expected ForeignKeyType, got %s", constraint.Type)
		}
		if constraint.Name != "fk_user_role" {
			t.Errorf("expected name fk_user_role, got %s", constraint.Name)
		}
		if len(constraint.Columns) != 1 || constraint.Columns[0] != "role_id" {
			t.Errorf("expected columns [role_id], got %v", constraint.Columns)
		}
		if constraint.Reference == nil {
			t.Fatal("expected reference to be set")
		}
		if constraint.Reference.Table != "roles" {
			t.Errorf("expected table 'roles', got %s", constraint.Reference.Table)
		}
		if len(constraint.Reference.Columns) != 1 || constraint.Reference.Columns[0] != "id" {
			t.Errorf("expected reference columns [id], got %v", constraint.Reference.Columns)
		}
		if constraint.Reference.OnDelete != "CASCADE" {
			t.Errorf("expected ON DELETE CASCADE, got %s", constraint.Reference.OnDelete)
		}
		if constraint.Reference.OnUpdate != "CASCADE" {
			t.Errorf("expected ON UPDATE CASCADE, got %s", constraint.Reference.OnUpdate)
		}
	})
}

func TestConstraintBuilderWithAlterTable(t *testing.T) {
	t.Run("add check constraint with builder", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddConstraintBuilder(ddl.NewConstraint().Check("chk_age", "age >= 0")).
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD CONSTRAINT chk_age CHECK (age >= 0)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add unique constraint with builder", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddConstraintBuilder(ddl.NewConstraint().Unique("idx_email", "email")).
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD CONSTRAINT idx_email UNIQUE (email)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add foreign key constraint with builder", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddConstraintBuilder(ddl.NewConstraint().ForeignKey("fk_user_role", "role_id").
				WithReference("roles", "id").
				WithOnDelete("CASCADE")).
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add foreign key constraint with ForeignKeyBuilder", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddForeignKey(
				ddl.ForeignKey("fk_user_role", "role_id").
					References("roles", "id").
					OnDelete("CASCADE"),
			).
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("add raw constraint", func(t *testing.T) {
		sql, _, err := ddl.AlterTable("users").
			AddConstraintBuilder(ddl.NewConstraint().Raw("chk_custom", "custom_expression")).
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "ALTER TABLE users ADD CONSTRAINT chk_custom CHECK (custom_expression)"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}

func TestConstraintBuilderWithCreateTable(t *testing.T) {
	t.Run("create table with check constraint using builder", func(t *testing.T) {
		sql, _, err := ddl.CreateTable("users").
			AddColumnWithType("id", "INT").
			AddColumnWithType("age", "INT").
			Check("chk_age", "age >= 0").
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT, age INT, CONSTRAINT chk_age CHECK (age >= 0))"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})

	t.Run("create table with unique constraint using builder", func(t *testing.T) {
		sql, _, err := ddl.CreateTable("users").
			AddColumnWithType("id", "INT").
			AddColumnWithType("email", "VARCHAR").
			Unique("idx_email", "email").
			WithDialect(NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT, email VARCHAR, CONSTRAINT idx_email UNIQUE (email))"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
	})
}
