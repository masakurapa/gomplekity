package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/masakurapa/gomplekity/internal/complexity"
	"github.com/masakurapa/gomplekity/internal/tree"
)

func main() {
	var (
		outputFile        = flag.String("output", "", "Output file path")
		targetDir         = flag.String("dir", ".", "Target directory to analyze")
		mediumThreshold   = flag.Int("medium", 10, "Medium complexity starts from this value (10+)")
		highThreshold     = flag.Int("high", 15, "High complexity starts from this value (15+)")
		criticalThreshold = flag.Int("critical", 20, "Critical complexity starts from this value (20+)")
		verbose           = flag.Bool("verbose", false, "Show detailed complexity analysis")
		help              = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		usage()
		return
	}

	if *verbose {
		fmt.Printf("Analyzing directory: %s\n", *targetDir)
		fmt.Printf("Complexity thresholds: Low<=%d, Medium≥%d, High≥%d, Critical≥%d\n", *mediumThreshold-1, *mediumThreshold, *highThreshold, *criticalThreshold)

		if *outputFile != "" {
			fmt.Printf("Output file: %s\n", *outputFile)
		}
	}

	// Create complexity analyzer
	analyzer := complexity.NewComplexityAnalyzer(*mediumThreshold, *highThreshold, *criticalThreshold)

	// Analyze the directory
	functions, err := analyzer.AnalyzeDirectory(*targetDir)
	if err != nil {
		fmt.Printf("Error analyzing directory: %v\n", err)
		return
	}

	// Print complexity report only if verbose
	if *verbose {
		PrintComplexityReport(functions, analyzer, *mediumThreshold, *highThreshold, *criticalThreshold)

		// Build and display tree structure
		complexityTree := analyzer.BuildComplexityTree(functions)
		fmt.Printf("\n")
		PrintTree(complexityTree)
	}

	// Generate tree visualization based on complexity
	generateTreeVisualization(functions, analyzer, *outputFile)
}

func usage() {
	fmt.Println("Gomplekity - Go Complexity Tree Visualizer")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  gomplekity [OPTIONS]")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	fmt.Println("  -output string")
	fmt.Println("        Output file path")
	fmt.Println("  -dir string")
	fmt.Println("        Target directory to analyze (default \".\")")
	fmt.Println("  -medium int")
	fmt.Println("        Medium complexity starts from this value (10+) (default 10)")
	fmt.Println("  -high int")
	fmt.Println("        High complexity starts from this value (15+) (default 15)")
	fmt.Println("  -critical int")
	fmt.Println("        Critical complexity starts from this value (20+) (default 20)")
	fmt.Println("  -verbose")
	fmt.Println("        Show detailed complexity analysis")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  gomplekity")
	fmt.Println("  gomplekity -dir ./src -output complexity.txt")
	fmt.Println("  gomplekity -medium 8 -high 12 -critical 16 -verbose")
}

// generateTreeVisualization generates a tree SVG based on complexity analysis
func generateTreeVisualization(functions []complexity.FunctionComplexity, analyzer *complexity.ComplexityAnalyzer, outputFile string) {

	// Calculate complexity distribution
	lowCount, mediumCount, highCount, criticalCount := 0, 0, 0, 0

	for _, fn := range functions {
		level := analyzer.GetComplexityLevel(fn.Complexity)
		switch level {
		case "low":
			lowCount++
		case "medium":
			mediumCount++
		case "high":
			highCount++
		case "critical":
			criticalCount++
		}
	}

	// Convert counts to ratios (green=low, yellow=medium, red=high, brown=critical)
	totalFunctions := len(functions)
	if totalFunctions == 0 {
		totalFunctions = 1 // Avoid division by zero
	}

	green := float64(lowCount) / float64(totalFunctions)
	yellow := float64(mediumCount) / float64(totalFunctions)
	red := float64(highCount) / float64(totalFunctions)
	brown := float64(criticalCount) / float64(totalFunctions)

	// Ensure minimum representation for each level if functions exist
	if lowCount > 0 && green < 0.1 {
		green = 0.1
	}
	if mediumCount > 0 && yellow < 0.1 {
		yellow = 0.1
	}
	if highCount > 0 && red < 0.1 {
		red = 0.1
	}
	if criticalCount > 0 && brown < 0.1 {
		brown = 0.1
	}

	// Generate the SVG tree
	svg := tree.Generate(green, yellow, red, brown)

	// Determine output filename
	filename := outputFile
	if filename == "" {
		filename = "complexity_tree.svg"
	}

	// Write to file
	err := os.WriteFile(filename, []byte(svg.String()), 0644)
	if err != nil {
		fmt.Printf("❌ Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("✅ Tree visualization saved to: %s\n", filename)
	fmt.Printf("📊 Color distribution: 🟢%.1f%% 🟡%.1f%% 🔴%.1f%% 🟤%.1f%%\n",
		green*100, yellow*100, red*100, brown*100)
}
