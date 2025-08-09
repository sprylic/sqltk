package ddl

import (
	"errors"
	"github.com/sprylic/sqltk/sqldebug"
	"strings"

	"github.com/sprylic/sqltk/sqldialect"
)

// CreateDatabaseBuilder builds CREATE DATABASE statements.
type CreateDatabaseBuilder struct {
	name        string
	ifNotExists bool
	charset     string
	collation   string
	options     []DatabaseOption
	err         error
	dialect     sqldialect.Dialect
}

// DatabaseOption represents a database option.
type DatabaseOption struct {
	Name  string
	Value string
}

// CreateDatabase creates a new CreateDatabaseBuilder for the given database name.
func CreateDatabase(name string) *CreateDatabaseBuilder {
	if name == "" {
		return &CreateDatabaseBuilder{err: errors.New("database name is required")}
	}
	return &CreateDatabaseBuilder{name: name}
}

// IfNotExists adds IF NOT EXISTS clause to prevent errors if database already exists.
func (b *CreateDatabaseBuilder) IfNotExists() *CreateDatabaseBuilder {
	if b.err != nil {
		return b
	}
	b.ifNotExists = true
	return b
}

// Charset sets the character set for the database.
func (b *CreateDatabaseBuilder) Charset(charset string) *CreateDatabaseBuilder {
	if b.err != nil {
		return b
	}
	b.charset = charset
	return b
}

// Collation sets the collation for the database.
func (b *CreateDatabaseBuilder) Collation(collation string) *CreateDatabaseBuilder {
	if b.err != nil {
		return b
	}
	b.collation = collation
	return b
}

// Option adds a custom database option.
func (b *CreateDatabaseBuilder) Option(name, value string) *CreateDatabaseBuilder {
	if b.err != nil {
		return b
	}
	b.options = append(b.options, DatabaseOption{Name: name, Value: value})
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *CreateDatabaseBuilder) WithDialect(d sqldialect.Dialect) *CreateDatabaseBuilder {
	b.dialect = d
	return b
}

// Build builds the SQL CREATE DATABASE query and returns the query string, arguments, and error if any.
func (b *CreateDatabaseBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.name == "" {
		return "", nil, errors.New("database name is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = sqldialect.GetDialect()
	}

	var parts []string
	parts = append(parts, "CREATE DATABASE")

	if b.ifNotExists {
		parts = append(parts, "IF NOT EXISTS")
	}

	parts = append(parts, dialect.QuoteIdent(b.name))

	// Add charset if specified
	if b.charset != "" {
		parts = append(parts, "CHARACTER SET", b.charset)
	}

	// Add collation if specified
	if b.collation != "" {
		parts = append(parts, "COLLATE", b.collation)
	}

	// Add custom options
	for _, opt := range b.options {
		parts = append(parts, opt.Name)
		if opt.Value != "" {
			parts = append(parts, opt.Value)
		}
	}

	sql := strings.Join(parts, " ")
	return sql, nil, nil
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *CreateDatabaseBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return sqldebug.InterpolateSQL(sql, args).GetUnsafeString()
}
