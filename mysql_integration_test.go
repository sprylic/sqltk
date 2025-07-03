//go:build integration

package stk

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sprylic/stk/ddl"
)

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
