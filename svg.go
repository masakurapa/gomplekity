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
}

// NewSVGGenerator creates a new SVG generator
func NewSVGGenerator(tree *ComplexityTree) *SVGGenerator {
	return &SVGGenerator{
		Width:  800,
		Height: 600,
		tree:   tree,
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
	
	// Draw tree
	if svg.tree.Root != nil {
		svg.drawTree(&builder)
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