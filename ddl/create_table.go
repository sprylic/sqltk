package ddl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sprylic/sqltk/shared"
	"github.com/sprylic/sqltk/sqlfunc"
)

// CreateTableBuilder builds SQL CREATE TABLE queries.
type CreateTableBuilder struct {
	tableName   string
	columns     []ColumnDef
	constraints []Constraint
	options     []TableOption // ENGINE, CHARSET, etc. in order
	ifNotExists bool
	temporary   bool
	err         error
	dialect     shared.Dialect
}

// CreateTable creates a new CreateTableBuilder for the given table.
func CreateTable(tableName string) *CreateTableBuilder {
	if tableName == "" {
		return &CreateTableBuilder{err: errors.New("table name is required")}
	}
	return &CreateTableBuilder{
		tableName: tableName,
		options:   make([]TableOption, 0),
	}
}

// ColumnBuilder builds a column definition.
type ColumnBuilder struct {
	def ColumnDef
	err error
}

// Column creates a new ColumnBuilder for the given column name.
func Column(name string) *ColumnBuilder {
	if name == "" {
		return &ColumnBuilder{err: errors.New("column name is required")}
	}
	return &ColumnBuilder{
		def: ColumnDef{Name: name},
	}
}

// BuildDef returns the built ColumnDef and any error.
func (cb *ColumnBuilder) BuildDef() (ColumnDef, error) {
	if cb.err != nil {
		return ColumnDef{}, cb.err
	}
	if cb.def.Type == "" {
		return ColumnDef{}, errors.New("column type is required")
	}
	return cb.def, nil
}

// AddColumn adds a column from a ColumnBuilder.
func (b *CreateTableBuilder) AddColumn(cb *ColumnBuilder) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	col, err := cb.BuildDef()
	if err != nil {
		b.err = err
		return b
	}
	b.columns = append(b.columns, col)
	return b
}

// AddColumns adds columns from a ColumnBuilder.
func (b *CreateTableBuilder) AddColumns(cbs ...*ColumnBuilder) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	for _, cb := range cbs {
		col, err := cb.BuildDef()
		if err != nil {
			b.err = err
		}
		b.columns = append(b.columns, col)
	}
	return b
}

// AddColumnWithType is a convenience method to add a column with just name and type.
func (b *CreateTableBuilder) AddColumnWithType(name, typ string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	col := ColumnDef{
		Name: name,
		Type: strings.ToUpper(typ),
	}
	b.columns = append(b.columns, col)
	return b
}

// Type sets the column type.
func (cb *ColumnBuilder) Type(typ string) *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	if typ == "" {
		cb.err = errors.New("column type is required")
		return cb
	}
	cb.def.Type = strings.ToUpper(typ)
	return cb
}

// Size sets the column size (for VARCHAR, INT, etc.).
func (cb *ColumnBuilder) Size(size int) *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	if size <= 0 {
		cb.err = errors.New("column size must be positive")
		return cb
	}
	cb.def.Size = &size
	return cb
}

// Precision sets the column precision and scale (for DECIMAL, NUMERIC, etc.).
func (cb *ColumnBuilder) Precision(precision, scale int) *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	if precision <= 0 {
		cb.err = errors.New("precision must be positive")
		return cb
	}
	if scale < 0 || scale > precision {
		cb.err = errors.New("scale must be between 0 and precision")
		return cb
	}
	cb.def.Precision = &precision
	cb.def.Scale = &scale
	return cb
}

// Nullable makes the column nullable.
func (cb *ColumnBuilder) Nullable() *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	nullable := true
	cb.def.Nullable = &nullable
	return cb
}

// NotNull makes the column not null.
func (cb *ColumnBuilder) NotNull() *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	nullable := false
	cb.def.Nullable = &nullable
	return cb
}

// Default sets the column default value.
func (cb *ColumnBuilder) Default(value interface{}) *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	cb.def.Default = value
	return cb
}

// AutoIncrement makes the column auto-incrementing.
func (cb *ColumnBuilder) AutoIncrement() *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	cb.def.AutoIncrement = true
	return cb
}

// PrimaryKey marks this column as a primary key.
func (cb *ColumnBuilder) PrimaryKey() *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	cb.def.IsPrimaryKey = true
	return cb
}

// Unique marks this column as unique.
func (cb *ColumnBuilder) Unique() *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	cb.def.IsUnique = true
	return cb
}

// Collation sets the column collation.
func (cb *ColumnBuilder) Collation(collation string) *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	cb.def.Collation = collation
	return cb
}

// Charset sets the column character set.
func (cb *ColumnBuilder) Charset(charset string) *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	cb.def.Charset = charset
	return cb
}

// Comment sets the column comment.
func (cb *ColumnBuilder) Comment(comment string) *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	cb.def.Comment = comment
	return cb
}

