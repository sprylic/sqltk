package cqb

import "github.com/sprylic/cqb/ddl"

// DDL package re-exports for backward compatibility

// Re-export DDL types and functions
type (
	ColumnDef     = ddl.ColumnDef
	Constraint    = ddl.Constraint
	ForeignKeyRef = ddl.ForeignKeyRef
	TableOption   = ddl.TableOption
)

// Re-export DDL builders
type (
	CreateTableBuilder = ddl.CreateTableBuilder
	DropTableBuilder   = ddl.DropTableBuilder
	AlterTableBuilder  = ddl.AlterTableBuilder
	CreateIndexBuilder = ddl.CreateIndexBuilder
	ColumnBuilder      = ddl.ColumnBuilder
	ForeignKeyBuilder  = ddl.ForeignKeyBuilder
)

// Re-export DDL functions
var (
	CreateTable = ddl.CreateTable
	DropTable   = ddl.DropTable
	AlterTable  = ddl.AlterTable
	CreateIndex = ddl.CreateIndex
	Column      = ddl.Column
	ForeignKey  = ddl.ForeignKey
)
