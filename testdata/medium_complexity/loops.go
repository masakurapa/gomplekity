package medium_complexity

import "fmt"

// FindPattern は配列内のパターンを探す関数 (複雑度: 8)
func FindPattern(data []int, pattern []int) []int {
	var positions []int
	
	for i := 0; i <= len(data)-len(pattern); i++ {
		match := true
		for j := 0; j < len(pattern); j++ {
			if data[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			positions = append(positions, i)
		}
	}
	
	return positions
}

// ProcessMatrix は行列を処理する関数 (複雑度: 9)
func ProcessMatrix(matrix [][]int) [][]int {
	result := make([][]int, len(matrix))
	
	for i := 0; i < len(matrix); i++ {
		result[i] = make([]int, len(matrix[i]))
		for j := 0; j < len(matrix[i]); j++ {
			if i == 0 || j == 0 {
				result[i][j] = matrix[i][j]
			} else if matrix[i][j] > 0 {
				result[i][j] = matrix[i][j] * 2
			} else if matrix[i][j] < 0 {
				result[i][j] = matrix[i][j] / 2
			} else {
				result[i][j] = 1
			}
		}
	}
	
	return result
}

// ValidateAndProcess は検証と処理を行う関数 (複雑度: 10)
func ValidateAndProcess(items []string) map[string]int {
	result := make(map[string]int)
	
	for _, item := range items {
		if item == "" {
			continue
		}
		
		if len(item) < 3 {
			result[item] = -1
			continue
		}
		
		score := 0
		for _, char := range item {
			if char >= 'a' && char <= 'z' {
				score += 1
			} else if char >= 'A' && char <= 'Z' {
				score += 2
			} else if char >= '0' && char <= '9' {
				score += 3
			}
		}
		
		if score > 50 {
			result[item] = 3
		} else if score > 20 {
			result[item] = 2
		} else {
			result[item] = 1
		}
	}
	
	return result
}