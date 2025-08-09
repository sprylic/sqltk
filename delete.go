package sqltk

import (
	"errors"
	"strings"

	"github.com/sprylic/sqltk/raw"
	"github.com/sprylic/sqltk/sqldebug"
	"github.com/sprylic/sqltk/sqldialect"
)

// DeleteBuilder builds SQL DELETE queries.
type DeleteBuilder struct {
	tableClauseString
	whereClause
	dialect sqldialect.Dialect // per-builder dialect, if set
}

// Delete creates a new DeleteBuilder for the given table.
func Delete(table string) *DeleteBuilder {
	b := &DeleteBuilder{}
	b.SetTable(table)
	return b
}

// Where adds a WHERE clause. Accepts a Condition.
func (b *DeleteBuilder) Where(cond Condition, args ...interface{}) *DeleteBuilder {
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

// WhereNull adds a WHERE clause for NULL check (column IS NULL).
func (b *DeleteBuilder) WhereNull(column string) *DeleteBuilder {
	b.Where(NewCond().IsNull(column))
	return b
}

// WhereNotNull adds a WHERE clause for NOT NULL check (column IS NOT NULL).
func (b *DeleteBuilder) WhereNotNull(column string) *DeleteBuilder {
	b.Where(NewCond().IsNotNull(column))
	return b
}

// WhereGreaterThan adds a WHERE clause for greater than comparison (column > value).
func (b *DeleteBuilder) WhereGreaterThan(column string, value interface{}) *DeleteBuilder {
	b.Where(NewCond().GreaterThan(column, value))
	return b
}

// WhereGreaterThanOrEqual adds a WHERE clause for greater than or equal comparison (column >= value).
func (b *DeleteBuilder) WhereGreaterThanOrEqual(column string, value interface{}) *DeleteBuilder {
	b.Where(NewCond().GreaterThanOrEqual(column, value))
	return b
}

// WhereLessThan adds a WHERE clause for less than comparison (column < value).
func (b *DeleteBuilder) WhereLessThan(column string, value interface{}) *DeleteBuilder {
	b.Where(NewCond().LessThan(column, value))
	return b
}

// WhereLessThanOrEqual adds a WHERE clause for less than or equal comparison (column <= value).
func (b *DeleteBuilder) WhereLessThanOrEqual(column string, value interface{}) *DeleteBuilder {
	b.Where(NewCond().LessThanOrEqual(column, value))
	return b
}

// WhereLike adds a WHERE clause for LIKE pattern matching (column LIKE pattern).
func (b *DeleteBuilder) WhereLike(column string, pattern string) *DeleteBuilder {
	b.Where(NewCond().Like(column, pattern))
	return b
}

// WhereNotLike adds a WHERE clause for NOT LIKE pattern matching (column NOT LIKE pattern).
func (b *DeleteBuilder) WhereNotLike(column string, pattern string) *DeleteBuilder {
	b.Where(NewCond().NotLike(column, pattern))
	return b
}

// WhereIn adds a WHERE clause for IN condition (column IN (values...)).
func (b *DeleteBuilder) WhereIn(column string, values ...interface{}) *DeleteBuilder {
	b.Where(NewCond().In(column, values...))
	return b
}

// WhereNotIn adds a WHERE clause for NOT IN condition (column NOT IN (values...)).
func (b *DeleteBuilder) WhereNotIn(column string, values ...interface{}) *DeleteBuilder {
	b.Where(NewCond().NotIn(column, values...))
	return b
}

// WhereBetween adds a WHERE clause for BETWEEN condition (column BETWEEN min AND max).
func (b *DeleteBuilder) WhereBetween(column string, min, max interface{}) *DeleteBuilder {
	b.Where(NewCond().Between(column, min, max))
	return b
}

// WhereNotBetween adds a WHERE clause for NOT BETWEEN condition (column NOT BETWEEN min AND max).
func (b *DeleteBuilder) WhereNotBetween(column string, min, max interface{}) *DeleteBuilder {
	b.Where(NewCond().NotBetween(column, min, max))
	return b
}

// WhereExists adds a WHERE clause for EXISTS condition (EXISTS (subquery)).
func (b *DeleteBuilder) WhereExists(subquery interface{}) *DeleteBuilder {
	b.Where(NewCond().Exists(subquery))
	return b
}

// WhereNotExists adds a WHERE clause for NOT EXISTS condition (NOT EXISTS (subquery)).
func (b *DeleteBuilder) WhereNotExists(subquery interface{}) *DeleteBuilder {
	b.Where(NewCond().NotExists(subquery))
	return b
}

// WhereColsEqual adds a WHERE clause for column equality (column1 = column2).
func (b *DeleteBuilder) WhereColsEqual(column1, column2 string) *DeleteBuilder {
	b.Where(raw.Raw(column1 + " = " + column2))
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *DeleteBuilder) WithDialect(d sqldialect.Dialect) *DeleteBuilder {
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
		dialect = sqldialect.GetDialect()
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
	return sqldebug.InterpolateSQL(sql, args).GetUnsafeString()
}
