package pgfunc

import (
	"fmt"
	"strings"

	"github.com/sprylic/sqltk/sqlfunc"
)

// Date and Time Functions
func Now() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("now()")
}

func CurrentTimestamp() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("CURRENT_TIMESTAMP")
}

func CurrentDate() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("CURRENT_DATE")
}

func CurrentTime() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("CURRENT_TIME")
}

func ClockTimestamp() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("clock_timestamp()")
}

func StatementTimestamp() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("statement_timestamp()")
}

func TransactionTimestamp() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("transaction_timestamp()")
}

func Extract(field interface{}, source interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("EXTRACT(%v FROM %v)", field, source))
}

func DatePart(field interface{}, source interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("date_part(%v, %v)", field, source))
}

func DateTrunc(field interface{}, source interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("date_trunc(%v, %v)", field, source))
}

func Age(timestamp interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("age(%v)", timestamp))
}

func AgeWithEnd(timestamp1, timestamp2 interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("age(%v, %v)", timestamp1, timestamp2))
}

// String Functions
func Concat(args ...interface{}) sqlfunc.SqlFunc {
	var argStrs []string
	for _, arg := range args {
		if err := sqlfunc.ValidateSqlFuncInput(arg); err != nil {
			panic(fmt.Sprintf("Concat: %v", err))
		}
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}
	return sqlfunc.SqlFunc("concat(" + strings.Join(argStrs, ", ") + ")")
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
	return sqlfunc.SqlFunc("concat_ws(" + strings.Join(argStrs, ", ") + ")")
}

