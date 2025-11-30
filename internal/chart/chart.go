package chart

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/odesaur/vex-tui/pkg/models"
)

const (
	ChartHeight = 15
	ChartWidth  = 60
)

type ChartType int

const (
	ChartBar ChartType = iota
	ChartLine
	ChartSparkline
	ChartPie
)

type ChartData struct {
	Labels []string
	Values []float64
	Title  string
}

// ExtractChartData extracts numeric data from selected cells
func ExtractChartData(sheet models.Sheet, startRow, startCol, endRow, endCol int) ChartData {
	data := ChartData{
		Labels: make([]string, 0),
		Values: make([]float64, 0),
	}

	// Extract labels from first column and values from second
	for row := startRow; row <= endRow && row < len(sheet.Rows); row++ {
		if startCol < len(sheet.Rows[row]) {
			label := sheet.Rows[row][startCol].Value
			if label == "" {
				label = fmt.Sprintf("Row %d", row+1)
			}
			data.Labels = append(data.Labels, label)

			// Try to get numeric value from next column
			if startCol+1 <= endCol && startCol+1 < len(sheet.Rows[row]) {
				valStr := sheet.Rows[row][startCol+1].Value
				if val, err := strconv.ParseFloat(valStr, 64); err == nil {
					data.Values = append(data.Values, val)
				} else {
					data.Values = append(data.Values, 0)
				}
			}
		}
	}

	return data
}

// RenderBarChart creates a beautiful ASCII bar chart
func RenderBarChart(data ChartData, style lipgloss.Style, accentColor, textColor lipgloss.Color) string {
	if len(data.Values) == 0 {
		return style.Render("No data to visualize")
	}

	var b strings.Builder

	// Title
	if data.Title != "" {
		title := lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Render(data.Title)
		b.WriteString(title + "\n\n")
	}

	maxVal := maxFloat(data.Values)
	if maxVal == 0 {
		maxVal = 1
	}

	maxLabelLen := 0
	for _, label := range data.Labels {
		if len(label) > maxLabelLen {
			maxLabelLen = len(label)
		}
	}
	if maxLabelLen > 15 {
		maxLabelLen = 15
	}

	barWidth := ChartWidth - maxLabelLen - 15

	for i, val := range data.Values {
		if i >= len(data.Labels) {
			break
		}

		// Label (truncated)
		label := data.Labels[i]
		if len(label) > maxLabelLen {
			label = label[:maxLabelLen-2] + ".."
		}
		labelStr := lipgloss.NewStyle().
			Foreground(textColor).
			Width(maxLabelLen).
			Render(label)

		// Bar
		barLen := int(float64(barWidth) * (val / maxVal))
		if barLen < 0 {
			barLen = 0
		}
		if barLen > barWidth {
			barLen = barWidth
		}

		bar := strings.Repeat("█", barLen)
		barStr := lipgloss.NewStyle().
			Foreground(accentColor).
			Render(bar)

		// Value
		valStr := lipgloss.NewStyle().
			Foreground(textColor).
			Render(fmt.Sprintf(" %.1f", val))

		b.WriteString(labelStr + " │ " + barStr + valStr + "\n")
	}

	// Scale
	b.WriteString(strings.Repeat(" ", maxLabelLen) + " └" + strings.Repeat("─", barWidth+2) + "\n")
	b.WriteString(strings.Repeat(" ", maxLabelLen) + "  0")
	b.WriteString(strings.Repeat(" ", barWidth-10))
	b.WriteString(fmt.Sprintf("%.1f\n", maxVal))

	return style.Render(b.String())
}

// RenderLineChart creates an ASCII line chart
func RenderLineChart(data ChartData, style lipgloss.Style, accentColor, textColor lipgloss.Color) string {
	if len(data.Values) == 0 {
		return style.Render("No data to visualize")
	}

	var b strings.Builder

	if data.Title != "" {
		title := lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Render(data.Title)
		b.WriteString(title + "\n\n")
	}

	height := ChartHeight
	width := ChartWidth

	minVal := minFloat(data.Values)
	maxVal := maxFloat(data.Values)
	if maxVal == minVal {
		maxVal = minVal + 1
	}

	// Create grid
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Plot points
	for i, val := range data.Values {
		x := int(float64(i) * float64(width-1) / float64(len(data.Values)-1))
		if x >= width {
			x = width - 1
		}

		normalized := (val - minVal) / (maxVal - minVal)
		y := height - 1 - int(normalized*float64(height-1))
		if y < 0 {
			y = 0
		}
		if y >= height {
			y = height - 1
		}

		grid[y][x] = '●'

		// Connect with line
		if i > 0 {
			prevVal := data.Values[i-1]
			prevX := int(float64(i-1) * float64(width-1) / float64(len(data.Values)-1))
			prevNormalized := (prevVal - minVal) / (maxVal - minVal)
			prevY := height - 1 - int(prevNormalized*float64(height-1))

			// Draw line between points
			if prevY < 0 {
				prevY = 0
			}
			if prevY >= height {
				prevY = height - 1
			}

			drawLine(grid, prevX, prevY, x, y)
		}
	}

	// Render grid with colors
	for i, row := range grid {
		lineStr := ""
		for _, char := range row {
			if char == '●' || char == '─' || char == '/' || char == '\\' || char == '|' {
				lineStr += lipgloss.NewStyle().Foreground(accentColor).Render(string(char))
			} else {
				lineStr += string(char)
			}
		}

		// Y-axis labels
		if i%3 == 0 {
			val := maxVal - (float64(i)/float64(height-1))*(maxVal-minVal)
			label := lipgloss.NewStyle().
				Foreground(textColor).
				Render(fmt.Sprintf("%6.1f│", val))
			b.WriteString(label + lineStr + "\n")
		} else {
			b.WriteString("      │" + lineStr + "\n")
		}
	}

	// X-axis
	b.WriteString("      └" + strings.Repeat("─", width) + "\n")

	return style.Render(b.String())
}

