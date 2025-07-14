package ddl

import (
	"errors"
	"strings"

	"github.com/sprylic/sqltk/shared"
)

// DropDatabaseBuilder builds DROP DATABASE statements.
type DropDatabaseBuilder struct {
	name     string
	ifExists bool
	cascade  bool
	err      error
	dialect  shared.Dialect
}

// DropDatabase creates a new DropDatabaseBuilder for the given database name.
func DropDatabase(name string) *DropDatabaseBuilder {
	if name == "" {
		return &DropDatabaseBuilder{err: errors.New("database name is required")}
	}
	return &DropDatabaseBuilder{name: name}
}

// IfExists adds IF EXISTS clause to prevent errors if database doesn't exist.
func (b *DropDatabaseBuilder) IfExists() *DropDatabaseBuilder {
	if b.err != nil {
		return b
	}
	b.ifExists = true
	return b
}

// Cascade adds CASCADE clause to drop dependent objects (PostgreSQL).
func (b *DropDatabaseBuilder) Cascade() *DropDatabaseBuilder {
	if b.err != nil {
		return b
	}
	b.cascade = true
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *DropDatabaseBuilder) WithDialect(d shared.Dialect) *DropDatabaseBuilder {
	b.dialect = d
	return b
}

// Build builds the SQL DROP DATABASE query and returns the query string, arguments, and error if any.
func (b *DropDatabaseBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.name == "" {
		return "", nil, errors.New("database name is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = shared.GetDialect()
	}

	var parts []string
	parts = append(parts, "DROP DATABASE")

	if b.ifExists {
		parts = append(parts, "IF EXISTS")
	}

	parts = append(parts, dialect.QuoteIdent(b.name))

	// Add CASCADE for PostgreSQL
	if b.cascade && dialect == shared.Postgres() {
		parts = append(parts, "CASCADE")
	}

	sql := strings.Join(parts, " ")
	return sql, nil, nil
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *DropDatabaseBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return shared.InterpolateSQL(sql, args).GetUnsafeString()
}
