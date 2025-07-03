package stk

import (
	"errors"
	"strings"
)

// InsertBuilder builds SQL INSERT queries.
type InsertBuilder struct {
	table   string
	columns []string
	values  [][]interface{}
	err     error
	dialect Dialect // per-builder dialect, if set
}

// Insert creates a new InsertBuilder for the given table.
func Insert(table string) *InsertBuilder {
	return &InsertBuilder{table: table}
}

// Columns sets the columns for the INSERT statement.
func (b *InsertBuilder) Columns(cols ...string) *InsertBuilder {
	if b.err != nil {
		return b
	}
	b.columns = append([]string{}, cols...)
	return b
}

// Values adds a row of values to insert. Call multiple times for multi-row insert.
func (b *InsertBuilder) Values(vals ...interface{}) *InsertBuilder {
	if b.err != nil {
		return b
	}
	if len(vals) != len(b.columns) {
		b.err = errors.New("Values: number of values must match number of columns")
		return b
	}
	b.values = append(b.values, vals)
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *InsertBuilder) WithDialect(d Dialect) *InsertBuilder {
	b.dialect = d
	return b
}

// Build builds the SQL INSERT query and returns the query string, arguments, and error if any.
func (b *InsertBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.table == "" {
		return "", nil, errors.New("Insert: table must be set")
	}
	if len(b.columns) == 0 {
		return "", nil, errors.New("Insert: columns must be set")
	}
	if len(b.values) == 0 {
		return "", nil, errors.New("Insert: at least one row of values must be set")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = GetDialect()
	}
	placeholderIdx := 1

	var sb strings.Builder
	args := make([]interface{}, 0, len(b.values)*len(b.columns))

	sb.WriteString("INSERT INTO ")
	sb.WriteString(dialect.QuoteIdent(b.table))
	sb.WriteString(" (")
	for i, col := range b.columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(dialect.QuoteIdent(col))
	}
	sb.WriteString(") VALUES ")

	for i, row := range b.values {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("(")
		for j := range row {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(dialect.Placeholder(placeholderIdx))
			placeholderIdx++
			args = append(args, row[j])
		}
		sb.WriteString(")")
	}

	return sb.String(), args, nil
}

// PostgresInsertBuilder extends InsertBuilder with RETURNING support for Postgres.
type PostgresInsertBuilder struct {
	*InsertBuilder
	returning []string
}

// NewPostgresInsert creates a new PostgresInsertBuilder for the given table.
func NewPostgresInsert(table string) *PostgresInsertBuilder {
	return &PostgresInsertBuilder{InsertBuilder: Insert(table).WithDialect(Postgres())}
}

// Returning adds a RETURNING clause (Postgres only).
func (b *PostgresInsertBuilder) Returning(cols ...string) *PostgresInsertBuilder {
	b.returning = append(b.returning, cols...)
	return b
}

// Build builds the SQL INSERT query with RETURNING (if set) and returns the query string, arguments, and error if any.
func (b *PostgresInsertBuilder) Build() (string, []interface{}, error) {
	sql, args, err := b.InsertBuilder.Build()
	if err != nil {
		return sql, args, err
	}
	if len(b.returning) > 0 {
		sql += " RETURNING " + strings.Join(b.returning, ", ")
	}
	return sql, args, nil
}

// Example usage:
//   pq := sq.NewPostgresInsert("users").Columns("name").Values("Alice").Returning("id")
//   sql, args, err := pq.Build()

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *InsertBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return InterpolateSQL(sql, args)
}
