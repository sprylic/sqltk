package stk

import (
	"errors"
	"fmt"
	"strings"
)

// whereClause holds shared WHERE clause logic for builders.
type whereClause struct {
	whereParam []string
	whereRaw   []string
	whereArgs  []interface{}
	err        error
}

func (w *whereClause) Where(cond interface{}, args ...interface{}) {
	if w.err != nil {
		return
	}
	switch c := cond.(type) {
	case Raw:
		w.whereRaw = append(w.whereRaw, string(c))
	case string:
		w.whereParam = append(w.whereParam, c)
		w.whereArgs = append(w.whereArgs, args...)
	case *ConditionBuilder:
		sql, condArgs, err := c.Build()
		if err != nil {
			w.err = fmt.Errorf("Where: condition builder error: %w", err)
			return
		}
		if sql != "" {
			w.whereParam = append(w.whereParam, sql)
			w.whereArgs = append(w.whereArgs, condArgs...)
		}
	default:
		w.err = fmt.Errorf("Where: cond must be string, sq.Raw, or *ConditionBuilder (got type %T)", cond)
	}
}

func (w *whereClause) WhereEqual(column string, value interface{}) {
	if value == nil {
		w.Where(column + " IS NULL")
		return
	}
	w.Where(column+" = ?", value)
}

func (w *whereClause) WhereNotEqual(column string, value interface{}) {
	if value == nil {
		w.Where(column + " IS NOT NULL")
		return
	}
	w.Where(column+" != ?", value)
}

func (w *whereClause) buildWhereSQL(dialect Dialect, placeholderIdx *int) (string, []interface{}) {
	var wheres []string
	if len(w.whereParam) > 0 {
		wheres = append(wheres, w.whereParam...)
	}
	if len(w.whereRaw) > 0 {
		wheres = append(wheres, w.whereRaw...)
	}
	if len(wheres) == 0 {
		return "", nil
	}
	whereSQL := strings.Join(wheres, " AND ")
	for strings.Contains(whereSQL, "?") && dialect.Placeholder(0) != "?" {
		whereSQL = strings.Replace(whereSQL, "?", dialect.Placeholder(*placeholderIdx), 1)
		(*placeholderIdx)++
	}
	return whereSQL, w.whereArgs
}

// tableClauseString holds shared table and error logic for builders with string table names.
type tableClauseString struct {
	table string
	err   error
}

func (t *tableClauseString) SetTable(table string) {
	if table == "" {
		t.err = errors.New("tableClauseString: table must be set")
	} else {
		t.table = table
	}
}

// tableClauseInterface holds shared table and error logic for builders with interface{} table names.
type tableClauseInterface struct {
	table interface{}
	err   error
}

func (t *tableClauseInterface) SetTable(table interface{}) {
	if table == nil || table == "" {
		t.err = errors.New("tableClauseInterface: table must be set")
	} else {
		t.table = table
	}
}
