package ddl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sprylic/stk/shared"
)

// AlterTableBuilder builds SQL ALTER TABLE queries.
type AlterTableBuilder struct {
	tableName  string
	operations []AlterOperation
	err        error
	dialect    shared.Dialect
}

// AlterOperation represents a single ALTER TABLE operation.
type AlterOperation struct {
	Type           AlterOperationType
	Column         string
	NewName        string
	NewType        string
	Size           *int
	Precision      *int
	Scale          *int
	Nullable       *bool
	Default        interface{}
	ConstraintName string
	Columns        []string
	Reference      *ForeignKeyRef
	CheckExpr      string
	IndexName      string
	ConstraintType ConstraintType
}

// AlterOperationType represents the type of ALTER TABLE operation.
type AlterOperationType string

const (
	AddColumnType      AlterOperationType = "ADD COLUMN"
	DropColumnType     AlterOperationType = "DROP COLUMN"
	RenameColumnType   AlterOperationType = "RENAME COLUMN"
	RenameTableType    AlterOperationType = "RENAME TO"
	ModifyColumnType   AlterOperationType = "MODIFY COLUMN"
	AddConstraintType  AlterOperationType = "ADD CONSTRAINT"
	DropConstraintType AlterOperationType = "DROP CONSTRAINT"
	AddIndexType       AlterOperationType = "ADD INDEX"
	DropIndexType      AlterOperationType = "DROP INDEX"
)

// AlterTable creates a new AlterTableBuilder for the given table.
func AlterTable(tableName string) *AlterTableBuilder {
	if tableName == "" {
		return &AlterTableBuilder{err: errors.New("table name is required")}
	}
	return &AlterTableBuilder{
		tableName:  tableName,
		operations: make([]AlterOperation, 0),
	}
}

// AddColumn adds a column to the table.
func (b *AlterTableBuilder) AddColumn(cb *ColumnBuilder) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	col, err := cb.BuildDef()
	if err != nil {
		b.err = err
		return b
	}
	b.operations = append(b.operations, AlterOperation{
		Type:      AddColumnType,
		Column:    col.Name,
		NewType:   col.Type,
		Size:      col.Size,
		Precision: col.Precision,
		Scale:     col.Scale,
		Nullable:  col.Nullable,
		Default:   col.Default,
	})
	return b
}

// AddColumnWithType is a convenience method to add a column with just name and type.
func (b *AlterTableBuilder) AddColumnWithType(name, typ string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	b.operations = append(b.operations, AlterOperation{
		Type:    AddColumnType,
		Column:  name,
		NewType: strings.ToUpper(typ),
	})
	return b
}

// DropColumn drops a column from the table.
func (b *AlterTableBuilder) DropColumn(columnName string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if columnName == "" {
		b.err = errors.New("column name is required")
		return b
	}
	b.operations = append(b.operations, AlterOperation{
		Type:   DropColumnType,
		Column: columnName,
	})
	return b
}

// RenameColumn renames a column in the table.
func (b *AlterTableBuilder) RenameColumn(oldName, newName string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if oldName == "" {
		b.err = errors.New("old column name is required")
		return b
	}
	if newName == "" {
		b.err = errors.New("new column name is required")
		return b
	}
	b.operations = append(b.operations, AlterOperation{
		Type:    RenameColumnType,
		Column:  oldName,
		NewName: newName,
	})
	return b
}

// RenameTable renames the table.
func (b *AlterTableBuilder) RenameTable(newName string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if newName == "" {
		b.err = errors.New("new table name is required")
		return b
	}
	b.operations = append(b.operations, AlterOperation{
		Type:    RenameTableType,
		NewName: newName,
	})
	return b
}

// ModifyColumn modifies an existing column.
func (b *AlterTableBuilder) ModifyColumn(cb *ColumnBuilder) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	col, err := cb.BuildDef()
	if err != nil {
		b.err = err
		return b
	}
	b.operations = append(b.operations, AlterOperation{
		Type:      ModifyColumnType,
		Column:    col.Name,
		NewType:   col.Type,
		Size:      col.Size,
		Precision: col.Precision,
		Scale:     col.Scale,
		Nullable:  col.Nullable,
		Default:   col.Default,
	})
	return b
}

// AddConstraint adds a constraint to the table.
func (b *AlterTableBuilder) AddConstraint(constraint Constraint) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	b.operations = append(b.operations, AlterOperation{
		Type:           AddConstraintType,
		ConstraintName: constraint.Name,
		Columns:        constraint.Columns,
		Reference:      constraint.Reference,
		CheckExpr:      constraint.CheckExpr,
		ConstraintType: constraint.Type,
	})
	return b
}

