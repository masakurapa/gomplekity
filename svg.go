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

// drawNaturalTree draws a natural tree structure with hierarchical approach
func (svg *SVGGenerator) drawNaturalTree(builder *strings.Builder) {
	// Tree base parameters
	rootX := float64(svg.Width) / 2
	rootY := float64(svg.Height) - 50 // Start from bottom
	
	// Main trunk parameters
	trunkHeight := 200.0
	trunkWidth := 20.0
	
	// Draw the main trunk
	svg.drawMainTrunk(builder, rootX, rootY, trunkWidth, trunkHeight)
	
	// Step 1: Create branch-and-leaf combinations for each Node
	svg.drawHierarchicalTree(builder, rootX, rootY-trunkHeight)
}

// drawHierarchicalTree draws a single trunk with all top directory functions as leaves
func (svg *SVGGenerator) drawHierarchicalTree(builder *strings.Builder, trunkX, trunkTopY float64) {
	if svg.tree.Root == nil || len(svg.tree.Root.Children) == 0 {
		return
	}
	
	// Collect all functions from all files in top directory
	var allFunctions []*TreeNode
	for _, fileNode := range svg.tree.Root.Children {
		allFunctions = append(allFunctions, fileNode.Children...)
	}
	
	// Draw all functions as leaves directly on the single trunk
	svg.drawLeavesCluster(builder, trunkX, trunkTopY, allFunctions)
}

// drawFileBranch draws a branch for a file with its functions as leaves at the tip
func (svg *SVGGenerator) drawFileBranch(builder *strings.Builder, trunkX, trunkTopY float64, index, totalFiles int, fileNode *TreeNode) {
	// Calculate branch position and angle
	heightOffset := 60.0 * float64(index) / float64(totalFiles) // Distribute along trunk
	branchStartY := trunkTopY + heightOffset
	
	// Alternate branches left and right
	var angle float64
	var branchLength float64 = 100.0
	
	if index%2 == 0 {
		// Left side branch
		angle = 135.0 // 135 degrees (upward left)
	} else {
		// Right side branch  
		angle = 45.0 // 45 degrees (upward right)
	}
	
	// Calculate branch end position
	angleRad := angle * math.Pi / 180
	branchEndX := trunkX + branchLength*math.Cos(angleRad)
	branchEndY := branchStartY - branchLength*math.Sin(angleRad)
	
	// Draw the file branch
	builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="8.0"/>`,
		trunkX, branchStartY, branchEndX, branchEndY))
	builder.WriteString("\n")
	
	// Draw all functions as leaves clustered at the branch tip
	functions := fileNode.Children
	svg.drawLeavesCluster(builder, branchEndX, branchEndY, functions)
}

// drawTopDirectoryTree draws a tree with only top directory files
func (svg *SVGGenerator) drawTopDirectoryTree(builder *strings.Builder, trunkX, trunkTopY float64, mainPackage *TreeNode) {
	functions := mainPackage.Children
	functionCount := len(functions)
	
	if functionCount == 0 {
		return
	}
	
	// Create a simple branching structure from the trunk top
	// Each function gets its own small branch from the trunk tip area
	
	branchStartY := trunkTopY
	branchLength := 40.0
	
	// Calculate positions around the trunk tip in a natural distribution
	for i, functionNode := range functions {
		// Create branches radiating from trunk tip in a circular pattern
		angle := float64(i) * 360.0 / float64(functionCount)
		
		// Convert to radians
		angleRad := angle * math.Pi / 180.0
		
		// Calculate branch end position
		branchEndX := trunkX + branchLength*math.Cos(angleRad)
		branchEndY := branchStartY + branchLength*math.Sin(angleRad)
		
		// Draw the small branch for this function
		builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="3.0"/>`,
			trunkX, branchStartY, branchEndX, branchEndY))
		builder.WriteString("\n")
		
		// Draw the leaf at the end of the branch
		svg.drawNaturalLeaf(builder, branchEndX, branchEndY, functionNode)
	}
}

