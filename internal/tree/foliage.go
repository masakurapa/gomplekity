package tree

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

func addFoliage(svg *strings.Builder, centerX, centerY, radius float64, colorRatio ColorRatio) {
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
		layerRadius := radius * (0.3 + float64(layer)*0.14) // Different layers at different radii, start smaller
		
		// Generate green leaves
		for i := 0; i < greenLeaves/5; i++ {
			generateLeafInArea(svg, centerX, centerY, colors["green"], layerRadius)
		}
		
		// Generate yellow leaves
		for i := 0; i < yellowLeaves/5; i++ {
			generateLeafInArea(svg, centerX, centerY, colors["yellow"], layerRadius)
		}
		
		// Generate red leaves
		for i := 0; i < redLeaves/5; i++ {
			generateLeafInArea(svg, centerX, centerY, colors["red"], layerRadius)
		}
		
		// Generate brown leaves
		for i := 0; i < brownLeaves/5; i++ {
			generateLeafInArea(svg, centerX, centerY, colors["brown"], layerRadius)
		}
	}
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

func generateFallenLeaf(svg *strings.Builder, x, y, size float64, color string, rotation float64) {
	// Create a fallen leaf on the ground with shadow
	leafWidth := size
	leafHeight := size * 1.2
	
	// Leaf shadow first
	svg.WriteString(fmt.Sprintf(`<g transform="translate(%.1f,%.1f) rotate(%.1f)">`, x+1, y+1, rotation))
	svg.WriteString(fmt.Sprintf(`<path d="M 0 %.1f Q %.1f %.1f %.1f 0 Q %.1f %.1f 0 %.1f Q %.1f %.1f %.1f 0 Q %.1f %.1f 0 %.1f Z" fill="#1b5e20" opacity="0.3"/>`,
		-leafHeight/2,
		leafWidth/3, -leafHeight/3, leafWidth/2,
		leafWidth/3, leafHeight/3, leafHeight/2,
		-leafWidth/3, leafHeight/3, -leafWidth/2,
		-leafWidth/3, -leafHeight/3, -leafHeight/2))
	svg.WriteString(`</g>`)
	
	// Main fallen leaf
	svg.WriteString(fmt.Sprintf(`<g transform="translate(%.1f,%.1f) rotate(%.1f)">`, x, y, rotation))
	svg.WriteString(fmt.Sprintf(`<path d="M 0 %.1f Q %.1f %.1f %.1f 0 Q %.1f %.1f 0 %.1f Q %.1f %.1f %.1f 0 Q %.1f %.1f 0 %.1f Z" fill="%s" opacity="%.2f"/>`,
		-leafHeight/2,
		leafWidth/3, -leafHeight/3, leafWidth/2,
		leafWidth/3, leafHeight/3, leafHeight/2,
		-leafWidth/3, leafHeight/3, -leafWidth/2,
		-leafWidth/3, -leafHeight/3, -leafHeight/2,
		color, 0.6+rand.Float64()*0.3))
	
	// Central vein on fallen leaf
	svg.WriteString(fmt.Sprintf(`<line x1="0" y1="%.1f" x2="0" y2="%.1f" stroke="#2e7d32" stroke-width="0.3" opacity="%.2f"/>`,
		-leafHeight/2, leafHeight/2, 0.4))
	
	svg.WriteString(`</g>`)
}