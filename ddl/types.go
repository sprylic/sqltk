package ddl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sprylic/sqltk/shared"
)

// ColumnDef represents a column definition in a CREATE TABLE statement.
type ColumnDef struct {
	Name          string
	Type          string
	Size          *int
	Precision     *int
	Scale         *int
	Nullable      *bool
	Default       interface{}
	AutoIncrement bool
	IsPrimaryKey  bool
	IsUnique      bool
	Collation     string
	Charset       string
	Comment       string
	OnUpdate      string
}

// ConstraintType represents the type of constraint.
type ConstraintType string

const (
	PrimaryKeyType ConstraintType = "PRIMARY KEY"
	ForeignKeyType ConstraintType = "FOREIGN KEY"
	UniqueType     ConstraintType = "UNIQUE"
	CheckType      ConstraintType = "CHECK"
	IndexType      ConstraintType = "INDEX"
)

// Constraint represents a table constraint.
type Constraint struct {
	Type      ConstraintType
	Name      string
	Columns   []string
	Reference *ForeignKeyRef
	CheckExpr string
}

// ForeignKeyRef represents a foreign key reference.
type ForeignKeyRef struct {
	Table    string
	Columns  []string
	OnDelete string
	OnUpdate string
}

// TableOption represents a table option (ENGINE, CHARSET, etc.).
type TableOption struct {
	Name  string
	Value string
}

// buildColumnSQL builds the SQL for a column definition.
func (c *ColumnDef) buildSQL(dialect shared.Dialect) (string, error) {
	if c.Name == "" {
		return "", errors.New("column name is required")
	}
	if c.Type == "" {
		return "", errors.New("column type is required")
	}

	var parts []string
	parts = append(parts, dialect.QuoteIdent(c.Name))

	// Type with size/precision
	typeSQL := c.Type
	if c.Size != nil {
		typeSQL += fmt.Sprintf("(%d)", *c.Size)
	} else if c.Precision != nil {
		if c.Scale != nil {
			typeSQL += fmt.Sprintf("(%d,%d)", *c.Precision, *c.Scale)
		} else {
			typeSQL += fmt.Sprintf("(%d)", *c.Precision)
		}
	}
	parts = append(parts, typeSQL)

	// Charset
	if c.Charset != "" {
		parts = append(parts, "CHARACTER SET", c.Charset)
	}

	// Collation
	if c.Collation != "" {
		parts = append(parts, "COLLATE", c.Collation)
	}

	// Nullable
	if c.Nullable != nil {
		if *c.Nullable {
			parts = append(parts, "NULL")
		} else {
			parts = append(parts, "NOT NULL")
		}
	}

	// Default
	if c.Default != nil {
		parts = append(parts, "DEFAULT", formatDefaultValue(c.Default, dialect))
	}

	// Auto increment
	if c.AutoIncrement {
		if dialect == shared.Postgres() {
			// For Postgres, change the type to SERIAL based on the original type
			parts = parts[:1] // Keep only the quoted column name
			switch strings.ToUpper(c.Type) {
			case "BIGINT":
				parts = append(parts, "BIGSERIAL")
			case "INT", "INTEGER":
				parts = append(parts, "SERIAL")
			default:
				parts = append(parts, "SERIAL") // Default to SERIAL for other types
			}
			if c.Nullable != nil && !*c.Nullable {
				parts = append(parts, "NOT NULL")
			}
		} else {
			// For MySQL and Standard, use AUTO_INCREMENT
			parts = append(parts, "AUTO_INCREMENT")
		}
	}

	// Comment
	if c.Comment != "" {
		parts = append(parts, "COMMENT", dialect.QuoteString(c.Comment))
	}

	// ON UPDATE (MySQL only, but safe to emit for others if set)
	if c.OnUpdate != "" {
		parts = append(parts, "ON UPDATE", c.OnUpdate)
	}

	return strings.Join(parts, " "), nil
}

// formatDefaultValue formats a default value for SQL.
func formatDefaultValue(value interface{}, dialect shared.Dialect) string {
	switch v := value.(type) {
	case shared.Raw:
		// Raw SQL - include directly without quotes
		return string(v)
	case string:
		// String literals - quote them
		return dialect.QuoteString(v)
	case nil:
		return "NULL"
	default:
		// Numbers, booleans, etc. - format as-is
		return fmt.Sprintf("%v", v)
	}
}

