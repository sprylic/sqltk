package cqb

import (
	"errors"
	"strings"
)

// UpdateBuilder builds SQL UPDATE queries.
type UpdateBuilder struct {
	table      string
	sets       []string
	setArgs    []interface{}
	whereParam []string
	whereRaw   []string
	whereArgs  []interface{}
	err        error
}

// Update creates a new UpdateBuilder for the given table.
func Update(table string) *UpdateBuilder {
	return &UpdateBuilder{table: table}
}

// Set adds a SET clause. Accepts column name and value.
func (b *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	if b.err != nil {
		return b
	}
	b.sets = append(b.sets, column+" = ?")
	b.setArgs = append(b.setArgs, value)
	return b
}

// SetRaw adds a raw SET clause (use with caution).
func (b *UpdateBuilder) SetRaw(expr string) *UpdateBuilder {
	if b.err != nil {
		return b
	}
	b.sets = append(b.sets, expr)
	return b
}

// Where adds a WHERE clause. Accepts either a condition string (with optional args) or a Raw type.
func (b *UpdateBuilder) Where(cond interface{}, args ...interface{}) *UpdateBuilder {
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

// Build builds the SQL UPDATE query and returns the query string, arguments, and error if any.
func (b *UpdateBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.table == "" {
		return "", nil, errors.New("Update: table must be set")
	}
	if len(b.sets) == 0 {
		return "", nil, errors.New("Update: at least one SET clause must be set")
	}

	dialect := getDialect()
	placeholderIdx := 1

	var sb strings.Builder
	args := append([]interface{}{}, b.setArgs...)

	sb.WriteString("UPDATE ")
	sb.WriteString(dialect.QuoteIdent(b.table))
	sb.WriteString(" SET ")

	setSQL := strings.Join(b.sets, ", ")
	if dialect.Placeholder(0) != "?" {
		for strings.Contains(setSQL, "?") && dialect.Placeholder(0) != "?" {
			setSQL = strings.Replace(setSQL, "?", dialect.Placeholder(placeholderIdx), 1)
			placeholderIdx++
		}
	}

	sb.WriteString(setSQL)

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
