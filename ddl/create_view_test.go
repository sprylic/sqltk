package ddl

import (
	"testing"

	"github.com/sprylic/sqltk/raw"
	"github.com/sprylic/sqltk/sqldialect"
)

func init() {
	sqldialect.SetDialect(sqldialect.NoQuoteIdent())
}

func TestCreateViewBuilder(t *testing.T) {
	t.Run("basic create view", func(t *testing.T) {
		q := CreateView("active_users").
			As(raw.Raw("SELECT id, name FROM users WHERE active = 1"))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE VIEW active_users AS SELECT id, name FROM users WHERE active = 1"

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

	t.Run("create view with or replace", func(t *testing.T) {
		q := CreateView("user_stats").
			OrReplace().
			As(raw.Raw("SELECT COUNT(*) as total_users FROM users"))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE OR REPLACE VIEW user_stats AS SELECT COUNT(*) as total_users FROM users"

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

	t.Run("create materialized view", func(t *testing.T) {
		q := CreateView("expensive_view").
			Materialized().
			As(raw.Raw("SELECT u.id, u.name, COUNT(o.id) as order_count FROM users u LEFT JOIN orders o ON u.id = o.user_id GROUP BY u.id, u.name"))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE MATERIALIZED VIEW expensive_view AS SELECT u.id, u.name, COUNT(o.id) as order_count FROM users u LEFT JOIN orders o ON u.id = o.user_id GROUP BY u.id, u.name"

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

	t.Run("create materialized view with or replace", func(t *testing.T) {
		q := CreateView("complex_stats").
			Materialized().
			OrReplace().
			As(raw.Raw("SELECT * FROM complex_calculation_view"))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE OR REPLACE MATERIALIZED VIEW complex_stats AS SELECT * FROM complex_calculation_view"

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

	t.Run("complex view with joins and aggregations", func(t *testing.T) {
		q := CreateView("user_order_summary").
			As(raw.Raw("SELECT u.id, u.name, u.email, COUNT(o.id) as total_orders, SUM(o.total) as total_spent, AVG(o.total) as avg_order_value FROM users u LEFT JOIN orders o ON u.id = o.user_id WHERE u.active = 1 GROUP BY u.id, u.name, u.email HAVING COUNT(o.id) > 0"))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE VIEW user_order_summary AS SELECT u.id, u.name, u.email, COUNT(o.id) as total_orders, SUM(o.total) as total_spent, AVG(o.total) as avg_order_value FROM users u LEFT JOIN orders o ON u.id = o.user_id WHERE u.active = 1 GROUP BY u.id, u.name, u.email HAVING COUNT(o.id) > 0"

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

	t.Run("builder as view definition", func(t *testing.T) {
		mockBuilder := &mockSelectBuilder{sql: "SELECT id FROM users"}
		q := CreateView("builder_view").As(mockBuilder)
		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE VIEW builder_view AS SELECT id FROM users"
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

func TestCreateViewBuilder_Errors(t *testing.T) {
	t.Run("empty view name", func(t *testing.T) {
		q := CreateView("").
			As("SELECT * FROM users")
		_, _, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for empty view name, got none")
		}
	})

	t.Run("empty view definition", func(t *testing.T) {
		q := CreateView("test_view").
			As("")
		_, _, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for empty view definition, got none")
		}
	})
}

func TestCreateViewBuilder_Dialect(t *testing.T) {
	t.Run("MySQL dialect", func(t *testing.T) {
		q := CreateView("active_users").
			As(raw.Raw("SELECT id, name FROM users WHERE active = 1"))

		sql, args, err := q.WithDialect(sqldialect.MySQL()).Build()
		wantSQL := "CREATE VIEW `active_users` AS SELECT id, name FROM users WHERE active = 1"

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
		q := CreateView("active_users").
			As(raw.Raw("SELECT id, name FROM users WHERE active = 1"))

		sql, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE VIEW \"active_users\" AS SELECT id, name FROM users WHERE active = 1"

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

	t.Run("Postgres dialect with or replace", func(t *testing.T) {
		q := CreateView("user_stats").
			OrReplace().
			As(raw.Raw("SELECT COUNT(*) as total_users FROM users"))

		sql, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE OR REPLACE VIEW \"user_stats\" AS SELECT COUNT(*) as total_users FROM users"

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

	t.Run("Postgres dialect with materialized view", func(t *testing.T) {
		q := CreateView("expensive_view").
			Materialized().
			As(raw.Raw("SELECT u.id, u.name, COUNT(o.id) as order_count FROM users u LEFT JOIN orders o ON u.id = o.user_id GROUP BY u.id, u.name"))

		sql, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE MATERIALIZED VIEW \"expensive_view\" AS SELECT u.id, u.name, COUNT(o.id) as order_count FROM users u LEFT JOIN orders o ON u.id = o.user_id GROUP BY u.id, u.name"

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

func TestCreateViewBuilder_WithDialect(t *testing.T) {
	t.Run("explicit dialect override", func(t *testing.T) {
		q := CreateView("test_view").
			As(raw.Raw("SELECT * FROM users"))

		sql, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE VIEW \"test_view\" AS SELECT * FROM users"

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

type mockSelectBuilder struct {
	sql  string
	args []interface{}
	err  error
}

func (m *mockSelectBuilder) Build() (string, []interface{}, error) {
	return m.sql, m.args, m.err
}
