//go:build integration

package stk

import (
	"database/sql"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sprylic/stk/ddl"
)

func TestMySQLIntegration(t *testing.T) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/"
	}

	// TODO DO NOT COMMIT
	dsn = dsn + "?allowNativePasswords=true"

	// Connect to MySQL without specifying a test database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skipf("skipping: failed to connect to mysql: %v", err)
	}
	defer db.Close()

	suffix := func() string {
		rand.Seed(time.Now().UnixNano())
		letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
		b := make([]rune, 8)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}()

	testDBName := "stk_test_db_" + suffix

	// Create test database
	createDB := ddl.CreateDatabase(testDBName).IfNotExists().Charset("utf8mb4").Collation("utf8mb4_unicode_ci")
	sqlStr, _, err := createDB.WithDialect(MySQL()).Build()
	if err != nil {
		t.Fatalf("create database build: %v", err)
	}
	_, err = db.Exec(sqlStr)
	if err != nil && !strings.Contains(err.Error(), "exists") {
		t.Fatalf("create database exec: %v", err)
	}

	// Connect to the test database
	testDSN := dsn
	if idx := strings.LastIndex(testDSN, "/"); idx != -1 {
		// Check if there's already a database name after the last /
		afterSlash := testDSN[idx+1:]
		if afterSlash == "" {
			// No database name, just add the test database
			testDSN = testDSN + testDBName
		} else {
			// There's already a database name, replace it
			testDSN = testDSN[:idx+1] + testDBName
		}
	} else {
		testDSN += "/" + testDBName
	}
	testDB, err := sql.Open("mysql", testDSN)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer func() {
		testDB.Close()
		dropDB := ddl.DropDatabase(testDBName).IfExists()
		sqlStr, _, _ := dropDB.WithDialect(MySQL()).Build()
		_, _ = db.Exec(sqlStr)
	}()

	// Test DDL operations
	t.Run("DDL Operations", func(t *testing.T) {
		testMySQLDDL(t, testDB)
	})

	// Test basic CRUD operations
	t.Run("Basic CRUD", func(t *testing.T) {
		testMySQLCRUD(t, testDB)
	})

	// Test advanced features
	t.Run("Advanced Features", func(t *testing.T) {
		testMySQLAdvanced(t, testDB)
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
			AddColumn(ddl.Column("created_at").Type("TIMESTAMP").Default(Raw("CURRENT_TIMESTAMP")).NotNull()).
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
			t.Fatalf("create table exec: %v\n%s", err, sqlStr)
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
			t.Fatalf("create index exec: %v\n%s", err, sqlStr)
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
			t.Fatalf("create view exec: %v\n%s", err, sqlStr)
		}
	})

	// Test ALTER TABLE
	t.Run("Alter Table", func(t *testing.T) {
		q := ddl.AlterTable("users").
			AddColumn(ddl.Column("updated_at").Type("TIMESTAMP").Default(Raw("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")).NotNull()).
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
			t.Fatalf("alter table exec: %v\n%s", err, sqlStr)
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
			t.Fatalf("insert exec: %v\n%s", err, sqlStr)
		}
		insertID, _ := result.LastInsertId()

		// Select
		q2 := Select("id", "name", "email").From("users").Where(NewStringCondition("id = ?", insertID))
		sqlStr, args, err = q2.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("select build: %v", err)
		}
		var id int64
		var name, email string
		err = db.QueryRow(sqlStr, args...).Scan(&id, &name, &email)
		if err != nil {
			t.Fatalf("select query: %v\n%s", err, sqlStr)
		}
		if name != "Bob" {
			t.Errorf("expected name Bob, got %s", name)
		}

		// Update
		q3 := Update("users").Set("age", 26).Where(NewStringCondition("id = ?", insertID))
		sqlStr, args, err = q3.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("update build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("update exec: %v\n%s", err, sqlStr)
		}

		// Delete
		q4 := Delete("users").Where(NewStringCondition("id = ?", insertID))
		sqlStr, args, err = q4.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("delete build: %v", err)
		}
		_, err = db.Exec(sqlStr, args...)
		if err != nil {
			t.Fatalf("delete exec: %v\n%s", err, sqlStr)
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
		subq := Select(Raw("AVG(age)")).From("users")
		subSQL, _, _ := subq.WithDialect(MySQL()).Build()
		q := Select("name", "age").From("users").Where(Raw("age > (" + subSQL + ")"))
		sqlStr, args, err := q.WithDialect(MySQL()).Build()
		if err != nil {
			t.Fatalf("complex select build: %v", err)
		}

		rows, err := db.Query(sqlStr, args...)
		if err != nil {
			t.Fatalf("complex select query: %v\n%s", err, sqlStr)
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