// OnUpdate sets the ON UPDATE action for the column default value (e.g., ON UPDATE CURRENT_TIMESTAMP).
func (cb *ColumnBuilder) OnUpdate(action interface{}) *ColumnBuilder {
	if cb.err != nil {
		return cb
	}
	var actionStr string
	switch v := action.(type) {
	case string:
		actionStr = v
	case sqlfunc.SqlFunc:
		actionStr = string(v)
	default:
		cb.err = fmt.Errorf("OnUpdate action must be string or sqlfunc.SqlFunc, got %T", action)
		return cb
	}
	cb.def.OnUpdate = actionStr
	return cb
}

// Table-level options
func (b *CreateTableBuilder) IfNotExists() *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	b.ifNotExists = true
	return b
}

func (b *CreateTableBuilder) Temporary() *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	b.temporary = true
	return b
}

func (b *CreateTableBuilder) Engine(engine string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	b.options = append(b.options, TableOption{Name: "ENGINE", Value: engine})
	return b
}

func (b *CreateTableBuilder) Charset(charset string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	b.options = append(b.options, TableOption{Name: "CHARACTER SET", Value: charset})
	return b
}

func (b *CreateTableBuilder) Collation(collation string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	b.options = append(b.options, TableOption{Name: "COLLATE", Value: collation})
	return b
}

func (b *CreateTableBuilder) Comment(comment string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	b.options = append(b.options, TableOption{Name: "COMMENT", Value: comment})
	return b
}

// Constraint methods
func (b *CreateTableBuilder) PrimaryKey(columns ...string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if len(columns) == 0 {
		b.err = errors.New("primary key must specify at least one column")
		return b
	}
	b.constraints = append(b.constraints, Constraint{
		Type:    PrimaryKeyType,
		Columns: columns,
	})
	return b
}

func (b *CreateTableBuilder) Unique(name string, columns ...string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if name == "" {
		b.err = errors.New("unique constraint name is required")
		return b
	}
	if len(columns) == 0 {
		b.err = errors.New("unique constraint must specify at least one column")
		return b
	}
	b.constraints = append(b.constraints, Constraint{
		Type:    UniqueType,
		Name:    name,
		Columns: columns,
	})
	return b
}

func (b *CreateTableBuilder) Check(name, expr string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if name == "" {
		b.err = errors.New("check constraint name is required")
		return b
	}
	if expr == "" {
		b.err = errors.New("check constraint expression is required")
		return b
	}
	b.constraints = append(b.constraints, Constraint{
		Type:      CheckType,
		Name:      name,
		CheckExpr: expr,
	})
	return b
}

func (b *CreateTableBuilder) Index(name string, columns ...string) *CreateTableBuilder {
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
	b.constraints = append(b.constraints, Constraint{
		Type:    IndexType,
		Name:    name,
		Columns: columns,
	})
	return b
}

// ForeignKeyBuilder builds foreign key constraints.
type ForeignKeyBuilder struct {
	parent     *CreateTableBuilder
	constraint Constraint
	err        error
}

// ForeignKey creates a standalone foreign key builder.
func ForeignKey(name string, columns ...string) *ForeignKeyBuilder {
	if name == "" {
		return &ForeignKeyBuilder{err: errors.New("foreign key name is required")}
	}
	if len(columns) == 0 {
		return &ForeignKeyBuilder{err: errors.New("foreign key must specify at least one column")}
	}
	return &ForeignKeyBuilder{
		constraint: Constraint{
			Type:    ForeignKeyType,
			Name:    name,
			Columns: columns,
		},
	}
}

// References sets the referenced table and columns.
func (fkb *ForeignKeyBuilder) References(table string, column string, columns ...string) *ForeignKeyBuilder {
	if fkb.err != nil {
		return fkb
	}
	if table == "" {
		fkb.err = errors.New("referenced table is required")
		return fkb
	}

	columns = append([]string{column}, columns...)

	fkb.constraint.Reference = &ForeignKeyRef{
		Table:   table,
		Columns: columns,
	}
	return fkb
}

// OnDelete sets the ON DELETE action.
func (fkb *ForeignKeyBuilder) OnDelete(action string) *ForeignKeyBuilder {
	if fkb.err != nil {
		return fkb
	}
	if fkb.constraint.Reference == nil {
		fkb.err = errors.New("must call References before OnDelete")
		return fkb
	}
	fkb.constraint.Reference.OnDelete = action
	return fkb
}

// OnUpdate sets the ON UPDATE action.
func (fkb *ForeignKeyBuilder) OnUpdate(action string) *ForeignKeyBuilder {
	if fkb.err != nil {
		return fkb
	}
	if fkb.constraint.Reference == nil {
		fkb.err = errors.New("must call References before OnUpdate")
		return fkb
	}
	fkb.constraint.Reference.OnUpdate = action
	return fkb
}

