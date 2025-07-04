package main

import (
	"fmt"
	"log"

	"github.com/sprylic/stk"
	"github.com/sprylic/stk/ddl"
)

func main() {
	// Set dialect (optional, MySQL is default)
	stk.SetDialect(stk.Standard())

	fmt.Println("=== TRUNCATE TABLE Examples ===")
	fmt.Println()

	// Basic truncate table
	truncateBasic := ddl.TruncateTable("users")
	sql, args, err := truncateBasic.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Basic truncate: %s\n", sql)
	fmt.Printf("Args: %v\n\n", args)

	// Truncate multiple tables
	truncateMultiple := ddl.TruncateTable("users", "orders", "products")
	sql, args, err = truncateMultiple.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Multiple tables: %s\n", sql)
	fmt.Printf("Args: %v\n\n", args)

	// Truncate with cascade
	truncateCascade := ddl.TruncateTable("users").Cascade()
	sql, args, err = truncateCascade.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("With cascade: %s\n", sql)
	fmt.Printf("Args: %v\n\n", args)

	// Truncate with restrict
	truncateRestrict := ddl.TruncateTable("users").Restrict()
	sql, args, err = truncateRestrict.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("With restrict: %s\n", sql)
	fmt.Printf("Args: %v\n\n", args)

	// PostgreSQL-specific examples
	fmt.Println("=== PostgreSQL Examples ===")
	fmt.Println()

	// Truncate with restart identity (PostgreSQL)
	truncateRestart := ddl.TruncateTable("users").Restart().WithDialect(stk.Postgres())
	sql, args, err = truncateRestart.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("With restart identity: %s\n", sql)
	fmt.Printf("Args: %v\n\n", args)

	// Truncate with continue identity (PostgreSQL)
	truncateContinue := ddl.TruncateTable("users").Continue().WithDialect(stk.Postgres())
	sql, args, err = truncateContinue.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("With continue identity: %s\n", sql)
	fmt.Printf("Args: %v\n\n", args)

	// Complex PostgreSQL example
	truncateComplex := ddl.TruncateTable("users", "orders").
		Restart().
		Cascade().
		WithDialect(stk.Postgres())
	sql, args, err = truncateComplex.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Complex (restart + cascade): %s\n", sql)
	fmt.Printf("Args: %v\n\n", args)

	// MySQL example with dialect quoting
	fmt.Println("=== MySQL Examples ===")
	fmt.Println()

	truncateMySQL := ddl.TruncateTable("users").WithDialect(stk.MySQL())
	sql, args, err = truncateMySQL.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("MySQL with quoting: %s\n", sql)
	fmt.Printf("Args: %v\n\n", args)

	// Error handling example
	fmt.Println("=== Error Handling ===")
	fmt.Println()

	// This will cause an error
	truncateError := ddl.TruncateTable()
	_, _, err = truncateError.Build()
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
}
