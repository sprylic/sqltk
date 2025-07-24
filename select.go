package sqltk

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sprylic/sqltk/sqlfunc"
)

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

// Where adds a WHERE clause to the query. Accepts a Condition.
func (b *SelectBuilder) Where(cond Condition, args ...interface{}) *SelectBuilder {
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

// WhereNull adds a WHERE clause for NULL check (column IS NULL).
func (b *SelectBuilder) WhereNull(column string) *SelectBuilder {
	b.Where(NewCond().IsNull(column))
	return b
}

// WhereNotNull adds a WHERE clause for NOT NULL check (column IS NOT NULL).
func (b *SelectBuilder) WhereNotNull(column string) *SelectBuilder {
	b.Where(NewCond().IsNotNull(column))
	return b
}

// WhereGreaterThan adds a WHERE clause for greater than comparison (column > value).
func (b *SelectBuilder) WhereGreaterThan(column string, value interface{}) *SelectBuilder {
	b.Where(NewCond().GreaterThan(column, value))
	return b
}

// WhereGreaterThanOrEqual adds a WHERE clause for greater than or equal comparison (column >= value).
func (b *SelectBuilder) WhereGreaterThanOrEqual(column string, value interface{}) *SelectBuilder {
	b.Where(NewCond().GreaterThanOrEqual(column, value))
	return b
}

// WhereLessThan adds a WHERE clause for less than comparison (column < value).
func (b *SelectBuilder) WhereLessThan(column string, value interface{}) *SelectBuilder {
	b.Where(NewCond().LessThan(column, value))
	return b
}

// WhereLessThanOrEqual adds a WHERE clause for less than or equal comparison (column <= value).
func (b *SelectBuilder) WhereLessThanOrEqual(column string, value interface{}) *SelectBuilder {
	b.Where(NewCond().LessThanOrEqual(column, value))
	return b
}

// WhereLike adds a WHERE clause for LIKE pattern matching (column LIKE pattern).
func (b *SelectBuilder) WhereLike(column string, pattern string) *SelectBuilder {
	b.Where(NewCond().Like(column, pattern))
	return b
}

// WhereNotLike adds a WHERE clause for NOT LIKE pattern matching (column NOT LIKE pattern).
func (b *SelectBuilder) WhereNotLike(column string, pattern string) *SelectBuilder {
	b.Where(NewCond().NotLike(column, pattern))
	return b
}

// WhereIn adds a WHERE clause for IN condition (column IN (values...)).
func (b *SelectBuilder) WhereIn(column string, values ...interface{}) *SelectBuilder {
	b.Where(NewCond().In(column, values...))
	return b
}

// WhereNotIn adds a WHERE clause for NOT IN condition (column NOT IN (values...)).
func (b *SelectBuilder) WhereNotIn(column string, values ...interface{}) *SelectBuilder {
	b.Where(NewCond().NotIn(column, values...))
	return b
}

// WhereBetween adds a WHERE clause for BETWEEN condition (column BETWEEN min AND max).
func (b *SelectBuilder) WhereBetween(column string, min, max interface{}) *SelectBuilder {
	b.Where(NewCond().Between(column, min, max))
	return b
}

// WhereNotBetween adds a WHERE clause for NOT BETWEEN condition (column NOT BETWEEN min AND max).
func (b *SelectBuilder) WhereNotBetween(column string, min, max interface{}) *SelectBuilder {
	b.Where(NewCond().NotBetween(column, min, max))
	return b
}

// WhereExists adds a WHERE clause for EXISTS condition (EXISTS (subquery)).
func (b *SelectBuilder) WhereExists(subquery interface{}) *SelectBuilder {
	b.Where(NewCond().Exists(subquery))
	return b
}

// WhereNotExists adds a WHERE clause for NOT EXISTS condition (NOT EXISTS (subquery)).
func (b *SelectBuilder) WhereNotExists(subquery interface{}) *SelectBuilder {
	b.Where(NewCond().NotExists(subquery))
	return b
}

// WhereColsEqual adds a WHERE clause for column equality (column1 = column2).
func (b *SelectBuilder) WhereColsEqual(column1, column2 string) *SelectBuilder {
	b.Where(Raw(column1 + " = " + column2))
	return b
}

// GroupBy adds a GROUP BY clause. Accepts either a column string or Raw.
func (b *SelectBuilder) GroupBy(expr ...interface{}) *SelectBuilder {
	if b.whereClause.err != nil || b.tableClauseInterface.err != nil {
		return b
	}

	for _, e := range expr {
		switch c := e.(type) {
		case sqlfunc.SqlFunc:
			b.groupByRaw = append(b.groupByRaw, string(c))
		case Raw:
			b.groupByRaw = append(b.groupByRaw, string(c))
		case string:
			b.groupBy = append(b.groupBy, c)
		default:
			b.whereClause.err = errors.New("GroupBy: expr must be string or sq.Raw")
		}
	}
	return b
}

// Having adds a HAVING clause. Accepts a Condition.
func (b *SelectBuilder) Having(cond Condition, args ...interface{}) *SelectBuilder {
	if b.whereClause.err != nil || b.tableClauseInterface.err != nil {
		return b
	}

	sql, condArgs, err := cond.BuildCondition()
	if err != nil {
		b.whereClause.err = fmt.Errorf("Having: condition error: %w", err)
		return b
	}
	if sql != "" {
		b.havingParam = append(b.havingParam, sql)
		b.havingArgs = append(b.havingArgs, condArgs...)
	}
	return b
}

// OrderBy adds an ORDER BY clause. Accepts either a column string or Raw.
func (b *SelectBuilder) OrderBy(expr interface{}) *SelectBuilder {
	if b.whereClause.err != nil || b.tableClauseInterface.err != nil {
		return b
	}
	switch c := expr.(type) {
	case sqlfunc.SqlFunc:
		b.orderByRaw = append(b.orderByRaw, string(c))
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
	joinTable interface{}
	err       error
}

// Join starts an INNER JOIN clause. Accepts a table, subquery, or alias.
//   - string: table name (optionally with alias, e.g. "users u")
//   - Raw: raw SQL for the table
//   - *SelectBuilder: subquery as table
//   - AliasExpr: alias for a table or subquery (use sqltk.Alias)
//
// Example usage:
//
//	Join("orders o")
//	Join(sqltk.Alias("orders", "o"))
//	Join(sqltk.Alias(sqltk.Select("id").From("orders"), "o"))
func (b *SelectBuilder) Join(table interface{}) *JoinBuilder {
	return &JoinBuilder{parent: b, joinType: "JOIN", joinTable: table}
}

// LeftJoin starts a LEFT JOIN clause. Accepts a table, subquery, or alias.
//   - string: table name (optionally with alias)
//   - Raw: raw SQL for the table
//   - *SelectBuilder: subquery as table
//   - AliasExpr: alias for a table or subquery (use sqltk.Alias)
//
// Example usage:
//
//	LeftJoin("orders o")
//	LeftJoin(sqltk.Alias("orders", "o"))
//	LeftJoin(sqltk.Alias(sqltk.Select("id").From("orders"), "o"))
func (b *SelectBuilder) LeftJoin(table interface{}) *JoinBuilder {
	return &JoinBuilder{parent: b, joinType: "LEFT JOIN", joinTable: table}
}

// RightJoin starts a RIGHT JOIN clause. Accepts a table, subquery, or alias.
//   - string: table name (optionally with alias)
//   - Raw: raw SQL for the table
//   - *SelectBuilder: subquery as table
//   - AliasExpr: alias for a table or subquery (use sqltk.Alias)
//
// Example usage:
//
//	RightJoin("orders o")
//	RightJoin(sqltk.Alias("orders", "o"))
//	RightJoin(sqltk.Alias(sqltk.Select("id").From("orders"), "o"))
func (b *SelectBuilder) RightJoin(table interface{}) *JoinBuilder {
	return &JoinBuilder{parent: b, joinType: "RIGHT JOIN", joinTable: table}
}

// FullJoin starts a FULL JOIN clause. Accepts a table, subquery, or alias.
//   - string: table name (optionally with alias)
//   - Raw: raw SQL for the table
//   - *SelectBuilder: subquery as table
//   - AliasExpr: alias for a table or subquery (use sqltk.Alias)
//
// Example usage:
//
//	FullJoin("orders o")
//	FullJoin(sqltk.Alias("orders", "o"))
//	FullJoin(sqltk.Alias(sqltk.Select("id").From("orders"), "o"))
func (b *SelectBuilder) FullJoin(table interface{}) *JoinBuilder {
	return &JoinBuilder{parent: b, joinType: "FULL JOIN", joinTable: table}
}

// On finalizes the JOIN ... ON ... clause and returns the parent SelectBuilder.
func (jb *JoinBuilder) On(left, right string) *SelectBuilder {
	if jb.err != nil {
		jb.parent.whereClause.err = jb.err
		return jb.parent
	}

	clause := jb.joinType + " "
	dialect := jb.parent.dialect
	if dialect == nil {
		dialect = GetDialect()
	}

	switch t := jb.joinTable.(type) {
	case string:
		clause += dialect.QuoteIdent(t)
	case Raw:
		clause += string(t)
	case *SelectBuilder:
		subSQL, subArgs, subErr := t.Build()
		if subErr != nil {
			jb.parent.whereClause.err = fmt.Errorf("join subquery error: %w", subErr)
			return jb.parent
		}
		clause += "(" + subSQL + ")"
		// Store the subquery args in the parent's whereClause for later use
		jb.parent.whereClause.whereArgs = append(jb.parent.whereClause.whereArgs, subArgs...)
	case AliasExpr:
		switch expr := t.Expr.(type) {
		case *SelectBuilder:
			subSQL, subArgs, subErr := expr.Build()
			if subErr != nil {
				jb.parent.whereClause.err = fmt.Errorf("join alias subquery error: %w", subErr)
				return jb.parent
			}
			clause += "(" + subSQL + ") AS " + t.Alias
			// Store the subquery args in the parent's whereClause for later use
			jb.parent.whereClause.whereArgs = append(jb.parent.whereClause.whereArgs, subArgs...)
		case string:
			clause += dialect.QuoteIdent(expr) + " AS " + t.Alias
		case Raw:
			clause += string(expr) + " AS " + t.Alias
		default:
			jb.parent.whereClause.err = fmt.Errorf("join alias: expr must be string, Raw, or *SelectBuilder (got %T)", expr)
			return jb.parent
		}
	default:
		jb.parent.whereClause.err = fmt.Errorf("join: table must be string, Raw, *SelectBuilder, or AliasExpr (got %T)", t)
		return jb.parent
	}

	clause += " ON " + left + " = " + right
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
		dialect = GetDialect()
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
				// Handle expressions with aliases (e.g., "COUNT(*) as count")
				if strings.Contains(strings.ToUpper(c), " AS ") {
					parts := strings.SplitN(c, " AS ", 2)
					if len(parts) == 2 {
						expr := strings.TrimSpace(parts[0])
						alias := strings.TrimSpace(parts[1])

						// Handle table-qualified column names in expressions
						if strings.Contains(expr, ".") {
							exprParts := strings.Split(expr, ".")
							for i, part := range exprParts {
								if i > 0 {
									sb.WriteString(".")
								}
								sb.WriteString(dialect.QuoteIdent(strings.TrimSpace(part)))
							}
						} else {
							sb.WriteString(expr)
						}
						sb.WriteString(" AS ")
						sb.WriteString(alias)
					} else {
						sb.WriteString(c)
					}
				} else if strings.Contains(c, ".") {
					// Handle table-qualified column names (e.g., "table.column")
					parts := strings.Split(c, ".")
					for i, part := range parts {
						if i > 0 {
							sb.WriteString(".")
						}
						sb.WriteString(dialect.QuoteIdent(strings.TrimSpace(part)))
					}
				} else {
					sb.WriteString(dialect.QuoteIdent(c))
				}
			case Raw:
				sb.WriteString(string(c))
			case sqlfunc.SqlFunc:
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
					sb.WriteString(c.Alias)
					args = append(args, subArgs...)
				case string:
					// Handle table-qualified column names in AliasExpr
					if strings.Contains(expr, ".") {
						parts := strings.Split(expr, ".")
						for i, part := range parts {
							if i > 0 {
								sb.WriteString(".")
							}
							sb.WriteString(dialect.QuoteIdent(strings.TrimSpace(part)))
						}
					} else {
						sb.WriteString(dialect.QuoteIdent(expr))
					}
					sb.WriteString(" AS ")
					sb.WriteString(c.Alias)
				case Raw:
					sb.WriteString(string(expr))
					sb.WriteString(" AS ")
					sb.WriteString(c.Alias)
				case sqlfunc.SqlFunc:
					sb.WriteString(string(expr))
					sb.WriteString(" AS ")
					sb.WriteString(c.Alias)
				default:
					err = errors.New("Alias: expr must be string, sq.Raw, *SelectBuilder, or sqlfunc.SqlFunc")
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
	case sqlfunc.SqlFunc:
		sb.WriteString(string(t))
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
			sb.WriteString(t.Alias)
			args = append(args, subArgs...)
		case string:
			sb.WriteString(dialect.QuoteIdent(expr))
			sb.WriteString(" AS ")
			sb.WriteString(t.Alias)
		case Raw:
			sb.WriteString(string(expr))
			sb.WriteString(" AS ")
			sb.WriteString(t.Alias)
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
	} else {
		// Even if there's no WHERE clause, we need to include any args from subqueries
		args = append(args, whereArgs...)
	}

	var groupBys []string
	if len(b.groupBy) > 0 {
		for _, g := range b.groupBy {
			if strings.Contains(g, ".") {
				// Handle table-qualified column names (e.g., "table.column")
				parts := strings.Split(g, ".")
				var quoted string
				for i, part := range parts {
					if i > 0 {
						quoted += "."
					}
					quoted += dialect.QuoteIdent(strings.TrimSpace(part))
				}
				groupBys = append(groupBys, quoted)
			} else {
				groupBys = append(groupBys, dialect.QuoteIdent(g))
			}
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
			// Handle expressions like 'total_amount DESC'
			if idx := strings.IndexAny(o, " "); idx > 0 {
				col := o[:idx]
				dir := strings.TrimSpace(o[idx+1:])
				if strings.Contains(col, ".") {
					parts := strings.Split(col, ".")
					var quoted string
					for i, part := range parts {
						if i > 0 {
							quoted += "."
						}
						quoted += dialect.QuoteIdent(strings.TrimSpace(part))
					}
					orderBys = append(orderBys, quoted+" "+dir)
				} else {
					orderBys = append(orderBys, dialect.QuoteIdent(col)+" "+dir)
				}
			} else if strings.Contains(o, ".") {
				parts := strings.Split(o, ".")
				var quoted string
				for i, part := range parts {
					if i > 0 {
						quoted += "."
					}
					quoted += dialect.QuoteIdent(strings.TrimSpace(part))
				}
				orderBys = append(orderBys, quoted)
			} else {
				orderBys = append(orderBys, dialect.QuoteIdent(o))
			}
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

// GetColumns returns all columns used in the builder and all subquery builders.
func (b *SelectBuilder) GetColumns() []string {
	var cols []string
	for _, col := range b.columns {
		switch col.(type) {
		case string:
			cols = append(cols, col.(string))
		case Raw:
			cols = append(cols, string(col.(Raw)))
		case sqlfunc.SqlFunc:
			cols = append(cols, string(col.(sqlfunc.SqlFunc)))
		case *SelectBuilder:
			cols = append(cols, col.(*SelectBuilder).GetColumns()...)
		case AliasExpr:
			cols = append(cols, col.(AliasExpr).Alias)
		}
	}
	return cols
}

// Compose combines this SelectBuilder with one or more other SelectBuilder instances.
// This merges columns, joins, where conditions, group by, having, order by, limit, and offset.
// The first builder's table and dialect are preserved.
// Example:
//
//	q1 := Select("id", "name").From("users").Where("active = ?", true)
//	q2 := Select("email").From("users").Where("verified = ?", true)
//	q := q1.Compose(q2) // Combines columns and merges where conditions
func (b *SelectBuilder) Compose(builders ...*SelectBuilder) *SelectBuilder {
	for _, other := range builders {
		if other == nil {
			continue
		}

		// Merge columns
		b.columns = append(b.columns, other.columns...)

		// Merge joins
		b.joinClauses = append(b.joinClauses, other.joinClauses...)

		// Merge where conditions
		if other.whereClause.err != nil {
			b.whereClause.err = other.whereClause.err
		} else {
			b.whereClause.whereParam = append(b.whereClause.whereParam, other.whereClause.whereParam...)
			b.whereClause.whereRaw = append(b.whereClause.whereRaw, other.whereClause.whereRaw...)
			b.whereClause.whereArgs = append(b.whereClause.whereArgs, other.whereClause.whereArgs...)
		}

		// Merge group by
		b.groupBy = append(b.groupBy, other.groupBy...)
		b.groupByRaw = append(b.groupByRaw, other.groupByRaw...)

		// Merge having
		b.havingParam = append(b.havingParam, other.havingParam...)
		b.havingRaw = append(b.havingRaw, other.havingRaw...)
		b.havingArgs = append(b.havingArgs, other.havingArgs...)

		// Merge order by
		b.orderBy = append(b.orderBy, other.orderBy...)
		b.orderByRaw = append(b.orderByRaw, other.orderByRaw...)

		// Use the most restrictive limit/offset
		if other.limitSet && (!b.limitSet || other.limit < b.limit) {
			b.limitSet = true
			b.limit = other.limit
		}
		if other.offsetSet && (!b.offsetSet || other.offset > b.offset) {
			b.offsetSet = true
			b.offset = other.offset
		}

		// Preserve distinct if any builder has it
		if other.distinct {
			b.distinct = true
		}
	}

	return b
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *SelectBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return InterpolateSQL(sql, args).GetUnsafeString()
}
