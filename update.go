package sqltk

import (
	"errors"
	"strings"
)

// UpdateBuilder builds SQL UPDATE queries.
type UpdateBuilder struct {
	tableClauseString
	sets    []string
	setArgs []interface{}
	whereClause
	dialect Dialect // per-builder dialect, if set
}

// Update creates a new UpdateBuilder for the given table.
func Update(table string) *UpdateBuilder {
	b := &UpdateBuilder{}
	b.SetTable(table)
	return b
}

func (t *UpdateBuilder) SetTable(table string) {
	if table == "" {
		t.tableClauseString.err = errors.New("table must be set")
	} else {
		t.table = table
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

// WithDialect sets the dialect for this builder instance.
func (b *UpdateBuilder) WithDialect(d Dialect) *UpdateBuilder {
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
		dialect = GetDialect()
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
	return InterpolateSQL(sql, args).GetUnsafeString()
}
