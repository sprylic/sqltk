package mysqlfunc

import (
	"fmt"
	"strings"

	"github.com/sprylic/sqltk/sqlfunc"
)

// Date and Time Functions
func CurrentTimestamp() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("CURRENT_TIMESTAMP")
}

func CurrentDate() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("CURRENT_DATE")
}

func CurrentTime() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("CURRENT_TIME")
}

func Now() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("NOW()")
}

func UnixTimestamp() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("UNIX_TIMESTAMP()")
}

// String Functions
func Concat(args ...interface{}) sqlfunc.SqlFunc {
	var argStrs []string
	for _, arg := range args {
		if err := sqlfunc.ValidateSqlFuncInput(arg); err != nil {
			// In production, you might want to panic or log this
			panic(fmt.Sprintf("Concat: %v", err))
		}
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}
	return sqlfunc.SqlFunc("CONCAT(" + strings.Join(argStrs, ", ") + ")")
}

func ConcatWs(separator string, args ...interface{}) sqlfunc.SqlFunc {
	var argStrs []string
	argStrs = append(argStrs, fmt.Sprintf("'%s'", separator))
	for _, arg := range args {
		if err := sqlfunc.ValidateSqlFuncInput(arg); err != nil {
			panic(fmt.Sprintf("ConcatWs: %v", err))
		}
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}
	return sqlfunc.SqlFunc("CONCAT_WS(" + strings.Join(argStrs, ", ") + ")")
}

func Substring(str interface{}, pos, length interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Substring: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(pos); err != nil {
		panic(fmt.Sprintf("Substring: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(length); err != nil {
		panic(fmt.Sprintf("Substring: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("SUBSTRING(%v, %v, %v)", str, pos, length))
}

func Left(str interface{}, length interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Left: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(length); err != nil {
		panic(fmt.Sprintf("Left: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("LEFT(%v, %v)", str, length))
}

func Right(str interface{}, length interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Right: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(length); err != nil {
		panic(fmt.Sprintf("Right: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("RIGHT(%v, %v)", str, length))
}

func Upper(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Upper: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("UPPER(%v)", str))
}

func Lower(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Lower: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("LOWER(%v)", str))
}

func Trim(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Trim: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("TRIM(%v)", str))
}

func Ltrim(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Ltrim: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("LTRIM(%v)", str))
}

func Rtrim(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Rtrim: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("RTRIM(%v)", str))
}

func Replace(str, from, to interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Replace: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(from); err != nil {
		panic(fmt.Sprintf("Replace: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(to); err != nil {
		panic(fmt.Sprintf("Replace: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("REPLACE(%v, %v, %v)", str, from, to))
}

func Reverse(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Reverse: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("REVERSE(%v)", str))
}

func Length(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Length: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("LENGTH(%v)", str))
}

func CharLength(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("CharLength: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("CHAR_LENGTH(%v)", str))
}

// Numeric Functions
func Abs(num interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("ABS(%v)", num))
}

func Ceiling(num interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("CEILING(%v)", num))
}

func Floor(num interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("FLOOR(%v)", num))
}

func Round(num interface{}, decimals ...interface{}) sqlfunc.SqlFunc {
	if len(decimals) > 0 {
		return sqlfunc.SqlFunc(fmt.Sprintf("ROUND(%v, %v)", num, decimals[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("ROUND(%v)", num))
}

func Truncate(num, decimals interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("TRUNCATE(%v, %v)", num, decimals))
}

func Mod(dividend, divisor interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("MOD(%v, %v)", dividend, divisor))
}

func Power(base, exponent interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("POWER(%v, %v)", base, exponent))
}

func Sqrt(num interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("SQRT(%v)", num))
}

func Random() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("RAND()")
}

// Aggregate Functions
func Count(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("COUNT(%v)", expr))
}

func Sum(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("SUM(%v)", expr))
}

func Avg(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("AVG(%v)", expr))
}

func Min(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("MIN(%v)", expr))
}

func Max(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("MAX(%v)", expr))
}

