package main

import (
	"os"
)

// WriteFile writes content to a file
func WriteFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}