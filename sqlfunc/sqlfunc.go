package sqlfunc

import (
	"fmt"
	"regexp"
	"strings"
)

type SqlFunc string

// ValidateSqlFuncInput checks if input is safe for SQL function generation
func ValidateSqlFuncInput(input interface{}) error {
	switch v := input.(type) {
	case string:
		return validatePattern(v)
	case SqlFunc:
		return validatePattern(string(v))
	default:
		// Other types (numbers, etc.) are generally safe
		return nil
	}
}

func validatePattern(pattern string) error {
	// Check for dangerous patterns that suggest SQL injection
	if strings.Contains(pattern, "';") || strings.Contains(pattern, "/*") || strings.Contains(pattern, "*/") {
		return fmt.Errorf("potentially unsafe input detected: %q", pattern)
	}
	// Check for SQL injection patterns with proper context
	// Only flag if it looks like a complete SQL statement, not just column names
	if regexp.MustCompile(`(?i)^(union|select|insert|update|delete|drop|create|alter|exec|execute)\s+`).MatchString(strings.TrimSpace(pattern)) {
		return fmt.Errorf("input contains SQL keywords that suggest injection: %q", pattern)
	}
	// Check for comment patterns that could be used for injection
	if strings.Contains(pattern, "--") || strings.Contains(pattern, "#") {
		return fmt.Errorf("input contains comment markers: %q", pattern)
	}
	return nil
}
