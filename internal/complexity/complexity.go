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

// TreeNode represents a node in the complexity tree
type TreeNode struct {
	Name       string
	NodeType   string // "root", "package", "function", "file"
	Complexity int
	Level      string // "low", "medium", "high", "critical"
	Color      string // "green", "yellow", "red", "brown"
	Children   []*TreeNode
	Parent     *TreeNode
}

// ComplexityTree represents the entire complexity tree structure
type ComplexityTree struct {
	Root *TreeNode
}

// ComplexityAnalyzer analyzes the cyclomatic complexity of Go files
type ComplexityAnalyzer struct {
	mediumThreshold   int
	highThreshold     int
	criticalThreshold int
}

// NewComplexityAnalyzer creates a new complexity analyzer
func NewComplexityAnalyzer(mediumThreshold, highThreshold, criticalThreshold int) *ComplexityAnalyzer {
	return &ComplexityAnalyzer{
		mediumThreshold:   mediumThreshold,
		highThreshold:     highThreshold,
		criticalThreshold: criticalThreshold,
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
	if complexity < ca.mediumThreshold {
		return "low"
	} else if complexity < ca.highThreshold {
		return "medium"
	} else if complexity < ca.criticalThreshold {
		return "high"
	}
	return "critical"
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
	case "critical":
		return "brown"
	default:
		return "gray"
	}
}

// BuildComplexityTree builds a tree structure from complexity data organized by files
func (ca *ComplexityAnalyzer) BuildComplexityTree(functions []FunctionComplexity) *ComplexityTree {
	// Create root node
	root := &TreeNode{
		Name:     "Project Root",
		NodeType: "root",
		Level:    "low",
		Color:    "green",
		Children: []*TreeNode{},
	}

	// Group functions by file
	fileMap := make(map[string][]FunctionComplexity)
	for _, fn := range functions {
		fileName := filepath.Base(fn.File)
		fileMap[fileName] = append(fileMap[fileName], fn)
	}

	// Create file nodes (branches)
	for fileName, fileFunctions := range fileMap {
		// Calculate file complexity statistics
		totalComplexity := 0
		for _, fn := range fileFunctions {
			totalComplexity += fn.Complexity
		}
		avgComplexity := float64(totalComplexity) / float64(len(fileFunctions))

		fileNode := &TreeNode{
			Name:       fileName,
			NodeType:   "file",
			Complexity: totalComplexity,
			Level:      ca.GetComplexityLevel(int(avgComplexity)),
			Color:      ca.GetComplexityColor(int(avgComplexity)),
			Children:   []*TreeNode{},
			Parent:     root,
		}

		// Create function nodes (leaves) for this file
		for _, fn := range fileFunctions {
			functionNode := &TreeNode{
				Name:       fn.Name,
				NodeType:   "function",
				Complexity: fn.Complexity,
				Level:      ca.GetComplexityLevel(fn.Complexity),
				Color:      ca.GetComplexityColor(fn.Complexity),
				Children:   []*TreeNode{},
				Parent:     fileNode,
			}
			fileNode.Children = append(fileNode.Children, functionNode)
		}

		root.Children = append(root.Children, fileNode)
	}

	return &ComplexityTree{Root: root}
}
