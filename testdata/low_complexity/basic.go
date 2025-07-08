package low_complexity

import "fmt"

// CheckPositive は正の数かチェックする関数 (複雑度: 2)
func CheckPositive(num int) bool {
	if num > 0 {
		return true
	}
	return false
}

// GetGrade は成績を判定する関数 (複雑度: 4)
func GetGrade(score int) string {
	if score >= 90 {
		return "A"
	} else if score >= 80 {
		return "B"
	} else if score >= 70 {
		return "C"
	}
	return "F"
}

// ProcessNumbers は数値を処理する関数 (複雑度: 3)
func ProcessNumbers(numbers []int) {
	for _, num := range numbers {
		if num%2 == 0 {
			fmt.Printf("Even: %d\n", num)
		} else {
			fmt.Printf("Odd: %d\n", num)
		}
	}
}