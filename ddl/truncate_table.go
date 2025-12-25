package ddl

import (
	"errors"
	"strings"

	"github.com/sprylic/sqltk/sqldebug"

	"github.com/sprylic/sqltk/sqldialect"
)

// TruncateTableBuilder builds SQL TRUNCATE TABLE queries.
type TruncateTableBuilder struct {
	tableNames []string
	cascade    bool
	restrict   bool
	restart    bool
	identity   bool
	err        error
	dialect    sqldialect.Dialect
}

// TruncateTable creates a new TruncateTableBuilder for the given table(s).
func TruncateTable(tableNames ...string) *TruncateTableBuilder {
	if len(tableNames) == 0 {
		return &TruncateTableBuilder{err: errors.New("at least one table name is required")}
	}
	for _, name := range tableNames {
		if name == "" {
			return &TruncateTableBuilder{err: errors.New("table name cannot be empty")}
		}
	}
	return &TruncateTableBuilder{
		tableNames: tableNames,
	}
}

// Cascade adds CASCADE to the TRUNCATE TABLE statement.
func (b *TruncateTableBuilder) Cascade() *TruncateTableBuilder {
	if b.err != nil {
		return b
	}
	b.cascade = true
	b.restrict = false // CASCADE and RESTRICT are mutually exclusive
	return b
}

// Restrict adds RESTRICT to the TRUNCATE TABLE statement.
func (b *TruncateTableBuilder) Restrict() *TruncateTableBuilder {
	if b.err != nil {
		return b
	}
	b.restrict = true
	b.cascade = false // CASCADE and RESTRICT are mutually exclusive
	return b
}

// Restart adds RESTART IDENTITY to the TRUNCATE TABLE statement (PostgreSQL).
func (b *TruncateTableBuilder) Restart() *TruncateTableBuilder {
	if b.err != nil {
		return b
	}
	b.restart = true
	b.identity = false // RESTART and CONTINUE are mutually exclusive
	return b
}

// Continue adds CONTINUE IDENTITY to the TRUNCATE TABLE statement (PostgreSQL).
func (b *TruncateTableBuilder) Continue() *TruncateTableBuilder {
	if b.err != nil {
		return b
	}
	b.identity = true
	b.restart = false // RESTART and CONTINUE are mutually exclusive
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *TruncateTableBuilder) WithDialect(d sqldialect.Dialect) *TruncateTableBuilder {
	if b.err != nil {
		return b
	}
	b.dialect = d
	return b
}

// Build builds the SQL TRUNCATE TABLE query and returns the query string, arguments, and error if any.
func (b *TruncateTableBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if len(b.tableNames) == 0 {
		return "", nil, errors.New("at least one table name is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = sqldialect.GetDialect() // Use global dialect instead of defaulting to MySQL
	}

	var sb strings.Builder
	args := []interface{}{}

	// TRUNCATE TABLE
	sb.WriteString("TRUNCATE TABLE ")

	// Table names
	quotedNames := make([]string, len(b.tableNames))
	for i, name := range b.tableNames {
		quotedNames[i] = dialect.QuoteIdent(name)
	}
	sb.WriteString(strings.Join(quotedNames, ", "))

	// PostgreSQL-specific options
	if dialect == sqldialect.Postgres() {
		var options []string

		if b.restart {
			options = append(options, "RESTART IDENTITY")
		} else if b.identity {
			options = append(options, "CONTINUE IDENTITY")
		}

		if b.cascade {
			options = append(options, "CASCADE")
		} else if b.restrict {
			options = append(options, "RESTRICT")
		}

		if len(options) > 0 {
			sb.WriteString(" " + strings.Join(options, " "))
		}
	} else {
		// MySQL and Standard dialects
		if b.cascade {
			sb.WriteString(" CASCADE")
		} else if b.restrict {
			sb.WriteString(" RESTRICT")
		}
	}

	return sb.String(), args, nil
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *TruncateTableBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return sqldebug.InterpolateSQL(sql, args).GetUnsafeString()
}
