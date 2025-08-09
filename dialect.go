package sqltk

import "github.com/sprylic/sqltk/sqldialect"

func SetDialect(dialect sqldialect.Dialect) {
	sqldialect.SetDialect(dialect)
}
