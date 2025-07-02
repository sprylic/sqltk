package cqb

import (
	"errors"
	"strings"
)

// Raw marks a string as raw SQL to be included directly in the query.
type Raw string

// SelectBuilder builds SQL SELECT queries.
type SelectBuilder struct {
	tableClauseInterface
	distinct    bool
	columns     []interface{} // string, Raw, or *SelectBuilder
	joinClauses []string
	whereClause
	groupBy     []string
	groupByRaw  []string
	havingParam []string
	havingRaw   []string
	havingArgs  []interface{}
	orderBy     []string
	orderByRaw  []string
	limitSet    bool
	limit       int
	offsetSet   bool
	offset      int
	dialect     Dialect // per-builder dialect, if set
}

// Distinct sets the DISTINCT flag for the SELECT query.
func (b *SelectBuilder) Distinct() *SelectBuilder {
	b.distinct = true
	return b
}

// Select creates a new SelectBuilder with the given columns. Columns can be string, Raw, or *SelectBuilder (for subqueries).
func Select(columns ...interface{}) *SelectBuilder {
	return &SelectBuilder{columns: columns}
}

// AddField adds columns to the query. Columns can be string, Raw, or *SelectBuilder (for subqueries).
func (b *SelectBuilder) AddField(fields ...interface{}) *SelectBuilder {
	b.columns = append(b.columns, fields...)
	return b
}

// From sets the table for the SELECT query. Accepts string, Raw, or *SelectBuilder (for subqueries).
func (b *SelectBuilder) From(table interface{}) *SelectBuilder {
	b.SetTable(table)
	return b
}

// Where adds a WHERE clause to the query. Accepts either a condition string (with optional args) or a Raw type.
func (b *SelectBuilder) Where(cond interface{}, args ...interface{}) *SelectBuilder {
	b.whereClause.Where(cond, args...)
	return b
}

// WhereEqual adds a WHERE clause for equality (column = value).
func (b *SelectBuilder) WhereEqual(column string, value interface{}) *SelectBuilder {
	b.whereClause.WhereEqual(column, value)
	return b
}

// WhereNotEqual adds a WHERE clause for inequality (column != value).
func (b *SelectBuilder) WhereNotEqual(column string, value interface{}) *SelectBuilder {
	b.whereClause.WhereNotEqual(column, value)
	return b
}

// GroupBy adds a GROUP BY clause. Accepts either a column string or Raw.
func (b *SelectBuilder) GroupBy(expr interface{}) *SelectBuilder {
	if b.whereClause.err != nil || b.tableClauseInterface.err != nil {
		return b
	}
	switch c := expr.(type) {
	case Raw:
		b.groupByRaw = append(b.groupByRaw, string(c))
	case string:
		b.groupBy = append(b.groupBy, c)
	default:
		b.whereClause.err = errors.New("GroupBy: expr must be string or sq.Raw")
	}
	return b
}

// Having adds a HAVING clause. Accepts either a condition string (with optional args) or Raw.
func (b *SelectBuilder) Having(cond interface{}, args ...interface{}) *SelectBuilder {
	if b.whereClause.err != nil || b.tableClauseInterface.err != nil {
		return b
	}
	switch c := cond.(type) {
	case Raw:
		b.havingRaw = append(b.havingRaw, string(c))
	case string:
		b.havingParam = append(b.havingParam, c)
		b.havingArgs = append(b.havingArgs, args...)
	default:
		b.whereClause.err = errors.New("Having: cond must be string or sq.Raw")
	}
	return b
}

// OrderBy adds an ORDER BY clause. Accepts either a column string or Raw.
func (b *SelectBuilder) OrderBy(expr interface{}) *SelectBuilder {
	if b.whereClause.err != nil || b.tableClauseInterface.err != nil {
		return b
	}
	switch c := expr.(type) {
	case Raw:
		b.orderByRaw = append(b.orderByRaw, string(c))
	case string:
		b.orderBy = append(b.orderBy, c)
	default:
		b.whereClause.err = errors.New("OrderBy: expr must be string or sq.Raw")
	}
	return b
}

// JoinBuilder is used for fluent JOIN ... ON ... chaining.
type JoinBuilder struct {
	parent    *SelectBuilder
	joinType  string
	joinTable string
}

// Join starts an INNER JOIN clause. Call .On(left, right) to specify the ON condition.
// Example: q := cqb.Select("u.id").From("users u").Join("orders o").On("o.user_id", "u.id")
func (b *SelectBuilder) Join(table string) *JoinBuilder {
	return &JoinBuilder{parent: b, joinType: "JOIN", joinTable: table}
}

