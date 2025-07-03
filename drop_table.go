package cqb

import (
	"errors"
	"strings"
)

// DropTableBuilder builds SQL DROP TABLE queries.
type DropTableBuilder struct {
	tableName string
	ifExists  bool
	cascade   bool
	restrict  bool
	dialect   Dialect
	err       error
}

// DropTable creates a new DropTableBuilder for the given table.
func DropTable(tableName string) *DropTableBuilder {
	if tableName == "" {
		return &DropTableBuilder{err: errors.New("table name is required")}
	}
	return &DropTableBuilder{tableName: tableName}
}

// IfExists adds IF EXISTS to the DROP TABLE statement.
func (b *DropTableBuilder) IfExists() *DropTableBuilder {
	b.ifExists = true
	return b
}

// Cascade adds CASCADE to the DROP TABLE statement.
func (b *DropTableBuilder) Cascade() *DropTableBuilder {
	b.cascade = true
	b.restrict = false
	return b
}

// Restrict adds RESTRICT to the DROP TABLE statement.
func (b *DropTableBuilder) Restrict() *DropTableBuilder {
	b.restrict = true
	b.cascade = false
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *DropTableBuilder) WithDialect(d Dialect) *DropTableBuilder {
	b.dialect = d
	return b
}

// Build builds the SQL DROP TABLE query and returns the query string, arguments, and error if any.
func (b *DropTableBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.tableName == "" {
		return "", nil, errors.New("table name is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = getDialect()
	}

	var sb strings.Builder
	args := []interface{}{}

	sb.WriteString("DROP TABLE ")
	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}
	sb.WriteString(dialect.QuoteIdent(b.tableName))
	if b.cascade {
		sb.WriteString(" CASCADE")
	} else if b.restrict {
		sb.WriteString(" RESTRICT")
	}

	return sb.String(), args, nil
}
