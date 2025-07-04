package sqltk

import (
	"testing"

	"github.com/sprylic/sqltk/ddl"
	"github.com/sprylic/sqltk/shared"
)

func TestCreateSchema(t *testing.T) {
	tests := []struct {
		name     string
		builder  *ddl.CreateSchemaBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic create schema",
			builder:  ddl.CreateSchema("testschema"),
			expected: "CREATE SCHEMA `testschema`",
		},
		{
			name:     "create schema with if not exists",
			builder:  ddl.CreateSchema("testschema").IfNotExists(),
			expected: "CREATE SCHEMA IF NOT EXISTS `testschema`",
		},
		{
			name:     "create schema with authorization",
			builder:  ddl.CreateSchema("testschema").Authorization("testuser"),
			expected: "CREATE SCHEMA `testschema` AUTHORIZATION `testuser`",
		},
		{
			name:     "create schema with if not exists and authorization",
			builder:  ddl.CreateSchema("testschema").IfNotExists().Authorization("testuser"),
			expected: "CREATE SCHEMA IF NOT EXISTS `testschema` AUTHORIZATION `testuser`",
		},
		{
			name:     "create schema with custom option",
			builder:  ddl.CreateSchema("testschema").Option("DEFAULT_CHARACTER_SET", "utf8mb4"),
			expected: "CREATE SCHEMA `testschema` DEFAULT_CHARACTER_SET utf8mb4",
		},
		{
			name:     "create schema with multiple options",
			builder:  ddl.CreateSchema("testschema").Option("DEFAULT_CHARACTER_SET", "utf8mb4").Option("DEFAULT_COLLATION", "utf8mb4_unicode_ci"),
			expected: "CREATE SCHEMA `testschema` DEFAULT_CHARACTER_SET utf8mb4 DEFAULT_COLLATION utf8mb4_unicode_ci",
		},
		{
			name:    "empty schema name",
			builder: ddl.CreateSchema(""),
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

func TestCreateSchemaPostgres(t *testing.T) {
	tests := []struct {
		name     string
		builder  *ddl.CreateSchemaBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic create schema",
			builder:  ddl.CreateSchema("testschema"),
			expected: `CREATE SCHEMA "testschema"`,
		},
		{
			name:     "create schema with if not exists",
			builder:  ddl.CreateSchema("testschema").IfNotExists(),
			expected: `CREATE SCHEMA IF NOT EXISTS "testschema"`,
		},
		{
			name:     "create schema with authorization",
			builder:  ddl.CreateSchema("testschema").Authorization("testuser"),
			expected: `CREATE SCHEMA "testschema" AUTHORIZATION "testuser"`,
		},
		{
			name:     "create schema with custom option",
			builder:  ddl.CreateSchema("testschema").Option("QUOTA", "100MB"),
			expected: `CREATE SCHEMA "testschema" QUOTA 100MB`,
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

func TestCreateSchemaNoQuoteIdent(t *testing.T) {
	tests := []struct {
		name     string
		builder  *ddl.CreateSchemaBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic create schema",
			builder:  ddl.CreateSchema("testschema"),
			expected: "CREATE SCHEMA testschema",
		},
		{
			name:     "create schema with if not exists",
			builder:  ddl.CreateSchema("testschema").IfNotExists(),
			expected: "CREATE SCHEMA IF NOT EXISTS testschema",
		},
		{
			name:     "create schema with authorization",
			builder:  ddl.CreateSchema("testschema").Authorization("testuser"),
			expected: "CREATE SCHEMA testschema AUTHORIZATION testuser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.builder.WithDialect(shared.NoQuoteIdent())
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

func TestCreateSchemaDebugSQL(t *testing.T) {
	builder := ddl.CreateSchema("testschema").IfNotExists().Authorization("testuser")
	builder.WithDialect(shared.MySQL())

	debugSQL := builder.DebugSQL()
	expected := "CREATE SCHEMA IF NOT EXISTS `testschema` AUTHORIZATION `testuser`"

	if debugSQL != expected {
		t.Errorf("expected debug SQL %q, got %q", expected, debugSQL)
	}
}
