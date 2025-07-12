package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

type ColorRatio struct {
	Green  float64
	Yellow float64
	Red    float64
	Brown  float64
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	// Command line flags
	var (
		green    = flag.Float64("green", 0.4, "Green leaves ratio (0.0-1.0)")
		yellow   = flag.Float64("yellow", 0.3, "Yellow leaves ratio (0.0-1.0)")
		red      = flag.Float64("red", 0.2, "Red leaves ratio (0.0-1.0)")
		brown    = flag.Float64("brown", 0.1, "Brown leaves ratio (0.0-1.0)")
		output   = flag.String("output", "tree.svg", "Output SVG filename")
		help     = flag.Bool("help", false, "Show usage help")
	)
	flag.Parse()
	
	if *help {
		fmt.Println("Tree SVG Generator")
		fmt.Println("Usage:")
		fmt.Println("  go run main.go [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println("  go run main.go -green=0.5 -yellow=0.2 -red=0.2 -brown=0.1 -output=autumn_tree.svg")
		return
	}
	
	// Validate and normalize ratios
	total := *green + *yellow + *red + *brown
	if total <= 0 {
		fmt.Println("Error: Total color ratio must be greater than 0")
		os.Exit(1)
	}
	
	// Normalize ratios to sum to 1.0
	colorRatio := ColorRatio{
		Green:  *green / total,
		Yellow: *yellow / total,
		Red:    *red / total,
		Brown:  *brown / total,
	}
	
	fmt.Printf("Color ratios: Green=%.1f%%, Yellow=%.1f%%, Red=%.1f%%, Brown=%.1f%%\n",
		colorRatio.Green*100, colorRatio.Yellow*100, colorRatio.Red*100, colorRatio.Brown*100)
	
	svg := generateTreeSVG(800, 600, colorRatio)
	
	file, err := os.Create(*output)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	
	file.WriteString(svg)
	fmt.Printf("Tree SVG generated: %s\n", *output)
}

func generateTreeSVG(width, height int, colorRatio ColorRatio) string {
	var svg strings.Builder
	
	svg.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`, width, height))
	svg.WriteString(`<defs>`)
	
	// Trunk gradient
	svg.WriteString(`<linearGradient id="trunkGrad" x1="0%" y1="0%" x2="100%" y2="0%">`)
	svg.WriteString(`<stop offset="0%" style="stop-color:#6d4c41;stop-opacity:1" />`)
	svg.WriteString(`<stop offset="50%" style="stop-color:#8d6e63;stop-opacity:1" />`)
	svg.WriteString(`<stop offset="100%" style="stop-color:#a1887f;stop-opacity:1" />`)
	svg.WriteString(`</linearGradient>`)
	
	
	
	svg.WriteString(`</defs>`)
	
	// Background
	svg.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#87ceeb"/>`, width, height))
	
	// Ground
	svg.WriteString(fmt.Sprintf(`<rect x="0" y="%d" width="%d" height="50" fill="#90ee90"/>`, height-50, width))
	
	
	// Trunk
	trunkCenterX := float64(width) / 2
	trunkBottomY := float64(height - 50)
	trunkTopY := float64(height - 200)
	trunkWidth := 60.0
	
	// Trunk shape (slightly tapered)
	svg.WriteString(fmt.Sprintf(`<path d="M %.1f %.1f Q %.1f %.1f %.1f %.1f L %.1f %.1f Q %.1f %.1f %.1f %.1f Z" fill="url(#trunkGrad)"/>`,
		trunkCenterX-trunkWidth/2, trunkBottomY,
		trunkCenterX, trunkTopY+50, 
		trunkCenterX-trunkWidth/3, trunkTopY,
		trunkCenterX+trunkWidth/3, trunkTopY,
		trunkCenterX, trunkTopY+50,
		trunkCenterX+trunkWidth/2, trunkBottomY))
	
	// Foliage area settings
	foliageCenterX := trunkCenterX
	foliageCenterY := trunkTopY - 80
	foliageRadius := 150.0
	
	
	// Add individual leaves to fill the entire foliage area
	totalLeaves := 700
	
	// Color definitions
	colors := map[string][]string{
		"green":  {"#4caf50", "#66bb6a", "#81c784"},
		"yellow": {"#ffeb3b", "#ffc107", "#ff9800"},
		"red":    {"#f44336", "#e53935", "#d32f2f"},
		"brown":  {"#8d6e63", "#6d4c41", "#5d4037"},
	}
	
	// Calculate number of leaves for each color
	greenLeaves := int(float64(totalLeaves) * colorRatio.Green)
	yellowLeaves := int(float64(totalLeaves) * colorRatio.Yellow)
	redLeaves := int(float64(totalLeaves) * colorRatio.Red)
	brownLeaves := int(float64(totalLeaves) * colorRatio.Brown)
	
	// Generate leaves in multiple layers for density
	for layer := 0; layer < 5; layer++ {
		layerRadius := foliageRadius * (0.3 + float64(layer)*0.14) // Different layers at different radii, start smaller
		
		// Generate green leaves
		for i := 0; i < greenLeaves/5; i++ {
			generateLeafInArea(&svg, foliageCenterX, foliageCenterY, colors["green"], layerRadius)
		}
		
		// Generate yellow leaves
		for i := 0; i < yellowLeaves/5; i++ {
			generateLeafInArea(&svg, foliageCenterX, foliageCenterY, colors["yellow"], layerRadius)
		}
		
		// Generate red leaves
		for i := 0; i < redLeaves/5; i++ {
			generateLeafInArea(&svg, foliageCenterX, foliageCenterY, colors["red"], layerRadius)
		}
		
		// Generate brown leaves
		for i := 0; i < brownLeaves/5; i++ {
			generateLeafInArea(&svg, foliageCenterX, foliageCenterY, colors["brown"], layerRadius)
		}
	}
	
	svg.WriteString(`</svg>`)
	return svg.String()
}

