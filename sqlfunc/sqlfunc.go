package sqlfunc

import (
	"fmt"
	"regexp"
	"strings"
)

type SqlFunc string

// validateSqlFuncInput checks if input is safe for SQL function generation
func ValidateSqlFuncInput(input interface{}) error {
	switch v := input.(type) {
	case string:
		// Check for dangerous patterns that suggest SQL injection
		if strings.Contains(v, "';") || strings.Contains(v, "/*") || strings.Contains(v, "*/") {
			return fmt.Errorf("potentially unsafe input detected: %q", v)
		}
		// Check for SQL injection patterns with proper context
		// Only flag if it looks like a complete SQL statement, not just column names
		if regexp.MustCompile(`(?i)^(union|select|insert|update|delete|drop|create|alter|exec|execute)\s+`).MatchString(strings.TrimSpace(v)) {
			return fmt.Errorf("input contains SQL keywords that suggest injection: %q", v)
		}
		// Check for comment patterns that could be used for injection
		if strings.Contains(v, "--") || strings.Contains(v, "#") {
			return fmt.Errorf("input contains comment markers: %q", v)
		}
	case SqlFunc:
		// For SqlFunc objects, check the string content
		if strings.Contains(string(v), "';") || strings.Contains(string(v), "/*") || strings.Contains(string(v), "*/") {
			return fmt.Errorf("potentially unsafe input detected: %q", v)
		}
		if regexp.MustCompile(`(?i)^(union|select|insert|update|delete|drop|create|alter|exec|execute)\s+`).MatchString(strings.TrimSpace(string(v))) {
			return fmt.Errorf("input contains SQL keywords that suggest injection: %q", v)
		}
		if strings.Contains(string(v), "--") || strings.Contains(string(v), "#") {
			return fmt.Errorf("input contains comment markers: %q", v)
		}
		return nil
	default:
		// Other types (numbers, etc.) are generally safe
		return nil
	}
	return nil
}