// RenderSparkline creates a compact sparkline
func RenderSparkline(data ChartData, accentColor lipgloss.Color) string {
	if len(data.Values) == 0 {
		return ""
	}

	chars := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	minVal := minFloat(data.Values)
	maxVal := maxFloat(data.Values)
	if maxVal == minVal {
		maxVal = minVal + 1
	}

	var b strings.Builder
	for _, val := range data.Values {
		normalized := (val - minVal) / (maxVal - minVal)
		idx := int(normalized * float64(len(chars)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(chars) {
			idx = len(chars) - 1
		}
		b.WriteRune(chars[idx])
	}

	return lipgloss.NewStyle().Foreground(accentColor).Render(b.String())
}

// RenderPieChart creates an ASCII pie chart
func RenderPieChart(data ChartData, style lipgloss.Style, colors []lipgloss.Color, textColor lipgloss.Color) string {
	if len(data.Values) == 0 {
		return style.Render("No data to visualize")
	}

	var b strings.Builder

	if data.Title != "" {
		title := lipgloss.NewStyle().
			Foreground(colors[0]).
			Bold(true).
			Render(data.Title)
		b.WriteString(title + "\n\n")
	}

	total := sumFloat(data.Values)
	if total == 0 {
		total = 1
	}

	radius := 10
	centerX, centerY := radius, radius

	// Create grid
	grid := make([][]int, radius*2+1)
	for i := range grid {
		grid[i] = make([]int, radius*2+1)
		for j := range grid[i] {
			grid[i][j] = -1
		}
	}

	// Fill pie chart
	currentAngle := 0.0
	for i, val := range data.Values {
		percentage := val / total
		angle := percentage * 360

		for y := 0; y <= radius*2; y++ {
			for x := 0; x <= radius*2; x++ {
				dx := x - centerX
				dy := y - centerY
				dist := math.Sqrt(float64(dx*dx + dy*dy))

				if dist <= float64(radius) {
					pointAngle := math.Atan2(float64(dy), float64(dx)) * 180 / math.Pi
					if pointAngle < 0 {
						pointAngle += 360
					}

					if pointAngle >= currentAngle && pointAngle < currentAngle+angle {
						grid[y][x] = i
					}
				}
			}
		}

		currentAngle += angle
	}

	// Render grid
	for _, row := range grid {
		for _, val := range row {
			if val >= 0 && val < len(colors) {
				b.WriteString(lipgloss.NewStyle().
					Foreground(colors[val%len(colors)]).
					Render("●"))
			} else {
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
	}

	// Legend
	b.WriteString("\n")
	for i, label := range data.Labels {
		if i >= len(data.Values) || i >= len(colors) {
			break
		}
		percentage := (data.Values[i] / total) * 100

		colorBox := lipgloss.NewStyle().
			Foreground(colors[i%len(colors)]).
			Render("■")

		labelStr := lipgloss.NewStyle().
			Foreground(textColor).
			Render(fmt.Sprintf(" %s: %.1f%%", label, percentage))

		b.WriteString(colorBox + labelStr + "\n")
	}

	return style.Render(b.String())
}

// Helper functions
func drawLine(grid [][]rune, x0, y0, x1, y1 int) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	for {
		if x0 >= 0 && x0 < len(grid[0]) && y0 >= 0 && y0 < len(grid) {
			if grid[y0][x0] == ' ' {
				if dx > dy {
					grid[y0][x0] = '─'
				} else if dy > dx {
					grid[y0][x0] = '|'
				} else if (sx > 0 && sy > 0) || (sx < 0 && sy < 0) {
					grid[y0][x0] = '\\'
				} else {
					grid[y0][x0] = '/'
				}
			}
		}

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func maxFloat(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	max := vals[0]
	for _, v := range vals {
		if v > max {
			max = v
		}
	}
	return max
}

func minFloat(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	min := vals[0]
	for _, v := range vals {
		if v < min {
			min = v
		}
	}
	return min
}

func sumFloat(vals []float64) float64 {
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
