package stk

import (
	"testing"

	"github.com/sprylic/stk/ddl"
	"github.com/sprylic/stk/shared"
)

func TestDropDatabase(t *testing.T) {
	tests := []struct {
		name     string
		builder  *ddl.DropDatabaseBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic drop database",
			builder:  ddl.DropDatabase("testdb"),
			expected: "DROP DATABASE `testdb`",
		},
		{
			name:     "drop database with if exists",
			builder:  ddl.DropDatabase("testdb").IfExists(),
			expected: "DROP DATABASE IF EXISTS `testdb`",
		},
		{
			name:    "empty database name",
			builder: ddl.DropDatabase(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.builder.WithDialect(shared.MySQL())
			sql, args, err := tt.builder.Build()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if sql != tt.expected {
				t.Errorf("expected SQL %q, got %q", tt.expected, sql)
			}

			if len(args) != 0 {
				t.Errorf("expected no arguments, got %d", len(args))
			}
		})
	}
}

func TestDropDatabasePostgres(t *testing.T) {
	tests := []struct {
		name     string
		builder  *ddl.DropDatabaseBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic drop database",
			builder:  ddl.DropDatabase("testdb"),
			expected: `DROP DATABASE "testdb"`,
		},
		{
			name:     "drop database with if exists",
			builder:  ddl.DropDatabase("testdb").IfExists(),
			expected: `DROP DATABASE IF EXISTS "testdb"`,
		},
		{
			name:     "drop database with cascade",
			builder:  ddl.DropDatabase("testdb").Cascade(),
			expected: `DROP DATABASE "testdb" CASCADE`,
		},
		{
			name:     "drop database with if exists and cascade",
			builder:  ddl.DropDatabase("testdb").IfExists().Cascade(),
			expected: `DROP DATABASE IF EXISTS "testdb" CASCADE`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.builder.WithDialect(shared.Postgres())
			sql, args, err := tt.builder.Build()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if sql != tt.expected {
				t.Errorf("expected SQL %q, got %q", tt.expected, sql)
			}

			if len(args) != 0 {
				t.Errorf("expected no arguments, got %d", len(args))
			}
		})
	}
}

func TestDropDatabaseStandard(t *testing.T) {
	tests := []struct {
		name     string
		builder  *ddl.DropDatabaseBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic drop database",
			builder:  ddl.DropDatabase("testdb"),
			expected: "DROP DATABASE testdb",
		},
		{
			name:     "drop database with if exists",
			builder:  ddl.DropDatabase("testdb").IfExists(),
			expected: "DROP DATABASE IF EXISTS testdb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.builder.WithDialect(shared.Standard())
			sql, args, err := tt.builder.Build()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if sql != tt.expected {
				t.Errorf("expected SQL %q, got %q", tt.expected, sql)
			}

			if len(args) != 0 {
				t.Errorf("expected no arguments, got %d", len(args))
			}
		})
	}
}

func TestDropDatabaseDebugSQL(t *testing.T) {
	builder := ddl.DropDatabase("testdb").IfExists().Cascade()
	builder.WithDialect(shared.Postgres())

	debugSQL := builder.DebugSQL()
	expected := `DROP DATABASE IF EXISTS "testdb" CASCADE`

	if debugSQL != expected {
		t.Errorf("expected debug SQL %q, got %q", expected, debugSQL)
	}
}
