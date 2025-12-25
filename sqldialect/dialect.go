package sqldialect

import (
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

// standardDialect uses ? for all placeholders and no identifier quoting (NoQuotes dialect).
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

// NoQuoteIdent returns a SQL dialect with no identifier quoting.
func NoQuoteIdent() Dialect { return &standardDialectInstance }

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
