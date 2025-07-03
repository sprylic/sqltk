package stk

import (
	"fmt"
	"strings"
)

// Condition represents a SQL condition that can be used in WHERE or HAVING clauses.
// This interface ensures type safety and prevents unsafe queries.
type Condition interface {
	// BuildCondition returns the SQL condition string and arguments.
	// This method is used internally by the query builders.
	BuildCondition() (string, []interface{}, error)
}

// ICondition is an alias for Condition for better readability in method signatures.
type ICondition = Condition

// StringCondition wraps a string condition for type safety.
type StringCondition struct {
	SQL  string
	Args []interface{}
}

// NewStringCondition creates a new StringCondition from a SQL string and arguments.
// This is the preferred way to create string-based conditions for better type safety.
func NewStringCondition(sql string, args ...interface{}) *StringCondition {
	return &StringCondition{SQL: sql, Args: args}
}

// BuildCondition implements the Condition interface.
func (sc *StringCondition) BuildCondition() (string, []interface{}, error) {
	return sc.SQL, sc.Args, nil
}

// RawCondition wraps a Raw SQL condition for type safety.
type RawCondition struct {
	SQL Raw
}

// NewRawCondition creates a new RawCondition from Raw SQL.
func NewRawCondition(sql Raw) *RawCondition {
	return &RawCondition{SQL: sql}
}

// AsCondition converts a Raw to a Condition for use in Where/Having clauses.
func AsCondition(r Raw) Condition {
	return &RawCondition{SQL: r}
}

// BuildCondition implements the Condition interface.
func (rc *RawCondition) BuildCondition() (string, []interface{}, error) {
	return string(rc.SQL), nil, nil
}

// ConditionBuilder provides a fluent API for building SQL conditions.
type ConditionBuilder struct {
	parts   []string
	args    []interface{}
	err     error
	dialect Dialect
}

// BuildCondition implements the Condition interface.
func (c *ConditionBuilder) BuildCondition() (string, []interface{}, error) {
	return c.Build()
}

// NewCond creates a new ConditionBuilder.
func NewCond() *ConditionBuilder {
	return &ConditionBuilder{}
}

// WithDialect sets the dialect for this condition builder.
func (c *ConditionBuilder) WithDialect(d Dialect) *ConditionBuilder {
	c.dialect = d
	return c
}

// getDialect returns the dialect, using global if not set.
func (c *ConditionBuilder) getDialect() Dialect {
	if c.dialect != nil {
		return c.dialect
	}
	return GetDialect()
}

// Where adds a simple WHERE condition.
func (c *ConditionBuilder) Where(column string, operator string, value interface{}) *ConditionBuilder {
	if c.err != nil {
		return c
	}

	dialect := c.getDialect()
	var quotedCol string

	// Handle table-qualified column names (e.g., "table.column")
	if strings.Contains(column, ".") {
		parts := strings.Split(column, ".")
		quotedParts := make([]string, len(parts))
		for i, part := range parts {
			quotedParts[i] = dialect.QuoteIdent(strings.TrimSpace(part))
		}
		quotedCol = strings.Join(quotedParts, ".")
	} else {
		quotedCol = dialect.QuoteIdent(column)
	}

	if value == nil {
		switch operator {
		case "=":
			c.parts = append(c.parts, quotedCol+" IS NULL")
		case "!=", "<>":
			c.parts = append(c.parts, quotedCol+" IS NOT NULL")
		default:
			c.err = fmt.Errorf("invalid operator %q for NULL value", operator)
		}
		return c
	}

	c.parts = append(c.parts, quotedCol+" "+operator+" ?")
	c.args = append(c.args, value)
	return c
}

// Equal adds an equality condition (column = value).
func (c *ConditionBuilder) Equal(column string, value interface{}) *ConditionBuilder {
	return c.Where(column, "=", value)
}

// NotEqual adds an inequality condition (column != value).
func (c *ConditionBuilder) NotEqual(column string, value interface{}) *ConditionBuilder {
	return c.Where(column, "!=", value)
}

// GreaterThan adds a greater than condition (column > value).
func (c *ConditionBuilder) GreaterThan(column string, value interface{}) *ConditionBuilder {
	return c.Where(column, ">", value)
}

