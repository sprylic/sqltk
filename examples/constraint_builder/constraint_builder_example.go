package main

import (
	"fmt"
	"log"

	"github.com/sprylic/sqltk/ddl"
	"github.com/sprylic/sqltk/sqldialect"
)

// This example demonstrates the unified constraint builder API.
// Instead of multiple helper functions, we now have a single NewConstraint()
// function that can be chained with different constraint types.
func main() {
	// Set the dialect (optional, defaults to MySQL)
	sqldialect.SetDialect(sqldialect.MySQL())

	// Example 1: Create a table with various constraints using the unified constraint builder API
	createTableSQL, _, err := ddl.CreateTable("users").
		AddColumn(ddl.Column("id").Type("INT").PrimaryKey()).
		AddColumn(ddl.Column("email").Type("VARCHAR").Size(255).NotNull()).
		AddColumn(ddl.Column("age").Type("INT")).
		AddColumn(ddl.Column("role_id").Type("INT")).
		Unique("idx_email", "email").
		Check("chk_age", "age >= 0").
		AddForeignKey(ddl.ForeignKey("fk_user_role", "role_id").
			References("roles", "id").
			OnDelete("CASCADE")).
		Build()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Create Table with Constraints ===")
	fmt.Println(createTableSQL)
	fmt.Println()

	// Example 2: Alter table to add constraints using the unified constraint builder
	alterTableSQL, _, err := ddl.AlterTable("users").
		AddConstraint(ddl.NewConstraint().Check("chk_email_format", "email LIKE '%@%'")).
		AddConstraint(ddl.NewConstraint().Index("idx_name_email", "name", "email")).
		Build()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Alter Table with Constraints ===")
	fmt.Println(alterTableSQL)
	fmt.Println()

	// Example 3: Create a table with composite primary key
	createTableWithCompositePK, _, err := ddl.CreateTable("user_permissions").
		AddColumn(ddl.Column("user_id").Type("INT")).
		AddColumn(ddl.Column("permission_id").Type("INT")).
		AddColumn(ddl.Column("granted_at").Type("TIMESTAMP")).
		PrimaryKey("user_id", "permission_id").
		AddForeignKey(ddl.ForeignKey("fk_user_permissions_user", "user_id").
			References("users", "id").
			OnDelete("CASCADE")).
		AddForeignKey(ddl.ForeignKey("fk_user_permissions_permission", "permission_id").
			References("permissions", "id").
			OnDelete("CASCADE")).
		Build()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Create Table with Composite Primary Key ===")
	fmt.Println(createTableWithCompositePK)
	fmt.Println()

	// Example 4: Raw constraint for complex expressions
	createTableWithRawConstraint, _, err := ddl.CreateTable("products").
		AddColumn(ddl.Column("id").Type("INT").PrimaryKey()).
		AddColumn(ddl.Column("price").Type("DECIMAL").Precision(10, 2)).
		AddColumn(ddl.Column("discount").Type("DECIMAL").Precision(10, 2)).
		AddColumn(ddl.Column("category").Type("VARCHAR").Size(50)).
		Check("chk_price_logic", "price >= 0 AND discount >= 0 AND price >= discount").
		Check("chk_category", "category IN ('electronics', 'clothing', 'books')").
		Build()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Create Table with Raw Constraint ===")
	fmt.Println(createTableWithRawConstraint)
	fmt.Println()

	// Example 5: Using constraint builder with fluent chaining
	constraintBuilder := ddl.NewConstraint().
		ForeignKey("fk_complex", "user_id", "tenant_id").
		WithReference("users", "id", "tenant_id").
		WithOnDelete("CASCADE").
		WithOnUpdate("CASCADE")

	alterTableWithComplexFK, _, err := ddl.AlterTable("orders").
		AddConstraint(constraintBuilder).
		Build()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Alter Table with Complex Foreign Key ===")
	fmt.Println(alterTableWithComplexFK)
	fmt.Println()

	// Example 6: Add raw constraint directly
	alterTableWithRawConstraint, _, err := ddl.AlterTable("users").
		AddRawConstraint("chk_complex_logic", "age >= 0 AND age <= 150 AND status IN ('active', 'inactive')").
		Build()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Alter Table with Raw Constraint ===")
	fmt.Println(alterTableWithRawConstraint)
}
