package ddl

import (
	"errors"
	"strings"

	"github.com/sprylic/sqltk/shared"
)

// DropViewBuilder builds DROP VIEW statements.
type DropViewBuilder struct {
	viewName string
	ifExists bool
	cascade  bool
	restrict bool
	err      error
	dialect  shared.Dialect
}

// DropView creates a new DROP VIEW builder.
func DropView(viewName string) *DropViewBuilder {
	if viewName == "" {
		return &DropViewBuilder{err: errors.New("view name is required")}
	}
	return &DropViewBuilder{
		viewName: viewName,
	}
}

// IfExists adds IF EXISTS to the DROP VIEW statement.
func (b *DropViewBuilder) IfExists() *DropViewBuilder {
	if b.err != nil {
		return b
	}
	b.ifExists = true
	return b
}

// Cascade adds CASCADE to the DROP VIEW statement.
func (b *DropViewBuilder) Cascade() *DropViewBuilder {
	if b.err != nil {
		return b
	}
	if b.restrict {
		b.err = errors.New("cannot use both CASCADE and RESTRICT")
		return b
	}
	b.cascade = true
	return b
}

// Restrict adds RESTRICT to the DROP VIEW statement.
func (b *DropViewBuilder) Restrict() *DropViewBuilder {
	if b.err != nil {
		return b
	}
	if b.cascade {
		b.err = errors.New("cannot use both CASCADE and RESTRICT")
		return b
	}
	b.restrict = true
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *DropViewBuilder) WithDialect(d shared.Dialect) *DropViewBuilder {
	if b.err != nil {
		return b
	}
	b.dialect = d
	return b
}

// Build builds the SQL DROP VIEW query and returns the query string, arguments, and error if any.
func (b *DropViewBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.viewName == "" {
		return "", nil, errors.New("view name is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = shared.GetDialect() // Use global dialect instead of defaulting to MySQL
	}

	var sb strings.Builder
	args := []interface{}{}

	// DROP VIEW
	sb.WriteString("DROP VIEW ")
	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}
	sb.WriteString(dialect.QuoteIdent(b.viewName))

	// CASCADE or RESTRICT
	if b.cascade {
		sb.WriteString(" CASCADE")
	} else if b.restrict {
		sb.WriteString(" RESTRICT")
	}

	return sb.String(), args, nil
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *DropViewBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return shared.InterpolateSQL(sql, args)
}
