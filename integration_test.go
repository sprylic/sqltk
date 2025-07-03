//go:build integration

package cqb

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/sprylic/cqb/ddl"
)

func TestPostgresIntegration(t *testing.T) {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres_test?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("skipping: failed to connect to postgres: %v", err)
	}
	defer db.Close()

	// Test DDL operations
	t.Run("DDL Operations", func(t *testing.T) {
		testPostgresDDL(t, db)
	})

	// Test basic CRUD operations
	t.Run("Basic CRUD", func(t *testing.T) {
		testPostgresCRUD(t, db)
	})

	// Test advanced features
	t.Run("Advanced Features", func(t *testing.T) {
		testPostgresAdvanced(t, db)
	})
}

func TestMySQLIntegration(t *testing.T) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/mysql_test"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skipf("skipping: failed to connect to mysql: %v", err)
	}
	defer db.Close()

	// Test DDL operations
	t.Run("DDL Operations", func(t *testing.T) {
		testMySQLDDL(t, db)
	})

	// Test basic CRUD operations
	t.Run("Basic CRUD", func(t *testing.T) {
		testMySQLCRUD(t, db)
	})

	// Test advanced features
	t.Run("Advanced Features", func(t *testing.T) {
		testMySQLAdvanced(t, db)
	})
}

