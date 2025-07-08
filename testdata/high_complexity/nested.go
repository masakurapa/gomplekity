package high_complexity

import (
	"fmt"
	"strings"
)

// ComplexValidation は複雑な検証を行う関数 (複雑度: 14)
func ComplexValidation(data map[string]interface{}, rules map[string]string) []string {
	var errors []string
	
	for key, rule := range rules {
		value, exists := data[key]
		if !exists {
			errors = append(errors, fmt.Sprintf("Missing field: %s", key))
			continue
		}
		
		switch rule {
		case "required":
			if value == nil || value == "" {
				errors = append(errors, fmt.Sprintf("Field %s is required", key))
			}
		case "string":
			if str, ok := value.(string); ok {
				if len(str) == 0 {
					errors = append(errors, fmt.Sprintf("Field %s cannot be empty", key))
				} else if len(str) > 255 {
					errors = append(errors, fmt.Sprintf("Field %s is too long", key))
				}
			} else {
				errors = append(errors, fmt.Sprintf("Field %s must be a string", key))
			}
		case "number":
			if num, ok := value.(float64); ok {
				if num < 0 {
					errors = append(errors, fmt.Sprintf("Field %s must be positive", key))
				} else if num > 1000000 {
					errors = append(errors, fmt.Sprintf("Field %s is too large", key))
				}
			} else {
				errors = append(errors, fmt.Sprintf("Field %s must be a number", key))
			}
		case "email":
			if str, ok := value.(string); ok {
				if !strings.Contains(str, "@") {
					errors = append(errors, fmt.Sprintf("Field %s must be a valid email", key))
				} else if strings.Count(str, "@") != 1 {
					errors = append(errors, fmt.Sprintf("Field %s must contain exactly one @", key))
				}
			} else {
				errors = append(errors, fmt.Sprintf("Field %s must be a string", key))
			}
		}
	}
	
	return errors
}