package cqb

import (
	"testing"

	"github.com/sprylic/cqb/ddl"
)

func init() {
	SetDialect(Standard())
}

func TestCreateViewBuilder(t *testing.T) {
	t.Run("basic create view", func(t *testing.T) {
		q := ddl.CreateView("active_users").
			As(Raw("SELECT id, name FROM users WHERE active = 1"))

		sql, args, err := q.WithDialect(Standard()).Build()
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
		q := ddl.CreateView("user_stats").
			OrReplace().
			As(Raw("SELECT COUNT(*) as total_users FROM users"))

		sql, args, err := q.WithDialect(Standard()).Build()
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
		q := ddl.CreateView("expensive_view").
			Materialized().
			As(Raw("SELECT u.id, u.name, COUNT(o.id) as order_count FROM users u LEFT JOIN orders o ON u.id = o.user_id GROUP BY u.id, u.name"))

		sql, args, err := q.WithDialect(Standard()).Build()
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
		q := ddl.CreateView("complex_stats").
			Materialized().
			OrReplace().
			As(Raw("SELECT * FROM complex_calculation_view"))

		sql, args, err := q.WithDialect(Standard()).Build()
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
		q := ddl.CreateView("user_order_summary").
			As(Raw("SELECT u.id, u.name, u.email, COUNT(o.id) as total_orders, SUM(o.total) as total_spent, AVG(o.total) as avg_order_value FROM users u LEFT JOIN orders o ON u.id = o.user_id WHERE u.active = 1 GROUP BY u.id, u.name, u.email HAVING COUNT(o.id) > 0"))

		sql, args, err := q.WithDialect(Standard()).Build()
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
		q := ddl.CreateView("builder_view").As(mockBuilder)
		sql, args, err := q.WithDialect(Standard()).Build()
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
		q := ddl.CreateView("").
			As("SELECT * FROM users")
		_, _, err := q.WithDialect(Standard()).Build()
		if err == nil {
			t.Errorf("expected error for empty view name, got none")
		}
	})

	t.Run("empty view definition", func(t *testing.T) {
		q := ddl.CreateView("test_view").
			As("")
		_, _, err := q.WithDialect(Standard()).Build()
		if err == nil {
			t.Errorf("expected error for empty view definition, got none")
		}
	})
}

func TestCreateViewBuilder_Dialect(t *testing.T) {
	t.Run("MySQL dialect", func(t *testing.T) {
		q := ddl.CreateView("active_users").
			As(Raw("SELECT id, name FROM users WHERE active = 1"))

		sql, args, err := q.WithDialect(MySQL()).Build()
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
		q := ddl.CreateView("active_users").
			As(Raw("SELECT id, name FROM users WHERE active = 1"))

		sql, args, err := q.WithDialect(Postgres()).Build()
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
		q := ddl.CreateView("user_stats").
			OrReplace().
			As(Raw("SELECT COUNT(*) as total_users FROM users"))

		sql, args, err := q.WithDialect(Postgres()).Build()
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
		q := ddl.CreateView("expensive_view").
			Materialized().
			As(Raw("SELECT u.id, u.name, COUNT(o.id) as order_count FROM users u LEFT JOIN orders o ON u.id = o.user_id GROUP BY u.id, u.name"))

		sql, args, err := q.WithDialect(Postgres()).Build()
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
		q := ddl.CreateView("test_view").
			As(Raw("SELECT * FROM users"))

		sql, args, err := q.WithDialect(Postgres()).Build()
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
