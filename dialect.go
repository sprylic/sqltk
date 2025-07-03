package stk

import "github.com/sprylic/stk/shared"

// Re-export shared types for backward compatibility
type (
	Dialect = shared.Dialect
	Raw     = shared.Raw
	PGJSON  = shared.PGJSON
	PGArray = shared.PGArray
)

// Re-export shared functions for backward compatibility
var (
	Standard       = shared.Standard
	MySQL          = shared.MySQL
	Postgres       = shared.Postgres
	SetDialect     = shared.SetDialect
	GetDialect     = shared.GetDialect
	InterpolateSQL = shared.InterpolateSQL
)
