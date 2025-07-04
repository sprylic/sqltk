package sqltk

import "github.com/sprylic/sqltk/shared"

// Re-export shared types for backward compatibility
type (
	Dialect = shared.Dialect
	Raw     = shared.Raw
	PGJSON  = shared.PGJSON
	PGArray = shared.PGArray
)

// Re-export shared functions for backward compatibility
var (
	NoQuoteIdent   = shared.NoQuoteIdent
	MySQL          = shared.MySQL
	Postgres       = shared.Postgres
	SetDialect     = shared.SetDialect
	GetDialect     = shared.GetDialect
	InterpolateSQL = shared.InterpolateSQL
)