func generateLeafInArea(svg *strings.Builder, centerX, centerY float64, colorSet []string, maxRadius float64) {
	// Random position within the foliage area with better distribution
	angle := rand.Float64() * 2 * math.Pi
	// Use square root to get more even distribution across the circular area
	distance := math.Sqrt(rand.Float64()) * maxRadius * 0.9 // Slightly reduce to 90% to keep within bounds
	x := centerX + distance*math.Cos(angle)
	y := centerY + distance*math.Sin(angle)
	
	// Random size with more variation
	size := 6 + rand.Float64()*18
	
	// Random color from the set
	color := colorSet[rand.Intn(len(colorSet))]
	
	// Random opacity for natural blending
	opacity := 0.6 + rand.Float64()*0.3
	
	// Random rotation for natural variation
	rotation := rand.Float64() * 360
	
	// Generate realistic leaf shape using SVG path
	generateLeafShape(svg, x, y, size, color, opacity, rotation)
}

func generateLeafShape(svg *strings.Builder, x, y, size float64, color string, opacity float64, rotation float64) {
	// Create a realistic leaf shape with stem
	leafWidth := size
	leafHeight := size * 1.4
	
	// Leaf shape path - elongated with pointed tip and indented sides
	svg.WriteString(fmt.Sprintf(`<g transform="translate(%.1f,%.1f) rotate(%.1f)">`, x, y, rotation))
	
	// Main leaf body
	svg.WriteString(fmt.Sprintf(`<path d="M 0 %.1f Q %.1f %.1f %.1f 0 Q %.1f %.1f 0 %.1f Q %.1f %.1f %.1f 0 Q %.1f %.1f 0 %.1f Z" fill="%s" opacity="%.2f"/>`,
		-leafHeight/2,
		leafWidth/3, -leafHeight/3, leafWidth/2,
		leafWidth/3, leafHeight/3, leafHeight/2,
		-leafWidth/3, leafHeight/3, -leafWidth/2,
		-leafWidth/3, -leafHeight/3, -leafHeight/2,
		color, opacity))
	
	// Leaf stem
	stemLength := size * 0.3
	svg.WriteString(fmt.Sprintf(`<line x1="0" y1="%.1f" x2="0" y2="%.1f" stroke="#8d6e63" stroke-width="1" opacity="%.2f"/>`,
		leafHeight/2, leafHeight/2+stemLength, opacity*0.8))
	
	// Central vein
	svg.WriteString(fmt.Sprintf(`<line x1="0" y1="%.1f" x2="0" y2="%.1f" stroke="#2e7d32" stroke-width="0.5" opacity="%.2f"/>`,
		-leafHeight/2, leafHeight/2, opacity*0.6))
	
	svg.WriteString(`</g>`)
}