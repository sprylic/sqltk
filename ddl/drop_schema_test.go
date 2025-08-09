package ddl

import (
	"testing"

	"github.com/sprylic/sqltk/sqldialect"
)

func TestDropSchema(t *testing.T) {
	tests := []struct {
		name     string
		builder  *DropSchemaBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic drop schema",
			builder:  DropSchema("testschema"),
			expected: "DROP SCHEMA `testschema`",
		},
		{
			name:     "drop schema with if exists",
			builder:  DropSchema("testschema").IfExists(),
			expected: "DROP SCHEMA IF EXISTS `testschema`",
		},
		{
			name:     "drop schema with cascade",
			builder:  DropSchema("testschema").Cascade(),
			expected: "DROP SCHEMA `testschema` CASCADE",
		},
		{
			name:     "drop schema with restrict",
			builder:  DropSchema("testschema").Restrict(),
			expected: "DROP SCHEMA `testschema` RESTRICT",
		},
		{
			name:     "drop schema with if exists and cascade",
			builder:  DropSchema("testschema").IfExists().Cascade(),
			expected: "DROP SCHEMA IF EXISTS `testschema` CASCADE",
		},
		{
			name:     "drop schema with if exists and restrict",
			builder:  DropSchema("testschema").IfExists().Restrict(),
			expected: "DROP SCHEMA IF EXISTS `testschema` RESTRICT",
		},
		{
			name:    "empty schema name",
			builder: DropSchema(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.builder.WithDialect(sqldialect.MySQL())
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

func TestDropSchemaPostgres(t *testing.T) {
	tests := []struct {
		name     string
		builder  *DropSchemaBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic drop schema",
			builder:  DropSchema("testschema"),
			expected: `DROP SCHEMA "testschema"`,
		},
		{
			name:     "drop schema with if exists",
			builder:  DropSchema("testschema").IfExists(),
			expected: `DROP SCHEMA IF EXISTS "testschema"`,
		},
		{
			name:     "drop schema with cascade",
			builder:  DropSchema("testschema").Cascade(),
			expected: `DROP SCHEMA "testschema" CASCADE`,
		},
		{
			name:     "drop schema with restrict",
			builder:  DropSchema("testschema").Restrict(),
			expected: `DROP SCHEMA "testschema" RESTRICT`,
		},
		{
			name:     "drop schema with if exists and cascade",
			builder:  DropSchema("testschema").IfExists().Cascade(),
			expected: `DROP SCHEMA IF EXISTS "testschema" CASCADE`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.builder.WithDialect(sqldialect.Postgres())
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

func TestDropSchemaNoQuoteIdent(t *testing.T) {
	tests := []struct {
		name     string
		builder  *DropSchemaBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic drop schema",
			builder:  DropSchema("testschema"),
			expected: "DROP SCHEMA testschema",
		},
		{
			name:     "drop schema with if exists",
			builder:  DropSchema("testschema").IfExists(),
			expected: "DROP SCHEMA IF EXISTS testschema",
		},
		{
			name:     "drop schema with cascade",
			builder:  DropSchema("testschema").Cascade(),
			expected: "DROP SCHEMA testschema CASCADE",
		},
		{
			name:     "drop schema with restrict",
			builder:  DropSchema("testschema").Restrict(),
			expected: "DROP SCHEMA testschema RESTRICT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.builder.WithDialect(sqldialect.NoQuoteIdent())
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

func TestDropSchemaDebugSQL(t *testing.T) {
	builder := DropSchema("testschema").IfExists().Cascade()
	builder.WithDialect(sqldialect.Postgres())

	debugSQL := builder.DebugSQL()
	expected := `DROP SCHEMA IF EXISTS "testschema" CASCADE`

	if debugSQL != expected {
		t.Errorf("expected debug SQL %q, got %q", expected, debugSQL)
	}
}