// Build finalizes the foreign key constraint and adds it to the table.
func (fkb *ForeignKeyBuilder) Build() *CreateTableBuilder {
	if fkb.err != nil {
		fkb.parent.err = fkb.err
		return fkb.parent
	}
	if fkb.constraint.Reference == nil {
		fkb.parent.err = errors.New("foreign key must specify referenced table")
		return fkb.parent
	}
	fkb.parent.constraints = append(fkb.parent.constraints, fkb.constraint)
	return fkb.parent
}

// AddForeignKey adds a foreign key constraint to the table.
func (b *CreateTableBuilder) AddForeignKey(fkb *ForeignKeyBuilder) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if fkb.err != nil {
		b.err = fkb.err
		return b
	}
	if fkb.constraint.Reference == nil {
		b.err = errors.New("foreign key must specify referenced table")
		return b
	}
	b.constraints = append(b.constraints, fkb.constraint)
	return b
}

// AddForeignKeys adds foriegn key constraints to the table
func (b *CreateTableBuilder) AddForeignKeys(fkbs ...*ForeignKeyBuilder) *CreateTableBuilder {
	if b.err != nil {
		return b
	}

	for _, fkb := range fkbs {
		b.AddForeignKey(fkb)
	}

	return b
}

// WithDialect sets the dialect for this builder instance.
func (b *CreateTableBuilder) WithDialect(d shared.Dialect) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	b.dialect = d
	return b
}

// Build builds the SQL CREATE TABLE query and returns the query string, arguments, and error if any.
func (b *CreateTableBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	if b.tableName == "" {
		return "", nil, errors.New("table name is required")
	}
	if len(b.columns) == 0 {
		return "", nil, errors.New("at least one column must be defined")
	}

	dialect := b.dialect
	if dialect == nil {
		dialect = shared.GetDialect() // Use global dialect instead of defaulting to MySQL
	}

	var sb strings.Builder
	args := []interface{}{}

	// CREATE TABLE
	sb.WriteString("CREATE ")
	if b.temporary {
		sb.WriteString("TEMPORARY ")
	}
	sb.WriteString("TABLE ")
	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	sb.WriteString(dialect.QuoteIdent(b.tableName))
	sb.WriteString(" (")

	// Columns
	columnSQLs := make([]string, 0, len(b.columns))
	for _, col := range b.columns {
		colSQL, err := col.buildSQL(dialect)
		if err != nil {
			return "", nil, fmt.Errorf("column %s: %w", col.Name, err)
		}
		columnSQLs = append(columnSQLs, colSQL)
	}

	// Handle columns marked as primary keys
	var primaryKeyColumns []string
	for _, col := range b.columns {
		if col.IsPrimaryKey {
			primaryKeyColumns = append(primaryKeyColumns, col.Name)
		}
	}

	// Handle columns marked as unique
	var uniqueColumns []string
	for _, col := range b.columns {
		if col.IsUnique {
			uniqueColumns = append(uniqueColumns, col.Name)
		}
	}

	// Constraints
	for _, constraint := range b.constraints {
		constraintSQL, err := constraint.buildSQL(dialect)
		if err != nil {
			return "", nil, fmt.Errorf("constraint: %w", err)
		}
		columnSQLs = append(columnSQLs, constraintSQL)
	}

	// Add primary key constraint for columns marked as primary keys
	if len(primaryKeyColumns) > 0 {
		primaryKeyConstraint := Constraint{
			Type:    PrimaryKeyType,
			Columns: primaryKeyColumns,
		}
		constraintSQL, err := primaryKeyConstraint.buildSQL(dialect)
		if err != nil {
			return "", nil, fmt.Errorf("primary key constraint: %w", err)
		}
		columnSQLs = append(columnSQLs, constraintSQL)
	}

	// Add unique constraints for columns marked as unique
	for _, colName := range uniqueColumns {
		uniqueConstraint := Constraint{
			Type:    UniqueType,
			Columns: []string{colName},
		}
		constraintSQL, err := uniqueConstraint.buildSQL(dialect)
		if err != nil {
			return "", nil, fmt.Errorf("unique constraint: %w", err)
		}
		columnSQLs = append(columnSQLs, constraintSQL)
	}

	sb.WriteString(strings.Join(columnSQLs, ", "))
	sb.WriteString(")")

	// Table options in order
	if len(b.options) > 0 {
		optionSQLs := make([]string, 0, len(b.options))
		for _, opt := range b.options {
			if opt.Value == "" {
				optionSQLs = append(optionSQLs, opt.Name)
			} else {
				if opt.Name == "COMMENT" {
					optionSQLs = append(optionSQLs, fmt.Sprintf("%s %s", opt.Name, dialect.QuoteString(opt.Value)))
				} else {
					optionSQLs = append(optionSQLs, fmt.Sprintf("%s %s", opt.Name, opt.Value))
				}
			}
		}
		sb.WriteString(" ")
		sb.WriteString(strings.Join(optionSQLs, " "))
	}

	// For PostgreSQL, generate triggers for OnUpdate columns
	if dialect == shared.Postgres() {
		triggerSQL := b.buildPostgresTriggers(dialect)
		if triggerSQL != "" {
			sb.WriteString(";\n")
			sb.WriteString(triggerSQL)
		}
	}

	return sb.String(), args, nil
}

