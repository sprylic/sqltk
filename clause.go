package cqb

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
	default:
		w.err = fmt.Errorf("Where: cond must be string or sq.Raw (got type %T)", cond)
	}
}

func (w *whereClause) WhereEqual(column string, value interface{}) {
	w.Where(column+" = ?", value)
}

func (w *whereClause) WhereNotEqual(column string, value interface{}) {
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

// InterpolateSQL returns the query with arguments interpolated for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func InterpolateSQL(query string, args []interface{}) string {
	if len(args) == 0 {
		return query
	}
	out := ""
	argIdx := 0
	for i := 0; i < len(query); i++ {
		if query[i] == '?' && argIdx < len(args) {
			out += formatArg(args[argIdx])
			argIdx++
		} else {
			out += string(query[i])
		}
	}
	return out
}

func formatArg(arg interface{}) string {
	switch v := arg.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case []byte:
		return "'" + strings.ReplaceAll(string(v), "'", "''") + "'"
	case int, int8, int16, int32, int64, float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}