// buildConstraintSQL builds the SQL for a constraint.
func (c *Constraint) buildSQL(dialect shared.Dialect) (string, error) {
	var parts []string

	switch c.Type {
	case PrimaryKeyType:
		parts = append(parts, "PRIMARY KEY")
		if len(c.Columns) > 0 {
			quotedCols := make([]string, len(c.Columns))
			for i, col := range c.Columns {
				quotedCols[i] = dialect.QuoteIdent(col)
			}
			parts = append(parts, "("+strings.Join(quotedCols, ", ")+")")
		}

	case UniqueType:
		if c.Name != "" {
			parts = append(parts, "CONSTRAINT", dialect.QuoteIdent(c.Name))
		}
		parts = append(parts, "UNIQUE")
		if len(c.Columns) > 0 {
			quotedCols := make([]string, len(c.Columns))
			for i, col := range c.Columns {
				quotedCols[i] = dialect.QuoteIdent(col)
			}
			parts = append(parts, "("+strings.Join(quotedCols, ", ")+")")
		}

	case CheckType:
		if c.Name != "" {
			parts = append(parts, "CONSTRAINT", dialect.QuoteIdent(c.Name))
		}
		parts = append(parts, "CHECK", "("+c.CheckExpr+")")

	case ForeignKeyType:
		if c.Name != "" {
			parts = append(parts, "CONSTRAINT", dialect.QuoteIdent(c.Name))
		}
		parts = append(parts, "FOREIGN KEY")
		if len(c.Columns) > 0 {
			quotedCols := make([]string, len(c.Columns))
			for i, col := range c.Columns {
				quotedCols[i] = dialect.QuoteIdent(col)
			}
			parts = append(parts, "("+strings.Join(quotedCols, ", ")+")")
		}
		if c.Reference != nil {
			parts = append(parts, "REFERENCES", dialect.QuoteIdent(c.Reference.Table))
			if len(c.Reference.Columns) > 0 {
				quotedRefCols := make([]string, len(c.Reference.Columns))
				for i, col := range c.Reference.Columns {
					quotedRefCols[i] = dialect.QuoteIdent(col)
				}
				parts = append(parts, "("+strings.Join(quotedRefCols, ", ")+")")
			}
			if c.Reference.OnDelete != "" {
				parts = append(parts, "ON DELETE", c.Reference.OnDelete)
			}
			if c.Reference.OnUpdate != "" {
				parts = append(parts, "ON UPDATE", c.Reference.OnUpdate)
			}
		}

	case IndexType:
		parts = append(parts, "INDEX")
		if c.Name != "" {
			parts = append(parts, dialect.QuoteIdent(c.Name))
		}
		if len(c.Columns) > 0 {
			quotedCols := make([]string, len(c.Columns))
			for i, col := range c.Columns {
				quotedCols[i] = dialect.QuoteIdent(col)
			}
			parts = append(parts, "("+strings.Join(quotedCols, ", ")+")")
		}

	default:
		return "", fmt.Errorf("unsupported constraint type: %s", c.Type)
	}

	return strings.Join(parts, " "), nil
}

// buildIndexSQL builds the SQL for an index definition.
func buildIndexSQL(dialect shared.Dialect, name string, columns []string) (string, error) {
	if name == "" {
		return "", errors.New("index name is required")
	}
	if len(columns) == 0 {
		return "", errors.New("at least one column is required for index")
	}

	quotedCols := make([]string, len(columns))
	for i, col := range columns {
		quotedCols[i] = dialect.QuoteIdent(col)
	}
	return "(" + strings.Join(quotedCols, ", ") + ")", nil
}

// ConstraintBuilder builds constraints with a fluent API.
type ConstraintBuilder struct {
	constraint Constraint
	err        error
}

// NewConstraint creates a new constraint builder.
func NewConstraint() *ConstraintBuilder {
	return &ConstraintBuilder{}
}

// Check creates a check constraint.
func (cb *ConstraintBuilder) Check(name, expr string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	cb.constraint = Constraint{
		Type:      CheckType,
		Name:      name,
		CheckExpr: expr,
	}
	return cb
}

// Unique creates a unique constraint.
func (cb *ConstraintBuilder) Unique(name string, columns ...string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	cb.constraint = Constraint{
		Type:    UniqueType,
		Name:    name,
		Columns: columns,
	}
	return cb
}

// PrimaryKey creates a primary key constraint.
func (cb *ConstraintBuilder) PrimaryKey(columns ...string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	cb.constraint = Constraint{
		Type:    PrimaryKeyType,
		Columns: columns,
	}
	return cb
}

// ForeignKey creates a foreign key constraint.
func (cb *ConstraintBuilder) ForeignKey(name string, columns ...string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	cb.constraint = Constraint{
		Type:    ForeignKeyType,
		Name:    name,
		Columns: columns,
	}
	return cb
}

// Index creates an index constraint.
func (cb *ConstraintBuilder) Index(name string, columns ...string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	cb.constraint = Constraint{
		Type:    IndexType,
		Name:    name,
		Columns: columns,
	}
	return cb
}

// Raw creates a raw constraint from a string (for advanced use cases).
func (cb *ConstraintBuilder) Raw(name, rawSQL string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	cb.constraint = Constraint{
		Type:      CheckType, // Default to CHECK for raw constraints
		Name:      name,
		CheckExpr: rawSQL,
	}
	return cb
}

// WithColumns sets the columns for the constraint.
func (cb *ConstraintBuilder) WithColumns(columns ...string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	cb.constraint.Columns = columns
	return cb
}

// WithCheckExpr sets the check expression for the constraint.
func (cb *ConstraintBuilder) WithCheckExpr(expr string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	cb.constraint.CheckExpr = expr
	return cb
}

// WithReference sets the foreign key reference for the constraint.
func (cb *ConstraintBuilder) WithReference(table string, columns ...string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	cb.constraint.Reference = &ForeignKeyRef{
		Table:   table,
		Columns: columns,
	}
	return cb
}

// WithOnDelete sets the ON DELETE action for the foreign key constraint.
func (cb *ConstraintBuilder) WithOnDelete(action string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	if cb.constraint.Reference == nil {
		cb.constraint.Reference = &ForeignKeyRef{}
	}
	cb.constraint.Reference.OnDelete = action
	return cb
}

// WithOnUpdate sets the ON UPDATE action for the foreign key constraint.
func (cb *ConstraintBuilder) WithOnUpdate(action string) *ConstraintBuilder {
	if cb.err != nil {
		return cb
	}
	if cb.constraint.Reference == nil {
		cb.constraint.Reference = &ForeignKeyRef{}
	}
	cb.constraint.Reference.OnUpdate = action
	return cb
}

// Build returns the built Constraint.
func (cb *ConstraintBuilder) Build() Constraint {
	return cb.constraint
}

// BuildPtr returns a pointer to the built Constraint.
func (cb *ConstraintBuilder) BuildPtr() *Constraint {
	return &cb.constraint
}