// DropConstraint drops a constraint from the table.
func (b *AlterTableBuilder) DropConstraint(constraintName string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if constraintName == "" {
		b.err = errors.New("constraint name is required")
		return b
	}
	b.operations = append(b.operations, AlterOperation{
		Type:           DropConstraintType,
		ConstraintName: constraintName,
	})
	return b
}

// AddIndex adds an index to the table.
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
	b.operations = append(b.operations, AlterOperation{
		Type:      AddIndexType,
		IndexName: name,
		Columns:   columns,
	})
	return b
}

// DropIndex drops an index from the table.
func (b *AlterTableBuilder) DropIndex(indexName string) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
	if indexName == "" {
		b.err = errors.New("index name is required")
		return b
	}
	b.operations = append(b.operations, AlterOperation{
		Type:      DropIndexType,
		IndexName: indexName,
	})
	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *AlterTableBuilder) WithDialect(d shared.Dialect) *AlterTableBuilder {
	if b.err != nil {
		return b
	}
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
	if len(b.operations) == 0 {
		return "", nil, errors.New("at least one operation is required")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = shared.GetDialect() // Use global dialect instead of defaulting to MySQL
	}

	var sb strings.Builder
	args := []interface{}{}

	// ALTER TABLE
	sb.WriteString("ALTER TABLE ")
	sb.WriteString(dialect.QuoteIdent(b.tableName))

	// Operations
	operationSQLs := make([]string, 0, len(b.operations))
	for _, op := range b.operations {
		opSQL, err := b.buildOperationSQL(op, dialect)
		if err != nil {
			return "", nil, fmt.Errorf("operation %s: %w", op.Type, err)
		}
		operationSQLs = append(operationSQLs, opSQL)
	}

	sb.WriteString(" ")
	sb.WriteString(strings.Join(operationSQLs, ", "))

	return sb.String(), args, nil
}

// buildOperationSQL builds the SQL for a single ALTER TABLE operation.
func (b *AlterTableBuilder) buildOperationSQL(op AlterOperation, dialect shared.Dialect) (string, error) {
	switch op.Type {
	case AddColumnType:
		col := ColumnDef{
			Name:      op.Column,
			Type:      op.NewType,
			Size:      op.Size,
			Precision: op.Precision,
			Scale:     op.Scale,
			Nullable:  op.Nullable,
			Default:   op.Default,
		}
		colSQL, err := col.buildSQL(dialect)
		if err != nil {
			return "", err
		}
		return "ADD COLUMN " + colSQL, nil

	case DropColumnType:
		return "DROP COLUMN " + dialect.QuoteIdent(op.Column), nil

	case RenameColumnType:
		return "RENAME COLUMN " + dialect.QuoteIdent(op.Column) + " TO " + dialect.QuoteIdent(op.NewName), nil

	case RenameTableType:
		return "RENAME TO " + dialect.QuoteIdent(op.NewName), nil

	case ModifyColumnType:
		col := ColumnDef{
			Name:      op.Column,
			Type:      op.NewType,
			Size:      op.Size,
			Precision: op.Precision,
			Scale:     op.Scale,
			Nullable:  op.Nullable,
			Default:   op.Default,
		}
		colSQL, err := col.buildSQL(dialect)
		if err != nil {
			return "", err
		}
		return "MODIFY COLUMN " + colSQL, nil

	case AddConstraintType:
		constraint := Constraint{
			Type:      op.ConstraintType,
			Name:      op.ConstraintName,
			Columns:   op.Columns,
			Reference: op.Reference,
			CheckExpr: op.CheckExpr,
		}
		constraintSQL, err := constraint.buildSQL(dialect)
		if err != nil {
			return "", err
		}
		return "ADD " + constraintSQL, nil

	case DropConstraintType:
		return "DROP CONSTRAINT " + dialect.QuoteIdent(op.ConstraintName), nil

	case AddIndexType:
		indexSQL, err := buildIndexSQL(dialect, op.IndexName, op.Columns)
		if err != nil {
			return "", err
		}
		return "ADD INDEX " + dialect.QuoteIdent(op.IndexName) + " " + indexSQL, nil

	case DropIndexType:
		return "DROP INDEX " + dialect.QuoteIdent(op.IndexName), nil

	default:
		return "", fmt.Errorf("unsupported operation type: %s", op.Type)
	}
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *AlterTableBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return shared.InterpolateSQL(sql, args)
}
