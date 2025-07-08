package main

import (
	"fmt"
	"math"
	"strings"
)

// SVGGenerator generates SVG visualizations of complexity trees
type SVGGenerator struct {
	Width  int
	Height int
	tree   *ComplexityTree
	Style  string // "hierarchical" or "natural"
}

// NewSVGGenerator creates a new SVG generator
func NewSVGGenerator(tree *ComplexityTree) *SVGGenerator {
	return &SVGGenerator{
		Width:  800,
		Height: 600,
		tree:   tree,
		Style:  "natural", // Default to natural tree style
	}
}

// GenerateSVG generates SVG content for the complexity tree
func (svg *SVGGenerator) GenerateSVG() string {
	var builder strings.Builder
	
	// SVG header
	builder.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`, svg.Width, svg.Height))
	builder.WriteString("\n")
	
	// Background
	builder.WriteString(`<rect width="100%" height="100%" fill="#f8f9fa"/>`)
	builder.WriteString("\n")
	
	// Title
	builder.WriteString(`<text x="400" y="30" text-anchor="middle" font-family="Arial" font-size="20" font-weight="bold" fill="#333">`)
	builder.WriteString("🌳 Complexity Tree Visualization")
	builder.WriteString("</text>\n")
	
	// Draw tree based on style
	if svg.tree.Root != nil {
		if svg.Style == "natural" {
			svg.drawNaturalTree(&builder)
		} else {
			svg.drawTree(&builder)
		}
	}
	
	// SVG footer
	builder.WriteString("</svg>\n")
	
	return builder.String()
}

// drawTree draws the entire tree structure
func (svg *SVGGenerator) drawTree(builder *strings.Builder) {
	rootX := float64(svg.Width) / 2
	rootY := 80.0
	
	// Draw root node
	svg.drawNode(builder, svg.tree.Root, rootX, rootY, 0, 0)
	
	// Draw package nodes
	packageCount := len(svg.tree.Root.Children)
	if packageCount > 0 {
		packageSpacing := float64(svg.Width-100) / float64(packageCount)
		
		for i, packageNode := range svg.tree.Root.Children {
			packageX := 50 + float64(i)*packageSpacing + packageSpacing/2
			packageY := rootY + 120
			
			// Draw connection from root to package
			svg.drawConnection(builder, rootX, rootY+30, packageX, packageY-30)
			
			// Draw package node
			svg.drawNode(builder, packageNode, packageX, packageY, 1, i)
			
			// Draw function nodes
			functionCount := len(packageNode.Children)
			if functionCount > 0 {
				functionSpacing := math.Min(packageSpacing/float64(functionCount), 80)
				startX := packageX - (float64(functionCount-1)*functionSpacing)/2
				
				for j, functionNode := range packageNode.Children {
					functionX := startX + float64(j)*functionSpacing
					functionY := packageY + 100
					
					// Draw connection from package to function
					svg.drawConnection(builder, packageX, packageY+30, functionX, functionY-20)
					
					// Draw function node
					svg.drawNode(builder, functionNode, functionX, functionY, 2, j)
				}
			}
		}
	}
}

// drawNode draws a single node
func (svg *SVGGenerator) drawNode(builder *strings.Builder, node *TreeNode, x, y float64, level, index int) {
	// Get color based on complexity level
	fillColor := svg.getNodeColor(node.Level)
	strokeColor := svg.getStrokeColor(node.Level)
	
	// Draw different shapes based on node type
	switch node.NodeType {
	case "root":
		// Draw trunk as rectangle
		svg.drawTrunk(builder, x, y, fillColor, strokeColor)
	case "package":
		// Draw branch as rounded rectangle
		svg.drawBranch(builder, x, y, fillColor, strokeColor)
	case "function":
		// Draw leaf shape
		svg.drawLeaf(builder, x, y, fillColor, strokeColor, node.Level)
	}
	
	// Draw node label
	var fontSize int
	switch node.NodeType {
	case "root":
		fontSize = 12
	case "package":
		fontSize = 10
	case "function":
		fontSize = 8
	}
	
	// Truncate long names
	displayName := node.Name
	if len(displayName) > 15 {
		displayName = displayName[:12] + "..."
	}
	
	builder.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" font-family="Arial" font-size="%d" fill="#333">`,
		x, y+float64(fontSize)/3, fontSize))
	builder.WriteString(displayName)
	builder.WriteString("</text>\n")
	
	// Draw complexity value for non-root nodes
	if node.NodeType != "root" {
		var offset float64
		if node.NodeType == "function" {
			offset = 25 // More space for leaf shape
		} else {
			offset = 20
		}
		builder.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" font-family="Arial" font-size="10" fill="#666">`,
			x, y+offset))
		builder.WriteString(fmt.Sprintf("%d", node.Complexity))
		builder.WriteString("</text>\n")
	}
}

