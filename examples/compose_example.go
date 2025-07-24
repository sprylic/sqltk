package main

import (
	"fmt"

	"github.com/sprylic/sqltk"
)

func main() {
	// Example 1: Basic composition - combining columns and where conditions
	fmt.Println("=== Example 1: Basic Composition ===")
	q1 := sqltk.Select("id", "name").From("users").Where(sqltk.NewStringCondition("active = ?", true))
	q2 := sqltk.Select("email").From("users").Where(sqltk.NewStringCondition("verified = ?", true))

	combined := q1.Compose(q2)
	sql, args, _ := combined.Build()
	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Args: %v\n\n", args)

	// Example 2: Composition with joins
	fmt.Println("=== Example 2: Composition with Joins ===")
	q3 := sqltk.Select("u.id", "u.name").From("users u").Join("posts p").On("p.user_id", "u.id")
	q4 := sqltk.Select("u.email").From("users u").LeftJoin("profiles pr").On("pr.user_id", "u.id")

	combined2 := q3.Compose(q4)
	sql2, _, _ := combined2.Build()
	fmt.Printf("SQL: %s\n\n", sql2)

	// Example 3: Composition with group by and having
	fmt.Println("=== Example 3: Composition with Group By and Having ===")
	q5 := sqltk.Select("user_id", "COUNT(*)").From("orders").GroupBy("user_id").Having(sqltk.NewStringCondition("COUNT(*) > ?", 5))
	q6 := sqltk.Select("SUM(amount)").From("orders").GroupBy("user_id").Having(sqltk.NewStringCondition("SUM(amount) > ?", 1000))

	combined3 := q5.Compose(q6)
	sql3, args3, _ := combined3.Build()
	fmt.Printf("SQL: %s\n", sql3)
	fmt.Printf("Args: %v\n\n", args3)

	// Example 4: Composition with limit and offset
	fmt.Println("=== Example 4: Composition with Limit and Offset ===")
	q7 := sqltk.Select("id").From("users").Limit(10)
	q8 := sqltk.Select("name").From("users").Offset(5)
	q9 := sqltk.Select("email").From("users").Limit(5).Offset(10)

	combined4 := q7.Compose(q8, q9)
	sql4, _, _ := combined4.Build()
	fmt.Printf("SQL: %s\n\n", sql4)

	// Example 5: Composition with distinct
	fmt.Println("=== Example 5: Composition with Distinct ===")
	q10 := sqltk.Select("id").From("users")
	q11 := sqltk.Select("name").From("users").Distinct()

	combined5 := q10.Compose(q11)
	sql5, _, _ := combined5.Build()
	fmt.Printf("SQL: %s\n\n", sql5)

	// Example 6: Composition with subqueries
	fmt.Println("=== Example 6: Composition with Subqueries ===")
	sub1 := sqltk.Select("COUNT(*)").From("posts").Where(sqltk.NewStringCondition("user_id = users.id"))
	sub2 := sqltk.Select("MAX(created_at)").From("posts").Where(sqltk.NewStringCondition("user_id = users.id"))

	q12 := sqltk.Select("id", sub1).From("users")
	q13 := sqltk.Select("name", sub2).From("users")

	combined6 := q12.Compose(q13)
	sql6, _, _ := combined6.Build()
	fmt.Printf("SQL: %s\n\n", sql6)

	// Example 7: Composition with aliases
	fmt.Println("=== Example 7: Composition with Aliases ===")
	q14 := sqltk.Select(sqltk.Alias("id", "user_id")).From("users")
	q15 := sqltk.Select(sqltk.Alias("name", "user_name")).From("users")

	combined7 := q14.Compose(q15)
	sql7, _, _ := combined7.Build()
	fmt.Printf("SQL: %s\n\n", sql7)
}
