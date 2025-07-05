package sqltk

import (
	"testing"

	"github.com/sprylic/sqltk/mysqlfunc"
	"github.com/sprylic/sqltk/pgfunc"
)

func init() {
	SetDialect(NoQuoteIdent())
}

func TestMySQLFunctions(t *testing.T) {
	t.Run("basic mysql functions", func(t *testing.T) {
		// Test date/time functions
		q := Select(mysqlfunc.CurrentTimestamp()).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT CURRENT_TIMESTAMP FROM users"
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

	t.Run("mysql string functions", func(t *testing.T) {
		// Test CONCAT function
		q := Select(Alias(mysqlfunc.Concat("first_name", "' '", "last_name"), "full_name")).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT CONCAT(first_name, ' ', last_name) AS full_name FROM users"
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

	t.Run("mysql aggregate functions", func(t *testing.T) {
		// Test COUNT function
		q := Select(Alias(mysqlfunc.Count("*"), "total_users")).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT COUNT(*) AS total_users FROM users"
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

	t.Run("mysql date formatting", func(t *testing.T) {
		// Test DATE_FORMAT function
		q := Select(Alias(mysqlfunc.DateFormat("created_at", "'%Y-%m-%d'"), "created_date")).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT DATE_FORMAT(created_at, '%Y-%m-%d') AS created_date FROM users"
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

	t.Run("mysql conditional functions", func(t *testing.T) {
		// Test IF function
		q := Select(Alias(mysqlfunc.If("active", "'Active'", "'Inactive'"), "status")).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT IF(active, 'Active', 'Inactive') AS status FROM users"
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

	t.Run("mysql numeric functions", func(t *testing.T) {
		// Test ROUND function
		q := Select(Alias(mysqlfunc.Round("price", 2), "rounded_price")).From("products")
		sql, args, err := q.Build()
		wantSQL := "SELECT ROUND(price, 2) AS rounded_price FROM products"
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

	t.Run("mysql complex query", func(t *testing.T) {
		// Test multiple functions in one query
		q := Select(
			Alias(mysqlfunc.Concat("first_name", "' '", "last_name"), "full_name"),
			Alias(mysqlfunc.DateFormat("created_at", "'%Y-%m-%d'"), "created_date"),
			Alias(mysqlfunc.Count("*"), "total"),
			Alias(mysqlfunc.If("active", "'Active'", "'Inactive'"), "status"),
		).From("users").GroupBy("first_name", "last_name")

		sql, args, err := q.Build()
		wantSQL := "SELECT CONCAT(first_name, ' ', last_name) AS full_name, DATE_FORMAT(created_at, '%Y-%m-%d') AS created_date, COUNT(*) AS total, IF(active, 'Active', 'Inactive') AS status FROM users GROUP BY first_name, last_name"
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

func TestPostgreSQLFunctions(t *testing.T) {
	t.Run("basic postgres functions", func(t *testing.T) {
		// Test date/time functions
		q := Select(pgfunc.Now()).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT now() FROM users"
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

	t.Run("postgres string functions", func(t *testing.T) {
		// Test CONCAT function
		q := Select(Alias(pgfunc.Concat("first_name", "' '", "last_name"), "full_name")).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT concat(first_name, ' ', last_name) AS full_name FROM users"
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

	t.Run("postgres aggregate functions", func(t *testing.T) {
		// Test COUNT function
		q := Select(Alias(pgfunc.Count("*"), "total_users")).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT count(*) AS total_users FROM users"
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

	t.Run("postgres date functions", func(t *testing.T) {
		// Test EXTRACT function
		q := Select(Alias(pgfunc.Extract("year", "created_at"), "created_year")).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT EXTRACT(year FROM created_at) AS created_year FROM users"
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

	t.Run("postgres conditional functions", func(t *testing.T) {
		// Test COALESCE function
		q := Select(Alias(pgfunc.Coalesce("nickname", "first_name"), "display_name")).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT coalesce(nickname, first_name) AS display_name FROM users"
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

	t.Run("postgres array functions", func(t *testing.T) {
		// Test ARRAY_AGG function
		q := Select(Alias(pgfunc.ArrayAgg("tag"), "all_tags")).From("posts").GroupBy("category")
		sql, args, err := q.Build()
		wantSQL := "SELECT array_agg(tag) AS all_tags FROM posts GROUP BY category"
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

	t.Run("postgres json functions", func(t *testing.T) {
		// Test JSON functions
		q := Select(Alias(pgfunc.JsonExtract("data", "'name'"), "user_name")).From("users")
		sql, args, err := q.Build()
		wantSQL := "SELECT json_extract_path_text(data, 'name') AS user_name FROM users"
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

	t.Run("postgres text search functions", func(t *testing.T) {
		// Test text search functions
		q := Select(Alias(pgfunc.ToTsvector("english", "content"), "search_vector")).From("articles")
		sql, args, err := q.Build()
		wantSQL := "SELECT to_tsvector(english, content) AS search_vector FROM articles"
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

	t.Run("postgres complex query", func(t *testing.T) {
		// Test multiple functions in one query
		q := Select(
			Alias(pgfunc.Concat("first_name", "' '", "last_name"), "full_name"),
			Alias(pgfunc.Extract("year", "created_at"), "created_year"),
			Alias(pgfunc.Count("*"), "total"),
			Alias(pgfunc.Coalesce("nickname", "first_name"), "display_name"),
		).From("users").GroupBy("first_name").GroupBy("last_name")

		sql, args, err := q.Build()
		wantSQL := "SELECT concat(first_name, ' ', last_name) AS full_name, EXTRACT(year FROM created_at) AS created_year, count(*) AS total, coalesce(nickname, first_name) AS display_name FROM users GROUP BY first_name, last_name"
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

func TestFunctionComparison(t *testing.T) {
	t.Run("compare mysql vs postgres functions", func(t *testing.T) {
		// Test that both MySQL and PostgreSQL functions work correctly

		// MySQL version
		mysqlQuery := Select(
			Alias(mysqlfunc.Concat("first_name", "' '", "last_name"), "full_name"),
			Alias(mysqlfunc.Count("*"), "total"),
		).From("users")

		mysqlSQL, mysqlArgs, mysqlErr := mysqlQuery.Build()
		if mysqlErr != nil {
			t.Fatalf("MySQL query error: %v", mysqlErr)
		}

		// PostgreSQL version
		postgresQuery := Select(
			Alias(pgfunc.Concat("first_name", "' '", "last_name"), "full_name"),
			Alias(pgfunc.Count("*"), "total"),
		).From("users")

		postgresSQL, postgresArgs, postgresErr := postgresQuery.Build()
		if postgresErr != nil {
			t.Fatalf("PostgreSQL query error: %v", postgresErr)
		}

		// Both should have no arguments
		if len(mysqlArgs) != 0 {
			t.Errorf("MySQL query has args %v, want none", mysqlArgs)
		}
		if len(postgresArgs) != 0 {
			t.Errorf("PostgreSQL query has args %v, want none", postgresArgs)
		}

		// Both should generate valid SQL
		if mysqlSQL == "" {
			t.Error("MySQL SQL is empty")
		}
		if postgresSQL == "" {
			t.Error("PostgreSQL SQL is empty")
		}

		// They should be different (different function names)
		if mysqlSQL == postgresSQL {
			t.Error("MySQL and PostgreSQL SQL should be different")
		}
	})
}

func TestFunctionWithWhereClause(t *testing.T) {
	t.Run("mysql function in where clause", func(t *testing.T) {
		// Test using MySQL function in WHERE clause
		q := Select("id", "name").From("users").Where(
			AsCondition(Raw("created_at > " + string(mysqlfunc.CurrentTimestamp()))),
		)

		sql, args, err := q.Build()
		wantSQL := "SELECT id, name FROM users WHERE created_at > CURRENT_TIMESTAMP"

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

	t.Run("postgres function in where clause", func(t *testing.T) {
		// Test using PostgreSQL function in WHERE clause
		q := Select("id", "name").From("users").Where(
			AsCondition(Raw("created_at > " + string(pgfunc.Now()))),
		)

		sql, args, err := q.Build()
		wantSQL := "SELECT id, name FROM users WHERE created_at > now()"

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

func TestFunctionInOrderBy(t *testing.T) {
	t.Run("mysql function in order by", func(t *testing.T) {
		// Test using MySQL function in ORDER BY
		q := Select("id", "name").From("users").OrderBy(mysqlfunc.Random())

		sql, args, err := q.Build()
		wantSQL := "SELECT id, name FROM users ORDER BY RAND()"

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

	t.Run("postgres function in order by", func(t *testing.T) {
		// Test using PostgreSQL function in ORDER BY
		q := Select("id", "name").From("users").OrderBy(pgfunc.Random())

		sql, args, err := q.Build()
		wantSQL := "SELECT id, name FROM users ORDER BY random()"

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

func TestFunctionInGroupBy(t *testing.T) {
	t.Run("mysql function in group by", func(t *testing.T) {
		// Test using MySQL function in GROUP BY
		q := Select("id").From("users").GroupBy(mysqlfunc.Year("created_at"))

		sql, args, err := q.Build()
		wantSQL := "SELECT id FROM users GROUP BY YEAR(created_at)"

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

	t.Run("postgres function in group by", func(t *testing.T) {
		// Test using PostgreSQL function in GROUP BY
		q := Select("id").From("users").GroupBy(pgfunc.Extract("year", "created_at"))

		sql, args, err := q.Build()
		wantSQL := "SELECT id FROM users GROUP BY EXTRACT(year FROM created_at)"

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
