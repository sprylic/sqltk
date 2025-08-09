package pgtypes

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

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
	if a.V == nil {
		return nil, nil
	}

	// Handle string slices specifically for PostgreSQL text arrays
	if strSlice, ok := a.V.([]string); ok {
		// Convert []string to PostgreSQL array format: {"value1","value2"}
		// Escape quotes in strings and wrap in curly braces
		escaped := make([]string, len(strSlice))
		for i, s := range strSlice {
			// Escape double quotes by doubling them
			escaped[i] = `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
		}
		result := "{" + strings.Join(escaped, ",") + "}"
		return result, nil
	}

	// For other types, let the driver handle it
	return a.V, nil
}
