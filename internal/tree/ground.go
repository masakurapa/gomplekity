package tree

import (
	"fmt"
	"math/rand"
	"strings"
)

func addGroundAndGrass(svg *strings.Builder, width, height int) {
	// Ground with depth and dimension
	// Add gradient base for ground depth
	svg.WriteString(`<defs>`)
	svg.WriteString(`<radialGradient id="groundDepth" cx="50%" cy="30%" r="70%">`)
	svg.WriteString(`<stop offset="0%" style="stop-color:#66bb6a;stop-opacity:1" />`)
	svg.WriteString(`<stop offset="70%" style="stop-color:#4caf50;stop-opacity:1" />`)
	svg.WriteString(`<stop offset="100%" style="stop-color:#2e7d32;stop-opacity:1" />`)
	svg.WriteString(`</radialGradient>`)
	svg.WriteString(`</defs>`)
	
	// Base ground with depth gradient
	svg.WriteString(fmt.Sprintf(`<rect x="0" y="%d" width="%d" height="30" fill="url(#groundDepth)"/>`, height-30, width))
	
	// Add subtle grass texture patterns on the base
	grassPatternColors := []string{"#43a047", "#388e3c", "#2e7d32", "#66bb6a"}
	
	// Create dense grass texture covering 90% of ground
	groundArea := float64(width * 30) // Total ground area
	grassCoverage := groundArea * 0.9 // 90% coverage
	grassPatches := int(grassCoverage / 40) // Each patch covers ~40 pixels
	
	for i := 0; i < grassPatches; i++ {
		patchX := rand.Float64() * float64(width)
		patchY := float64(height-30) + rand.Float64()*30
		
		// Create dense grass tuft
		tufts := 5 + rand.Intn(8) // More grass per patch
		for j := 0; j < tufts; j++ {
			grassX := patchX + (rand.Float64()-0.5)*8
			grassY := patchY + (rand.Float64()-0.5)*6
			grassHeight := 1.5 + rand.Float64()*5
			grassBend := (rand.Float64() - 0.5) * 3
			
			color := grassPatternColors[rand.Intn(len(grassPatternColors))]
			
			svg.WriteString(fmt.Sprintf(`<path d="M %.1f %.1f Q %.1f %.1f %.1f %.1f" stroke="%s" stroke-width="1.2" fill="none" opacity="0.5"/>`,
				grassX, grassY,
				grassX+grassBend*0.5, grassY-grassHeight*0.6,
				grassX+grassBend, grassY-grassHeight,
				color))
		}
	}
	
	// Add even more fine grass details for density
	for i := 0; i < int(float64(grassPatches)*0.5); i++ {
		x := rand.Float64() * float64(width)
		y := float64(height-30) + rand.Float64()*30
		
		// Very short grass for base texture
		microGrassHeight := 1 + rand.Float64()*2
		microGrassBend := (rand.Float64() - 0.5) * 1
		color := grassPatternColors[rand.Intn(len(grassPatternColors))]
		
		svg.WriteString(fmt.Sprintf(`<path d="M %.1f %.1f Q %.1f %.1f %.1f %.1f" stroke="%s" stroke-width="0.8" fill="none" opacity="0.3"/>`,
			x, y,
			x+microGrassBend*0.5, y-microGrassHeight*0.7,
			x+microGrassBend, y-microGrassHeight,
			color))
	}
	
	// Add fallen leaves on ground for natural appearance
	fallenLeafColors := []string{"#4caf50", "#66bb6a", "#43a047", "#ffeb3b", "#ffc107", "#8d6e63"}
	
	for i := 0; i < 15; i++ {
		leafX := rand.Float64() * float64(width)
		leafY := float64(height-25) + rand.Float64()*25
		leafSize := 8 + rand.Float64()*12
		leafRotation := rand.Float64() * 360
		leafColor := fallenLeafColors[rand.Intn(len(fallenLeafColors))]
		
		// Generate fallen leaf with shadow
		generateFallenLeaf(svg, leafX, leafY, leafSize, leafColor, leafRotation)
	}
	
	// Add wind-blown grass on top
	addWindBlownGrass(svg, width, height)
}

func addWindBlownGrass(svg *strings.Builder, width, height int) {
	grassColors := []string{"#2e7d32", "#388e3c", "#43a047", "#66bb6a"}
	
	// Generate wind-blown grass blades
	for i := 0; i < 200; i++ {
		x := rand.Float64() * float64(width)
		grassHeight := 5 + rand.Float64()*15
		
		// Wind direction (slight rightward bend)
		windForce := 0.3 + rand.Float64()*0.4
		bendAmount := grassHeight * windForce
		
		// Create curved grass blade using quadratic curve
		startX := x
		startY := float64(height - 30)
		controlX := x + bendAmount*0.6
		controlY := startY - grassHeight*0.7
		endX := x + bendAmount
		endY := startY - grassHeight
		
		color := grassColors[rand.Intn(len(grassColors))]
		
		svg.WriteString(fmt.Sprintf(`<path d="M %.1f %.1f Q %.1f %.1f %.1f %.1f" stroke="%s" stroke-width="%.1f" fill="none" opacity="%.2f"/>`,
			startX, startY, controlX, controlY, endX, endY,
			color, 0.5+rand.Float64()*0.8, 0.7+rand.Float64()*0.3))
	}
	
	// Add some grass clusters for density
	for i := 0; i < 50; i++ {
		clusterX := rand.Float64() * float64(width)
		clusterSize := 3 + rand.Intn(5)
		
		for j := 0; j < clusterSize; j++ {
			x := clusterX + (rand.Float64()-0.5)*8
			grassHeight := 3 + rand.Float64()*10
			
			windBend := grassHeight * (0.2 + rand.Float64()*0.3)
			
			startX := x
			startY := float64(height - 30)
			endX := x + windBend
			endY := startY - grassHeight
			
			color := grassColors[rand.Intn(len(grassColors))]
			
			svg.WriteString(fmt.Sprintf(`<path d="M %.1f %.1f Q %.1f %.1f %.1f %.1f" stroke="%s" stroke-width="%.1f" fill="none" opacity="%.2f"/>`,
				startX, startY, startX+windBend*0.5, startY-grassHeight*0.6, endX, endY,
				color, 0.6+rand.Float64()*0.6, 0.6+rand.Float64()*0.4))
		}
	}
}