// GreaterThanOrEqual adds a greater than or equal condition (column >= value).
func (c *ConditionBuilder) GreaterThanOrEqual(column string, value interface{}) *ConditionBuilder {
	return c.Where(column, ">=", value)
}

// LessThan adds a less than condition (column < value).
func (c *ConditionBuilder) LessThan(column string, value interface{}) *ConditionBuilder {
	return c.Where(column, "<", value)
}

// LessThanOrEqual adds a less than or equal condition (column <= value).
func (c *ConditionBuilder) LessThanOrEqual(column string, value interface{}) *ConditionBuilder {
	return c.Where(column, "<=", value)
}

// Like adds a LIKE condition (column LIKE pattern).
func (c *ConditionBuilder) Like(column string, pattern string) *ConditionBuilder {
	return c.Where(column, "LIKE", pattern)
}

// NotLike adds a NOT LIKE condition (column NOT LIKE pattern).
func (c *ConditionBuilder) NotLike(column string, pattern string) *ConditionBuilder {
	return c.Where(column, "NOT LIKE", pattern)
}

// In adds an IN condition (column IN (values...)).
func (c *ConditionBuilder) In(column string, values ...interface{}) *ConditionBuilder {
	if c.err != nil {
		return c
	}

	if len(values) == 0 {
		c.err = fmt.Errorf("IN condition requires at least one value")
		return c
	}

	dialect := c.getDialect()
	var quotedCol string

	// Handle table-qualified column names (e.g., "table.column")
	if strings.Contains(column, ".") {
		parts := strings.Split(column, ".")
		quotedParts := make([]string, len(parts))
		for i, part := range parts {
			quotedParts[i] = dialect.QuoteIdent(strings.TrimSpace(part))
		}
		quotedCol = strings.Join(quotedParts, ".")
	} else {
		quotedCol = dialect.QuoteIdent(column)
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	c.parts = append(c.parts, quotedCol+" IN ("+strings.Join(placeholders, ", ")+")")
	c.args = append(c.args, values...)
	return c
}

// NotIn adds a NOT IN condition (column NOT IN (values...)).
func (c *ConditionBuilder) NotIn(column string, values ...interface{}) *ConditionBuilder {
	if c.err != nil {
		return c
	}

	if len(values) == 0 {
		c.err = fmt.Errorf("NOT IN condition requires at least one value")
		return c
	}

	dialect := c.getDialect()
	var quotedCol string

	// Handle table-qualified column names (e.g., "table.column")
	if strings.Contains(column, ".") {
		parts := strings.Split(column, ".")
		quotedParts := make([]string, len(parts))
		for i, part := range parts {
			quotedParts[i] = dialect.QuoteIdent(strings.TrimSpace(part))
		}
		quotedCol = strings.Join(quotedParts, ".")
	} else {
		quotedCol = dialect.QuoteIdent(column)
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	c.parts = append(c.parts, quotedCol+" NOT IN ("+strings.Join(placeholders, ", ")+")")
	c.args = append(c.args, values...)
	return c
}

// Between adds a BETWEEN condition (column BETWEEN min AND max).
func (c *ConditionBuilder) Between(column string, min, max interface{}) *ConditionBuilder {
	if c.err != nil {
		return c
	}

	dialect := c.getDialect()
	var quotedCol string

	// Handle table-qualified column names (e.g., "table.column")
	if strings.Contains(column, ".") {
		parts := strings.Split(column, ".")
		quotedParts := make([]string, len(parts))
		for i, part := range parts {
			quotedParts[i] = dialect.QuoteIdent(strings.TrimSpace(part))
		}
		quotedCol = strings.Join(quotedParts, ".")
	} else {
		quotedCol = dialect.QuoteIdent(column)
	}

	c.parts = append(c.parts, quotedCol+" BETWEEN ? AND ?")
	c.args = append(c.args, min, max)
	return c
}

// NotBetween adds a NOT BETWEEN condition (column NOT BETWEEN min AND max).
func (c *ConditionBuilder) NotBetween(column string, min, max interface{}) *ConditionBuilder {
	if c.err != nil {
		return c
	}

	dialect := c.getDialect()
	var quotedCol string

	// Handle table-qualified column names (e.g., "table.column")
	if strings.Contains(column, ".") {
		parts := strings.Split(column, ".")
		quotedParts := make([]string, len(parts))
		for i, part := range parts {
			quotedParts[i] = dialect.QuoteIdent(strings.TrimSpace(part))
		}
		quotedCol = strings.Join(quotedParts, ".")
	} else {
		quotedCol = dialect.QuoteIdent(column)
	}

	c.parts = append(c.parts, quotedCol+" NOT BETWEEN ? AND ?")
	c.args = append(c.args, min, max)
	return c
}

// IsNull adds an IS NULL condition (column IS NULL).
func (c *ConditionBuilder) IsNull(column string) *ConditionBuilder {
	if c.err != nil {
		return c
	}

	dialect := c.getDialect()
	var quotedCol string

	// Handle table-qualified column names (e.g., "table.column")
	if strings.Contains(column, ".") {
		parts := strings.Split(column, ".")
		quotedParts := make([]string, len(parts))
		for i, part := range parts {
			quotedParts[i] = dialect.QuoteIdent(strings.TrimSpace(part))
		}
		quotedCol = strings.Join(quotedParts, ".")
	} else {
		quotedCol = dialect.QuoteIdent(column)
	}

	c.parts = append(c.parts, quotedCol+" IS NULL")
	return c
}

// IsNotNull adds an IS NOT NULL condition (column IS NOT NULL).
func (c *ConditionBuilder) IsNotNull(column string) *ConditionBuilder {
	if c.err != nil {
		return c
	}

	dialect := c.getDialect()
	var quotedCol string

	// Handle table-qualified column names (e.g., "table.column")
	if strings.Contains(column, ".") {
		parts := strings.Split(column, ".")
		quotedParts := make([]string, len(parts))
		for i, part := range parts {
			quotedParts[i] = dialect.QuoteIdent(strings.TrimSpace(part))
		}
		quotedCol = strings.Join(quotedParts, ".")
	} else {
		quotedCol = dialect.QuoteIdent(column)
	}

	c.parts = append(c.parts, quotedCol+" IS NOT NULL")
	return c
}

// Exists adds an EXISTS condition (EXISTS (subquery)).
func (c *ConditionBuilder) Exists(subquery interface{}) *ConditionBuilder {
	if c.err != nil {
		return c
	}

	var sql string
	var args []interface{}
	var err error

	switch sq := subquery.(type) {
	case *SelectBuilder:
		sql, args, err = sq.Build()
		if err != nil {
			c.err = fmt.Errorf("exists subquery error: %w", err)
			return c
		}
	case Raw:
		sql = string(sq)
	default:
		c.err = fmt.Errorf("exists: subquery must be *SelectBuilder or Raw (got %T)", subquery)
		return c
	}

	c.parts = append(c.parts, "EXISTS ("+sql+")")
	c.args = append(c.args, args...)
	return c
}

// NotExists adds a NOT EXISTS condition (NOT EXISTS (subquery)).
func (c *ConditionBuilder) NotExists(subquery interface{}) *ConditionBuilder {
	if c.err != nil {
		return c
	}

	var sql string
	var args []interface{}
	var err error

	switch sq := subquery.(type) {
	case *SelectBuilder:
		sql, args, err = sq.Build()
		if err != nil {
			c.err = fmt.Errorf("not exists subquery error: %w", err)
			return c
		}
	case Raw:
		sql = string(sq)
	default:
		c.err = fmt.Errorf("not exists: subquery must be *SelectBuilder or Raw (got %T)", subquery)
		return c
	}

	c.parts = append(c.parts, "NOT EXISTS ("+sql+")")
	c.args = append(c.args, args...)
	return c
}

// Case adds a CASE WHEN condition.
func (c *ConditionBuilder) Case() *CaseBuilder {
	return &CaseBuilder{parent: c}
}

// Raw adds a raw SQL condition.
func (c *ConditionBuilder) Raw(sql string, args ...interface{}) *ConditionBuilder {
	if c.err != nil {
		return c
	}

	c.parts = append(c.parts, sql)
	c.args = append(c.args, args...)
	return c
}

// And combines conditions with AND.
func (c *ConditionBuilder) And(other *ConditionBuilder) *ConditionBuilder {
	if c.err != nil {
		return c
	}
	if other.err != nil {
		c.err = other.err
		return c
	}

	if len(other.parts) == 0 {
		return c
	}

	if len(c.parts) == 0 {
		c.parts = other.parts
		c.args = other.args
		return c
	}

	// Combine conditions with AND
	c.parts = append(c.parts, other.parts...)
	c.args = append(c.args, other.args...)
	return c
}

// Or combines conditions with OR.
func (c *ConditionBuilder) Or(other *ConditionBuilder) *ConditionBuilder {
	if c.err != nil {
		return c
	}
	if other.err != nil {
		c.err = other.err
		return c
	}

	if len(other.parts) == 0 {
		return c
	}

	if len(c.parts) == 0 {
		c.parts = other.parts
		c.args = other.args
		return c
	}

	// Wrap both sides in parentheses and combine with OR
	c.parts = []string{"(" + strings.Join(c.parts, " AND ") + ") OR (" + strings.Join(other.parts, " AND ") + ")"}
	c.args = append(c.args, other.args...)
	return c
}

// Build returns the SQL condition string and arguments.
func (c *ConditionBuilder) Build() (string, []interface{}, error) {
	if c.err != nil {
		return "", nil, c.err
	}

	if len(c.parts) == 0 {
		return "", nil, nil
	}

	return strings.Join(c.parts, " AND "), c.args, nil
}

// String returns the condition as a string (for debugging).
func (c *ConditionBuilder) String() string {
	sql, args, err := c.Build()
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	if len(args) == 0 {
		return sql
	}
	return InterpolateSQL(sql, args)
}

// CaseBuilder provides a fluent API for building CASE WHEN expressions.
type CaseBuilder struct {
	parent    *ConditionBuilder
	whenParts []string
	whenArgs  []interface{}
	elsePart  string
	elseArgs  []interface{}
	err       error
}

// When adds a WHEN clause to the CASE expression.
func (cb *CaseBuilder) When(condition interface{}, result interface{}) *CaseBuilder {
	if cb.err != nil {
		return cb
	}

	var condSQL string
	var condArgs []interface{}
	var err error

	switch c := condition.(type) {
	case *ConditionBuilder:
		condSQL, condArgs, err = c.Build()
		if err != nil {
			cb.err = fmt.Errorf("case when condition error: %w", err)
			return cb
		}
	case Raw:
		condSQL = string(c)
	case string:
		condSQL = c
	default:
		cb.err = fmt.Errorf("case when: condition must be *ConditionBuilder, Raw, or string (got %T)", condition)
		return cb
	}

	whenClause := "WHEN " + condSQL + " THEN ?"
	cb.whenParts = append(cb.whenParts, whenClause)
	cb.whenArgs = append(cb.whenArgs, condArgs...)
	cb.whenArgs = append(cb.whenArgs, result)

	return cb
}

// Else adds an ELSE clause to the CASE expression.
func (cb *CaseBuilder) Else(value interface{}) *CaseBuilder {
	if cb.err != nil {
		return cb
	}

	cb.elsePart = "ELSE ?"
	cb.elseArgs = []interface{}{value}
	return cb
}

// End finalizes the CASE expression and returns to the parent ConditionBuilder.
func (cb *CaseBuilder) End() *ConditionBuilder {
	if cb.err != nil {
		cb.parent.err = cb.err
		return cb.parent
	}

	if len(cb.whenParts) == 0 {
		cb.parent.err = fmt.Errorf("case expression must have at least one WHEN clause")
		return cb.parent
	}

	caseSQL := "CASE " + strings.Join(cb.whenParts, " ")
	if cb.elsePart != "" {
		caseSQL += " " + cb.elsePart
	}
	caseSQL += " END"

	cb.parent.parts = append(cb.parent.parts, caseSQL)
	cb.parent.args = append(cb.parent.args, cb.whenArgs...)
	cb.parent.args = append(cb.parent.args, cb.elseArgs...)

	return cb.parent
}
