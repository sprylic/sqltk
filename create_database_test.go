package sqltk

import (
	"testing"

	"github.com/sprylic/sqltk/ddl"
	"github.com/sprylic/sqltk/shared"
)

func TestCreateDatabase(t *testing.T) {
	tests := []struct {
		name     string
		builder  *ddl.CreateDatabaseBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic create database",
			builder:  ddl.CreateDatabase("testdb"),
			expected: "CREATE DATABASE `testdb`",
		},
		{
			name:     "create database with if not exists",
			builder:  ddl.CreateDatabase("testdb").IfNotExists(),
			expected: "CREATE DATABASE IF NOT EXISTS `testdb`",
		},
		{
			name:     "create database with charset",
			builder:  ddl.CreateDatabase("testdb").Charset("utf8mb4"),
			expected: "CREATE DATABASE `testdb` CHARACTER SET utf8mb4",
		},
		{
			name:     "create database with collation",
			builder:  ddl.CreateDatabase("testdb").Collation("utf8mb4_unicode_ci"),
			expected: "CREATE DATABASE `testdb` COLLATE utf8mb4_unicode_ci",
		},
		{
			name:     "create database with charset and collation",
			builder:  ddl.CreateDatabase("testdb").Charset("utf8mb4").Collation("utf8mb4_unicode_ci"),
			expected: "CREATE DATABASE `testdb` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		},
		{
			name:     "create database with custom option",
			builder:  ddl.CreateDatabase("testdb").Option("ENCRYPTION", "Y"),
			expected: "CREATE DATABASE `testdb` ENCRYPTION Y",
		},
		{
			name:     "create database with multiple options",
			builder:  ddl.CreateDatabase("testdb").Option("ENCRYPTION", "Y").Option("READ ONLY", "1"),
			expected: "CREATE DATABASE `testdb` ENCRYPTION Y READ ONLY 1",
		},
		{
			name:    "empty database name",
			builder: ddl.CreateDatabase(""),
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

func TestCreateDatabasePostgres(t *testing.T) {
	tests := []struct {
		name     string
		builder  *ddl.CreateDatabaseBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic create database",
			builder:  ddl.CreateDatabase("testdb"),
			expected: `CREATE DATABASE "testdb"`,
		},
		{
			name:     "create database with if not exists",
			builder:  ddl.CreateDatabase("testdb").IfNotExists(),
			expected: `CREATE DATABASE IF NOT EXISTS "testdb"`,
		},
		{
			name:     "create database with charset",
			builder:  ddl.CreateDatabase("testdb").Charset("UTF8"),
			expected: `CREATE DATABASE "testdb" CHARACTER SET UTF8`,
		},
		{
			name:     "create database with collation",
			builder:  ddl.CreateDatabase("testdb").Collation("en_US.utf8"),
			expected: `CREATE DATABASE "testdb" COLLATE en_US.utf8`,
		},
		{
			name:     "create database with custom option",
			builder:  ddl.CreateDatabase("testdb").Option("TEMPLATE", "template0"),
			expected: `CREATE DATABASE "testdb" TEMPLATE template0`,
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

func TestCreateDatabaseNoQuoteIdent(t *testing.T) {
	tests := []struct {
		name     string
		builder  *ddl.CreateDatabaseBuilder
		expected string
		wantErr  bool
	}{
		{
			name:     "basic create database",
			builder:  ddl.CreateDatabase("testdb"),
			expected: "CREATE DATABASE testdb",
		},
		{
			name:     "create database with if not exists",
			builder:  ddl.CreateDatabase("testdb").IfNotExists(),
			expected: "CREATE DATABASE IF NOT EXISTS testdb",
		},
		{
			name:     "create database with charset",
			builder:  ddl.CreateDatabase("testdb").Charset("UTF8"),
			expected: "CREATE DATABASE testdb CHARACTER SET UTF8",
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

func TestCreateDatabaseDebugSQL(t *testing.T) {
	builder := ddl.CreateDatabase("testdb").Charset("utf8mb4").Collation("utf8mb4_unicode_ci")
	builder.WithDialect(shared.MySQL())

	debugSQL := builder.DebugSQL()
	expected := "CREATE DATABASE `testdb` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"

	if debugSQL != expected {
		t.Errorf("expected debug SQL %q, got %q", expected, debugSQL)
	}
}