// LeftJoin starts a LEFT JOIN clause. Call .On(left, right) to specify the ON condition.
func (b *SelectBuilder) LeftJoin(table string) *JoinBuilder {
	return &JoinBuilder{parent: b, joinType: "LEFT JOIN", joinTable: table}
}

// RightJoin starts a RIGHT JOIN clause. Call .On(left, right) to specify the ON condition.
func (b *SelectBuilder) RightJoin(table string) *JoinBuilder {
	return &JoinBuilder{parent: b, joinType: "RIGHT JOIN", joinTable: table}
}

// FullJoin starts a FULL JOIN clause. Call .On(left, right) to specify the ON condition.
func (b *SelectBuilder) FullJoin(table string) *JoinBuilder {
	return &JoinBuilder{parent: b, joinType: "FULL JOIN", joinTable: table}
}

// On finalizes the JOIN ... ON ... clause and returns the parent SelectBuilder.
func (jb *JoinBuilder) On(left, right string) *SelectBuilder {
	clause := jb.joinType + " " + jb.joinTable + " ON " + left + " = " + right
	jb.parent.joinClauses = append(jb.parent.joinClauses, clause)
	return jb.parent
}

// Limit sets a LIMIT clause.
func (b *SelectBuilder) Limit(n int) *SelectBuilder {
	b.limitSet = true
	b.limit = n
	return b
}

// Offset sets an OFFSET clause.
func (b *SelectBuilder) Offset(n int) *SelectBuilder {
	b.offsetSet = true
	b.offset = n
	return b
}

// AliasExpr represents an aliased SQL expression (column, subquery, or table).
type AliasExpr struct {
	Expr  interface{}
	Alias string
}

// Alias creates an aliased SQL expression for use in columns or FROM.
func Alias(expr interface{}, alias string) AliasExpr {
	return AliasExpr{Expr: expr, Alias: alias}
}

// WithDialect sets the dialect for this builder instance.
func (b *SelectBuilder) WithDialect(d Dialect) *SelectBuilder {
	b.dialect = d
	return b
}