func GroupConcat(expr interface{}, separator ...string) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(expr); err != nil {
		panic(fmt.Sprintf("GroupConcat: %v", err))
	}
	if len(separator) > 0 {
		return sqlfunc.SqlFunc(fmt.Sprintf("GROUP_CONCAT(%v SEPARATOR '%s')", expr, separator[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("GROUP_CONCAT(%v)", expr))
}

// Conditional Functions
func If(condition, trueVal, falseVal interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(condition); err != nil {
		panic(fmt.Sprintf("If: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(trueVal); err != nil {
		panic(fmt.Sprintf("If: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(falseVal); err != nil {
		panic(fmt.Sprintf("If: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("IF(%v, %v, %v)", condition, trueVal, falseVal))
}

func IfNull(expr, nullVal interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(expr); err != nil {
		panic(fmt.Sprintf("IfNull: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(nullVal); err != nil {
		panic(fmt.Sprintf("IfNull: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("IFNULL(%v, %v)", expr, nullVal))
}

func NullIf(expr1, expr2 interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(expr1); err != nil {
		panic(fmt.Sprintf("NullIf: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(expr2); err != nil {
		panic(fmt.Sprintf("NullIf: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("NULLIF(%v, %v)", expr1, expr2))
}

// Type Conversion Functions
func Cast(expr interface{}, asType string) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(expr); err != nil {
		panic(fmt.Sprintf("Cast: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("CAST(%v AS %s)", expr, asType))
}

func Convert(expr interface{}, asType string) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(expr); err != nil {
		panic(fmt.Sprintf("Convert: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("CONVERT(%v, %s)", expr, asType))
}

// JSON Functions (MySQL 5.7+)
func JsonExtract(jsonDoc, path interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonDoc); err != nil {
		panic(fmt.Sprintf("JsonExtract: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(path); err != nil {
		panic(fmt.Sprintf("JsonExtract: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("JSON_EXTRACT(%v, %v)", jsonDoc, path))
}

func JsonUnquote(jsonVal interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonVal); err != nil {
		panic(fmt.Sprintf("JsonUnquote: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("JSON_UNQUOTE(%v)", jsonVal))
}

func JsonLength(jsonDoc interface{}, path ...interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonDoc); err != nil {
		panic(fmt.Sprintf("JsonLength: %v", err))
	}
	if len(path) > 0 {
		if err := sqlfunc.ValidateSqlFuncInput(path[0]); err != nil {
			panic(fmt.Sprintf("JsonLength: %v", err))
		}
		return sqlfunc.SqlFunc(fmt.Sprintf("JSON_LENGTH(%v, %v)", jsonDoc, path[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("JSON_LENGTH(%v)", jsonDoc))
}

func JsonKeys(jsonDoc interface{}, path ...interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonDoc); err != nil {
		panic(fmt.Sprintf("JsonKeys: %v", err))
	}
	if len(path) > 0 {
		if err := sqlfunc.ValidateSqlFuncInput(path[0]); err != nil {
			panic(fmt.Sprintf("JsonKeys: %v", err))
		}
		return sqlfunc.SqlFunc(fmt.Sprintf("JSON_KEYS(%v, %v)", jsonDoc, path[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("JSON_KEYS(%v)", jsonDoc))
}

func JsonContains(jsonDoc, val interface{}, path ...interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonDoc); err != nil {
		panic(fmt.Sprintf("JsonContains: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(val); err != nil {
		panic(fmt.Sprintf("JsonContains: %v", err))
	}
	if len(path) > 0 {
		if err := sqlfunc.ValidateSqlFuncInput(path[0]); err != nil {
			panic(fmt.Sprintf("JsonContains: %v", err))
		}
		return sqlfunc.SqlFunc(fmt.Sprintf("JSON_CONTAINS(%v, %v, %v)", jsonDoc, val, path[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("JSON_CONTAINS(%v, %v)", jsonDoc, val))
}

func JsonSearch(jsonDoc, oneOrAll, searchStr interface{}, escapeChar ...interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonDoc); err != nil {
		panic(fmt.Sprintf("JsonSearch: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(oneOrAll); err != nil {
		panic(fmt.Sprintf("JsonSearch: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(searchStr); err != nil {
		panic(fmt.Sprintf("JsonSearch: %v", err))
	}
	if len(escapeChar) > 0 {
		if err := sqlfunc.ValidateSqlFuncInput(escapeChar[0]); err != nil {
			panic(fmt.Sprintf("JsonSearch: %v", err))
		}
		return sqlfunc.SqlFunc(fmt.Sprintf("JSON_SEARCH(%v, %v, %v, %v)", jsonDoc, oneOrAll, searchStr, escapeChar[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("JSON_SEARCH(%v, %v, %v)", jsonDoc, oneOrAll, searchStr))
}

// Encryption Functions
func Md5(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Md5: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("MD5(%v)", str))
}

func Sha1(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Sha1: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("SHA1(%v)", str))
}

func Sha2(str, hashLength interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Sha2: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(hashLength); err != nil {
		panic(fmt.Sprintf("Sha2: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("SHA2(%v, %v)", str, hashLength))
}

func Password(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Password: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("PASSWORD(%v)", str))
}

// Information Functions
func LastInsertId() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("LAST_INSERT_ID()")
}

func RowCount() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("ROW_COUNT()")
}

func FoundRows() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("FOUND_ROWS()")
}

func ConnectionId() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("CONNECTION_ID()")
}

func Database() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("DATABASE()")
}

func User() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("USER()")
}

func Version() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("VERSION()")
}

// Date/Time Manipulation Functions
func DateAdd(date interface{}, interval string, expr interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("DateAdd: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(expr); err != nil {
		panic(fmt.Sprintf("DateAdd: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("DATE_ADD(%v, INTERVAL %v %s)", date, expr, interval))
}

func DateSub(date interface{}, interval string, expr interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("DateSub: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(expr); err != nil {
		panic(fmt.Sprintf("DateSub: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("DATE_SUB(%v, INTERVAL %v %s)", date, expr, interval))
}

func DateDiff(date1, date2 interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date1); err != nil {
		panic(fmt.Sprintf("DateDiff: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(date2); err != nil {
		panic(fmt.Sprintf("DateDiff: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("DATEDIFF(%v, %v)", date1, date2))
}

func TimeDiff(time1, time2 interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(time1); err != nil {
		panic(fmt.Sprintf("TimeDiff: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(time2); err != nil {
		panic(fmt.Sprintf("TimeDiff: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("TIMEDIFF(%v, %v)", time1, time2))
}

func Year(date interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("Year: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("YEAR(%v)", date))
}

func Month(date interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("Month: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("MONTH(%v)", date))
}

func Day(date interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("Day: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("DAY(%v)", date))
}

func Hour(time interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(time); err != nil {
		panic(fmt.Sprintf("Hour: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("HOUR(%v)", time))
}

func Minute(time interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(time); err != nil {
		panic(fmt.Sprintf("Minute: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("MINUTE(%v)", time))
}

func Second(time interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(time); err != nil {
		panic(fmt.Sprintf("Second: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("SECOND(%v)", time))
}

func DayOfWeek(date interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("DayOfWeek: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("DAYOFWEEK(%v)", date))
}

func DayOfYear(date interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("DayOfYear: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("DAYOFYEAR(%v)", date))
}

func Week(date interface{}, mode ...interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("Week: %v", err))
	}
	if len(mode) > 0 {
		if err := sqlfunc.ValidateSqlFuncInput(mode[0]); err != nil {
			panic(fmt.Sprintf("Week: %v", err))
		}
		return sqlfunc.SqlFunc(fmt.Sprintf("WEEK(%v, %v)", date, mode[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("WEEK(%v)", date))
}

func MonthName(date interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("MonthName: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("MONTHNAME(%v)", date))
}

func DayName(date interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("DayName: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("DAYNAME(%v)", date))
}

// Formatting Functions
func DateFormat(date, format interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(date); err != nil {
		panic(fmt.Sprintf("DateFormat: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(format); err != nil {
		panic(fmt.Sprintf("DateFormat: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("DATE_FORMAT(%v, %v)", date, format))
}

func TimeFormat(time, format interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(time); err != nil {
		panic(fmt.Sprintf("TimeFormat: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(format); err != nil {
		panic(fmt.Sprintf("TimeFormat: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("TIME_FORMAT(%v, %v)", time, format))
}

func Format(num, decimals interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(num); err != nil {
		panic(fmt.Sprintf("Format: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(decimals); err != nil {
		panic(fmt.Sprintf("Format: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("FORMAT(%v, %v)", num, decimals))
}

// String Search Functions
func Like(str, pattern interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Like: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(pattern); err != nil {
		panic(fmt.Sprintf("Like: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("%v LIKE %v", str, pattern))
}

func Regexp(str, pattern interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Regexp: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(pattern); err != nil {
		panic(fmt.Sprintf("Regexp: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("%v REGEXP %v", str, pattern))
}

func Soundex(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Soundex: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("SOUNDEX(%v)", str))
}

// Mathematical Constants
func Pi() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("PI()")
}

func E() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("E()")
}
