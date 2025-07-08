package low_complexity

import "fmt"

// SimpleFunction は最も単純な関数 (複雑度: 1)
func SimpleFunction() {
	fmt.Println("Hello, World!")
}

// AddTwoNumbers は2つの数値を足すだけの関数 (複雑度: 1)
func AddTwoNumbers(a, b int) int {
	return a + b
}

// GreetUser は単純な挨拶関数 (複雑度: 1)
func GreetUser(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}