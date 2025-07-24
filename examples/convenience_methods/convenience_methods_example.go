package main

import (
	"fmt"

	"github.com/sprylic/sqltk"
)

func main() {
	// Example 1: Basic convenience methods
	fmt.Println("=== Basic Convenience Methods ===")
	q1 := sqltk.Select("id", "name", "email").
		From("users").
		WhereEqual("active", true).
		WhereNotNull("email").
		WhereGreaterThan("age", 18).
		WhereLessThan("age", 65).
		WhereIn("status", "active", "pending").
		WhereLike("name", "%john%").
		WhereBetween("score", 70, 100)

	sql1, args1, _ := q1.Build()
	fmt.Printf("SQL: %s\n", sql1)
	fmt.Printf("Args: %v\n\n", args1)

	// Example 2: Column-to-column comparison
	fmt.Println("=== Column-to-Column Comparison ===")
	q2 := sqltk.Select("u.id", "u.name", "o.amount").
		From("users u").
		Join("orders o").On("u.id", "o.user_id").
		WhereColsEqual("u.status", "o.status").
		WhereGreaterThan("o.amount", 100)

	sql2, args2, _ := q2.Build()
	fmt.Printf("SQL: %s\n", sql2)
	fmt.Printf("Args: %v\n\n", args2)

	// Example 3: EXISTS with subqueries
	fmt.Println("=== EXISTS with Subqueries ===")
	sub1 := sqltk.Select("1").From("orders").WhereColsEqual("user_id", "users.id")
	sub2 := sqltk.Select("1").From("posts").WhereColsEqual("user_id", "users.id")

	q3 := sqltk.Select("id", "name").From("users").
		WhereExists(sub1).
		WhereNotExists(sub2).
		WhereEqual("active", true)

	sql3, args3, _ := q3.Build()
	fmt.Printf("SQL: %s\n", sql3)
	fmt.Printf("Args: %v\n\n", args3)

	// Example 4: Complex combination
	fmt.Println("=== Complex Combination ===")
	q4 := sqltk.Select("id", "name", "email").
		From("users").
		WhereEqual("active", true).
		WhereNotNull("email").
		WhereGreaterThan("age", 18).
		WhereLessThan("age", 65).
		WhereIn("status", "active", "pending").
		WhereLike("name", "%john%").
		WhereBetween("score", 70, 100).
		WhereNull("deleted_at").
		WhereNotLike("email", "%spam%").
		WhereNotIn("role", "admin", "moderator").
		WhereNotBetween("last_login", "2020-01-01", "2020-12-31")

	sql4, args4, _ := q4.Build()
	fmt.Printf("SQL: %s\n", sql4)
	fmt.Printf("Args: %v\n\n", args4)

	// Example 5: Update with convenience methods
	fmt.Println("=== Update with Convenience Methods ===")
	update := sqltk.Update("users").
		Set("last_login", "NOW()").
		Set("login_count", "login_count + 1").
		WhereEqual("id", 123).
		WhereNotNull("email").
		WhereGreaterThan("login_count", 0)

	sql5, args5, _ := update.Build()
	fmt.Printf("SQL: %s\n", sql5)
	fmt.Printf("Args: %v\n\n", args5)

	// Example 6: Delete with convenience methods
	fmt.Println("=== Delete with Convenience Methods ===")
	delete := sqltk.Delete("users").
		WhereEqual("active", false).
		WhereNull("email").
		WhereLessThan("last_login", "2020-01-01").
		WhereIn("status", "inactive", "suspended")

	sql6, args6, _ := delete.Build()
	fmt.Printf("SQL: %s\n", sql6)
	fmt.Printf("Args: %v\n", args6)
}