// Build builds the SQL query and returns the query string, arguments, and error if any invalid type is encountered.
func (b *SelectBuilder) Build() (string, []interface{}, error) {
	if b.tableClauseInterface.err != nil {
		return "", nil, b.tableClauseInterface.err
	}
	if b.whereClause.err != nil {
		return "", nil, b.whereClause.err
	}
	var sb strings.Builder
	var err error
	args := []interface{}{}

	dialect := b.dialect
	if dialect == nil {
		dialect = getDialect()
	}
	placeholderIdx := 1

	sb.WriteString("SELECT ")
	if b.distinct {
		sb.WriteString("DISTINCT ")
	}
	if len(b.columns) == 0 {
		sb.WriteString("*")
	} else {
		for i, col := range b.columns {
			if i > 0 {
				sb.WriteString(", ")
			}
			switch c := col.(type) {
			case string:
				sb.WriteString(dialect.QuoteIdent(c))
			case Raw:
				sb.WriteString(string(c))
			case *SelectBuilder:
				subSQL, subArgs, subErr := c.Build()
				if subErr != nil {
					err = subErr
				}
				sb.WriteString("(")
				sb.WriteString(subSQL)
				sb.WriteString(")")
				args = append(args, subArgs...)
			case AliasExpr:
				switch expr := c.Expr.(type) {
				case *SelectBuilder:
					subSQL, subArgs, subErr := expr.Build()
					if subErr != nil {
						err = subErr
					}
					sb.WriteString("(")
					sb.WriteString(subSQL)
					sb.WriteString(") AS ")
					sb.WriteString(dialect.QuoteIdent(c.Alias))
					args = append(args, subArgs...)
				case string:
					sb.WriteString(dialect.QuoteIdent(expr))
					sb.WriteString(" AS ")
					sb.WriteString(dialect.QuoteIdent(c.Alias))
				case Raw:
					sb.WriteString(string(expr))
					sb.WriteString(" AS ")
					sb.WriteString(dialect.QuoteIdent(c.Alias))
				default:
					err = errors.New("Alias: expr must be string, sq.Raw, or *SelectBuilder")
				}
			default:
				err = errors.New("Select: column must be string, sq.Raw, *SelectBuilder, or sq.AliasExpr")
			}
		}
	}
	sb.WriteString(" FROM ")
	switch t := b.tableClauseInterface.table.(type) {
	case string:
		sb.WriteString(dialect.QuoteIdent(t))
	case Raw:
		sb.WriteString(string(t))
	case *SelectBuilder:
		subSQL, subArgs, subErr := t.Build()
		if subErr != nil {
			err = subErr
		}
		sb.WriteString("(")
		sb.WriteString(subSQL)
		sb.WriteString(")")
		args = append(args, subArgs...)
	case AliasExpr:
		switch expr := t.Expr.(type) {
		case *SelectBuilder:
			subSQL, subArgs, subErr := expr.Build()
			if subErr != nil {
				err = subErr
			}
			sb.WriteString("(")
			sb.WriteString(subSQL)
			sb.WriteString(") AS ")
			sb.WriteString(dialect.QuoteIdent(t.Alias))
			args = append(args, subArgs...)
		case string:
			sb.WriteString(dialect.QuoteIdent(expr))
			sb.WriteString(" AS ")
			sb.WriteString(dialect.QuoteIdent(t.Alias))
		case Raw:
			sb.WriteString(string(expr))
			sb.WriteString(" AS ")
			sb.WriteString(dialect.QuoteIdent(t.Alias))
		default:
			err = errors.New("Alias: expr must be string, sq.Raw, or *SelectBuilder")
		}
	default:
		err = errors.New("From: table must be string, sq.Raw, *SelectBuilder, or sq.AliasExpr")
	}

	if len(b.joinClauses) > 0 {
		sb.WriteString(" ")
		sb.WriteString(strings.Join(b.joinClauses, " "))
	}

	whereSQL, whereArgs := b.whereClause.buildWhereSQL(dialect, &placeholderIdx)
	if whereSQL != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	var groupBys []string
	if len(b.groupBy) > 0 {
		for _, g := range b.groupBy {
			groupBys = append(groupBys, dialect.QuoteIdent(g))
		}
	}
	if len(b.groupByRaw) > 0 {
		groupBys = append(groupBys, b.groupByRaw...)
	}
	if len(groupBys) > 0 {
		sb.WriteString(" GROUP BY ")
		sb.WriteString(strings.Join(groupBys, ", "))
	}

	var havings []string
	if len(b.havingParam) > 0 {
		havings = append(havings, b.havingParam...)
	}
	if len(b.havingRaw) > 0 {
		havings = append(havings, b.havingRaw...)
	}
	if len(havings) > 0 {
		sb.WriteString(" HAVING ")
		havingSQL := strings.Join(havings, " AND ")
		for strings.Contains(havingSQL, "?") && dialect.Placeholder(0) != "?" {
			havingSQL = strings.Replace(havingSQL, "?", dialect.Placeholder(placeholderIdx), 1)
			placeholderIdx++
		}
		sb.WriteString(havingSQL)
		args = append(args, b.havingArgs...)
	}

	var orderBys []string
	if len(b.orderBy) > 0 {
		for _, o := range b.orderBy {
			orderBys = append(orderBys, dialect.QuoteIdent(o))
		}
	}
	if len(b.orderByRaw) > 0 {
		orderBys = append(orderBys, b.orderByRaw...)
	}
	if len(orderBys) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.Join(orderBys, ", "))
	}

	if b.limitSet {
		sb.WriteString(" LIMIT ")
		sb.WriteString(intToString(b.limit))
	}
	if b.offsetSet {
		sb.WriteString(" OFFSET ")
		sb.WriteString(intToString(b.offset))
	}

	if err != nil {
		return sb.String(), args, err
	}
	return sb.String(), args, nil
}

// intToString is a helper to convert int to string without importing strconv for this small use case.
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	bp := len(b)
	for n > 0 {
		bp--
		b[bp] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		bp--
		b[bp] = '-'
	}
	return string(b[bp:])
}

// MustSQL is a helper for internal use to get SQL for subqueries in joins (ignores args and errors).
func (b *SelectBuilder) MustSQL() string {
	sql, _, _ := b.Build()
	return sql
}

// SelectFragment is a function that composes or modifies a SelectBuilder.
type SelectFragment func(*SelectBuilder) *SelectBuilder

// Compose applies one or more SelectFragment functions to the builder.
// Example:
//
//	isActive := func(b *SelectBuilder) *SelectBuilder { return b.Where("active = ?", true) }
//	q := Select("id").From("users").Compose(isActive)
func (b *SelectBuilder) Compose(fns ...SelectFragment) *SelectBuilder {
	for _, fn := range fns {
		b = fn(b)
	}
	return b
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *SelectBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return InterpolateSQL(sql, args)
}
