package sqldebug

import (
	"fmt"
	"strings"
)

type UnsafeSqlString string

func (s UnsafeSqlString) GetUnsafeString() string {
	return string(s)
}

// InterpolateSQL interpolates arguments into a SQL query for debugging/logging only.
// DO NOT use the result for execution (not safe against SQL injection).
func InterpolateSQL(query string, args []interface{}) UnsafeSqlString {
	if len(args) == 0 {
		return UnsafeSqlString(query)
	}

	// Simple interpolation - replace ? with values
	result := query
	argIndex := 0

	for i := 0; i < len(result) && argIndex < len(args); i++ {
		if result[i] == '?' {
			arg := args[argIndex]
			var argStr string

			switch v := arg.(type) {
			case string:
				argStr = "'" + strings.ReplaceAll(v, "'", "''") + "'"
			case nil:
				argStr = "NULL"
			default:
				argStr = fmt.Sprintf("%v", v)
			}

			result = result[:i] + argStr + result[i+1:]
			argIndex++
		}
	}

	return UnsafeSqlString(result)
}
