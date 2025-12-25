package raw

// Raw represents raw SQL that should be included directly without quoting.
type Raw string

// BuildCondition implements the Condition interface.
func (r Raw) BuildCondition() (string, []interface{}, error) {
	return string(r), nil, nil
}

// RawCondition wraps a Raw SQL condition for type safety.
type RawCondition struct {
	SQL Raw
}

func Cond(sql string) *RawCondition {
	return &RawCondition{SQL: Raw(sql)}
}

// BuildCondition implements the Condition interface.
func (rc *RawCondition) BuildCondition() (string, []interface{}, error) {
	return string(rc.SQL), nil, nil
}
