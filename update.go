package sqltk

import (
	"errors"
	"github.com/sprylic/sqltk/raw"
	"github.com/sprylic/sqltk/sqldebug"
	"strings"

	"github.com/sprylic/sqltk/sqldialect"
)

// UpdateBuilder builds SQL UPDATE queries.
type UpdateBuilder struct {
	tableClauseString
	sets    []string
	setArgs []interface{}
	whereClause
	dialect sqldialect.Dialect // per-builder dialect, if set
}

// Update creates a new UpdateBuilder for the given table.
func Update(table string) *UpdateBuilder {
	b := &UpdateBuilder{}
	b.SetTable(table)
	return b
}

func (b *UpdateBuilder) SetTable(table string) {
	if table == "" {
		b.tableClauseString.err = errors.New("table must be set")
	} else {
		b.table = table
	}
}

// Set adds a SET clause. Accepts column name and value.
func (b *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	if b.whereClause.err != nil {
		return b
	}
	b.sets = append(b.sets, column+" = ?")
	b.setArgs = append(b.setArgs, value)
	return b
}

// SetRaw adds a raw SET clause (use with caution).
func (b *UpdateBuilder) SetRaw(expr string) *UpdateBuilder {
	if b.whereClause.err != nil {
		return b
	}
	b.sets = append(b.sets, expr)
	return b
}

// Where adds a WHERE clause. Accepts a Condition.
func (b *UpdateBuilder) Where(cond Condition, args ...interface{}) *UpdateBuilder {
	b.whereClause.Where(cond, args...)
	return b
}

// WhereEqual adds a WHERE clause for equality (column = value).
func (b *UpdateBuilder) WhereEqual(column string, value interface{}) *UpdateBuilder {
	b.Where(NewStringCondition(column+" = ?", value))
	return b
}

// WhereNotEqual adds a WHERE clause for inequality (column != value).
func (b *UpdateBuilder) WhereNotEqual(column string, value interface{}) *UpdateBuilder {
	b.Where(NewStringCondition(column+" != ?", value))
	return b
}

// WhereNull adds a WHERE clause for NULL check (column IS NULL).
func (b *UpdateBuilder) WhereNull(column string) *UpdateBuilder {
	b.Where(NewCond().IsNull(column))
	return b
}

// WhereNotNull adds a WHERE clause for NOT NULL check (column IS NOT NULL).
func (b *UpdateBuilder) WhereNotNull(column string) *UpdateBuilder {
	b.Where(NewCond().IsNotNull(column))
	return b
}

// WhereGreaterThan adds a WHERE clause for greater than comparison (column > value).
func (b *UpdateBuilder) WhereGreaterThan(column string, value interface{}) *UpdateBuilder {
	b.Where(NewCond().GreaterThan(column, value))
	return b
}

// WhereGreaterThanOrEqual adds a WHERE clause for greater than or equal comparison (column >= value).
func (b *UpdateBuilder) WhereGreaterThanOrEqual(column string, value interface{}) *UpdateBuilder {
	b.Where(NewCond().GreaterThanOrEqual(column, value))
	return b
}

// WhereLessThan adds a WHERE clause for less than comparison (column < value).
func (b *UpdateBuilder) WhereLessThan(column string, value interface{}) *UpdateBuilder {
	b.Where(NewCond().LessThan(column, value))
	return b
}

// WhereLessThanOrEqual adds a WHERE clause for less than or equal comparison (column <= value).
func (b *UpdateBuilder) WhereLessThanOrEqual(column string, value interface{}) *UpdateBuilder {
	b.Where(NewCond().LessThanOrEqual(column, value))
	return b
}

// WhereLike adds a WHERE clause for LIKE pattern matching (column LIKE pattern).
func (b *UpdateBuilder) WhereLike(column string, pattern string) *UpdateBuilder {
	b.Where(NewCond().Like(column, pattern))
	return b
}

// WhereNotLike adds a WHERE clause for NOT LIKE pattern matching (column NOT LIKE pattern).
func (b *UpdateBuilder) WhereNotLike(column string, pattern string) *UpdateBuilder {
	b.Where(NewCond().NotLike(column, pattern))
	return b
}

// WhereIn adds a WHERE clause for IN condition (column IN (values...)).
func (b *UpdateBuilder) WhereIn(column string, values ...interface{}) *UpdateBuilder {
	b.Where(NewCond().In(column, values...))
	return b
}

