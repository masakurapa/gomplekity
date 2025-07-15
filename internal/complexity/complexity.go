package complexity

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

// AnalyzeTopDirectoryOnly analyzes only Go files in the specified directory (no subdirectories)
func (ca *ComplexityAnalyzer) AnalyzeTopDirectoryOnly(dir string) ([]FunctionComplexity, error) {
	var functions []FunctionComplexity

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Skip non-Go files
		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		// Skip test files
		if strings.HasSuffix(file.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		funcs, err := ca.analyzeFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze file %s: %w", filePath, err)
		}

		functions = append(functions, funcs...)
	}

	return functions, nil
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