// drawNaturalTree draws a natural tree structure from bottom to top
func (svg *SVGGenerator) drawNaturalTree(builder *strings.Builder) {
	// Tree base parameters
	trunkBaseX := float64(svg.Width) / 2
	trunkBaseY := float64(svg.Height) - 50 // Start from bottom
	trunkHeight := 200.0
	trunkWidth := 60.0
	
	// Draw main trunk
	svg.drawNaturalTrunk(builder, trunkBaseX, trunkBaseY, trunkWidth, trunkHeight)
	
	// Calculate branch positions
	branches := svg.calculateBranchPositions(trunkBaseX, trunkBaseY, trunkHeight)
	
	// Draw branches and leaves
	for i, branch := range branches {
		svg.drawNaturalBranch(builder, branch, i)
	}
}

// BranchInfo represents information about a branch
type BranchInfo struct {
	StartX, StartY float64
	EndX, EndY     float64
	Angle          float64
	PackageNode    *TreeNode
	Functions      []*TreeNode
}

// calculateBranchPositions calculates positions for branches based on packages
func (svg *SVGGenerator) calculateBranchPositions(trunkX, trunkBaseY, trunkHeight float64) []BranchInfo {
	var branches []BranchInfo
	
	if svg.tree.Root == nil || len(svg.tree.Root.Children) == 0 {
		return branches
	}
	
	packages := svg.tree.Root.Children
	packageCount := len(packages)
	
	// Calculate branch positions along the trunk
	for i, packageNode := range packages {
		// Branch height on trunk (distribute more evenly)
		branchHeight := trunkBaseY - (trunkHeight * (float64(i+1) / float64(packageCount+1)))
		
		// Branch angle (alternate left and right, upward growth)
		var angle float64
		var branchSide string
		
		if i%2 == 0 {
			// Left side branches: upward angles from 160° to 130° (top-left quadrant)
			baseAngle := 160.0
			variation := float64(i/2) * 6.0
			angle = baseAngle - variation
			branchSide = "left"
		} else {
			// Right side branches: upward angles from 20° to 50° (top-right quadrant)
			baseAngle := 20.0
			variation := float64(i/2) * 6.0
			angle = baseAngle + variation
			branchSide = "right"
		}
		
		// Branch length (limited to reasonable size)
		baseBranchLength := 60.0 + float64(len(packageNode.Children))*8
		// Add some variation to make it more natural
		var branchLength float64
		if branchSide == "left" {
			branchLength = baseBranchLength + float64(i%3)*10
		} else {
			branchLength = baseBranchLength + float64((i+1)%3)*10
		}
		
		// Limit branch length to prevent going off screen
		if branchLength > 120.0 {
			branchLength = 120.0
		}
		
		// Calculate branch end position
		angleRad := angle * math.Pi / 180
		endX := trunkX + branchLength*math.Cos(angleRad)
		// For upward branches, we want negative Y (SVG coordinates)
		endY := branchHeight - branchLength*math.Sin(angleRad)
		
		branch := BranchInfo{
			StartX:      trunkX,
			StartY:      branchHeight,
			EndX:        endX,
			EndY:        endY,
			Angle:       angle,
			PackageNode: packageNode,
			Functions:   packageNode.Children,
		}
		
		branches = append(branches, branch)
	}
	
	return branches
}

// drawNaturalTrunk draws the main trunk of the tree
func (svg *SVGGenerator) drawNaturalTrunk(builder *strings.Builder, baseX, baseY, width, height float64) {
	// Draw trunk as a tapered rectangle (wider at bottom)
	topWidth := width * 0.6
	
	// Create path for tapered trunk
	path := fmt.Sprintf(`M %.1f,%.1f L %.1f,%.1f L %.1f,%.1f L %.1f,%.1f Z`,
		baseX-width/2, baseY,           // Bottom left
		baseX+width/2, baseY,           // Bottom right
		baseX+topWidth/2, baseY-height, // Top right
		baseX-topWidth/2, baseY-height) // Top left
	
	// Draw main trunk
	builder.WriteString(fmt.Sprintf(`<path d="%s" fill="#8B4513" stroke="#654321" stroke-width="2"/>`, path))
	builder.WriteString("\n")
	
	// Add realistic bark texture and wood grain
	svg.drawTrunkTexture(builder, baseX, baseY, width, height, topWidth)
}

