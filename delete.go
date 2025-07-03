package cqb

import (
	"errors"
	"strings"
)

// DeleteBuilder builds SQL DELETE queries.
type DeleteBuilder struct {
	tableClauseString
	whereClause
	dialect Dialect // per-builder dialect, if set
}

// Delete creates a new DeleteBuilder for the given table.
func Delete(table string) *DeleteBuilder {
	b := &DeleteBuilder{}
	b.SetTable(table)
	return b
}

// Where adds a WHERE clause. Accepts either a condition string (with optional args) or a Raw type.
func (b *DeleteBuilder) Where(cond interface{}, args ...interface{}) *DeleteBuilder {
	b.whereClause.Where(cond, args...)
	return b
}

// WhereEqual adds a WHERE clause for equality (column = value).
func (b *DeleteBuilder) WhereEqual(column string, value interface{}) *DeleteBuilder {
	b.whereClause.WhereEqual(column, value)
	return b
}

// WhereNotEqual adds a WHERE clause for inequality (column != value).
func (b *DeleteBuilder) WhereNotEqual(column string, value interface{}) *DeleteBuilder {
	b.whereClause.WhereNotEqual(column, value)
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *DeleteBuilder) WithDialect(d Dialect) *DeleteBuilder {
	b.dialect = d
	return b
}

// Build builds the SQL DELETE query and returns the query string, arguments, and error if any.
func (b *DeleteBuilder) Build() (string, []interface{}, error) {
	if b.tableClauseString.err != nil {
		return "", nil, b.tableClauseString.err
	}
	if b.whereClause.err != nil {
		return "", nil, b.whereClause.err
	}
	if b.tableClauseString.table == "" {
		return "", nil, errors.New("Delete: table must be set")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = GetDialect()
	}
	placeholderIdx := 1

	var sb strings.Builder
	args := []interface{}{}

	sb.WriteString("DELETE FROM ")
	sb.WriteString(dialect.QuoteIdent(b.tableClauseString.table))

	whereSQL, whereArgs := b.whereClause.buildWhereSQL(dialect, &placeholderIdx)
	if whereSQL != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	return sb.String(), args, nil
}

// PostgresDeleteBuilder extends DeleteBuilder with RETURNING support for Postgres.
type PostgresDeleteBuilder struct {
	*DeleteBuilder
	returning []string
}

// NewPostgresDelete creates a new PostgresDeleteBuilder for the given table.
func NewPostgresDelete(table string) *PostgresDeleteBuilder {
	return &PostgresDeleteBuilder{DeleteBuilder: Delete(table)}
}

// Returning adds a RETURNING clause (Postgres only).
func (b *PostgresDeleteBuilder) Returning(cols ...string) *PostgresDeleteBuilder {
	b.returning = append(b.returning, cols...)
	return b
}

// Build builds the SQL DELETE query with RETURNING (if set) and returns the query string, arguments, and error if any.
func (b *PostgresDeleteBuilder) Build() (string, []interface{}, error) {
	sql, args, err := b.DeleteBuilder.Build()
	if err != nil {
		return sql, args, err
	}
	if len(b.returning) > 0 {
		sql += " RETURNING " + strings.Join(b.returning, ", ")
	}
	return sql, args, nil
}

// Example usage:
//   pq := sq.NewPostgresDelete("users").Where("id = ?", 1).Returning("id")
//   sql, args, err := pq.Build()

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *DeleteBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return InterpolateSQL(sql, args)
}
