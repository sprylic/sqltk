package main

import (
	"fmt"

	"github.com/sprylic/sqltk"
	"github.com/sprylic/sqltk/mysqlfunc"
	"github.com/sprylic/sqltk/pgfunc"
)

func main() {
	// Set dialect to not quote identifiers for cleaner output
	sqltk.SetDialect(sqltk.NoQuoteIdent())

	fmt.Println("=== MySQL Functions Example ===")

	// Basic date/time functions
	mysqlQuery1 := sqltk.Select(mysqlfunc.CurrentTimestamp()).From("users")
	sql1, args1, _ := mysqlQuery1.Build()
	fmt.Printf("MySQL Current Timestamp: %s\n", sql1)

	// String concatenation
	mysqlQuery2 := sqltk.Select(
		sqltk.Alias(mysqlfunc.Concat("first_name", "' '", "last_name"), "full_name"),
	).From("users")
	sql2, args2, _ := mysqlQuery2.Build()
	fmt.Printf("MySQL Concat: %s\n", sql2)

	// Aggregate functions
	mysqlQuery3 := sqltk.Select(
		sqltk.Alias(mysqlfunc.Count("*"), "total_users"),
		sqltk.Alias(mysqlfunc.Avg("age"), "avg_age"),
	).From("users")
	sql3, args3, _ := mysqlQuery3.Build()
	fmt.Printf("MySQL Aggregates: %s\n", sql3)

	// Conditional functions
	mysqlQuery4 := sqltk.Select(
		sqltk.Alias(mysqlfunc.If("active", "'Active'", "'Inactive'"), "status"),
	).From("users")
	sql4, args4, _ := mysqlQuery4.Build()
	fmt.Printf("MySQL Conditional: %s\n", sql4)

	// Date formatting
	mysqlQuery5 := sqltk.Select(
		sqltk.Alias(mysqlfunc.DateFormat("created_at", "'%Y-%m-%d'"), "created_date"),
	).From("users")
	sql5, args5, _ := mysqlQuery5.Build()
	fmt.Printf("MySQL Date Format: %s\n", sql5)

	fmt.Println("\n=== PostgreSQL Functions Example ===")

	// Basic date/time functions
	pgQuery1 := sqltk.Select(pgfunc.Now()).From("users")
	pgSQL1, pgArgs1, _ := pgQuery1.Build()
	fmt.Printf("PostgreSQL Now: %s\n", pgSQL1)

	// String concatenation
	pgQuery2 := sqltk.Select(
		sqltk.Alias(pgfunc.Concat("first_name", "' '", "last_name"), "full_name"),
	).From("users")
	pgSQL2, pgArgs2, _ := pgQuery2.Build()
	fmt.Printf("PostgreSQL Concat: %s\n", pgSQL2)

	// Aggregate functions
	pgQuery3 := sqltk.Select(
		sqltk.Alias(pgfunc.Count("*"), "total_users"),
		sqltk.Alias(pgfunc.Avg("age"), "avg_age"),
	).From("users")
	pgSQL3, pgArgs3, _ := pgQuery3.Build()
	fmt.Printf("PostgreSQL Aggregates: %s\n", pgSQL3)

	// Conditional functions
	pgQuery4 := sqltk.Select(
		sqltk.Alias(pgfunc.Coalesce("nickname", "first_name"), "display_name"),
	).From("users")
	pgSQL4, pgArgs4, _ := pgQuery4.Build()
	fmt.Printf("PostgreSQL Conditional: %s\n", pgSQL4)

	// Date extraction
	pgQuery5 := sqltk.Select(
		sqltk.Alias(pgfunc.Extract("year", "created_at"), "created_year"),
	).From("users")
	pgSQL5, pgArgs5, _ := pgQuery5.Build()
	fmt.Printf("PostgreSQL Date Extract: %s\n", pgSQL5)

	// Array aggregation
	pgQuery6 := sqltk.Select(
		sqltk.Alias(pgfunc.ArrayAgg("tag"), "all_tags"),
	).From("posts").GroupBy("category")
	pgSQL6, pgArgs6, _ := pgQuery6.Build()
	fmt.Printf("PostgreSQL Array Agg: %s\n", pgSQL6)

	fmt.Println("\n=== Function Usage in WHERE Clauses ===")

	// Using functions in WHERE clauses
	mysqlWhereQuery := sqltk.Select("id", "name").From("users").Where(
		sqltk.AsCondition(sqltk.Raw("created_at > " + string(mysqlfunc.CurrentTimestamp()))),
	)
	mysqlWhereSQL, mysqlWhereArgs, _ := mysqlWhereQuery.Build()
	fmt.Printf("MySQL WHERE with function: %s\n", mysqlWhereSQL)

	pgWhereQuery := sqltk.Select("id", "name").From("users").Where(
		sqltk.AsCondition(sqltk.Raw("created_at > " + string(pgfunc.Now()))),
	)
	pgWhereSQL, pgWhereArgs, _ := pgWhereQuery.Build()
	fmt.Printf("PostgreSQL WHERE with function: %s\n", pgWhereSQL)

	fmt.Println("\n=== Function Usage in ORDER BY ===")

	// Using functions in ORDER BY
	mysqlOrderQuery := sqltk.Select("id", "name").From("users").OrderBy(mysqlfunc.Random())
	mysqlOrderSQL, mysqlOrderArgs, _ := mysqlOrderQuery.Build()
	fmt.Printf("MySQL ORDER BY with function: %s\n", mysqlOrderSQL)

	pgOrderQuery := sqltk.Select("id", "name").From("users").OrderBy(pgfunc.Random())
	pgOrderSQL, pgOrderArgs, _ := pgOrderQuery.Build()
	fmt.Printf("PostgreSQL ORDER BY with function: %s\n", pgOrderSQL)

	// Print argument counts to show that functions don't generate parameters
	fmt.Printf("\nAll queries have 0 arguments: MySQL=%d, PostgreSQL=%d\n",
		len(args1)+len(args2)+len(args3)+len(args4)+len(args5)+len(mysqlWhereArgs)+len(mysqlOrderArgs),
		len(pgArgs1)+len(pgArgs2)+len(pgArgs3)+len(pgArgs4)+len(pgArgs5)+len(pgArgs6)+len(pgWhereArgs)+len(pgOrderArgs))

}
