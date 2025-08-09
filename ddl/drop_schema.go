package ddl

import (
	"errors"
	"github.com/sprylic/sqltk/sqldebug"
	"strings"

	"github.com/sprylic/sqltk/sqldialect"
)

// DropSchemaBuilder builds DROP SCHEMA statements.
type DropSchemaBuilder struct {
	name     string
	ifExists bool
	cascade  bool
	restrict bool
	err      error
	dialect  sqldialect.Dialect
}

// DropSchema creates a new DropSchemaBuilder for the given schema name.
func DropSchema(name string) *DropSchemaBuilder {
	if name == "" {
		return &DropSchemaBuilder{err: errors.New("schema name is required")}
	}
	return &DropSchemaBuilder{name: name}
}

// IfExists adds IF EXISTS clause to prevent errors if schema doesn't exist.
func (b *DropSchemaBuilder) IfExists() *DropSchemaBuilder {
	if b.err != nil {
		return b
	}
	b.ifExists = true
	return b
}

// Cascade adds CASCADE clause to drop dependent objects.
func (b *DropSchemaBuilder) Cascade() *DropSchemaBuilder {
	if b.err != nil {
		return b
	}
	b.cascade = true
	b.restrict = false
	return b
}

// Restrict adds RESTRICT clause to prevent dropping if dependent objects exist.
func (b *DropSchemaBuilder) Restrict() *DropSchemaBuilder {
	if b.err != nil {
		return b
	}
	b.restrict = true
	b.cascade = false
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *DropSchemaBuilder) WithDialect(d sqldialect.Dialect) *DropSchemaBuilder {
	b.dialect = d
	return b
}

// Build builds the SQL DROP SCHEMA query and returns the query string, arguments, and error if any.
func (b *DropSchemaBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.name == "" {
		return "", nil, errors.New("schema name is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = sqldialect.GetDialect()
	}

	var parts []string
	parts = append(parts, "DROP SCHEMA")

	if b.ifExists {
		parts = append(parts, "IF EXISTS")
	}

	parts = append(parts, dialect.QuoteIdent(b.name))

	// Add CASCADE or RESTRICT
	if b.cascade {
		parts = append(parts, "CASCADE")
	} else if b.restrict {
		parts = append(parts, "RESTRICT")
	}

	sql := strings.Join(parts, " ")
	return sql, nil, nil
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *DropSchemaBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return sqldebug.InterpolateSQL(sql, args).GetUnsafeString()
}