// drawMainPackageBranch draws the main package as a short branch from trunk tip with leaves clustered at end
func (svg *SVGGenerator) drawMainPackageBranch(builder *strings.Builder, trunkX, trunkTopY float64, packageNode *TreeNode) {
	// Create a short upward branch from trunk tip
	branchLength := 50.0
	branchEndX := trunkX
	branchEndY := trunkTopY - branchLength
	
	// Draw the main package branch (short vertical extension)
	builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="6.0"/>`,
		trunkX, trunkTopY, branchEndX, branchEndY))
	builder.WriteString("\n")
	
	// Draw all functions as leaves clustered at the branch tip
	functions := packageNode.Children
	svg.drawLeavesCluster(builder, branchEndX, branchEndY, functions)
}

// drawSubPackageBranch draws a subdirectory package as a branch from trunk with leaves at tip
func (svg *SVGGenerator) drawSubPackageBranch(builder *strings.Builder, trunkX, trunkTopY float64, index, totalPackages int, packageNode *TreeNode) {
	// Calculate branch position and angle
	heightOffset := 60.0 * float64(index) / float64(totalPackages) // Distribute along trunk
	branchStartY := trunkTopY + heightOffset
	
	// Alternate branches left and right
	var angle float64
	var branchLength float64 = 80.0
	
	if index%2 == 0 {
		// Left side branch
		angle = 135.0 // 135 degrees (upward left)
	} else {
		// Right side branch  
		angle = 45.0 // 45 degrees (upward right)
	}
	
	// Calculate branch end position
	angleRad := angle * math.Pi / 180
	branchEndX := trunkX + branchLength*math.Cos(angleRad)
	branchEndY := branchStartY - branchLength*math.Sin(angleRad)
	
	// Draw the package branch
	builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="8.0"/>`,
		trunkX, branchStartY, branchEndX, branchEndY))
	builder.WriteString("\n")
	
	// Draw all functions as leaves clustered at the branch tip
	functions := packageNode.Children
	svg.drawLeavesCluster(builder, branchEndX, branchEndY, functions)
}

// drawLeavesCluster draws a cluster of function leaves naturally spreading from the branch tip
func (svg *SVGGenerator) drawLeavesCluster(builder *strings.Builder, tipX, tipY float64, functions []*TreeNode) {
	functionCount := len(functions)
	if functionCount == 0 {
		return
	}
	
	// Create a natural spreading pattern from the branch tip
	// Use different spreading patterns based on number of leaves
	
	if functionCount == 1 {
		// Single leaf directly at tip
		svg.drawNaturalLeaf(builder, tipX, tipY, functions[0])
	} else if functionCount <= 3 {
		// Small cluster - simple spread
		for i, functionNode := range functions {
			angle := float64(i) * 120.0 // 120 degrees apart
			radius := 15.0 // Reduced from 20.0
			
			angleRad := angle * math.Pi / 180
			leafX := tipX + radius*math.Cos(angleRad)
			leafY := tipY + radius*math.Sin(angleRad)
			
			svg.drawNaturalLeaf(builder, leafX, leafY, functionNode)
		}
	} else {
		// Larger cluster - natural spreading in organic pattern with tighter packing
		for i, functionNode := range functions {
			// Create organic spreading with varied distances - much tighter
			baseRadius := 8.0 // Reduced from 18.0 for tighter clustering
			radiusVariation := float64(i%6) * 6.0 // Vary radius naturally, more layers
			radius := baseRadius + radiusVariation
			
			// Distribute around the tip with some randomness
			angle := float64(i) * 360.0 / float64(functionCount)
			angleVariation := (float64(i%5) - 2) * 10.0 // Reduced variation for more uniform distribution
			angle += angleVariation
			
			angleRad := angle * math.Pi / 180
			leafX := tipX + radius*math.Cos(angleRad)
			leafY := tipY + radius*math.Sin(angleRad)
			
			svg.drawNaturalLeaf(builder, leafX, leafY, functionNode)
		}
	}
}

// BranchInfo represents information about a branch in the new unified system
type BranchInfo struct {
	StartX, StartY float64
	EndX, EndY     float64
	Angle          float64
	Length         float64
	Thickness      float64
	Node           *TreeNode    // The node this branch represents
	SubBranches    []BranchInfo // Child branches
	IsLeafBranch   bool         // True if this branch should have leaves
}

