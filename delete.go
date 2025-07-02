package cqb

import (
	"errors"
	"strings"
)

// DeleteBuilder builds SQL DELETE queries.
type DeleteBuilder struct {
	table      string
	whereParam []string
	whereRaw   []string
	whereArgs  []interface{}
	err        error
}

// Delete creates a new DeleteBuilder for the given table.
func Delete(table string) *DeleteBuilder {
	return &DeleteBuilder{table: table}
}

// Where adds a WHERE clause. Accepts either a condition string (with optional args) or a Raw type.
func (b *DeleteBuilder) Where(cond interface{}, args ...interface{}) *DeleteBuilder {
	if b.err != nil {
		return b
	}
	switch c := cond.(type) {
	case Raw:
		b.whereRaw = append(b.whereRaw, string(c))
	case string:
		b.whereParam = append(b.whereParam, c)
		b.whereArgs = append(b.whereArgs, args...)
	default:
		b.err = errors.New("Where: cond must be string or sq.Raw")
	}
	return b
}

// WhereEqual adds a WHERE clause for equality (column = value).
func (b *DeleteBuilder) WhereEqual(column string, value interface{}) *DeleteBuilder {
	return b.Where(column+" = ?", value)
}

// WhereNotEqual adds a WHERE clause for inequality (column != value).
func (b *DeleteBuilder) WhereNotEqual(column string, value interface{}) *DeleteBuilder {
	return b.Where(column+" != ?", value)
}

// Build builds the SQL DELETE query and returns the query string, arguments, and error if any.
func (b *DeleteBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.table == "" {
		return "", nil, errors.New("Delete: table must be set")
	}

	dialect := getDialect()
	placeholderIdx := 1

	var sb strings.Builder
	args := []interface{}{}

	sb.WriteString("DELETE FROM ")
	sb.WriteString(dialect.QuoteIdent(b.table))

	var wheres []string
	if len(b.whereParam) > 0 {
		wheres = append(wheres, b.whereParam...)
	}
	if len(b.whereRaw) > 0 {
		wheres = append(wheres, b.whereRaw...)
	}
	if len(wheres) > 0 {
		sb.WriteString(" WHERE ")
		whereSQL := strings.Join(wheres, " AND ")
		for strings.Contains(whereSQL, "?") && dialect.Placeholder(0) != "?" {
			whereSQL = strings.Replace(whereSQL, "?", dialect.Placeholder(placeholderIdx), 1)
			placeholderIdx++
		}
		sb.WriteString(whereSQL)
		args = append(args, b.whereArgs...)
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
