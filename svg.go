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
	// Node appearance based on type and complexity
	var radius float64
	var strokeWidth int
	
	switch node.NodeType {
	case "root":
		radius = 25
		strokeWidth = 3
	case "package":
		radius = 20
		strokeWidth = 2
	case "function":
		radius = 15
		strokeWidth = 1
	}
	
	// Get color based on complexity level
	fillColor := svg.getNodeColor(node.Level)
	strokeColor := svg.getStrokeColor(node.Level)
	
	// Draw circle
	builder.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" stroke="%s" stroke-width="%d"/>`,
		x, y, radius, fillColor, strokeColor, strokeWidth))
	builder.WriteString("\n")
	
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
		builder.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" font-family="Arial" font-size="10" fill="#666">`,
			x, y+radius+15))
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