// drawTrunkTexture adds realistic bark texture and wood grain to the trunk
func (svg *SVGGenerator) drawTrunkTexture(builder *strings.Builder, baseX, baseY, width, height, topWidth float64) {
	// Add vertical wood grain lines
	grainCount := 5
	for i := 0; i < grainCount; i++ {
		// Calculate position across trunk width
		t := float64(i) / float64(grainCount-1)
		
		// Bottom position
		bottomX := baseX - width/2 + width*t
		// Top position (accounting for taper)
		topX := baseX - topWidth/2 + topWidth*t
		
		// Add slight curve to make it more natural
		midX := bottomX + (topX-bottomX)*0.5 + float64(i%2)*8 - 4
		
		// Draw curved grain line
		grainPath := fmt.Sprintf(`M %.1f,%.1f Q %.1f,%.1f %.1f,%.1f`,
			bottomX, baseY,
			midX, baseY-height/2,
			topX, baseY-height)
		
		builder.WriteString(fmt.Sprintf(`<path d="%s" fill="none" stroke="#654321" stroke-width="1" opacity="0.6"/>`, grainPath))
		builder.WriteString("\n")
	}
	
	// Add horizontal bark texture lines
	barkLines := 8
	for i := 1; i < barkLines; i++ {
		y := baseY - (height * float64(i) / float64(barkLines))
		
		// Calculate width at this height
		currentWidth := width - (width-topWidth)*(float64(i)/float64(barkLines))
		
		// Add slight irregularity
		leftX := baseX - currentWidth/2 + float64((i%3-1))*2
		rightX := baseX + currentWidth/2 - float64((i%3-1))*2
		
		builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#654321" stroke-width="0.8" opacity="0.4"/>`,
			leftX, y, rightX, y))
		builder.WriteString("\n")
	}
	
	// Add bark bumps and knots
	svg.addBarkDetails(builder, baseX, baseY, width, height, topWidth)
}

// addBarkDetails adds small bark details like knots and bumps
func (svg *SVGGenerator) addBarkDetails(builder *strings.Builder, baseX, baseY, width, height, topWidth float64) {
	// Add a few knots
	knots := []struct{ x, y, size float64 }{
		{baseX - 8, baseY - height*0.3, 4},
		{baseX + 12, baseY - height*0.6, 3},
		{baseX - 6, baseY - height*0.8, 2},
	}
	
	for _, knot := range knots {
		// Draw knot as small ellipse
		builder.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#654321" opacity="0.8"/>`,
			knot.x, knot.y, knot.size, knot.size*0.7))
		builder.WriteString("\n")
		
		// Add highlight on knot
		builder.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#A0522D" opacity="0.6"/>`,
			knot.x-0.5, knot.y-0.5, knot.size*0.6, knot.size*0.4))
		builder.WriteString("\n")
	}
	
	// Add some bark ridges
	ridges := []struct{ x, y, width, height float64 }{
		{baseX - 15, baseY - height*0.2, 8, 15},
		{baseX + 10, baseY - height*0.4, 6, 12},
		{baseX - 5, baseY - height*0.7, 4, 8},
	}
	
	for _, ridge := range ridges {
		// Draw ridge as small rounded rectangle
		builder.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="2" ry="2" fill="#654321" opacity="0.5"/>`,
			ridge.x-ridge.width/2, ridge.y-ridge.height/2, ridge.width, ridge.height))
		builder.WriteString("\n")
	}
}

// drawNaturalBranch draws a branch with its leaves
func (svg *SVGGenerator) drawNaturalBranch(builder *strings.Builder, branch BranchInfo, index int) {
	// Draw branch line
	builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="%.1f"/>`,
		branch.StartX, branch.StartY, branch.EndX, branch.EndY, math.Max(8.0-float64(index)*0.5, 3.0)))
	builder.WriteString("\n")
	
	// Draw leaves along the branch
	if len(branch.Functions) > 0 {
		svg.drawLeavesOnBranch(builder, branch)
	}
}

// drawLeavesOnBranch draws leaves along a branch
func (svg *SVGGenerator) drawLeavesOnBranch(builder *strings.Builder, branch BranchInfo) {
	functionCount := len(branch.Functions)
	
	for i, functionNode := range branch.Functions {
		// Calculate leaf position along the branch
		t := float64(i+1) / float64(functionCount+1)
		leafX := branch.StartX + (branch.EndX-branch.StartX)*t
		leafY := branch.StartY + (branch.EndY-branch.StartY)*t
		
		// Add some randomness to leaf position
		offsetX := float64((i%3-1)) * 15 // -15, 0, or 15
		offsetY := float64((i%2)) * 10   // 0 or 10
		
		leafX += offsetX
		leafY += offsetY
		
		// Draw leaf based on complexity
		svg.drawNaturalLeaf(builder, leafX, leafY, functionNode)
	}
}

