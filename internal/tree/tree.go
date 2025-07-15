package tree

import (
	"fmt"
	"strings"
)

// ColorRatio represents the ratio of different leaf colors
type ColorRatio struct {
	Green  float64
	Yellow float64
	Red    float64
	Brown  float64
}

// Generate creates an SVG tree with specified color ratios
func Generate(green, yellow, red, brown float64) *strings.Builder {

	// Validate and normalize ratios
	total := green + yellow + red + brown
	if total <= 0 {
		total = 1.0
		green, yellow, red, brown = 0.4, 0.3, 0.2, 0.1 // Default values
	}

	colorRatio := ColorRatio{
		Green:  green / total,
		Yellow: yellow / total,
		Red:    red / total,
		Brown:  brown / total,
	}

	return generateTreeSVG(500, 400, colorRatio)
}

func generateTreeSVG(width, height int, colorRatio ColorRatio) *strings.Builder {
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
	svg.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#e1f5fe"/>`, width, height))

	// Ground with depth and dimension
	addGroundAndGrass(&svg, width, height)

	// Trunk
	trunkCenterX := float64(width) / 2
	trunkBottomY := float64(height - 30)
	trunkTopY := float64(height - 150)
	trunkWidth := 40.0

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
	foliageCenterY := trunkTopY - 30 // Move foliage down to cover trunk top
	foliageRadius := 120.0

	// Add individual leaves to fill the entire foliage area
	addFoliage(&svg, foliageCenterX, foliageCenterY, foliageRadius, colorRatio)

	svg.WriteString(`</svg>`)
	return &svg
}