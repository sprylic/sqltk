package sqlfunc

import (
	"testing"
)

func TestValidateSqlFuncInput(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		// Safe inputs
		{"column name", "created_at", false},
		{"simple string", "name", false},
		{"number", 42, false},
		{"sqlfunc", SqlFunc("UPPER(name)"), false},
		{"quoted string", "'%Y-%m-%d'", false},

		// Dangerous inputs
		{"sql injection with semicolon", "'; DROP TABLE users; --", true},
		{"sql injection with comment", "/* DROP TABLE users */", true},
		{"sql injection with union", "UNION SELECT * FROM users", true},
		{"sql injection with select", "SELECT * FROM users", true},
		{"sql injection with drop", "DROP TABLE users", true},
		{"sql injection with comment marker", "name -- DROP TABLE users", true},
		{"sql injection with hash comment", "name # DROP TABLE users", true},
		{"sqlfunc with injection", SqlFunc("'; DROP TABLE users; --"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSqlFuncInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSqlFuncInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
