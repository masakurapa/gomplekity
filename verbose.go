package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/masakurapa/gomplekity/internal/complexity"
)

// PrintComplexityReport prints a formatted complexity report
func PrintComplexityReport(functions []complexity.FunctionComplexity, analyzer *complexity.ComplexityAnalyzer, mediumThreshold, highThreshold, criticalThreshold int) {
	fmt.Printf("🌳 Complexity Analysis Report\n")
	fmt.Printf("================================\n")
	fmt.Printf("Thresholds: Low < %d, Medium ≥ %d, High ≥ %d, Critical ≥ %d\n\n",
		mediumThreshold, mediumThreshold, highThreshold, criticalThreshold)

	// Calculate package statistics
	packages := calculatePackageComplexity(functions)

	fmt.Printf("📦 Package Statistics:\n")
	for packageName, pkg := range packages {
		fmt.Printf("  %s: avg=%.1f, max=%d, min=%d, total=%d (%d functions)\n",
			packageName, pkg.AverageComplexity, pkg.MaxComplexity, pkg.MinComplexity,
			pkg.TotalComplexity, len(pkg.Functions))
	}
	fmt.Printf("\n🔍 Function Details:\n")

	lowCount, mediumCount, highCount, criticalCount := 0, 0, 0, 0

	for _, fn := range functions {
		level := analyzer.GetComplexityLevel(fn.Complexity)

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
		case "critical":
			emoji = "🟤"
			criticalCount++
		}

		fmt.Printf("%s %s (%s): %d - %s:%d\n",
			emoji, fn.Name, level, fn.Complexity, fn.File, fn.Line)
	}

	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("🟢 Low complexity: %d functions\n", lowCount)
	fmt.Printf("🟡 Medium complexity: %d functions\n", mediumCount)
	fmt.Printf("🔴 High complexity: %d functions\n", highCount)
	fmt.Printf("🟤 Critical complexity: %d functions\n", criticalCount)
	fmt.Printf("📈 Total functions: %d\n", len(functions))
}

// PackageComplexity represents the complexity statistics of a package
type PackageComplexity struct {
	PackageName       string
	Functions         []complexity.FunctionComplexity
	TotalComplexity   int
	AverageComplexity float64
	MaxComplexity     int
	MinComplexity     int
}

// calculatePackageComplexity calculates package-level complexity statistics
func calculatePackageComplexity(functions []complexity.FunctionComplexity) map[string]PackageComplexity {
	packageMap := make(map[string][]complexity.FunctionComplexity)

	// Group functions by package (extracted from file path)
	for _, fn := range functions {
		// Extract package name from file path
		packageName := filepath.Dir(fn.File)
		if packageName == "." {
			packageName = "main"
		}

		packageMap[packageName] = append(packageMap[packageName], fn)
	}

	packages := make(map[string]PackageComplexity)

	for packageName, packageFunctions := range packageMap {
		if len(packageFunctions) == 0 {
			continue
		}

		total := 0
		min := packageFunctions[0].Complexity
		max := packageFunctions[0].Complexity

		for _, fn := range packageFunctions {
			total += fn.Complexity
			if fn.Complexity < min {
				min = fn.Complexity
			}
			if fn.Complexity > max {
				max = fn.Complexity
			}
		}

		average := float64(total) / float64(len(packageFunctions))

		packages[packageName] = PackageComplexity{
			PackageName:       packageName,
			Functions:         packageFunctions,
			TotalComplexity:   total,
			AverageComplexity: average,
			MaxComplexity:     max,
			MinComplexity:     min,
		}
	}

	return packages
}

// PrintTree prints the tree structure for debugging
func PrintTree(tree *complexity.ComplexityTree) {
	fmt.Printf("🌳 Complexity Tree Structure\n")
	fmt.Printf("=============================\n")
	printNode(tree.Root, 0)
	fmt.Println()
}

// printNode recursively prints tree nodes with indentation
func printNode(node *complexity.TreeNode, depth int) {
	indent := strings.Repeat("  ", depth)

	var emoji string
	switch node.Level {
	case "low":
		emoji = "🟢"
	case "medium":
		emoji = "🟡"
	case "high":
		emoji = "🔴"
	case "critical":
		emoji = "🟤"
	default:
		emoji = "⚪"
	}

	complexityInfo := ""
	if node.NodeType != "root" {
		complexityInfo = fmt.Sprintf(" (complexity: %d)", node.Complexity)
	}

	fmt.Printf("%s%s %s [%s]%s\n", indent, emoji, node.Name, node.NodeType, complexityInfo)

	for _, child := range node.Children {
		printNode(child, depth+1)
	}
}
