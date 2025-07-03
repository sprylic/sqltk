//go:build exclude

package main

import (
	"fmt"
	"log"

	"github.com/sprylic/stk/ddl"
	"github.com/sprylic/stk/shared"
)

func main() {
	fmt.Println("=== Database and Schema DDL Examples ===")

	// Set default dialect to MySQL
	shared.SetDialect(shared.MySQL())

	// Example 1: Create Database
	fmt.Println("1. CREATE DATABASE Examples:")

	// Basic database creation
	createDB := ddl.CreateDatabase("myapp_db")
	sql, _, err := createDB.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Basic: %s\n", sql)

	// Database with charset and collation
	createDBWithCharset := ddl.CreateDatabase("myapp_db").
		IfNotExists().
		Charset("utf8mb4").
		Collation("utf8mb4_unicode_ci")
	sql, _, err = createDBWithCharset.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   With charset: %s\n", sql)

	// Database with custom options
	createDBWithOptions := ddl.CreateDatabase("myapp_db").
		IfNotExists().
		Option("ENCRYPTION", "Y").
		Option("READ_ONLY", "0")
	sql, _, err = createDBWithOptions.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   With options: %s\n", sql)

	// Example 2: Drop Database
	fmt.Println("\n2. DROP DATABASE Examples:")

	// Basic drop
	dropDB := ddl.DropDatabase("myapp_db")
	sql, _, err = dropDB.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Basic: %s\n", sql)

	// Drop with if exists
	dropDBIfExists := ddl.DropDatabase("myapp_db").IfExists()
	sql, _, err = dropDBIfExists.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   With IF EXISTS: %s\n", sql)

	// Example 3: Create Schema (PostgreSQL style)
	fmt.Println("\n3. CREATE SCHEMA Examples (PostgreSQL):")

	// Set dialect to PostgreSQL
	shared.SetDialect(shared.Postgres())

	// Basic schema creation
	createSchema := ddl.CreateSchema("myapp_schema")
	sql, _, err = createSchema.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Basic: %s\n", sql)

	// Schema with authorization
	createSchemaWithAuth := ddl.CreateSchema("myapp_schema").
		IfNotExists().
		Authorization("myapp_user")
	sql, _, err = createSchemaWithAuth.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   With authorization: %s\n", sql)

	// Schema with custom options
	createSchemaWithOptions := ddl.CreateSchema("myapp_schema").
		IfNotExists().
		Option("QUOTA", "100MB")
	sql, _, err = createSchemaWithOptions.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   With options: %s\n", sql)

	// Example 4: Drop Schema
	fmt.Println("\n4. DROP SCHEMA Examples:")

	// Basic drop
	dropSchema := ddl.DropSchema("myapp_schema")
	sql, _, err = dropSchema.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Basic: %s\n", sql)

	// Drop with cascade
	dropSchemaCascade := ddl.DropSchema("myapp_schema").
		IfExists().
		Cascade()
	sql, _, err = dropSchemaCascade.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   With CASCADE: %s\n", sql)

	// Drop with restrict
	dropSchemaRestrict := ddl.DropSchema("myapp_schema").
		IfExists().
		Restrict()
	sql, _, err = dropSchemaRestrict.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   With RESTRICT: %s\n", sql)

	// Example 5: Complete workflow
	fmt.Println("\n5. Complete Database Setup Workflow:")

	// Switch back to MySQL for the workflow
	shared.SetDialect(shared.MySQL())

	// Step 1: Create database
	workflowDB := ddl.CreateDatabase("production_db").
		IfNotExists().
		Charset("utf8mb4").
		Collation("utf8mb4_unicode_ci")
	sql, _, err = workflowDB.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Step 1 - Create DB: %s\n", sql)

	// Step 2: Create schema (MySQL doesn't support CREATE SCHEMA like PostgreSQL)
	// In MySQL, schemas are equivalent to databases
	fmt.Printf("   Step 2 - Note: MySQL doesn't support CREATE SCHEMA\n")

	// Step 3: Create tables (example)
	fmt.Printf("   Step 3 - Create tables using existing DDL features\n")

	// Step 4: Cleanup
	cleanupDB := ddl.DropDatabase("production_db").IfExists()
	sql, _, err = cleanupDB.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Step 4 - Cleanup: %s\n", sql)

	fmt.Println()
	fmt.Println("=== End of Examples ===")
}