// WhereNotIn adds a WHERE clause for NOT IN condition (column NOT IN (values...)).
func (b *UpdateBuilder) WhereNotIn(column string, values ...interface{}) *UpdateBuilder {
	b.Where(NewCond().NotIn(column, values...))
	return b
}

// WhereBetween adds a WHERE clause for BETWEEN condition (column BETWEEN min AND max).
func (b *UpdateBuilder) WhereBetween(column string, min, max interface{}) *UpdateBuilder {
	b.Where(NewCond().Between(column, min, max))
	return b
}

// WhereNotBetween adds a WHERE clause for NOT BETWEEN condition (column NOT BETWEEN min AND max).
func (b *UpdateBuilder) WhereNotBetween(column string, min, max interface{}) *UpdateBuilder {
	b.Where(NewCond().NotBetween(column, min, max))
	return b
}

// WhereExists adds a WHERE clause for EXISTS condition (EXISTS (subquery)).
func (b *UpdateBuilder) WhereExists(subquery interface{}) *UpdateBuilder {
	b.Where(NewCond().Exists(subquery))
	return b
}

// WhereNotExists adds a WHERE clause for NOT EXISTS condition (NOT EXISTS (subquery)).
func (b *UpdateBuilder) WhereNotExists(subquery interface{}) *UpdateBuilder {
	b.Where(NewCond().NotExists(subquery))
	return b
}

// WhereColsEqual adds a WHERE clause for column equality (column1 = column2).
func (b *UpdateBuilder) WhereColsEqual(column1, column2 string) *UpdateBuilder {
	b.Where(raw.Raw(column1 + " = " + column2))
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *UpdateBuilder) WithDialect(d sqldialect.Dialect) *UpdateBuilder {
	b.dialect = d
	return b
}

// Build builds the SQL UPDATE query and returns the query string, arguments, and error if any.
func (b *UpdateBuilder) Build() (string, []interface{}, error) {
	if b.tableClauseString.err != nil {
		return "", nil, b.tableClauseString.err
	}
	if b.whereClause.err != nil {
		return "", nil, b.whereClause.err
	}
	if b.tableClauseString.table == "" {
		return "", nil, errors.New("Update: table must be set")
	}
	if len(b.sets) == 0 {
		return "", nil, errors.New("Update: at least one SET clause must be set")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = sqldialect.GetDialect()
	}
	placeholderIdx := 1

	var sb strings.Builder
	args := append([]interface{}{}, b.setArgs...)

	sb.WriteString("UPDATE ")
	sb.WriteString(dialect.QuoteIdent(b.tableClauseString.table))
	sb.WriteString(" SET ")

	setSQL := strings.Join(b.sets, ", ")
	if dialect.Placeholder(0) != "?" {
		for strings.Contains(setSQL, "?") && dialect.Placeholder(0) != "?" {
			setSQL = strings.Replace(setSQL, "?", dialect.Placeholder(placeholderIdx), 1)
			placeholderIdx++
		}
	}

	sb.WriteString(setSQL)

	whereSQL, whereArgs := b.buildWhereSQL(dialect, &placeholderIdx)
	if whereSQL != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	return sb.String(), args, nil
}

// PostgresUpdateBuilder extends UpdateBuilder with RETURNING support for Postgres.
type PostgresUpdateBuilder struct {
	*UpdateBuilder
	returning []string
}

// NewPostgresUpdate creates a new PostgresUpdateBuilder for the given table.
func NewPostgresUpdate(table string) *PostgresUpdateBuilder {
	return &PostgresUpdateBuilder{UpdateBuilder: Update(table)}
}

// Returning adds a RETURNING clause (Postgres only).
func (b *PostgresUpdateBuilder) Returning(cols ...string) *PostgresUpdateBuilder {
	b.returning = append(b.returning, cols...)
	return b
}

// Build builds the SQL UPDATE query with RETURNING (if set) and returns the query string, arguments, and error if any.
func (b *PostgresUpdateBuilder) Build() (string, []interface{}, error) {
	sql, args, err := b.UpdateBuilder.Build()
	if err != nil {
		return sql, args, err
	}
	if len(b.returning) > 0 {
		sql += " RETURNING " + strings.Join(b.returning, ", ")
	}
	return sql, args, nil
}

// Example usage:
//   pq := sq.NewPostgresUpdate("users").Set("name", "Alice").Where("id = ?", 1).Returning("id")
//   sql, args, err := pq.Build()

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *UpdateBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return sqldebug.InterpolateSQL(sql, args).GetUnsafeString()
}
