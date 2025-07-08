package main

import (
	"flag"
	"fmt"
)

func main() {
	var (
		outputFormat = flag.String("format", "svg", "Output format (svg, html, png)")
		outputFile   = flag.String("output", "", "Output file path")
		targetDir    = flag.String("dir", ".", "Target directory to analyze")
		lowThreshold = flag.Int("low", 5, "Low complexity threshold")
		midThreshold = flag.Int("mid", 10, "Medium complexity threshold")
		help         = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		usage()
		return
	}

	fmt.Printf("Analyzing directory: %s\n", *targetDir)
	fmt.Printf("Output format: %s\n", *outputFormat)
	fmt.Printf("Complexity thresholds: Low=%d, Medium=%d, High=%d+\n", *lowThreshold, *midThreshold, *midThreshold+1)

	if *outputFile != "" {
		fmt.Printf("Output file: %s\n", *outputFile)
	}

	// Create complexity analyzer
	analyzer := NewComplexityAnalyzer(*lowThreshold, *midThreshold)

	// Analyze the directory
	functions, err := analyzer.AnalyzeDirectory(*targetDir)
	if err != nil {
		fmt.Printf("Error analyzing directory: %v\n", err)
		return
	}

	// Print complexity report
	analyzer.PrintComplexityReport(functions)

	// Build and display tree structure
	tree := analyzer.BuildComplexityTree(functions)
	fmt.Printf("\n")
	tree.PrintTree()

	// Generate visualization output
	fmt.Printf("\n🎨 Generating %s visualization...\n", *outputFormat)
	
	switch *outputFormat {
	case "svg":
		outputFileName := *outputFile
		if outputFileName == "" {
			outputFileName = "complexity.svg"
		}
		
		svgGenerator := NewSVGGenerator(tree)
		err := svgGenerator.SaveSVG(outputFileName)
		if err != nil {
			fmt.Printf("Error generating SVG: %v\n", err)
			return
		}
		
		fmt.Printf("✅ SVG visualization saved to: %s\n", outputFileName)
		
	case "html":
		fmt.Println("🔧 HTML output coming soon...")
		
	case "png":
		fmt.Println("🔧 PNG output coming soon...")
		
	default:
		fmt.Printf("❌ Unsupported output format: %s\n", *outputFormat)
	}
}

func usage() {
	fmt.Println("Gomplekity - Go Complexity Tree Visualizer")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  gomplekity [OPTIONS]")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	fmt.Println("  -format string")
	fmt.Println("        Output format (svg, html, png) (default \"svg\")")
	fmt.Println("  -output string")
	fmt.Println("        Output file path")
	fmt.Println("  -dir string")
	fmt.Println("        Target directory to analyze (default \".\")")
	fmt.Println("  -low int")
	fmt.Println("        Low complexity threshold (default 5)")
	fmt.Println("  -mid int")
	fmt.Println("        Medium complexity threshold (default 10)")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  gomplekity")
	fmt.Println("  gomplekity -dir ./src -format html -output complexity.html")
	fmt.Println("  gomplekity -low 3 -mid 7 -format svg")
}
