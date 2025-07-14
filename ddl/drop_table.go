package ddl

import (
	"errors"
	"strings"

	"github.com/sprylic/sqltk/shared"
)

// DropTableBuilder builds SQL DROP TABLE queries.
type DropTableBuilder struct {
	tableNames []string
	ifExists   bool
	cascade    bool
	restrict   bool
	err        error
	dialect    shared.Dialect
}

// DropTable creates a new DropTableBuilder for the given table(s).
func DropTable(tableNames ...string) *DropTableBuilder {
	if len(tableNames) == 0 {
		return &DropTableBuilder{err: errors.New("at least one table name is required")}
	}
	for _, name := range tableNames {
		if name == "" {
			return &DropTableBuilder{err: errors.New("table name cannot be empty")}
		}
	}
	return &DropTableBuilder{
		tableNames: tableNames,
	}
}

// IfExists adds IF EXISTS to the DROP TABLE statement.
func (b *DropTableBuilder) IfExists() *DropTableBuilder {
	if b.err != nil {
		return b
	}
	b.ifExists = true
	return b
}

// Cascade adds CASCADE to the DROP TABLE statement.
func (b *DropTableBuilder) Cascade() *DropTableBuilder {
	if b.err != nil {
		return b
	}
	b.cascade = true
	b.restrict = false // CASCADE and RESTRICT are mutually exclusive
	return b
}

// Restrict adds RESTRICT to the DROP TABLE statement.
func (b *DropTableBuilder) Restrict() *DropTableBuilder {
	if b.err != nil {
		return b
	}
	b.restrict = true
	b.cascade = false // CASCADE and RESTRICT are mutually exclusive
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *DropTableBuilder) WithDialect(d shared.Dialect) *DropTableBuilder {
	if b.err != nil {
		return b
	}
	b.dialect = d
	return b
}

// Build builds the SQL DROP TABLE query and returns the query string, arguments, and error if any.
func (b *DropTableBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if len(b.tableNames) == 0 {
		return "", nil, errors.New("at least one table name is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = shared.GetDialect() // Use global dialect instead of defaulting to MySQL
	}

	var sb strings.Builder
	args := []interface{}{}

	// DROP TABLE
	sb.WriteString("DROP TABLE ")
	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}

	// Table names
	quotedNames := make([]string, len(b.tableNames))
	for i, name := range b.tableNames {
		quotedNames[i] = dialect.QuoteIdent(name)
	}
	sb.WriteString(strings.Join(quotedNames, ", "))

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
func (b *DropTableBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return shared.InterpolateSQL(sql, args).GetUnsafeString()
}
