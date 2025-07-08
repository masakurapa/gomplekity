package medium_complexity

import "fmt"

// ValidateInput は入力値を検証する関数 (複雑度: 6)
func ValidateInput(input string, minLen, maxLen int) bool {
	if input == "" {
		return false
	}
	if len(input) < minLen {
		return false
	}
	if len(input) > maxLen {
		return false
	}
	
	// 特殊文字チェック
	for _, char := range input {
		if char < 32 || char > 126 {
			return false
		}
	}
	
	return true
}

// CalculateDiscount は割引を計算する関数 (複雑度: 7)
func CalculateDiscount(price float64, customerType string, itemCount int) float64 {
	discount := 0.0
	
	if customerType == "premium" {
		discount = 0.2
	} else if customerType == "regular" {
		discount = 0.1
	}
	
	if itemCount >= 10 {
		discount += 0.05
	} else if itemCount >= 5 {
		discount += 0.02
	}
	
	if price > 1000 {
		discount += 0.03
	}
	
	return price * (1 - discount)
}