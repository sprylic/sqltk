package ddl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sprylic/cqb/shared"
)

// CreateViewBuilder builds CREATE VIEW statements.
type CreateViewBuilder struct {
	viewName     string
	selectSQL    string
	orReplace    bool
	materialized bool // For future materialized view support
	err          error
	dialect      shared.Dialect
}

// CreateView creates a new CREATE VIEW builder.
func CreateView(viewName string) *CreateViewBuilder {
	if viewName == "" {
		return &CreateViewBuilder{err: errors.New("view name is required")}
	}
	return &CreateViewBuilder{
		viewName: viewName,
	}
}

// As sets the view definition using a builder or Raw.
// Expected arguments:
//   - Raw: for literal SQL strings (e.g., Raw("SELECT * FROM users"))
//   - Any value implementing Build() (string, []interface{}, error): for composable query builders
//     (e.g., SelectBuilder, or any custom builder with a Build method)
//
// Examples:
//   - As(Raw("SELECT id, name FROM users WHERE active = 1"))
//   - As(Select("id", "name").From("users").Where("active = 1"))
func (b *CreateViewBuilder) As(query interface{}) *CreateViewBuilder {
	if b.err != nil {
		return b
	}
	if query == nil {
		b.err = errors.New("view definition is required")
		return b
	}

	switch v := query.(type) {
	case shared.Raw:
		b.selectSQL = string(v)
		return b
	case *shared.Raw:
		b.selectSQL = string(*v)
		return b
	case interface {
		Build() (string, []interface{}, error)
	}:
		sql, args, err := v.Build()
		if err != nil {
			b.err = fmt.Errorf("failed to build view definition: %w", err)
			return b
		}
		if len(args) > 0 {
			b.err = errors.New("view definition builder must not use arguments")
			return b
		}
		b.selectSQL = sql
		return b
	default:
		b.err = errors.New("As() expects a builder with Build() or Raw")
		return b
	}
}

// OrReplace adds OR REPLACE to the CREATE VIEW statement.
func (b *CreateViewBuilder) OrReplace() *CreateViewBuilder {
	if b.err != nil {
		return b
	}
	b.orReplace = true
	return b
}

// Materialized marks this as a materialized view (for future use).
func (b *CreateViewBuilder) Materialized() *CreateViewBuilder {
	if b.err != nil {
		return b
	}
	b.materialized = true
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *CreateViewBuilder) WithDialect(d shared.Dialect) *CreateViewBuilder {
	if b.err != nil {
		return b
	}
	b.dialect = d
	return b
}

// Build builds the SQL CREATE VIEW query and returns the query string, arguments, and error if any.
func (b *CreateViewBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.viewName == "" {
		return "", nil, errors.New("view name is required")
	}
	if b.selectSQL == "" {
		return "", nil, errors.New("view definition is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = shared.GetDialect() // Use global dialect instead of defaulting to MySQL
	}

	var sb strings.Builder
	args := []interface{}{}

	// CREATE VIEW
	sb.WriteString("CREATE ")
	if b.orReplace {
		sb.WriteString("OR REPLACE ")
	}
	if b.materialized {
		sb.WriteString("MATERIALIZED ")
	}
	sb.WriteString("VIEW ")
	sb.WriteString(dialect.QuoteIdent(b.viewName))
	sb.WriteString(" AS ")
	sb.WriteString(b.selectSQL)

	return sb.String(), args, nil
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *CreateViewBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return shared.InterpolateSQL(sql, args)
}
