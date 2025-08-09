//go:build integration

package sqltk

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sprylic/sqltk/pgfunc"
	"github.com/sprylic/sqltk/pgtypes"
	"github.com/sprylic/sqltk/raw"
	"github.com/sprylic/sqltk/sqldialect"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/sprylic/sqltk/ddl"
)

func TestPostgresIntegration(t *testing.T) {
	// Connect to the default database (usually 'postgres')
	defaultDSN := os.Getenv("POSTGRES_DSN")
	if defaultDSN == "" {
		defaultDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	defaultDB, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		t.Skipf("skipping: failed to connect to default postgres: %v", err)
	}
	defer defaultDB.Close()

	suffix := func() string {
		rand.Seed(time.Now().UnixNano())
		letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
		b := make([]rune, 8)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}()

	testDBName := "sqltk_test_db_" + suffix
	testSchema := "sqltk_test_schema_" + suffix

	// Create test database (PostgreSQL doesn't support IF NOT EXISTS for CREATE DATABASE)
	createDB := ddl.CreateDatabase(testDBName)
	sqlStr, _, err := createDB.WithDialect(sqldialect.Postgres()).Build()
	if err != nil {
		t.Fatalf("create database build: %v", err)
	}
	_, err = defaultDB.Exec(sqlStr)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("create database exec: %v", err)
	}

	// Connect to the test database
	// Parse the original DSN to extract username and password
	// Format: postgres://username:password@host:port/database?params
	dsnParts := strings.Split(defaultDSN, "@")
	if len(dsnParts) != 2 {
		t.Fatalf("invalid DSN format: %s", defaultDSN)
	}

	// Extract the protocol and credentials part
	protocolCreds := dsnParts[0]
	// Remove the protocol part (postgres://)
	creds := strings.TrimPrefix(protocolCreds, "postgres://")

	// Extract the host and params part
	hostParams := dsnParts[1]
	// Split on / to separate host:port from database and query params
	parts := strings.SplitN(hostParams, "/", 2)
	if len(parts) < 2 {
		t.Fatalf("invalid DSN format: %s", defaultDSN)
	}

	hostPort := parts[0]
	dbAndParams := parts[1]

	// Split database and query parameters
	dbParams := strings.SplitN(dbAndParams, "?", 2)
	var params string
	if len(dbParams) > 1 {
		params = dbParams[1]
	}

	// Reconstruct DSN with test database
	testDSN := fmt.Sprintf("postgres://%s@%s/%s", creds, hostPort, testDBName)
	if params != "" {
		testDSN += "?" + params
	}

	db, err := sql.Open("postgres", testDSN)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer func() {
		db.Close()
		// Drop test database (must disconnect all clients first)
		_, _ = defaultDB.Exec(fmt.Sprintf("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s' AND pid <> pg_backend_pid()", testDBName))
		dropDB := ddl.DropDatabase(testDBName).IfExists().Cascade()
		sqlStr, _, _ := dropDB.WithDialect(sqldialect.Postgres()).Build()
		_, _ = defaultDB.Exec(sqlStr)
	}()

	// Create test schema
	createSchema := ddl.CreateSchema(testSchema).IfNotExists()
	sqlStr, _, err = createSchema.WithDialect(sqldialect.Postgres()).Build()
	if err != nil {
		t.Fatalf("create schema build: %v", err)
	}
	_, err = db.Exec(sqlStr)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("create schema exec: %v", err)
	}
	defer func() {
		dropSchema := ddl.DropSchema(testSchema).IfExists().Cascade()
		sqlStr, _, _ := dropSchema.WithDialect(sqldialect.Postgres()).Build()
		_, _ = db.Exec(sqlStr)
	}()

	// Set search_path to the test schema
	_, err = db.Exec(fmt.Sprintf("SET search_path TO %s", testSchema))
	if err != nil {
		t.Fatalf("set search_path: %v", err)
	}

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

		sqlStr, args, err := q.WithDialect(sqldialect.Postgres()).Build()
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
		sqlStr, args, err := q.WithDialect(sqldialect.Postgres()).Build()
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
		sqlStr, args, err := q.WithDialect(sqldialect.Postgres()).Build()
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
			AddColumn(ddl.Column("updated_at").Type("TIMESTAMP").Default(pgfunc.Now())).
			AddConstraint(
				ddl.NewConstraint().Unique("idx_email_age", "email", "age"),
			)
		sqlStr, args, err := q.WithDialect(sqldialect.Postgres()).Build()
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
		// Use PGArray for PostgreSQL array type
		tags := []string{"admin", "verified"}

		pq := NewPostgresInsert("users")
		pq.InsertBuilder = pq.InsertBuilder.Columns("name", "email", "age", "data", "tags").
			Values("Alice", "alice@example.com", 30, pgtypes.PGJSON{V: jsonData}, pgtypes.PGArray{V: tags})
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
			t.Fatalf("insert query: %v\n%s", err, sqlStr)
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

		// TODO: Update to not use raw subquery or string condition
		// Complex query with join, subquery, and aggregation
		// Use raw SQL for subquery to avoid SelectBuilder type issues
		subq := raw.Raw("(SELECT AVG(amount) FROM orders)")
		q := Select("u.name", "COUNT(o.id) as order_count", "SUM(o.amount) as total_amount").
			From(Alias("users", "u")).
			LeftJoin("orders o").On("o.user_id", "u.id").
			Where(NewStringCondition("o.amount > " + string(subq))).
			GroupBy("u.name").
			OrderBy("total_amount DESC")

		sqlStr, args, err := q.WithDialect(sqldialect.Postgres()).Build()
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

	// Advanced UPDATE with RETURNING
	t.Run("Advanced Update", func(t *testing.T) {
		// Insert a test user first
		_, err := db.Exec(`INSERT INTO users (name, email, age) VALUES ('Bob', 'bob@example.com', 25)`)
		if err != nil {
			t.Fatalf("insert test user: %v", err)
		}

		// Update with RETURNING
		q := Update("users").
			Set("age", 26).
			Set("updated_at", raw.Raw("NOW()")).
			Where(NewStringCondition("name = ?", "Bob"))

		sqlStr, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		if err != nil {
			t.Fatalf("update build: %v", err)
		}

		// Add RETURNING clause manually for Postgres
		sqlStr += " RETURNING id, name, age, updated_at"

		var id int
		var name string
		var age int
		var updatedAt interface{}
		err = db.QueryRow(sqlStr, args...).Scan(&id, &name, &age, &updatedAt)
		if err != nil {
			t.Fatalf("update query: %v", err)
		}

		if name != "Bob" || age != 26 {
			t.Errorf("expected Bob with age 26, got %s with age %d", name, age)
		}
	})

	// DELETE with RETURNING
	t.Run("Delete with Returning", func(t *testing.T) {
		// Insert a test user first
		_, err := db.Exec(`INSERT INTO users (name, email, age) VALUES ('Charlie', 'charlie@example.com', 30)`)
		if err != nil {
			t.Fatalf("insert test user: %v", err)
		}

		// Delete with RETURNING
		q := Delete("users").Where(NewStringCondition("name = ?", "Charlie"))

		sqlStr, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		if err != nil {
			t.Fatalf("delete build: %v", err)
		}

		// Add RETURNING clause manually for Postgres
		sqlStr += " RETURNING id, name, email"

		var id int
		var name, email string
		err = db.QueryRow(sqlStr, args...).Scan(&id, &name, &email)
		if err != nil {
			t.Fatalf("delete query: %v", err)
		}

		if name != "Charlie" || email != "charlie@example.com" {
			t.Errorf("expected Charlie with email charlie@example.com, got %s with email %s", name, email)
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
		sqlStr, args, err := q2.WithDialect(sqldialect.Postgres()).Build()
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
