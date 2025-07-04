package ddl

import (
	"errors"
	"strings"

	"github.com/sprylic/sqltk/shared"
)

// CreateSchemaBuilder builds CREATE SCHEMA statements.
type CreateSchemaBuilder struct {
	name          string
	ifNotExists   bool
	authorization string
	options       []SchemaOption
	err           error
	dialect       shared.Dialect
}

// SchemaOption represents a schema option.
type SchemaOption struct {
	Name  string
	Value string
}

// CreateSchema creates a new CreateSchemaBuilder for the given schema name.
func CreateSchema(name string) *CreateSchemaBuilder {
	if name == "" {
		return &CreateSchemaBuilder{err: errors.New("schema name is required")}
	}
	return &CreateSchemaBuilder{name: name}
}

// IfNotExists adds IF NOT EXISTS clause to prevent errors if schema already exists.
func (b *CreateSchemaBuilder) IfNotExists() *CreateSchemaBuilder {
	if b.err != nil {
		return b
	}
	b.ifNotExists = true
	return b
}

// Authorization sets the authorization (owner) for the schema.
func (b *CreateSchemaBuilder) Authorization(user string) *CreateSchemaBuilder {
	if b.err != nil {
		return b
	}
	b.authorization = user
	return b
}

// Option adds a custom schema option.
func (b *CreateSchemaBuilder) Option(name, value string) *CreateSchemaBuilder {
	if b.err != nil {
		return b
	}
	b.options = append(b.options, SchemaOption{Name: name, Value: value})
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *CreateSchemaBuilder) WithDialect(d shared.Dialect) *CreateSchemaBuilder {
	b.dialect = d
	return b
}

// Build builds the SQL CREATE SCHEMA query and returns the query string, arguments, and error if any.
func (b *CreateSchemaBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.name == "" {
		return "", nil, errors.New("schema name is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = shared.GetDialect()
	}

	var parts []string
	parts = append(parts, "CREATE SCHEMA")

	if b.ifNotExists {
		parts = append(parts, "IF NOT EXISTS")
	}

	parts = append(parts, dialect.QuoteIdent(b.name))

	// Add authorization if specified
	if b.authorization != "" {
		parts = append(parts, "AUTHORIZATION", dialect.QuoteIdent(b.authorization))
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
func (b *CreateSchemaBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return shared.InterpolateSQL(sql, args)
}