// createSubdirectoryBranches creates branches for subdirectories only (excluding main package)
func (svg *SVGGenerator) createSubdirectoryBranches(trunkX, trunkTopY float64) []BranchInfo {
	var branches []BranchInfo
	
	if svg.tree.Root == nil || len(svg.tree.Root.Children) == 0 {
		return branches
	}
	
	// Filter out the main package - only create branches for subdirectories
	var subPackages []*TreeNode
	for _, packageNode := range svg.tree.Root.Children {
		if packageNode.Name != "main" || packageNode.NodeType != "package" {
			subPackages = append(subPackages, packageNode)
		}
	}
	
	packageCount := len(subPackages)
	if packageCount == 0 {
		return branches
	}
	
	// Distribute branches along the trunk from top to bottom
	for i, packageNode := range subPackages {
		// Calculate branch position along the trunk
		heightRatio := 1.0 - float64(i)/float64(packageCount) // Start from top
		branchY := trunkTopY + 60.0*heightRatio // Space for branch distribution
		
		// Alternate branches left and right from the trunk
		var angle float64
		if i%2 == 0 {
			// Left side branches: upward angles from 135° to 165°
			angle = 135.0 + float64(i/2)*5.0
		} else {
			// Right side branches: upward angles from 45° to 15°
			angle = 45.0 - float64(i/2)*5.0
		}
		
		// Create branch for subdirectory package
		branch := svg.createPackageBranch(trunkX, branchY, angle, packageNode)
		branches = append(branches, branch)
	}
	
	return branches
}

// createBranchSystem creates a unified branch system for the entire tree (deprecated)
func (svg *SVGGenerator) createBranchSystem(trunkX, trunkTopY float64) []BranchInfo {
	// This function is deprecated - use createSubdirectoryBranches for new trunk/branch design
	return svg.createSubdirectoryBranches(trunkX, trunkTopY)
}

// createPackageBranch creates a branch for a package with functions clustered at tip
func (svg *SVGGenerator) createPackageBranch(startX, startY, angle float64, packageNode *TreeNode) BranchInfo {
	// Main branch length based on complexity
	baseLength := 80.0 + float64(len(packageNode.Children))*8.0
	if baseLength > 120.0 {
		baseLength = 120.0
	}
	
	// Calculate end position
	angleRad := angle * math.Pi / 180
	endX := startX + baseLength*math.Cos(angleRad)
	endY := startY - baseLength*math.Sin(angleRad)
	
	// Create main branch with all functions as direct leaves at tip
	mainBranch := BranchInfo{
		StartX:      startX,
		StartY:      startY,
		EndX:        endX,
		EndY:        endY,
		Angle:       angle,
		Length:      baseLength,
		Thickness:   8.0,
		Node:        packageNode,
		SubBranches: []BranchInfo{},
		IsLeafBranch: true, // This branch will have clustered leaves at tip
	}
	
	return mainBranch
}

// createFunctionBranch creates a leaf-bearing branch for a function
func (svg *SVGGenerator) createFunctionBranch(startX, startY, angle float64, functionNode *TreeNode) BranchInfo {
	// Function branch length (shorter than main branches)
	length := 30.0 + float64(functionNode.Complexity)*2.0
	if length > 50.0 {
		length = 50.0
	}
	
	// Calculate end position
	angleRad := angle * math.Pi / 180
	endX := startX + length*math.Cos(angleRad)
	endY := startY - length*math.Sin(angleRad)
	
	return BranchInfo{
		StartX:      startX,
		StartY:      startY,
		EndX:        endX,
		EndY:        endY,
		Angle:       angle,
		Length:      length,
		Thickness:   3.0,
		Node:        functionNode,
		SubBranches: []BranchInfo{},
		IsLeafBranch: true,
	}
}

// drawBranchSystem draws the entire unified branch system
func (svg *SVGGenerator) drawBranchSystem(builder *strings.Builder, branches []BranchInfo) {
	for _, branch := range branches {
		svg.drawBranchRecursive(builder, branch)
	}
}

