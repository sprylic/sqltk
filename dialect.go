package cqb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// Dialect defines SQL dialect-specific behavior for placeholders and identifier quoting.
type Dialect interface {
	Placeholder(n int) string
	QuoteIdent(ident string) string
	QuoteString(s string) string
}

// standardDialect uses ? for all placeholders and no identifier quoting.
type standardDialect struct{}

func (standardDialect) Placeholder(n int) string       { return "?" }
func (standardDialect) QuoteIdent(ident string) string { return ident }
func (standardDialect) QuoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// mySQLDialect uses ? for all placeholders and backticks for identifier quoting.
type mySQLDialect struct{}

func (mySQLDialect) Placeholder(n int) string       { return "?" }
func (mySQLDialect) QuoteIdent(ident string) string { return "`" + ident + "`" }
func (mySQLDialect) QuoteString(s string) string    { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }

// postgresDialect uses $n for placeholders and double quotes for identifier quoting.
type postgresDialect struct{}

func (postgresDialect) Placeholder(n int) string       { return "$" + fmt.Sprint(n) }
func (postgresDialect) QuoteIdent(ident string) string { return "\"" + ident + "\"" }
func (postgresDialect) QuoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

var (
	standardDialectInstance = standardDialect{}
	mySQLDialectInstance    = mySQLDialect{}
	postgresDialectInstance = postgresDialect{}

	dialectMu     sync.RWMutex
	globalDialect Dialect = &mySQLDialectInstance
)

// Standard returns the standard SQL dialect (no quoting, not the default).
func Standard() Dialect { return &standardDialectInstance }

// MySQL returns the MySQL SQL dialect (default).
func MySQL() Dialect { return &mySQLDialectInstance }

// Postgres returns the Postgres SQL dialect.
func Postgres() Dialect { return &postgresDialectInstance }

// SetDialect sets the global SQL dialect for all builders.
func SetDialect(d Dialect) {
	dialectMu.Lock()
	defer dialectMu.Unlock()
	globalDialect = d
}

// getDialect returns the current global SQL dialect.
func getDialect() Dialect {
	dialectMu.RLock()
	defer dialectMu.RUnlock()
	return globalDialect
}

// PGJSON wraps a value for JSON encoding in Postgres queries.
type PGJSON struct {
	V interface{}
}

// Value implements driver.Valuer for PGJSON.
func (j PGJSON) Value() (driver.Value, error) {
	if j.V == nil {
		return nil, nil
	}
	return json.Marshal(j.V)
}

// Example usage:
//   pq := sq.NewPostgresInsert("users").Columns("name", "data").Values("Alice", sq.PGJSON(map[string]interface{}{...})).Returning("id")
//   sql, args, err := pq.Build()

// PGArray wraps a value for Postgres array encoding in queries.
type PGArray struct {
	V interface{}
}

// Value implements driver.Valuer for PGArray.
func (a PGArray) Value() (driver.Value, error) {
	// For simple types, database/sql/driver will handle slices as Postgres arrays.
	return a.V, nil
}

// Example usage:
//   pq := sq.NewPostgresInsert("users").Columns("tags").Values(sq.PGArray([]string{"foo", "bar"})).Returning("id")
//   sql, args, err := pq.Build()