// buildPostgresTriggers generates PostgreSQL triggers for columns with OnUpdate
func (b *CreateTableBuilder) buildPostgresTriggers(dialect shared.Dialect) string {
	var triggers []string

	for _, col := range b.columns {
		if col.OnUpdate != "" {
			// Generate trigger function name
			triggerFuncName := fmt.Sprintf("%s_%s_update_trigger", b.tableName, col.Name)
			triggerName := fmt.Sprintf("tr_%s_%s_update", b.tableName, col.Name)

			// Create trigger function
			triggerFunc := fmt.Sprintf(`
CREATE OR REPLACE FUNCTION %s()
RETURNS TRIGGER AS $$
BEGIN
    NEW.%s = %s;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;`,
				dialect.QuoteIdent(triggerFuncName),
				dialect.QuoteIdent(col.Name),
				col.OnUpdate)

			// Create trigger
			trigger := fmt.Sprintf(`
CREATE OR REPLACE TRIGGER %s
    BEFORE UPDATE ON %s
    FOR EACH ROW
    EXECUTE FUNCTION %s();`,
				dialect.QuoteIdent(triggerName),
				dialect.QuoteIdent(b.tableName),
				dialect.QuoteIdent(triggerFuncName))

			if b.ifNotExists {
				// Use a DO block to check for trigger existence
				doBlock := fmt.Sprintf(`
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = '%s'
    ) THEN
        %s
    END IF;
END$$;`,
					triggerName,
					trigger)
				triggers = append(triggers, triggerFunc, doBlock)
			} else {
				triggers = append(triggers, triggerFunc, trigger)
			}
		}
	}

	return strings.Join(triggers, "\n")
}

// DebugSQL returns the SQL with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func (b *CreateTableBuilder) DebugSQL() string {
	sql, args, _ := b.Build()
	return shared.InterpolateSQL(sql, args)
}

// Column constraint methods that can be chained after convenience methods

// NotNull makes the last added column not null.
func (b *CreateTableBuilder) NotNull() *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if len(b.columns) == 0 {
		b.err = errors.New("no column to modify")
		return b
	}
	nullable := false
	b.columns[len(b.columns)-1].Nullable = &nullable
	return b
}

// Nullable makes the last added column nullable.
func (b *CreateTableBuilder) Nullable() *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if len(b.columns) == 0 {
		b.err = errors.New("no column to modify")
		return b
	}
	nullable := true
	b.columns[len(b.columns)-1].Nullable = &nullable
	return b
}

// Default sets the default value for the last added column.
func (b *CreateTableBuilder) Default(value interface{}) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if len(b.columns) == 0 {
		b.err = errors.New("no column to modify")
		return b
	}
	b.columns[len(b.columns)-1].Default = value
	return b
}

// AutoIncrement makes the last added column auto-incrementing.
func (b *CreateTableBuilder) AutoIncrement() *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if len(b.columns) == 0 {
		b.err = errors.New("no column to modify")
		return b
	}
	b.columns[len(b.columns)-1].AutoIncrement = true
	return b
}

// ColumnCollation sets the collation for the last added column.
func (b *CreateTableBuilder) ColumnCollation(collation string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if len(b.columns) == 0 {
		b.err = errors.New("no column to modify")
		return b
	}
	b.columns[len(b.columns)-1].Collation = collation
	return b
}

// ColumnCharset sets the character set for the last added column.
func (b *CreateTableBuilder) ColumnCharset(charset string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if len(b.columns) == 0 {
		b.err = errors.New("no column to modify")
		return b
	}
	b.columns[len(b.columns)-1].Charset = charset
	return b
}

// ColumnComment sets the comment for the last added column.
func (b *CreateTableBuilder) ColumnComment(comment string) *CreateTableBuilder {
	if b.err != nil {
		return b
	}
	if len(b.columns) == 0 {
		b.err = errors.New("no column to modify")
		return b
	}
	b.columns[len(b.columns)-1].Comment = comment
	return b
}