// drawNaturalLeaf draws a single leaf in natural style
func (svg *SVGGenerator) drawNaturalLeaf(builder *strings.Builder, x, y float64, functionNode *TreeNode) {
	fillColor := svg.getNodeColor(functionNode.Level)
	strokeColor := svg.getStrokeColor(functionNode.Level)
	
	// Draw leaf shape based on complexity level
	switch functionNode.Level {
	case "low":
		svg.drawNaturalHealthyLeaf(builder, x, y, fillColor, strokeColor)
	case "medium":
		svg.drawNaturalCautionLeaf(builder, x, y, fillColor, strokeColor)
	case "high":
		svg.drawNaturalDangerLeaf(builder, x, y, fillColor, strokeColor)
	}
	
	// Add function name as tooltip
	builder.WriteString(fmt.Sprintf(`<title>%s (complexity: %d)</title>`, functionNode.Name, functionNode.Complexity))
	builder.WriteString("\n")
}

// drawNaturalHealthyLeaf draws a healthy leaf
func (svg *SVGGenerator) drawNaturalHealthyLeaf(builder *strings.Builder, x, y float64, fillColor, strokeColor string) {
	// Draw as natural leaf shape
	path := fmt.Sprintf(`M %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Z`,
		x, y-12,        // Top point
		x+8, y-8,       // Right curve control
		x+12, y,        // Right point
		x+8, y+8,       // Bottom right control
		x, y+6,         // Bottom point
		x-8, y+8,       // Bottom left control
		x-12, y,        // Left point
		x-8, y-8,       // Top left control
		x, y-12)        // Back to top
	
	builder.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" stroke="%s" stroke-width="1"/>`, path, fillColor, strokeColor))
	builder.WriteString("\n")
}

// drawNaturalCautionLeaf draws a caution leaf
func (svg *SVGGenerator) drawNaturalCautionLeaf(builder *strings.Builder, x, y float64, fillColor, strokeColor string) {
	// Slightly different shape for caution
	path := fmt.Sprintf(`M %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Z`,
		x, y-10,        // Top point
		x+10, y-6,      // Right curve control
		x+14, y+2,      // Right point
		x+6, y+10,      // Bottom right control
		x, y+8,         // Bottom point
		x-6, y+10,      // Bottom left control
		x-14, y+2,      // Left point
		x-10, y-6,      // Top left control
		x, y-10)        // Back to top
	
	builder.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" stroke="%s" stroke-width="1"/>`, path, fillColor, strokeColor))
	builder.WriteString("\n")
}

// drawNaturalDangerLeaf draws a danger leaf (withered)
func (svg *SVGGenerator) drawNaturalDangerLeaf(builder *strings.Builder, x, y float64, fillColor, strokeColor string) {
	// Jagged, withered shape
	path := fmt.Sprintf(`M %.1f,%.1f L %.1f,%.1f L %.1f,%.1f L %.1f,%.1f L %.1f,%.1f L %.1f,%.1f L %.1f,%.1f L %.1f,%.1f L %.1f,%.1f L %.1f,%.1f Z`,
		x, y-8,         // Top
		x+4, y-6,       // Right top
		x+10, y-2,      // Right
		x+6, y+2,       // Right bottom
		x+8, y+6,       // Bottom right
		x, y+4,         // Bottom
		x-8, y+6,       // Bottom left
		x-6, y+2,       // Left bottom
		x-10, y-2,      // Left
		x-4, y-6)       // Back to top
	
	builder.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" stroke="%s" stroke-width="1"/>`, path, fillColor, strokeColor))
	builder.WriteString("\n")
}

// drawConnection draws a line between two nodes
func (svg *SVGGenerator) drawConnection(builder *strings.Builder, x1, y1, x2, y2 float64) {
	builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#666" stroke-width="1"/>`,
		x1, y1, x2, y2))
	builder.WriteString("\n")
}