func testPostgresDDL(t *testing.T, db *sql.DB) {
	// Clean up any existing tables
	_, _ = db.Exec(`DROP TABLE IF EXISTS orders CASCADE`)
	_, _ = db.Exec(`DROP TABLE IF EXISTS users CASCADE`)
	_, _ = db.Exec(`DROP VIEW IF EXISTS user_stats CASCADE`)
	_, _ = db.Exec(`DROP INDEX IF EXISTS idx_users_email`)

	// Test CREATE TABLE with all features
	t.Run("Create Table", func(t *testing.T) {
		q := ddl.CreateTable("users").
			AddColumn(ddl.Column("id").Type("SERIAL").NotNull()).
			AddColumn(ddl.Column("name").Type("VARCHAR").Size(255).NotNull()).
			AddColumn(ddl.Column("email").Type("VARCHAR").Size(255)).
			AddColumn(ddl.Column("age").Type("INTEGER")).
			AddColumn(ddl.Column("data").Type("JSONB")).
			AddColumn(ddl.Column("tags").Type("TEXT[]")).
			AddColumn(ddl.Column("created_at").Type("TIMESTAMP").Default("NOW()")).
			PrimaryKey("id").
			Unique("idx_email", "email").
			Check("chk_age", "age >= 0")

		sqlStr, args, err := q.WithDialect(Postgres()).Build()
		if err != nil {
			t.Fatalf("create table build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("create table exec: %v", err)
		}
	})

	// Test CREATE INDEX
	t.Run("Create Index", func(t *testing.T) {
		q := ddl.CreateIndex("idx_users_name", "users").Columns("name")
		sqlStr, args, err := q.WithDialect(Postgres()).Build()
		if err != nil {
			t.Fatalf("create index build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("create index exec: %v", err)
		}
	})

	// Test CREATE VIEW
	t.Run("Create View", func(t *testing.T) {
		subq := Select("name", "COUNT(*) as count").From("users").GroupBy("name")
		q := ddl.CreateView("user_stats").As(subq)
		sqlStr, args, err := q.WithDialect(Postgres()).Build()
		if err != nil {
			t.Fatalf("create view build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("create view exec: %v", err)
		}
	})

	// Test ALTER TABLE
	t.Run("Alter Table", func(t *testing.T) {
		q := ddl.AlterTable("users").
			AddColumn(ddl.Column("updated_at").Type("TIMESTAMP").Default("NOW()")).
			AddConstraint(ddl.Constraint{
				Type:    ddl.UniqueType,
				Name:    "idx_name_age",
				Columns: []string{"name", "age"},
			})
		sqlStr, args, err := q.WithDialect(Postgres()).Build()
		if err != nil {
			t.Fatalf("alter table build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("alter table exec: %v", err)
		}
	})
}

func testPostgresCRUD(t *testing.T, db *sql.DB) {
	// Insert with RETURNING and PGJSON/PGArray
	t.Run("Insert with Advanced Features", func(t *testing.T) {
		jsonData := map[string]interface{}{"preferences": map[string]string{"theme": "dark"}}
		tags := []string{"admin", "verified"}

		pq := NewPostgresInsert("users")
		pq.InsertBuilder = pq.InsertBuilder.Columns("name", "email", "age", "data", "tags").Values("Alice", "alice@example.com", 30, PGJSON{jsonData}, PGArray{tags})
		pq = pq.Returning("id", "name", "data")

		sqlStr, args, err := pq.Build()
		if err != nil {
			t.Fatalf("insert build: %v", err)
		}

		var id int
		var name string
		var dataBytes []byte
		err = db.QueryRow(sqlStr, args...).Scan(&id, &name, &dataBytes)
		if err != nil {
			t.Fatalf("insert query: %v", err)
		}

		if name != "Alice" {
			t.Errorf("expected name Alice, got %s", name)
		}

		// Verify JSON data
		var data map[string]interface{}
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}
		if data["preferences"].(map[string]interface{})["theme"] != "dark" {
			t.Errorf("expected theme dark, got %v", data["preferences"])
		}
	})

	// Advanced SELECT with joins, subqueries, and aggregations
	t.Run("Advanced Select", func(t *testing.T) {
		// Create orders table for join test
		_, _ = db.Exec(`DROP TABLE IF EXISTS orders`)
		_, err := db.Exec(`CREATE TABLE orders (id SERIAL PRIMARY KEY, user_id INTEGER, amount DECIMAL(10,2), created_at TIMESTAMP DEFAULT NOW())`)
		if err != nil {
			t.Fatalf("create orders table: %v", err)
		}

		// Insert some orders
		_, err = db.Exec(`INSERT INTO orders (user_id, amount) VALUES (1, 100.50), (1, 200.75), (2, 150.25)`)
		if err != nil {
			t.Fatalf("insert orders: %v", err)
		}

		// Complex query with join, subquery, and aggregation
		subq := Select("user_id").From("orders").GroupBy("user_id").Having("SUM(amount) > ?", 250)
		q := Select("u.name", "COUNT(o.id) as order_count", "SUM(o.amount) as total_amount").
			From("users u").
			LeftJoin("orders o").On("u.id", "o.user_id").
			Where("u.id IN (?)", subq).
			GroupBy("u.id").GroupBy("u.name").
			OrderBy("total_amount DESC").
			Limit(10)

		sqlStr, args, err := q.WithDialect(Postgres()).Build()
		if err != nil {
			t.Fatalf("select build: %v", err)
		}

		rows, err := db.Query(sqlStr, args...)
		if err != nil {
			t.Fatalf("select query: %v", err)
		}
		defer rows.Close()

		var count int
		for rows.Next() {
			var name string
			var orderCount int
			var totalAmount float64
			if err := rows.Scan(&name, &orderCount, &totalAmount); err != nil {
				t.Fatalf("scan: %v", err)
			}
			count++
		}
		if count == 0 {
			t.Error("expected at least one result from complex query")
		}
	})

	// Test UPDATE with complex WHERE conditions
	t.Run("Advanced Update", func(t *testing.T) {
		q := NewPostgresUpdate("users")
		q.UpdateBuilder = q.UpdateBuilder.Set("age", 31).Set("updated_at", Raw("NOW()")).Where("name = ?", "Alice").Where("age < ?", 35)
		q = q.Returning("id", "name", "age")

		sqlStr, args, err := q.Build()
		if err != nil {
			t.Fatalf("update build: %v", err)
		}

		var id int
		var name string
		var age int
		err = db.QueryRow(sqlStr, args...).Scan(&id, &name, &age)
		if err != nil {
			t.Fatalf("update query: %v", err)
		}

		if age != 31 {
			t.Errorf("expected age 31, got %d", age)
		}
	})

	// Test DELETE with RETURNING
	t.Run("Delete with Returning", func(t *testing.T) {
		q := NewPostgresDelete("users")
		q.DeleteBuilder = q.DeleteBuilder.Where("name = ?", "Alice")
		q = q.Returning("id", "name")
		sqlStr, args, err := q.Build()
		if err != nil {
			t.Fatalf("delete build: %v", err)
		}

		var id int
		var name string
		err = db.QueryRow(sqlStr, args...).Scan(&id, &name)
		if err != nil {
			t.Fatalf("delete query: %v", err)
		}

		if name != "Alice" {
			t.Errorf("expected name Alice, got %s", name)
		}
	})
}

func testPostgresAdvanced(t *testing.T, db *sql.DB) {
	// Test error handling with actual database operations
	t.Run("Error Handling", func(t *testing.T) {
		// Test invalid table name
		q := Select("name").From("")
		_, _, err := q.Build()
		if err == nil {
			t.Error("expected error for empty table name")
		}

		// Test querying non-existent column (should not error during build, but during execution)
		q2 := Select("invalid_column").From("users")
		sqlStr, args, err := q2.WithDialect(Postgres()).Build()
		if err != nil {
			t.Fatalf("unexpected error for valid select: %v", err)
		}

		// This should fail during execution, not during build
		_, err = db.Query(sqlStr, args...)
		if err == nil {
			t.Error("expected error when querying non-existent column")
		}
	})
}

func testMySQLDDL(t *testing.T, db *sql.DB) {
	// Clean up any existing tables
	_, _ = db.Exec(`DROP TABLE IF EXISTS orders`)
	_, _ = db.Exec(`DROP TABLE IF EXISTS users`)
	_, _ = db.Exec(`DROP VIEW IF EXISTS user_stats`)
	_, _ = db.Exec(`DROP INDEX IF EXISTS idx_users_email ON users`)

	// Test CREATE TABLE with MySQL-specific features
	t.Run("Create Table", func(t *testing.T) {
		q := ddl.CreateTable("users").
			AddColumn(ddl.Column("id").Type("INT").AutoIncrement().NotNull()).
			AddColumn(ddl.Column("name").Type("VARCHAR").Size(255).NotNull()).
			AddColumn(ddl.Column("email").Type("VARCHAR").Size(255)).
			AddColumn(ddl.Column("age").Type("INT")).
			AddColumn(ddl.Column("created_at").Type("TIMESTAMP").Default("CURRENT_TIMESTAMP")).
			PrimaryKey("id").
			Unique("idx_email", "email").
			Check("chk_age", "age >= 0").
			Engine("InnoDB").
			Charset("utf8mb4").
			Collation("utf8mb4_unicode_ci")

		sqlStr, args, err := q.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("create table build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("create table exec: %v", err)
		}
	})

	// Test CREATE INDEX
	t.Run("Create Index", func(t *testing.T) {
		q := ddl.CreateIndex("idx_users_name", "users").Columns("name")
		sqlStr, args, err := q.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("create index build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("create index exec: %v", err)
		}
	})

	// Test CREATE VIEW
	t.Run("Create View", func(t *testing.T) {
		subq := Select("name", "COUNT(*) as count").From("users").GroupBy("name")
		q := ddl.CreateView("user_stats").As(subq)
		sqlStr, args, err := q.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("create view build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("create view exec: %v", err)
		}
	})

	// Test ALTER TABLE
	t.Run("Alter Table", func(t *testing.T) {
		q := ddl.AlterTable("users").
			AddColumn(ddl.Column("updated_at").Type("TIMESTAMP").Default("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")).
			AddConstraint(ddl.Constraint{
				Type:    ddl.UniqueType,
				Name:    "idx_name_age",
				Columns: []string{"name", "age"},
			})
		sqlStr, args, err := q.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("alter table build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("alter table exec: %v", err)
		}
	})
}

func testMySQLCRUD(t *testing.T, db *sql.DB) {
	// Basic CRUD operations
	t.Run("Basic CRUD", func(t *testing.T) {
		// Insert
		q := Insert("users").Columns("name", "email", "age").Values("Bob", "bob@example.com", 25)
		sqlStr, args, err := q.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("insert build: %v", err)
		}
		result, err := db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("insert exec: %v", err)
		}
		insertID, _ := result.LastInsertId()

		// Select
		q2 := Select("id", "name", "email").From("users").Where("id = ?", insertID)
		sqlStr, args, err = q2.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("select build: %v", err)
		}
		var id int64
		var name, email string
		err = db.QueryRow(sqlStr, args...).Scan(&id, &name, &email)
		if err != nil {
			t.Fatalf("select query: %v", err)
		}
		if name != "Bob" {
			t.Errorf("expected name Bob, got %s", name)
		}

		// Update
		q3 := Update("users").Set("age", 26).Where("id = ?", insertID)
		sqlStr, args, err = q3.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("update build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("update exec: %v", err)
		}

		// Delete
		q4 := Delete("users").Where("id = ?", insertID)
		sqlStr, args, err = q4.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("delete build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("delete exec: %v", err)
		}
	})
}

func testMySQLAdvanced(t *testing.T, db *sql.DB) {
	// Test complex queries
	t.Run("Complex Queries", func(t *testing.T) {
		// Insert test data
		_, err := db.Exec(`INSERT INTO users (name, email, age) VALUES ('Alice', 'alice@example.com', 30), ('Charlie', 'charlie@example.com', 35)`)
		if err != nil {
			t.Fatalf("insert test data: %v", err)
		}

		// Complex query with subquery and aggregation
		subq := Select("AVG(age)").From("users")
		q := Select("name", "age").From("users").Where("age > (?)", subq)
		sqlStr, args, err := q.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("complex select build: %v", err)
		}

		rows, err := db.Query(sqlStr, args...)
		if err != nil {
			t.Fatalf("complex select query: %v", err)
		}
		defer rows.Close()

		var count int
		for rows.Next() {
			var name string
			var age int
			if err := rows.Scan(&name, &age); err != nil {
				t.Fatalf("scan: %v", err)
			}
			count++
		}
		if count == 0 {
			t.Error("expected at least one result from complex query")
		}
	})
}