func Substring(str interface{}, from interface{}, forArg ...interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Substring: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(from); err != nil {
		panic(fmt.Sprintf("Substring: %v", err))
	}
	if len(forArg) > 0 {
		if err := sqlfunc.ValidateSqlFuncInput(forArg[0]); err != nil {
			panic(fmt.Sprintf("Substring: %v", err))
		}
		return sqlfunc.SqlFunc(fmt.Sprintf("substring(%v from %v for %v)", str, from, forArg[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("substring(%v from %v)", str, from))
}

func Left(str interface{}, n interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Left: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(n); err != nil {
		panic(fmt.Sprintf("Left: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("left(%v, %v)", str, n))
}

func Right(str interface{}, n interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Right: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(n); err != nil {
		panic(fmt.Sprintf("Right: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("right(%v, %v)", str, n))
}

func Upper(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Upper: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("upper(%v)", str))
}

func Lower(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Lower: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("lower(%v)", str))
}

func Initcap(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Initcap: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("initcap(%v)", str))
}

func Trim(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Trim: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("trim(%v)", str))
}

func Ltrim(str interface{}, chars ...interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Ltrim: %v", err))
	}
	if len(chars) > 0 {
		if err := sqlfunc.ValidateSqlFuncInput(chars[0]); err != nil {
			panic(fmt.Sprintf("Ltrim: %v", err))
		}
		return sqlfunc.SqlFunc(fmt.Sprintf("ltrim(%v, %v)", str, chars[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("ltrim(%v)", str))
}

func Rtrim(str interface{}, chars ...interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Rtrim: %v", err))
	}
	if len(chars) > 0 {
		if err := sqlfunc.ValidateSqlFuncInput(chars[0]); err != nil {
			panic(fmt.Sprintf("Rtrim: %v", err))
		}
		return sqlfunc.SqlFunc(fmt.Sprintf("rtrim(%v, %v)", str, chars[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("rtrim(%v)", str))
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
	return sqlfunc.SqlFunc(fmt.Sprintf("replace(%v, %v, %v)", str, from, to))
}

func Reverse(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Reverse: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("reverse(%v)", str))
}

func Length(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Length: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("length(%v)", str))
}

func CharLength(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("CharLength: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("char_length(%v)", str))
}

func Position(substring, string interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(substring); err != nil {
		panic(fmt.Sprintf("Position: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(string); err != nil {
		panic(fmt.Sprintf("Position: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("position(%v in %v)", substring, string))
}

func Substr(str interface{}, from interface{}, count ...interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Substr: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(from); err != nil {
		panic(fmt.Sprintf("Substr: %v", err))
	}
	if len(count) > 0 {
		if err := sqlfunc.ValidateSqlFuncInput(count[0]); err != nil {
			panic(fmt.Sprintf("Substr: %v", err))
		}
		return sqlfunc.SqlFunc(fmt.Sprintf("substr(%v, %v, %v)", str, from, count[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("substr(%v, %v)", str, from))
}

// Numeric Functions
func Abs(num interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("abs(%v)", num))
}

func Ceiling(num interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("ceiling(%v)", num))
}

func Floor(num interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("floor(%v)", num))
}

func Round(num interface{}, decimals ...interface{}) sqlfunc.SqlFunc {
	if len(decimals) > 0 {
		return sqlfunc.SqlFunc(fmt.Sprintf("round(%v, %v)", num, decimals[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("round(%v)", num))
}

func Trunc(num interface{}, decimals ...interface{}) sqlfunc.SqlFunc {
	if len(decimals) > 0 {
		return sqlfunc.SqlFunc(fmt.Sprintf("trunc(%v, %v)", num, decimals[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("trunc(%v)", num))
}

func Mod(dividend, divisor interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("mod(%v, %v)", dividend, divisor))
}

func Power(base, exponent interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("power(%v, %v)", base, exponent))
}

func Sqrt(num interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("sqrt(%v)", num))
}

func Random() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("random()")
}

func Pi() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("pi()")
}

func E() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("e()")
}

// Aggregate Functions
func Count(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("count(%v)", expr))
}

func Sum(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("sum(%v)", expr))
}

func Avg(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("avg(%v)", expr))
}

func Min(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("min(%v)", expr))
}

func Max(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("max(%v)", expr))
}

func StringAgg(expr interface{}, delimiter string) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("string_agg(%v, '%s')", expr, delimiter))
}

func ArrayAgg(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("array_agg(%v)", expr))
}

func JsonAgg(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("json_agg(%v)", expr))
}

func JsonbAgg(expr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("jsonb_agg(%v)", expr))
}

// Conditional Functions
func Coalesce(args ...interface{}) sqlfunc.SqlFunc {
	var argStrs []string
	for _, arg := range args {
		if err := sqlfunc.ValidateSqlFuncInput(arg); err != nil {
			panic(fmt.Sprintf("Coalesce: %v", err))
		}
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}
	return sqlfunc.SqlFunc("coalesce(" + strings.Join(argStrs, ", ") + ")")
}

func NullIf(expr1, expr2 interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(expr1); err != nil {
		panic(fmt.Sprintf("NullIf: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(expr2); err != nil {
		panic(fmt.Sprintf("NullIf: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("nullif(%v, %v)", expr1, expr2))
}

func Greatest(args ...interface{}) sqlfunc.SqlFunc {
	var argStrs []string
	for _, arg := range args {
		if err := sqlfunc.ValidateSqlFuncInput(arg); err != nil {
			panic(fmt.Sprintf("Greatest: %v", err))
		}
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}
	return sqlfunc.SqlFunc("greatest(" + strings.Join(argStrs, ", ") + ")")
}

func Least(args ...interface{}) sqlfunc.SqlFunc {
	var argStrs []string
	for _, arg := range args {
		if err := sqlfunc.ValidateSqlFuncInput(arg); err != nil {
			panic(fmt.Sprintf("Least: %v", err))
		}
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}
	return sqlfunc.SqlFunc("least(" + strings.Join(argStrs, ", ") + ")")
}

// Type Conversion Functions
func Cast(expr interface{}, asType string) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(expr); err != nil {
		panic(fmt.Sprintf("Cast: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("cast(%v as %s)", expr, asType))
}

func Convert(expr interface{}, asType string) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(expr); err != nil {
		panic(fmt.Sprintf("Convert: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("convert(%v, %s)", expr, asType))
}

// JSON Functions
func JsonExtract(jsonDoc, path interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonDoc); err != nil {
		panic(fmt.Sprintf("JsonExtract: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(path); err != nil {
		panic(fmt.Sprintf("JsonExtract: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("json_extract_path_text(%v, %v)", jsonDoc, path))
}

func JsonExtractPath(jsonDoc, path interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonDoc); err != nil {
		panic(fmt.Sprintf("JsonExtractPath: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(path); err != nil {
		panic(fmt.Sprintf("JsonExtractPath: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("json_extract_path(%v, %v)", jsonDoc, path))
}

func JsonbExtractPath(jsonDoc, path interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonDoc); err != nil {
		panic(fmt.Sprintf("JsonbExtractPath: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(path); err != nil {
		panic(fmt.Sprintf("JsonbExtractPath: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("jsonb_extract_path(%v, %v)", jsonDoc, path))
}

func JsonTypeof(jsonVal interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonVal); err != nil {
		panic(fmt.Sprintf("JsonTypeof: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("json_typeof(%v)", jsonVal))
}

func JsonbTypeof(jsonVal interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonVal); err != nil {
		panic(fmt.Sprintf("JsonbTypeof: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("jsonb_typeof(%v)", jsonVal))
}

func JsonLength(jsonVal interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonVal); err != nil {
		panic(fmt.Sprintf("JsonLength: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("json_array_length(%v)", jsonVal))
}

func JsonbLength(jsonVal interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonVal); err != nil {
		panic(fmt.Sprintf("JsonbLength: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("jsonb_array_length(%v)", jsonVal))
}

func JsonKeys(jsonVal interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonVal); err != nil {
		panic(fmt.Sprintf("JsonKeys: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("json_object_keys(%v)", jsonVal))
}

func JsonbKeys(jsonVal interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonVal); err != nil {
		panic(fmt.Sprintf("JsonbKeys: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("jsonb_object_keys(%v)", jsonVal))
}

func JsonContains(jsonDoc, val interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonDoc); err != nil {
		panic(fmt.Sprintf("JsonContains: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(val); err != nil {
		panic(fmt.Sprintf("JsonContains: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("json_contains(%v, %v)", jsonDoc, val))
}

func JsonbContains(jsonDoc, val interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(jsonDoc); err != nil {
		panic(fmt.Sprintf("JsonbContains: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(val); err != nil {
		panic(fmt.Sprintf("JsonbContains: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("jsonb_contains(%v, %v)", jsonDoc, val))
}

// Array Functions
func ArrayLength(arr interface{}, dim ...interface{}) sqlfunc.SqlFunc {
	if len(dim) > 0 {
		return sqlfunc.SqlFunc(fmt.Sprintf("array_length(%v, %v)", arr, dim[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("array_length(%v)", arr))
}

func ArrayUpper(arr interface{}, dim ...interface{}) sqlfunc.SqlFunc {
	if len(dim) > 0 {
		return sqlfunc.SqlFunc(fmt.Sprintf("array_upper(%v, %v)", arr, dim[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("array_upper(%v)", arr))
}

func ArrayLower(arr interface{}, dim ...interface{}) sqlfunc.SqlFunc {
	if len(dim) > 0 {
		return sqlfunc.SqlFunc(fmt.Sprintf("array_lower(%v, %v)", arr, dim[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("array_lower(%v)", arr))
}

func ArrayDims(arr interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("array_dims(%v)", arr))
}

func ArrayToString(arr interface{}, delimiter string) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("array_to_string(%v, '%s')", arr, delimiter))
}

func StringToArray(str, delimiter interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("string_to_array(%v, %v)", str, delimiter))
}

// Encryption Functions
func Md5(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Md5: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("md5(%v)", str))
}

func Sha256(str interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(str); err != nil {
		panic(fmt.Sprintf("Sha256: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("sha256(%v)", str))
}

func Crypt(password, salt interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(password); err != nil {
		panic(fmt.Sprintf("Crypt: %v", err))
	}
	if err := sqlfunc.ValidateSqlFuncInput(salt); err != nil {
		panic(fmt.Sprintf("Crypt: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("crypt(%v, %v)", password, salt))
}

func GenSalt(method interface{}) sqlfunc.SqlFunc {
	if err := sqlfunc.ValidateSqlFuncInput(method); err != nil {
		panic(fmt.Sprintf("GenSalt: %v", err))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("gen_salt('%s')", method))
}

// Information Functions
func CurrentDatabase() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("current_database()")
}

func CurrentUser() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("current_user")
}

func SessionUser() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("session_user")
}

func Version() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("version()")
}

func PgVersion() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("pg_version()")
}

func PgVersionNum() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("pg_version_num()")
}

func PgBackendPid() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("pg_backend_pid()")
}

func PgPostmasterStartTime() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("pg_postmaster_start_time()")
}

// Text Search Functions
func ToTsvector(config, document interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("to_tsvector(%v, %v)", config, document))
}

func ToTsquery(config, query interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("to_tsquery(%v, %v)", config, query))
}

func PlainToTsquery(config, query interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("plainto_tsquery(%v, %v)", config, query))
}

func WebsearchToTsquery(config, query interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("websearch_to_tsquery(%v, %v)", config, query))
}

func TsRank(vector, query interface{}, weights ...interface{}) sqlfunc.SqlFunc {
	if len(weights) > 0 {
		return sqlfunc.SqlFunc(fmt.Sprintf("ts_rank(%v, %v, %v)", vector, query, weights[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("ts_rank(%v, %v)", vector, query))
}

func TsRankCd(vector, query interface{}, weights ...interface{}) sqlfunc.SqlFunc {
	if len(weights) > 0 {
		return sqlfunc.SqlFunc(fmt.Sprintf("ts_rank_cd(%v, %v, %v)", vector, query, weights[0]))
	}
	return sqlfunc.SqlFunc(fmt.Sprintf("ts_rank_cd(%v, %v)", vector, query))
}

// UUID Functions
func GenRandomUuid() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("gen_random_uuid()")
}

func UuidGenerateV4() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("uuid_generate_v4()")
}

// Network Functions
func InetClientAddr() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("inet_client_addr()")
}

func InetClientPort() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("inet_client_port()")
}

func InetServerAddr() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("inet_server_addr()")
}

func InetServerPort() sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc("inet_server_port()")
}

// Formatting Functions
func ToChar(timestamp, format interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("to_char(%v, %v)", timestamp, format))
}

func ToDate(str, format interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("to_date(%v, %v)", str, format))
}

func ToTimestamp(str, format interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("to_timestamp(%v, %v)", str, format))
}

func ToNumber(str, format interface{}) sqlfunc.SqlFunc {
	return sqlfunc.SqlFunc(fmt.Sprintf("to_number(%v, %v)", str, format))
}
