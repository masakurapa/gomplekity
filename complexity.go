package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/fzipp/gocyclo"
)

// FunctionComplexity represents the complexity of a single function
type FunctionComplexity struct {
	Name       string
	File       string
	Line       int
	Column     int
	Complexity int
}

// ComplexityAnalyzer analyzes the cyclomatic complexity of Go files
type ComplexityAnalyzer struct {
	lowThreshold    int
	mediumThreshold int
}

// NewComplexityAnalyzer creates a new complexity analyzer
func NewComplexityAnalyzer(lowThreshold, mediumThreshold int) *ComplexityAnalyzer {
	return &ComplexityAnalyzer{
		lowThreshold:    lowThreshold,
		mediumThreshold: mediumThreshold,
	}
}

// AnalyzeDirectory analyzes all Go files in the given directory
func (ca *ComplexityAnalyzer) AnalyzeDirectory(dir string) ([]FunctionComplexity, error) {
	var functions []FunctionComplexity

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files for now
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		funcs, err := ca.analyzeFile(path)
		if err != nil {
			return fmt.Errorf("failed to analyze file %s: %w", path, err)
		}

		functions = append(functions, funcs...)
		return nil
	})

	return functions, err
}

// analyzeFile analyzes a single Go file
func (ca *ComplexityAnalyzer) analyzeFile(filename string) ([]FunctionComplexity, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var stats gocyclo.Stats
	stats = gocyclo.AnalyzeASTFile(node, fset, stats)
	
	var functions []FunctionComplexity
	for _, stat := range stats {
		functions = append(functions, FunctionComplexity{
			Name:       stat.FuncName,
			File:       filename,
			Line:       stat.Pos.Line,
			Column:     stat.Pos.Column,
			Complexity: stat.Complexity,
		})
	}

	return functions, nil
}

// GetComplexityLevel returns the complexity level based on thresholds
func (ca *ComplexityAnalyzer) GetComplexityLevel(complexity int) string {
	if complexity <= ca.lowThreshold {
		return "low"
	} else if complexity <= ca.mediumThreshold {
		return "medium"
	}
	return "high"
}

// GetComplexityColor returns the color for the complexity level
func (ca *ComplexityAnalyzer) GetComplexityColor(complexity int) string {
	switch ca.GetComplexityLevel(complexity) {
	case "low":
		return "green"
	case "medium":
		return "yellow"
	case "high":
		return "red"
	default:
		return "gray"
	}
}

// PrintComplexityReport prints a formatted complexity report
func (ca *ComplexityAnalyzer) PrintComplexityReport(functions []FunctionComplexity) {
	fmt.Printf("🌳 Complexity Analysis Report\n")
	fmt.Printf("================================\n")
	fmt.Printf("Thresholds: Low ≤ %d, Medium ≤ %d, High > %d\n\n", 
		ca.lowThreshold, ca.mediumThreshold, ca.mediumThreshold)

	lowCount, mediumCount, highCount := 0, 0, 0

	for _, fn := range functions {
		level := ca.GetComplexityLevel(fn.Complexity)
		
		var emoji string
		switch level {
		case "low":
			emoji = "🟢"
			lowCount++
		case "medium":
			emoji = "🟡"
			mediumCount++
		case "high":
			emoji = "🔴"
			highCount++
		}

		fmt.Printf("%s %s (%s): %d - %s:%d\n", 
			emoji, fn.Name, level, fn.Complexity, fn.File, fn.Line)
	}

	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("🟢 Low complexity: %d functions\n", lowCount)
	fmt.Printf("🟡 Medium complexity: %d functions\n", mediumCount)
	fmt.Printf("🔴 High complexity: %d functions\n", highCount)
	fmt.Printf("📈 Total functions: %d\n", len(functions))
}