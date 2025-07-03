package cqb

import (
	"errors"
	"strings"
)

// AlterTableOpType represents the type of ALTER TABLE operation.
type AlterTableOpType int

const (
	AlterAddColumn AlterTableOpType = iota
	AlterDropColumn
	AlterRenameColumn
	AlterModifyColumn
	AlterAddConstraint
	AlterDropConstraint
	AlterRenameTable
	AlterAddIndex
	AlterDropIndex
)

// AlterTableOp represents a single ALTER TABLE operation.
type AlterTableOp struct {
	Type       AlterTableOpType
	Column     *ColumnDef  // for add/modify
	ColumnName string      // for drop
	OldName    string      // for rename column
	NewName    string      // for rename column/table
	Constraint *Constraint // for add/drop constraint
	IndexName  string      // for add/drop index
	IndexCols  []string    // for add index
}

// AlterTableBuilder builds SQL ALTER TABLE queries.
type AlterTableBuilder struct {
	tableName string
	ops       []AlterTableOp
	dialect   Dialect
	err       error
}

// AlterTable creates a new AlterTableBuilder for the given table.
func AlterTable(tableName string) *AlterTableBuilder {
	if tableName == "" {
		return &AlterTableBuilder{err: errors.New("table name is required")}
	}
	return &AlterTableBuilder{tableName: tableName}
}

// AddColumn adds an ADD COLUMN operation.
func (b *AlterTableBuilder) AddColumn(cb *ColumnBuilder) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	col, err := cb.BuildDef()
	if err != nil {
		b.err = err
		return b
	}
	b.ops = append(b.ops, AlterTableOp{Type: AlterAddColumn, Column: &col})
	return b
}

// DropColumn adds a DROP COLUMN operation.
func (b *AlterTableBuilder) DropColumn(name string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if name == "" {
		b.err = errors.New("column name is required")
		return b
	}
	b.ops = append(b.ops, AlterTableOp{Type: AlterDropColumn, ColumnName: name})
	return b
}

// RenameColumn adds a RENAME COLUMN operation.
func (b *AlterTableBuilder) RenameColumn(oldName, newName string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if oldName == "" || newName == "" {
		b.err = errors.New("old and new column names are required")
		return b
	}
	b.ops = append(b.ops, AlterTableOp{Type: AlterRenameColumn, OldName: oldName, NewName: newName})
	return b
}

// RenameTable adds a RENAME TO operation.
func (b *AlterTableBuilder) RenameTable(newName string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if newName == "" {
		b.err = errors.New("new table name is required")
		return b
	}
	b.ops = append(b.ops, AlterTableOp{Type: AlterRenameTable, NewName: newName})
	return b
}

// ModifyColumn adds a MODIFY COLUMN operation.
func (b *AlterTableBuilder) ModifyColumn(cb *ColumnBuilder) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	col, err := cb.BuildDef()
	if err != nil {
		b.err = err
		return b
	}
	b.ops = append(b.ops, AlterTableOp{Type: AlterModifyColumn, Column: &col})
	return b
}

// AddConstraint adds an ADD CONSTRAINT operation.
func (b *AlterTableBuilder) AddConstraint(constraint *Constraint) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if constraint == nil {
		b.err = errors.New("constraint is required")
		return b
	}
	b.ops = append(b.ops, AlterTableOp{Type: AlterAddConstraint, Constraint: constraint})
	return b
}

// DropConstraint adds a DROP CONSTRAINT operation.
func (b *AlterTableBuilder) DropConstraint(name string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if name == "" {
		b.err = errors.New("constraint name is required")
		return b
	}
	b.ops = append(b.ops, AlterTableOp{Type: AlterDropConstraint, Constraint: &Constraint{Name: name}})
	return b
}

// AddIndex adds an ADD INDEX operation.
func (b *AlterTableBuilder) AddIndex(name string, columns ...string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if name == "" {
		b.err = errors.New("index name is required")
		return b
	}
	if len(columns) == 0 {
		b.err = errors.New("at least one column is required for index")
		return b
	}
	b.ops = append(b.ops, AlterTableOp{Type: AlterAddIndex, IndexName: name, IndexCols: columns})
	return b
}

// DropIndex adds a DROP INDEX operation.
func (b *AlterTableBuilder) DropIndex(name string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if name == "" {
		b.err = errors.New("index name is required")
		return b
	}
	b.ops = append(b.ops, AlterTableOp{Type: AlterDropIndex, IndexName: name})
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *AlterTableBuilder) WithDialect(d Dialect) *AlterTableBuilder {
	b.dialect = d
	return b
}

// Build builds the SQL ALTER TABLE query and returns the query string, arguments, and error if any.
func (b *AlterTableBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.tableName == "" {
		return "", nil, errors.New("table name is required")
	}
	if len(b.ops) == 0 {
		return "", nil, errors.New("no alter operations specified")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = getDialect()
	}

	var sb strings.Builder
	args := []interface{}{}

	sb.WriteString("ALTER TABLE ")
	sb.WriteString(dialect.QuoteIdent(b.tableName))

	opSQLs := make([]string, 0, len(b.ops))
	for _, op := range b.ops {
		switch op.Type {
		case AlterAddColumn:
			colSQL, err := op.Column.buildSQL(dialect)
			if err != nil {
				return "", nil, err
			}
			opSQLs = append(opSQLs, "ADD COLUMN "+colSQL)
		case AlterDropColumn:
			opSQLs = append(opSQLs, "DROP COLUMN "+dialect.QuoteIdent(op.ColumnName))
		case AlterRenameColumn:
			opSQLs = append(opSQLs, "RENAME COLUMN "+dialect.QuoteIdent(op.OldName)+" TO "+dialect.QuoteIdent(op.NewName))
		case AlterRenameTable:
			opSQLs = append(opSQLs, "RENAME TO "+dialect.QuoteIdent(op.NewName))
		case AlterModifyColumn:
			colSQL, err := op.Column.buildSQL(dialect)
			if err != nil {
				return "", nil, err
			}
			opSQLs = append(opSQLs, "MODIFY COLUMN "+colSQL)
		case AlterAddConstraint:
			constraintSQL, err := op.Constraint.buildSQL(dialect)
			if err != nil {
				return "", nil, err
			}
			opSQLs = append(opSQLs, "ADD "+constraintSQL)
		case AlterDropConstraint:
			opSQLs = append(opSQLs, "DROP CONSTRAINT "+dialect.QuoteIdent(op.Constraint.Name))
		case AlterAddIndex:
			indexSQL, err := buildIndexSQL(dialect, op.IndexName, op.IndexCols)
			if err != nil {
				return "", nil, err
			}
			opSQLs = append(opSQLs, "ADD INDEX "+dialect.QuoteIdent(op.IndexName)+" "+indexSQL)
		case AlterDropIndex:
			opSQLs = append(opSQLs, "DROP INDEX "+dialect.QuoteIdent(op.IndexName))
		}
	}

	sb.WriteString(" ")
	sb.WriteString(strings.Join(opSQLs, ", "))

	return sb.String(), args, nil
}

// buildIndexSQL builds the SQL for an index definition.
func buildIndexSQL(dialect Dialect, name string, columns []string) (string, error) {
	if len(columns) == 0 {
		return "", errors.New("at least one column is required for index")
	}

	quotedCols := make([]string, len(columns))
	for i, col := range columns {
		quotedCols[i] = dialect.QuoteIdent(col)
	}

	return "(" + strings.Join(quotedCols, ", ") + ")", nil
}
