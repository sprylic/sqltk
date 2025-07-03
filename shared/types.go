package shared

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

// GetDialect returns the current global SQL dialect.
func GetDialect() Dialect {
	dialectMu.RLock()
	defer dialectMu.RUnlock()
	return globalDialect
}

// Raw represents raw SQL that should be included directly without quoting.
type Raw string

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

// PGArray wraps a value for Postgres array encoding in queries.
type PGArray struct {
	V interface{}
}

// Value implements driver.Valuer for PGArray.
func (a PGArray) Value() (driver.Value, error) {
	// For simple types, database/sql/driver will handle slices as Postgres arrays.
	return a.V, nil
}

// InterpolateSQL interpolates arguments into a SQL query for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func InterpolateSQL(query string, args []interface{}) string {
	if len(args) == 0 {
		return query
	}

	// Simple interpolation - replace ? with values
	result := query
	argIndex := 0

	for i := 0; i < len(result) && argIndex < len(args); i++ {
		if result[i] == '?' {
			arg := args[argIndex]
			var argStr string

			switch v := arg.(type) {
			case string:
				argStr = "'" + strings.ReplaceAll(v, "'", "''") + "'"
			case nil:
				argStr = "NULL"
			default:
				argStr = fmt.Sprintf("%v", v)
			}

			result = result[:i] + argStr + result[i+1:]
			argIndex++
		}
	}

	return result
}
