package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/masakurapa/gomplekity/internal/complexity"
	"github.com/masakurapa/gomplekity/internal/tree"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
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
		svgOutput         = flag.Bool("svg", false, "Generate SVG output instead of PNG")
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
	generateTreeVisualization(functions, analyzer, *outputFile, *svgOutput)
}

func usage() {
	fmt.Println("Gomplekity - Go Complexity Tree Visualizer")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  gomplekity [OPTIONS]")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	fmt.Println("  -output string")
	fmt.Println("        Output file path (extension determines format: .svg or .png)")
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
	fmt.Println("  -svg")
	fmt.Println("        Generate SVG output instead of PNG (default is PNG)")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  gomplekity")
	fmt.Println("  gomplekity -dir ./src -output complexity.png")
	fmt.Println("  gomplekity -dir ./src -output complexity.svg -svg")
	fmt.Println("  gomplekity -medium 8 -high 12 -critical 16 -verbose")
}

// generateTreeVisualization generates a tree visualization based on complexity analysis
func generateTreeVisualization(functions []complexity.FunctionComplexity, analyzer *complexity.ComplexityAnalyzer, outputFile string, svgOutput bool) {

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

	// Determine output filename and format
	filename := outputFile
	if filename == "" {
		if svgOutput {
			filename = "complexity_tree.svg"
		} else {
			filename = "complexity_tree.png"
		}
	} else {
		// Check if output format matches filename extension
		ext := strings.ToLower(filepath.Ext(filename))
		if ext == ".svg" {
			svgOutput = true
		} else if ext == ".png" {
			svgOutput = false
		}
	}

	if svgOutput {
		// Write SVG to file
		err := os.WriteFile(filename, []byte(svg.String()), 0644)
		if err != nil {
			fmt.Printf("❌ Error writing SVG file: %v\n", err)
			return
		}
	} else {
		// Convert SVG to PNG
		err := convertSVGToPNG(svg.String(), filename)
		if err != nil {
			fmt.Printf("❌ Error writing PNG file: %v\n", err)
			return
		}
	}

	fmt.Printf("✅ Tree visualization saved to: %s\n", filename)
	fmt.Printf("📊 Color distribution: 🟢%.1f%% 🟡%.1f%% 🔴%.1f%% 🟤%.1f%%\n",
		green*100, yellow*100, red*100, brown*100)
}

// convertSVGToPNG converts SVG string to PNG and saves it to file
func convertSVGToPNG(svgContent, filename string) error {
	// Fix gradients in SVG content before parsing
	fixedSVG := fixGradientsInSVG(svgContent)
	
	// Parse SVG content
	icon, err := oksvg.ReadIconStream(strings.NewReader(fixedSVG))
	if err != nil {
		return fmt.Errorf("failed to parse SVG: %v", err)
	}
	
	// Set up rendering dimensions (original size)
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	if w == 0 || h == 0 {
		w, h = 500, 400 // Default size
	}
	
	// Use original scale for proper sizing
	scale := 1.0
	
	// Create image
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	
	// Create scanner and raster
	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scanner)
	
	// Render SVG to image with scaling
	icon.Draw(raster, scale)
	
	// Create PNG file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create PNG file: %v", err)
	}
	defer file.Close()
	
	// Write PNG
	err = png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("failed to write PNG: %v", err)
	}
	
	return nil
}

// fixGradientsInSVG replaces gradient fills with solid colors
func fixGradientsInSVG(svgContent string) string {
	// Replace trunk gradient with solid brown color
	trunkGradientRe := regexp.MustCompile(`url\(#trunkGrad\)`)
	svgContent = trunkGradientRe.ReplaceAllString(svgContent, `#8d6e63`)
	
	// Replace ground gradient with solid green color  
	groundGradientRe := regexp.MustCompile(`url\(#groundDepth\)`)
	svgContent = groundGradientRe.ReplaceAllString(svgContent, `#4caf50`)
	
	// Remove gradient definitions to reduce file size
	gradientDefRe := regexp.MustCompile(`<defs>.*?</defs>`)
	svgContent = gradientDefRe.ReplaceAllString(svgContent, ``)
	
	return svgContent
}