// drawTrunk draws a trunk shape for the root node
func (svg *SVGGenerator) drawTrunk(builder *strings.Builder, x, y float64, fillColor, strokeColor string) {
	width := 40.0
	height := 50.0
	
	// Draw trunk as rounded rectangle
	builder.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="5" ry="5" fill="%s" stroke="%s" stroke-width="3"/>`,
		x-width/2, y-height/2, width, height, "#8B4513", "#654321"))
	builder.WriteString("\n")
}

// drawBranch draws a branch shape for package nodes
func (svg *SVGGenerator) drawBranch(builder *strings.Builder, x, y float64, fillColor, strokeColor string) {
	width := 60.0
	height := 20.0
	
	// Draw branch as rounded rectangle
	builder.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" ry="10" fill="%s" stroke="%s" stroke-width="2"/>`,
		x-width/2, y-height/2, width, height, "#8B4513", "#654321"))
	builder.WriteString("\n")
}

// drawLeaf draws a leaf shape for function nodes
func (svg *SVGGenerator) drawLeaf(builder *strings.Builder, x, y float64, fillColor, strokeColor string, level string) {
	// Draw leaf as ellipse or different shapes based on complexity
	switch level {
	case "low":
		// Healthy green leaf - simple ellipse
		svg.drawHealthyLeaf(builder, x, y, fillColor, strokeColor)
	case "medium":
		// Caution yellow leaf - slightly different shape
		svg.drawCautionLeaf(builder, x, y, fillColor, strokeColor)
	case "high":
		// Danger red leaf - withered/jagged shape
		svg.drawDangerLeaf(builder, x, y, fillColor, strokeColor)
	}
}

// drawHealthyLeaf draws a healthy green leaf
func (svg *SVGGenerator) drawHealthyLeaf(builder *strings.Builder, x, y float64, fillColor, strokeColor string) {
	// Draw as a simple ellipse
	builder.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="18" ry="12" fill="%s" stroke="%s" stroke-width="1"/>`,
		x, y, fillColor, strokeColor))
	builder.WriteString("\n")
	
	// Add leaf vein
	builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1" opacity="0.6"/>`,
		x-10, y, x+10, y, strokeColor))
	builder.WriteString("\n")
}

// drawCautionLeaf draws a caution yellow leaf
func (svg *SVGGenerator) drawCautionLeaf(builder *strings.Builder, x, y float64, fillColor, strokeColor string) {
	// Draw as a slightly pointed ellipse
	builder.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="16" ry="14" fill="%s" stroke="%s" stroke-width="1"/>`,
		x, y, fillColor, strokeColor))
	builder.WriteString("\n")
	
	// Add leaf vein
	builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1" opacity="0.6"/>`,
		x-8, y, x+8, y, strokeColor))
	builder.WriteString("\n")
	
	// Add small warning dots
	builder.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="2" fill="%s" opacity="0.8"/>`,
		x-5, y-3, strokeColor))
	builder.WriteString("\n")
}

// drawDangerLeaf draws a danger red leaf (withered)
func (svg *SVGGenerator) drawDangerLeaf(builder *strings.Builder, x, y float64, fillColor, strokeColor string) {
	// Draw as a jagged/withered shape using path
	path := fmt.Sprintf(`M %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Q %.1f,%.1f %.1f,%.1f Z`,
		x-15, y, x-8, y-10, x, y-5, x+8, y-12, x+15, y, x+8, y+10, x, y+8, x-8, y+12, x-15, y)
	
	builder.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" stroke="%s" stroke-width="1"/>`,
		path, fillColor, strokeColor))
	builder.WriteString("\n")
	
	// Add cracks/stress lines
	builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1" opacity="0.8"/>`,
		x-5, y-5, x+5, y+5, strokeColor))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1" opacity="0.8"/>`,
		x-5, y+5, x+5, y-5, strokeColor))
	builder.WriteString("\n")
}

// getNodeColor returns the fill color based on complexity level
func (svg *SVGGenerator) getNodeColor(level string) string {
	switch level {
	case "low":
		return "#d4edda"  // Light green
	case "medium":
		return "#fff3cd"  // Light yellow
	case "high":
		return "#f8d7da"  // Light red
	default:
		return "#e9ecef"  // Light gray
	}
}

// getStrokeColor returns the stroke color based on complexity level
func (svg *SVGGenerator) getStrokeColor(level string) string {
	switch level {
	case "low":
		return "#28a745"  // Green
	case "medium":
		return "#ffc107"  // Yellow/Orange
	case "high":
		return "#dc3545"  // Red
	default:
		return "#6c757d"  // Gray
	}
}

// SaveSVG saves the SVG to a file
func (svg *SVGGenerator) SaveSVG(filename string) error {
	content := svg.GenerateSVG()
	return WriteFile(filename, content)
}