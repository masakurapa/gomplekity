package high_complexity

import (
	"fmt"
	"sort"
	"strings"
)

// SuperComplexProcessor は非常に複雑な処理を行う関数 (複雑度: 18)
func SuperComplexProcessor(input []map[string]interface{}, config map[string]string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	
	for _, item := range input {
		processedItem := make(map[string]interface{})
		
		for key, value := range item {
			configKey := "process_" + key
			if processingRule, exists := config[configKey]; exists {
				switch processingRule {
				case "uppercase":
					if str, ok := value.(string); ok {
						processedItem[key] = strings.ToUpper(str)
					} else {
						return nil, fmt.Errorf("cannot uppercase non-string value for key %s", key)
					}
				case "lowercase":
					if str, ok := value.(string); ok {
						processedItem[key] = strings.ToLower(str)
					} else {
						return nil, fmt.Errorf("cannot lowercase non-string value for key %s", key)
					}
				case "double":
					if num, ok := value.(float64); ok {
						processedItem[key] = num * 2
					} else if intVal, ok := value.(int); ok {
						processedItem[key] = intVal * 2
					} else {
						return nil, fmt.Errorf("cannot double non-numeric value for key %s", key)
					}
				case "validate":
					if str, ok := value.(string); ok {
						if len(str) < 3 {
							return nil, fmt.Errorf("validation failed for key %s: too short", key)
						}
						if len(str) > 50 {
							return nil, fmt.Errorf("validation failed for key %s: too long", key)
						}
						if strings.Contains(str, "invalid") {
							return nil, fmt.Errorf("validation failed for key %s: contains invalid content", key)
						}
						processedItem[key] = str
					} else {
						return nil, fmt.Errorf("cannot validate non-string value for key %s", key)
					}
				case "sort":
					if slice, ok := value.([]interface{}); ok {
						var sortedSlice []interface{}
						for _, elem := range slice {
							if str, ok := elem.(string); ok {
								sortedSlice = append(sortedSlice, str)
							} else {
								return nil, fmt.Errorf("cannot sort non-string elements for key %s", key)
							}
						}
						sort.Slice(sortedSlice, func(i, j int) bool {
							return sortedSlice[i].(string) < sortedSlice[j].(string)
						})
						processedItem[key] = sortedSlice
					} else {
						return nil, fmt.Errorf("cannot sort non-slice value for key %s", key)
					}
				default:
					processedItem[key] = value
				}
			} else {
				processedItem[key] = value
			}
		}
		
		// 追加の後処理
		if len(processedItem) == 0 {
			continue
		}
		
		// 必須フィールドチェック
		requiredFields := []string{"id", "name"}
		for _, field := range requiredFields {
			if _, exists := processedItem[field]; !exists {
				return nil, fmt.Errorf("missing required field: %s", field)
			}
		}
		
		result = append(result, processedItem)
	}
	
	return result, nil
}