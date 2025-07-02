//go:build integration

package cqb

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func TestPostgresIntegration(t *testing.T) {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/sq_postgres_db?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("skipping: failed to connect to postgres: %v", err)
	}
	defer db.Close()
	SetDialect(Postgres())

	_, _ = db.Exec(`DROP TABLE IF EXISTS users`)
	_, err = db.Exec(`CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT, age INT)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	// Insert
	q := Insert("users").Columns("name", "age").Values("Alice", 30).Values("Bob", 25)
	sqlStr, args, err := q.Build()
	if err != nil {
		t.Fatalf("insert build: %v", err)
	}
	_, err = db.Exec(sqlStr, args...)
	if err != nil {
		t.Fatalf("insert exec: %v", err)
	}

	// Select
	q2 := Select("id", "name", "age").From("users").Where("age > ?", 20)
	sqlStr, args, err = q2.Build()
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
		var id, age int
		var name string
		if err := rows.Scan(&id, &name, &age); err != nil {
			t.Fatalf("scan: %v", err)
		}
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 users, got %d", count)
	}

	// Update
	q3 := Update("users").Set("age", 31).Where("name = ?", "Alice")
	sqlStr, args, err = q3.Build()
	if err != nil {
		t.Fatalf("update build: %v", err)
	}
	_, err = db.Exec(sqlStr, args...)
	if err != nil {
		t.Fatalf("update exec: %v", err)
	}

	// Delete
	q4 := Delete("users").Where("name = ?", "Bob")
	sqlStr, args, err = q4.Build()
	if err != nil {
		t.Fatalf("delete build: %v", err)
	}
	_, err = db.Exec(sqlStr, args...)
	if err != nil {
		t.Fatalf("delete exec: %v", err)
	}
}

func TestMySQLIntegration(t *testing.T) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/sq_mysql_db"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skipf("skipping: failed to connect to mysql: %v", err)
	}
	defer db.Close()
	SetDialect(MySQL())

	_, _ = db.Exec(`DROP TABLE IF EXISTS users`)
	_, err = db.Exec(`CREATE TABLE users (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), age INT)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	// Insert
	q := Insert("users").Columns("name", "age").Values("Alice", 30).Values("Bob", 25)
	sqlStr, args, err := q.Build()
	if err != nil {
		t.Fatalf("insert build: %v", err)
	}
	_, err = db.Exec(sqlStr, args...)
	if err != nil {
		t.Fatalf("insert exec: %v", err)
	}

	// Select
	q2 := Select("id", "name", "age").From("users").Where("age > ?", 20)
	sqlStr, args, err = q2.Build()
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
		var id, age int
		var name string
		if err := rows.Scan(&id, &name, &age); err != nil {
			t.Fatalf("scan: %v", err)
		}
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 users, got %d", count)
	}

	// Update
	q3 := Update("users").Set("age", 31).Where("name = ?", "Alice")
	sqlStr, args, err = q3.Build()
	if err != nil {
		t.Fatalf("update build: %v", err)
	}
	_, err = db.Exec(sqlStr, args...)
	if err != nil {
		t.Fatalf("update exec: %v", err)
	}

	// Delete
	q4 := Delete("users").Where("name = ?", "Bob")
	sqlStr, args, err = q4.Build()
	if err != nil {
		t.Fatalf("delete build: %v", err)
	}
	_, err = db.Exec(sqlStr, args...)
	if err != nil {
		t.Fatalf("delete exec: %v", err)
	}
}
