package ddl

import (
	"testing"

	"github.com/sprylic/sqltk/sqldialect"
)

func TestDropDatabase(t *testing.T) {
	tests := []struct {
		name     string
		builder  *DropDatabaseBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic drop database",
			builder:  DropDatabase("testdb"),
			expected: "DROP DATABASE `testdb`",
		},
		{
			name:     "drop database with if exists",
			builder:  DropDatabase("testdb").IfExists(),
			expected: "DROP DATABASE IF EXISTS `testdb`",
		},
		{
			name:    "empty database name",
			builder: DropDatabase(""),
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

func TestDropDatabasePostgres(t *testing.T) {
	tests := []struct {
		name     string
		builder  *DropDatabaseBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic drop database",
			builder:  DropDatabase("testdb"),
			expected: `DROP DATABASE "testdb"`,
		},
		{
			name:     "drop database with if exists",
			builder:  DropDatabase("testdb").IfExists(),
			expected: `DROP DATABASE IF EXISTS "testdb"`,
		},
		{
			name:     "drop database with cascade",
			builder:  DropDatabase("testdb").Cascade(),
			expected: `DROP DATABASE "testdb" CASCADE`,
		},
		{
			name:     "drop database with if exists and cascade",
			builder:  DropDatabase("testdb").IfExists().Cascade(),
			expected: `DROP DATABASE IF EXISTS "testdb" CASCADE`,
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

func TestDropDatabaseNoQuoteIdent(t *testing.T) {
	tests := []struct {
		name     string
		builder  *DropDatabaseBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic drop database",
			builder:  DropDatabase("testdb"),
			expected: "DROP DATABASE testdb",
		},
		{
			name:     "drop database with if exists",
			builder:  DropDatabase("testdb").IfExists(),
			expected: "DROP DATABASE IF EXISTS testdb",
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

func TestDropDatabaseDebugSQL(t *testing.T) {
	builder := DropDatabase("testdb").IfExists().Cascade()
	builder.WithDialect(sqldialect.Postgres())

	debugSQL := builder.DebugSQL()
	expected := `DROP DATABASE IF EXISTS "testdb" CASCADE`

	if debugSQL != expected {
		t.Errorf("expected debug SQL %q, got %q", expected, debugSQL)
	}
}