// drawBranchRecursive recursively draws a branch and its sub-branches
func (svg *SVGGenerator) drawBranchRecursive(builder *strings.Builder, branch BranchInfo) {
	// Draw the branch line
	builder.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="%.1f"/>`,
		branch.StartX, branch.StartY, branch.EndX, branch.EndY, branch.Thickness))
	builder.WriteString("\n")
	
	// If this is a leaf branch, draw leaves at the tip
	if branch.IsLeafBranch {
		svg.drawLeavesAtTip(builder, branch)
	}
	
	// Draw sub-branches
	for _, subBranch := range branch.SubBranches {
		svg.drawBranchRecursive(builder, subBranch)
	}
}

// drawLeavesAtTip draws leaves concentrated at the branch tip
func (svg *SVGGenerator) drawLeavesAtTip(builder *strings.Builder, branch BranchInfo) {
	// For package branches, draw all functions as leaves clustered at tip
	if branch.IsLeafBranch && branch.Node.NodeType == "package" {
		functions := branch.Node.Children
		functionCount := len(functions)
		
		if functionCount == 0 {
			return
		}
		
		// Create circular distribution of leaves around the branch tip
		for i, functionNode := range functions {
			// Calculate position in circle around tip
			radius := 15.0 + float64(i%3)*8.0 // Vary radius for natural look
			angle := float64(i) * 2 * math.Pi / float64(functionCount)
			
			// Add some randomness for natural distribution
			angleOffset := (float64(i%5) - 2) * 0.3
			angle += angleOffset
			
			leafX := branch.EndX + radius*math.Cos(angle)
			leafY := branch.EndY + radius*math.Sin(angle)
			
			// Draw the leaf for this function
			svg.drawNaturalLeaf(builder, leafX, leafY, functionNode)
		}
	}
}

// calculateBranchPositions is deprecated - use createBranchSystem instead
func (svg *SVGGenerator) calculateBranchPositions(trunkX, trunkBaseY, trunkHeight float64) []BranchInfo {
	// This function is no longer used in the unified branch system
	// All branch creation is now handled by createBranchSystem
	return []BranchInfo{}
}

// drawTrunkLeaves draws leaves for top directory functions clustered at trunk tip
func (svg *SVGGenerator) drawTrunkLeaves(builder *strings.Builder, trunkX, trunkBaseY, trunkHeight float64) {
	if svg.tree.Root == nil {
		return
	}
	
	// Find the main package (top directory functions)
	var mainPackage *TreeNode
	for _, packageNode := range svg.tree.Root.Children {
		if packageNode.Name == "main" && packageNode.NodeType == "package" {
			mainPackage = packageNode
			break
		}
	}
	
	if mainPackage == nil || len(mainPackage.Children) == 0 {
		return
	}
	
	// Calculate trunk tip position - this is the actual top of the trunk
	trunkTipX := trunkX
	trunkTipY := trunkBaseY - trunkHeight
	
	// Draw leaves clustered tightly at the very tip of the trunk
	functions := mainPackage.Children
	functionCount := len(functions)
	
	// Use smaller radius and position leaves closer to the tip
	baseRadius := 20.0
	
	for i, functionNode := range functions {
		// Create concentric circles for clustering at tip
		layer := i / 6 // 6 leaves per layer
		radius := baseRadius + float64(layer)*12.0
		
		// Distribute leaves in each layer
		leavesInLayer := 6
		if layer == 0 && functionCount < 6 {
			leavesInLayer = functionCount
		}
		layerIndex := i % leavesInLayer
		
		angle := float64(layerIndex) * 2 * math.Pi / float64(leavesInLayer)
		
		// Add slight variation for natural look
		angleOffset := (float64(i%3) - 1) * 0.2
		angle += angleOffset
		
		// Position leaves just above the trunk tip
		leafX := trunkTipX + radius*math.Cos(angle)
		leafY := trunkTipY - 10.0 + radius*math.Sin(angle) // Move up 10 pixels from tip
		
		// Draw the leaf for this function
		svg.drawNaturalLeaf(builder, leafX, leafY, functionNode)
	}
}

// drawMainTrunk draws the main vertical trunk representing the top directory
func (svg *SVGGenerator) drawMainTrunk(builder *strings.Builder, baseX, baseY, width, height float64) {
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

// drawNaturalBranch and drawLeavesOnBranch are deprecated in unified branch system
// These functions are replaced by drawBranchRecursive and drawLeavesAtTip

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