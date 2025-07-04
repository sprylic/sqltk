package main

import (
	"fmt"

	"github.com/sprylic/stk"
)

func main() {
	fmt.Println("=== Type-Safe Conditions Example ===")
	fmt.Println()

	// Example 1: Type-safe string condition
	fmt.Println("1. Type-safe string condition:")
	stringCond := stk.NewStringCondition("active = ? AND age > ?", true, 18)
	q1 := stk.Select("id", "name").From("users").Where(stringCond)
	sql1, args1, _ := q1.Build()
	fmt.Printf("SQL: %s\n", sql1)
	fmt.Printf("Args: %v\n\n", args1)

	// Example 2: Type-safe raw condition
	fmt.Println("2. Type-safe raw condition:")
	rawCond := stk.NewRawCondition(stk.Raw("id = 1"))
	q2 := stk.Select("id", "name").From("users").Where(rawCond)
	sql2, args2, _ := q2.Build()
	fmt.Printf("SQL: %s\n", sql2)
	fmt.Printf("Args: %v\n\n", args2)

	// Example 3: ConditionBuilder (implements Condition interface)
	fmt.Println("3. ConditionBuilder (implements Condition interface):")
	condBuilder := stk.NewCond().
		Equal("active", true).
		And(stk.NewCond().GreaterThan("age", 18)).
		And(stk.NewCond().In("status", "active", "pending"))
	q3 := stk.Select("id", "name").From("users").Where(condBuilder)
	sql3, args3, _ := q3.Build()
	fmt.Printf("SQL: %s\n", sql3)
	fmt.Printf("Args: %v\n\n", args3)

	// Example 4: Complex condition with multiple types
	fmt.Println("4. Complex condition combining different types:")
	complexCond := stk.NewCond().
		Equal("active", true).
		And(stk.NewCond().GreaterThan("age", 18)).
		Or(stk.NewCond().Equal("vip", true))

	// Use in UPDATE
	updateQ := stk.Update("users").Set("last_login", "NOW()").Where(complexCond)
	updateSQL, updateArgs, _ := updateQ.Build()
	fmt.Printf("UPDATE SQL: %s\n", updateSQL)
	fmt.Printf("UPDATE Args: %v\n\n", updateArgs)

	// Example 5: Raw conditions now require AsCondition wrapper
	fmt.Println("5. Raw conditions now require AsCondition wrapper:")
	rawCond2 := stk.AsCondition(stk.Raw("id = 1"))
	q5 := stk.Select("id").From("users").Where(rawCond2)
	sql5, args5, _ := q5.Build()
	fmt.Printf("SQL: %s\n", sql5)
	fmt.Printf("Args: %v\n\n", args5)

	// Example 6: Compile-time type safety
	fmt.Println("6. Compile-time type safety:")
	fmt.Println("   - Invalid types (like int) are caught at compile time")
	fmt.Println("   - Raw SQL must be wrapped with AsCondition()")
	fmt.Println("   - String conditions must use NewStringCondition()")
	fmt.Println("   - This prevents runtime errors and improves safety")
	fmt.Println()

	// Example 7: Using in HAVING clause
	fmt.Println("7. Using in HAVING clause:")
	havingCond := stk.NewCond().GreaterThan("COUNT(*)", 10)
	havingQ := stk.Select("department", "COUNT(*) as count").
		From("employees").
		GroupBy("department").
		Having(havingCond)
	havingSQL, havingArgs, _ := havingQ.Build()
	fmt.Printf("HAVING SQL: %s\n", havingSQL)
	fmt.Printf("HAVING Args: %v\n", havingArgs)
}
