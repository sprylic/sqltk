package cqb

import (
	"errors"
	"strings"
)

// CreateIndexBuilder builds CREATE INDEX statements.
type CreateIndexBuilder struct {
	indexName   string
	tableName   string
	columns     []string
	unique      bool
	ifNotExists bool
	dialect     Dialect
	err         error
}

// CreateIndex starts building a CREATE INDEX statement.
func CreateIndex(indexName string) *CreateIndexBuilder {
	if indexName == "" {
		return &CreateIndexBuilder{err: errors.New("index name is required")}
	}
	return &CreateIndexBuilder{
		indexName: indexName,
	}
}

// On sets the table name for the index.
func (b *CreateIndexBuilder) On(tableName string) *CreateIndexBuilder {
	if b.err != nil {
		return b
	}
	if tableName == "" {
		b.err = errors.New("table name is required")
		return b
	}
	b.tableName = tableName
	return b
}

// Columns sets the columns for the index.
func (b *CreateIndexBuilder) Columns(columns ...string) *CreateIndexBuilder {
	if b.err != nil {
		return b
	}
	if len(columns) == 0 {
		b.err = errors.New("at least one column must be specified")
		return b
	}
	b.columns = columns
	return b
}

// Unique makes the index unique.
func (b *CreateIndexBuilder) Unique() *CreateIndexBuilder {
	if b.err != nil {
		return b
	}
	b.unique = true
	return b
}

// IfNotExists adds IF NOT EXISTS to the statement.
func (b *CreateIndexBuilder) IfNotExists() *CreateIndexBuilder {
	if b.err != nil {
		return b
	}
	b.ifNotExists = true
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *CreateIndexBuilder) WithDialect(d Dialect) *CreateIndexBuilder {
	if b.err != nil {
		return b
	}
	b.dialect = d
	return b
}

// Build builds the SQL CREATE INDEX query and returns the query string, arguments, and error if any.
func (b *CreateIndexBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.indexName == "" {
		return "", nil, errors.New("index name is required")
	}
	if b.tableName == "" {
		return "", nil, errors.New("table name is required")
	}
	if len(b.columns) == 0 {
		return "", nil, errors.New("at least one column must be specified")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = getDialect()
	}

	var sb strings.Builder
	args := []interface{}{}

	// CREATE [UNIQUE] INDEX
	sb.WriteString("CREATE ")
	if b.unique {
		sb.WriteString("UNIQUE ")
	}
	sb.WriteString("INDEX ")

	// IF NOT EXISTS
	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}

	// Index name
	sb.WriteString(dialect.QuoteIdent(b.indexName))

	// ON table
	sb.WriteString(" ON ")
	sb.WriteString(dialect.QuoteIdent(b.tableName))

	// Columns
	sb.WriteString(" (")
	quotedColumns := make([]string, len(b.columns))
	for i, col := range b.columns {
		quotedColumns[i] = dialect.QuoteIdent(col)
	}
	sb.WriteString(strings.Join(quotedColumns, ", "))
	sb.WriteString(")")

	return sb.String(), args, nil
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *CreateIndexBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return InterpolateSQL(sql, args)
}
