package ddl

import (
	"errors"
	"strings"

	"github.com/sprylic/sqltk/shared"
)

// CreateIndexBuilder builds SQL CREATE INDEX queries.
type CreateIndexBuilder struct {
	indexName   string
	tableName   string
	columns     []string
	unique      bool
	ifNotExists bool
	err         error
	dialect     shared.Dialect
}

// CreateIndex creates a new CreateIndexBuilder for the given index and table.
func CreateIndex(indexName, tableName string) *CreateIndexBuilder {
	if indexName == "" {
		return &CreateIndexBuilder{err: errors.New("index name is required")}
	}
	if tableName == "" {
		return &CreateIndexBuilder{err: errors.New("table name is required")}
	}
	return &CreateIndexBuilder{
		indexName: indexName,
		tableName: tableName,
		columns:   make([]string, 0),
	}
}

// Columns adds columns to the index.
func (b *CreateIndexBuilder) Columns(columns ...string) *CreateIndexBuilder {
	if b.err != nil {
		return b
	}
	if len(columns) == 0 {
		b.err = errors.New("at least one column is required")
		return b
	}
	for _, col := range columns {
		if col == "" {
			b.err = errors.New("column name cannot be empty")
			return b
		}
	}
	b.columns = append(b.columns, columns...)
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

// IfNotExists adds IF NOT EXISTS to the CREATE INDEX statement.
func (b *CreateIndexBuilder) IfNotExists() *CreateIndexBuilder {
	if b.err != nil {
		return b
	}
	b.ifNotExists = true
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *CreateIndexBuilder) WithDialect(d shared.Dialect) *CreateIndexBuilder {
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
		return "", nil, errors.New("at least one column is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = shared.GetDialect() // Use global dialect instead of defaulting to MySQL
	}

	var sb strings.Builder
	args := []interface{}{}

	// CREATE INDEX
	sb.WriteString("CREATE ")
	if b.unique {
		sb.WriteString("UNIQUE ")
	}
	sb.WriteString("INDEX ")
	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	sb.WriteString(dialect.QuoteIdent(b.indexName))
	sb.WriteString(" ON ")
	sb.WriteString(dialect.QuoteIdent(b.tableName))

	// Columns
	quotedCols := make([]string, len(b.columns))
	for i, col := range b.columns {
		quotedCols[i] = dialect.QuoteIdent(col)
	}
	sb.WriteString(" (")
	sb.WriteString(strings.Join(quotedCols, ", "))
	sb.WriteString(")")

	return sb.String(), args, nil
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *CreateIndexBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return shared.InterpolateSQL(sql, args).GetUnsafeString()
}
