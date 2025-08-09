package main

import (
	"fmt"

	"github.com/sprylic/sqltk/ddl"
	"github.com/sprylic/sqltk/sqldialect"
)

func main() {
	// Example: Creating a table with OnUpdate for automatic timestamp updates

	// MySQL version - uses ON UPDATE CURRENT_TIMESTAMP
	fmt.Println("=== MySQL Version ===")
	mysqlTable := ddl.CreateTable("users").
		AddColumn(ddl.Column("id").Type("INT").PrimaryKey()).
		AddColumn(ddl.Column("name").Type("VARCHAR").Size(255).NotNull()).
		AddColumn(ddl.Column("email").Type("VARCHAR").Size(255)).
		AddColumn(ddl.Column("updated_at").Type("TIMESTAMP").OnUpdate("CURRENT_TIMESTAMP"))

	mysqlSQL, _, err := mysqlTable.WithDialect(sqldialect.MySQL()).Build()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(mysqlSQL)
	fmt.Println()

	// PostgreSQL version - uses triggers for equivalent functionality
	fmt.Println("=== PostgreSQL Version ===")
	postgresTable := ddl.CreateTable("users").
		AddColumn(ddl.Column("id").Type("INT").PrimaryKey()).
		AddColumn(ddl.Column("name").Type("VARCHAR").Size(255).NotNull()).
		AddColumn(ddl.Column("email").Type("VARCHAR").Size(255)).
		AddColumn(ddl.Column("updated_at").Type("TIMESTAMP").OnUpdate("CURRENT_TIMESTAMP"))

	postgresSQL, _, err := postgresTable.WithDialect(sqldialect.Postgres()).Build()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(postgresSQL)
	fmt.Println()

	// Example with multiple OnUpdate columns
	fmt.Println("=== PostgreSQL with Multiple OnUpdate Columns ===")
	multiTable := ddl.CreateTable("products").
		AddColumn(ddl.Column("id").Type("INT").PrimaryKey()).
		AddColumn(ddl.Column("name").Type("VARCHAR").Size(255).NotNull()).
		AddColumn(ddl.Column("updated_at").Type("TIMESTAMP").OnUpdate("CURRENT_TIMESTAMP")).
		AddColumn(ddl.Column("modified_at").Type("TIMESTAMP").OnUpdate("NOW()"))

	multiSQL, _, err := multiTable.WithDialect(sqldialect.Postgres()).Build()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(multiSQL)
	fmt.Println()

	// Example with IfNotExists - safe for repeated execution
	fmt.Println("=== PostgreSQL with IfNotExists (Safe for Repeated Execution) ===")
	safeTable := ddl.CreateTable("users").
		IfNotExists().
		AddColumn(ddl.Column("id").Type("INT").PrimaryKey()).
		AddColumn(ddl.Column("name").Type("VARCHAR").Size(255).NotNull()).
		AddColumn(ddl.Column("updated_at").Type("TIMESTAMP").OnUpdate("CURRENT_TIMESTAMP"))

	safeSQL, _, err := safeTable.WithDialect(sqldialect.Postgres()).Build()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(safeSQL)
}
